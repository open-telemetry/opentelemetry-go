package async

import (
	"context"

	metric "go.opentelemetry.io/otel/metric2"
	asyncfloat64metric "go.opentelemetry.io/otel/metric2/async/float64"
	asyncint64metric "go.opentelemetry.io/otel/metric2/async/int64"
)

type Meter struct {
}

type Callback struct {
}

func (m Meter) Callback(func(context.Context), ...metric.Instrument) Callback {
	return Callback{}
}

func (m Meter) Integer() asyncint64metric.Meter {
	return asyncint64metric.Meter{}
}

func (m Meter) FloatingPoint() asyncfloat64metric.Meter {
	return asyncfloat64metric.Meter{}
}
