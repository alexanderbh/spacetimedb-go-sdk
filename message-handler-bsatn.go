package spacetimedb

import (
	"fmt"
)

func (db *DBConnection) parseBsantMessage(msg []byte) error {
	reader := NewBinaryReader(msg)

	// Handle compression
	compression := reader.ReadU8()
	switch compression {
	case CompressionTypeNone:
		// No compression, read the message directly
	case CompressionTypeBrotly:
		return fmt.Errorf("Brotly compression is not implemented yet")
	case CompressionTypeGzip:
		return fmt.Errorf("Gzip compression is not implemented yet")
	default:
		return fmt.Errorf("unknown compression type: %d", compression)
	}

	// Read the message type
	deserializedMsg, err := ServerMessage_GetAlgebraicType().Deserialize(reader)
	if err != nil {
		return fmt.Errorf("failed to deserialize message: %w", err)
	}
	if deserializedMsg == nil {
		return fmt.Errorf("deserialized message is nil")
	}
	fmt.Printf("Received message: %s\n\n", deserializedMsg)

	typedMsg := deserializedMsg.(map[string]any)

	// Assert assumption that the message contains a single entry is correect
	if len(typedMsg) != 1 {
		return fmt.Errorf("expected a single entry in the deserialized message, got %d entries", len(typedMsg))
	}

	for key := range typedMsg {
		switch key {
		case "IdentityToken":
			identityToken, ok := typedMsg[key].(map[string]any)
			if !ok {
				return fmt.Errorf("failed to cast %s to map[string]interface{}", key)
			}

			idToken, err := NewIdentityTokenFromMap(identityToken)
			if err != nil {
				return fmt.Errorf("failed to create IdentityToken from map: %w", err)
			}
			fmt.Printf("Received IdentityToken: %+v\n\n", idToken)

			db.IsConnected = true
			db.Identity = idToken.Identity
			if db.Token == "" && idToken.Token != "" {
				db.Token = idToken.Token
			}
			db.ConnectionId = idToken.ConnectionId
			if db.OnConnect != nil {
				db.OnConnect(db, idToken.Identity, idToken.Token, idToken.ConnectionId)
			}
		}
	}

	return nil
}
