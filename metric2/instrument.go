package metric2

import (
	"go.opentelemetry.io/otel/metric/number"
)

type Measurement struct {
	Instrument
	number.Number
}

type Instrument interface {
}
