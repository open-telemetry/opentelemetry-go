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
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"sync"

	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"

	"github.com/prometheus/client_golang/prometheus"

	"go.opentelemetry.io/otel/api/core"
	export "go.opentelemetry.io/otel/sdk/export/metric"
)

const (
	prefixSplitter = "+"
)

type metricKey struct {
	desc    *export.Descriptor
	encoded string
}

// Exporter is an implementation of metric.Exporter that sends metrics to
// Prometheus.
type Exporter struct {
	sync.RWMutex

	handler http.Handler

	registerer prometheus.Registerer
	gatherer   prometheus.Gatherer

	counters map[metricKey]prometheus.Counter
	gauges   map[metricKey]prometheus.Gauge

	counterVecs map[*export.Descriptor]*prometheus.CounterVec
	gaugeVecs   map[*export.Descriptor]*prometheus.GaugeVec
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
}

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
		registerer: opts.Registerer,
		gatherer:   opts.Gatherer,
		handler:    promhttp.HandlerFor(opts.Gatherer, promhttp.HandlerOpts{}),

		counters: make(map[metricKey]prometheus.Counter),
		gauges:   make(map[metricKey]prometheus.Gauge),

		counterVecs: make(map[*export.Descriptor]*prometheus.CounterVec),
		gaugeVecs:   make(map[*export.Descriptor]*prometheus.GaugeVec),
	}, nil
}

// Export exports the provide metric record to prometheus.
func (e *Exporter) Export(_ context.Context, checkpointSet export.CheckpointSet) error {
	var forEachError error
	checkpointSet.ForEach(func(record export.Record) {
		agg := record.Aggregator()

		desc := record.Descriptor()
		mKey := metricKey{
			desc:    desc,
			encoded: record.Labels().Encoded(),
		}

		if sum, ok := agg.(aggregator.Sum); ok {
			value, err := sum.Sum()
			if err != nil {
				// TODO: handle this better when we have a more
				//  sophisticated error handler mechanism for this ForEach method.
				forEachError = err
				return
			}

			c, err := e.getCounter(record, mKey)
			if err != nil {
				// TODO: handle this better when we have a more
				//  sophisticated error handler mechanism for this ForEach method.
				forEachError = err
				return
			}

			desc := record.Descriptor()
			c.Add(value.CoerceToFloat64(desc.NumberKind()))
		}

		if gauge, ok := agg.(aggregator.LastValue); ok {
			lv, _, err := gauge.LastValue()
			if err != nil {
				// TODO: handle this better when we have a more
				//  sophisticated error handler mechanism for this ForEach method.
				forEachError = err
				return
			}

			g, err := e.getGauge(record, mKey)
			if err != nil {
				// TODO: handle this better when we have a more
				//  sophisticated error handler mechanism for this ForEach method.
				forEachError = err
				return
			}

			desc := record.Descriptor()
			g.Set(lv.CoerceToFloat64(desc.NumberKind()))
		}
	})

	return forEachError
}

func (e *Exporter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.handler.ServeHTTP(w, r)
}

func (e *Exporter) getCounter(record export.Record, mKey metricKey) (prometheus.Counter, error) {
	e.Lock()
	defer e.Unlock()
	if c, ok := e.counters[mKey]; ok {
		return c, nil
	}

	desc := record.Descriptor()
	counterVec, err := e.getCounterVec(desc, record.Labels())
	if err != nil {
		return nil, err
	}

	counter, err := counterVec.GetMetricWith(labelsToTags(record.Labels()))
	if err != nil {
		return nil, err
	}

	e.counters[mKey] = counter
	return counter, nil
}

func (e *Exporter) getGauge(record export.Record, mKey metricKey) (prometheus.Gauge, error) {
	e.Lock()
	defer e.Unlock()
	if g, ok := e.gauges[mKey]; ok {
		return g, nil
	}

	desc := record.Descriptor()
	gaugeVec, err := e.getGaugeVec(desc, record.Labels())
	if err != nil {
		return nil, err
	}

	gauge, err := gaugeVec.GetMetricWith(labelsToTags(record.Labels()))
	if err != nil {
		return nil, err
	}
	e.gauges[mKey] = gauge
	return gauge, nil
}

func (e *Exporter) getCounterVec(desc *export.Descriptor, labels export.Labels) (*prometheus.CounterVec, error) {
	if c, ok := e.counterVecs[desc]; ok {
		return c, nil
	}

	tagKeys := getTagKeys(labels.Ordered())
	c := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: Sanitize(desc.Name()),
			Help: desc.Description(),
		},
		tagKeys,
	)

	if err := e.registerer.Register(c); err != nil {
		return nil, err
	}

	e.counterVecs[desc] = c
	return c, nil
}

func (e *Exporter) getGaugeVec(desc *export.Descriptor, labels export.Labels) (*prometheus.GaugeVec, error) {
	if g, ok := e.gaugeVecs[desc]; ok {
		return g, nil
	}

	tagKeys := getTagKeys(labels.Ordered())
	g := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: Sanitize(desc.Name()),
			Help: desc.Description(),
		},
		tagKeys,
	)

	if err := e.registerer.Register(g); err != nil {
		return nil, err
	}

	e.gaugeVecs[desc] = g
	return g, nil
}

func getTagKeys(keys []core.KeyValue) []string {
	tagKeys := make([]string, 0, len(keys))
	for _, kv := range keys {
		tagKeys = append(tagKeys, Sanitize(string(kv.Key)))
	}
	return tagKeys
}

func labelsToTags(labels export.Labels) map[string]string {
	tags := make(map[string]string, labels.Len())
	for _, label := range labels.Ordered() {
		tags[Sanitize(string(label.Key))] = label.Value.Emit()
	}
	return tags
}
