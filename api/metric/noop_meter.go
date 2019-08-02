package metric

import (
	"context"

	"go.opentelemetry.io/api/core"
)

type noopMeter struct{}

type noopMetric struct{}

var _ Meter = noopMeter{}

var _ Float64Gauge = noopMetric{}

func (noopMeter) GetFloat64Gauge(ctx context.Context, gauge *Float64GaugeHandle, labels ...core.KeyValue) Float64Gauge {
	return noopMetric{}
}

func (noopMetric) Set(ctx context.Context, value float64, labels ...core.KeyValue) {
}
