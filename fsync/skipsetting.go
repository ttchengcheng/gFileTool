package fsync

import (
	"os"
	"strings"
)

// SkipSetting is a struct who hold the information of the folder/files that should be skipped
type SkipSetting struct {
	skipFolders      map[string]struct{}
	skipFiles        map[string]struct{}
	folderExceptions map[string]struct{}
	fileExceptions   map[string]struct{}
}

// Parse a skip setting string
func (ss *SkipSetting) Parse(str string) error {
	if ss == nil {
		panic("nil SkipSetting in Parse()")
	}
	ss.skipFolders = map[string]struct{}{}
	ss.skipFiles = map[string]struct{}{}
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
		// start with "!" means "this is an exception, don't skip this"
		case '!':
			reversed = true
		}

		if line[size-1] == '/' { // folder
			if reversed {
				ss.folderExceptions[line] = struct{}{}
			} else {
				ss.skipFolders[line[:size-1]] = struct{}{}
			}
		} else { // file
			if reversed {
				ss.fileExceptions[line] = struct{}{}
			} else {
				ss.skipFiles[line] = struct{}{}
			}
		}
	}
	return nil
}

// IsSkipped is a function to check whether a path should be skipped
func (ss *SkipSetting) IsSkipped(path string, fileInfo os.FileInfo) bool {
	if fileInfo.IsDir() {
		if _, ok := ss.folderExceptions[fileInfo.Name()]; ok {
			return false
		}
		if _, ok := ss.skipFolders[fileInfo.Name()]; ok {
			return true
		}
	} else {
		if _, ok := ss.fileExceptions[fileInfo.Name()]; ok {
			return false
		}
		if _, ok := ss.skipFiles[fileInfo.Name()]; ok {
			return true
		}
	}
	return false
}
