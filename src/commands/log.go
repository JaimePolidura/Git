package commands

import (
	"fmt"
	"git/src/objects"
	"git/src/repository"
	"git/src/utils"
)

// Log Args: main.go log <sha>
func Log(args []string) {
	if len(args) != 3 {
		utils.ExitError("Invalid arguments: log <sha>")
	}

	currentSha := args[2]
	currentRepository, err := repository.FindCurrentRepository(utils.CurrentPath())
	currentCommitHasParent := true

	if err != nil {
		utils.ExitError(err.Error())
	}

	for currentCommitHasParent {
		commit := getGitCommitObject(currentRepository, currentSha)
		if !commit.HasParent() {
			currentCommitHasParent = false
		}

		printCommit(commit, currentSha)

		currentSha = commit.Parent
	}
}

func printCommit(commitObject objects.CommitObject, sha string) {
	fmt.Println("commit " + sha)
	fmt.Println("Author " + commitObject.Author)
	fmt.Println("")
	fmt.Println("\t" + string(commitObject.Data()))
	fmt.Println("")
}

func getGitCommitObject(currentRepository *repository.Repository, sha string) objects.CommitObject {
	gitObject, err := currentRepository.ReadObject(sha)
	if err != nil {
		utils.ExitError("Object with SHA " + sha + " not found")
	}
	if gitObject.Type() != objects.COMMIT {
		utils.ExitError("Object has invalid type. Only COMMIT types are allowed")
	}

	return gitObject.(objects.CommitObject)
}
