package spacetimedb

import "fmt"

type UpdateStatus struct {
	Status any // Union type
}

type UpdateStatusComitted struct {
	DatabaseUpdate *DatabaseUpdate
}

type UpdateStatusFailed struct {
	ErrorMessage string
}

type UpdateStatusOutOfEnergy struct {
}

func (it *UpdateStatus) Deserialize(reader *BinaryReader) error {
	unionType := reader.ReadU8()
	switch unionType {
	case 0x00:
		status := &UpdateStatusComitted{}
		status.DatabaseUpdate = &DatabaseUpdate{}
		if err := status.DatabaseUpdate.Deserialize(reader); err != nil {
			return err
		}
		it.Status = status
	case 0x01:
		failed := &UpdateStatusFailed{}
		failed.ErrorMessage = reader.ReadString()
		it.Status = failed
	case 0x02:
		it.Status = &UpdateStatusOutOfEnergy{}
	default:
		return fmt.Errorf("UpdateStatus.Deserialize: unknown union type 0x%02x", unionType)
	}
	return nil
}

func (it *UpdateStatus) String() string {
	switch status := it.Status.(type) {
	case *UpdateStatusComitted:
		return fmt.Sprintf("UpdateStatusComitted: %s", status.DatabaseUpdate.String())
	case *UpdateStatusFailed:
		return fmt.Sprintf("UpdateStatusFailed: %s", status.ErrorMessage)
	case *UpdateStatusOutOfEnergy:
		return "UpdateStatusOutOfEnergy"
	default:
		return fmt.Sprintf("Unknown UpdateStatus type: %T", it.Status)
	}
}
