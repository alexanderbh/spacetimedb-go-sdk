package spacetimedb

import (
	"encoding/binary"
	"math"
	"math/big"
)

// BinaryReader helps to read binary data from a byte slice.
type BinaryReader struct {
	buffer []byte
	offset int
}

// NewBinaryReader creates a new BinaryReader with the given input byte slice.
// []byte is a slice so it will not copy the actual data, just the reference to it.
func NewBinaryReader(input []byte) *BinaryReader {
	return &BinaryReader{
		buffer: input,
		offset: 0,
	}
}

// Offset returns the current reading offset.
func (br *BinaryReader) Offset() int {
	return br.offset
}

func (br *BinaryReader) checkBounds(needed int) {
	if br.offset+needed > len(br.buffer) {
		panic("BinaryReader: read out of bounds")
	}
}

// ReadUInt8Array reads a U32 for length, then that many bytes.
func (br *BinaryReader) ReadUInt8Array() []byte {
	length := br.ReadU32()
	br.checkBounds(int(length))
	value := make([]byte, length)
	copy(value, br.buffer[br.offset:br.offset+int(length)])
	br.offset += int(length)
	return value
}

// ReadBool reads a single byte as a boolean (non-zero is true).
func (br *BinaryReader) ReadBool() bool {
	br.checkBounds(1)
	value := br.buffer[br.offset]
	br.offset += 1
	return value != 0
}

// ReadByte reads a single byte.
func (br *BinaryReader) ReadByte() byte {
	br.checkBounds(1)
	value := br.buffer[br.offset]
	br.offset += 1
	return value
}

// ReadBytes reads a specified number of bytes.
func (br *BinaryReader) ReadBytes(length int) []byte {
	br.checkBounds(length)
	value := make([]byte, length)
	copy(value, br.buffer[br.offset:br.offset+length])
	br.offset += length
	return value
}

// ReadI8 reads a signed 8-bit integer.
func (br *BinaryReader) ReadI8() int8 {
	br.checkBounds(1)
	value := int8(br.buffer[br.offset])
	br.offset += 1
	return value
}

// ReadU8 reads an unsigned 8-bit integer.
func (br *BinaryReader) ReadU8() uint8 {
	br.checkBounds(1)
	value := br.buffer[br.offset]
	br.offset += 1
	return value
}

// ReadI16 reads a signed 16-bit integer (little-endian).
func (br *BinaryReader) ReadI16() int16 {
	br.checkBounds(2)
	value := binary.LittleEndian.Uint16(br.buffer[br.offset : br.offset+2])
	br.offset += 2
	return int16(value)
}

// ReadU16 reads an unsigned 16-bit integer (little-endian).
func (br *BinaryReader) ReadU16() uint16 {
	br.checkBounds(2)
	value := binary.LittleEndian.Uint16(br.buffer[br.offset : br.offset+2])
	br.offset += 2
	return value
}

// ReadI32 reads a signed 32-bit integer (little-endian).
func (br *BinaryReader) ReadI32() int32 {
	br.checkBounds(4)
	value := binary.LittleEndian.Uint32(br.buffer[br.offset : br.offset+4])
	br.offset += 4
	return int32(value)
}

// ReadU32 reads an unsigned 32-bit integer (little-endian).
func (br *BinaryReader) ReadU32() uint32 {
	br.checkBounds(4)
	value := binary.LittleEndian.Uint32(br.buffer[br.offset : br.offset+4])
	br.offset += 4
	return value
}

// ReadI64 reads a signed 64-bit integer (little-endian).
func (br *BinaryReader) ReadI64() int64 {
	br.checkBounds(8)
	value := binary.LittleEndian.Uint64(br.buffer[br.offset : br.offset+8])
	br.offset += 8
	return int64(value)
}

// ReadU64 reads an unsigned 64-bit integer (little-endian).
func (br *BinaryReader) ReadU64() uint64 {
	br.checkBounds(8)
	value := binary.LittleEndian.Uint64(br.buffer[br.offset : br.offset+8])
	br.offset += 8
	return value
}

// ReadU128 reads an unsigned 128-bit integer (little-endian).
func (br *BinaryReader) ReadU128() *big.Int {
	br.checkBounds(16)
	lowerPart := binary.LittleEndian.Uint64(br.buffer[br.offset : br.offset+8])
	upperPart := binary.LittleEndian.Uint64(br.buffer[br.offset+8 : br.offset+16])
	br.offset += 16

	result := new(big.Int).SetUint64(upperPart)
	result.Lsh(result, 64)
	result.Add(result, new(big.Int).SetUint64(lowerPart))
	return result
}

// ReadI128 reads a signed 128-bit integer (little-endian).
// Matches TypeScript logic: lower 64 bits as unsigned, upper 64 bits as signed.
func (br *BinaryReader) ReadI128() *big.Int {
	br.checkBounds(16)
	lowerPart := binary.LittleEndian.Uint64(br.buffer[br.offset : br.offset+8])
	upperPartSigned := int64(binary.LittleEndian.Uint64(br.buffer[br.offset+8 : br.offset+16]))
	br.offset += 16

	upperBigInt := new(big.Int).SetInt64(upperPartSigned)
	upperBigInt.Lsh(upperBigInt, 64)

	lowerBigInt := new(big.Int).SetUint64(lowerPart)

	result := new(big.Int).Add(upperBigInt, lowerBigInt)
	return result
}

// ReadU256 reads an unsigned 256-bit integer (little-endian).
func (br *BinaryReader) ReadU256() *big.Int {
	br.checkBounds(32)
	p0 := binary.LittleEndian.Uint64(br.buffer[br.offset : br.offset+8])
	p1 := binary.LittleEndian.Uint64(br.buffer[br.offset+8 : br.offset+16])
	p2 := binary.LittleEndian.Uint64(br.buffer[br.offset+16 : br.offset+24])
	p3 := binary.LittleEndian.Uint64(br.buffer[br.offset+24 : br.offset+32])
	br.offset += 32

	valP3 := new(big.Int).SetUint64(p3)
	valP3.Lsh(valP3, 192)

	valP2 := new(big.Int).SetUint64(p2)
	valP2.Lsh(valP2, 128)

	valP1 := new(big.Int).SetUint64(p1)
	valP1.Lsh(valP1, 64)

	valP0 := new(big.Int).SetUint64(p0)

	result := new(big.Int).Add(valP3, valP2)
	result.Add(result, valP1)
	result.Add(result, valP0)
	return result
}

// ReadI256 reads a signed 256-bit integer (little-endian).
// Matches TypeScript logic: p0, p1, p2 as unsigned 64-bit, p3 as signed 64-bit.
func (br *BinaryReader) ReadI256() *big.Int {
	br.checkBounds(32)
	p0 := binary.LittleEndian.Uint64(br.buffer[br.offset : br.offset+8])
	p1 := binary.LittleEndian.Uint64(br.buffer[br.offset+8 : br.offset+16])
	p2 := binary.LittleEndian.Uint64(br.buffer[br.offset+16 : br.offset+24])
	p3Signed := int64(binary.LittleEndian.Uint64(br.buffer[br.offset+24 : br.offset+32]))
	br.offset += 32

	valP3 := new(big.Int).SetInt64(p3Signed)
	valP3.Lsh(valP3, 192)

	valP2 := new(big.Int).SetUint64(p2)
	valP2.Lsh(valP2, 128)

	valP1 := new(big.Int).SetUint64(p1)
	valP1.Lsh(valP1, 64)

	valP0 := new(big.Int).SetUint64(p0)

	result := new(big.Int).Add(valP3, valP2)
	result.Add(result, valP1)
	result.Add(result, valP0)
	return result
}

// ReadF32 reads a 32-bit float (little-endian).
func (br *BinaryReader) ReadF32() float32 {
	br.checkBounds(4)
	bits := binary.LittleEndian.Uint32(br.buffer[br.offset : br.offset+4])
	br.offset += 4
	return math.Float32frombits(bits)
}

// ReadF64 reads a 64-bit float (little-endian).
func (br *BinaryReader) ReadF64() float64 {
	br.checkBounds(8)
	bits := binary.LittleEndian.Uint64(br.buffer[br.offset : br.offset+8])
	br.offset += 8
	return math.Float64frombits(bits)
}

// ReadString reads a U32 for length, then that many bytes as a UTF-8 string.
func (br *BinaryReader) ReadString() string {
	length := br.ReadU32()
	br.checkBounds(int(length))
	value := string(br.buffer[br.offset : br.offset+int(length)])
	br.offset += int(length)
	return value
}
