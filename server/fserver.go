package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
)

var defaultPort = "8801"

func showHelp(msg string) {
	if len(msg) > 0 {
		fmt.Println("fserver: ", msg)
		fmt.Println("--------------------------------------------------")
		fmt.Println()
	}

	fmt.Println("Usage  : fserver <directory> [<port>]")
	fmt.Println("Example: fserver /Users/hoolai/Documents/Book 7070")
	fmt.Println()
	fmt.Println("default of <port> is ", defaultPort)
	os.Exit(0)
}

func main() {
	args := os.Args[1:]

	// 1st argument
	if len(args) < 1 {
		showHelp("The 1st argument can't be omitted")
	}

	dir := args[0]
	if dir == "-h" || dir == "-help" ||
		dir == "--h" || dir == "--help" {
		showHelp("")
	}

	if fi, err := os.Stat(dir); err != nil || !fi.IsDir() {
		showHelp("The 1st argument should be a valid directory")
	}

	// 2nd argument
	port := defaultPort
	if len(args) > 1 {
		port = args[1]
		re := regexp.MustCompile("[0-9]+")
		if !re.MatchString(port) {
			showHelp("The 2nd argument should be a valid number")
		}
	}

	// start server
	panic(http.ListenAndServe(":"+port, http.FileServer(http.Dir(dir))))
}
