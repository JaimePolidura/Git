package commands

import (
	"fmt"
	"git/src/repository"
	"git/src/utils"
)

// CheckIgnore Args main.go check-ignore <file> [more files...]
func CheckIgnore(args []string) {
	if len(args) < 3 {
		utils.ExitError("Invalid arguments: check-ignore <file> [more files...]")
	}

	currentRepository, _, err := repository.FindCurrentRepository(utils.CurrentPath())
	if err != nil {
		utils.ExitError(err.Error())
	}

	for _, fileNameToCheck := range args[2:] {
		ignored, err := currentRepository.IsIgnored(fileNameToCheck)
		if err != nil {
			utils.ExitError(err.Error())
		}
		if ignored {
			fmt.Println(fileNameToCheck)
		}
	}
}
