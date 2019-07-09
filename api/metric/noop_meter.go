package metric

import (
	"context"

	"github.com/open-telemetry/opentelemetry-go/api/core"
)

type noopMetric struct{}

func (noopMeter) GetFloat64Gauge(gauge *Float64GaugeRegistration, value float64, labels ...core.KeyValue) Float64Gauge {
	return noopMetric{}
}

func (noopMetric) Set(ctx context.Context, value float64, labels ...core.KeyValue) {
}
