package spacetimedb

import "fmt"

type ClientMessage struct {
	Message any // Union type (go please!)
}

func (sm *ClientMessage) Serialize(writer *BinaryWriter) error {
	switch v := sm.Message.(type) {
	case *CallReducer:
		writer.WriteU8(0x00) // Type identifier for CallReducer
		return v.Serialize(writer)
	case *Subscribe:
		writer.WriteU8(0x01) // Type identifier for Subscribe
		return v.Serialize(writer)
	}
	return fmt.Errorf("unsupported message type when serializing ClientMessage: %T", sm.Message)
}
