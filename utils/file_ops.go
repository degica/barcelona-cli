package utils

import (
	"io/ioutil"
	"os"
)

type FileOps struct {
}

func (_ FileOps) FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (_ FileOps) ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func (_ FileOps) WriteFile(path string, content []byte) error {
	return ioutil.WriteFile(path, content, 0600)
}
