package fsync

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

// Indexer
type Indexer struct {
	path string
	file Filer
	data map[string][]string
}

const (
	fieldPath = iota
	fieldSize
	fieldModTime
	fieldChecksum
)

// Load
func (indexer *Indexer) Load() error {
	if indexer == nil {
		return errors.New("indexer is nil")
	}
	if indexer.data == nil {
		indexer.data = make(map[string][]string)
	}

	file, err := os.Open(indexer.path)
	if err != nil {
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ",")
		if len(fields) <= fieldChecksum {
			continue
		}
		indexer.data[fields[fieldPath]] = fields[:fieldChecksum+1]
	}

	if err = scanner.Err(); err != nil {
		return err
	}
	return nil
}

// Save
func (indexer *Indexer) Save() (err error) {
	if indexer == nil {
		return errors.New("indexer is nil")
	}
	file, err := os.Open(indexer.path)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
		file, err = os.Create(indexer.path)
		if err != nil {
			return
		}
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, list := range indexer.data {
		writer.WriteString(strings.Join(list, ",") + "\n")
	}

	return
}

// Checksum: get the checksum for given file.
// return cached checksum if os.FileInfo.Size and os.FileInfo.ModTime matches
// or recalculate checksum and return
func (indexer *Indexer) Checksum(path string, info os.FileInfo) (checksum string) {
	if indexer == nil {
		return
	}
	checksum = ""
	if info.IsDir() {
		return
	}

	if indexer.data == nil {
		indexer.data = make(map[string][]string)
	}
	size := strconv.FormatInt(info.Size(), 10)
	modTime := info.ModTime().String()
	index, ok := indexer.data[path]
	if !ok || len(index) <= fieldChecksum || index[fieldSize] != size || index[fieldModTime] != modTime {
		checksum, _ = indexer.file.Checksum(path)
		indexer.data[path] = []string{path, size, modTime, checksum}
		return
	}

	checksum = index[fieldChecksum]
	return
}
