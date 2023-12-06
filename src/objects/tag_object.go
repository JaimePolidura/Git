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

func CreateTagObject(objectTag string, tag string, tagger string) *Object {
	keyValue := utils.CreateNavigationMap[string, string]()
	keyValue.Put("objectTag", objectTag)
	keyValue.Put("tagger", tagger)
	keyValue.Put("tag", tag)

	return &Object{
		Type: TAG,
		SerializableGitObject: TagObject{
			ObjectTag: objectTag,
			Tag:       tag,
			Tagger:    tagger,
			keyValue:  keyValue,
		},
	}
}

func deserializeTagObject(toDeserialize []byte) (TagObject, error) {
	deserializedKeyValue, _ := keyValueListDeserialize(toDeserialize)
	if allContained := deserializedKeyValue.Contains("objectTag", "tagger", "tag"); !allContained {
		return TagObject{}, errors.New("Invalid key value format. Missing fields")
	}

	tagObject := TagObject{
		ObjectTag: deserializedKeyValue.Get("objectTag"),
		Tagger:    deserializedKeyValue.Get("tagger"),
		Tag:       deserializedKeyValue.Get("tag"),
		keyValue:  deserializedKeyValue,
	}

	return tagObject, nil
}

func (c TagObject) Serialize() []byte {
	return keyValueListSerialize(c.keyValue)
}
