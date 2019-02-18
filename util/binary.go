package util

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"os"
)

func SaveData(path string, data interface{}) error {
	byteData, err := serialize(data)
	if err != nil {
		return err
	}

	return writeBytes(path, byteData)
}

func LoadData(path string, data interface{}) error {
	byteData, err := readBytes(path)
	if err != nil {
		return err
	}

	return deserialize(byteData, data)
}

func RemoveData(path string) error {
	return os.RemoveAll(path)
}

func serialize(data interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(data); err != nil {
		return make([]byte, 0), err
	}

	return buf.Bytes(), nil
}

func deserialize(byteData []byte, v interface{}) error {
	buf := bytes.NewBuffer(byteData)
	if err := gob.NewDecoder(buf).Decode(v); err != nil {
		return err
	}

	return nil
}

func writeBytes(path string, data []byte) error {
	absPath, err := EnsurePath(path)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(absPath, data, 0600)
}

func readBytes(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}
