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
	attributes := []attribute.KeyValue{attribute.String("foo", "bar")}
	// Float64 Instruments
	t.Run("Float64", func(*testing.T) {
		t.Run("Counter", func(*testing.T) {
			delegate := &sfCounter{}
			f := func(v float64) { delegate.WithAttributes(attributes...).Add(t.Context(), v) }
			testFloat64ConcurrentSafe(f, delegate.setDelegate)
		})

		t.Run("UpDownCounter", func(*testing.T) {
			delegate := &sfUpDownCounter{}
			f := func(v float64) { delegate.WithAttributes(attributes...).Add(t.Context(), v) }
			testFloat64ConcurrentSafe(f, delegate.setDelegate)
		})

		t.Run("Histogram", func(*testing.T) {
			delegate := &sfHistogram{}
			f := func(v float64) { delegate.WithAttributes(attributes...).Record(t.Context(), v) }
			testFloat64ConcurrentSafe(f, delegate.setDelegate)
		})

		t.Run("Gauge", func(*testing.T) {
			delegate := &sfGauge{}
			f := func(v float64) { delegate.WithAttributes(attributes...).Record(t.Context(), v) }
			testFloat64ConcurrentSafe(f, delegate.setDelegate)
		})
	})

	// Int64 Instruments

	t.Run("Int64", func(*testing.T) {
		t.Run("Counter", func(*testing.T) {
			delegate := &siCounter{}
			f := func(v int64) { delegate.WithAttributes(attributes...).Add(t.Context(), v) }
			testInt64ConcurrentSafe(f, delegate.setDelegate)
		})

		t.Run("UpDownCounter", func(*testing.T) {
			delegate := &siUpDownCounter{}
			f := func(v int64) { delegate.WithAttributes(attributes...).Add(t.Context(), v) }
			testInt64ConcurrentSafe(f, delegate.setDelegate)
		})

		t.Run("Histogram", func(*testing.T) {
			delegate := &siHistogram{}
			f := func(v int64) { delegate.WithAttributes(attributes...).Record(t.Context(), v) }
			testInt64ConcurrentSafe(f, delegate.setDelegate)
		})

		t.Run("Gauge", func(*testing.T) {
			delegate := &siGauge{}
			f := func(v int64) { delegate.WithAttributes(attributes...).Record(t.Context(), v) }
			testInt64ConcurrentSafe(f, delegate.setDelegate)
		})
	})
}

type testFloat64Counter struct {
	count      int
	attributes []attribute.KeyValue

	embedded.Float64Counter
}

func (i *testFloat64Counter) Add(context.Context, float64, ...metric.AddOption) {
	i.count++
}

func (i *testFloat64Counter) Record(context.Context, float64, ...metric.RecordOption) {
	i.count++
}

func (*testFloat64Counter) Enabled(context.Context) bool {
	return true
}

func (i *testFloat64Counter) WithAttributes(kvs ...attribute.KeyValue) metric.Float64Counter {
	return &testFloat64Counter{count: i.count, attributes: append(i.attributes, kvs...)}
}

type testFloat64UpDownCounter struct {
	count      int
	attributes []attribute.KeyValue

	embedded.Float64UpDownCounter
}

func (i *testFloat64UpDownCounter) Add(context.Context, float64, ...metric.AddOption) {
	i.count++
}

func (i *testFloat64UpDownCounter) Record(context.Context, float64, ...metric.RecordOption) {
	i.count++
}

func (*testFloat64UpDownCounter) Enabled(context.Context) bool {
	return true
}

func (i *testFloat64UpDownCounter) WithAttributes(kvs ...attribute.KeyValue) metric.Float64UpDownCounter {
	return &testFloat64UpDownCounter{count: i.count, attributes: append(i.attributes, kvs...)}
}

type testFloat64Histogram struct {
	count      int
	attributes []attribute.KeyValue

	embedded.Float64Histogram
}

func (i *testFloat64Histogram) Add(context.Context, float64, ...metric.AddOption) {
	i.count++
}

func (i *testFloat64Histogram) Record(context.Context, float64, ...metric.RecordOption) {
	i.count++
}

func (*testFloat64Histogram) Enabled(context.Context) bool {
	return true
}

func (i *testFloat64Histogram) WithAttributes(kvs ...attribute.KeyValue) metric.Float64Histogram {
	return &testFloat64Histogram{count: i.count, attributes: append(i.attributes, kvs...)}
}

type testFloat64Gauge struct {
	count      int
	attributes []attribute.KeyValue

	embedded.Float64Gauge
}

func (i *testFloat64Gauge) Add(context.Context, float64, ...metric.AddOption) {
	i.count++
}

func (i *testFloat64Gauge) Record(context.Context, float64, ...metric.RecordOption) {
	i.count++
}

func (*testFloat64Gauge) Enabled(context.Context) bool {
	return true
}

func (i *testFloat64Gauge) WithAttributes(kvs ...attribute.KeyValue) metric.Float64Gauge {
	return &testFloat64Gauge{count: i.count, attributes: append(i.attributes, kvs...)}
}

type testFloat64Observable struct {
	count      int
	attributes []attribute.KeyValue

	metric.Float64Observable
	embedded.Float64ObservableCounter
	embedded.Float64ObservableUpDownCounter
	embedded.Float64ObservableGauge
}

func (i *testFloat64Observable) observe(kvs ...attribute.KeyValue) {
	i.count++
	i.attributes = kvs
}

func (i *testFloat64Observable) Add(context.Context, float64, ...metric.AddOption) {
	i.count++
}

func (i *testFloat64Observable) Record(context.Context, float64, ...metric.RecordOption) {
	i.count++
}

func (*testFloat64Observable) Enabled(context.Context) bool {
	return true
}

func (i *testFloat64Observable) WithAttributes(kvs ...attribute.KeyValue) metric.Float64Observable {
	return &testFloat64Observable{count: i.count, attributes: append(i.attributes, kvs...)}
}

type testInt64Counter struct {
	count      int
	attributes []attribute.KeyValue

	embedded.Int64Counter
}

func (i *testInt64Counter) Add(context.Context, int64, ...metric.AddOption) {
	i.count++
}

func (i *testInt64Counter) Record(context.Context, int64, ...metric.RecordOption) {
	i.count++
}

func (*testInt64Counter) Enabled(context.Context) bool {
	return true
}

func (i *testInt64Counter) WithAttributes(kvs ...attribute.KeyValue) metric.Int64Counter {
	return &testInt64Counter{count: i.count, attributes: append(i.attributes, kvs...)}
}

type testInt64UpDownCounter struct {
	count      int
	attributes []attribute.KeyValue

	embedded.Int64UpDownCounter
}

func (i *testInt64UpDownCounter) Add(context.Context, int64, ...metric.AddOption) {
	i.count++
}

func (i *testInt64UpDownCounter) Record(context.Context, int64, ...metric.RecordOption) {
	i.count++
}

func (*testInt64UpDownCounter) Enabled(context.Context) bool {
	return true
}

func (i *testInt64UpDownCounter) WithAttributes(kvs ...attribute.KeyValue) metric.Int64UpDownCounter {
	return &testInt64UpDownCounter{count: i.count, attributes: append(i.attributes, kvs...)}
}

type testInt64Histogram struct {
	count      int
	attributes []attribute.KeyValue

	embedded.Int64Histogram
}

func (i *testInt64Histogram) Add(context.Context, int64, ...metric.AddOption) {
	i.count++
}

func (i *testInt64Histogram) Record(context.Context, int64, ...metric.RecordOption) {
	i.count++
}

func (*testInt64Histogram) Enabled(context.Context) bool {
	return true
}

func (i *testInt64Histogram) WithAttributes(kvs ...attribute.KeyValue) metric.Int64Histogram {
	return &testInt64Histogram{count: i.count, attributes: append(i.attributes, kvs...)}
}

type testInt64Gauge struct {
	count      int
	attributes []attribute.KeyValue

	embedded.Int64Gauge
}

func (i *testInt64Gauge) Add(context.Context, int64, ...metric.AddOption) {
	i.count++
}

func (i *testInt64Gauge) Record(context.Context, int64, ...metric.RecordOption) {
	i.count++
}

func (*testInt64Gauge) Enabled(context.Context) bool {
	return true
}

func (i *testInt64Gauge) WithAttributes(kvs ...attribute.KeyValue) metric.Int64Gauge {
	return &testInt64Gauge{count: i.count, attributes: append(i.attributes, kvs...)}
}

type testInt64Observable struct {
	count      int
	attributes []attribute.KeyValue

	metric.Int64Observable
	embedded.Int64ObservableCounter
	embedded.Int64ObservableUpDownCounter
	embedded.Int64ObservableGauge
}

func (i *testInt64Observable) observe(kvs ...attribute.KeyValue) {
	i.count++
	i.attributes = kvs
}

func (i *testInt64Observable) Add(context.Context, int64, ...metric.AddOption) {
	i.count++
}

func (i *testInt64Observable) Record(context.Context, int64, ...metric.RecordOption) {
	i.count++
}

func (i *testInt64Observable) WithAttributes(kvs ...attribute.KeyValue) metric.Int64Observable {
	return &testInt64Observable{count: i.count, attributes: append(i.attributes, kvs...)}
}
