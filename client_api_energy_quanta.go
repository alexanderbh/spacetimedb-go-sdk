package spacetimedb

import "math/big"

type EnergyQuanta struct {
	Quanta big.Int
}

func (it *EnergyQuanta) Deserialize(reader *BinaryReader) error {

	it.Quanta = *reader.ReadU128()

	return nil
}
