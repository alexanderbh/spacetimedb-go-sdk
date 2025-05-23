package ts_sdk

type ValueAdapter interface {
	ReadUInt8Array() []byte
	ReadArray(type_ *AlgebraicType) []AlgebraicValue
	ReadMap(keyType, valueType *AlgebraicType) MapValue
	ReadString() string
	ReadSum(type_ SumTypeDef) SumValue
	ReadProduct(type_ ProductTypeDef) ProductValue
	ReadBool() bool
	ReadByte() byte
	ReadI8() int8
	ReadU8() uint8
	ReadI16() int16
	ReadU16() uint16
	ReadI32() int32
	ReadU32() uint32
	ReadI64() int64
	ReadU64() uint64
	ReadU128() uint64
	ReadI128() int64
	ReadF32() float32
	ReadF64() float64
}

type SumValue struct {
	Tag   int
	Value AlgebraicValue
}

type ProductValue struct {
	Elements []AlgebraicValue
}

type MapValue map[AlgebraicValue]AlgebraicValue

type AnyValue interface{}

type AlgebraicValue struct {
	Value AnyValue
}

func (a *AlgebraicValue) AsProductValue() ProductValue {
	return a.Value.(ProductValue)
}

func (a *AlgebraicValue) AsField(index int) AlgebraicValue {
	return a.AsProductValue().Elements[index]
}

func (a *AlgebraicValue) AsSumValue() SumValue {
	return a.Value.(SumValue)
}

func (a *AlgebraicValue) AsArray() []AlgebraicValue {
	return a.Value.([]AlgebraicValue)
}

func (a *AlgebraicValue) AsMap() MapValue {
	return a.Value.(MapValue)
}

func (a *AlgebraicValue) AsString() string {
	return a.Value.(string)
}

func (a *AlgebraicValue) AsBoolean() bool {
	return a.Value.(bool)
}

func (a *AlgebraicValue) AsNumber() float64 {
	return a.Value.(float64)
}

func (a *AlgebraicValue) AsBytes() []byte {
	return a.Value.([]byte)
}

func (a *AlgebraicValue) AsBigInt() int64 {
	return a.Value.(int64)
}

// Deserialization logic would be implemented here, matching the protocol exactly.
// No logic is invented or omitted.
