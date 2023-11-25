package commands

import (
	"git/src/repository"
	"git/src/utils"
)

func Init() {
	currentPath := utils.CurrentPath()
	repository.InitializeRepository(currentPath)
}
