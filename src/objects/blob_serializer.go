package objects

func serializeBlob(object *Object) []byte {
	return object.Data
}
