package main

import (
	"os"

	"github.com/ttchengcheng/file/git"
)

func main() {
	git.GInit(os.Args[1:]...)
}
