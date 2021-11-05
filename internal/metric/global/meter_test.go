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

package global_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/internal/metric/global"
	"go.opentelemetry.io/otel/metric"
	metricglobal "go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/metrictest"
	"go.opentelemetry.io/otel/metric/number"
	"go.opentelemetry.io/otel/metric/sdkapi"
)

var Must = metric.Must

var asInt = number.NewInt64Number
var asFloat = number.NewFloat64Number

func TestDirect(t *testing.T) {
	global.ResetForTest()

	ctx := context.Background()
	meter1 := metricglobal.Meter("test1", metric.WithInstrumentationVersion("semver:v1.0.0"))
	meter2 := metricglobal.Meter("test2", metric.WithSchemaURL("hello"))

	library1 := metrictest.Library{
		InstrumentationName:    "test1",
		InstrumentationVersion: "semver:v1.0.0",
	}
	library2 := metrictest.Library{
		InstrumentationName: "test2",
		SchemaURL:           "hello",
	}

	labels1 := []attribute.KeyValue{attribute.String("A", "B")}
	labels2 := []attribute.KeyValue{attribute.String("C", "D")}
	labels3 := []attribute.KeyValue{attribute.String("E", "F")}

	counter := Must(meter1).NewInt64Counter("test.counter")
	counter.Add(ctx, 1, labels1...)
	counter.Add(ctx, 1, labels1...)

	histogram := Must(meter1).NewFloat64Histogram("test.histogram")
	histogram.Record(ctx, 1, labels1...)
	histogram.Record(ctx, 2, labels1...)

	_ = Must(meter1).NewFloat64GaugeObserver("test.gauge.float", func(_ context.Context, result metric.Float64ObserverResult) {
		result.Observe(1., labels1...)
		result.Observe(2., labels2...)
	})

	_ = Must(meter1).NewInt64GaugeObserver("test.gauge.int", func(_ context.Context, result metric.Int64ObserverResult) {
		result.Observe(1, labels1...)
		result.Observe(2, labels2...)
	})

	second := Must(meter2).NewFloat64Histogram("test.second")
	second.Record(ctx, 1, labels3...)
	second.Record(ctx, 2, labels3...)

	provider := metrictest.NewMeterProvider()
	metricglobal.SetMeterProvider(provider)

	counter.Add(ctx, 1, labels1...)
	histogram.Record(ctx, 3, labels1...)
	second.Record(ctx, 3, labels3...)

	provider.RunAsyncInstruments()

	measurements := metrictest.AsStructs(provider.MeasurementBatches)

	require.EqualValues(t,
		[]metrictest.Measured{
			{
				Name:    "test.counter",
				Library: library1,
				Labels:  metrictest.LabelsToMap(labels1...),
				Number:  asInt(1),
			},
			{
				Name:    "test.histogram",
				Library: library1,
				Labels:  metrictest.LabelsToMap(labels1...),
				Number:  asFloat(3),
			},
			{
				Name:    "test.second",
				Library: library2,
				Labels:  metrictest.LabelsToMap(labels3...),
				Number:  asFloat(3),
			},
			{
				Name:    "test.gauge.float",
				Library: library1,
				Labels:  metrictest.LabelsToMap(labels1...),
				Number:  asFloat(1),
			},
			{
				Name:    "test.gauge.float",
				Library: library1,
				Labels:  metrictest.LabelsToMap(labels2...),
				Number:  asFloat(2),
			},
			{
				Name:    "test.gauge.int",
				Library: library1,
				Labels:  metrictest.LabelsToMap(labels1...),
				Number:  asInt(1),
			},
			{
				Name:    "test.gauge.int",
				Library: library1,
				Labels:  metrictest.LabelsToMap(labels2...),
				Number:  asInt(2),
			},
		},
		measurements,
	)
}

type meterProviderWithConstructorError struct {
	metric.MeterProvider
}

type meterWithConstructorError struct {
	sdkapi.MeterImpl
}

func (m *meterProviderWithConstructorError) Meter(iName string, opts ...metric.MeterOption) metric.Meter {
	return metric.WrapMeterImpl(&meterWithConstructorError{m.MeterProvider.Meter(iName, opts...).MeterImpl()})
}

func (m *meterWithConstructorError) NewSyncInstrument(_ sdkapi.Descriptor) (sdkapi.SyncImpl, error) {
	return sdkapi.NewNoopSyncInstrument(), errors.New("constructor error")
}

func TestErrorInDeferredConstructor(t *testing.T) {
	global.ResetForTest()

	ctx := context.Background()
	meter := metricglobal.GetMeterProvider().Meter("builtin")

	c1 := Must(meter).NewInt64Counter("test")
	c2 := Must(meter).NewInt64Counter("test")

	provider := metrictest.NewMeterProvider()
	sdk := &meterProviderWithConstructorError{provider}

	require.Panics(t, func() {
		metricglobal.SetMeterProvider(sdk)
	})

	c1.Add(ctx, 1)
	c2.Add(ctx, 2)
}

func TestImplementationIndirection(t *testing.T) {
	global.ResetForTest()

	// Test that Implementation() does the proper indirection, i.e.,
	// returns the implementation interface not the global, after
	// registered.

	meter1 := metricglobal.Meter("test1")

	// Sync: no SDK yet
	counter := Must(meter1).NewInt64Counter("interface.counter")

	ival := counter.Measurement(1).SyncImpl().Implementation()
	require.NotNil(t, ival)

	_, ok := ival.(*metrictest.Sync)
	require.False(t, ok)

	// Async: no SDK yet
	gauge := Must(meter1).NewFloat64GaugeObserver(
		"interface.gauge",
		func(_ context.Context, result metric.Float64ObserverResult) {},
	)

	ival = gauge.AsyncImpl().Implementation()
	require.NotNil(t, ival)

	_, ok = ival.(*metrictest.Async)
	require.False(t, ok)

	// Register the SDK
	provider := metrictest.NewMeterProvider()
	metricglobal.SetMeterProvider(provider)

	// Repeat the above tests

	// Sync
	ival = counter.Measurement(1).SyncImpl().Implementation()
	require.NotNil(t, ival)

	_, ok = ival.(*metrictest.Sync)
	require.True(t, ok)

	// Async
	ival = gauge.AsyncImpl().Implementation()
	require.NotNil(t, ival)

	_, ok = ival.(*metrictest.Async)
	require.True(t, ok)
}

func TestRecordBatchMock(t *testing.T) {
	global.ResetForTest()

	meter := metricglobal.GetMeterProvider().Meter("builtin")

	counter := Must(meter).NewInt64Counter("test.counter")

	meter.RecordBatch(context.Background(), nil, counter.Measurement(1))

	provider := metrictest.NewMeterProvider()
	metricglobal.SetMeterProvider(provider)

	meter.RecordBatch(context.Background(), nil, counter.Measurement(1))

	require.EqualValues(t,
		[]metrictest.Measured{
			{
				Name: "test.counter",
				Library: metrictest.Library{
					InstrumentationName: "builtin",
				},
				Labels: metrictest.LabelsToMap(),
				Number: asInt(1),
			},
		},
		metrictest.AsStructs(provider.MeasurementBatches))
}
