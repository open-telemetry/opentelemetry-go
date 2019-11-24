// Copyright 2019, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prometheus

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.opentelemetry.io/otel/api/core"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
)

type metricKey struct {
	desc    *export.Descriptor
	encoded string
}

// Exporter is an implementation of metric.Exporter that sends metrics to
// Prometheus.
type Exporter struct {
	handler http.Handler

	registerer prometheus.Registerer
	gatherer   prometheus.Gatherer

	counters           counters
	gauges             gauges
	histograms         histograms
	summaries          summaries
	measureAggregation MeasureAggregation
}

var _ export.Exporter = &Exporter{}
var _ http.Handler = &Exporter{}

// Options is a set of options for the tally reporter.
type Options struct {
	// Registry is the prometheus registry that will be used as the default Registerer and
	// Gatherer if these are not specified.
	//
	// If not set a new empty Registry is created.
	Registry *prometheus.Registry

	// Registerer is the prometheus registerer to register
	// metrics with.
	//
	// If not specified the Registry will be used as default.
	Registerer prometheus.Registerer

	// Gatherer is the prometheus gatherer to gather
	// metrics with.
	//
	// If not specified the Registry will be used as default.
	Gatherer prometheus.Gatherer

	// DefaultHistogramBuckets is the default histogram buckets
	// to use. Use nil to specify the system-default histogram buckets.
	DefaultHistogramBuckets []float64

	// DefaultSummaryObjectives is the default summary objectives
	// to use. Use nil to specify the system-default summary objectives.
	DefaultSummaryObjectives map[float64]float64

	// MeasureAggregation defines how metric.Measure are exported.
	// Possible values are 'Histogram' or 'Summary'.
	// The default export representation for measures is Histograms.
	MeasureAggregation MeasureAggregation
}

type MeasureAggregation int

const (
	Histogram MeasureAggregation = iota
	Summary
)

// NewExporter returns a new prometheus exporter for prometheus metrics.
func NewExporter(opts Options) (*Exporter, error) {
	if opts.Registry == nil {
		opts.Registry = prometheus.NewRegistry()
	}

	if opts.Registerer == nil {
		opts.Registerer = opts.Registry
	}

	if opts.Gatherer == nil {
		opts.Gatherer = opts.Registry
	}

	return &Exporter{
		registerer:         opts.Registerer,
		gatherer:           opts.Gatherer,
		handler:            promhttp.HandlerFor(opts.Gatherer, promhttp.HandlerOpts{}),
		measureAggregation: opts.MeasureAggregation,

		counters:   newCounters(opts.Registerer),
		gauges:     newGauges(opts.Registerer),
		histograms: newHistograms(opts.Registerer, opts.DefaultHistogramBuckets),
		summaries:  newSummaries(opts.Registerer, opts.DefaultSummaryObjectives),
	}, nil
}

// Export exports the provide metric record to prometheus.
func (e *Exporter) Export(_ context.Context, checkpointSet export.CheckpointSet) error {
	var forEachError error
	checkpointSet.ForEach(func(record export.Record) {
		agg := record.Aggregator()

		mKey := metricKey{
			desc:    record.Descriptor(),
			encoded: record.Labels().Encoded(),
		}

		if points, ok := agg.(aggregator.Points); ok {
			observerExporter := e.histograms.export
			if e.measureAggregation == Summary {
				observerExporter = e.summaries.export
			}

			err := observerExporter(points, record, mKey)
			if err != nil {
				forEachError = err
			}
			return
		}

		if sum, ok := agg.(aggregator.Sum); ok {
			err := e.counters.export(sum, record, mKey)
			if err != nil {
				forEachError = err
			}
			return
		}

		if gauge, ok := agg.(aggregator.LastValue); ok {
			err := e.gauges.export(gauge, record, mKey)
			if err != nil {
				forEachError = err
			}
			return
		}
	})

	return forEachError
}

func (e *Exporter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.handler.ServeHTTP(w, r)
}

func labelsKeys(kvs []core.KeyValue) []string {
	keys := make([]string, 0, len(kvs))
	for _, kv := range kvs {
		keys = append(keys, sanitize(string(kv.Key)))
	}
	return keys
}

func labelValues(labels export.Labels) []string {
	// TODO(paivagustavo): parse the labels.Encoded() instead of calling `Emit()` directly
	//  this would avoid unnecessary allocations.
	values := make([]string, 0, labels.Len())
	for _, label := range labels.Ordered() {
		values = append(values, label.Value.Emit())
	}
	return values
}
