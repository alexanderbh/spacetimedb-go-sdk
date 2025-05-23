package spacetimedb

import (
	"log"

	"github.com/gorilla/websocket"
)

// CallReducer calls a reducer on your SpacetimeDB module.
func (db *DBConnection) CallReducer(reducerName string, args []byte, flags string) {
	var flagsNum int
	switch flags {
	case "FullUpdate":
		flagsNum = 0
	case "NoSuccessNotify":
		flagsNum = 1
	default:
		// Defaulting to FullUpdate
		flagsNum = 0
		log.Printf("CallReducer: Unknown flags value '%s', defaulting to FullUpdate (0)", flags)
	}

	// Construct the CallReducer message
	callReducerMsg := CallReducerMessage{
		Reducer:   reducerName,
		Args:      args,
		RequestID: 0, // Based on TypeScript SDK's usage
		Flags:     uint8(flagsNum),
	}

	// Serialize the message
	// Assuming NewBinaryWriter and its methods are defined in the current spacetimedb package
	// and behave similarly to the TypeScript SDK's BinaryWriter.
	writer := NewBinaryWriter()

	// ClientMessage variant tag for CallReducer is 0 (as per sum type definition in TS)
	writer.WriteByte(0)

	// Serialize CallReducerMessage fields
	writer.WriteString(callReducerMsg.Reducer)  // Assumes WriteString prefixes with length
	writer.WriteUInt8Array(callReducerMsg.Args) // Corrected method name
	writer.WriteU32(callReducerMsg.RequestID)   // Corrected method name
	writer.WriteByte(callReducerMsg.Flags)

	messageBytes := writer.GetBuffer() // Corrected method name

	// Send the message
	if err := db.sendMessage(websocket.BinaryMessage, messageBytes); err != nil {
		log.Printf("Error sending CallReducer message for '%s': %v", reducerName, err)
		// Further error handling (e.g., callbacks, events) could be added here if needed.
	} else {
		log.Printf("CallReducer message sent for '%s'", reducerName) // Uncomment for debugging
	}
}
