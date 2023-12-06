package commands

import (
	"git/src/repository"
	"git/src/utils"
)

//Initializes git repository
func Init() {
	currentPath := utils.CurrentPath()
	repository.InitializeRepository(currentPath)
}
