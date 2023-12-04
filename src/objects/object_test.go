package objects

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTreeObject_Deserialize(t *testing.T) {
	bytes := []byte("100777 src" + string('\x00') + "a5fa5f1s1sa5a5f1sf1s100357 README.md" + string('\x00') + "saaa5f1sf15fa5f1s15s")
	deserialized, err := deserializeTreeObject(bytes)

	assert.Equal(t, err, nil)
	firstEntry := deserialized.Entries[0]
	secondEntry := deserialized.Entries[1]
	assert.Equal(t, firstEntry.Mode, 100777)
	assert.Equal(t, firstEntry.Path, "src")
	assert.Equal(t, firstEntry.Sha, "a5fa5f1s1sa5a5f1sf1s")

	assert.Equal(t, secondEntry.Mode, 100357)
	assert.Equal(t, secondEntry.Path, "README.md")
	assert.Equal(t, secondEntry.Sha, "saaa5f1sf15fa5f1s15s")
}

func TestTreeObject_Serialize(t *testing.T) {
	bytes := []byte("100777 src" + string('\x00') + "a5fa5f1s1sa5a5f1sf1s")
	deserialized, _ := deserializeTreeObject(bytes)

	assert.Equal(t, deserialized.Serialize(), bytes)
}

func TestCommitObject_Deserialize(t *testing.T) {
	bytes := []byte("tree 29ff16c9c14e2652b22f8b78bb08a5a07930c147\nparent 206941306e8a8af65b66eaaaea388a7ae24d49a0\n" +
		"author Thibault Polge <thibault@thb.lt> 1527025023 +0200\ncommitter Thibault Polge <thibault@thb.lt> 1527025044 +0200\n\nCreate first commit")

	object, err := deserializeCommitObject(bytes)

	//fmt.Println(err.Error())

	assert.True(t, err == nil)
	assert.Equal(t, object.Tree, "29ff16c9c14e2652b22f8b78bb08a5a07930c147")
	assert.Equal(t, object.Parent, "206941306e8a8af65b66eaaaea388a7ae24d49a0")
	assert.Equal(t, object.Author, "Thibault Polge <thibault@thb.lt> 1527025023 +0200")
	assert.Equal(t, object.Committer, "Thibault Polge <thibault@thb.lt> 1527025044 +0200")
	assert.Equal(t, object.Message, "Create first commit")
	assert.Equal(t, object.keyValue.Keys(), []string{"tree", "parent", "author", "committer"})
}

func TestCommitObject_Serialize(t *testing.T) {
	bytes := []byte("tree 29ff16c9c14e2652b22f8b78bb08a5a07930c147\nparent 206941306e8a8af65b66eaaaea388a7ae24d49a0\n" +
		"author Thibault Polge <thibault@thb.lt> 1527025023 +0200\ncommitter Thibault Polge <thibault@thb.lt> 1527025044 +0200\n\nCreate first commit")
	object, _ := deserializeCommitObject(bytes)

	serializedCommit := object.Serialize()

	assert.Equal(t, serializedCommit, bytes)
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
