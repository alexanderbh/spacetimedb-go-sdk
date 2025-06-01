package spacetimedb

import (
	"fmt"
)

func (db *DBConnection) parseBsantMessage(msg []byte) error {
	reader := NewBinaryReader(msg)

	// Handle compression
	compression := reader.ReadU8()
	switch compression {
	case CompressionTypeNone:
		// No compression, read the message directly
	case CompressionTypeBrotly:
		return fmt.Errorf("Brotly compression is not implemented yet")
	case CompressionTypeGzip:
		return fmt.Errorf("Gzip compression is not implemented yet")
	default:
		return fmt.Errorf("unknown compression type: %d", compression)
	}

	// Read the message type
	serverMsg := &ServerMessage{}
	if err := serverMsg.Deserialize(reader); err != nil {
		return fmt.Errorf("failed to deserialize server message: %w", err)
	}

	fmt.Printf("Received message: %s\n\n", serverMsg)

	switch msg := serverMsg.Message.(type) {
	case *IdentityToken:

		fmt.Printf("Received IdentityToken: %#v\n\n", serverMsg.Message)

		db.IsConnected = true
		db.Identity = msg.Identity
		if db.Token == "" && msg.Token != "" {
			db.Token = msg.Token
		}
		db.ConnectionId = msg.ConnectionId
		if db.OnConnect != nil {
			db.OnConnect(db, msg.Identity, msg.Token, msg.ConnectionId)
		}
	case *TransactionUpdate:
		fmt.Printf("Received TransactionUpdate:\n")
		fmt.Printf("  Status: %s\n", msg.Status.String())
		fmt.Printf("  Timestamp: %#v\n", msg.Timestamp)
		fmt.Printf("  CallerIdentity: %#v\n", msg.CallerIdentity)
		fmt.Printf("  CallerConnectionId: %#v\n", msg.CallerConnectionId)
		fmt.Printf("  ReducerCall: %#v\n", msg.ReducerCall)
		fmt.Printf("  EnergyQuantaUsed: %#v\n", msg.EnergyQuantaUsed)
		fmt.Printf("  TotalHostExecutionDuration: %#v\n", msg.TotalHostExecutionDuration)
	case *InitialSubscription:
		fmt.Printf("Received InitialSubscription:\n")
		fmt.Printf("  RequestId: %d\n", msg.RequestId)
		fmt.Printf("  TotalHostExecutionDuration: %s\n", msg.TotalHostExecutionDuration.String())
		if msg.DatabaseUpdate != nil {
			fmt.Printf("  %s\n", msg.DatabaseUpdate.String())
			for _, tableUpdate := range msg.DatabaseUpdate.Tables {
				if tableUpdate == nil || tableUpdate.NumRows == 0 {
					continue
				}
				for _, update := range tableUpdate.Updates {
					if update == nil {
						continue
					}
					if db.TableNameMap[tableUpdate.TableName] == nil {
						fmt.Printf("  Table %s not found in TableNameMap\n", tableUpdate.TableName)
						continue
					}
					reader := NewBinaryReader(update.Inserts.RowsData)
					// While reader is not at the end loop through the rows
					for reader.offset < len(reader.buffer) {
						row, err := db.TableNameMap[tableUpdate.TableName].DeserializeRow(reader)
						if err != nil {
							fmt.Println("  Error deserializing row:", err)
							continue
						}
						err = db.TableNameMap[tableUpdate.TableName].Insert(row)
						if err != nil {
							fmt.Println("  Error inserting row:", err)
							continue
						}
					}
				}
			}
		} else {
			fmt.Println("  DatabaseUpdate is nil")
		}
	}

	return nil
}
