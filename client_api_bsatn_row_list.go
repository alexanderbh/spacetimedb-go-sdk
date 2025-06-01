package spacetimedb

type BsatnRowList struct {
	SizeHint *RowSizeHint
	RowsData []byte
}

func (it *BsatnRowList) Deserialize(reader *BinaryReader) error {

	it.SizeHint = &RowSizeHint{}
	it.SizeHint.Deserialize(reader)

	it.RowsData = reader.ReadUInt8Array()

	return nil
}

func (it *BsatnRowList) String() string {
	result := "BsatnRowList:\n"
	if it.SizeHint != nil {
		result += "  SizeHint: " + it.SizeHint.String() + "\n"
	} else {
		result += "  SizeHint: <nil>\n"
	}
	if it.RowsData != nil {
		result += "  RowsData: <byte array of length " + U64ToHexString(uint64(len(it.RowsData))) + ">\n"
	} else {
		result += "  RowsData: <nil>\n"
	}
	return result
}
