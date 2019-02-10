package util

import (
	"log"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}

func EnsurePath(path string) (absPath string, err error) {
	absPath, err = filepath.Abs(path)
	if err != nil {
		return "", err
	}

	dir := filepath.Dir(absPath)
	exist, err := Exists(dir)
	if err != nil {
		return "", err
	}

	if !exist {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return "", err
		}
	}

	return
}

func HomeDir() string {
	dir, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	return dir
}
