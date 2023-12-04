package commands

import (
	"fmt"
	"git/src/repository"
	"git/src/utils"
	"strconv"
)

// LsFiles Args: maing.go ls-files
func LsFiles(args []string) {
	if len(args) != 2 {
		utils.ExitError("Invalid arguments: ls-files")
	}

	currentRepository, _, err := repository.FindCurrentRepository(utils.CurrentPath())
	if err != nil {
		utils.ExitError(err.Error())
	}

	index, err := currentRepository.ReadIndex()
	if err != nil {
		utils.ExitError("Cannot read index: " + err.Error())
	}

	fmt.Println("Index file format version", strconv.Itoa(int(index.Version)), "containing", strconv.Itoa(len(index.Entries)), "entries")

	for _, entry := range index.Entries {
		fmt.Println("\t", entry.FullPathName, "inode:", entry.Ino, "device:", entry.Dev, "size:", entry.Fsize, "sha:", entry.Sha)
	}

}
