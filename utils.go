package spacetimedb

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
	"regexp"
	"strings"
	"unicode"
)

// ToPascalCase converts a string to PascalCase.
// It replicates the behavior of the TypeScript version:
// - Replaces "-a" or "_a" with "A" (case-insensitive for 'a').
// - Capitalizes the first letter of the result.
func ToPascalCase(s string) string {
	re := regexp.MustCompile(`([-_][a-zA-Z])`)
	str := re.ReplaceAllStringFunc(s, func(match string) string {
		// match is the full matched text e.g., "-a" or "_B"
		// We want to uppercase the character *after* the hyphen or underscore.
		if len(match) > 1 {
			return strings.ToUpper(string(match[1]))
		}
		return "" // Should not happen with the regex
	})

	if len(str) == 0 {
		return ""
	}
	runes := []rune(str)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// DeepEqual is a wrapper around reflect.DeepEqual.
// The original TypeScript version compares objects by keys and recursively.
// reflect.DeepEqual provides a robust and comprehensive deep comparison for Go types,
// which is the idiomatic way to achieve this in Go.
func DeepEqual(obj1, obj2 interface{}) bool {
	return reflect.DeepEqual(obj1, obj2)
}

// Uint8ArrayToHexString converts a byte array to a hex string.
// IMPORTANT: It reverses the byte array before encoding, matching the TS behavior.
func Uint8ArrayToHexString(array []byte) string {
	if array == nil {
		return ""
	}
	// Create a new slice for the reversed array to avoid modifying the input slice if it's used elsewhere.
	reversedArray := make([]byte, len(array))
	for i := 0; i < len(array); i++ {
		reversedArray[i] = array[len(array)-1-i]
	}
	return hex.EncodeToString(reversedArray)
}

// HexStringToUint8Array converts a hex string to a byte array.
// IMPORTANT:
//   - It expects the hex string to represent exactly 32 bytes.
//     If not, or if decoding fails, it returns an empty byte array.
//   - It reverses the resulting byte array.
//
// This matches the specific behavior of the TypeScript version.
func HexStringToUint8Array(str string) []byte {
	if strings.HasPrefix(str, "0x") {
		str = str[2:]
	}

	data, err := hex.DecodeString(str)
	if err != nil {
		return []byte{} // Return empty slice on decoding error
	}

	if len(data) != 32 {
		return []byte{} // Return empty slice if not 32 bytes
	}

	// Reverse the data
	reversedData := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		reversedData[i] = data[len(data)-1-i]
	}
	return reversedData
}

// Uint8ArrayToU128 converts a 16-byte array to a *big.Int.
// It relies on NewBinaryReader and ReadU128 methods (assumed to be in this package).
func Uint8ArrayToU128(array []byte) (*big.Int, error) {
	if len(array) != 16 {
		return nil, fmt.Errorf("byte array is not 16 bytes long: got %d bytes, expected 16", len(array))
	}
	reader := NewBinaryReader(array) // Assumes NewBinaryReader(array []byte) *BinaryReader
	val := reader.ReadU128()         // Assumes ReadU128() *big.Int
	return val, nil
}

// Uint8ArrayToU256 converts a 32-byte array to a *big.Int.
// It relies on NewBinaryReader and ReadU256 methods (assumed to be in this package).
func Uint8ArrayToU256(array []byte) (*big.Int, error) {
	if len(array) != 32 {
		return nil, fmt.Errorf("byte array is not 32 bytes long: got %d bytes, expected 32", len(array))
	}
	reader := NewBinaryReader(array)
	val := reader.ReadU256() // Assumes ReadU256() *big.Int
	return val, nil
}

// HexStringToU128 converts a hex string to a *big.Int (128-bit).
// It uses HexStringToUint8Array and Uint8ArrayToU128.
// Note: Due to HexStringToUint8Array's specific 32-byte output requirement (it returns
// an empty slice if the input hex doesn't decode to 32 bytes), this function
// will effectively fail (as Uint8ArrayToU128 expects 16 bytes) if the input hex string 'str'
// does not decode to exactly 32 bytes. This matches the behavior derived from the
// literal translation of the TS functions.
func HexStringToU128(str string) (*big.Int, error) {
	byteArray := HexStringToUint8Array(str)
	return Uint8ArrayToU128(byteArray)
}

// HexStringToU256 converts a hex string to a *big.Int (256-bit).
// It uses HexStringToUint8Array and Uint8ArrayToU256.
// This function works as intended if 'str' is a hex representation of 32 bytes,
// as HexStringToUint8Array will provide a 32-byte array to Uint8ArrayToU256.
func HexStringToU256(str string) (*big.Int, error) {
	byteArray := HexStringToUint8Array(str)
	return Uint8ArrayToU256(byteArray)
}

// U128ToUint8Array converts a *big.Int to a 16-byte array.
// It relies on NewBinaryWriter, WriteU128, and Bytes methods (assumed to be in this package).
func U128ToUint8Array(data *big.Int) ([]byte, error) {
	writer := NewBinaryWriter(16) // Assumes NewBinaryWriter(initialCapacity int) *BinaryWriter
	writer.WriteU128(data)        // Assumes WriteU128(val *big.Int)
	// No explicit error return from WriteU128 in the provided Go code for BinaryWriter
	buffer := writer.GetBuffer() // Assumes GetBuffer() []byte returns the written content

	// The TS version relies on the BinaryWriter to produce a buffer of the correct length (16 bytes).
	// We replicate this reliance. If writer.GetBuffer() could return a buffer of a different length
	// (e.g. larger, with padding, or if WriteU128 doesn't guarantee 16 bytes),
	// this might not be 16 bytes. For an "exact" translation, no extra checks are added here.
	return buffer, nil
}

// U256ToUint8Array converts a *big.Int to a 32-byte array.
// It relies on NewBinaryWriter, WriteU256, and Bytes methods (assumed to be in this package).
func U256ToUint8Array(data *big.Int) ([]byte, error) {
	writer := NewBinaryWriter(32)
	writer.WriteU256(data) // Assumes WriteU256(val *big.Int)
	// No explicit error return from WriteU256
	buffer := writer.GetBuffer()
	return buffer, nil
}

// U128ToHexString converts a *big.Int (128-bit) to its hex string representation.
// The resulting hex string corresponds to a reversed byte array, due to Uint8ArrayToHexString's behavior.
func U128ToHexString(data *big.Int) (string, error) {
	byteArray, err := U128ToUint8Array(data)
	if err != nil {
		return "", fmt.Errorf("failed to convert U128 to byte array for hex string: %w", err)
	}
	return Uint8ArrayToHexString(byteArray), nil
}

// U256ToHexString converts a *big.Int (256-bit) to its hex string representation.
// The resulting hex string corresponds to a reversed byte array, due to Uint8ArrayToHexString's behavior.
func U256ToHexString(data *big.Int) (string, error) {
	byteArray, err := U256ToUint8Array(data)
	if err != nil {
		return "", fmt.Errorf("failed to convert U256 to byte array for hex string: %w", err)
	}
	return Uint8ArrayToHexString(byteArray), nil
}

func U32ToHexString(data uint32) string {
	// Convert uint32 to byte array
	byteArray := make([]byte, 4)
	byteArray[0] = byte(data >> 24)
	byteArray[1] = byte(data >> 16)
	byteArray[2] = byte(data >> 8)
	byteArray[3] = byte(data)

	// Convert to hex string
	return Uint8ArrayToHexString(byteArray)
}

func U64ToHexString(data uint64) string {
	// Convert uint64 to byte array
	byteArray := make([]byte, 8)
	byteArray[0] = byte(data >> 56)
	byteArray[1] = byte(data >> 48)
	byteArray[2] = byte(data >> 40)
	byteArray[3] = byte(data >> 32)
	byteArray[4] = byte(data >> 24)
	byteArray[5] = byte(data >> 16)
	byteArray[6] = byte(data >> 8)
	byteArray[7] = byte(data)

	// Convert to hex string
	return Uint8ArrayToHexString(byteArray)
}
