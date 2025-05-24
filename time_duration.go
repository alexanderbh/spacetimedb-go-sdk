package spacetimedb

import "math/big"

// TimeDuration represents a difference between two points in time, in microseconds.
type TimeDuration struct {
	Micros *big.Int
}

var microsPerMillis = big.NewInt(1000)

// NewTimeDuration creates a new TimeDuration from microseconds.
func NewTimeDuration(micros *big.Int) *TimeDuration {
	return &TimeDuration{Micros: new(big.Int).Set(micros)}
}

// Micros returns the number of microseconds in the duration.
func (td *TimeDuration) GetMicros() *big.Int {
	return new(big.Int).Set(td.Micros)
}

// Millis returns the number of milliseconds in the duration (as int64, truncating any remainder).
func (td *TimeDuration) Millis() int64 {
	if td.Micros == nil {
		return 0
	}
	return new(big.Int).Div(td.Micros, microsPerMillis).Int64()
}

// FromMillis creates a new TimeDuration from milliseconds.
func FromMillis(millis int64) *TimeDuration {
	micros := new(big.Int).Mul(big.NewInt(millis), microsPerMillis)
	return NewTimeDuration(micros)
}
