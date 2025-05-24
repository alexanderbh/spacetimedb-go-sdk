package spacetimedb

type ParseableType[T any] interface {
	Deserialize(reader *BinaryReader) T
}

func ParseValue[T any](ty ParseableType[T], src []byte) T {
	reader := NewBinaryReader(src)
	return ty.Deserialize(reader)
}
