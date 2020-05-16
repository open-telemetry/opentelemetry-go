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

	"go.opentelemetry.io/otel/api/metric"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/label"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	integrator "go.opentelemetry.io/otel/sdk/metric/integrator/simple"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
)

// Exporter is an implementation of metric.Exporter that sends metrics to
// Prometheus.
type Exporter struct {
	handler http.Handler

	registerer prometheus.Registerer
	gatherer   prometheus.Gatherer

	snapshot export.CheckpointSet
	onError  func(error)

	defaultSummaryQuantiles    []float64
	defaultHistogramBoundaries []metric.Number
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

	// DefaultHistogramBoundaries defines the default histogram bucket
	// boundaries.
	DefaultHistogramBoundaries []metric.Number

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
		handler:                    promhttp.HandlerFor(config.Gatherer, promhttp.HandlerOpts{}),
		registerer:                 config.Registerer,
		gatherer:                   config.Gatherer,
		defaultSummaryQuantiles:    config.DefaultSummaryQuantiles,
		defaultHistogramBoundaries: config.DefaultHistogramBoundaries,
		onError:                    config.OnError,
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
// chaining a NewRawExporter into the recommended selectors and integrators.
func NewExportPipeline(config Config, period time.Duration) (*push.Controller, http.HandlerFunc, error) {
	selector := simple.NewWithHistogramDistribution(config.DefaultHistogramBoundaries)
	exporter, err := NewRawExporter(config)
	if err != nil {
		return nil, nil, err
	}

	// Prometheus needs to use a stateful integrator since counters (and histogram since they are a collection of Counters)
	// are cumulative (i.e., monotonically increasing values) and should not be resetted after each export.
	//
	// Prometheus uses this approach to be resilient to scrape failures.
	// If a Prometheus server tries to scrape metrics from a host and fails for some reason,
	// it could try again on the next scrape and no data would be lost, only resolution.
	//
	// Gauges (or LastValues) and Summaries are an exception to this and have different behaviors.
	integrator := integrator.New(selector, true)
	pusher := push.New(integrator, exporter, period)
	pusher.Start()

	return pusher, exporter.ServeHTTP, nil
}

// Export exports the provide metric record to prometheus.
func (e *Exporter) Export(_ context.Context, _ *resource.Resource, checkpointSet export.CheckpointSet) error {
	// TODO: Use the resource value in this exporter.
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

	err := c.exp.snapshot.ForEach(func(record export.Record) error {
		agg := record.Aggregator()
		numberKind := record.Descriptor().NumberKind()
		labels := labelValues(record.Labels())
		desc := c.toDesc(&record)

		if hist, ok := agg.(aggregator.Histogram); ok {
			if err := c.exportHistogram(ch, hist, numberKind, desc, labels); err != nil {
				return fmt.Errorf("exporting histogram: %w", err)
			}
		} else if dist, ok := agg.(aggregator.Distribution); ok {
			// TODO: summaries values are never being resetted.
			//  As measurements are recorded, new records starts to have less impact on these summaries.
			//  We should implement an solution that is similar to the Prometheus Clients
			//  using a rolling window for summaries could be a solution.
			//
			//  References:
			// 	https://www.robustperception.io/how-does-a-prometheus-summary-work
			//  https://github.com/prometheus/client_golang/blob/fa4aa9000d2863904891d193dea354d23f3d712a/prometheus/summary.go#L135
			if err := c.exportSummary(ch, dist, numberKind, desc, labels); err != nil {
				return fmt.Errorf("exporting summary: %w", err)
			}
		} else if sum, ok := agg.(aggregator.Sum); ok {
			if err := c.exportCounter(ch, sum, numberKind, desc, labels); err != nil {
				return fmt.Errorf("exporting counter: %w", err)
			}
		} else if lastValue, ok := agg.(aggregator.LastValue); ok {
			if err := c.exportLastValue(ch, lastValue, numberKind, desc, labels); err != nil {
				return fmt.Errorf("exporting last value: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		c.exp.onError(err)
	}
}

func (c *collector) exportLastValue(ch chan<- prometheus.Metric, lvagg aggregator.LastValue, kind metric.NumberKind, desc *prometheus.Desc, labels []string) error {
	lv, _, err := lvagg.LastValue()
	if err != nil {
		return fmt.Errorf("error retrieving last value: %w", err)
	}

	m, err := prometheus.NewConstMetric(desc, prometheus.GaugeValue, lv.CoerceToFloat64(kind), labels...)
	if err != nil {
		return fmt.Errorf("error creating constant metric: %w", err)
	}

	ch <- m
	return nil
}

func (c *collector) exportCounter(ch chan<- prometheus.Metric, sum aggregator.Sum, kind metric.NumberKind, desc *prometheus.Desc, labels []string) error {
	v, err := sum.Sum()
	if err != nil {
		return fmt.Errorf("error retrieving counter: %w", err)
	}

	m, err := prometheus.NewConstMetric(desc, prometheus.CounterValue, v.CoerceToFloat64(kind), labels...)
	if err != nil {
		return fmt.Errorf("error creating constant metric: %w", err)
	}

	ch <- m
	return nil
}

func (c *collector) exportSummary(ch chan<- prometheus.Metric, dist aggregator.Distribution, kind metric.NumberKind, desc *prometheus.Desc, labels []string) error {
	count, err := dist.Count()
	if err != nil {
		return fmt.Errorf("error retrieving count: %w", err)
	}

	var sum metric.Number
	sum, err = dist.Sum()
	if err != nil {
		return fmt.Errorf("error retrieving distribution sum: %w", err)
	}

	quantiles := make(map[float64]float64)
	for _, quantile := range c.exp.defaultSummaryQuantiles {
		q, _ := dist.Quantile(quantile)
		quantiles[quantile] = q.CoerceToFloat64(kind)
	}

	m, err := prometheus.NewConstSummary(desc, uint64(count), sum.CoerceToFloat64(kind), quantiles, labels...)
	if err != nil {
		return fmt.Errorf("error creating constant summary: %w", err)
	}

	ch <- m
	return nil
}

func (c *collector) exportHistogram(ch chan<- prometheus.Metric, hist aggregator.Histogram, kind metric.NumberKind, desc *prometheus.Desc, labels []string) error {
	buckets, err := hist.Histogram()
	if err != nil {
		return fmt.Errorf("error retrieving histogram: %w", err)
	}
	sum, err := hist.Sum()
	if err != nil {
		return fmt.Errorf("error retrieving sum: %w", err)
	}

	var totalCount uint64
	// counts maps from the bucket upper-bound to the cumulative count.
	// The bucket with upper-bound +inf is not included.
	counts := make(map[float64]uint64, len(buckets.Boundaries))
	for i := range buckets.Boundaries {
		boundary := buckets.Boundaries[i].CoerceToFloat64(kind)
		totalCount += buckets.Counts[i].AsUint64()
		counts[boundary] = totalCount
	}
	// Include the +inf bucket in the total count.
	totalCount += buckets.Counts[len(buckets.Counts)-1].AsUint64()

	m, err := prometheus.NewConstHistogram(desc, totalCount, sum.CoerceToFloat64(kind), counts, labels...)
	if err != nil {
		return fmt.Errorf("error creating constant histogram: %w", err)
	}

	ch <- m
	return nil
}

func (c *collector) toDesc(record *export.Record) *prometheus.Desc {
	desc := record.Descriptor()
	labels := labelsKeys(record.Labels())
	return prometheus.NewDesc(sanitize(desc.Name()), desc.Description(), labels, nil)
}

func (e *Exporter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.handler.ServeHTTP(w, r)
}

func labelsKeys(labels *label.Set) []string {
	iter := labels.Iter()
	keys := make([]string, 0, iter.Len())
	for iter.Next() {
		kv := iter.Label()
		keys = append(keys, sanitize(string(kv.Key)))
	}
	return keys
}

func labelValues(labels *label.Set) []string {
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
