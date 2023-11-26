package objects

import (
	"errors"
	"strconv"
	"strings"
)

type ObjectType string

const (
	COMMIT ObjectType = "commit"
	BLOB   ObjectType = "blob"
	TREE   ObjectType = "tree"
	TAG    ObjectType = "tag"
)

type Object struct {
	Type   ObjectType
	Length int
	Data   []byte
}

func (o *Object) serializeHeader() []byte {
	return []byte(string(o.Type) + " " + strconv.Itoa(o.Length) + string('\x00'))
}

func GetObjectTypeByString(objectTypeString string) (ObjectType, error) {
	switch strings.ToLower(objectTypeString) {
	case "commit":
		return COMMIT, nil
	case "blob":
		return BLOB, nil
	case "tree":
		return TREE, nil
	case "TAG":
		return TAG, nil
	default:
		return "", errors.New("ObjectType " + objectTypeString + " not found")
	}
}

func Serialize(object *Object) []byte {
	bytes := object.serializeHeader()

	switch object.Type {
	case COMMIT:
		return nil
	case BLOB:
		return append(bytes, serializeBlob(object)...)
	case TREE:
		return nil
	case TAG:
		return nil
	}

	panic("wtf?")
}
