package spacetimedb

import "fmt"

type ServerMessage struct {
	Message any // Union type (go please!)
}

func (sm *ServerMessage) Deserialize(reader *BinaryReader) error {
	unionType := reader.ReadU8()
	switch unionType {
	case 0x00:
		fmt.Println("ServerMessage.Deserialize: type 0x00 is not implemented yet")
	case 0x01:
		transactionUpdate := &TransactionUpdate{}
		if err := transactionUpdate.Deserialize(reader); err != nil {
			return fmt.Errorf("failed to deserialize TransactionUpdate: %w", err)
		}
		sm.Message = transactionUpdate
	case 0x02:
		fmt.Println("ServerMessage.Deserialize: type 0x02 is not implemented yet")
	case 0x03:
		identityToken := &IdentityToken{}
		if err := identityToken.Deserialize(reader); err != nil {
			return fmt.Errorf("failed to deserialize IdentityToken: %w", err)
		}
		sm.Message = identityToken
	}

	return nil
}
