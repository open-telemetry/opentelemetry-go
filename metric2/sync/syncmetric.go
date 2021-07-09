package sync

import (
	syncfloat64metric "go.opentelemetry.io/otel/metric2/sync/float64"
	syncint64metric "go.opentelemetry.io/otel/metric2/sync/int64"
)

type Meter struct{}

func (m Meter) Integer() syncint64metric.Meter {
	return syncint64metric.Meter{}
}

func (m Meter) FloatingPoint() syncfloat64metric.Meter {
	return syncfloat64metric.Meter{}
}
