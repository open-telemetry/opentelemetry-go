// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package global

import (
	"sync"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/noop"
	"go.opentelemetry.io/otel/metric"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
)

type nonComparableTracerProvider struct {
	trace.TracerProvider

	nonComparable func() //nolint:unused  // This is not called.
}

type nonComparableMeterProvider struct {
	metric.MeterProvider

	nonComparable func() //nolint:unused  // This is not called.
}

func TestSetTracerProvider(t *testing.T) {
	t.Run("Set With default is a noop", func(t *testing.T) {
		ResetForTest(t)
		SetTracerProvider(TracerProvider())

		tp, ok := TracerProvider().(*tracerProvider)
		if !ok {
			t.Fatal("Global TracerProvider should be the default tracer provider")
		}

		if tp.delegate != nil {
			t.Fatal("tracer provider should not delegate when setting itself")
		}
	})

	t.Run("First Set() should replace the delegate", func(t *testing.T) {
		ResetForTest(t)

		SetTracerProvider(tracenoop.NewTracerProvider())

		_, ok := TracerProvider().(*tracerProvider)
		if ok {
			t.Fatal("Global TracerProvider was not changed")
		}
	})

	t.Run("Set() should delegate existing TracerProviders", func(t *testing.T) {
		ResetForTest(t)

		tp := TracerProvider()
		SetTracerProvider(tracenoop.NewTracerProvider())

		ntp := tp.(*tracerProvider)

		if ntp.delegate == nil {
			t.Fatal("The delegated tracer providers should have a delegate")
		}
	})

	t.Run("non-comparable types should not panic", func(t *testing.T) {
		ResetForTest(t)

		tp := nonComparableTracerProvider{}
		SetTracerProvider(tp)
		assert.NotPanics(t, func() { SetTracerProvider(tp) })
	})
}

func TestSetTextMapPropagator(t *testing.T) {
	t.Run("Set With default is a noop", func(t *testing.T) {
		ResetForTest(t)
		SetTextMapPropagator(TextMapPropagator())

		tmp, ok := TextMapPropagator().(*textMapPropagator)
		if !ok {
			t.Fatal("Global TextMapPropagator should be the default propagator")
		}

		if tmp.delegate != nil {
			t.Fatal("TextMapPropagator should not delegate when setting itself")
		}
	})

	t.Run("First Set() should replace the delegate", func(t *testing.T) {
		ResetForTest(t)

		SetTextMapPropagator(propagation.TraceContext{})

		_, ok := TextMapPropagator().(*textMapPropagator)
		if ok {
			t.Fatal("Global TextMapPropagator was not changed")
		}
	})

	t.Run("Set() should delegate existing propagators", func(t *testing.T) {
		ResetForTest(t)

		p := TextMapPropagator()
		SetTextMapPropagator(propagation.TraceContext{})

		np := p.(*textMapPropagator)

		if np.delegate == nil {
			t.Fatal("The delegated TextMapPropagators should have a delegate")
		}
	})

	t.Run("non-comparable types should not panic", func(t *testing.T) {
		ResetForTest(t)

		// A composite TextMapPropagator is not comparable.
		prop := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{})
		SetTextMapPropagator(prop)
		assert.NotPanics(t, func() { SetTextMapPropagator(prop) })
	})
}

func TestSetMeterProvider(t *testing.T) {
	t.Run("Set With default is a noop", func(t *testing.T) {
		ResetForTest(t)

		SetMeterProvider(MeterProvider())

		mp, ok := MeterProvider().(*meterProvider)
		if !ok {
			t.Fatal("Global MeterProvider should be the default meter provider")
		}

		if mp.delegate != nil {
			t.Fatal("meter provider should not delegate when setting itself")
		}
	})

	t.Run("First Set() should replace the delegate", func(t *testing.T) {
		ResetForTest(t)

		SetMeterProvider(metricnoop.NewMeterProvider())

		_, ok := MeterProvider().(*meterProvider)
		if ok {
			t.Fatal("Global MeterProvider was not changed")
		}
	})

	t.Run("Set() should delegate existing Meter Providers", func(t *testing.T) {
		ResetForTest(t)

		mp := MeterProvider()

		SetMeterProvider(metricnoop.NewMeterProvider())

		dmp := mp.(*meterProvider)

		if dmp.delegate == nil {
			t.Fatal("The delegated meter providers should have a delegate")
		}
	})

	t.Run("non-comparable types should not panic", func(t *testing.T) {
		ResetForTest(t)

		mp := nonComparableMeterProvider{}
		SetMeterProvider(mp)
		assert.NotPanics(t, func() { SetMeterProvider(mp) })
	})
}

func TestSetLoggerProvider(t *testing.T) {
	reset := func() {
		globalLoggerProvider = defaultLoggerProvider()
		delegateLoggerOnce = sync.Once{}
	}

	t.Run("Set With default is a noop", func(t *testing.T) {
		t.Cleanup(reset)

		t.Cleanup(func(orig logr.Logger) func() {
			SetLogger(testr.New(t)) // Don't pollute output.
			return func() { SetLogger(orig) }
		}(GetLogger()))
		SetLoggerProvider(LoggerProvider())

		provider, ok := LoggerProvider().(*loggerProvider)
		if !ok {
			t.Fatal("Global LoggerProvider should be the default logger provider")
		}
		if provider.delegate != nil {
			t.Fatal("logger provider should not delegate when setting itself")
		}
	})

	t.Run("First Set() should replace the delegate", func(t *testing.T) {
		t.Cleanup(reset)

		SetLoggerProvider(noop.NewLoggerProvider())
		if _, ok := LoggerProvider().(*loggerProvider); ok {
			t.Fatal("Global LoggerProvider was not changed")
		}
	})

	t.Run("Set() should delegate existing Logger Providers", func(t *testing.T) {
		t.Cleanup(reset)

		provider := LoggerProvider()
		SetLoggerProvider(noop.NewLoggerProvider())

		if del := provider.(*loggerProvider); del.delegate == nil {
			t.Fatal("The delegated logger providers should have a delegate")
		}
	})

	t.Run("non-comparable types should not panic", func(t *testing.T) {
		t.Cleanup(reset)

		type nonComparableLoggerProvider struct {
			log.LoggerProvider
			noCmp [0]func() //nolint:unused  // This is indeed used.
		}

		provider := nonComparableLoggerProvider{}
		SetLoggerProvider(provider)
		assert.NotPanics(t, func() { SetLoggerProvider(provider) })
	})
}
