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

package prometheus // import "go.opentelemetry.io/otel/exporters/metric/prometheus"

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.opentelemetry.io/otel/global"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/number"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/controller/pull"
	"go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

// Exporter supports Prometheus pulls.  It does not implement the
// sdk/export/metric.Exporter interface--instead it creates a pull
// controller and reads the latest checkpointed data on-scrape.
type Exporter struct {
	handler http.Handler

	registerer prometheus.Registerer
	gatherer   prometheus.Gatherer

	// lock protects access to the controller. The controller
	// exposes its own lock, but using a dedicated lock in this
	// struct allows the exporter to potentially support multiple
	// controllers (e.g., with different resources).
	lock       sync.RWMutex
	controller *pull.Controller

	defaultSummaryQuantiles    []float64
	defaultHistogramBoundaries []float64
}

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
	DefaultHistogramBoundaries []float64
}

// NewExportPipeline sets up a complete export pipeline with the recommended setup,
// using the recommended selector and standard processor.  See the pull.Options.
func NewExportPipeline(config Config, options ...pull.Option) (*Exporter, error) {
	if config.Registry == nil {
		config.Registry = prometheus.NewRegistry()
	}

	if config.Registerer == nil {
		config.Registerer = config.Registry
	}

	if config.Gatherer == nil {
		config.Gatherer = config.Registry
	}

	e := &Exporter{
		handler:                    promhttp.HandlerFor(config.Gatherer, promhttp.HandlerOpts{}),
		registerer:                 config.Registerer,
		gatherer:                   config.Gatherer,
		defaultSummaryQuantiles:    config.DefaultSummaryQuantiles,
		defaultHistogramBoundaries: config.DefaultHistogramBoundaries,
	}

	c := &collector{
		exp: e,
	}
	e.SetController(config, options...)
	if err := config.Registerer.Register(c); err != nil {
		return nil, fmt.Errorf("cannot register the collector: %w", err)
	}

	return e, nil
}

// InstallNewPipeline instantiates a NewExportPipeline and registers it globally.
// Typically called as:
//
// 	hf, err := prometheus.InstallNewPipeline(prometheus.Config{...})
//
// 	if err != nil {
// 		...
// 	}
// 	http.HandleFunc("/metrics", hf)
// 	defer pipeline.Stop()
// 	... Done
func InstallNewPipeline(config Config, options ...pull.Option) (*Exporter, error) {
	exp, err := NewExportPipeline(config, options...)
	if err != nil {
		return nil, err
	}
	global.SetMeterProvider(exp.MeterProvider())
	return exp, nil
}

// SetController sets up a standard *pull.Controller as the metric provider
// for this exporter.
func (e *Exporter) SetController(config Config, options ...pull.Option) {
	e.lock.Lock()
	defer e.lock.Unlock()

	e.controller = pull.New(
		basic.New(
			simple.NewWithHistogramDistribution(config.DefaultHistogramBoundaries),
			e,
			basic.WithMemory(true),
		),
		options...,
	)
}

// MeterProvider returns the MeterProvider of this exporter.
func (e *Exporter) MeterProvider() metric.MeterProvider {
	return e.controller.MeterProvider()
}

// Controller returns the controller object that coordinates collection for the SDK.
func (e *Exporter) Controller() *pull.Controller {
	e.lock.RLock()
	defer e.lock.RUnlock()
	return e.controller
}

func (e *Exporter) ExportKindFor(desc *metric.Descriptor, kind aggregation.Kind) export.ExportKind {
	// NOTE: Summary values should use Delta aggregation, then be
	// combined into a sliding window, see the TODO below.
	// NOTE: Prometheus also supports a "GaugeDelta" exposition format,
	// which is expressed as a delta histogram.  Need to understand if this
	// should be a default behavior for ValueRecorder/ValueObserver.
	return export.CumulativeExportKindSelector().ExportKindFor(desc, kind)
}

func (e *Exporter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.handler.ServeHTTP(w, r)
}

// collector implements prometheus.Collector interface.
type collector struct {
	exp *Exporter
}

var _ prometheus.Collector = (*collector)(nil)

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	c.exp.lock.RLock()
	defer c.exp.lock.RUnlock()

	_ = c.exp.Controller().ForEach(c.exp, func(record export.Record) error {
		var labelKeys []string
		mergeLabels(record, &labelKeys, nil)
		ch <- c.toDesc(record, labelKeys)
		return nil
	})
}

// Collect exports the last calculated CheckpointSet.
//
// Collect is invoked whenever prometheus.Gatherer is also invoked.
// For example, when the HTTP endpoint is invoked by Prometheus.
func (c *collector) Collect(ch chan<- prometheus.Metric) {
	c.exp.lock.RLock()
	defer c.exp.lock.RUnlock()

	ctrl := c.exp.Controller()
	if err := ctrl.Collect(context.Background()); err != nil {
		global.Handle(err)
	}

	err := ctrl.ForEach(c.exp, func(record export.Record) error {
		agg := record.Aggregation()
		numberKind := record.Descriptor().NumberKind()
		instrumentKind := record.Descriptor().InstrumentKind()

		var labelKeys, labels []string
		mergeLabels(record, &labelKeys, &labels)

		desc := c.toDesc(record, labelKeys)

		if hist, ok := agg.(aggregation.Histogram); ok {
			if err := c.exportHistogram(ch, hist, numberKind, desc, labels); err != nil {
				return fmt.Errorf("exporting histogram: %w", err)
			}
		} else if dist, ok := agg.(aggregation.Distribution); ok {
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
		} else if sum, ok := agg.(aggregation.Sum); ok && instrumentKind.Monotonic() {
			if err := c.exportMonotonicCounter(ch, sum, numberKind, desc, labels); err != nil {
				return fmt.Errorf("exporting monotonic counter: %w", err)
			}
		} else if sum, ok := agg.(aggregation.Sum); ok && !instrumentKind.Monotonic() {
			if err := c.exportNonMonotonicCounter(ch, sum, numberKind, desc, labels); err != nil {
				return fmt.Errorf("exporting non monotonic counter: %w", err)
			}
		} else if lastValue, ok := agg.(aggregation.LastValue); ok {
			if err := c.exportLastValue(ch, lastValue, numberKind, desc, labels); err != nil {
				return fmt.Errorf("exporting last value: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		global.Handle(err)
	}
}

func (c *collector) exportLastValue(ch chan<- prometheus.Metric, lvagg aggregation.LastValue, kind number.Kind, desc *prometheus.Desc, labels []string) error {
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

func (c *collector) exportNonMonotonicCounter(ch chan<- prometheus.Metric, sum aggregation.Sum, kind number.Kind, desc *prometheus.Desc, labels []string) error {
	v, err := sum.Sum()
	if err != nil {
		return fmt.Errorf("error retrieving counter: %w", err)
	}

	m, err := prometheus.NewConstMetric(desc, prometheus.GaugeValue, v.CoerceToFloat64(kind), labels...)
	if err != nil {
		return fmt.Errorf("error creating constant metric: %w", err)
	}

	ch <- m
	return nil
}

func (c *collector) exportMonotonicCounter(ch chan<- prometheus.Metric, sum aggregation.Sum, kind number.Kind, desc *prometheus.Desc, labels []string) error {
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

func (c *collector) exportSummary(ch chan<- prometheus.Metric, dist aggregation.Distribution, kind number.Kind, desc *prometheus.Desc, labels []string) error {
	count, err := dist.Count()
	if err != nil {
		return fmt.Errorf("error retrieving count: %w", err)
	}

	var sum number.Number
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

func (c *collector) exportHistogram(ch chan<- prometheus.Metric, hist aggregation.Histogram, kind number.Kind, desc *prometheus.Desc, labels []string) error {
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
		boundary := buckets.Boundaries[i]
		totalCount += uint64(buckets.Counts[i])
		counts[boundary] = totalCount
	}
	// Include the +inf bucket in the total count.
	totalCount += uint64(buckets.Counts[len(buckets.Counts)-1])

	m, err := prometheus.NewConstHistogram(desc, totalCount, sum.CoerceToFloat64(kind), counts, labels...)
	if err != nil {
		return fmt.Errorf("error creating constant histogram: %w", err)
	}

	ch <- m
	return nil
}

func (c *collector) toDesc(record export.Record, labelKeys []string) *prometheus.Desc {
	desc := record.Descriptor()
	return prometheus.NewDesc(sanitize(desc.Name()), desc.Description(), labelKeys, nil)
}

// mergeLabels merges the export.Record's labels and resources into a
// single set, giving precedence to the record's labels in case of
// duplicate keys.  This outputs one or both of the keys and the
// values as a slice, and either argument may be nil to avoid
// allocating an unnecessary slice.
func mergeLabels(record export.Record, keys, values *[]string) {
	if keys != nil {
		*keys = make([]string, 0, record.Labels().Len()+record.Resource().Len())
	}
	if values != nil {
		*values = make([]string, 0, record.Labels().Len()+record.Resource().Len())
	}

	// Duplicate keys are resolved by taking the record label value over
	// the resource value.
	mi := label.NewMergeIterator(record.Labels(), record.Resource().LabelSet())
	for mi.Next() {
		label := mi.Label()
		if keys != nil {
			*keys = append(*keys, sanitize(string(label.Key)))
		}
		if values != nil {
			*values = append(*values, label.Value.Emit())
		}
	}
}
