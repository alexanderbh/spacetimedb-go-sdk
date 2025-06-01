package spacetimedb

import "fmt"

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
	result := ""
	if it.SizeHint != nil {
		result += "      " + it.SizeHint.String() + "\n"
	}
	if it.RowsData != nil {
		result += "      <byte array of length " + fmt.Sprint(len(it.RowsData)) + ">\n"
	}
	return result
}
