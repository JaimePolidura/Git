package index

import (
	"encoding/binary"
	"io"
	"io/ioutil"
	"os"
)

type IndexObject struct {
	Version uint32
	Entries map[string]IndexEntry
}

type IndexEntry struct {
	Ctime        uint64
	Mtime        uint64
	Dev          uint32
	Ino          uint32
	ModeType     uint32
	ModePerms    uint32
	Uid          uint32
	Gid          uint32
	Fsize        uint32
	Sha          string
	FullPathName string
}

func CreateIndexEntry(stats os.FileInfo, pathRelativeRepo string, sha string) IndexEntry {
	return IndexEntry{
		Ctime:        uint64(stats.ModTime().UnixNano()),
		Mtime:        uint64(stats.ModTime().UnixNano()),
		Dev:          0,
		Ino:          0,
		ModeType:     0x08,
		ModePerms:    0x01A,
		Uid:          1,
		Gid:          1,
		Fsize:        uint32(stats.Size()),
		Sha:          sha,
		FullPathName: pathRelativeRepo,
	}
}

func (self *IndexObject) Serialize() []byte {
	bytes := make([]byte, 0)

	bytes = binary.BigEndian.AppendUint32(bytes, self.Version)
	bytes = binary.BigEndian.AppendUint32(bytes, uint32(len(self.Entries)))

	for _, entry := range self.Entries {
		serializedEntryBytes := entry.Serialize()
		bytes = append(bytes, serializedEntryBytes...)
	}

	return bytes
}

func (self *IndexEntry) Serialize() []byte {
	bytes := make([]byte, 0)

	bytes = binary.BigEndian.AppendUint64(bytes, self.Ctime)
	bytes = binary.BigEndian.AppendUint64(bytes, self.Mtime)

	bytes = binary.BigEndian.AppendUint32(bytes, self.Dev)
	bytes = binary.BigEndian.AppendUint32(bytes, self.Ino)
	bytes = binary.BigEndian.AppendUint32(bytes, self.ModeType<<12|self.ModePerms)
	bytes = binary.BigEndian.AppendUint32(bytes, self.Uid)
	bytes = binary.BigEndian.AppendUint32(bytes, self.Gid)
	bytes = binary.BigEndian.AppendUint32(bytes, self.Fsize)

	bytes = append(bytes, []byte(self.Sha)...)

	bytes = binary.BigEndian.AppendUint16(bytes, uint16(len(self.FullPathName)))
	bytes = append(bytes, []byte(self.FullPathName)...)

	return bytes
}

func Deserialize(reader io.Reader) (*IndexObject, error) {
	allBytes, err := ioutil.ReadAll(reader)

	if err != nil {
		return nil, nil
	}
	if len(allBytes) == 0 {
		return &IndexObject{Version: 0, Entries: make(map[string]IndexEntry)}, nil
	}

	version := binary.BigEndian.Uint32(allBytes[:4])
	count := binary.BigEndian.Uint32(allBytes[4:8])

	entries := make(map[string]IndexEntry)
	content := allBytes[8:]
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
	mode := binary.BigEndian.Uint32(content[offset+24 : offset+28])
	modeType := mode >> 12
	modePerms := mode & 0x1FF
	uid := binary.BigEndian.Uint32(content[offset+28 : offset+32])
	gid := binary.BigEndian.Uint32(content[offset+32 : offset+36])
	fsize := binary.BigEndian.Uint32(content[offset+36 : offset+40])
	sha := string(content[offset+40 : offset+80])
	nameLength := binary.BigEndian.Uint16(content[offset+80 : offset+82])
	name := string(content[offset+82 : offset+82+int(nameLength)])

	offset += 82 + int(nameLength)

	return IndexEntry{
		Ctime:        ctime,
		Mtime:        mtime,
		Dev:          device,
		Ino:          ino,
		ModeType:     modeType,
		ModePerms:    modePerms,
		Uid:          uid,
		Gid:          gid,
		Fsize:        fsize,
		Sha:          sha,
		FullPathName: name,
	}, offset
}
