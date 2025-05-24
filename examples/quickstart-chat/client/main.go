package main

import (
	"fmt"

	"github.com/alexanderbh/spacetimedb-go-sdk"
)

func main() {

	db := spacetimedb.NewDBConnection("wss://maincloud.spacetimedb.com/v1/database/go-sdk-test/subscribe?compression=None")
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
