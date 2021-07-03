package time

import (
	lib "time"
)

type Clock interface {
	Now() lib.Time
	Since(t lib.Time) lib.Duration
}

type DefaultClock struct {
}

func (c DefaultClock) Now() lib.Time {
	return lib.Now()
}
func (c DefaultClock) Since(t lib.Time) lib.Duration {
	return lib.Since(t)
}

var _ Clock = DefaultClock{}
