package spacetimedb

import (
	"fmt"
	"math/big"
	"reflect"
)

// AlgebraicTypeKind defines the various kinds of algebraic types.
type AlgebraicTypeKind int

const (
	// UnitKind represents an empty value, often for variants without data or standalone unit type.
	UnitKind AlgebraicTypeKind = iota
	// BoolKind represents a boolean value.
	BoolKind
	// U8Kind represents an unsigned 8-bit integer.
	U8Kind
	// I8Kind represents a signed 8-bit integer.
	I8Kind
	// U16Kind represents an unsigned 16-bit integer.
	U16Kind
	// I16Kind represents a signed 16-bit integer.
	I16Kind
	// U32Kind represents an unsigned 32-bit integer.
	U32Kind
	// I32Kind represents a signed 32-bit integer.
	I32Kind
	// U64Kind represents an unsigned 64-bit integer.
	U64Kind
	// I64Kind represents a signed 64-bit integer.
	I64Kind
	// F32Kind represents a 32-bit floating-point number.
	F32Kind
	// F64Kind represents a 64-bit floating-point number.
	F64Kind
	// StringKind represents a UTF-8 string.
	StringKind
	// BytesKind represents an array of bytes.
	BytesKind
	// U128Kind represents an unsigned 128-bit integer.
	U128Kind
	// I128Kind represents a signed 128-bit integer.
	I128Kind
	// U256Kind represents an unsigned 256-bit integer.
	U256Kind
	// I256Kind represents a signed 256-bit integer.
	I256Kind

	// Domain-specific Kinds, often wrappers around primitives or specific structures
	// TimestampKind represents a point in time (e.g., uint64 microseconds since epoch).
	TimestampKind
	// TimeDurationKind represents a duration (e.g., int64 microseconds).
	TimeDurationKind
	// ConnectionIdKind represents a connection identifier (typically U128).
	ConnectionIdKind // Already maps to U128Kind in practice via spacetimedb.ConnectionId
	// IdentityKind represents a user identity (typically U256).
	IdentityKind // Already maps to U256Kind in practice via spacetimedb.Identity

	// Structural Kinds
	// SumKind represents a sum type (discriminated union).
	SumKind
	// StructKind represents a product type (record/struct).
	StructKind
	// OptionKind is handled as a special SumKind with "some" and "none" variants.
	// ArrayKind represents an array of elements of the same type, with special handling for byte arrays.
	ArrayKind
	// MapKind represents a map/dictionary.
	MapKind
	// TupleKind represents an ordered, fixed-size collection of elements of potentially different types.
	TupleKind
)

// AlgebraicType represents the definition of a data type.
type AlgebraicType struct {
	Kind AlgebraicTypeKind
	Name string // Optional: For named types, e.g., "MyStruct", "ScheduleAt"

	// Details for structural kinds
	SumDetails struct {
		Variants []*SumTypeVariant
	}
	StructDetails struct {
		Fields []StructField // Fields are ordered
	}
	OptionDetails struct { // Only if OptionKind is made distinct and not just a SumType pattern
		SomeType *AlgebraicType
	}
	ArrayDetails struct {
		ElementType *AlgebraicType
	}
	MapDetails struct {
		KeyType   *AlgebraicType
		ValueType *AlgebraicType
	}
	TupleDetails struct {
		ElementTypes []*AlgebraicType // Elements are ordered
	}
}

// SumTypeVariant defines a variant within a sum type.
type SumTypeVariant struct {
	Name          string
	AlgebraicType *AlgebraicType // The type of the data this variant holds
}

// NewSumTypeVariant creates a new SumTypeVariant.
func NewSumTypeVariant(name string, algebraicType *AlgebraicType) *SumTypeVariant {
	return &SumTypeVariant{Name: name, AlgebraicType: algebraicType}
}

// StructField defines a field within a struct type.
type StructField struct {
	Name string
	Type *AlgebraicType
}

// Factory functions for AlgebraicType instances

// CreateUnitType creates an AlgebraicType for the unit type.
func CreateUnitType() *AlgebraicType { return &AlgebraicType{Kind: UnitKind} }

// CreateBoolType creates an AlgebraicType for the boolean type.
func CreateBoolType() *AlgebraicType { return &AlgebraicType{Kind: BoolKind} }

// CreateU8Type creates an AlgebraicType for the U8 type.
func CreateU8Type() *AlgebraicType { return &AlgebraicType{Kind: U8Kind} }

// CreateI8Type creates an AlgebraicType for the I8 type.
func CreateI8Type() *AlgebraicType { return &AlgebraicType{Kind: I8Kind} }

// CreateU16Type creates an AlgebraicType for the U16 type.
func CreateU16Type() *AlgebraicType { return &AlgebraicType{Kind: U16Kind} }

// CreateI16Type creates an AlgebraicType for the I16 type.
func CreateI16Type() *AlgebraicType { return &AlgebraicType{Kind: I16Kind} }

// CreateU32Type creates an AlgebraicType for the U32 type.
func CreateU32Type() *AlgebraicType { return &AlgebraicType{Kind: U32Kind} }

// CreateI32Type creates an AlgebraicType for the I32 type.
func CreateI32Type() *AlgebraicType { return &AlgebraicType{Kind: I32Kind} }

// CreateU64Type creates an AlgebraicType for the U64 type.
func CreateU64Type() *AlgebraicType { return &AlgebraicType{Kind: U64Kind} }

// CreateI64Type creates an AlgebraicType for the I64 type.
func CreateI64Type() *AlgebraicType { return &AlgebraicType{Kind: I64Kind} }

// CreateF32Type creates an AlgebraicType for the F32 type.
func CreateF32Type() *AlgebraicType { return &AlgebraicType{Kind: F32Kind} }

// CreateF64Type creates an AlgebraicType for the F64 type.
func CreateF64Type() *AlgebraicType { return &AlgebraicType{Kind: F64Kind} }

// CreateStringType creates an AlgebraicType for the String type.
func CreateStringType() *AlgebraicType { return &AlgebraicType{Kind: StringKind} }

// CreateBytesType creates an AlgebraicType for the Bytes type.
func CreateBytesType() *AlgebraicType { return &AlgebraicType{Kind: BytesKind} }

// CreateU128Type creates an AlgebraicType for the U128 type.
func CreateU128Type() *AlgebraicType { return &AlgebraicType{Kind: U128Kind} }

// CreateI128Type creates an AlgebraicType for the I128 type.
func CreateI128Type() *AlgebraicType { return &AlgebraicType{Kind: I128Kind} }

// CreateU256Type creates an AlgebraicType for the U256 type.
func CreateU256Type() *AlgebraicType { return &AlgebraicType{Kind: U256Kind} }

// CreateI256Type creates an AlgebraicType for the I256 type.
func CreateI256Type() *AlgebraicType { return &AlgebraicType{Kind: I256Kind} }

// CreateTimestampType creates an AlgebraicType for Timestamp (represented as U64 micros).
func CreateTimestampType() *AlgebraicType { return &AlgebraicType{Kind: TimestampKind} }

// CreateTimeDurationType creates an AlgebraicType for TimeDuration (represented as I64 micros).
func CreateTimeDurationType() *AlgebraicType { return &AlgebraicType{Kind: TimeDurationKind} }

// CreateConnectionIdType creates an AlgebraicType for ConnectionId (represented as U128).
func CreateConnectionIdType() *AlgebraicType { return &AlgebraicType{Kind: ConnectionIdKind} }

// CreateIdentityType creates an AlgebraicType for Identity (represented as U256).
func CreateIdentityType() *AlgebraicType { return &AlgebraicType{Kind: IdentityKind} }

// CreateSumType creates an AlgebraicType for a sum type.
func CreateSumType(name string, variants []*SumTypeVariant) *AlgebraicType {
	return &AlgebraicType{
		Name:       name,
		Kind:       SumKind,
		SumDetails: struct{ Variants []*SumTypeVariant }{Variants: variants},
	}
}

// CreateOptionType creates an AlgebraicType for an option type.
// Option<T> is represented as a SumType with variants "some"(T) and "none"(Unit).
func CreateOptionType(name string, someType *AlgebraicType) *AlgebraicType {
	optionName := name
	if optionName == "" {
		optionName = fmt.Sprintf("Option<%s>", someType.Name) // Or some other naming convention
	}
	return CreateSumType(optionName, []*SumTypeVariant{
		NewSumTypeVariant("some", someType),
		NewSumTypeVariant("none", CreateUnitType()),
	})
}

// CreateStructType creates an AlgebraicType for a struct type.
func CreateStructType(name string, fields []StructField) *AlgebraicType {
	return &AlgebraicType{
		Name:          name,
		Kind:          StructKind,
		StructDetails: struct{ Fields []StructField }{Fields: fields},
	}
}

// CreateArrayType creates an AlgebraicType for an array type.
func CreateArrayType(name string, elementType *AlgebraicType) *AlgebraicType {
	arrName := name
	if arrName == "" {
		if elementType != nil && elementType.Name != "" {
			arrName = fmt.Sprintf("Array<%s>", elementType.Name)
		} else {
			arrName = "Array_UnknownElement"
		}
	}
	return &AlgebraicType{
		Name:         arrName,
		Kind:         ArrayKind,
		ArrayDetails: struct{ ElementType *AlgebraicType }{ElementType: elementType},
	}
}

// CreateMapType creates an AlgebraicType for a map type.
func CreateMapType(name string, keyType, valueType *AlgebraicType) *AlgebraicType {
	mapName := name
	if mapName == "" {
		mapName = fmt.Sprintf("Map<%s, %s>", keyType.Name, valueType.Name)
	}
	return &AlgebraicType{
		Name: mapName,
		Kind: MapKind,
		MapDetails: struct {
			KeyType   *AlgebraicType
			ValueType *AlgebraicType
		}{KeyType: keyType, ValueType: valueType},
	}
}

// CreateTupleType creates an AlgebraicType for a tuple type.
func CreateTupleType(name string, elementTypes []*AlgebraicType) *AlgebraicType {
	return &AlgebraicType{
		Name:         name,
		Kind:         TupleKind,
		TupleDetails: struct{ ElementTypes []*AlgebraicType }{ElementTypes: elementTypes},
	}
}

// Serialize writes the given value to the BinaryWriter according to the AlgebraicType definition.
func (at *AlgebraicType) Serialize(writer *BinaryWriter, value interface{}) error {
	switch at.Kind {
	case UnitKind:
		return nil // Unit type serializes to nothing
	case BoolKind:
		v, ok := value.(bool)
		if !ok {
			return fmt.Errorf("%s: expected bool, got %T", at.Name, value)
		}
		writer.WriteBool(v)
	case U8Kind:
		v, ok := value.(uint8)
		if !ok {
			return fmt.Errorf("%s: expected uint8, got %T", at.Name, value)
		}
		writer.WriteU8(v)
	case I8Kind:
		v, ok := value.(int8)
		if !ok {
			return fmt.Errorf("%s: expected int8, got %T", at.Name, value)
		}
		writer.WriteI8(v)
	case U16Kind:
		v, ok := value.(uint16)
		if !ok {
			return fmt.Errorf("%s: expected uint16, got %T", at.Name, value)
		}
		writer.WriteU16(v)
	case I16Kind:
		v, ok := value.(int16)
		if !ok {
			return fmt.Errorf("%s: expected int16, got %T", at.Name, value)
		}
		writer.WriteI16(v)
	case U32Kind:
		v, ok := value.(uint32)
		if !ok {
			return fmt.Errorf("%s: expected uint32, got %T", at.Name, value)
		}
		writer.WriteU32(v)
	case I32Kind:
		v, ok := value.(int32)
		if !ok {
			return fmt.Errorf("%s: expected int32, got %T", at.Name, value)
		}
		writer.WriteI32(v)
	case U64Kind:
		v, ok := value.(uint64)
		if !ok {
			return fmt.Errorf("%s: expected uint64, got %T", at.Name, value)
		}
		writer.WriteU64(v)
	case I64Kind:
		v, ok := value.(int64)
		if !ok {
			return fmt.Errorf("%s: expected int64, got %T", at.Name, value)
		}
		writer.WriteI64(v)
	case F32Kind:
		v, ok := value.(float32)
		if !ok {
			return fmt.Errorf("%s: expected float32, got %T", at.Name, value)
		}
		writer.WriteF32(v)
	case F64Kind:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("%s: expected float64, got %T", at.Name, value)
		}
		writer.WriteF64(v)
	case StringKind:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("%s: expected string, got %T", at.Name, value)
		}
		writer.WriteString(v)
	case BytesKind:
		v, ok := value.([]byte)
		if !ok {
			return fmt.Errorf("%s: expected []byte, got %T", at.Name, value)
		}
		writer.WriteUInt8Array(v)
	case U128Kind:
		v, ok := value.(*big.Int) // Assuming U128 is represented as *big.Int for writer
		if !ok {
			return fmt.Errorf("%s: expected *big.Int for U128Kind, got %T", at.Name, value)
		}
		writer.WriteU128(v)
	case I128Kind:
		v, ok := value.(*big.Int)
		if !ok {
			return fmt.Errorf("%s: expected *big.Int for I128Kind, got %T", at.Name, value)
		}
		writer.WriteI128(v)
	case U256Kind:
		v, ok := value.(*big.Int)
		if !ok {
			return fmt.Errorf("%s: expected *big.Int for U256Kind, got %T", at.Name, value)
		}
		writer.WriteU256(v)
	case I256Kind:
		v, ok := value.(*big.Int)
		if !ok {
			return fmt.Errorf("%s: expected *big.Int for I256Kind, got %T", at.Name, value)
		}
		writer.WriteI256(v)

	case TimestampKind:
		v, ok := value.(*Timestamp)
		if !ok {
			return fmt.Errorf("%s: expected *Timestamp, got %T", at.Name, value)
		}
		writer.WriteI64(v.MicrosSinceUnixEpoch())
	case TimeDurationKind:
		// spacetimedb.TimeDuration uses *big.Int. If protocol is I64:
		td, ok := value.(*TimeDuration) // Using existing spacetimedb.TimeDuration
		if !ok {
			return fmt.Errorf("%s: expected *spacetimedb.TimeDuration, got %T", at.Name, value)
		}
		// Assuming protocol uses I64 for TimeDuration micros. This might need adjustment
		// if the protocol expects a big int representation directly.
		// For now, let's assume I64 for simplicity as per earlier thoughts.
		// This requires TimeDuration.Micros to be convertible to int64.
		if td.Micros == nil {
			return fmt.Errorf("%s: TimeDuration.Micros is nil", at.Name)
		}
		writer.WriteI64(td.Micros.Int64()) // Potential precision loss if Micros > max int64
	case ConnectionIdKind:
		cid, ok := value.(*ConnectionId) // Using existing spacetimedb.ConnectionId
		if !ok {
			return fmt.Errorf("%s: expected *spacetimedb.ConnectionId, got %T", at.Name, value)
		}
		if cid.data == nil {
			return fmt.Errorf("%s: ConnectionId.data is nil", at.Name)
		}
		writer.WriteU128(cid.data)
	case IdentityKind:
		id, ok := value.(*Identity) // Using existing spacetimedb.Identity
		if !ok {
			return fmt.Errorf("%s: expected *spacetimedb.Identity, got %T", at.Name, value)
		}
		if id.data == nil {
			return fmt.Errorf("%s: Identity.data is nil", at.Name)
		}
		writer.WriteU256(id.data)

	case SumKind:
		// Special handling for Option<T> (sum type with "some" and "none")
		isOption := len(at.SumDetails.Variants) == 2 &&
			at.SumDetails.Variants[0].Name == "some" &&
			at.SumDetails.Variants[1].Name == "none"

		if isOption {
			if value != nil { // Some(value)
				writer.WriteU8(0) // Tag for "some"
				variantType := at.SumDetails.Variants[0].AlgebraicType
				if variantType == nil {
					return fmt.Errorf("%s: 'some' variant type is nil", at.Name)
				}
				return variantType.Serialize(writer, value)
			}
			// None
			writer.WriteU8(1) // Tag for "none"
			// "none" variant's type is Unit, which serializes to nothing.
			return nil
		}

		// General sum type serialization
		sumMap, ok := value.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%s: value for SumKind must be map[string]interface{}, got %T", at.Name, value)
		}
		if len(sumMap) != 1 {
			return fmt.Errorf("%s: sum type map must have one entry, got %d", at.Name, len(sumMap))
		}

		for variantNameInMap, payload := range sumMap { // Iterates once
			found := false
			for i, variant := range at.SumDetails.Variants {
				if variant.Name == variantNameInMap {
					writer.WriteU8(uint8(i)) // Write tag (index)
					if variant.AlgebraicType == nil {
						return fmt.Errorf("%s: variant '%s' type is nil", at.Name, variant.Name)
					}
					err := variant.AlgebraicType.Serialize(writer, payload)
					if err != nil {
						return fmt.Errorf("%s: failed to serialize payload for variant '%s': %w", at.Name, variant.Name, err)
					}
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("%s: variant '%s' not found in definition", at.Name, variantNameInMap)
			}
			break
		}
	case ArrayKind:
		valSlice := reflect.ValueOf(value)
		if valSlice.Kind() != reflect.Slice {
			return fmt.Errorf("%s: expected slice for ArrayKind, got %T", at.Name, value)
		}

		if at.ArrayDetails.ElementType == nil {
			return fmt.Errorf("%s: array element type is nil for ArrayKind", at.Name)
		}

		if at.ArrayDetails.ElementType.Kind == U8Kind {
			byteSlice, ok := value.([]byte)
			if !ok {
				if valSlice.Type().Elem().Kind() == reflect.Uint8 {
					tempSlice := make([]byte, valSlice.Len())
					for i := 0; i < valSlice.Len(); i++ {
						uintVal := valSlice.Index(i).Uint()
						if uintVal > 255 {
							return fmt.Errorf("%s: element at index %d (%v) out of byte range for ArrayKind<U8>", at.Name, i, uintVal)
						}
						tempSlice[i] = byte(uintVal)
					}
					byteSlice = tempSlice
					ok = true
				}
				if !ok {
					return fmt.Errorf("%s: expected []byte or slice of uint8 for ArrayKind<U8>, got %T", at.Name, value)
				}
			}
			writer.WriteUInt8Array(byteSlice)
		} else {
			length := valSlice.Len()
			writer.WriteU32(uint32(length))
			for i := 0; i < length; i++ {
				elem := valSlice.Index(i).Interface()
				if err := at.ArrayDetails.ElementType.Serialize(writer, elem); err != nil {
					return fmt.Errorf("%s: failed to serialize array element %d: %w", at.Name, i, err)
				}
			}
		}
	case StructKind:
		valMap, ok := value.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%s: expected map[string]interface{} for StructKind, got %T", at.Name, value)
		}
		for _, field := range at.StructDetails.Fields {
			fieldValue, exists := valMap[field.Name]
			if !exists {
				return fmt.Errorf("%s: field '%s' missing in value for StructKind", at.Name, field.Name)
			}
			if field.Type == nil {
				return fmt.Errorf("%s: type for field '%s' is nil", at.Name, field.Name)
			}
			if err := field.Type.Serialize(writer, fieldValue); err != nil {
				return fmt.Errorf("%s: failed to serialize field '%s': %w", at.Name, field.Name, err)
			}
		}
	case MapKind:
		valMap := reflect.ValueOf(value)
		if valMap.Kind() != reflect.Map {
			return fmt.Errorf("%s: expected map for MapKind, got %T", at.Name, value)
		}
		writer.WriteU32(uint32(valMap.Len()))
		if at.MapDetails.KeyType == nil {
			return fmt.Errorf("%s: map key type is nil", at.Name)
		}
		if at.MapDetails.ValueType == nil {
			return fmt.Errorf("%s: map value type is nil", at.Name)
		}

		iter := valMap.MapRange()
		for iter.Next() {
			k := iter.Key().Interface()
			v := iter.Value().Interface()
			if err := at.MapDetails.KeyType.Serialize(writer, k); err != nil {
				return fmt.Errorf("%s: failed to serialize map key: %w", at.Name, err)
			}
			if err := at.MapDetails.ValueType.Serialize(writer, v); err != nil {
				return fmt.Errorf("%s: failed to serialize map value: %w", at.Name, err)
			}
		}
	case TupleKind:
		valSlice, ok := value.([]interface{})
		if !ok {
			return fmt.Errorf("%s: expected []interface{} for TupleKind, got %T", at.Name, value)
		}
		if len(valSlice) != len(at.TupleDetails.ElementTypes) {
			return fmt.Errorf("%s: tuple value length %d does not match definition length %d", at.Name, len(valSlice), len(at.TupleDetails.ElementTypes))
		}
		for i, elemType := range at.TupleDetails.ElementTypes {
			if elemType == nil {
				return fmt.Errorf("%s: tuple element type at index %d is nil", at.Name, i)
			}
			if err := elemType.Serialize(writer, valSlice[i]); err != nil {
				return fmt.Errorf("%s: failed to serialize tuple element %d: %w", at.Name, i, err)
			}
		}
	default:
		return fmt.Errorf("%s: serialization not implemented for kind %v", at.Name, at.Kind)
	}
	return nil
}

// Deserialize reads from the BinaryReader and constructs a Go value according to the AlgebraicType.
// Note: BinaryReader methods (e.g., ReadU8) in spacetimedb currently panic on error.
// This Deserialize method assumes they might return (value, error) or that panics are handled upstream.
// For robustness, BinaryReader methods should ideally return errors.
func (at *AlgebraicType) Deserialize(reader *BinaryReader) (interface{}, error) {
	switch at.Kind {
	case UnitKind:
		return nil, nil
	case BoolKind:
		return reader.ReadBool(), nil
	case U8Kind:
		return reader.ReadU8(), nil
	case I8Kind:
		return reader.ReadI8(), nil
	case U16Kind:
		return reader.ReadU16(), nil
	case I16Kind:
		return reader.ReadI16(), nil
	case U32Kind:
		return reader.ReadU32(), nil
	case I32Kind:
		return reader.ReadI32(), nil
	case U64Kind:
		return reader.ReadU64(), nil
	case I64Kind:
		return reader.ReadI64(), nil
	case F32Kind:
		return reader.ReadF32(), nil
	case F64Kind:
		return reader.ReadF64(), nil
	case StringKind:
		return reader.ReadString(), nil
	case BytesKind:
		return reader.ReadUInt8Array(), nil
	case U128Kind:
		return reader.ReadU128(), nil
	case I128Kind:
		return reader.ReadI128(), nil
	case U256Kind:
		return reader.ReadU256(), nil
	case I256Kind:
		return reader.ReadI256(), nil

	case TimestampKind:
		micros := reader.ReadI64()
		// return SATSTimestamp{MicrosSinceUnixEpoch: micros}, nil
		return NewTimestamp(micros), nil
	case TimeDurationKind:
		// Assuming protocol uses I64 for TimeDuration micros.
		microsVal := reader.ReadI64()
		return NewTimeDuration(big.NewInt(microsVal)), nil // Returns *spacetimedb.TimeDuration
	case ConnectionIdKind:
		data := reader.ReadU128()
		return NewConnectionId(data), nil // Returns *spacetimedb.ConnectionId
	case IdentityKind:
		data := reader.ReadU256()
		// NewIdentity takes interface{} (string or *big.Int). We have *big.Int.
		id, idErr := NewIdentity(data)
		if idErr != nil {
			return nil, fmt.Errorf("%s: failed to create Identity from *big.Int: %w", at.Name, idErr)
		}
		return id, nil

	case SumKind:
		isOption := len(at.SumDetails.Variants) == 2 &&
			at.SumDetails.Variants[0].Name == "some" &&
			at.SumDetails.Variants[1].Name == "none"

		tag := reader.ReadU8()

		if isOption {
			if tag == 0 { // Some
				variantType := at.SumDetails.Variants[0].AlgebraicType
				if variantType == nil {
					return nil, fmt.Errorf("%s: 'some' variant type is nil", at.Name)
				}
				return variantType.Deserialize(reader)
			} else if tag == 1 { // None
				// "none" variant's type is Unit, deserializes to nil.
				return nil, nil
			}
			return nil, fmt.Errorf("%s: invalid tag %d for option type", at.Name, tag)
		}

		// General sum type deserialization
		if int(tag) >= len(at.SumDetails.Variants) {
			return nil, fmt.Errorf("%s: invalid tag %d for sum type with %d variants", at.Name, tag, len(at.SumDetails.Variants))
		}
		variant := at.SumDetails.Variants[tag]
		if variant.AlgebraicType == nil {
			return nil, fmt.Errorf("%s: variant '%s' (tag %d) type is nil", at.Name, variant.Name, tag)
		}
		payload, err := variant.AlgebraicType.Deserialize(reader)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to deserialize payload for variant '%s': %w", at.Name, variant.Name, err)
		}
		return map[string]interface{}{variant.Name: payload}, nil
	case ArrayKind:
		if at.ArrayDetails.ElementType == nil {
			return nil, fmt.Errorf("%s: array element type is nil for ArrayKind", at.Name)
		}

		if at.ArrayDetails.ElementType.Kind == U8Kind {
			// Assuming ReadUInt8Array reads length and then bytes, and returns []byte or error
			// If ReadUInt8Array doesn't return an error, and error handling is different, adjust this.
			bytes := reader.ReadUInt8Array() // If this can fail, it should return an error
			// if reader has an error state: if reader.Error() != nil { return nil, reader.Error() }
			return bytes, nil // Assuming `bytes` is the correct return type e.g. []byte
		} else {
			length := reader.ReadU32() // If this can fail, it should return an error
			// if reader has an error state: if reader.Error() != nil { return nil, reader.Error() }

			elemType := at.ArrayDetails.ElementType.GoType()
			var slice reflect.Value
			if elemType != nil {
				slice = reflect.MakeSlice(reflect.SliceOf(elemType), int(length), int(length))
			} else {
				objSlice := make([]interface{}, length)
				for i := uint32(0); i < length; i++ {
					elem, errDeserialize := at.ArrayDetails.ElementType.Deserialize(reader)
					if errDeserialize != nil {
						return nil, fmt.Errorf("%s: failed to deserialize array element %d: %w", at.Name, i, errDeserialize)
					}
					objSlice[i] = elem
				}
				return objSlice, nil
			}

			for i := uint32(0); i < length; i++ {
				elem, errDeserialize := at.ArrayDetails.ElementType.Deserialize(reader)
				if errDeserialize != nil {
					return nil, fmt.Errorf("%s: failed to deserialize array element %d: %w", at.Name, i, errDeserialize)
				}
				if reflect.TypeOf(elem).AssignableTo(elemType) {
					slice.Index(int(i)).Set(reflect.ValueOf(elem))
				} else if elemType.Kind() == reflect.Interface && reflect.TypeOf(elem).Implements(elemType) {
					slice.Index(int(i)).Set(reflect.ValueOf(elem))
				} else if reflect.ValueOf(elem).Type().ConvertibleTo(elemType) {
					slice.Index(int(i)).Set(reflect.ValueOf(elem).Convert(elemType))
				} else {
					return nil, fmt.Errorf("%s: deserialized element type %T not assignable or convertible to slice element type %v", at.Name, elem, elemType)
				}
			}
			return slice.Interface(), nil
		}
	case StructKind:
		resultMap := make(map[string]interface{})
		for _, field := range at.StructDetails.Fields {
			if field.Type == nil {
				return nil, fmt.Errorf("%s: type for field '%s' is nil", at.Name, field.Name)
			}
			fieldValue, err := field.Type.Deserialize(reader)
			if err != nil {
				return nil, fmt.Errorf("%s: failed to deserialize field '%s': %w", at.Name, field.Name, err)
			}
			resultMap[field.Name] = fieldValue
		}
		return resultMap, nil
	case MapKind:
		count := reader.ReadU32()
		resultMap := make(map[interface{}]interface{}) // Go map keys must be comparable
		if at.MapDetails.KeyType == nil {
			return nil, fmt.Errorf("%s: map key type is nil", at.Name)
		}
		if at.MapDetails.ValueType == nil {
			return nil, fmt.Errorf("%s: map value type is nil", at.Name)
		}

		for i := uint32(0); i < count; i++ {
			key, err := at.MapDetails.KeyType.Deserialize(reader)
			if err != nil {
				return nil, fmt.Errorf("%s: failed to deserialize map key for entry %d: %w", at.Name, i, err)
			}
			// Ensure key is comparable for Go map
			// This is a runtime check; ideally, MapKeyType should be restricted to comparable types at definition.
			// For simplicity, this example assumes deserialized keys will be comparable.

			value, err := at.MapDetails.ValueType.Deserialize(reader)
			if err != nil {
				return nil, fmt.Errorf("%s: failed to deserialize map value for entry %d: %w", at.Name, i, err)
			}
			resultMap[key] = value
		}
		return resultMap, nil
	case TupleKind:
		numElements := len(at.TupleDetails.ElementTypes)
		resultSlice := make([]interface{}, numElements)
		for i, elemType := range at.TupleDetails.ElementTypes {
			if elemType == nil {
				return nil, fmt.Errorf("%s: tuple element type at index %d is nil", at.Name, i)
			}
			elemValue, err := elemType.Deserialize(reader)
			if err != nil {
				return nil, fmt.Errorf("%s: failed to deserialize tuple element %d: %w", at.Name, i, err)
			}
			resultSlice[i] = elemValue
		}
		return resultSlice, nil
	default:
		return nil, fmt.Errorf("%s: deserialization not implemented for kind %v", at.Name, at.Kind)
	}
}

// Helper to check if an AlgebraicType is a Unit type.
// This was mentioned in the TS SumType.serialize fallback logic.
func (at *AlgebraicType) isUnit() bool {
	return at.Kind == UnitKind
}

// GoType returns the Go reflect.Type corresponding to the AlgebraicType.
// This is a simplified version for use in Deserialize.
func (at *AlgebraicType) GoType() reflect.Type {
	switch at.Kind {
	case UnitKind:
		return reflect.TypeOf(nil) // Or a specific struct{}{} type
	case BoolKind:
		return reflect.TypeOf(false)
	case U8Kind:
		return reflect.TypeOf(uint8(0))
	case I8Kind:
		return reflect.TypeOf(int8(0))
	case U16Kind:
		return reflect.TypeOf(uint16(0))
	case I16Kind:
		return reflect.TypeOf(int16(0))
	case U32Kind:
		return reflect.TypeOf(uint32(0))
	case I32Kind:
		return reflect.TypeOf(int32(0))
	case U64Kind:
		return reflect.TypeOf(uint64(0))
	case I64Kind:
		return reflect.TypeOf(int64(0))
	case F32Kind:
		return reflect.TypeOf(float32(0))
	case F64Kind:
		return reflect.TypeOf(float64(0))
	case StringKind:
		return reflect.TypeOf("")
	case BytesKind: // This kind might be redundant if ArrayKind<U8> handles it
		return reflect.TypeOf([]byte{})
	case U128Kind, I128Kind, U256Kind, I256Kind:
		return reflect.TypeOf((*big.Int)(nil))
	case TimestampKind:
		return reflect.TypeOf((*Timestamp)(nil))
	case TimeDurationKind:
		return reflect.TypeOf((*TimeDuration)(nil))
	case ConnectionIdKind:
		return reflect.TypeOf((*ConnectionId)(nil))
	case IdentityKind:
		return reflect.TypeOf((*Identity)(nil))
	// For ArrayKind, this method would be called on at.ArrayDetails.ElementType
	// So, the cases above for primitive types are what ArrayDetails.ElementType.GoType() would return.
	default:
		return nil
	}
}
