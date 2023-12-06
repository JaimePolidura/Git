package commands

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"git/src/index"
	"git/src/repository"
	"git/src/utils"
	"io/ioutil"
	"os"
	"strings"
)

// main.go add <pathds>
func Add(args []string) {
	if len(args) < 3 {
		utils.ExitError("Invalid args: add <file names in current path...>")
	}
	currentPath := utils.CurrentPath()
	currentRepository, _, err := repository.FindCurrentRepository(currentPath)
	utils.CheckError(err)
	indexRepository, err := currentRepository.ReadIndex()
	utils.CheckError(err)

	pathsToAdd := args[2:]

	for _, pathToAdd := range pathsToAdd {
		isAbsolute := strings.HasPrefix(pathToAdd, "/")
		allSubfilesMode := pathToAdd == "."
		if !allSubfilesMode && isAbsolute && !strings.Contains(pathToAdd, currentRepository.WorkTree) {
			fmt.Println("Cannot add " + pathToAdd + " doesnt belong to repository")
		}
		if !allSubfilesMode && !utils.CheckFileOrDirExists(pathToAdd) {
			fmt.Println("Cannot add " + pathToAdd + " doest exist")
		}

		pathInRepository := currentRepository.GetPathFileInRepository(pathToAdd)

		if allSubfilesMode {
			addSubfiles(currentRepository, indexRepository, pathInRepository)
		} else {
			add(currentRepository, indexRepository, pathInRepository)
		}
	}

	if err := currentRepository.WriteIndex(indexRepository); err != nil {
		fmt.Println("Cannot write to INDEX: " + err.Error())
	}
}

func addSubfiles(currentRepository *repository.Repository, indexObject *index.IndexObject, pathRelativeRepo string) *index.IndexObject {
	children, _ := os.ReadDir(pathRelativeRepo)
	for _, child := range children {
		add(currentRepository, indexObject, utils.Paths(pathRelativeRepo, child.Name()))
	}

	return indexObject
}

func add(currentRepository *repository.Repository, indexObject *index.IndexObject, pathRelativeRepo string) {
	stat, err := os.Stat(pathRelativeRepo)
	if err != nil {
		fmt.Println("Cannot get stat info of file " + pathRelativeRepo)
	}

	if ignored, _ := currentRepository.IsIgnored(pathRelativeRepo); ignored {
		return
	}

	if stat.IsDir() {
		addSubfiles(currentRepository, indexObject, pathRelativeRepo)
	} else {
		addFile(indexObject, pathRelativeRepo, stat)
	}
}

func addFile(indexObject *index.IndexObject, pathRelativeRepo string, stat os.FileInfo) {
	indexEntry, indexEntryExists := indexObject.Entries[pathRelativeRepo]

	if indexEntryExists {
		modified := stat.ModTime().UnixNano() > int64(indexEntry.Ctime) || stat.ModTime().UnixNano() > int64(indexEntry.Mtime)
		if modified {
			indexObject.Entries[pathRelativeRepo] = index.CreateIndexEntry(stat, pathRelativeRepo, getSha(pathRelativeRepo))
		}
	} else {
		fmt.Println(pathRelativeRepo)
		indexObject.Entries[pathRelativeRepo] = index.CreateIndexEntry(stat, pathRelativeRepo, getSha(pathRelativeRepo))
	}
}

func getSha(filePath string) string {
	file, err := os.Open(filePath)
	defer file.Close()
	utils.CheckError(err)

	bytes, err := ioutil.ReadAll(file)
	utils.CheckError(err)

	sha1Hasher := sha1.New()
	sha1Hasher.Write(bytes)

	return hex.EncodeToString(sha1Hasher.Sum(nil))
}
