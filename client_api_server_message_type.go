package spacetimedb

var serverMessageAt *AlgebraicType

func ServerMessage_GetAlgebraicType() *AlgebraicType {
	if at == nil {
		at = CreateSumType("ServerMessage", []*SumTypeVariant{
			{Name: "id0", AlgebraicType: IdentityToken_GetAlgebraicType()},           // TODO: change index 0
			{Name: "id1", AlgebraicType: IdentityToken_GetAlgebraicType()},           // TODO: change index 1
			{Name: "id2", AlgebraicType: IdentityToken_GetAlgebraicType()},           // TODO: change index 2
			{Name: "IdentityToken", AlgebraicType: IdentityToken_GetAlgebraicType()}, // type: 0x03
		})
	}
	return at
}
