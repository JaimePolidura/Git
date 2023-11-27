package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func JoinStrings(strings []string) string {
	resultString := ""
	for _, actualString := range strings {
		resultString = resultString + actualString
	}

	return resultString
}

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

func FindIndex(bytes []byte, offset int, char uint8) int {
	for i := offset; i < len(bytes); i++ {
		if bytes[i] == char {
			return i
		}
	}

	return -1
}

func ReadUntil(bytes []byte, offset int, char uint8) ([]byte, int, error) {
	untilEof := char == 0
	if offset >= len(bytes) {
		return nil, offset, errors.New("Incorrect bytes size")
	}

	result := make([]byte, 0)
	actual := bytes[offset]

	for actual != char {
		result = append(result, actual)

		offset = offset + 1

		if offset >= len(bytes) && !untilEof {
			return result, offset, errors.New("Unexpected EOF")
		}
		if offset >= len(bytes) && untilEof {
			return result, offset + 1, nil
		}

		actual = bytes[offset]
	}

	return result, offset + 1, nil
}

func CurrentPath() string {
	currentPath, err := os.Getwd()
	Check(err, "Cannot get the current path")
	return currentPath
}

func Path(path string, fileName string) string {
	return filepath.Join(path, fileName)
}

func Paths(paths ...string) string {
	return filepath.Join(paths...)
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

func CheckFileOrDirExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else {
		return false
	}
}
