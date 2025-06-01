package spacetimedb

import (
	"encoding/binary"
	"math"
	"math/big"
)

const defaultInitialBufferSize = 1024

// BinaryWriter helps to write binary data into a byte slice.
type BinaryWriter struct {
	buffer []byte
	offset int
}

// NewBinaryWriter creates a new BinaryWriter with an initial buffer size.
func NewBinaryWriter(initialSize ...int) *BinaryWriter {
	size := defaultInitialBufferSize
	if len(initialSize) > 0 && initialSize[0] > 0 {
		size = initialSize[0]
	}
	return &BinaryWriter{
		buffer: make([]byte, size),
		offset: 0,
	}
}

func (bw *BinaryWriter) expandBuffer(additionalCapacity int) {
	minCapacity := bw.offset + additionalCapacity
	if minCapacity <= len(bw.buffer) {
		return
	}
	newCapacity := len(bw.buffer) * 2
	if newCapacity < minCapacity {
		newCapacity = minCapacity
	}
	newBuffer := make([]byte, newCapacity)
	copy(newBuffer, bw.buffer[:bw.offset])
	bw.buffer = newBuffer
}

// GetBuffer returns the written bytes as a slice.
func (bw *BinaryWriter) GetBuffer() []byte {
	return bw.buffer[:bw.offset]
}

// WriteUInt8Array writes a U32 for length, then the byte array.
func (bw *BinaryWriter) WriteUInt8Array(value []byte) {
	length := uint32(len(value))
	bw.WriteU32(length)
	bw.expandBuffer(len(value))
	copy(bw.buffer[bw.offset:], value)
	bw.offset += len(value)
}

// WriteBool writes a boolean as a single byte (1 for true, 0 for false).
func (bw *BinaryWriter) WriteBool(value bool) {
	bw.expandBuffer(1)
	if value {
		bw.buffer[bw.offset] = 1
	} else {
		bw.buffer[bw.offset] = 0
	}
	bw.offset += 1
}

// WriteByte writes a single byte.
func (bw *BinaryWriter) WriteByte(value byte) {
	bw.expandBuffer(1)
	bw.buffer[bw.offset] = value
	bw.offset += 1
}

// WriteI8 writes a signed 8-bit integer.
func (bw *BinaryWriter) WriteI8(value int8) {
	bw.expandBuffer(1)
	bw.buffer[bw.offset] = byte(value)
	bw.offset += 1
}

// WriteU8 writes an unsigned 8-bit integer.
func (bw *BinaryWriter) WriteU8(value uint8) {
	bw.expandBuffer(1)
	bw.buffer[bw.offset] = value
	bw.offset += 1
}

// WriteI16 writes a signed 16-bit integer (little-endian).
func (bw *BinaryWriter) WriteI16(value int16) {
	bw.expandBuffer(2)
	binary.LittleEndian.PutUint16(bw.buffer[bw.offset:], uint16(value))
	bw.offset += 2
}

// WriteU16 writes an unsigned 16-bit integer (little-endian).
func (bw *BinaryWriter) WriteU16(value uint16) {
	bw.expandBuffer(2)
	binary.LittleEndian.PutUint16(bw.buffer[bw.offset:], value)
	bw.offset += 2
}

// WriteI32 writes a signed 32-bit integer (little-endian).
func (bw *BinaryWriter) WriteI32(value int32) {
	bw.expandBuffer(4)
	binary.LittleEndian.PutUint32(bw.buffer[bw.offset:], uint32(value))
	bw.offset += 4
}

// WriteU32 writes an unsigned 32-bit integer (little-endian).
func (bw *BinaryWriter) WriteU32(value uint32) {
	bw.expandBuffer(4)
	binary.LittleEndian.PutUint32(bw.buffer[bw.offset:], value)
	bw.offset += 4
}

// WriteI64 writes a signed 64-bit integer (little-endian).
func (bw *BinaryWriter) WriteI64(value int64) {
	bw.expandBuffer(8)
	binary.LittleEndian.PutUint64(bw.buffer[bw.offset:], uint64(value))
	bw.offset += 8
}

// WriteU64 writes an unsigned 64-bit integer (little-endian).
func (bw *BinaryWriter) WriteU64(value uint64) {
	bw.expandBuffer(8)
	binary.LittleEndian.PutUint64(bw.buffer[bw.offset:], value)
	bw.offset += 8
}

// WriteU128 writes an unsigned 128-bit integer (little-endian).
func (bw *BinaryWriter) WriteU128(value *big.Int) {
	bw.expandBuffer(16)
	bytes := make([]byte, 16)
	// Ensure the big.Int is treated as unsigned and get its bytes in big-endian.
	// Then reverse to little-endian if necessary, or fill parts.
	valBytes := value.Bytes()
	// If the number is smaller than 16 bytes, it will be padded with leading zeros (big-endian).
	// We need to place these bytes at the correct position for little-endian.
	start := 16 - len(valBytes)
	for i := 0; i < len(valBytes); i++ {
		bytes[start+i] = valBytes[i] // Place in big-endian order first
	}
	// Reverse for little-endian (or handle parts directly)
	// The TypeScript code writes lowerPart then upperPart.
	// lowerPart = value & 0xFFFFFFFFFFFFFFFF
	// upperPart = value >> 64
	var lowerPart, upperPart big.Int
	mask := new(big.Int).SetBytes([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})

	lowerPart.And(value, mask)
	upperPart.Rsh(value, 64)

	binary.LittleEndian.PutUint64(bw.buffer[bw.offset:], lowerPart.Uint64())
	binary.LittleEndian.PutUint64(bw.buffer[bw.offset+8:], upperPart.Uint64())
	bw.offset += 16
}

// WriteI128 writes a signed 128-bit integer (little-endian).
// Matches TypeScript logic: lower 64 bits as unsigned, upper 64 bits as signed.
func (bw *BinaryWriter) WriteI128(value *big.Int) {
	bw.expandBuffer(16)
	// TypeScript: lowerPart = value & BigInt('0xFFFFFFFFFFFFFFFF'); (unsigned)
	//             upperPart = value >> BigInt(64); (signed)
	var lowerPartUnsigned uint64
	var upperPartSigned int64

	// For lower part (unsigned 64 bits)
	bytes := value.Bytes() // Big-endian bytes
	tempLower := make([]byte, 8)
	if len(bytes) >= 8 {
		copy(tempLower, bytes[len(bytes)-8:])
	} else {
		copy(tempLower[8-len(bytes):], bytes)
	}
	lowerPartUnsigned = binary.BigEndian.Uint64(tempLower)

	// For upper part (signed 64 bits)
	tempUpperVal := new(big.Int).Rsh(value, 64)
	upperPartSigned = tempUpperVal.Int64() // This handles sign correctly

	binary.LittleEndian.PutUint64(bw.buffer[bw.offset:], lowerPartUnsigned)
	binary.LittleEndian.PutUint64(bw.buffer[bw.offset+8:], uint64(upperPartSigned))
	bw.offset += 16
}

// WriteU256 writes an unsigned 256-bit integer (little-endian).
func (bw *BinaryWriter) WriteU256(value *big.Int) {
	bw.expandBuffer(32)
	bytes := make([]byte, 32) // Ensure 32 bytes for the number
	valBytes := value.Bytes() // These are big-endian

	// Copy valBytes into the end of our 32-byte slice to handle numbers smaller than 256 bits.
	// This effectively zero-pads the most significant bytes if valBytes is shorter.
	offset := 32 - len(valBytes)
	copy(bytes[offset:], valBytes)

	// Write in little-endian order (p0, p1, p2, p3)
	binary.LittleEndian.PutUint64(bw.buffer[bw.offset+0:], binary.BigEndian.Uint64(bytes[24:32])) // p0
	binary.LittleEndian.PutUint64(bw.buffer[bw.offset+8:], binary.BigEndian.Uint64(bytes[16:24])) // p1
	binary.LittleEndian.PutUint64(bw.buffer[bw.offset+16:], binary.BigEndian.Uint64(bytes[8:16])) // p2
	binary.LittleEndian.PutUint64(bw.buffer[bw.offset+24:], binary.BigEndian.Uint64(bytes[0:8]))  // p3
	bw.offset += 32
}

// WriteI256 writes a signed 256-bit integer (little-endian).
// Matches TypeScript: p0, p1, p2 as unsigned, p3 as signed.
func (bw *BinaryWriter) WriteI256(value *big.Int) {
	bw.expandBuffer(32)

	mask := new(big.Int).SetBytes([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})

	var p0, p1, p2 big.Int
	p0.And(value, mask)
	p1.Rsh(value, 64)
	p1.And(&p1, mask)
	p2.Rsh(value, 128)
	p2.And(&p2, mask)

	p3 := new(big.Int).Rsh(value, 192)

	binary.LittleEndian.PutUint64(bw.buffer[bw.offset+0:], p0.Uint64())
	binary.LittleEndian.PutUint64(bw.buffer[bw.offset+8:], p1.Uint64())
	binary.LittleEndian.PutUint64(bw.buffer[bw.offset+16:], p2.Uint64())
	binary.LittleEndian.PutUint64(bw.buffer[bw.offset+24:], uint64(p3.Int64())) // p3 is signed

	bw.offset += 32
}

// WriteF32 writes a 32-bit float (little-endian).
func (bw *BinaryWriter) WriteF32(value float32) {
	bw.expandBuffer(4)
	binary.LittleEndian.PutUint32(bw.buffer[bw.offset:], math.Float32bits(value))
	bw.offset += 4
}

// WriteF64 writes a 64-bit float (little-endian).
func (bw *BinaryWriter) WriteF64(value float64) {
	bw.expandBuffer(8)
	binary.LittleEndian.PutUint64(bw.buffer[bw.offset:], math.Float64bits(value))
	bw.offset += 8
}

// WriteString writes a U32 for length, then the UTF-8 encoded string.
func (bw *BinaryWriter) WriteString(value string) {
	encodedString := []byte(value)
	length := uint32(len(encodedString))
	bw.WriteU32(length)
	bw.expandBuffer(len(encodedString))
	copy(bw.buffer[bw.offset:], encodedString)
	bw.offset += len(encodedString)
}

func WriteArray[T any](bw *BinaryWriter, values []T, writeFunc func(*BinaryWriter, T)) {
	bw.WriteU32(uint32(len(values)))
	for _, value := range values {
		writeFunc(bw, value)
	}
}
