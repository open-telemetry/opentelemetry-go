package metric

import (
	"context"

	"github.com/open-telemetry/opentelemetry-go/api/core"
)

type noopMetric struct{}

func (noopMeter) GetFloat64Gauge(ctx context.Context, gauge *Float64GaugeHandle, labels ...core.KeyValue) Float64Gauge {
	return noopMetric{}
}

func (noopMetric) Set(ctx context.Context, value float64, labels ...core.KeyValue) {
}
