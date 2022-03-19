package core

import (
	"os"
)

func GetRootPath() (rootPath string) {
	rootPath, err := os.Getwd()
	if err != nil {
		return "."
	}
	return
}
