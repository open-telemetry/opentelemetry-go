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

package global

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/metric"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
)

type nonComparableTracerProvider struct {
	trace.TracerProvider

	nonComparable func() //nolint:structcheck,unused  // This is not called.
}

type nonComparableMeterProvider struct {
	metric.MeterProvider

	nonComparable func() //nolint:structcheck,unused  // This is not called.
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
