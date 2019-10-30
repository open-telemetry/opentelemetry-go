package prometheus

import (
	"context"

	"go.opentelemetry.io/sdk/export"
)

// Exporter is an implementation of metric.Exporter that sends metrics to
// Prometheus.
type Exporter struct {
}

var _ export.MetricBatcher = (*Exporter)(nil)

// NewExporter returns a new prometheus exporter.
func NewExporter() *Exporter {
	return &Exporter{}
}

// AggregatorFor returns the metric aggregator used for the particular exporter.
func (e *Exporter) AggregatorFor(record export.MetricRecord) export.MetricAggregator {
	return nil
}

// Export exports the provide metric record to prometheus.
func (e *Exporter) Export(
	ctx context.Context,
	record export.MetricRecord,
	aggregator export.MetricAggregator) {

}
