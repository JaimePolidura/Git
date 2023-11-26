package commands

import (
	"git/src/repository"
	"git/src/utils"
	"os"
)

// CatFile Args: main.go cat-file <sha>
func CatFile(args []string) {
	if len(args) != 3 {
		utils.ExitError("Invalid arguemnts cat-file <type> <sha>")
	}

	sha := args[2]

	currentPath := utils.CurrentPath()
	repository, err := repository.FindCurrentRepository(currentPath)

	if err != nil {
		utils.ExitError(err.Error())
	}

	object, err := repository.ReadObject(sha)
	if err != nil {
		utils.ExitError("Cannot read object: " + err.Error())
	}

	os.Stdout.Write(object.Data)
}
