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
	"bytes"
	"context"
	"sort"
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"go.opentelemetry.io/otel/api/core"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/counter"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/gauge"
)

const (
	prefixSplitter  = '+'
	keyPairSplitter = ','
)

type metricID *export.Descriptor

// Exporter is an implementation of metric.Exporter that sends metrics to
// Prometheus.
type Exporter struct {
	sync.RWMutex
	registerer prometheus.Registerer
	gatherer   prometheus.Gatherer

	counters map[metricID]prometheus.Counter
	gauges   map[metricID]prometheus.Gauge

	onRegisterError func(e error)
	counterVecs     map[string]*prometheus.CounterVec
	gaugeVecs       map[string]*prometheus.GaugeVec
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

		counters: make(map[metricID]prometheus.Counter),
		gauges:   make(map[metricID]prometheus.Gauge),

		counterVecs: make(map[string]*prometheus.CounterVec),
		gaugeVecs:   make(map[string]*prometheus.GaugeVec),
	}
}

// AggregatorFor returns the metric aggregator used for the particular exporter.
func (e *Exporter) AggregatorFor(record export.Record) export.Aggregator {
	switch record.Descriptor().MetricKind() {
	case export.CounterKind:
		return counter.New()
	case export.GaugeKind:
		return gauge.New()
	}
	return nil
}

// Export exports the provide metric record to prometheus.
func (e *Exporter) Export(
	ctx context.Context,
	record export.Record,
	aggregator export.Aggregator,
) {
	switch record.Descriptor().MetricKind() {
	case export.CounterKind:
		e.exportCounter(record, aggregator)
	case export.GaugeKind:
		e.exportGauge(record, aggregator)
	}
}

func (e *Exporter) exportCounter(record export.Record, aggregator export.Aggregator) {
	c, err := e.getCounter(record)
	if err != nil {
		// TODO: log a warning here?
		return
	}
	// Retrieve the counter value from the aggregator and add it.
	if agg, ok := aggregator.(*counter.Aggregator); ok {
		c.Add(float64(agg.AsNumber()))
	}
}

func (e *Exporter) exportGauge(record export.Record, aggregator export.Aggregator) {
	g, err := e.getGauge(record)
	if err != nil {
		// TODO: log a warning here?
		return
	}
	// Retrieve the gauge value from the aggregator and set it.
	if agg, ok := aggregator.(*gauge.Aggregator); ok {
		g.Set(float64(agg.AsNumber()))
	}
}

func (e *Exporter) getCounter(record export.Record) (prometheus.Counter, error) {
	e.Lock()
	defer e.Unlock()

	desc := record.Descriptor()
	if c, ok := e.counters[desc]; ok {
		return c, nil
	}

	counterVec, err := e.getCounterVec(desc)
	if err != nil {
		return nil, err
	}

	counter := counterVec.With(labelsToTags(record.Labels()))

	e.counters[desc] = counter
	return counter, nil
}

func (e *Exporter) getGauge(record export.Record) (prometheus.Gauge, error) {
	e.Lock()
	defer e.Unlock()

	desc := record.Descriptor()
	if g, ok := e.gauges[desc]; ok {
		return g, nil
	}

	gaugeVec, err := e.getGaugeVec(desc)
	if err != nil {
		return nil, err
	}

	gauge := gaugeVec.With(labelsToTags(record.Labels()))

	e.gauges[desc] = gauge
	return gauge, nil
}

func (e *Exporter) getCounterVec(desc *export.Descriptor) (*prometheus.CounterVec, error) {
	id, tagKeys := getCanonicalID(desc)

	e.Lock()
	defer e.Unlock()

	if c, ok := e.counterVecs[id]; ok {
		return c, nil
	}

	c := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: desc.Name(),
			Help: desc.Description(),
		},
		tagKeys,
	)

	if err := e.registerer.Register(c); err != nil {
		return nil, err
	}

	e.counterVecs[id] = c
	return c, nil
}

func (e *Exporter) getGaugeVec(desc *export.Descriptor) (*prometheus.GaugeVec, error) {
	id, tagKeys := getCanonicalID(desc)

	e.Lock()
	defer e.Unlock()

	if g, ok := e.gaugeVecs[id]; ok {
		return g, nil
	}

	g := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: desc.Name(),
			Help: desc.Description(),
		},
		tagKeys,
	)

	if err := e.registerer.Register(g); err != nil {
		return nil, err
	}

	e.gaugeVecs[id] = g
	return g, nil
}

func getTagKeys(keys []core.Key) []string {
	tagKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		tagKeys = append(tagKeys, string(key))
	}
	return tagKeys
}

func getCanonicalID(desc *export.Descriptor) (string, []string) {
	tagKeys := getTagKeys(desc.Keys())
	sort.Strings(tagKeys)
	return generateKey(desc.Name(), tagKeys), tagKeys
}

func labelsToTags(labels []core.KeyValue) map[string]string {
	tags := make(map[string]string, len(labels))
	for _, label := range labels {
		tags[string(label.Key)] = label.Value.AsString()
	}
	return tags
}

func generateKey(name string, keys []string) string {
	// TODO: pool these objects.
	var buf bytes.Buffer
	buf.WriteString(name)
	buf.WriteByte(prefixSplitter)

	sortedKeysLen := len(keys)
	for i := 0; i < sortedKeysLen; i++ {
		buf.WriteString(keys[i])
		if i != sortedKeysLen-1 {
			buf.WriteByte(keyPairSplitter)
		}
	}
	return buf.String()
}
