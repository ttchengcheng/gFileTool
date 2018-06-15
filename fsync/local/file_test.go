package local

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ttchengcheng/file/fsync"
)

func TestFile_Checksum(t *testing.T) {
	f := File{}
	cases := []struct {
		path, checksum string
		wantErr        bool
	}{
		{"testdata/a.txt", "fff8f4ad15963784e898d2f76987c87908755def", false},
		{"testdata/big.binary", "0e2aa6b139224b64b41dc9ad6c3c7f124f45b1c1", false},
		{"testdata/empty", "da39a3ee5e6b4b0d3255bfef95601890afd80709", false},
		{"testdata/not-exist", "", true},
	}

	for _, c := range cases {
		checksum, err := f.Checksum(c.path)

		switch {
		case c.wantErr && err == nil:
			t.Error("want an error, but no error happens")
		case !c.wantErr && err != nil:
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
		dir := c.file.Dir()

		switch {
		case dir != c.dir:
			t.Errorf("Dir() of [%v] is [%s], want [%s]", c.file, dir, c.dir)
		}
	}
}

func TestFile_Copy(t *testing.T) {
	f := File{}

	cases := []struct {
		src, des  string
		wantErr   bool
		cleanList []string
	}{
		{"testdata/a.txt", "test_files/aa.txt", false, []string{}},
		{"testdata/a.txt", "test_files/aa/aa.txt", false, []string{"test_files/aa", "test_files/aa/aa.txt"}},
		{"testdata/a.txt", "test_files/中文/aa.txt", false, []string{"test_files/中文", "test_files/中文/aa.txt"}},
		{"testdata/b.txt", "test_files/中文/aa.txt", true, []string{}},
	}

	for _, c := range cases {
		srcPath, _ := filepath.Abs(c.src)
		desPath, _ := filepath.Abs(c.des)

		err := f.Copy(srcPath, desPath, 0666)

		switch {
		case c.wantErr && err == nil:
			t.Error("want an error, but no error happens")
		case !c.wantErr && err != nil:
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

	is := fsync.IgnoreSetting{}
	is.Parse("empty\nignore/")

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "", fields: fields{Path: "testdata", IgnoreSetting: &is}, wantErr: false},
		{name: "", fields: fields{Path: "aa", IgnoreSetting: &is}, wantErr: false},
		{name: "", fields: fields{Path: "aa", IgnoreSetting: &is}, wantErr: false},
		{name: "", fields: fields{Path: "aa", IgnoreSetting: &is}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				Path:          tt.fields.Path,
				IgnoreSetting: tt.fields.IgnoreSetting,
			}
			var list = fsync.FileList{}
			if err := f.GetList(&list); (err != nil) != tt.wantErr {
				t.Errorf("File.GetList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
