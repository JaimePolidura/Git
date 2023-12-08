package objects

import (
	"errors"
	"fmt"
	"git/src/utils"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

type ObjectType string

const (
	COMMIT ObjectType = "commit"
	BLOB   ObjectType = "blob"
	TREE   ObjectType = "tree"
	TAG    ObjectType = "tag"
	ANY    ObjectType = "none" //Used only for
)

type Object struct {
	SerializableGitObject
	Type ObjectType
}

type SerializableGitObject interface {
	Serialize() []byte
}

func (o Object) Serialize() []byte {
	serialized := o.SerializableGitObject.Serialize()
	header := []byte(string(o.Type) + " " + strconv.Itoa(len(serialized)) + string('\x00'))

	return append(header, serialized...)
}

func DeserializeObject(reader io.Reader) (Object, error) {
	commonObject, pendingToDeserialize, err := deserializeObjectCommonHeader(reader)

	if err != nil {
		return *commonObject, err
	}

	fmt.Println("----------------------------")
	fmt.Println(string(pendingToDeserialize))

	var gitObject SerializableGitObject
	switch commonObject.Type {
	case BLOB:
		gitObject, err = deserializeBlobObject(pendingToDeserialize)
	case COMMIT:
		gitObject, err = deserializeCommitObject(pendingToDeserialize)
	case TREE:
		gitObject, err = deserializeTreeObject(pendingToDeserialize)
	case TAG:
		gitObject, err = deserializeTagObject(pendingToDeserialize)
	}

	commonObject.SerializableGitObject = gitObject

	return *commonObject, err
}

func deserializeObjectCommonHeader(reader io.Reader) (*Object, []byte, error) {
	bytesDecompressed, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, []byte{}, err
	}

	objectTypeBytes, offset, err := utils.ReadUntil(bytesDecompressed, 0, 32)
	if err != nil {
		return nil, []byte{}, err
	}
	objectType, err := getObjectTypeByString(string(objectTypeBytes))
	if err != nil {
		return nil, []byte{}, err
	}

	_, offset, err = utils.ReadUntil(bytesDecompressed, offset, 0)
	if err != nil {
		return nil, []byte{}, err
	}
	restData := bytesDecompressed[offset:]

	return &Object{Type: objectType}, restData, nil
}

func getObjectTypeByString(objectTypeString string) (ObjectType, error) {
	switch strings.ToLower(objectTypeString) {
	case "commit":
		return COMMIT, nil
	case "blob":
		return BLOB, nil
	case "tree":
		return TREE, nil
	case "tag":
		return TAG, nil
	default:
		return "", errors.New("ObjectType " + objectTypeString + " not found")
	}
}

func keyValueListSerialize(kvMap *utils.NavigationMap[string, string]) []byte {
	result := ""

	for _, key := range kvMap.Keys() {
		value := kvMap.Get(key)
		result = result + key + " " + value + "\n"
	}

	return []byte(result + "\n")
}

func keyValueListDeserialize(bytes []byte) (*utils.NavigationMap[string, string], []byte) {
	return keyValueListParserDeserializeRecursive(bytes, 0, utils.CreateNavigationMap[string, string]())
}

func keyValueListParserDeserializeRecursive(bytes []byte, offset int, parsed *utils.NavigationMap[string, string]) (*utils.NavigationMap[string, string], []byte) {
	indexEndKey := utils.FindIndex(bytes, offset, 32)
	indexEndValue := utils.FindIndex(bytes, offset, 10)

	if indexEndValue > indexEndKey && indexEndKey > 0 && indexEndValue > 0 {
		key := string(bytes[offset:indexEndKey])
		value := string(bytes[indexEndKey+1 : indexEndValue])
		parsed.Put(key, value)

		return keyValueListParserDeserializeRecursive(bytes, indexEndValue+1, parsed)
	} else { //Blank line -> end of key/value
		return parsed, bytes[offset+1:]
	}
}
