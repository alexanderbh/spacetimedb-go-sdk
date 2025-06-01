package spacetimedb

type DatabaseUpdate struct {
	Tables []*TableUpdate
}

func (it *DatabaseUpdate) Deserialize(reader *BinaryReader) error {

	tables := ReadArray(reader, func() *TableUpdate {
		table := &TableUpdate{}
		table.Deserialize(reader)
		return table
	})
	it.Tables = tables

	return nil
}

func (it *DatabaseUpdate) String() string {
	result := "DatabaseUpdate:\n"
	for _, table := range it.Tables {
		result += "  " + table.String() + "\n"
	}
	return result
}
