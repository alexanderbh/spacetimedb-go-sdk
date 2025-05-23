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
	dialer.Subprotocols = []string{"v1.json.spacetimedb"}
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
				log.Printf("Received text message: %s\n", rawMessage)
			}
			if messageType == websocket.BinaryMessage {
				log.Printf("Received binary message: %x\n", rawMessage)
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
