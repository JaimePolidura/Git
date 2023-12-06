package index

import (
	"strings"
)

type IndexObjectTreeNode struct {
	Entry    IndexEntry
	Children map[string]*IndexObjectTreeNode //Keys are the name of the files, not the pathds
	Root     bool
	Name     string //Name of the file, not the path
}

func createRootTreeNode() *IndexObjectTreeNode {
	return &IndexObjectTreeNode{
		Entry:    IndexEntry{},
		Children: make(map[string]*IndexObjectTreeNode, 0),
		Root:     true,
		Name:     "",
	}
}

func createChildDirNode(dirName string) *IndexObjectTreeNode {
	return &IndexObjectTreeNode{
		Entry:    IndexEntry{},
		Children: make(map[string]*IndexObjectTreeNode, 0),
		Root:     false,
		Name:     dirName,
	}
}

func createChildFileNode(fileName string, entry IndexEntry) *IndexObjectTreeNode {
	return &IndexObjectTreeNode{
		Entry:    entry,
		Children: make(map[string]*IndexObjectTreeNode, 0),
		Root:     false,
		Name:     fileName,
	}
}

func (self *IndexObject) ToTree() *IndexObjectTreeNode {
	root := createRootTreeNode()

	for pathIndexEntry, indexEntry := range self.Entries {
		splitedBySep := strings.Split(pathIndexEntry, "/")

		if len(splitedBySep) > 1 {
			parents := splitedBySep[:len(splitedBySep)-1]
			child := splitedBySep[len(splitedBySep)-1]
			lastNode := root

			for _, parentPath := range parents {
				parentInLastNode, parentInLastNodeAlreadyCreated := lastNode.Children[parentPath]

				if !parentInLastNodeAlreadyCreated {
					parentNode := createChildDirNode(parentPath)
					lastNode.Children[parentPath] = parentNode
					lastNode = parentNode
				} else {
					lastNode = parentInLastNode
				}
			}

			lastNode.Children[child] = createChildFileNode(child, indexEntry)
		} else {
			root.Children[pathIndexEntry] = createChildFileNode(indexEntry.FullPathName, indexEntry)
		}
	}

	return root
}
