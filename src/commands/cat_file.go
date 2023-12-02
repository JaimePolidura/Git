package commands

import (
	"git/src/objects"
	"git/src/repository"
	"git/src/utils"
	"os"
)

// CatFile Args: main.go cat-file <sha>
func CatFile(args []string) {
	if len(args) != 3 {
		utils.ExitError("Invalid arguments cat-file <type> <sha>")
	}

	sha := args[2]
	currentPath := utils.CurrentPath()
	currentRepository, _, err := repository.FindCurrentRepository(currentPath)

	if err != nil {
		utils.ExitError(err.Error())
	}

	object, err := currentRepository.ReadObject(sha, objects.NONE)
	if err != nil {
		utils.ExitError("Cannot read object: " + err.Error())
	}

	_, _ = os.Stdout.Write(object.Data())
}
