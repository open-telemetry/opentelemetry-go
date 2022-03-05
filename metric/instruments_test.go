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

package metric_test

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

func TestInstruments_AsyncFloat64Counter(t *testing.T) {
	ctx := context.Background()
	config := prometheus.Config{}

	c := controller.New(
		processor.NewFactory(
			selector.NewWithInexpensiveDistribution(),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
	)
	exporter, err := prometheus.New(config, c)
	if err != nil {
		log.Panicf("failed to initialize prometheus exporter %v", err)
	}

	meterProvider := exporter.MeterProvider()
	meter := meterProvider.Meter("test-meter")
	i := metric.Instruments{Meter: meter}
	counter, _ := i.AsyncFloat64Counter("test-AsyncFloat64Counter")
	counter.Observe(ctx, 1)
	require.NoError(t, c.Stop(ctx))
}

func TestInstruments_AsyncFloat64Gauge(t *testing.T) {
	ctx := context.Background()
	config := prometheus.Config{}

	c := controller.New(
		processor.NewFactory(
			selector.NewWithInexpensiveDistribution(),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
	)
	exporter, err := prometheus.New(config, c)
	if err != nil {
		log.Panicf("failed to initialize prometheus exporter %v", err)
	}

	meterProvider := exporter.MeterProvider()
	meter := meterProvider.Meter("test-meter")
	i := metric.Instruments{Meter: meter}
	gauge, _ := i.AsyncFloat64Gauge("test-AsyncFloat64Gauge")
	gauge.Observe(ctx, 1)
	require.NoError(t, c.Stop(ctx))
}

func TestInstruments_AsyncFloat64UpDownCounter(t *testing.T) {
	ctx := context.Background()
	config := prometheus.Config{}

	c := controller.New(
		processor.NewFactory(
			selector.NewWithInexpensiveDistribution(),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
	)
	exporter, err := prometheus.New(config, c)
	if err != nil {
		log.Panicf("failed to initialize prometheus exporter %v", err)
	}

	meterProvider := exporter.MeterProvider()
	meter := meterProvider.Meter("test-meter")
	i := metric.Instruments{Meter: meter}
	counter, _ := i.AsyncFloat64UpDownCounter("test-AsyncFloat64UpDownCounter")
	counter.Observe(ctx, 1)
	require.NoError(t, c.Stop(ctx))
}

func TestInstruments_AsyncInt64Counter(t *testing.T) {
	ctx := context.Background()
	config := prometheus.Config{}

	c := controller.New(
		processor.NewFactory(
			selector.NewWithInexpensiveDistribution(),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
	)
	exporter, err := prometheus.New(config, c)
	if err != nil {
		log.Panicf("failed to initialize prometheus exporter %v", err)
	}

	meterProvider := exporter.MeterProvider()
	meter := meterProvider.Meter("test-meter")
	i := metric.Instruments{Meter: meter}
	counter, _ := i.AsyncInt64Counter("test-AsyncInt64Counter")
	counter.Observe(ctx, 1)
	require.NoError(t, c.Stop(ctx))
}

func TestInstruments_AsyncInt64Gauge(t *testing.T) {
	ctx := context.Background()
	config := prometheus.Config{}

	c := controller.New(
		processor.NewFactory(
			selector.NewWithInexpensiveDistribution(),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
	)
	exporter, err := prometheus.New(config, c)
	if err != nil {
		log.Panicf("failed to initialize prometheus exporter %v", err)
	}

	meterProvider := exporter.MeterProvider()
	meter := meterProvider.Meter("test-meter")
	i := metric.Instruments{Meter: meter}
	gauge, _ := i.AsyncInt64Gauge("test-AsyncInt64Gauge")
	gauge.Observe(ctx, 1)
	require.NoError(t, c.Stop(ctx))
}

func TestInstruments_AsyncInt64UpDownCounter(t *testing.T) {
	ctx := context.Background()
	config := prometheus.Config{}

	c := controller.New(
		processor.NewFactory(
			selector.NewWithInexpensiveDistribution(),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
	)
	exporter, err := prometheus.New(config, c)
	if err != nil {
		log.Panicf("failed to initialize prometheus exporter %v", err)
	}

	meterProvider := exporter.MeterProvider()
	meter := meterProvider.Meter("test-meter")
	i := metric.Instruments{Meter: meter}
	counter, _ := i.AsyncInt64UpDownCounter("test-AsyncInt64UpDownCounter")
	counter.Observe(ctx, 1)
	require.NoError(t, c.Stop(ctx))
}

func TestInstruments_SyncFloat64Counter(t *testing.T) {
	ctx := context.Background()
	config := prometheus.Config{}

	c := controller.New(
		processor.NewFactory(
			selector.NewWithInexpensiveDistribution(),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
	)
	exporter, err := prometheus.New(config, c)
	if err != nil {
		log.Panicf("failed to initialize prometheus exporter %v", err)
	}

	meterProvider := exporter.MeterProvider()
	meter := meterProvider.Meter("test-meter")
	i := metric.Instruments{Meter: meter}
	counter, _ := i.SyncFloat64Counter("test-SyncFloat64Counter")
	counter.Add(ctx, 1)
	require.NoError(t, c.Stop(ctx))
}

func TestInstruments_SyncFloat64Histogram(t *testing.T) {
	ctx := context.Background()
	config := prometheus.Config{
		DefaultHistogramBoundaries: []float64{1, 2, 5, 10, 20, 50},
	}

	c := controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries(config.DefaultHistogramBoundaries),
			),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
	)
	exporter, err := prometheus.New(config, c)
	if err != nil {
		log.Panicf("failed to initialize prometheus exporter %v", err)
	}

	meterProvider := exporter.MeterProvider()
	meter := meterProvider.Meter("test-meter")
	i := metric.Instruments{Meter: meter}
	record, _ := i.SyncFloat64Histogram("test-SyncFloat64Histogram")
	record.Record(ctx, 1)
	require.NoError(t, c.Stop(ctx))
}

func TestInstruments_SyncFloat64UpDownCounter(t *testing.T) {
	ctx := context.Background()
	config := prometheus.Config{}

	c := controller.New(
		processor.NewFactory(
			selector.NewWithInexpensiveDistribution(),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
	)
	exporter, err := prometheus.New(config, c)
	if err != nil {
		log.Panicf("failed to initialize prometheus exporter %v", err)
	}

	meterProvider := exporter.MeterProvider()
	meter := meterProvider.Meter("test-meter")
	i := metric.Instruments{Meter: meter}
	counter, _ := i.SyncFloat64UpDownCounter("test-SyncFloat64UpDownCounter")
	counter.Add(ctx, 1)
	require.NoError(t, c.Stop(ctx))
}

func TestInstruments_SyncInt64Counter(t *testing.T) {
	ctx := context.Background()
	config := prometheus.Config{}

	c := controller.New(
		processor.NewFactory(
			selector.NewWithInexpensiveDistribution(),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
	)
	exporter, err := prometheus.New(config, c)
	if err != nil {
		log.Panicf("failed to initialize prometheus exporter %v", err)
	}

	meterProvider := exporter.MeterProvider()
	meter := meterProvider.Meter("test-meter")
	i := metric.Instruments{Meter: meter}
	counter, _ := i.SyncInt64Counter("test-SyncInt64Counter")
	counter.Add(ctx, 1)
	require.NoError(t, c.Stop(ctx))
}

func TestInstruments_SyncInt64Histogram(t *testing.T) {
	ctx := context.Background()
	config := prometheus.Config{
		DefaultHistogramBoundaries: []float64{1, 2, 5, 10, 20, 50},
	}

	c := controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries(config.DefaultHistogramBoundaries),
			),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
	)
	exporter, err := prometheus.New(config, c)
	if err != nil {
		log.Panicf("failed to initialize prometheus exporter %v", err)
	}

	meterProvider := exporter.MeterProvider()
	meter := meterProvider.Meter("test-meter")
	i := metric.Instruments{Meter: meter}
	record, _ := i.SyncInt64Histogram("test-SyncInt64Histogram")
	record.Record(ctx, 1)
	require.NoError(t, c.Stop(ctx))
}

func TestInstruments_SyncInt64UpDownCounter(t *testing.T) {
	ctx := context.Background()
	config := prometheus.Config{}

	c := controller.New(
		processor.NewFactory(
			selector.NewWithInexpensiveDistribution(),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
	)
	exporter, err := prometheus.New(config, c)
	if err != nil {
		log.Panicf("failed to initialize prometheus exporter %v", err)
	}

	meterProvider := exporter.MeterProvider()
	meter := meterProvider.Meter("test-meter")
	i := metric.Instruments{Meter: meter}
	counter, _ := i.SyncInt64UpDownCounter("test-SyncInt64UpDownCounter")
	counter.Add(ctx, 1)
	require.NoError(t, c.Stop(ctx))
}
