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

package otel

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
)

type NoopMeterProvider struct{}

type noopInstrument struct{}
type noopBoundInstrument struct{}
type NoopSync struct{ noopInstrument }
type NoopAsync struct{ noopInstrument }

var _ MeterProvider = NoopMeterProvider{}
var _ SyncImpl = NoopSync{}
var _ BoundSyncImpl = noopBoundInstrument{}
var _ AsyncImpl = NoopAsync{}

func (NoopMeterProvider) Meter(_ string, _ ...MeterOption) Meter {
	return Meter{}
}

func (noopInstrument) Implementation() interface{} { return nil }

func (noopInstrument) Descriptor() Descriptor { return Descriptor{} }

func (noopBoundInstrument) RecordOne(context.Context, Number) {}

func (noopBoundInstrument) Unbind() {}

func (NoopSync) Bind([]label.KeyValue) BoundSyncImpl { return noopBoundInstrument{} }

func (NoopSync) RecordOne(context.Context, Number, []label.KeyValue) {}

type noopSpan struct{}

var _ Span = noopSpan{}

// SpanContext returns an invalid span context.
func (noopSpan) SpanContext() SpanContext { return EmptySpanContext() }

// IsRecording always returns false for NoopSpan.
func (noopSpan) IsRecording() bool { return false }

// SetStatus does nothing.
func (noopSpan) SetStatus(status codes.Code, msg string) {}

// SetError does nothing.
func (noopSpan) SetError(v bool) {}

// SetAttributes does nothing.
func (noopSpan) SetAttributes(attributes ...label.KeyValue) {}

// SetAttribute does nothing.
func (noopSpan) SetAttribute(k string, v interface{}) {}

// End does nothing.
func (noopSpan) End(options ...SpanOption) {}

// RecordError does nothing.
func (noopSpan) RecordError(ctx context.Context, err error, opts ...ErrorOption) {}

// Tracer returns noop implementation of Tracer.
func (noopSpan) Tracer() Tracer { return noopTracer{} }

// AddEvent does nothing.
func (noopSpan) AddEvent(ctx context.Context, name string, attrs ...label.KeyValue) {}

// AddEventWithTimestamp does nothing.
func (noopSpan) AddEventWithTimestamp(ctx context.Context, timestamp time.Time, name string, attrs ...label.KeyValue) {
}

// SetName does nothing.
func (noopSpan) SetName(name string) {}

type noopTracer struct{}

var _ Tracer = noopTracer{}

// Start starts a noop span.
func (noopTracer) Start(ctx context.Context, name string, opts ...SpanOption) (context.Context, Span) {
	span := noopSpan{}
	return ContextWithSpan(ctx, span), span
}

type noopProvider struct{}

var _ TracerProvider = noopProvider{}

// Tracer returns noop implementation of Tracer.
func (p noopProvider) Tracer(_ string, _ ...TracerOption) Tracer {
	return noopTracer{}
}

// NoopProvider returns a noop implementation of TracerProvider.
// The Tracer and Spans created from the noop provider will
// also be noop.
func NoopProvider() TracerProvider {
	return noopProvider{}
}
