package spacetimedb

import "fmt"

type InitialSubscription struct {
	DatabaseUpdate             *DatabaseUpdate
	RequestId                  uint32
	TotalHostExecutionDuration *TimeDuration
}

func (it *InitialSubscription) Deserialize(reader *BinaryReader) error {

	it.DatabaseUpdate = &DatabaseUpdate{}
	if err := it.DatabaseUpdate.Deserialize(reader); err != nil {
		return fmt.Errorf("failed to deserialize DatabaseUpdate: %w", err)
	}
	it.RequestId = reader.ReadU32()
	it.TotalHostExecutionDuration = NewTimeDuration(reader.ReadI64())

	return nil
}
