// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package global

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
	"go.opentelemetry.io/otel/metric/noop"
)

func testFloat64ConcurrentSafe(interact func(float64), setDelegate func(metric.Meter)) {
	done := make(chan struct{})
	finish := make(chan struct{})
	go func() {
		defer close(done)
		for {
			interact(1)
			select {
			case <-finish:
				return
			default:
			}
		}
	}()

	setDelegate(noop.NewMeterProvider().Meter(""))
	close(finish)
	<-done
}

func testInt64ConcurrentSafe(interact func(int64), setDelegate func(metric.Meter)) {
	done := make(chan struct{})
	finish := make(chan struct{})
	go func() {
		defer close(done)
		for {
			interact(1)
			select {
			case <-finish:
				return
			default:
			}
		}
	}()

	setDelegate(noop.NewMeterProvider().Meter(""))
	close(finish)
	<-done
}

func TestAsyncInstrumentSetDelegateConcurrentSafe(t *testing.T) {
	// Float64 Instruments
	t.Run("Float64", func(t *testing.T) {
		t.Run("Counter", func(*testing.T) {
			delegate := &afCounter{}
			f := func(float64) { _ = delegate.unwrap() }
			testFloat64ConcurrentSafe(f, delegate.setDelegate)
		})

		t.Run("UpDownCounter", func(*testing.T) {
			delegate := &afUpDownCounter{}
			f := func(float64) { _ = delegate.unwrap() }
			testFloat64ConcurrentSafe(f, delegate.setDelegate)
		})

		t.Run("Gauge", func(*testing.T) {
			delegate := &afGauge{}
			f := func(float64) { _ = delegate.unwrap() }
			testFloat64ConcurrentSafe(f, delegate.setDelegate)
		})
	})

	// Int64 Instruments

	t.Run("Int64", func(t *testing.T) {
		t.Run("Counter", func(*testing.T) {
			delegate := &aiCounter{}
			f := func(int64) { _ = delegate.unwrap() }
			testInt64ConcurrentSafe(f, delegate.setDelegate)
		})

		t.Run("UpDownCounter", func(*testing.T) {
			delegate := &aiUpDownCounter{}
			f := func(int64) { _ = delegate.unwrap() }
			testInt64ConcurrentSafe(f, delegate.setDelegate)
		})

		t.Run("Gauge", func(*testing.T) {
			delegate := &aiGauge{}
			f := func(int64) { _ = delegate.unwrap() }
			testInt64ConcurrentSafe(f, delegate.setDelegate)
		})
	})
}

func TestSyncInstrumentSetDelegateConcurrentSafe(t *testing.T) {
	// Float64 Instruments
	t.Run("Float64", func(*testing.T) {
		t.Run("Counter", func(*testing.T) {
			delegate := &sfCounter{}
			f := func(v float64) { delegate.Add(t.Context(), v) }
			testFloat64ConcurrentSafe(f, delegate.setDelegate)
		})

		t.Run("UpDownCounter", func(*testing.T) {
			delegate := &sfUpDownCounter{}
			f := func(v float64) { delegate.Add(t.Context(), v) }
			testFloat64ConcurrentSafe(f, delegate.setDelegate)
		})

		t.Run("Histogram", func(*testing.T) {
			delegate := &sfHistogram{}
			f := func(v float64) { delegate.Record(t.Context(), v) }
			testFloat64ConcurrentSafe(f, delegate.setDelegate)
		})

		t.Run("Gauge", func(*testing.T) {
			delegate := &sfGauge{}
			f := func(v float64) { delegate.Record(t.Context(), v) }
			testFloat64ConcurrentSafe(f, delegate.setDelegate)
		})
	})

	// Int64 Instruments

	t.Run("Int64", func(*testing.T) {
		t.Run("Counter", func(*testing.T) {
			delegate := &siCounter{}
			f := func(v int64) { delegate.Add(t.Context(), v) }
			testInt64ConcurrentSafe(f, delegate.setDelegate)
		})

		t.Run("UpDownCounter", func(*testing.T) {
			delegate := &siUpDownCounter{}
			f := func(v int64) { delegate.Add(t.Context(), v) }
			testInt64ConcurrentSafe(f, delegate.setDelegate)
		})

		t.Run("Histogram", func(*testing.T) {
			delegate := &siHistogram{}
			f := func(v int64) { delegate.Record(t.Context(), v) }
			testInt64ConcurrentSafe(f, delegate.setDelegate)
		})

		t.Run("Gauge", func(*testing.T) {
			delegate := &siGauge{}
			f := func(v int64) { delegate.Record(t.Context(), v) }
			testInt64ConcurrentSafe(f, delegate.setDelegate)
		})
	})
}

type testFloat64Counter struct {
	count int

	embedded.Float64Counter
}

func (i *testFloat64Counter) observe() {
	i.count++
}

func (i *testFloat64Counter) Add(context.Context, float64, ...metric.AddOption) {
	i.count++
}

func (i *testFloat64Counter) Record(context.Context, float64, ...metric.RecordOption) {
	i.count++
}

func (i *testFloat64Counter) Enabled(context.Context) bool {
	return true
}

func (i *testFloat64Counter) WithAttributes(...attribute.KeyValue) metric.Float64Counter {
	return i
}

type testFloat64UpDownCounter struct {
	count int

	embedded.Float64UpDownCounter
}

func (i *testFloat64UpDownCounter) observe() {
	i.count++
}

func (i *testFloat64UpDownCounter) Add(context.Context, float64, ...metric.AddOption) {
	i.count++
}

func (i *testFloat64UpDownCounter) Record(context.Context, float64, ...metric.RecordOption) {
	i.count++
}

func (i *testFloat64UpDownCounter) Enabled(context.Context) bool {
	return true
}

func (i *testFloat64UpDownCounter) WithAttributes(...attribute.KeyValue) metric.Float64UpDownCounter {
	return i
}

type testFloat64Histogram struct {
	count int

	embedded.Float64Histogram
}

func (i *testFloat64Histogram) observe() {
	i.count++
}

func (i *testFloat64Histogram) Add(context.Context, float64, ...metric.AddOption) {
	i.count++
}

func (i *testFloat64Histogram) Record(context.Context, float64, ...metric.RecordOption) {
	i.count++
}

func (i *testFloat64Histogram) Enabled(context.Context) bool {
	return true
}

func (i *testFloat64Histogram) WithAttributes(...attribute.KeyValue) metric.Float64Histogram {
	return i
}

type testFloat64Gauge struct {
	count int

	embedded.Float64Gauge
}

func (i *testFloat64Gauge) observe() {
	i.count++
}

func (i *testFloat64Gauge) Add(context.Context, float64, ...metric.AddOption) {
	i.count++
}

func (i *testFloat64Gauge) Record(context.Context, float64, ...metric.RecordOption) {
	i.count++
}

func (i *testFloat64Gauge) Enabled(context.Context) bool {
	return true
}

func (i *testFloat64Gauge) WithAttributes(...attribute.KeyValue) metric.Float64Gauge {
	return i
}

type testFloat64Observable struct {
	count int

	metric.Float64Observable
	embedded.Float64ObservableCounter
	embedded.Float64ObservableUpDownCounter
	embedded.Float64ObservableGauge
}

func (i *testFloat64Observable) observe() {
	i.count++
}

func (i *testFloat64Observable) Add(context.Context, float64, ...metric.AddOption) {
	i.count++
}

func (i *testFloat64Observable) Record(context.Context, float64, ...metric.RecordOption) {
	i.count++
}

func (i *testFloat64Observable) Enabled(context.Context) bool {
	return true
}

func (i *testFloat64Observable) WithAttributes(...attribute.KeyValue) metric.Float64Observable {
	return i
}

type testInt64Counter struct {
	count int

	embedded.Int64Counter
}

func (i *testInt64Counter) observe() {
	i.count++
}

func (i *testInt64Counter) Add(context.Context, int64, ...metric.AddOption) {
	i.count++
}

func (i *testInt64Counter) Record(context.Context, int64, ...metric.RecordOption) {
	i.count++
}

func (i *testInt64Counter) Enabled(context.Context) bool {
	return true
}

func (i *testInt64Counter) WithAttributes(...attribute.KeyValue) metric.Int64Counter {
	return i
}

type testInt64UpDownCounter struct {
	count int

	embedded.Int64UpDownCounter
}

func (i *testInt64UpDownCounter) observe() {
	i.count++
}

func (i *testInt64UpDownCounter) Add(context.Context, int64, ...metric.AddOption) {
	i.count++
}

func (i *testInt64UpDownCounter) Record(context.Context, int64, ...metric.RecordOption) {
	i.count++
}

func (i *testInt64UpDownCounter) Enabled(context.Context) bool {
	return true
}

func (i *testInt64UpDownCounter) WithAttributes(...attribute.KeyValue) metric.Int64UpDownCounter {
	return i
}

type testInt64Histogram struct {
	count int

	embedded.Int64Histogram
}

func (i *testInt64Histogram) observe() {
	i.count++
}

func (i *testInt64Histogram) Add(context.Context, int64, ...metric.AddOption) {
	i.count++
}

func (i *testInt64Histogram) Record(context.Context, int64, ...metric.RecordOption) {
	i.count++
}

func (i *testInt64Histogram) Enabled(context.Context) bool {
	return true
}

func (i *testInt64Histogram) WithAttributes(...attribute.KeyValue) metric.Int64Histogram {
	return i
}

type testInt64Gauge struct {
	count int

	embedded.Int64Gauge
}

func (i *testInt64Gauge) observe() {
	i.count++
}

func (i *testInt64Gauge) Add(context.Context, int64, ...metric.AddOption) {
	i.count++
}

func (i *testInt64Gauge) Record(context.Context, int64, ...metric.RecordOption) {
	i.count++
}

func (i *testInt64Gauge) Enabled(context.Context) bool {
	return true
}

func (i *testInt64Gauge) WithAttributes(...attribute.KeyValue) metric.Int64Gauge {
	return i
}

type testInt64Observable struct {
	count int

	metric.Int64Observable
	embedded.Int64ObservableCounter
	embedded.Int64ObservableUpDownCounter
	embedded.Int64ObservableGauge
}

func (i *testInt64Observable) observe() {
	i.count++
}

func (i *testInt64Observable) Add(context.Context, int64, ...metric.AddOption) {
	i.count++
}

func (i *testInt64Observable) Record(context.Context, int64, ...metric.RecordOption) {
	i.count++
}

func (i *testInt64Observable) WithAttributes(...attribute.KeyValue) metric.Int64Observable {
	return i
}
