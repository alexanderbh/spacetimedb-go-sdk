package module_bindings

import (
	"fmt"

	"github.com/alexanderbh/spacetimedb-go-sdk"
)

func SetName(conn *spacetimedb.DBConnection, name string) error {
	writer := spacetimedb.NewBinaryWriter(1024)
	writer.WriteString(name)
	err := conn.CallReducer(
		"set_name",
		writer.GetBuffer(),
		0,
		0,
	)
	if err != nil {
		return fmt.Errorf("Error calling set_name reducer: %v\n", err)
	}
	return nil
}
