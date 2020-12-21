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

package trace

import (
	"context"
	"testing"
)

func TestNewNoopTracerProvider(t *testing.T) {
	got, want := NewNoopTracerProvider(), noopTracerProvider{}
	if got != want {
		t.Errorf("NewNoopTracerProvider() returned %#v, want %#v", got, want)
	}
}

func TestNoopTracerProviderTracer(t *testing.T) {
	tp := NewNoopTracerProvider()
	got, want := tp.Tracer(""), noopTracer{}
	if got != want {
		t.Errorf("noopTracerProvider.Tracer() returned %#v, want %#v", got, want)
	}
}

func TestNoopTracerStart(t *testing.T) {
	ctx := context.Background()
	tracer := NewNoopTracerProvider().Tracer("test instrumentation")

	var span Span
	ctx, span = tracer.Start(ctx, "span name")
	got, ok := span.(noopSpan)
	if !ok {
		t.Fatalf("noopTracer.Start() returned a non-noopSpan: %#v", span)
	}
	want := noopSpan{}
	if got != want {
		t.Errorf("noopTracer.Start() returned %#v, want %#v", got, want)
	}
	got, ok = SpanFromContext(ctx).(noopSpan)
	if !ok {
		t.Fatal("noopTracer.Start() did not set span as current in returned context")
	}
	if got != want {
		t.Errorf("noopTracer.Start() current span in returned context set to %#v, want %#v", got, want)
	}
}

func TestNoopSpan(t *testing.T) {
	tracer := NewNoopTracerProvider().Tracer("test instrumentation")
	_, s := tracer.Start(context.Background(), "test span")
	span := s.(noopSpan)

	if got, want := span.SpanContext(), (SpanContext{}); !assertSpanContextEqual(got, want) {
		t.Errorf("span.SpanContext() returned %#v, want %#v", got, want)
	}

	if got, want := span.IsRecording(), false; got != want {
		t.Errorf("span.IsRecording() returned %#v, want %#v", got, want)
	}

	if got, want := span.Tracer(), tracer; got != want {
		t.Errorf("span.Tracer() returned %#v, want %#v", got, want)
	}
}
