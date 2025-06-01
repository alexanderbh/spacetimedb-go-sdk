package spacetimedb

import "fmt"

type CallReducer struct {
	Reducer   string
	Args      []byte
	RequestId uint32
	Flags     uint8
}

func (cr *CallReducer) Serialize(writer *BinaryWriter) error {
	writer.WriteString(cr.Reducer)
	writer.WriteUInt8Array(cr.Args)
	writer.WriteU32(cr.RequestId)
	writer.WriteU8(cr.Flags)
	return nil
}

func (conn *DBConnection) CallReducer(reducer string, args []byte, requestId uint32, flags uint8) error {

	clientMsg := &ClientMessage{
		Message: &CallReducer{
			Reducer:   reducer,
			Args:      args,
			RequestId: requestId,
			Flags:     flags,
		},
	}

	writer := NewBinaryWriter()

	clientMsg.Serialize(writer)

	msg := writer.GetBuffer()
	fmt.Printf("Sending ClientMessage: %s\n", clientMsg)
	return conn.SendMessage(msg)
}
