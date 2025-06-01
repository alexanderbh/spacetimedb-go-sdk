package spacetimedb

type TableCache[T any] struct {
	Rows map[string]T
}

type Table interface {
	DeserializeRow(reader *BinaryReader) (any, error)
	Insert(row any) error
}

type TableNameMap = map[string]Table
