package server_message

import (
	"github.com/alexanderbh/spacetimedb-go-sdk"
	"github.com/alexanderbh/spacetimedb-go-sdk/client_api/identity_token"
)

var at *spacetimedb.AlgebraicType

func GetAlgebraicType() *spacetimedb.AlgebraicType {
	if at == nil {
		at = spacetimedb.CreateSumType("ServerMessage", []*spacetimedb.SumTypeVariant{
			{Name: "id0", AlgebraicType: identity_token.GetAlgebraicType()},           // TODO: change index 0
			{Name: "id1", AlgebraicType: identity_token.GetAlgebraicType()},           // TODO: change index 1
			{Name: "id2", AlgebraicType: identity_token.GetAlgebraicType()},           // TODO: change index 2
			{Name: "IdentityToken", AlgebraicType: identity_token.GetAlgebraicType()}, // type: 0x03
		})
	}
	return at
}
