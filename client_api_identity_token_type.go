package spacetimedb

import "fmt"

type IdentityTokenType struct {
	Identity     *Identity     `json:"identity"`
	Token        string        `json:"token"`
	ConnectionId *ConnectionId `json:"connectionId"`
}

var at *AlgebraicType

func IdentityToken_GetAlgebraicType() *AlgebraicType {
	if at == nil {
		at = CreateStructType("IdentityToken", []StructField{
			{Name: "identity", Type: CreateIdentityType()},
			{Name: "token", Type: CreateStringType()},
			{Name: "connectionId", Type: CreateConnectionIdType()},
		})
	}
	return at
}

func NewIdentityTokenFromMap(m map[string]any) (*IdentityTokenType, error) {
	identity, ok := m["identity"].(*Identity)
	if !ok {
		return nil, fmt.Errorf("failed to cast identity to *Identity: %T", m["identity"])
	}

	token, ok := m["token"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to cast token to string")
	}
	connectionId, ok := m["connectionId"].(*ConnectionId)
	if !ok {
		return nil, fmt.Errorf("failed to cast connectionId to *ConnectionId: %T", m["connectionId"])
	}
	return &IdentityTokenType{
		Identity:     identity,
		Token:        token,
		ConnectionId: connectionId,
	}, nil
}
