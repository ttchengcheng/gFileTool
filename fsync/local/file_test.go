package local

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ttchengcheng/file/fsync"
)

func curWorkingPath() string {
	_, currentFilePath, _, _ := runtime.Caller(0)
	return path.Dir(currentFilePath)
}

func TestFile_Checksum(t *testing.T) {
	f := File{}
	cases := []struct {
		path, checksum string
		hasError       bool
	}{
		{"testdata/a.txt", "fff8f4ad15963784e898d2f76987c87908755def", false},
		{"testdata/big.binary", "0e2aa6b139224b64b41dc9ad6c3c7f124f45b1c1", false},
		{"testdata/empty", "da39a3ee5e6b4b0d3255bfef95601890afd80709", false},
		{"testdata/not-exist", "", true},
	}

	for _, c := range cases {
		// path := filepath.Join(folder, c.path)
		checksum, err := f.Checksum(c.path)

		switch {
		case c.hasError && err == nil:
			t.Error("want an error, but no error happens")
		case !c.hasError && err != nil:
			t.Errorf("want no error, but error [%s] happens", err.Error())
		case checksum != c.checksum:
			t.Errorf("checksum of [%s] is [%s], want [%s]", c.path, checksum, c.checksum)
		}
	}
}

func TestFile_Dir(t *testing.T) {
	ss := &fsync.IgnoreSetting{}
	cases := []struct {
		file *File
		dir  string
	}{
		{&File{}, ""},
		{&File{"", nil}, ""},
		{&File{"", ss}, ""},
		{&File{"*&^", nil}, "*&^"},
		{&File{"aaa/bbb", nil}, "aaa/bbb"},
		{&File{"aaa/bbb", ss}, "aaa/bbb"},
		{nil, ""},
	}

	for _, c := range cases {
		// path := filepath.Join(folder, c.path)
		dir := c.file.Dir()

		switch {
		case dir != c.dir:
			t.Errorf("Dir() of [%v] is [%s], want [%s]", c.file, dir, c.dir)
		}
	}
}

func TestFile_Copy(t *testing.T) {
	f := File{}

	curPath := curWorkingPath()
	cases := []struct {
		src, des  string
		hasError  bool
		cleanList []string
	}{
		{"testdata/a.txt", "test_files/aa.txt", false, []string{}},
		{"testdata/a.txt", "test_files/aa/aa.txt", false, []string{"test_files/aa", "test_files/aa/aa.txt"}},
		{"testdata/a.txt", "test_files/中文/aa.txt", false, []string{"test_files/中文", "test_files/中文/aa.txt"}},
		{"testdata/b.txt", "test_files/中文/aa.txt", true, []string{}},
	}

	for _, c := range cases {
		srcPath := filepath.Join(curPath, c.src)
		desPath := filepath.Join(curPath, c.des)

		err := f.Copy(srcPath, desPath, 0666)

		switch {
		case c.hasError && err == nil:
			t.Error("want an error, but no error happens")
		case !c.hasError && err != nil:
			t.Errorf("want no error, but error [%s] happens", err.Error())
		}

		if err == nil {
			fileInfo, statErr := os.Stat(desPath)
			if os.IsNotExist(statErr) {
				t.Errorf("Copy [%s] to [%s], no error returned, but file is not copied", srcPath, desPath)
			}
			if fileInfo.IsDir() {
				t.Errorf("Copy [%s] to [%s], no error returned, but result is a directory, not a file", srcPath, desPath)
			}

			// clean up
			for _, fpath := range c.cleanList {
				f.Remove(fpath)
			}
		}
	}
}

func TestFile_GetList(t *testing.T) {
	type fields struct {
		Path          string
		IgnoreSetting *fsync.IgnoreSetting
	}
	type args struct {
		list *fsync.FileList
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				Path:          tt.fields.Path,
				IgnoreSetting: tt.fields.IgnoreSetting,
			}
			if err := f.GetList(tt.args.list); (err != nil) != tt.wantErr {
				t.Errorf("File.GetList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
