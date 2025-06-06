package spacetimedb

type Subscribe struct {
	QueryStrings []string
	RequestId    uint32
}

func (cr *Subscribe) Serialize(writer *BinaryWriter) error {
	WriteArray(writer, cr.QueryStrings, func(writer *BinaryWriter, item string) {
		writer.WriteString(item)
	})
	writer.WriteU32(cr.RequestId)
	return nil
}

func (conn *DBConnection) Subscribe(queryStrings ...string) error {

	clientMsg := &ClientMessage{
		Message: &Subscribe{
			QueryStrings: queryStrings,
			RequestId:    0, // TODO: Support request IDs
		},
	}

	writer := NewBinaryWriter()

	clientMsg.Serialize(writer)

	msg := writer.GetBuffer()
	conn.Logger("Sending ClientMessage: %s", clientMsg)
	return conn.SendMessage(msg)
}
