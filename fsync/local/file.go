package local

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ttchengcheng/file/fsync"
)

// File of local disk
type File struct {
	Path          string
	IgnoreSetting *fsync.IgnoreSetting
}

// Checksum is a function to calculate sha1 with path of the file
func (*File) Checksum(path string) (string, error) {
	const buffSize = 1024
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buff := make([]byte, buffSize)
	sh1 := sha1.New()
	for {
		readLen, err := file.Read(buff)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		sh1.Write(buff[:readLen])
	}
	return fmt.Sprintf("%x", sh1.Sum(nil)), nil
}

// Dir returns the path string
func (f *File) Dir() string {
	if f == nil {
		return ""
	}
	return f.Path
}

// Dir returns the path string
func (f *File) GetList(list *fsync.FileList) error {
	if f == nil {
		return nil
	}
	return filepath.Walk(f.Path, func(path string, fileInfo os.FileInfo, err error) error {
		if f.IgnoreSetting != nil && f.IgnoreSetting.IsIgnored(path, fileInfo) {
			if fileInfo.IsDir() {
				return filepath.SkipDir
			} else {
				return nil
			}
		}
		cleanPath := filepath.ToSlash(path[len(f.Path):])
		(*list)[cleanPath] = fileInfo
		return nil
	})
}

func (f *File) Copy(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	dir, _ := filepath.Split(dst)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := f.Mkdir(dir, 0777); err != nil {
			return err
		}
	}

	out, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func (f *File) Remove(path string) error {
	if f == nil {
		return nil
	}
	return os.RemoveAll(path)
}

func (f *File) Mkdir(path string, mode os.FileMode) error {
	return os.MkdirAll(path, mode)
}
