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
	Object

	ObjectTag string
	Tagger    string
	Tag       string

	keyValue *utils.NavigationMap[string, string]
}

func (c TagObject) Type() ObjectType {
	return c.Object.Type
}

func (c TagObject) Length() int {
	return c.Object.Length
}

func (c TagObject) Data() []byte {
	return c.Object.Data
}

func (c TagObject) HasParent() bool {
	return true //TODO
}

func deserializeTagObject(commonObject *Object, toDeserialize []byte) (TagObject, error) {
	deserializedKeyValue, remainingData := keyValueListDeserialize(toDeserialize)
	if allContained := deserializedKeyValue.Contains("object", "tagger", "tag"); !allContained {
		return TagObject{}, errors.New("Invalid key value format. Missing fields")
	}

	tagObject := TagObject{
		Object:    *commonObject,
		ObjectTag: deserializedKeyValue.Get("object"),
		Tagger:    deserializedKeyValue.Get("tagger"),
		Tag:       deserializedKeyValue.Get("tag"),
		keyValue:  deserializedKeyValue,
	}

	tagObject.Object.Data = remainingData
	tagObject.Object.Length = len(remainingData)

	return tagObject, nil
}

func (c TagObject) serializeSpecificData() []byte {
	return append(keyValueListSerialize(c.keyValue), c.Object.Data...)
}
