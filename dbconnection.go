package spacetimedb

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// CallReducerMessage mirrors the structure needed for a CallReducer operation.
// Based on ts-sdk-reference/src/client_api/call_reducer_type.ts
type CallReducerMessage struct {
	Reducer   string
	Args      []byte
	RequestID uint32
	Flags     uint8
}

type DBConnection struct {
	host    string
	isAlive bool
	conn    *websocket.Conn
}

func NewDBConnection(host string) *DBConnection {
	return &DBConnection{
		host: host,
	}
}

func (db *DBConnection) Connect() error {
	if db.host == "" {
		return fmt.Errorf("host cannot be empty")
	}

	dialer := websocket.DefaultDialer
	dialer.Subprotocols = []string{"v1.bsatn.spacetimedb"}
	c, _, err := dialer.Dial(db.host, nil)

	if err != nil {
		return fmt.Errorf("failed to connect to websocket: %w", err)
	}
	db.conn = c
	db.isAlive = true
	fmt.Printf("Connected to websocket at %s\n", db.host)

	// Handle messages in a separate goroutine
	go func() {
		defer func() {
			db.isAlive = false
			if db.conn != nil {
				db.conn.Close() // Ensure connection is closed if read loop exits
			}
		}()
		for {
			if db.conn == nil {
				log.Println("connection is nil, exiting message read loop")
				return
			}
			messageType, rawMessage, err := db.conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading message: %v\n", err)
				return
			}
			if messageType == websocket.TextMessage {
				log.Printf("Received text message for some reason: %s\n", rawMessage)
			}
			if messageType == websocket.BinaryMessage {
				log.Printf("Received binary message: %x\n", rawMessage)
				// Handle binary message
				// Assuming NewBinaryReader and its methods are defined in the current spacetimedb package
				// and behave similarly to the TypeScript SDK's BinaryReader.
				reader := NewBinaryReader(rawMessage)
				// Read the message type (first byte)
				messageType := reader.ReadByte()

				log.Printf("Message type: %d\n", messageType)
			}
			if messageType == websocket.CloseMessage {
				log.Println("Received close message, closing connection")
				return
			}
			if messageType == websocket.PongMessage {
				log.Println("Received pong message")
			}
		}
	}()
	return nil
}

func (db *DBConnection) Close() {
	if db.conn != nil {
		err := db.conn.Close()
		if err != nil {
			fmt.Printf("Error closing connection: %v\n", err)
		} else {
			fmt.Println("Connection closed")
		}
	}
}

// sendMessage sends a message over the websocket connection.
func (db *DBConnection) sendMessage(messageType int, data []byte) error {
	if db.conn == nil {
		return fmt.Errorf("cannot send message: not connected")
	}
	err := db.conn.WriteMessage(messageType, data)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	return nil
}

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
