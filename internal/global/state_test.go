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
	t.Cleanup(ResetForTest)

	t.Run("Set With default is a noop", func(t *testing.T) {
		ResetForTest()
		SetTracerProvider(TracerProvider())

		_, ok := TracerProvider().(*tracerProvider)
		if !ok {
			t.Error("Global Tracer Provider should be the default tracer provider")
			return
		}
	})

	t.Run("First Set() should replace the delegate", func(t *testing.T) {
		ResetForTest()

		SetTracerProvider(trace.NewNoopTracerProvider())

		_, ok := TracerProvider().(*tracerProvider)
		if ok {
			t.Error("Global Tracer Provider was not changed")
			return
		}
	})

	t.Run("Set() should delegate existing Tracer Providers", func(t *testing.T) {
		ResetForTest()

		tp := TracerProvider()
		SetTracerProvider(trace.NewNoopTracerProvider())

		ntp := tp.(*tracerProvider)

		if ntp.delegate == nil {
			t.Error("The delegated tracer providers should have a delegate")
		}
	})
}

func TestSetTextMapPropagator(t *testing.T) {
	t.Cleanup(ResetForTest)

	t.Run("Set With default is a noop", func(t *testing.T) {
		ResetForTest()
		SetTextMapPropagator(TextMapPropagator())

		_, ok := TextMapPropagator().(*textMapPropagator)
		if !ok {
			t.Error("Global TextMap Propagator should be the default propagator")
			return
		}
	})

	t.Run("First Set() should replace the delegate", func(t *testing.T) {
		ResetForTest()

		SetTextMapPropagator(propagation.TraceContext{})

		_, ok := TextMapPropagator().(*textMapPropagator)
		if ok {
			t.Error("Global TextMap Propagator was not changed")
			return
		}
	})

	t.Run("Set() should delegate existing propagators", func(t *testing.T) {
		ResetForTest()

		p := TextMapPropagator()
		SetTextMapPropagator(propagation.TraceContext{})

		np := p.(*textMapPropagator)

		if np.delegate == nil {
			t.Error("The delegated TextMap propagators should have a delegate")
		}
	})
}
