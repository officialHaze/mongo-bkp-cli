package util

import (
	"log"
	"os"
)

func IsDir(dirpath string) bool {
	fileinfo, err := os.Stat(dirpath)
	if err != nil {
		log.Println(err)
		return false
	}

	return fileinfo.IsDir()
}
