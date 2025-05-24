package spacetimedb

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
)

type JSONIdentity struct {
	Identity string `json:"__identity__"`
}

type JSONMessageConnectionID struct {
	ConnectionID *big.Int `json:"__connection_id__"`
}

type JSONIdentityToken struct {
	Identity     JSONIdentity            `json:"identity"`
	Token        string                  `json:"token"`
	ConnectionID JSONMessageConnectionID `json:"connection_id"`
}

func (db *DBConnection) parseJsonMessage(msg []byte) error {
	var rawMessages map[string]json.RawMessage
	err := json.Unmarshal(msg, &rawMessages)
	if err != nil {
		return fmt.Errorf("error unmarshalling raw message: %w", err)
	}

	for key, rawMessage := range rawMessages {
		switch key {
		case "IdentityToken":
			var identityToken JSONIdentityToken
			err := json.Unmarshal(rawMessage, &identityToken)
			if err != nil {
				return fmt.Errorf("error unmarshalling IdentityToken: %w", err)
			}
			// TODO: Process the unmarshalled IdentityToken
			log.Printf("Received IdentityToken: %+v\n\n", identityToken)
		default:
			return fmt.Errorf("unknown message type: %s", key)
		}
		// Assuming one top-level key per message. Still not sure about this.
		break
	}

	return nil
}
