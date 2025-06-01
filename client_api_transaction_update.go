package spacetimedb

import "fmt"

type TransactionUpdate struct {
	Status                     *UpdateStatus
	Timestamp                  *Timestamp
	CallerIdentity             *Identity
	CallerConnectionId         *ConnectionId
	ReducerCall                *ReducerCallInfo
	EnergyQuantaUsed           *EnergyQuanta
	TotalHostExecutionDuration *TimeDuration
}

func (it *TransactionUpdate) Deserialize(reader *BinaryReader) error {

	it.Status = &UpdateStatus{}
	if err := it.Status.Deserialize(reader); err != nil {
		return fmt.Errorf("TransactionUpdate.Deserialize: failed to deserialize Status: %w", err)
	}

	it.Timestamp = NewTimestamp(reader.ReadI64())

	it.CallerIdentity = &Identity{}
	if err := it.CallerIdentity.Deserialize(reader); err != nil {
		return fmt.Errorf("TransactionUpdate.Deserialize: failed to deserialize CallerIdentity: %w", err)
	}
	it.CallerConnectionId = &ConnectionId{}
	if err := it.CallerConnectionId.Deserialize(reader); err != nil {
		return fmt.Errorf("TransactionUpdate.Deserialize: failed to deserialize CallerConnectionId: %w", err)
	}
	it.ReducerCall = &ReducerCallInfo{}
	if err := it.ReducerCall.Deserialize(reader); err != nil {
		return fmt.Errorf("TransactionUpdate.Deserialize: failed to deserialize ReducerCall: %w", err)
	}
	it.EnergyQuantaUsed = &EnergyQuanta{}
	if err := it.EnergyQuantaUsed.Deserialize(reader); err != nil {
		return fmt.Errorf("TransactionUpdate.Deserialize: failed to deserialize EnergyQuantaUsed: %w", err)
	}
	it.TotalHostExecutionDuration = NewTimeDuration(reader.ReadI64())

	return nil
}
