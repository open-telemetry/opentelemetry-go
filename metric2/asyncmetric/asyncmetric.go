package asyncmetric

import (
	"context"

	"go.opentelemetry.io/otel/metric2/instrument"
	"go.opentelemetry.io/otel/metric2/asyncmetric/asyncfloat64metric"
	"go.opentelemetry.io/otel/metric2/asyncmetric/asyncint64metric"
)

type Meter struct {
}

type Callback struct {
}

func (m Meter) Callback(func(context.Context), ...instrument.Instrument) Callback {
	return Callback{}
}

func (m Meter) Integer() asyncint64metric.Meter {
	return asyncint64metric.Meter{}
}

func (m Meter) FloatingPoint() asyncfloat64metric.Meter {
	return asyncfloat64metric.Meter{}
}
