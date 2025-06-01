package spacetimedb

// TimeDuration represents a difference between two points in time, in microseconds.
type TimeDuration struct {
	Micros int64
}

var microsPerMillis int64 = 1000

func NewTimeDuration(micros int64) *TimeDuration {
	return &TimeDuration{Micros: micros}
}

func (td *TimeDuration) Millis() int64 {
	return td.Micros / microsPerMillis
}

func FromMillis(millis int64) *TimeDuration {
	micros := millis * microsPerMillis
	return NewTimeDuration(micros)
}
