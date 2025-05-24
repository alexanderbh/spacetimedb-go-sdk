package spacetimedb

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

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

	identity     *Identity
	token        string
	connectionId *ConnectionId
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
	//dialer.Subprotocols = []string{"v1.json.spacetimedb"}
	dialer.Subprotocols = []string{"v1.bsatn.spacetimedb"}
	c, _, err := dialer.Dial(db.host, nil)

	if err != nil {
		return fmt.Errorf("failed to connect to websocket: %w", err)
	}
	db.conn = c
	db.isAlive = true
	fmt.Printf("Connected to websocket at %s\n", db.host)

	go func() {
		defer func() {
			db.isAlive = false
			if db.conn != nil {
				db.conn.Close()
			}
		}()
		for {
			if db.conn == nil {
				log.Println("connection is nil, exiting message read loop")
				return
			}
			messageType, rawMessage, err := db.conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading message: %v\n\n", err)
				return
			}
			if messageType == websocket.TextMessage {
				log.Printf("Received text message: %s\n\n", rawMessage)
				db.parseJsonMessage(rawMessage)
			}
			if messageType == websocket.BinaryMessage {
				log.Printf("Received binary message: %x\n\n", rawMessage)
				err = db.parseBsantMessage(rawMessage)
				if err != nil {
					log.Printf("Error parsing binary message: %v\n\n", err)
				}
			}
			if messageType == websocket.CloseMessage {
				log.Println("Received close message, closing connection")
				return
			}
			if messageType == websocket.PongMessage {
				log.Print("Received pong message\n\n")
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
