package console

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

// Prompt is a
func Prompt(title string) (input string) {
	fmt.Print(title + " ")
	fmt.Scanln(&input)

	return
}

func Shell(path string, args ...string) error {
	return Exec("sh", "-c", path+" "+strings.Join(args, " "))
}

func Exec(path string, args ...string) error {
	cmd := exec.Command(path, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func AppendFile(src string, des string) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	var f *os.File
	_, err = os.Stat(des)
	if os.IsNotExist(err) {
		f, err = os.Create(des)
	} else {
		f, err = os.OpenFile(des, os.O_APPEND|os.O_WRONLY, 0600)
	}

	if err != nil {
		return err
	}

	defer f.Close()
	f.Write(data)

	return nil
}

func Outputf(format string, a ...interface{}) {
	color.Cyan(format, a...)
}

func Errorf(format string, a ...interface{}) {
	color.Red(format, a...)
}

func Error(err *error) {
	if err == nil {
		return
	}
	Errorf("%s\n", (*err).Error())
}
