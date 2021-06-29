package batch

import (
	"go.opentelemetry.io/otel/metric/number"
	"go.opentelemetry.io/otel/metric2/instrument"
)

type Measurement struct {
	instrument.Instrument
	number.Number
}
