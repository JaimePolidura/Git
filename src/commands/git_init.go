package commands

import (
	"git/src/repository"
	"git/src/utils"
	"os"
)

func Init() {
	currentPath, err := os.Getwd()
	utils.Check(err, "Cannot get the current path")
	repository.InitiliazeRepository(currentPath)
}
