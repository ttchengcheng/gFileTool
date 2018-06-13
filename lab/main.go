package main

import (
	"fmt"
	"os"
)

func main() {
	fi, _ := os.Stat("/Users/hoolai/Documents/build/hmotion/debug/resources")
	a := int(os.ModeSymlink)
	b := int(fi.Mode() & os.ModeSymlink)
	if a == b {
		fmt.Println("is link")
	}
	fmt.Println(fi.Name())
	fmt.Println(fi.IsDir())
}
