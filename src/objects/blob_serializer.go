package objects

type BlobObject struct {
	Object
}

func CreateBlobObject(data []byte) *BlobObject {
	return &BlobObject{
		Object{Type: BLOB, Data: data, Length: len(data)},
	}
}

func deserializeBlobObject(commonObject *Object, toDeserialize []byte) (*BlobObject, error) {
	blobObject := &BlobObject{Object: *commonObject}
	blobObject.Data = toDeserialize
	return blobObject, nil
}

func (c *BlobObject) Serialize() []byte {
	return append(c.serializeHeader(), c.Data...)
}
