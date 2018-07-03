package console

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

// Prompt is a
func Prompt(title string) (input string) {
	fmt.Print(title + " ")
	fmt.Scanln(&input)

	return
}

func Exec(path string, args ...string) error {
	// out, _ := exec.Command(path, args...).Output()
	// fmt.Printf("%s\n", out)

	path, err := exec.LookPath(path)
	if err != nil {
		return err
	}

	p, err := os.StartProcess(path, args, &os.ProcAttr{})
	if err != nil {
		return err
	}
	_, err = p.Wait()

	return err
	// cmd := exec.Command(path, args...)
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// cmd.Run()
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
