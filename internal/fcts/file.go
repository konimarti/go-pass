package fcts

import (
	"os"
)

func IsDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	if fileInfo.IsDir() {
		return true, nil
	}
	return false, nil
}

func IsFile(path string) (bool, error) {
	ok, err := IsDir(path)
	return !ok, err
}

func FileExists(path string) bool {
	ok, err := IsFile(path)
	return ok && err == nil
}

func CheckSneakyPath(path string) string {
	// FIXME: implement checks
	return path
}
