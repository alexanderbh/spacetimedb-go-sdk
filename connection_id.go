package spacetimedb

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// ConnectionId is a unique identifier for a client connected to a database.
type ConnectionId struct {
	data *big.Int
}

// NewConnectionId creates a new ConnectionID with the given big.Int data.
// The provided data is used directly.
func NewConnectionId(data *big.Int) *ConnectionId {
	if data == nil {
		// Ensure data is not nil, similar to how bigint in TS is non-nullable.
		// Or, document that data must not be nil.
		// For now, let's initialize to 0 if nil is passed, though this might hide issues.
		// A better approach might be to panic or return an error if data is nil.
		// However, internal calls like RandomConnectionID ensure data is non-nil.
		return &ConnectionId{data: big.NewInt(0)}
	}
	return &ConnectionId{data: data}
}

// IsZero checks if the ConnectionID is zero.
func (cid *ConnectionId) IsZero() bool {
	if cid == nil || cid.data == nil {
		return true // Consider a nil ConnectionID or nil data as effectively zero or an invalid state.
	}
	return cid.data.Cmp(big.NewInt(0)) == 0
}

// NullIfZero returns nil if the ConnectionID is zero, otherwise returns the ConnectionID.
func NullIfZero(addr *ConnectionId) *ConnectionId {
	if addr == nil || addr.IsZero() {
		return nil
	}
	return addr
}

// randomPseudoByte generates a random integer in [0, 254], mimicking Math.floor(Math.random() * 0xff).
func randomPseudoByte() (uint8, error) {
	// big.NewInt(255) means the range for rand.Int is [0, 254].
	n, err := rand.Int(rand.Reader, big.NewInt(255))
	if err != nil {
		return 0, err
	}
	return uint8(n.Uint64()), nil // n.Uint64() is safe as n is small.
}

// RandomConnectionId creates a new random ConnectionID.
// It replicates the TypeScript logic of building a 128-bit number
// from 16 "bytes", each in the range [0, 254].
func RandomConnectionId() (*ConnectionId, error) {
	pseudoBytes := make([]byte, 16)
	for i := 0; i < 16; i++ {
		pb, err := randomPseudoByte()
		if err != nil {
			return nil, fmt.Errorf("failed to generate pseudo-byte for ConnectionID: %w", err)
		}
		pseudoBytes[i] = pb
	}
	// new(big.Int).SetBytes interprets pseudoBytes as a big-endian unsigned integer.
	data := new(big.Int).SetBytes(pseudoBytes)
	return NewConnectionId(data), nil
}

// IsEqual compares two ConnectionIDs for equality.
func (cid *ConnectionId) IsEqual(other *ConnectionId) bool {
	if cid == other { // Handles both being nil
		return true
	}
	if cid == nil || other == nil || cid.data == nil || other.data == nil {
		return false // One is nil, or internal data is nil
	}
	return cid.data.Cmp(other.data) == 0
}

// ToHexString converts the ConnectionID to a hexadecimal string.
// It relies on U128ToHexString from the utils package.
func (cid *ConnectionId) ToHexString() (string, error) {
	if cid == nil || cid.data == nil {
		return "", fmt.Errorf("cannot convert nil ConnectionID or ConnectionID with nil data to hex string")
	}
	return U128ToHexString(cid.data)
}

// ToUint8Array converts the ConnectionID to a byte array.
// It relies on U128ToUint8Array from the utils package.
func (cid *ConnectionId) ToUint8Array() ([]byte, error) {
	if cid == nil || cid.data == nil {
		return nil, fmt.Errorf("cannot convert nil ConnectionID or ConnectionID with nil data to byte array")
	}
	return U128ToUint8Array(cid.data)
}

// ConnectionIDFromString parses a ConnectionID from a hexadecimal string.
// It relies on HexStringToU128 from the utils package.
func ConnectionIDFromString(str string) (*ConnectionId, error) {
	data, err := HexStringToU128(str)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ConnectionID from string: %w", err)
	}
	if data == nil {
		// HexStringToU128 should ideally not return (nil, nil)
		return nil, fmt.Errorf("HexStringToU128 returned nil data without error for string: %s", str)
	}
	return NewConnectionId(data), nil
}

// ConnectionIDFromStringOrNull parses a ConnectionID from a hexadecimal string,
// returning (nil, nil) if the parsed ID is zero.
func ConnectionIDFromStringOrNull(str string) (*ConnectionId, error) {
	cid, err := ConnectionIDFromString(str)
	if err != nil {
		return nil, err // Error during parsing
	}
	// ConnectionIDFromString should ensure cid and cid.data are not nil if err is nil.
	if cid.IsZero() {
		return nil, nil // Parsed to zero, return (nil, nil)
	}
	return cid, nil
}

// GetData returns the internal *big.Int data.
// This is provided if direct access to the *big.Int is needed,
// similar to accessing the `data` property in the TypeScript version.
func (cid *ConnectionId) GetData() *big.Int {
	if cid == nil {
		return nil
	}
	return cid.data
}

// Serialize serializes the ConnectionID to a BinaryWriter.
func (cid *ConnectionId) Serialize(writer *BinaryWriter) error {
	panic("ConnectionId.Serialize: not implemented")
}

// Deserialize deserializes a ConnectionID from a BinaryReader.
func (cid *ConnectionId) Deserialize(reader *BinaryReader) error {
	data := reader.ReadU128()
	cid.data = data
	return nil
}
