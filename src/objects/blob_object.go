package objects

type BlobObject struct {
	Object
}

func CreateBlobObject(data []byte) *BlobObject {
	return &BlobObject{
		Object{Type: BLOB, Data: data, Length: len(data)},
	}
}

func deserializeBlobObject(commonObject *Object, toDeserialize []byte) (BlobObject, error) {
	blobObject := BlobObject{Object: *commonObject}
	blobObject.Object.Data = toDeserialize
	return blobObject, nil
}

func (b BlobObject) Type() ObjectType {
	return b.Object.Type
}

func (b BlobObject) Length() int {
	return b.Object.Length
}

func (b BlobObject) Data() []byte {
	return b.Object.Data
}

func (b BlobObject) Serialize() []byte {
	return append(b.serializeHeader(), b.Object.Data...)
}
