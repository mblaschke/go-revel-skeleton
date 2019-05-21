package app

import (
	"os"
)

func IsDirectory(path string) (bool) {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

func IsRegularFile(path string) (bool) {
	fileInfo, _ := os.Stat(path)
	return fileInfo.Mode().IsRegular()
}
