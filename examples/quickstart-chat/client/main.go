package main

import (
	"fmt"

	"github.com/alexanderbh/spacetimedb-go-sdk"
	"github.com/alexanderbh/spacetimedb-go-sdk/examples/quickstart-chat/client/module_bindings"
)

func main() {
	db := spacetimedb.NewDBConnection(
		spacetimedb.WithHost("wss://maincloud.spacetimedb.com"),
		spacetimedb.WithNameOrIdentity("go-sdk-test"),
		spacetimedb.WithOnConnect(onConnect),
		spacetimedb.WithOnDisconnect(onDisconnect),
		spacetimedb.WithTableNameMap(module_bindings.Tables),
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

	err = module_bindings.SetName(conn, "Setname called! Updates?")

	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}

	err = conn.Subscribe("SELECT * FROM user")
	if err != nil {
		fmt.Println("Error subscribing to query:", err)
		return
	}
}

func onDisconnect(conn *spacetimedb.DBConnection) {
	fmt.Printf("Disconnected from database.\n")
}
