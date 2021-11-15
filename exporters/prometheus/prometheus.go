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

package prometheus // import "go.opentelemetry.io/otel/exporters/prometheus"

// Note that this package does not support a way to export Prometheus
// Summary data points, removed in PR#1412.

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/number"
	"go.opentelemetry.io/otel/metric/sdkapi"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/resource"
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
	controller *controller.Controller

	defaultHistogramBoundaries []float64
}

// ErrUnsupportedAggregator is returned for unrepresentable aggregator
// types.
var ErrUnsupportedAggregator = fmt.Errorf("unsupported aggregator type")

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

	// DefaultHistogramBoundaries defines the default histogram bucket
	// boundaries.
	DefaultHistogramBoundaries []float64
}

// New returns a new Prometheus exporter using the configured metric
// controller.  See controller.New().
func New(config Config, controller *controller.Controller) (*Exporter, error) {
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
		controller:                 controller,
		defaultHistogramBoundaries: config.DefaultHistogramBoundaries,
	}

	c := &collector{
		exp: e,
	}
	if err := config.Registerer.Register(c); err != nil {
		return nil, fmt.Errorf("cannot register the collector: %w", err)
	}
	return e, nil
}

// MeterProvider returns the MeterProvider of this exporter.
func (e *Exporter) MeterProvider() metric.MeterProvider {
	return e.controller
}

// Controller returns the controller object that coordinates collection for the SDK.
func (e *Exporter) Controller() *controller.Controller {
	e.lock.RLock()
	defer e.lock.RUnlock()
	return e.controller
}

// TemporalityFor implements TemporalitySelector.
func (e *Exporter) TemporalityFor(desc *sdkapi.Descriptor, kind aggregation.Kind) aggregation.Temporality {
	return aggregation.CumulativeTemporalitySelector().TemporalityFor(desc, kind)
}

// ServeHTTP implements http.Handler.
func (e *Exporter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.handler.ServeHTTP(w, r)
}

// collector implements prometheus.Collector interface.
type collector struct {
	exp *Exporter
}

var _ prometheus.Collector = (*collector)(nil)

// Describe implements prometheus.Collector.
func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	c.exp.lock.RLock()
	defer c.exp.lock.RUnlock()

	_ = c.exp.Controller().ForEach(func(_ instrumentation.Library, reader export.Reader) error {
		return reader.ForEach(c.exp, func(record export.Record) error {
			var labelKeys []string
			mergeLabels(record, c.exp.controller.Resource(), &labelKeys, nil)
			ch <- c.toDesc(record, labelKeys)
			return nil
		})
	})
}

// Collect exports the last calculated Reader state.
//
// Collect is invoked whenever prometheus.Gatherer is also invoked.
// For example, when the HTTP endpoint is invoked by Prometheus.
func (c *collector) Collect(ch chan<- prometheus.Metric) {
	c.exp.lock.RLock()
	defer c.exp.lock.RUnlock()

	ctrl := c.exp.Controller()
	if err := ctrl.Collect(context.Background()); err != nil {
		otel.Handle(err)
	}

	err := ctrl.ForEach(func(_ instrumentation.Library, reader export.Reader) error {
		return reader.ForEach(c.exp, func(record export.Record) error {

			agg := record.Aggregation()
			numberKind := record.Descriptor().NumberKind()
			instrumentKind := record.Descriptor().InstrumentKind()

			var labelKeys, labels []string
			mergeLabels(record, c.exp.controller.Resource(), &labelKeys, &labels)

			desc := c.toDesc(record, labelKeys)

			if hist, ok := agg.(aggregation.Histogram); ok {
				if err := c.exportHistogram(ch, hist, numberKind, desc, labels); err != nil {
					return fmt.Errorf("exporting histogram: %w", err)
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
			} else {
				return fmt.Errorf("%w: %s", ErrUnsupportedAggregator, agg.Kind())
			}
			return nil
		})
	})
	if err != nil {
		otel.Handle(err)
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
func mergeLabels(record export.Record, res *resource.Resource, keys, values *[]string) {
	if keys != nil {
		*keys = make([]string, 0, record.Labels().Len()+res.Len())
	}
	if values != nil {
		*values = make([]string, 0, record.Labels().Len()+res.Len())
	}

	// Duplicate keys are resolved by taking the record label value over
	// the resource value.
	mi := attribute.NewMergeIterator(record.Labels(), res.Set())
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
