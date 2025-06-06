package spacetimedb

import (
	"context"
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
)

type DBConnection struct {
	Host           string
	NameOrIdentity string
	IsConnected    bool
	WS             *websocket.Conn

	ctx    context.Context
	cancel context.CancelFunc

	Identity     *Identity
	Token        string
	ConnectionId *ConnectionId

	Compression uint8

	TableNameMap TableNameMap

	OnConnect    func(conn *DBConnection, identity *Identity, token string, connectionId *ConnectionId)
	OnDisconnect func(*DBConnection)

	Logger func(format string, args ...interface{})
}

const (
	CompressionTypeNone   uint8 = 0
	CompressionTypeBrotly uint8 = 1
	CompressionTypeGzip   uint8 = 2
)

func NewDBConnection(opts ...DBConnectionOption) *DBConnection {
	conn := &DBConnection{
		Host: "wss://maincloud.spacetimedb.com",
	}

	for _, opt := range opts {
		opt(conn)
	}

	return conn
}

type DBConnectionOption func(*DBConnection)

func WithHost(host string) DBConnectionOption {
	return func(opts *DBConnection) {
		opts.Host = host
	}
}
func WithNameOrIdentity(nameOrIdentity string) DBConnectionOption {
	return func(opts *DBConnection) {
		opts.NameOrIdentity = nameOrIdentity
	}
}
func WithOnConnect(onConnect func(conn *DBConnection, identity *Identity, token string, connectionId *ConnectionId)) DBConnectionOption {
	return func(opts *DBConnection) {
		opts.OnConnect = onConnect
	}
}
func WithOnDisconnect(onDisconnect func(*DBConnection)) DBConnectionOption {
	return func(opts *DBConnection) {
		opts.OnDisconnect = onDisconnect
	}
}
func WithTableNameMap(tableNameMap TableNameMap) DBConnectionOption {
	return func(opts *DBConnection) {
		opts.TableNameMap = tableNameMap
	}
}
func WithLogger(logger func(format string, args ...interface{})) DBConnectionOption {
	return func(opts *DBConnection) {
		opts.Logger = logger
	}
}
func WithStdLogger() DBConnectionOption {
	return func(opts *DBConnection) {
		opts.Logger = func(format string, args ...interface{}) {
			fmt.Printf(format+"\n", args...)
		}
	}
}

func (db *DBConnection) Connect() error {
	if db.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}

	db.ctx, db.cancel = context.WithCancel(context.Background())

	dialer := websocket.DefaultDialer
	//dialer.Subprotocols = []string{"v1.json.spacetimedb"}
	dialer.Subprotocols = []string{"v1.bsatn.spacetimedb"}
	url, err := url.JoinPath(db.Host, "v1", "database", db.NameOrIdentity, "subscribe")

	if err != nil {
		return fmt.Errorf("failed to join URL path: %w", err)
	}

	c, _, err := dialer.DialContext(db.ctx, url+"?compression=None", nil) // TODO: Add compression support
	if err != nil {
		return fmt.Errorf("failed to connect to websocket: %w", err)
	}

	db.WS = c
	db.Logger("Connected to websocket at %s", db.Host)

	go func() {
		defer func() {
			db.IsConnected = false
			if db.WS != nil {
				db.WS.Close()
			}
		}()
		for {
			select {
			case <-db.ctx.Done():
				db.Logger("context cancelled, exiting message read loop")
				return
			default:
				if db.WS == nil {
					db.Logger("connection is nil, exiting message read loop")
					return
				}
				messageType, rawMessage, err := db.WS.ReadMessage()
				if err != nil {
					select {
					case <-db.ctx.Done():
						db.Logger("context cancelled, exiting message read loop after read error")
					default:
						db.Logger("Error reading message: %v", err)
					}
					return
				}
				if messageType == websocket.TextMessage {
					db.Logger("Received text message: %s", rawMessage)
				}
				if messageType == websocket.BinaryMessage {
					db.Logger("Received binary message: %x", rawMessage)
					err = db.parseBsantMessage(rawMessage)
					if err != nil {
						db.Logger("Error parsing binary message: %v", err)
					}
				}
				if messageType == websocket.CloseMessage {
					db.Logger("Received close message, closing connection")
					return
				}
				if messageType == websocket.PongMessage {
					db.Logger("Received pong message")
				}
			}
		}
	}()
	return nil
}

func (db *DBConnection) Close() {
	if db.cancel != nil {
		db.cancel()
	}
	if db.WS != nil {
		err := db.WS.Close()
		if err != nil {
			db.Logger("Error closing connection: %v", err)
		} else {
			db.Logger("Connection closed")
		}
	}
}

func (db *DBConnection) SendMessage(data []byte) error {
	if db.WS == nil {
		return fmt.Errorf("cannot send message: not connected")
	}
	err := db.WS.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	return nil
}
