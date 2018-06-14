package fsync

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Sync struct {
	Source, Destination Filer
	KeepOther           bool
}

type taskFn func() error

func task(fn taskFn, format string, a ...interface{}) {
	fmt.Printf(format+"...", a...)
	if err := fn(); err != nil {
		fmt.Println("failed.")
		panic(err)
	} else {
		fmt.Println("done.")
	}
}

func (sc *Sync) checkFolder() {
	if sc == nil {
		panic(errors.New("sync object is nil"))
	}
	// source folder should exist and be a folder
	srcInfoRoot, err := os.Stat(sc.Source.Dir())
	task(func() error {
		if err != nil {
			return err
		}
		if !srcInfoRoot.IsDir() {
			return errors.New("source path is not a valid folder")
		}
		return nil
	}, "check source path")

	// destination folder should be a folder and exist
	desInfoRoot, err := os.Stat(sc.Destination.Dir())
	if os.IsNotExist(err) {
		task(func() error {
			return sc.Destination.Mkdir(sc.Destination.Dir(), srcInfoRoot.Mode().Perm())
		}, "creating [%s]", sc.Destination.Dir())
	} else {
		if !desInfoRoot.IsDir() {
			task(func() error {
				return sc.Destination.Remove(sc.Destination.Dir())
			}, "removing [%s]", sc.Destination.Dir())
			task(func() error {
				return sc.Destination.Mkdir(sc.Destination.Dir(), srcInfoRoot.Mode().Perm())
			}, "creating [%s]", sc.Destination.Dir())
		}
	}
}

func (sc *Sync) readFileList() (listSource FileList, listDestination FileList) {
	if sc == nil {
		panic(errors.New("sync object is nil"))
	}

	task(func() error {
		listSource = make(FileList)
		return sc.Source.GetList(&listSource)
	}, "reading list [%s]", sc.Source.Dir())

	task(func() error {
		listDestination = make(FileList)
		return sc.Destination.GetList(&listDestination)
	}, "reading list [%s]", sc.Destination.Dir())

	return
}

func (sc *Sync) Run() {
	if sc == nil {
		panic(errors.New("sync object is nil"))
	}

	// recover-panic
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println()
			fmt.Println("all done!")
		}
	}()

	if sc.Source.Dir() == sc.Destination.Dir() {
		return
	}

	sc.checkFolder()
	listSource, listDestination := sc.readFileList()

	indexPath := filepath.Join(sc.Destination.Dir(), ".sync-cache")
	indexes := Indexer{indexPath, sc.Source, nil}
	task(func() error {
		return indexes.Load()
	}, "loading indexes [%s]", indexPath)

	defer task(func() error {
		return indexes.Save()
	}, "writing indexes [%s]", indexPath)

	removeDirs := make([]string, 0)
	for path, srcInfo := range listSource {
		srcPath := filepath.Join(sc.Source.Dir(), path)
		desPath := filepath.Join(sc.Destination.Dir(), path)

		if desInfo, ok := listDestination[path]; !ok {
			if srcInfo.IsDir() {
				task(func() error {
					return sc.Destination.Mkdir(desPath, srcInfo.Mode().Perm())
				}, "creating [%s]", desPath)
			} else {
				task(func() error {
					return sc.Destination.Copy(srcPath, desPath, srcInfo.Mode().Perm())
				}, "copying [%s]", desPath)
			}
		} else {
			if desInfo.IsDir() {
				if !srcInfo.IsDir() {
					task(func() error {
						return sc.Destination.Remove(desPath)
					}, "removing [%s]", desPath)

					removeDirs = append(removeDirs, path+"/")
				}
			} else {
				if indexes.Checksum(srcPath, srcInfo) != indexes.Checksum(desPath, desInfo) {
					task(func() error {
						return sc.Destination.Copy(srcPath, desPath, srcInfo.Mode().Perm())
					}, "copying [%s]", desPath)
				}
			}
			delete(listDestination, path)
		}
	}

	if !sc.KeepOther {
		for path := range listDestination {
			desPath := filepath.Join(sc.Destination.Dir(), path)

			task(func() error {
				return sc.Destination.Remove(desPath)
			}, "removing [%s]", desPath)
		}
	}
}
