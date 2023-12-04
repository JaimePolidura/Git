package objects

import (
	"errors"
	"git/src/utils"
)

type CommitObject struct {
	Tree      string
	Parent    string
	Author    string
	Committer string
	Message   string

	keyValue *utils.NavigationMap[string, string]
}

func CreateCommitObject(treeSha string, parent string, author string, message string) *Object {
	commitObject := CommitObject{
		Tree:      treeSha,
		Parent:    parent,
		Author:    author,
		Committer: author,
		Message:   message,
		keyValue:  &utils.NavigationMap[string, string]{},
	}
	commitObject.keyValue.Put("tree", treeSha)
	commitObject.keyValue.Put("parent", parent)
	commitObject.keyValue.Put("author", author)
	commitObject.keyValue.Put("commiter", author)
	commitObject.keyValue.Put("message", message)

	return &Object{
		Type:                  COMMIT,
		SerializableGitObject: commitObject,
	}
}

func (c CommitObject) HasParent() bool {
	return true //TODO
}

func deserializeCommitObject(toDeserialize []byte) (CommitObject, error) {
	deserializedKeyValue, remainingData := keyValueListDeserialize(toDeserialize)
	if allContained := deserializedKeyValue.Contains("tree", "parent", "author", "committer"); !allContained {
		return CommitObject{}, errors.New("invalid key value format. Missing fields")
	}

	commitObject := CommitObject{
		Tree:      deserializedKeyValue.Get("tree"),
		Parent:    deserializedKeyValue.Get("parent"),
		Author:    deserializedKeyValue.Get("author"),
		Committer: deserializedKeyValue.Get("committer"),
		Message:   string(remainingData),
		keyValue:  deserializedKeyValue,
	}

	return commitObject, nil
}

func (c CommitObject) Serialize() []byte {
	return append(keyValueListSerialize(c.keyValue), []byte(c.Message)...)
}
