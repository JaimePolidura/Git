package objects

import (
	"errors"
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

type GitObject interface {
	Type() ObjectType
	Data() []byte

	serializeSpecificData() []byte
}

type Object struct {
	Type   ObjectType
	Length int
	Data   []byte
}

func serializeHeader(gitObject GitObject, length int) []byte {
	return []byte(string(gitObject.Type()) + " " + strconv.Itoa(length) + string('\x00'))
}

func SerializeObject(object GitObject) []byte {
	serialized := object.serializeSpecificData()
	header := serializeHeader(object, len(serialized))

	return append(header, serialized...)
}

func DeserializeObjectWithType[T GitObject](reader io.Reader) (T, error) {
	gitObject, err := DeserializeObject(reader)
	return gitObject.(T), err
}

func DeserializeObject(reader io.Reader) (GitObject, error) {
	commonObject, pendingToDeserialize, err := DeserializeObjectCommonHeader(reader)
	var gitObject GitObject

	if err != nil {
		return gitObject, err
	}

	switch commonObject.Type {
	case BLOB:
		gitObject, err = deserializeBlobObject(commonObject, pendingToDeserialize)
	case COMMIT:
		gitObject, err = deserializeCommitObject(commonObject, pendingToDeserialize)
	case TREE:
		gitObject, err = deserializeTreeObject(commonObject, pendingToDeserialize)
	case TAG:
		gitObject, err = deserializeTagObject(commonObject, pendingToDeserialize)
	}

	return gitObject, err
}

func DeserializeObjectCommonHeader(reader io.Reader) (*Object, []byte, error) {
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

	objectLengthBytes, offset, err := utils.ReadUntil(bytesDecompressed, offset, 0)
	if err != nil {
		return nil, []byte{}, err
	}
	objectLengthInt, err := strconv.Atoi(string(objectLengthBytes))
	if err != nil {
		return nil, []byte{}, err
	}
	restData, offset, err := utils.ReadUntil(bytesDecompressed, offset, 0)
	if err != nil {
		return nil, []byte{}, err
	}

	return &Object{Type: objectType, Length: objectLengthInt}, restData, nil
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
