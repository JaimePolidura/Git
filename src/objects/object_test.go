package objects

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
