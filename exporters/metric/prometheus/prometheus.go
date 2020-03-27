// Copyright The OpenTelemetry Authors
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
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/global"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/batcher/defaultkeys"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
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

// Config is a set of configs for the tally reporter.
type Config struct {
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

// NewRawExporter returns a new prometheus exporter for prometheus metrics
// for use in a pipeline.
func NewRawExporter(config Config) (*Exporter, error) {
	if config.Registry == nil {
		config.Registry = prometheus.NewRegistry()
	}

	if config.Registerer == nil {
		config.Registerer = config.Registry
	}

	if config.Gatherer == nil {
		config.Gatherer = config.Registry
	}

	if config.OnError == nil {
		config.OnError = func(err error) {
			fmt.Println(err.Error())
		}
	}

	e := &Exporter{
		handler:                 promhttp.HandlerFor(config.Gatherer, promhttp.HandlerOpts{}),
		registerer:              config.Registerer,
		gatherer:                config.Gatherer,
		defaultSummaryQuantiles: config.DefaultSummaryQuantiles,
		onError:                 config.OnError,
	}

	c := newCollector(e)
	if err := config.Registerer.Register(c); err != nil {
		config.OnError(fmt.Errorf("cannot register the collector: %w", err))
	}

	return e, nil
}

// InstallNewPipeline instantiates a NewExportPipeline and registers it globally.
// Typically called as:
//
// 	pipeline, hf, err := prometheus.InstallNewPipeline(prometheus.Config{...})
//
// 	if err != nil {
// 		...
// 	}
// 	http.HandleFunc("/metrics", hf)
// 	defer pipeline.Stop()
// 	... Done
func InstallNewPipeline(config Config) (*push.Controller, http.HandlerFunc, error) {
	controller, hf, err := NewExportPipeline(config, time.Minute)
	if err != nil {
		return controller, hf, err
	}
	global.SetMeterProvider(controller)
	return controller, hf, err
}

// NewExportPipeline sets up a complete export pipeline with the recommended setup,
// chaining a NewRawExporter into the recommended selectors and batchers.
func NewExportPipeline(config Config, period time.Duration) (*push.Controller, http.HandlerFunc, error) {
	selector := simple.NewWithExactMeasure()
	exporter, err := NewRawExporter(config)
	if err != nil {
		return nil, nil, err
	}

	// Prometheus needs to use a stateful batcher since counters (and histogram since they are a collection of Counters)
	// are cumulative (i.e., monotonically increasing values) and should not be resetted after each export.
	//
	// Prometheus uses this approach to be resilient to scrape failures.
	// If a Prometheus server tries to scrape metrics from a host and fails for some reason,
	// it could try again on the next scrape and no data would be lost, only resolution.
	//
	// Gauges (or LastValues) and Summaries are an exception to this and have different behaviors.
	batcher := defaultkeys.New(selector, export.NewDefaultLabelEncoder(), true)
	pusher := push.New(batcher, exporter, period)
	pusher.Start()

	return pusher, exporter.ServeHTTP, nil
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

	_ = c.exp.snapshot.ForEach(func(record export.Record) error {
		ch <- c.toDesc(&record)
		return nil
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

	_ = c.exp.snapshot.ForEach(func(record export.Record) error {
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
		} else if lastValue, ok := agg.(aggregator.LastValue); ok {
			c.exportLastValue(ch, lastValue, numberKind, desc, labels)
		}
		return nil
	})
}

func (c *collector) exportLastValue(ch chan<- prometheus.Metric, lvagg aggregator.LastValue, kind core.NumberKind, desc *prometheus.Desc, labels []string) {
	lv, _, err := lvagg.LastValue()
	if err != nil {
		c.exp.onError(err)
		return
	}

	m, err := prometheus.NewConstMetric(desc, prometheus.GaugeValue, lv.CoerceToFloat64(kind), labels...)
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

func (c *collector) toDesc(record *export.Record) *prometheus.Desc {
	desc := record.Descriptor()
	labels := labelsKeys(record.Labels())
	return prometheus.NewDesc(sanitize(desc.Name()), desc.Description(), labels, nil)
}

func (e *Exporter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.handler.ServeHTTP(w, r)
}

func labelsKeys(labels export.Labels) []string {
	iter := labels.Iter()
	keys := make([]string, 0, iter.Len())
	for iter.Next() {
		kv := iter.Label()
		keys = append(keys, sanitize(string(kv.Key)))
	}
	return keys
}

func labelValues(labels export.Labels) []string {
	// TODO(paivagustavo): parse the labels.Encoded() instead of calling `Emit()` directly
	//  this would avoid unnecessary allocations.
	iter := labels.Iter()
	values := make([]string, 0, iter.Len())
	for iter.Next() {
		label := iter.Label()
		values = append(values, label.Value.Emit())
	}
	return values
}
