package identity_token

import (
	"github.com/alexanderbh/spacetimedb-go-sdk"
)

type IdentityTokenType struct {
	Identity     *spacetimedb.Identity     `json:"identity"`
	Token        string                    `json:"token"`
	ConnectionId *spacetimedb.ConnectionId `json:"connectionId"`
}

var at *spacetimedb.AlgebraicType

func GetAlgebraicType() *spacetimedb.AlgebraicType {
	if at == nil {
		at = spacetimedb.CreateStructType("IdentityToken", []spacetimedb.StructField{
			{Name: "identity", Type: spacetimedb.CreateIdentityType()},
			{Name: "token", Type: spacetimedb.CreateStringType()},
			{Name: "connectionId", Type: spacetimedb.CreateConnectionIdType()},
		})
	}
	return at
}
