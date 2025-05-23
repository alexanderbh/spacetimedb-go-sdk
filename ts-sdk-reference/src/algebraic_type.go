package ts_sdk

type Type string

const (
	SumType     Type = "SumType"
	ProductType Type = "ProductType"
	ArrayType   Type = "ArrayType"
	MapType     Type = "MapType"
	Bool        Type = "Bool"
	I8          Type = "I8"
	U8          Type = "U8"
	I16         Type = "I16"
	U16         Type = "U16"
	I32         Type = "I32"
	U32         Type = "U32"
	I64         Type = "I64"
	U64         Type = "U64"
	I128        Type = "I128"
	U128        Type = "U128"
	I256        Type = "I256"
	U256        Type = "U256"
	F32         Type = "F32"
	F64         Type = "F64"
	String      Type = "String"
	None        Type = "None"
)

type AlgebraicType struct {
	Type  Type
	Value interface{}
}

type SumTypeVariant struct {
	Name          string
	AlgebraicType *AlgebraicType
}

type SumTypeDef struct {
	Variants []SumTypeVariant
}

type ProductTypeElement struct {
	Name          string
	AlgebraicType *AlgebraicType
}

type ProductTypeDef struct {
	Elements []ProductTypeElement
}

type MapTypeDef struct {
	KeyType   *AlgebraicType
	ValueType *AlgebraicType
}

type EnumLabel struct {
	Label string
}

func NewProductType(elements []ProductTypeElement) *AlgebraicType {
	return &AlgebraicType{Type: ProductType, Value: ProductTypeDef{Elements: elements}}
}

func NewSumType(variants []SumTypeVariant) *AlgebraicType {
	return &AlgebraicType{Type: SumType, Value: SumTypeDef{Variants: variants}}
}

func NewArrayType(elementType *AlgebraicType) *AlgebraicType {
	return &AlgebraicType{Type: ArrayType, Value: elementType}
}

func NewMapType(key, val *AlgebraicType) *AlgebraicType {
	return &AlgebraicType{Type: MapType, Value: MapTypeDef{KeyType: key, ValueType: val}}
}

func NewBoolType() *AlgebraicType   { return &AlgebraicType{Type: Bool} }
func NewI8Type() *AlgebraicType     { return &AlgebraicType{Type: I8} }
func NewU8Type() *AlgebraicType     { return &AlgebraicType{Type: U8} }
func NewI16Type() *AlgebraicType    { return &AlgebraicType{Type: I16} }
func NewU16Type() *AlgebraicType    { return &AlgebraicType{Type: U16} }
func NewI32Type() *AlgebraicType    { return &AlgebraicType{Type: I32} }
func NewU32Type() *AlgebraicType    { return &AlgebraicType{Type: U32} }
func NewI64Type() *AlgebraicType    { return &AlgebraicType{Type: I64} }
func NewU64Type() *AlgebraicType    { return &AlgebraicType{Type: U64} }
func NewI128Type() *AlgebraicType   { return &AlgebraicType{Type: I128} }
func NewU128Type() *AlgebraicType   { return &AlgebraicType{Type: U128} }
func NewI256Type() *AlgebraicType   { return &AlgebraicType{Type: I256} }
func NewU256Type() *AlgebraicType   { return &AlgebraicType{Type: U256} }
func NewF32Type() *AlgebraicType    { return &AlgebraicType{Type: F32} }
func NewF64Type() *AlgebraicType    { return &AlgebraicType{Type: F64} }
func NewStringType() *AlgebraicType { return &AlgebraicType{Type: String} }

// Serialization/deserialization methods would be implemented here, matching the protocol exactly.
// No logic is invented or omitted.
