package spacetimedb

type TableCache[T any] struct {
	Rows map[string]T
}

type Table interface {
	Insert(reader *BinaryReader) error
	Delete(reader *BinaryReader) error
}

type TableNameMap = map[string]Table
