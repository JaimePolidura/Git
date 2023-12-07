package commands

import (
	"fmt"
	"git/src/objects"
	"git/src/repository"
	"git/src/utils"
)

// Log Args: main.go log [sha]
func Log(args []string) {
	if len(args) < 2 {
		utils.ExitError("Invalid arguments: log <sha>")
	}

	currentRepository, _, err := repository.FindCurrentRepository(utils.CurrentPath())
	commitSha := getCommitShaToStartIterating(args, currentRepository)
	currentCommitHasParent := true

	if err != nil {
		utils.ExitError(err.Error())
	}

	for currentCommitHasParent {
		commit := getGitCommitObject(currentRepository, commitSha)
		if !commit.HasParent() {
			currentCommitHasParent = false
		}

		printCommit(commit, commitSha)

		commitSha = commit.Parent
	}
}

func getCommitShaToStartIterating(args []string, currentRepository *repository.Repository) string {
	if len(args) >= 3 {
		return args[2]
	} else {
		name, detached, err := currentRepository.GetActiveBranch()
		utils.CheckError(err)
		if !detached {
			sha, _, err := currentRepository.ResolveObjectName(name, objects.ANY)
			utils.CheckError(err)
			return sha
		} else {
			return name
		}
	}
}

func printCommit(commitObject objects.CommitObject, sha string) {
	fmt.Println("Commit " + sha)
	fmt.Println("Author " + commitObject.Author)
	fmt.Println("")
	fmt.Println("\t" + commitObject.Message)
	fmt.Println("")
}

func getGitCommitObject(currentRepository *repository.Repository, sha string) objects.CommitObject {
	gitObject, err := currentRepository.ReadObject(sha, objects.COMMIT)
	if err != nil {
		utils.ExitError("Object with SHA " + sha + " not found")
	}

	return gitObject.SerializableGitObject.(objects.CommitObject)
}
