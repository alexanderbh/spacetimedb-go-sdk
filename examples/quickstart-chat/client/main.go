package main

import (
	"log"
	"os"
	"quickstart-chat/module_bindings"

	"github.com/alexanderbh/bubbleapp/app"
	"github.com/alexanderbh/spacetimedb-go-sdk"
	tea "github.com/charmbracelet/bubbletea/v2"
)

func main() {
	db := spacetimedb.NewDBConnection(
		spacetimedb.WithHost("wss://maincloud.spacetimedb.com"),
		spacetimedb.WithNameOrIdentity("go-sdk-test"),
		spacetimedb.WithOnConnect(onConnect),
		spacetimedb.WithOnDisconnect(onDisconnect),
		spacetimedb.WithTableNameMap(module_bindings.Tables),
		spacetimedb.WithLogger(Logger),
	)
	err := db.Connect()
	if err != nil {
		log.Fatalln("Error connecting to database:", err)
		return
	}
	defer db.Close()

	c := app.NewCtx()
	bubbleApp := app.New(c, NewRoot(db))
	p := tea.NewProgram(bubbleApp, tea.WithAltScreen(), tea.WithMouseAllMotion())
	bubbleApp.SetTeaProgram(p)
	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}

func onConnect(db *spacetimedb.DBConnection, identity *spacetimedb.Identity, token string, connectionId *spacetimedb.ConnectionId) {

	db.Logger("Connected to database with identity: %s", identity.ToHexString())
	db.Logger("Token: %s", token)
	connId, err := connectionId.ToHexString()
	if err != nil {
		db.Logger("Error converting connection ID to hex string:", err)
	} else {
		db.Logger("Connection ID: %s", connId)
	}

	err = module_bindings.SetName(db, "Setname called with this")

	if err != nil {
		log.Println("Error sending message:", err)
		return
	}

	err = db.Subscribe("SELECT * FROM user")
	if err != nil {
		log.Println("Error subscribing to query:", err)
		return
	}
}

func onDisconnect(db *spacetimedb.DBConnection) {
	db.Logger("Disconnected from database.")
}
