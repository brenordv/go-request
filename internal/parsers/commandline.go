package parsers

import (
	"errors"
	"fmt"
	"os"
)

func GetConfigFilenames() []string {
	args := os.Args[1:]
	var files []string

	for _, arg := range args {
		if _, err := os.Stat(arg); errors.Is(err, os.ErrNotExist) {
			fmt.Printf("Provided runtime  config file could not be found: %s\n", arg)
			continue
		}
		files = append(files, arg)
	}

	return files
}