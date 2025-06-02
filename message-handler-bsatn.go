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

	//fmt.Printf("Received message: %s\n\n", serverMsg)

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
		fmt.Printf("  Reducer:\t%s\n", msg.ReducerCall.String())
		switch status := msg.Status.Status.(type) {
		case *UpdateStatusComitted:
			fmt.Println("  Status:\tSuccess")
			db.handleTableUpdates(status.DatabaseUpdate.Tables)
		case *UpdateStatusFailed:
			fmt.Println("  Status:\tFailed")
			fmt.Printf("  Error:\t%s\n", status.ErrorMessage)
		}
	case *InitialSubscription:
		fmt.Printf("Received InitialSubscription:\n")
		fmt.Printf("  RequestId: %d\n", msg.RequestId)
		fmt.Printf("  TotalHostExecutionDuration: %s\n", msg.TotalHostExecutionDuration.String())
		if msg.DatabaseUpdate != nil && msg.DatabaseUpdate.Tables != nil {
			db.handleTableUpdates(msg.DatabaseUpdate.Tables)
		}
	}

	return nil
}

func (db *DBConnection) handleTableUpdates(updates []*TableUpdate) error {
	for _, tableUpdate := range updates {
		if tableUpdate == nil || tableUpdate.NumRows == 0 {
			continue
		}
		if db.TableNameMap[tableUpdate.TableName] == nil {
			return fmt.Errorf("table %s not found in TableNameMap", tableUpdate.TableName)
		}
		for _, update := range tableUpdate.Updates {
			if update == nil {
				continue
			}
			reader := NewBinaryReader(update.Inserts.RowsData)
			// While reader is not at the end loop through the rows
			for reader.offset < len(reader.buffer) {
				err := db.TableNameMap[tableUpdate.TableName].Insert(reader)
				if err != nil {
					return fmt.Errorf("error inserting row: %w", err)
				}
			}

			reader = NewBinaryReader(update.Deletes.RowsData)
			// While reader is not at the end loop through the rows
			for reader.offset < len(reader.buffer) {
				err := db.TableNameMap[tableUpdate.TableName].Delete(reader)
				if err != nil {
					return fmt.Errorf("error deleting row: %w", err)
				}
			}
		}
	}

	return nil
}
