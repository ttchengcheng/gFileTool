package fsync

import (
	"os"
	"strings"
)

// IgnoreSetting is a struct who hold the information of the folder/files that should be ignored
type IgnoreSetting struct {
	ignoreFolders    map[string]struct{}
	ignoreFiles      map[string]struct{}
	folderExceptions map[string]struct{}
	fileExceptions   map[string]struct{}
}

// Parse a IgnoreSetting string
func (ss *IgnoreSetting) Parse(str string) error {
	if ss == nil {
		panic("nil IgnoreSetting in Parse()")
	}
	ss.ignoreFolders = map[string]struct{}{}
	ss.ignoreFiles = map[string]struct{}{}
	ss.folderExceptions = map[string]struct{}{}
	ss.fileExceptions = map[string]struct{}{}

	lines := strings.Split(str, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		size := len(line)
		reversed := false

		// ignore empty string
		if size == 0 {
			continue
		}
		switch line[0] {
		// ignore comment that start with "//"
		case '/':
			if size > 1 && line[1] == '/' {
				continue
			}
		// start with "!" means "this is an exception, don't ignore this"
		case '!':
			reversed = true
		}

		if line[size-1] == '/' { // folder
			if reversed {
				ss.folderExceptions[line] = struct{}{}
			} else {
				ss.ignoreFolders[line[:size-1]] = struct{}{}
			}
		} else { // file
			if reversed {
				ss.fileExceptions[line] = struct{}{}
			} else {
				ss.ignoreFiles[line] = struct{}{}
			}
		}
	}
	return nil
}

// IsIgnored is a function to check whether a path should be ignored
func (ss *IgnoreSetting) IsIgnored(path string, fileInfo os.FileInfo) bool {
	if fileInfo.IsDir() {
		if _, ok := ss.folderExceptions[fileInfo.Name()]; ok {
			return false
		}
		if _, ok := ss.ignoreFolders[fileInfo.Name()]; ok {
			return true
		}
	} else {
		if _, ok := ss.fileExceptions[fileInfo.Name()]; ok {
			return false
		}
		if _, ok := ss.ignoreFiles[fileInfo.Name()]; ok {
			return true
		}
	}
	return false
}
