package objects

import (
	"bytes"
	"git/src/utils"
	"sort"
	"strings"
)

type TreeObject struct {
	Entries []TreeEntry
}

type TreeEntry struct {
	Sha  string
	Path string
}

func (t TreeObject) Serialize() []byte {
	t.sortEntries()
	var bufferResult bytes.Buffer

	for _, entry := range t.Entries {
		_, _ = bufferResult.Write(entry.serialize())
	}

	return bufferResult.Bytes()
}

func (t TreeObject) sortEntries() {
	sort.Slice(t.Entries, func(i, j int) bool {
		entryA := t.Entries[i]
		entryB := t.Entries[j]

		return entryA.formatPathToSort() < entryB.formatPathToSort()
	})
}

func (t TreeEntry) serialize() []byte {
	return []byte(t.Path + "\x00" + t.Sha)
}

func (t TreeEntry) formatPathToSort() string {
	if t.IsDir() {
		return t.Path + "/"
	} else {
		return t.Path
	}
}

func (t TreeEntry) IsDir() bool {
	return strings.HasPrefix(t.Path, "10")
}

func deserializeTreeObject(toDeserialize []byte) (TreeObject, error) {
	entries := make([]TreeEntry, 0)
	actualOffset := 0

	for len(toDeserialize) > actualOffset {
		if treeEntryDeserialized, newOffset, err := deserializeTreeObjectEntry(toDeserialize, actualOffset); err == nil {
			entries = append(entries, treeEntryDeserialized)
			actualOffset = newOffset
		} else {
			return TreeObject{}, err
		}
	}

	return TreeObject{Entries: entries}, nil
}

func deserializeTreeObjectEntry(bytes []byte, offset int) (TreeEntry, int, error) {
	pathBytes, offset, err := utils.ReadUntil(bytes, offset, 0)
	if err != nil {
		return TreeEntry{}, -1, err
	}
	shaBytes := bytes[offset : offset+40]

	offset = offset + 40

	return TreeEntry{
		Sha:  string(shaBytes),
		Path: string(pathBytes),
	}, offset, nil
}
