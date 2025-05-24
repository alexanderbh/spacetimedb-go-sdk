package main

import (
	"fmt"

	"github.com/alexanderbh/spacetimedb-go-sdk"
)

func main() {
	db := spacetimedb.NewDBConnection(
		spacetimedb.WithHost("wss://maincloud.spacetimedb.com"),
		spacetimedb.WithNameOrIdentity("go-sdk-test"),
		spacetimedb.WithOnConnect(onConnect),
		spacetimedb.WithOnDisconnect(onDisconnect),
	)
	err := db.Connect()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()

	fmt.Print("Press Enter to exit...\n\n")
	var input string
	fmt.Scanln(&input)
}

func onConnect(conn *spacetimedb.DBConnection, identity *spacetimedb.Identity, token string, connectionId *spacetimedb.ConnectionId) {
	fmt.Printf("Connected to database with identity: %s\n", identity.ToHexString())
	fmt.Printf("Token: %s\n", token)
	connId, err := connectionId.ToHexString()
	if err != nil {
		fmt.Println("Error converting connection ID to hex string:", err)
	} else {
		fmt.Printf("Connection ID: %s\n", connId)
	}

	argsWriter := spacetimedb.NewBinaryWriter()
	spacetimedb.CreateStringType().Serialize(argsWriter, "NameTest")

	taggedType, err := spacetimedb.NewMapFromClientMessage(&spacetimedb.CallReducer{
		Reducer:   "set_name",
		Args:      argsWriter.GetBuffer(),
		RequestId: 1,
		Flags:     0,
	})
	if err != nil {
		fmt.Println("Error creating CallReducer tagged type:", err)
		return
	}

	writer := spacetimedb.NewBinaryWriter()
	spacetimedb.ClientMessage_GetAlgebraicType().Serialize(writer, taggedType)

	msg := writer.GetBuffer()
	err = conn.SendMessage(msg)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}
}

func onDisconnect(conn *spacetimedb.DBConnection) {
	fmt.Printf("Disconnected from database.\n")
}
