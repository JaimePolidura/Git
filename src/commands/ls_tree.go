package commands

import (
	"fmt"
	"git/src/objects"
	"git/src/repository"
	"git/src/utils"
)

// LsTree Args: main.go ls-tree -r <sha>
func LsTree(args []string) {
	if len(args) != 3 {
		utils.ExitError("Invalid arguments: ls-tree -r <sha>>")
	}

	sha := args[3]

	currentRepository, err := repository.FindCurrentRepository(utils.CurrentPath())
	if err != nil {
		utils.ExitError(err.Error())
	}

	gitObject := getTreeGitObject(currentRepository, sha)
	
	printEntriesRecursive(currentRepository, gitObject.Entries)
}

func printEntriesRecursive(repository *repository.Repository, entries []objects.TreeEntry) {
	for _, entry := range entries {
		if entry.IsDir() {
			dirGitObject := getTreeGitObject(repository, entry.Sha)
			printEntriesRecursive(repository, dirGitObject.Entries)
		} else {
			fmt.Println(entry.Path + " " + entry.Sha)
		}
	}
}

func getTreeGitObject(repository *repository.Repository, sha string) objects.TreeObject {
	gitObject, err := repository.ReadObject(sha)
	if err != nil {
		utils.ExitError("Object with SHA " + sha + " not found")
	}
	if gitObject.Type() != objects.TREE {
		utils.ExitError("Object has invalid type. Only TREE types are allowed")
	}

	return gitObject.(objects.TreeObject)
}
