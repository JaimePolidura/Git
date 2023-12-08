package utils

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func CreateDirIfNotExists(path string, fileName string) {
	fullPathFileName := Path(path, fileName)
	Check(os.Mkdir(fullPathFileName, os.ModePerm), "Cannot create file "+fullPathFileName)
}

func CreateFileIfNotExists(path string, fileName string) {
	fileFullPath := Path(path, fileName)

	if _, err := os.Stat(fileFullPath); os.IsNotExist(err) {
		file, err := os.Create(fileFullPath)
		Check(err, "Cannot create file "+fileFullPath)
		file.Close()
	}
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
			return result, offset, nil
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

func CheckError(err error) {
	if err != nil {
		ExitError(err.Error())
	}
}

func ExitError(message string) {
	fmt.Fprintf(os.Stderr, message+"\n")
	os.Exit(1)
}

func IsValidGitHash(hash string) bool {
	match, err := regexp.MatchString("^[0-9A-Fa-f]{4,40}$", hash)
	return err == nil && match
}

func CheckFileOrDirExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else {
		return false
	}
}

func GetAllSubfiles(path string) map[string]string {
	dirFs := os.DirFS(path)
	results := make(map[string]string)
	fs.WalkDir(dirFs, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			results[path] = path
		}

		return nil
	})

	return results
}

func BoolToUint16(value bool) uint16 {
	if value {
		return 0x1
	} else {
		return 0x0
	}
}

func SanitizePath(path string) string {
	return strings.TrimRight(strings.Trim(path, " "), "\n")
}
