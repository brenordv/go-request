package utils

import (
	"os"
	"path"
)


func GetAppDir() (string, error) {
	execPath, err := os.Executable()
	return path.Dir(execPath), err
}