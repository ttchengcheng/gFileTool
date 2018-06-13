/*
TODO:
. improve output
. file index/checksum
. test case
*/
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"errors"

	"github.com/ttchengcheng/file/fsync"
	"github.com/ttchengcheng/file/fsync/cache"
	"github.com/ttchengcheng/file/fsync/local"
)

type taskFn func() error

func group(name string) {
	fmt.Printf("\n%s\n\n", name)
}

func task(fn taskFn, format string, a ...interface{}) {
	const indent = "-"
	fmt.Printf(indent+format, a...)
	if err := fn(); err != nil {
		fmt.Println("failed.")
		panic(err)
	} else {
		fmt.Println("done.")
	}
}

func checkFolder(src fsync.Filer, des fsync.Filer) {
	// source folder should exist and be a folder
	srcInfoRoot, err := os.Stat(src.Dir())
	if err != nil {
		panic(err)
	}
	if !srcInfoRoot.IsDir() {
		panic(errors.New("source path is not a valid folder"))
	}

	// destination folder should be a folder and exist
	desInfoRoot, err := os.Stat(des.Dir())
	if os.IsNotExist(err) {
		des.Mkdir(des.Dir(), srcInfoRoot.Mode().Perm())
	}
	if !desInfoRoot.IsDir() {
		des.Remove(des.Dir())
		des.Mkdir(des.Dir(), srcInfoRoot.Mode().Perm())
	}
}

func readFileList(src fsync.Filer, des fsync.Filer) (listSource fsync.FileList, listDestination fsync.FileList) {
	task(func() error {
		listSource = make(fsync.FileList)
		return src.GetList(&listSource)
	}, "[%s]...", src.Dir())

	task(func() error {
		listDestination = make(fsync.FileList)
		return des.GetList(&listDestination)
	}, "[%s]...", des.Dir())

	return
}

func sync(src fsync.Filer, des fsync.Filer) {
	// recover-panic
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("\nALL DONE!")
		}
	}()

	group("check source/destination folders")
	checkFolder(src, des)

	group("Get file list")
	listSource, listDestination := readFileList(src, des)

	group("Sync files")
	indexes := cache.Indexer{}
	indexPath := filepath.Join(des.Dir(), ".sync-cache")
	task(func() error {
		return indexes.Load(indexPath)
	}, "loading indexes [%s]...", indexPath)

	defer task(func() error {
		return indexes.Save()
	}, "writing index file [%s]...", indexPath)

	removeDirs := make([]string, 0)
	for path, srcInfo := range listSource {
		srcPath := filepath.Join(src.Dir(), path)
		desPath := filepath.Join(des.Dir(), path)

		if desInfo, ok := listDestination[path]; !ok {
			if srcInfo.IsDir() {
				task(func() error {
					return des.Mkdir(desPath, srcInfo.Mode().Perm())
				}, "creating [%s]...", desPath)
			} else {
				task(func() error {
					return des.Copy(srcPath, desPath, srcInfo.Mode().Perm())
				}, "copying [%s]...", desPath)
			}
		} else {
			if desInfo.IsDir() {
				if !srcInfo.IsDir() {
					task(func() error {
						return des.Remove(desPath)
					}, "removing [%s]...", desPath)

					removeDirs = append(removeDirs, path+"/")
				}
			} else {
				if indexes.Checksum(srcPath, srcInfo) != indexes.Checksum(desPath, desInfo) {
					task(func() error {
						return des.Copy(srcPath, desPath, srcInfo.Mode().Perm())
					}, "copying [%s]...", desPath)
				}
			}
			delete(listDestination, path)
		}
	}

	group("remove unmatched files")
	for path := range listDestination {
		desPath := filepath.Join(des.Dir(), path)

		task(func() error {
			return des.Remove(desPath)
		}, "removing [%s]...", desPath)
	}
}

func main() {
	srcDir := "/Users/hoolai/git/hScene"
	desDir := "/Users/hoolai/Downloads/temp"
	filter := &fsync.SkipSetting{}
	filter.Parse(".git/\n.sync-cache")

	src := local.File{srcDir, filter}
	des := local.File{desDir, filter}
	sync(&src, &des)
}
