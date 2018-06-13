package main

import "net/http"

func main() {
	panic(http.ListenAndServe(":8801", http.FileServer(http.Dir("/Users/hoolai/Downloads"))))
}