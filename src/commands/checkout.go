package commands

import (
	"git/src/objects"
	"git/src/repository"
	"git/src/utils"
	"os"
)

// Checkout Args: main.go checkout <sha>
func Checkout(args []string) {
	if len(args) != 3 {
		utils.ExitError("Invalid arguments checkout <sha>")
	}

	currentRepository, repositoryPath, err := repository.FindCurrentRepository(utils.CurrentPath())
	utils.CheckError(err)

	objectNameUnResolved := args[2]
	sha, isHead, err := currentRepository.ResolveObjectName(objectNameUnResolved, objects.ANY)
	utils.CheckError(err)

	commitGitObjet := getCommitObject(currentRepository, sha)
	treeObject := getTreeGitObject(currentRepository, commitGitObjet.Tree)

	restoreRecursive(currentRepository, treeObject, repositoryPath)
	updateHead(isHead, currentRepository, sha, objectNameUnResolved)
}

func updateHead(isHead bool, currentRepository *repository.Repository, sha string, objectNameUnResolved string) {
	if isHead {
		currentRepository.WriteToHead("ref: refs/heads/" + objectNameUnResolved)
	} else {
		currentRepository.WriteToHead(sha)
	}
}

func restoreRecursive(currentRepository *repository.Repository, tree objects.TreeObject, currentPath string) {
	for _, treeEntry := range tree.Entries {
		pathEntry := utils.Paths(currentPath, treeEntry.Path)
		entryExistsInFS := utils.CheckFileOrDirExists(pathEntry)

		if !entryExistsInFS {
			createTreeEntryInFS(treeEntry, pathEntry)
		}

		if !treeEntry.IsDir() {
			blobGitObject := getBlobObject(currentRepository, treeEntry.Sha)
			file, err := os.Open(pathEntry)
			defer file.Close()
			utils.Check(err, "Cannot open file "+pathEntry)

			utils.Check(os.Truncate(pathEntry, 0), "Cannot clear file content: "+pathEntry)
			_, err = file.Write(blobGitObject.Data)
			utils.Check(err, "Cannot write to file "+pathEntry)

		} else {
			entryTreeObject := getTreeGitObject(currentRepository, treeEntry.Sha)
			restoreRecursive(currentRepository, entryTreeObject, pathEntry)
		}
	}
}

func createTreeEntryInFS(treeEntry objects.TreeEntry, fullPathEntry string) {
	if treeEntry.IsDir() {
		utils.Check(os.Mkdir(fullPathEntry, os.FileMode(0777)), "Cannot create directory: "+fullPathEntry)
	} else {
		file, err := os.Create(fullPathEntry)
		utils.Check(err, "Cannot create file in "+fullPathEntry)
		file.Close()
	}
}

func getBlobObject(currentRepository *repository.Repository, sha string) objects.BlobObject {
	blobGitObjet, err := currentRepository.ReadBlobObject(sha)
	utils.Check(err, "Commit object not found or object type is not BLOB")
	return blobGitObjet
}

func getCommitObject(currentRepository *repository.Repository, sha string) objects.CommitObject {
	commitGitObjet, err := currentRepository.ReadCommitObject(sha)
	utils.Check(err, "Commit object not found or object type is not COMMIT")
	return commitGitObjet
}

func getTreeObject(currentRepository *repository.Repository, sha string) objects.TreeObject {
	treeObject, err := currentRepository.ReadTreeObject(sha)
	utils.Check(err, "Commit object might be corrupted. It doesnt point to a Tree object")
	return treeObject
}
