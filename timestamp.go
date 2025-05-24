package spacetimedb

import (
	"fmt"
	"math/big"
	"time"
)

const (
	microsPerMilli = 1000
	// MaxGoTimeMicros is the maximum number of microseconds that can be represented by Go's time.Time.
	// This corresponds to time.Unix(0, math.MaxInt64).UnixMicro() but calculated directly.
	// math.MaxInt64 / 1000 (for nanos to micros)
	maxGoTimeMicros = 9223372036854775 // (math.MaxInt64 / 1000) roughly, actual max time.UnixMicro()
	// MinGoTimeMicros is the minimum number of microseconds that can be represented by Go's time.Time.
	// This corresponds to time.Unix(0, math.MinInt64).UnixMicro()
	minGoTimeMicros = -9223372036854776 // (math.MinInt64 / 1000) roughly, actual min time.UnixMicro()
)

// Timestamp represents a point in time, represented as a number of microseconds since the Unix epoch.
type Timestamp struct {
	microsSinceUnixEpoch int64
}

// NewTimestamp creates a new Timestamp.
func NewTimestamp(micros int64) *Timestamp {
	return &Timestamp{microsSinceUnixEpoch: micros}
}

// MicrosSinceUnixEpoch returns the number of microseconds since the Unix epoch.
func (t *Timestamp) MicrosSinceUnixEpoch() int64 {
	return t.microsSinceUnixEpoch
}

// UnixEpoch is the midnight at the beginning of January 1, 1970, UTC.
var UnixEpoch = NewTimestamp(0)

// Now returns a Timestamp representing the current moment in time.
func Now() *Timestamp {
	return FromDate(time.Now())
}

// FromDate returns a Timestamp representing the same point in time as date.
func FromDate(date time.Time) *Timestamp {
	micros := date.UnixMicro()
	return NewTimestamp(micros)
}

// ToDate returns a time.Time representing approximately the same point in time as this Timestamp.
// Returns an error if the Timestamp is outside the range representable by Go's time.Time.
func (t *Timestamp) ToDate() (time.Time, error) {
	if t.microsSinceUnixEpoch > maxGoTimeMicros || t.microsSinceUnixEpoch < minGoTimeMicros {
		return time.Time{}, fmt.Errorf("Timestamp %d is outside of the representable range of Go's time.Time", t.microsSinceUnixEpoch)
	}
	// time.UnixMicro directly takes microseconds
	return time.UnixMicro(t.microsSinceUnixEpoch), nil
}

// ToBigInt returns the timestamp as a *big.Int representing microseconds since the Unix epoch.
// This is for compatibility with the original TS version's bigint type if needed for other operations,
// though the primary storage is int64.
func (t *Timestamp) ToBigInt() *big.Int {
	return big.NewInt(t.microsSinceUnixEpoch)
}

// FromBigInt creates a Timestamp from a *big.Int representing microseconds since the Unix epoch.
// Returns an error if the value is out of int64 range.
func FromBigInt(micros *big.Int) (*Timestamp, error) {
	if !micros.IsInt64() {
		return nil, fmt.Errorf("big.Int value %s is out of range for int64", micros.String())
	}
	return NewTimestamp(micros.Int64()), nil
}
