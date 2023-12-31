package commands

import (
	"fmt"
	"git/src/index"
	"git/src/objects"
	"git/src/repository"
	"git/src/utils"
	"io/ioutil"
	"os"
	"strings"
)

// Commit Args: main.go commit -m <message...>
func Commit(args []string) {
	if len(args) < 4 {
		utils.ExitError("Invalid arguments: commit -m <message...>")
	}

	currentRepository, _, err := repository.FindCurrentRepository(utils.CurrentPath())
	if err != nil {
		utils.ExitError(err.Error())
	}

	if _, detached, _ := currentRepository.GetActiveBranch(); detached {
		utils.ExitError("You cannot commit changes while you are in a detached branch. You will have to checkout to head")
	}

	commitMessage := strings.Trim(strings.Join(args[3:], " "), "\"")

	index, err := currentRepository.ReadIndex()
	utils.CheckError(err)

	rootTreeSha := createBlobsAndTrees(index.ToTree(), currentRepository)

	commitSha := createCommitObject(rootTreeSha, commitMessage, currentRepository)

	fmt.Println("Commited changes:", commitSha)
}

func createCommitObject(treeSha string, commitMessage string, currentRepository *repository.Repository) string {
	parentCommit := getParentCommit(currentRepository)
	commitObject := objects.CreateCommitObject(treeSha, parentCommit, "Jaime", commitMessage)

	commitSha, err := currentRepository.WriteObject(commitObject)
	utils.CheckError(err)

	currentBranch, _, err := currentRepository.GetActiveBranch()
	utils.CheckError(err)

	file, err := os.OpenFile(utils.Paths(currentRepository.GitDir, "refs", "heads", currentBranch), os.O_WRONLY, 07777)
	defer file.Close()
	utils.CheckError(err)

	_, err = file.Write([]byte(commitSha))
	utils.CheckError(err)

	return commitSha
}

func getParentCommit(currentRepository *repository.Repository) string {
	head, _, err := currentRepository.ResolveObjectName("HEAD", objects.ANY)
	firstCommit := repository.IsErrorTypeNoCommitError(err)

	if firstCommit {
		return objects.NO_PARENT_COMMIT_SHA
	} else {
		return head
	}
}

func createBlobsAndTrees(node *index.IndexObjectTreeNode, repository *repository.Repository) string {
	objectsCreated := make(map[string]*index.IndexObjectTreeNode) //Sha -> indexEntry tree node

	for _, child := range node.Children {
		sha := createAndGetSha(child, repository)
		objectsCreated[sha] = child
	}

	return createTreeObject(objectsCreated, node, repository)
}

func createAndGetSha(node *index.IndexObjectTreeNode, repository *repository.Repository) string {
	if len(node.Children) == 0 { //is file
		indexEntryNode := node.Entry
		file, err := os.Open(indexEntryNode.FullPathName)
		defer file.Close()
		utils.CheckError(err)
		bytes, err := ioutil.ReadAll(file)
		utils.CheckError(err)

		blobObject := objects.CreateBlobObject(bytes)
		sha, err := repository.WriteObject(blobObject)
		utils.CheckError(err)

		return sha
	} else {
		return createBlobsAndTrees(node, repository)
	}
}

func createTreeObject(children map[string]*index.IndexObjectTreeNode, parent *index.IndexObjectTreeNode, repository *repository.Repository) string {
	treeEntries := make([]objects.TreeEntry, 0)

	for sha, indexEntryTreeNode := range children {
		a := objects.TreeEntry{
			Sha:  sha,
			Path: indexEntryTreeNode.Name,
		}

		treeEntries = append(treeEntries, a)
	}

	treeObject := &objects.Object{
		Type:                  objects.TREE,
		SerializableGitObject: objects.TreeObject{Entries: treeEntries},
	}

	treeSha, err := repository.WriteObject(treeObject)
	utils.CheckError(err)

	return treeSha
}
