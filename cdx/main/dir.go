package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ttchengcheng/file/console"
)

func dataPath() string {
	pwd, err := os.Executable()
	if err != nil {
		panic(err)
	}
	dir, _ := filepath.Split(pwd)
	return filepath.Join(dir, "dir-list.json")
}

func save(shortcuts *map[string]string) error {
	data, err := json.Marshal(shortcuts)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dataPath(), data, 0600)
}

func load() map[string]string {
	v := map[string]string{}

	jsonFile := dataPath()
	fi, err := os.Stat(jsonFile)
	if os.IsNotExist(err) {
		return v
	}

	if fi.IsDir() {
		return v
	}
	data, err := ioutil.ReadFile(jsonFile)

	if err = json.Unmarshal(data, &v); err != nil {
		v = map[string]string{}
	}
	return v
}

func set(shortcut string, path string) error {
	shortcuts := load()
	shortcuts[shortcut] = path
	return save(&shortcuts)
}

func unset(shortcut string) error {
	shortcuts := load()
	delete(shortcuts, shortcut)
	return save(&shortcuts)
}

func add(shortcut string, path string) {
	shortcuts := load()
	if _, ok := shortcuts[shortcut]; ok {
		response := console.Prompt("Data already exist as above, replace?(y/n)")
		if strings.Compare(strings.ToUpper(response), "Y") != 0 {
			return
		}
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		console.Error(&err)
		return
	}
	err = set(shortcut, absPath)
	if err != nil {
		console.Error(&err)
		return
	}

	console.Outputf("1 entry added:\n%s | %s\n", shortcut, absPath)
}
func remove(shortcut string) {
	shortcuts := load()
	path, ok := shortcuts[shortcut]
	if !ok {
		console.Outputf("entry [%s] not found.", shortcut)
		return
	}

	err := unset(shortcut)
	if err != nil {
		console.Error(&err)
		return
	}

	console.Outputf("1 entry removed:\n%s | %s\n", shortcut, path)
}

func list() {
	shortcuts := load()
	for shortcut, path := range shortcuts {
		console.Outputf("%s | %s\n", shortcut, path)
	}
}

func lookup(shortcut string) (string, bool) {
	shortcuts := load()
	path, ok := shortcuts[shortcut]
	return path, ok
}

func showUsage() {

}
func main() {
	wrongUsage := false
	defer func() {
		if wrongUsage {
			showUsage()
			return
		}
	}()

	args := os.Args[1:]
	if len(args) < 1 {
		wrongUsage = true
		return
	}

	for i, count := 0, len(args); i < count; i++ {
		arg := args[i]
		if len(arg) < 1 {
			continue
		}

		if arg[0] == '-' {
			switch arg {
			case "-a", "--add", "--append":
				if i >= count-2 {
					wrongUsage = true
				} else {
					add(args[i+1], args[i+2])
				}
			case "-d", "-r", "--delete", "--remove":
				if i >= count-1 {
					wrongUsage = true
				} else {
					remove(args[i+1])
				}
			case "-l", "--list":
				list()
			case "-h", "--help":
				showUsage()
			}
		} else {
			lookup(arg)
		}
	}
}
