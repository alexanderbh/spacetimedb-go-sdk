package spacetimedb

type ReducerCallInfo struct {
	ReducerName string
	ReducerID   uint32
	Args        []byte
	RequestID   uint32
}

func (it *ReducerCallInfo) Deserialize(reader *BinaryReader) error {

	it.ReducerName = reader.ReadString()
	it.ReducerID = reader.ReadU32()
	it.Args = reader.ReadUInt8Array()
	it.RequestID = reader.ReadU32()

	return nil
}
