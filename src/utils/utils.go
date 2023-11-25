package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateDirIfNotExists(path string, fileName string) {
	fullPathFileName := Path(path, fileName)
	Check(os.Mkdir(fullPathFileName, os.ModePerm), "Cannot create file "+fullPathFileName)
}

func CreateFileIfNotExists(path string, fileName string) {
	fileFullPath := Path(path, fileName)
	file, err := os.Create(fileFullPath)
	Check(err, "Cannot create file "+fileFullPath)
	file.Close()
}

func CreateFileIfNotExistsWithContent(path string, fileName string, initContent string) {
	fileFullPath := Path(path, fileName)
	file, err := os.Create(fileFullPath)
	Check(err, "Cannot create file "+fileFullPath)
	defer file.Close()

	_, err = file.Write([]byte(initContent))
	Check(err, "Cannot write "+initContent+" to "+fileFullPath)
}

func Path(path string, fileName string) string {
	return filepath.Join(path, fileName)
}

func Check(err error, message string) {
	if err != nil {
		ExitError(message)
	}
}

func ExitError(message string) {
	fmt.Fprintf(os.Stderr, message+"\n")
	os.Exit(1)
}
