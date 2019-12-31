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
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.opentelemetry.io/otel/api/core"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
)

// Exporter is an implementation of metric.Exporter that sends metrics to
// Prometheus.
type Exporter struct {
	handler http.Handler

	registerer prometheus.Registerer
	gatherer   prometheus.Gatherer

	snapshot export.CheckpointSet
	onError  func(error)

	defaultSummaryQuantiles []float64
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

	// DefaultSummaryQuantiles is the default summary quantiles
	// to use. Use nil to specify the system-default summary quantiles.
	DefaultSummaryQuantiles []float64

	// OnError is a function that handle errors that may occur while exporting metrics.
	// TODO: This should be refactored or even removed once we have a better error handling mechanism.
	OnError func(error)
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

	if opts.OnError == nil {
		opts.OnError = func(err error) {
			fmt.Println(err.Error())
		}
	}

	e := &Exporter{
		handler:                 promhttp.HandlerFor(opts.Gatherer, promhttp.HandlerOpts{}),
		registerer:              opts.Registerer,
		gatherer:                opts.Gatherer,
		defaultSummaryQuantiles: opts.DefaultSummaryQuantiles,
	}

	c := newCollector(e)
	if err := opts.Registerer.Register(c); err != nil {
		opts.OnError(fmt.Errorf("cannot register the collector: %w", err))
	}

	return e, nil
}

// Export exports the provide metric record to prometheus.
func (e *Exporter) Export(_ context.Context, checkpointSet export.CheckpointSet) error {
	e.snapshot = checkpointSet
	return nil
}

// collector implements prometheus.Collector interface.
type collector struct {
	exp *Exporter
}

var _ prometheus.Collector = (*collector)(nil)

func newCollector(exporter *Exporter) *collector {
	return &collector{
		exp: exporter,
	}
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	if c.exp.snapshot == nil {
		return
	}

	c.exp.snapshot.ForEach(func(record export.Record) {
		ch <- c.toDesc(&record)
	})
}

// Collect exports the last calculated CheckpointSet.
//
// Collect is invoked whenever prometheus.Gatherer is also invoked.
// For example, when the HTTP endpoint is invoked by Prometheus.
func (c *collector) Collect(ch chan<- prometheus.Metric) {
	if c.exp.snapshot == nil {
		return
	}

	c.exp.snapshot.ForEach(func(record export.Record) {
		agg := record.Aggregator()
		numberKind := record.Descriptor().NumberKind()
		labels := labelValues(record.Labels())
		desc := c.toDesc(&record)

		// TODO: implement histogram export when the histogram aggregation is done.
		//  https://github.com/open-telemetry/opentelemetry-go/issues/317

		if dist, ok := agg.(aggregator.Distribution); ok {
			// TODO: summaries values are never being resetted.
			//  As measures are recorded, new records starts to have less impact on these summaries.
			//  We should implement an solution that is similar to the Prometheus Clients
			//  using a rolling window for summaries could be a solution.
			//
			//  References:
			// 	https://www.robustperception.io/how-does-a-prometheus-summary-work
			//  https://github.com/prometheus/client_golang/blob/fa4aa9000d2863904891d193dea354d23f3d712a/prometheus/summary.go#L135
			c.exportSummary(ch, dist, numberKind, desc, labels)
		} else if sum, ok := agg.(aggregator.Sum); ok {
			c.exportCounter(ch, sum, numberKind, desc, labels)
		} else if gauge, ok := agg.(aggregator.LastValue); ok {
			c.exportGauge(ch, gauge, numberKind, desc, labels)
		}
	})
}

func (c *collector) exportGauge(ch chan<- prometheus.Metric, gauge aggregator.LastValue, kind core.NumberKind, desc *prometheus.Desc, labels []string) {
	lastValue, _, err := gauge.LastValue()
	if err != nil {
		c.exp.onError(err)
		return
	}

	m, err := prometheus.NewConstMetric(desc, prometheus.GaugeValue, lastValue.CoerceToFloat64(kind), labels...)
	if err != nil {
		c.exp.onError(err)
		return
	}

	ch <- m
}

func (c *collector) exportCounter(ch chan<- prometheus.Metric, sum aggregator.Sum, kind core.NumberKind, desc *prometheus.Desc, labels []string) {
	v, err := sum.Sum()
	if err != nil {
		c.exp.onError(err)
		return
	}

	m, err := prometheus.NewConstMetric(desc, prometheus.CounterValue, v.CoerceToFloat64(kind), labels...)
	if err != nil {
		c.exp.onError(err)
		return
	}

	ch <- m
}

func (c *collector) exportSummary(ch chan<- prometheus.Metric, dist aggregator.Distribution, kind core.NumberKind, desc *prometheus.Desc, labels []string) {
	count, err := dist.Count()
	if err != nil {
		c.exp.onError(err)
		return
	}

	var sum core.Number
	sum, err = dist.Sum()
	if err != nil {
		c.exp.onError(err)
		return
	}

	quantiles := make(map[float64]float64)
	for _, quantile := range c.exp.defaultSummaryQuantiles {
		q, _ := dist.Quantile(quantile)
		quantiles[quantile] = q.CoerceToFloat64(kind)
	}

	m, err := prometheus.NewConstSummary(desc, uint64(count), sum.CoerceToFloat64(kind), quantiles, labels...)
	if err != nil {
		c.exp.onError(err)
		return
	}

	ch <- m
}

func (c *collector) toDesc(metric *export.Record) *prometheus.Desc {
	desc := metric.Descriptor()
	labels := labelsKeys(metric.Labels())
	return prometheus.NewDesc(sanitize(desc.Name()), desc.Description(), labels, nil)
}

func (e *Exporter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.handler.ServeHTTP(w, r)
}

func labelsKeys(labels export.Labels) []string {
	keys := make([]string, 0, labels.Len())
	for _, kv := range labels.Ordered() {
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
