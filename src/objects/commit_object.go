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

	keyValue *utils.NavigationMap[string, string]
}

func deserializeCommitObject(commonObject *Object, toDeserialize []byte) (*CommitObject, error) {
	deserializedKeyValue, data := KeyValueListDeserialize(toDeserialize)
	if allContained := deserializedKeyValue.Contains("tree", "parent", "author", "committer", "remaining"); !allContained {
		return nil, errors.New("Invalid key value format. Missing fields")
	}

	commitObject := &CommitObject{
		Object:    *commonObject,
		Tree:      deserializedKeyValue.Get("tree"),
		Parent:    deserializedKeyValue.Get("parent"),
		Author:    deserializedKeyValue.Get("author"),
		Committer: deserializedKeyValue.Get("committer"),
		keyValue:  deserializedKeyValue,
	}

	commitObject.Data = data

	return commitObject, nil
}

func (c *CommitObject) Serialize() []byte {
	header := c.serializeHeader()
	keyValueSerialized := KeyValueListSerialize(c.keyValue)

	return append(append(header, keyValueSerialized...), c.Data...)
}
