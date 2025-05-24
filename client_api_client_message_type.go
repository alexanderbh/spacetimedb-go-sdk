package spacetimedb

import "fmt"

var clientMessageAt *AlgebraicType

func ClientMessage_GetAlgebraicType() *AlgebraicType {
	if clientMessageAt == nil {
		clientMessageAt = CreateSumType("ClientMessage", []*SumTypeVariant{
			{Name: "CallReducer", AlgebraicType: CallReducer_GetAlgebraicType()},
		})
	}
	return clientMessageAt
}

func NewMapFromClientMessage(msg any) (map[string]any, error) {
	if msg == nil {
		return nil, nil
	}

	switch v := msg.(type) {
	case *CallReducer:
		m, err := NewMapFromCallReducer(v)
		if err != nil {
			return nil, err
		}
		return map[string]any{"CallReducer": m}, nil
	default:
		return nil, fmt.Errorf("unsupported message type: %T", v)
	}
}
