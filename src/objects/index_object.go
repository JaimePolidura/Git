package objects

import (
	"encoding/binary"
	"io"
	"io/ioutil"
	"math"
	"strconv"
)

type IndexObject struct {
	Version uint32
	Entries []IndexEntry
}

type IndexEntry struct {
	Ctime           uint64
	Mtime           uint64
	Dev             uint32
	Ino             uint32
	ModeType        uint32
	ModePerms       uint32
	Uid             uint32
	Gid             uint32
	Fsize           uint32
	Sha             string
	FlagAssumeValid bool
	FlagState       bool
	FullPathName    string
}

func (i *IndexObject) deserialize(reader io.Reader) IndexObject {
	allBytes, _ := ioutil.ReadAll(reader)

	version := binary.BigEndian.Uint32(allBytes[4:8])
	count := binary.BigEndian.Uint32(allBytes[8:12])

	entries := make([]IndexEntry, count)
	content := allBytes[12:]
	offset := 0

	for i := 0; i < int(count); i++ {
		entry, newOffset := deserializeIndexEntry(content, offset)
		offset = newOffset
		entries = append(entries, entry)
	}

	return IndexObject{Version: version, Entries: entries}
}

func deserializeIndexEntry(content []byte, offset int) (IndexEntry, int) {
	ctime := binary.BigEndian.Uint64(content[offset : offset+8])
	mtime := binary.BigEndian.Uint64(content[offset+8 : offset+16])

	device := binary.BigEndian.Uint32(content[offset+16 : offset+20])
	ino := binary.BigEndian.Uint32(content[offset+20 : offset+24])
	mode := binary.BigEndian.Uint32(content[offset+26 : offset+28])
	modeType := mode >> 12
	modePerms := mode & 0x1FF
	uid := binary.BigEndian.Uint32(content[offset+28 : offset+32])
	gid := binary.BigEndian.Uint32(content[offset+32 : offset+36])
	fsize := binary.BigEndian.Uint32(content[offset+36 : offset+40])
	sha := strconv.FormatUint(uint64(binary.BigEndian.Uint32(content[offset+40:offset+60])), 16)

	flags := binary.BigEndian.Uint32(content[offset+60 : offset+62])
	flag_assume_valid := (flags & 0x8000) != 0
	flag_state := (flags & 0x3000) != 0
	name_length := flags & 0xFFF
	offset += 62

	name := string(content[offset : offset+int(name_length)])
	offset += int(name_length) + 1

	offset = 8 * int(math.Ceil(float64(offset)/8))

	return IndexEntry{
		Ctime:           ctime,
		Mtime:           mtime,
		Dev:             device,
		Ino:             ino,
		ModeType:        modeType,
		ModePerms:       modePerms,
		Uid:             uid,
		Gid:             gid,
		Fsize:           fsize,
		Sha:             sha,
		FlagAssumeValid: flag_assume_valid,
		FlagState:       flag_state,
		FullPathName:    name,
	}, offset
}
