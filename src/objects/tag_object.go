package objects

import (
	"errors"
	"git/src/utils"
)

type Reference struct {
	NamePath string
	Value    string
}

type TagObject struct {
	ObjectTag string
	Tagger    string
	Tag       string

	keyValue *utils.NavigationMap[string, string]
}

func deserializeTagObject(toDeserialize []byte) (TagObject, error) {
	deserializedKeyValue, _ := keyValueListDeserialize(toDeserialize)
	if allContained := deserializedKeyValue.Contains("object", "tagger", "tag"); !allContained {
		return TagObject{}, errors.New("Invalid key value format. Missing fields")
	}

	tagObject := TagObject{
		ObjectTag: deserializedKeyValue.Get("object"),
		Tagger:    deserializedKeyValue.Get("tagger"),
		Tag:       deserializedKeyValue.Get("tag"),
		keyValue:  deserializedKeyValue,
	}

	return tagObject, nil
}

func (c TagObject) Serialize() []byte {
	return keyValueListSerialize(c.keyValue)
}
