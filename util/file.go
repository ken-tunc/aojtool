package util

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

var CacheDir = filepath.Join(HomeDir(), ".cache", "aojtool")

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

func EnsurePath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
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

	return absPath, nil
}

func HomeDir() string {
	dir, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	return dir
}

func ReadFile(path string) (string, error) {
	byteContent, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(byteContent), nil
}

func RemoveContents(dirname string) error {
	dir, err := ioutil.ReadDir(dirname)
	if err != nil {
		return err
	}

	for _, subDir := range dir {
		err = os.RemoveAll(filepath.Join(dirname, subDir.Name()))
		if err != nil {
			return err
		}
	}

	return nil
}
