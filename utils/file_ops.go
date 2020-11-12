package utils

import (
	"os"
)

type FileOps struct {
}

func (_ FileOps) FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
