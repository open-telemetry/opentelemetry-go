package metric2

import (
	"go.opentelemetry.io/otel/metric/number"
)

// Measurement represents a single measurement made on a specific
// instrument.  The encapsulated number's Kind matches the instrument
// definition.
type Measurement struct {
	instrument Instrument
	number     number.Number
}

// Instrument is a base
type Instrument interface {
	// Implementation returns the underlying implementation of the
	// instrument, which allows the SDK to gain access to its own
	// representation especially from a `Measurement`.
	// Implementation() interface{}

	// Descriptor returns a copy of the instrument's Descriptor.
	// Descriptor() metric.Descriptor
}
