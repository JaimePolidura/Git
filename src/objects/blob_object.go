package objects

type BlobObject struct {
	Data []byte
}

func CreateBlobObject(data []byte) *Object {
	return &Object{
		SerializableGitObject: BlobObject{Data: data},
		Type:                  BLOB,
	}
}

func deserializeBlobObject(toDeserialize []byte) (BlobObject, error) {
	blobObject := BlobObject{Data: toDeserialize}
	return blobObject, nil
}

func (b BlobObject) Serialize() []byte {
	return b.Data
}
