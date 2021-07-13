package sdkapi

import "go.opentelemetry.io/otel/metric/number"

// Measurement represents a single measurement made on a specific
// instrument.  The encapsulated number's Kind matches the instrument
// definition.
type Measurement struct {
	instrument Instrument
	number     number.Number
}

type Instrument interface {
}
