package spacetimedb

type TableUpdate struct {
	TableID   uint32
	TableName string
	NumRows   uint64
	Updates   []*QueryUpdate
}

func (it *TableUpdate) Deserialize(reader *BinaryReader) error {

	it.TableID = reader.ReadU32()
	it.TableName = reader.ReadString()

	it.NumRows = reader.ReadU64()

	updates := ReadArray(reader, func() *QueryUpdate {
		update := &CompressableQueryUpdate{}
		update.Deserialize(reader)
		return update.Update
	})
	it.Updates = updates

	return nil
}

func (it *TableUpdate) String() string {
	result := "TableUpdate:\n"
	result += "    TableID: " + U32ToHexString(it.TableID) + "\n"
	result += "    TableName: " + it.TableName + "\n"
	result += "    NumRows: " + U64ToHexString(it.NumRows) + "\n"
	for _, update := range it.Updates {
		result += update.String() + "\n"
	}
	return result
}
