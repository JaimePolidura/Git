package index

import (
	"encoding/binary"
	"git/src/utils"
	"io"
	"io/ioutil"
	"math"
	"os"
	"strconv"
)

type IndexObject struct {
	Version uint32
	Entries map[string]IndexEntry
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

func CreateIndexEntry(stats os.FileInfo, pathRelativeRepo string, sha string) IndexEntry {
	return IndexEntry{
		Ctime:           uint64(stats.ModTime().UnixNano()),
		Mtime:           uint64(stats.ModTime().UnixNano()),
		Dev:             0,
		Ino:             0,
		ModeType:        0x08,
		ModePerms:       0x01A,
		Uid:             1,
		Gid:             1,
		Fsize:           uint32(stats.Size()),
		Sha:             sha,
		FlagAssumeValid: false,
		FlagState:       false,
		FullPathName:    pathRelativeRepo,
	}
}

func (self *IndexObject) Serialize() []byte {
	bytes := make([]byte, 4)

	binary.BigEndian.AppendUint32(bytes, self.Version)
	binary.BigEndian.AppendUint32(bytes, uint32(len(self.Entries)))

	offset := 0
	for _, entry := range self.Entries {
		serializedEntryBytes := entry.Serialize()
		bytes = append(bytes, serializedEntryBytes...)
		offset += len(serializedEntryBytes)

		if offset%8 != 0 {
			pad := 8 - (offset % 8)
			for i := 0; i < pad; i++ {
				bytes = append(bytes, 0x00)
			}
			offset += 8
		}
	}

	return bytes
}

func (self *IndexEntry) Serialize() []byte {
	bytes := make([]byte, 0)

	binary.BigEndian.AppendUint64(bytes, self.Ctime)
	binary.BigEndian.AppendUint64(bytes, self.Mtime)
	binary.BigEndian.AppendUint32(bytes, self.Dev)

	binary.BigEndian.AppendUint32(bytes, self.Dev)
	binary.BigEndian.AppendUint32(bytes, self.Ino)
	binary.BigEndian.AppendUint32(bytes, self.ModeType<<12|self.ModePerms)
	binary.BigEndian.AppendUint32(bytes, self.Uid)
	binary.BigEndian.AppendUint32(bytes, self.Gid)
	binary.BigEndian.AppendUint32(bytes, self.Fsize)
	bytes = append(bytes, []byte(self.Sha)...)

	nameLength := uint16(len(self.FullPathName))

	flagAssumeValid := uint16(0)
	if self.FlagAssumeValid {
		flagAssumeValid = 0x1 << 15
	}
	binary.BigEndian.AppendUint16(bytes, uint16(flagAssumeValid|utils.BoolToUint16(self.FlagState)|nameLength))

	bytes = append(bytes, []byte(self.FullPathName)...)
	bytes = append(bytes, 0x00)

	return bytes
}

func Deserialize(reader io.Reader) (*IndexObject, error) {
	allBytes, err := ioutil.ReadAll(reader)

	if err != nil {
		return nil, nil
	}

	version := binary.BigEndian.Uint32(allBytes[4:8])
	count := binary.BigEndian.Uint32(allBytes[8:12])

	entries := make(map[string]IndexEntry)
	content := allBytes[12:]
	offset := 0

	for i := 0; i < int(count); i++ {
		entry, newOffset := deserializeIndexEntry(content, offset)
		offset = newOffset
		entries[entry.FullPathName] = entry
	}

	return &IndexObject{Version: version, Entries: entries}, nil
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
	flagAssumeValid := (flags & 0x8000) != 0
	flagState := (flags & 0x3000) != 0
	nameLength := flags & 0xFFF
	offset += 62

	name := string(content[offset : offset+int(nameLength)])
	offset += int(nameLength) + 1

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
		FlagAssumeValid: flagAssumeValid,
		FlagState:       flagState,
		FullPathName:    name,
	}, offset
}
