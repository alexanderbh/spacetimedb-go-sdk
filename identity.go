package spacetimedb

import (
	"fmt"
	"math/big"
)

// Identity is a unique identifier for a user connected to a database.
type Identity struct {
	data *big.Int
}

// NewIdentity creates a new Identity.
// 'data' can be a hexadecimal string or a *big.Int.
func NewIdentity(data interface{}) (*Identity, error) {
	switch v := data.(type) {
	case string:
		val, err := HexStringToU256(v)
		if err != nil {
			return nil, fmt.Errorf("invalid hex string for Identity: %w", err)
		}
		return &Identity{data: val}, nil
	case *big.Int:
		return &Identity{data: v}, nil
	default:
		return nil, fmt.Errorf("unsupported type for Identity")
	}
}

// Data returns the underlying bigint value.
func (id *Identity) Data() *big.Int {
	return id.data
}

// IsEqual compares two identities for equality.
func (id *Identity) IsEqual(other *Identity) bool {
	return id.ToHexString() == other.ToHexString()
}

// ToHexString prints the identity as a hexadecimal string.
func (id *Identity) ToHexString() string {
	s, _ := U256ToHexString(id.data)
	return s
}

// ToUint8Array converts the address to a byte array.
func (id *Identity) ToUint8Array() []byte {
	arr, _ := U256ToUint8Array(id.data)
	return arr
}

// FromString parses an Identity from a hexadecimal string.
func FromString(str string) (*Identity, error) {
	return NewIdentity(str)
}
