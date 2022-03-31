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

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func TestSetTracerProvider(t *testing.T) {
	t.Run("Set With default is a noop", func(t *testing.T) {
		ResetForTest(t)
		SetTracerProvider(TracerProvider())

		tp, ok := TracerProvider().(*tracerProvider)
		if !ok {
			t.Fatal("Global Tracer Provider should be the default tracer provider")
		}

		if tp.delegate != nil {
			t.Fatal("tracer provider should not delegate when setting itself")
		}
	})

	t.Run("First Set() should replace the delegate", func(t *testing.T) {
		ResetForTest(t)

		SetTracerProvider(trace.NewNoopTracerProvider())

		_, ok := TracerProvider().(*tracerProvider)
		if ok {
			t.Fatal("Global Tracer Provider was not changed")
		}
	})

	t.Run("Set() should delegate existing Tracer Providers", func(t *testing.T) {
		ResetForTest(t)

		tp := TracerProvider()
		SetTracerProvider(trace.NewNoopTracerProvider())

		ntp := tp.(*tracerProvider)

		if ntp.delegate == nil {
			t.Fatal("The delegated tracer providers should have a delegate")
		}
	})
}

func TestSetTextMapPropagator(t *testing.T) {
	t.Run("Set With default is a noop", func(t *testing.T) {
		ResetForTest(t)
		SetTextMapPropagator(TextMapPropagator())

		tmp, ok := TextMapPropagator().(*textMapPropagator)
		if !ok {
			t.Fatal("Global TextMap Propagator should be the default propagator")
		}

		if tmp.delegate != nil {
			t.Fatal("TextMap propagator should not delegate when setting itself")
		}
	})

	t.Run("First Set() should replace the delegate", func(t *testing.T) {
		ResetForTest(t)

		SetTextMapPropagator(propagation.TraceContext{})

		_, ok := TextMapPropagator().(*textMapPropagator)
		if ok {
			t.Fatal("Global TextMap Propagator was not changed")
		}
	})

	t.Run("Set() should delegate existing propagators", func(t *testing.T) {
		ResetForTest(t)

		p := TextMapPropagator()
		SetTextMapPropagator(propagation.TraceContext{})

		np := p.(*textMapPropagator)

		if np.delegate == nil {
			t.Fatal("The delegated TextMap propagators should have a delegate")
		}
	})
}
