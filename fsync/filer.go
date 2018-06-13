package fsync

import "os"

type FileList map[string]os.FileInfo

type Filer interface {
	Dir() string
	GetList(list *FileList) error
	Copy(src, dst string, mode os.FileMode) error
	Remove(path string) error
	Mkdir(path string, mode os.FileMode) error
	Checksum(path string) (string, error)
}
