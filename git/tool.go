package git

import (
	"os"
	"path/filepath"

	"strings"

	"fmt"

	"github.com/ttchengcheng/file/console"
)

// IgnoreFile ...
func IgnoreFile(key string) (string, bool) {
	value, ok := map[string]string{
		"vsc":  "VisualStudioCode",
		"qt":   "Qt",
		"py":   "Python",
		"objc": "Objective-C",
		"go":   "Go",
		"node": "Node",
		"rust": "Rust",
	}[key]

	if !ok {
		return "", false
	}
	return filepath.Join("/Users/hoolai/git/gitignore", value+".gitignore"), ok
}

// GInit is a
func GInit(templates ...string) {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	gitDir := filepath.Join(pwd, ".git")
	gitIgnoreFile := filepath.Join(pwd, ".gitignore")

	fmt.Println("file is : ", gitIgnoreFile)

	exist := func(path string) bool {
		fmt.Println("check existing [", path, "]")
		_, err = os.Stat(path)
		return !os.IsNotExist(err)
	}

	// clean up git settings
	if exist(gitDir) || exist(gitIgnoreFile) {
		input := console.Prompt("git files exist, continue? (y/n)")
		if strings.Compare(strings.ToUpper(input), "Y") != 0 {
			os.Exit(0)
		}

		GClean()
	}

	// git init
	console.Exec("git", "init")

	// copy contents to .gitignore
	for _, template := range templates {
		path, ok := IgnoreFile(template)
		if !ok {
			continue
		}

		console.AppendFile(path, gitIgnoreFile)
	}

	// vi .gitignore
	console.Exec("vi", ".gitignore")

	// git add .
	console.Exec("git", "add", ".")

	// git commit -m "init"
	console.Exec("git", "commit", "-m", `"init"`)
}

// GClean is a
func GClean() {
	console.Exec("rm", "-fr", ".git")
	console.Exec("rm", "-fr", ".gitignore")
}
