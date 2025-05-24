package spacetimedb

import "fmt"

type CallReducer struct {
	Reducer   string
	Args      []byte
	RequestId uint32
	Flags     uint8
}

var callReducerAt *AlgebraicType

func CallReducer_GetAlgebraicType() *AlgebraicType {
	if callReducerAt == nil {
		callReducerAt = CreateStructType("CallReducer", []StructField{
			{Name: "reducer", Type: CreateStringType()},
			{Name: "args", Type: CreateArrayType("Args", CreateU8Type())},
			{Name: "requestId", Type: CreateU32Type()},
			{Name: "flags", Type: CreateU8Type()},
		})
	}
	return callReducerAt
}

func NewMapFromCallReducer(reducer *CallReducer) (map[string]any, error) {
	if reducer == nil {
		return nil, fmt.Errorf("nil CallReducer provided")
	}

	return map[string]any{
		"reducer":   reducer.Reducer,
		"args":      reducer.Args,
		"requestId": reducer.RequestId,
		"flags":     reducer.Flags,
	}, nil
}

func NewCallReducerFromMap(m map[string]any) (*CallReducer, error) {
	if m == nil {
		return nil, fmt.Errorf("nil map provided")
	}

	reducer, ok := m["reducer"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'reducer' field")
	}

	args, ok := m["args"].([]byte)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'args' field")
	}

	requestId, ok := m["requestId"].(uint32)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'requestId' field")
	}

	flags, ok := m["flags"].(uint8)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'flags' field")
	}

	return &CallReducer{
		Reducer:   reducer,
		Args:      args,
		RequestId: requestId,
		Flags:     flags,
	}, nil
}
