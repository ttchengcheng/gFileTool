package console

import (
	"path/filepath"
	"testing"
)

func TestAppendFile(t *testing.T) {
	src, _ := filepath.Abs("helper_test.go")
	des, _ := filepath.Abs("new.txt")
	AppendFile(src, des)

	Exec("rm", "new.txt")
}

func TestExec(t *testing.T) {
	Exec("ls", "-l", "-a")
}
