package objects

import (
	"bytes"
	"fmt"
	"git/src/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTreeObject_Serialize(t *testing.T) {
	expectedBytes := []byte("tree 108" + string('\x00') + "100357 README.md" + string('\x00') + "saaa5f1sf15fa5f1s15ssaaa5f1sf15fa5f1s15s100777 src" + string('\x00') + "a5fa5f1s1sa5fa5f1s1sa5a5f1sf1sa5a5f1sf1s")
	object := Object{
		Type: TREE,
		SerializableGitObject: TreeObject{
			Entries: []TreeEntry{
				{Mode: 100777, Sha: "a5fa5f1s1sa5fa5f1s1sa5a5f1sf1sa5a5f1sf1s", Path: "src"},
				{Mode: 100357, Sha: "saaa5f1sf15fa5f1s15ssaaa5f1sf15fa5f1s15s", Path: "README.md"},
			},
		},
	}

	serialized := object.Serialize()

	fmt.Println(string(serialized))

	assert.Equal(t, serialized, expectedBytes)
}

func TestTreeObject_TreeDeserialize(t *testing.T) {
	serializedBytes := []byte("tree 108" + string('\x00') + "100357 README.md" + string('\x00') + "saaa5f1sf15fa5f1s15ssaaa5f1sf15fa5f1s15s100777 src" + string('\x00') + "a5fa5f1s1sa5fa5f1s1sa5a5f1sf1sa5a5f1sf1s")
	expectedObject := Object{
		Type: TREE,
		SerializableGitObject: TreeObject{
			Entries: []TreeEntry{
				{Mode: 100357, Sha: "saaa5f1sf15fa5f1s15ssaaa5f1sf15fa5f1s15s", Path: "README.md"},
				{Mode: 100777, Sha: "a5fa5f1s1sa5fa5f1s1sa5a5f1sf1sa5a5f1sf1s", Path: "src"},
			},
		},
	}

	actualObject, err := DeserializeObject(bytes.NewReader(serializedBytes))

	assert.Nil(t, err)
	assert.Equal(t, expectedObject, actualObject)
}

func TestCommitObject_Deserialize(t *testing.T) {
	serializedBytes := []byte("commit 231" + string('\x00') + "tree 29ff16c9c14e2652b22f8b78bb08a5a07930c147\nparent 206941306e8a8af65b66eaaaea388a7ae24d49a0\n" +
		"author Thibault Polge <thibault@thb.lt> 1527025023 +0200\ncommitter Thibault Polge <thibault@thb.lt> 1527025044 +0200\n\nCreate first commit")

	actualObject, err := DeserializeObject(bytes.NewReader(serializedBytes))

	assert.True(t, err == nil)
	assert.Equal(t, actualObject.SerializableGitObject.(CommitObject).Tree, "29ff16c9c14e2652b22f8b78bb08a5a07930c147")
	assert.Equal(t, actualObject.SerializableGitObject.(CommitObject).Parent, "206941306e8a8af65b66eaaaea388a7ae24d49a0")
	assert.Equal(t, actualObject.SerializableGitObject.(CommitObject).Author, "Thibault Polge <thibault@thb.lt> 1527025023 +0200")
	assert.Equal(t, actualObject.SerializableGitObject.(CommitObject).Committer, "Thibault Polge <thibault@thb.lt> 1527025044 +0200")
	assert.Equal(t, actualObject.SerializableGitObject.(CommitObject).Message, "Create first commit")
	assert.Equal(t, actualObject.SerializableGitObject.(CommitObject).keyValue.Keys(), []string{"tree", "parent", "author", "committer"})
}

func TestCommitObject_Serialize(t *testing.T) {
	objectToSerializeKeyValue := utils.CreateNavigationMap[string, string]()
	objectToSerializeKeyValue.Put("tree", "29ff16c9c14e2652b22f8b78bb08a5a07930c147")
	objectToSerializeKeyValue.Put("parent", "206941306e8a8af65b66eaaaea388a7ae24d49a0")
	objectToSerializeKeyValue.Put("author", "Thibault Polge <thibault@thb.lt> 1527025023 +0200")
	objectToSerializeKeyValue.Put("committer", "Thibault Polge <thibault@thb.lt> 1527025044 +0200")
	objectToSerialize := Object{
		Type: COMMIT,
		SerializableGitObject: CommitObject{
			Message:  "Create first commit",
			keyValue: objectToSerializeKeyValue,
		},
	}

	expectedSerializedBytes := []byte("commit 231" + string('\x00') + "tree 29ff16c9c14e2652b22f8b78bb08a5a07930c147\nparent 206941306e8a8af65b66eaaaea388a7ae24d49a0\n" +
		"author Thibault Polge <thibault@thb.lt> 1527025023 +0200\ncommitter Thibault Polge <thibault@thb.lt> 1527025044 +0200\n\nCreate first commit")

	serialized := objectToSerialize.Serialize()

	fmt.Println(string(serialized))

	assert.Equal(t, serialized, expectedSerializedBytes)
}

/*
	To Parse

tree 29ff16c9c14e2652b22f8b78bb08a5a07930c147
parent 206941306e8a8af65b66eaaaea388a7ae24d49a0
author Thibault Polge <thibault@thb.lt> 1527025023 +0200
committer Thibault Polge <thibault@thb.lt> 1527025044 +0200

Create first commit
*/
func TestKeyValueListDeserialize(t *testing.T) {
	parsed, remaining := keyValueListDeserialize([]byte("tree 29ff16c9c14e2652b22f8b78bb08a5a07930c147\nparent 206941306e8a8af65b66eaaaea388a7ae24d49a0\nauthor " +
		"Thibault Polge <thibault@thb.lt> 1527025023 +0200\ncommitter Thibault Polge <thibault@thb.lt> 1527025044 +0200\n\nCreate first commit"))

	assert.Equal(t, parsed.Size(), 4)
	assert.Equal(t, parsed.Get("tree"), "29ff16c9c14e2652b22f8b78bb08a5a07930c147")
	assert.Equal(t, parsed.Get("parent"), "206941306e8a8af65b66eaaaea388a7ae24d49a0")
	assert.Equal(t, parsed.Get("author"), "Thibault Polge <thibault@thb.lt> 1527025023 +0200")
	assert.Equal(t, parsed.Get("committer"), "Thibault Polge <thibault@thb.lt> 1527025044 +0200")
	assert.Equal(t, string(remaining), "Create first commit")
}

func TestKeyValueListSerialize(t *testing.T) {
	bytes := []byte("tree 29ff16c9c14e2652b22f8b78bb08a5a07930c147\nparent 206941306e8a8af65b66eaaaea388a7ae24d49a0\nauthor " +
		"Thibault Polge <thibault@thb.lt> 1527025023 +0200\ncommitter Thibault Polge <thibault@thb.lt> 1527025044 +0200\n\n")
	parsed, _ := keyValueListDeserialize(bytes)

	assert.Equal(t, bytes, keyValueListSerialize(parsed))
}
