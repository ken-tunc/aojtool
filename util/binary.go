package util

import (
	"bytes"
	"encoding/gob"
)

func Serialize(data interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(data); err != nil {
		return make([]byte, 0), err
	}

	return buf.Bytes(), nil
}

func Deserialize(byteData []byte, v interface{}) error {
	buf := bytes.NewBuffer(byteData)
	if err := gob.NewDecoder(buf).Decode(v); err != nil {
		return err
	}

	return nil
}
