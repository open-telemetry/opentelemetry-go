package prometheus

import (
	"context"
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	export "go.opentelemetry.io/otel/sdk/export/metric"
)

type metricID string

// Exporter is an implementation of metric.Exporter that sends metrics to
// Prometheus.
type Exporter struct {
	sync.RWMutex
	registerer      prometheus.Registerer
	gatherer        prometheus.Gatherer
	onRegisterError func(e error)
	counters        map[metricID]*prometheus.CounterVec
	gauges          map[metricID]*prometheus.GaugeVec
}

var _ export.Batcher = (*Exporter)(nil)

// Options is a set of options for the tally reporter.
type Options struct {
	// Registerer is the prometheus registerer to register
	// metrics with. Use nil to specify the default registerer.
	Registerer prometheus.Registerer

	// Gatherer is the prometheus gatherer to gather
	// metrics with. Use nil to specify the default gatherer.
	Gatherer prometheus.Gatherer

	// DefaultHistogramBuckets is the default histogram buckets
	// to use. Use nil to specify the default histogram buckets.
	DefaultHistogramBuckets []float64

	// DefaultSummaryObjectives is the default summary objectives
	// to use. Use nil to specify the default summary objectives.
	DefaultSummaryObjectives map[float64]float64

	// OnRegisterError defines a method to call to when registering
	// a metric with the registerer fails. Use nil to specify
	// to panic by default when registering fails.
	OnRegisterError func(err error)
}

// NewExporter returns a new prometheus exporter for prometheus metrics.
func NewExporter(opts Options) *Exporter {
	if opts.Registerer == nil {
		opts.Registerer = prometheus.DefaultRegisterer
	} else {
		// A specific registerer was set, check if it's a registry and if
		// no gatherer was set, then use that as the gatherer
		if reg, ok := opts.Registerer.(*prometheus.Registry); ok && opts.Gatherer == nil {
			opts.Gatherer = reg
		}
	}
	if opts.Gatherer == nil {
		opts.Gatherer = prometheus.DefaultGatherer
	}
	if opts.OnRegisterError == nil {
		opts.OnRegisterError = func(err error) {
			panic(err)
		}
	}

	return &Exporter{
		registerer:      opts.Registerer,
		gatherer:        opts.Gatherer,
		onRegisterError: opts.OnRegisterError,
		counters:        make(map[metricID]*prometheus.CounterVec),
		gauges:          make(map[metricID]*prometheus.GaugeVec),
	}
}

// AggregatorFor returns the metric aggregator used for the particular exporter.
func (e *Exporter) AggregatorFor(record export.Record) export.Aggregator {
	return nil
}

// Export exports the provide metric record to prometheus.
func (e *Exporter) Export(
	ctx context.Context,
	record export.Record,
	aggregator export.Aggregator) {

}
