package io

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

// PathExists check if the path exists
func PathExists(p string) bool {
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		return false
	}

	if err != nil {
		panic(err)
	}

	return true
}

// FileListBySuffix Recursively retrieve a list of all files with a specified extension in a directory.
func FileListBySuffix(p, ends string) (fileList []string) {
	err := filepath.Walk(p, func(p string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		if ends == "" || ends == path.Ext(p) {
			fileList = append(fileList, p)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return
}

// FileListByPattern Recursively retrieve a list of all files matching a specified pattern in a directory.
func FileListByPattern(p, pattern string) (fileList []string) {
	err := filepath.Walk(p, func(p string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		if strings.Count(path.Base(p), pattern) != 0 {
			fileList = append(fileList, p)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return
}

// DirListByPath Recursively retrieve a list of full paths of all subdirectories in a specified directory.
func DirListByPath(path string) (err error, dirList []string) {
	err = filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		dirList = append(dirList, path)
		return nil
	})

	return
}
