package spacetimedb

import (
	"encoding/json"
	"fmt"
	"math/big"
)

// Identity represents the __identity__ field
type Identity struct {
	Identity string `json:"__identity__"`
}

// MessageConnectionID represents the __connection_id__ field in the message
type MessageConnectionID struct {
	ConnectionID *big.Int `json:"__connection_id__"`
}

// IdentityToken represents the top-level JSON structure
type IdentityToken struct {
	Identity     Identity            `json:"identity"`
	Token        string              `json:"token"`
	ConnectionID MessageConnectionID `json:"connection_id"`
}

func (db *DBConnection) handleMessage(msg []byte) error {
	var rawMessages map[string]json.RawMessage
	err := json.Unmarshal(msg, &rawMessages)
	if err != nil {
		return fmt.Errorf("error unmarshalling raw message: %w", err)
	}

	for key, rawMessage := range rawMessages {
		switch key {
		case "IdentityToken":
			var identityToken IdentityToken
			err := json.Unmarshal(rawMessage, &identityToken)
			if err != nil {
				return fmt.Errorf("error unmarshalling IdentityToken: %w", err)
			}
			// TODO: Process the unmarshalled IdentityToken
			fmt.Printf("Received IdentityToken: %+v\\n", identityToken)
		default:
			return fmt.Errorf("unknown message type: %s", key)
		}
		// Assuming one top-level key per message as per examples
		break
	}

	return nil
}
