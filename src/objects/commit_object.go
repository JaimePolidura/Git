package objects

import (
	"errors"
	"git/src/utils"
)

type CommitObject struct {
	Object
	Tree      string
	Parent    string
	Author    string
	Committer string
	Message   string

	keyValue *utils.NavigationMap[string, string]
}

func CreateCommitObject(treeSha string, parent string, author string, message string) CommitObject {
	commitObject := CommitObject{
		Object:    Object{Type: COMMIT},
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

	return commitObject
}

func (c CommitObject) Type() ObjectType {
	return c.Object.Type
}

func (c CommitObject) Length() int {
	return c.Object.Length
}

func (c CommitObject) Data() []byte {
	return c.Object.Data
}

func (c CommitObject) HasParent() bool {
	return true //TODO
}

func deserializeCommitObject(commonObject *Object, toDeserialize []byte) (CommitObject, error) {
	deserializedKeyValue, remainingData := keyValueListDeserialize(toDeserialize)
	if allContained := deserializedKeyValue.Contains("tree", "parent", "author", "committer"); !allContained {
		return CommitObject{}, errors.New("Invalid key value format. Missing fields")
	}

	commitObject := CommitObject{
		Object:    *commonObject,
		Tree:      deserializedKeyValue.Get("tree"),
		Parent:    deserializedKeyValue.Get("parent"),
		Author:    deserializedKeyValue.Get("author"),
		Committer: deserializedKeyValue.Get("committer"),
		Message:   string(remainingData),
		keyValue:  deserializedKeyValue,
	}

	commitObject.Object.Data = remainingData
	commitObject.Object.Length = len(remainingData)

	return commitObject, nil
}

func (c CommitObject) serializeSpecificData() []byte {
	return append(keyValueListSerialize(c.keyValue), c.Object.Data...)
}
