package asyncmetric

import (
	"context"

	"go.opentelemetry.io/otel/metric2/asyncmetric/asyncfloat64metric"
	"go.opentelemetry.io/otel/metric2/asyncmetric/asyncint64metric"
)

type Meter struct {
}

type Callback struct {
}

type Instrument interface {
}

// GOAL: Use RecordBatch inside async contexts, or individual Add()/Observe().
// RecordBatch uses structs, so we only need to define async singletons here.

func (m Meter) Callback(func(context.Context), ...Instrument) Callback {
	return Callback{}
}

func (m Meter) Integer() asyncint64metric.Meter {
	return asyncint64metric.Meter{}
}

func (m Meter) FloatingPoint() asyncfloat64metric.Meter {
	return asyncfloat64metric.Meter{}
}
