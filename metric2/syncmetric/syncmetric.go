package syncmetric

import (
	syncfloat64metric "go.opentelemetry.io/otel/metric2/syncmetric/float64"
	syncint64metric "go.opentelemetry.io/otel/metric2/syncmetric/int64"
)

type Meter struct{}

func (m Meter) Integer() syncint64metric.Meter {
	return syncint64metric.Meter{}
}

func (m Meter) FloatingPoint() syncfloat64metric.Meter {
	return syncfloat64metric.Meter{}
}
