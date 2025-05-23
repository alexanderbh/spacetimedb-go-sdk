package main

import (
	"fmt"

	"github.com/alexanderbh/spacetimedb-go-sdk"
)

func main() {

	db := spacetimedb.NewDBConnection("wss://maincloud.spacetimedb.com/v1/database/go-sdk-test/subscribe")
	err := db.Connect()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()

	fmt.Println("Press Enter to exit...")
	var input string
	fmt.Scanln(&input)
}
