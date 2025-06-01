package spacetimedb

import "fmt"

type RowSizeHint struct {
	RowSizeHint any
}

type RowSizeHintFixedSize struct {
	FixedSize uint16
}

func NewRowSizeHintFixedSize(fixedSize uint16) *RowSizeHintFixedSize {
	return &RowSizeHintFixedSize{FixedSize: fixedSize}
}

type RowSizeHintRowOffsets struct {
	RowOffsets []uint64
}

func NewRowSizeHintRowOffsets(rowOffsets []uint64) *RowSizeHintRowOffsets {
	return &RowSizeHintRowOffsets{RowOffsets: rowOffsets}
}

func (it *RowSizeHint) Deserialize(reader *BinaryReader) error {
	unionType := reader.ReadU8()
	switch unionType {
	case 0x00:
		it.RowSizeHint = NewRowSizeHintFixedSize(reader.ReadU16())
	case 0x01:
		it.RowSizeHint = NewRowSizeHintRowOffsets(ReadArray(reader, reader.ReadU64))
	default:
		return fmt.Errorf("RowSizeHint.Deserialize: unknown union type 0x%02x", unionType)
	}
	return nil
}

func (it *RowSizeHint) String() string {
	switch hint := it.RowSizeHint.(type) {
	case *RowSizeHintFixedSize:
		return fmt.Sprintf("FixedSize=%d", hint.FixedSize)
	case *RowSizeHintRowOffsets:
		return fmt.Sprintf("RowOffsets=%+v", hint.RowOffsets)
	default:
		return fmt.Sprintf("Unknown RowSizeHint type: %T", it.RowSizeHint)
	}
}
