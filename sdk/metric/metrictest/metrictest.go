package metrictest

import (
	"context"

	"go.opentelemetry.io/otel/sdk/metric/reader"
)

type Exporter struct {
	reader.Producer
}

var _ reader.Exporter = &Exporter{}

func NewExporter() *Exporter {
	return &Exporter{}
}

func (t *Exporter) Register(producer reader.Producer) {
	t.Producer = producer
}

func (*Exporter) Flush(context.Context) error { return nil }

func (*Exporter) Shutdown(context.Context) error { return nil }
