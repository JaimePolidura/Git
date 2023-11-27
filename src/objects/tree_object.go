package objects

import (
	"bytes"
	"git/src/utils"
	"sort"
	"strconv"
	"strings"
)

type TreeObject struct {
	Object

	Entries []TreeEntry
}

type TreeEntry struct {
	Mode int
	Sha  string
	Path string
}

func (t TreeObject) serializeSpecificData() []byte {
	t.shortEntries()
	var bufferResult bytes.Buffer

	for _, entry := range t.Entries {
		_, _ = bufferResult.Write(entry.serialize())
	}

	return bufferResult.Bytes()
}

func (t TreeObject) shortEntries() {
	sort.Slice(t.Entries, func(i, j int) bool {
		entryA := t.Entries[i]
		entryB := t.Entries[j]

		return entryA.formatPathToSort() < entryB.formatPathToSort()
	})
}

func (t TreeEntry) serialize() []byte {
	return []byte(strconv.Itoa(t.Mode) + " " + t.Path + "\x00" + t.Sha)
}

func (t TreeEntry) formatPathToSort() string {
	if t.IsDir() {
		return t.Path + "/"
	} else {
		return t.Path
	}
}

func (t TreeEntry) GetPermissions() int {
	return t.Mode & 0x0F
}

func (t TreeEntry) IsDir() bool {
	return strings.HasPrefix(t.Path, "10")
}

func deserializeTreeObject(commonObject *Object, toDeserialize []byte) (TreeObject, error) {
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

	commonObject.Data = toDeserialize

	return TreeObject{
		Entries: entries,
		Object:  *commonObject,
	}, nil
}

func deserializeTreeObjectEntry(bytes []byte, offset int) (TreeEntry, int, error) {
	modeBytes := bytes[offset : offset+6]
	pathBytes, offset, err := utils.ReadUntil(bytes, offset+7, 0)
	if err != nil {
		return TreeEntry{}, -1, err
	}
	shaBytes := bytes[offset : offset+20]

	offset = offset + 20

	modeString, err := strconv.Atoi(string(modeBytes))
	if err != nil {
		return TreeEntry{}, -1, err
	}

	return TreeEntry{
		Mode: modeString,
		Sha:  string(shaBytes),
		Path: string(pathBytes),
	}, offset, nil
}

func (t TreeObject) Type() ObjectType {
	return t.Object.Type
}

func (t TreeObject) Length() int {
	return t.Object.Length
}

func (t TreeObject) Data() []byte {
	return t.Object.Data
}
