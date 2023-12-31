package commands

import (
	"fmt"
	"git/src/objects"
	"git/src/repository"
	"git/src/utils"
)

// RevParse Args: main.go rev-parse <name>
func RevParse(args []string) {
	if len(args) != 3 {
		utils.ExitError("Invalid args: rev-parse <name>")
	}

	currentRepository, _, err := repository.FindCurrentRepository(utils.CurrentPath())
	if err != nil {
		utils.ExitError(err.Error())
	}

	objectName := args[2]
	hash, _, err := currentRepository.ResolveObjectName(objectName, objects.ANY)

	if err != nil {
		utils.ExitError("Reference: " + objectName + " not found")
	}

	fmt.Println(hash)
}
