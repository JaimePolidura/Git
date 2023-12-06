package commands

import (
	"fmt"
	"git/src/index"
	"git/src/objects"
	"git/src/repository"
	"git/src/utils"
	"os"
)

func Status() {
	currentRepository, _, err := repository.FindCurrentRepository(utils.CurrentPath())
	if err != nil {
		utils.ExitError(err.Error())
	}
	repositoryIndex, err := currentRepository.ReadIndex()
	if err != nil {
		utils.ExitError("No commits haven been made in this repository")
	}

	printBranchStatus(currentRepository)
	printChangesBetweenHeadAndIndex(currentRepository, repositoryIndex)
	printChangesBetweenWorktreeAndIndex(currentRepository, repositoryIndex)
}

// files not in stagging area
func printChangesBetweenWorktreeAndIndex(repository *repository.Repository, index *index.IndexObject) {
	fmt.Println("Changes not stagged for commit:")

	fileNamesInWorkTree := utils.GetAllSubfiles(repository.WorkTree)

	for _, entry := range index.Entries {
		fileExists := utils.CheckFileOrDirExists(entry.FullPathName)

		if fileExists {
			stats, err := os.Stat(entry.FullPathName)
			utils.Check(err, "Cannot get stats from file "+entry.FullPathName)
			if stats.ModTime().UnixNano() > int64(entry.Ctime) || stats.ModTime().UnixNano() > int64(entry.Mtime) {
				fmt.Println("modified " + entry.FullPathName)
			}
		}
		if !fileExists {
			fmt.Println(" deleted " + entry.FullPathName)
		}

		if _, contained := fileNamesInWorkTree[entry.FullPathName]; contained {
			delete(fileNamesInWorkTree, entry.FullPathName)
		}
	}

	fmt.Println("\nUntracked files:")
	for untrackedFilePath, _ := range fileNamesInWorkTree {
		fmt.Println(" ", untrackedFilePath)
	}
}

// Stagging area compared to head
func printChangesBetweenHeadAndIndex(repository *repository.Repository, index *index.IndexObject) {
	fmt.Println("Changes to be commited:")

	treeObjectMapHead := getTreeObjectMapFromHEAD(repository)

	for _, indexEntry := range index.Entries {
		shaInHead, containedInHead := treeObjectMapHead[indexEntry.FullPathName]
		if containedInHead && shaInHead != indexEntry.Sha {
			fmt.Println(" modified " + indexEntry.FullPathName)
		}
		if containedInHead {
			delete(treeObjectMapHead, indexEntry.FullPathName)
		}
		if !containedInHead {
			fmt.Println(" added " + indexEntry.FullPathName)
		}
	}

	for key, _ := range treeObjectMapHead {
		fmt.Println(" deleted " + key)
	}
}

func getTreeObjectMapFromHEAD(repository *repository.Repository) map[string]string {
	treeHeadCommitSha, err := repository.ResolveObjectName("HEAD", objects.TREE)
	if err != nil {
		utils.ExitError("Cannot get HEAD reference: " + err.Error())
	}
	treeObject, err := repository.ReadTreeObject(treeHeadCommitSha)
	if err != nil {
		utils.ExitError("Cannot get tree object from sha : " + treeHeadCommitSha + " error: " + err.Error())
	}

	results := make(map[string]string)
	getTreeObjectMapFromHEADRecursive(repository, treeObject.Entries, "", results)

	return results
}

func getTreeObjectMapFromHEADRecursive(repository *repository.Repository, actualTreeEntries []objects.TreeEntry, prevPath string, results map[string]string) {
	for _, actualTreeEntry := range actualTreeEntries {
		actualPath := utils.Path(prevPath, actualTreeEntry.Path)

		if actualTreeEntry.IsDir() {
			subTreeObject, err := repository.ReadTreeObject(actualTreeEntry.Sha)
			if err != nil {
				utils.ExitError("Cannot get tree object from sha : " + actualTreeEntry.Path + " error: " + err.Error())
			}

			getTreeObjectMapFromHEADRecursive(repository, subTreeObject.Entries, actualPath, results)
		} else {
			results[actualPath] = actualTreeEntry.Sha
		}
	}
}

func printBranchStatus(repository *repository.Repository) {
	branchName, detatched, err := repository.GetActiveBranch()
	if err != nil {
		utils.ExitError(err.Error())
	}

	if detatched {
		fmt.Println("HEAD detached at", detatched)
	} else {
		fmt.Println("On branch", branchName)
	}
}
