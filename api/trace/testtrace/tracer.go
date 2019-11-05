// Copyright 2019, OpenTelemetry Authors
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

package testtrace

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/trace"
)

var _ trace.Tracer = (*Tracer)(nil)

type Tracer struct {
	lock      *sync.Mutex
	generator Generator
	spans     []*Span
}

func NewTracer(opts ...TracerOption) *Tracer {
	c := newTracerConfig(opts...)

	return &Tracer{
		lock:      &sync.Mutex{},
		generator: c.generator,
	}
}

func (t *Tracer) Start(ctx context.Context, name string, opts ...trace.SpanOption) (context.Context, trace.Span) {
	var c trace.SpanOptions

	for _, opt := range opts {
		opt(&c)
	}

	var traceID core.TraceID

	if parentSpanContext := c.Relation.SpanContext; parentSpanContext.IsValid() {
		traceID = parentSpanContext.TraceID
	} else if parentSpanContext := trace.CurrentSpan(ctx).SpanContext(); parentSpanContext.IsValid() {
		traceID = parentSpanContext.TraceID
	} else {
		traceID = t.generator.TraceID()
	}

	spanID := t.generator.SpanID()

	span := &Span{
		lock:   &sync.Mutex{},
		tracer: t,
		spanContext: core.SpanContext{
			TraceID: traceID,
			SpanID:  spanID,
		},
		name:       name,
		attributes: c.Attributes,
		startTime:  time.Now(),
	}

	t.lock.Lock()

	t.spans = append(t.spans, span)

	t.lock.Unlock()

	return trace.SetCurrentSpan(ctx, span), span
}

func (t *Tracer) WithSpan(ctx context.Context, name string, body func(ctx context.Context) error) error {
	ctx, _ = t.Start(ctx, name)

	return body(ctx)
}

func (t *Tracer) Spans() []*Span {
	return t.spans
}

type TracerOption func(*tracerConfig)

func TracerWithGenerator(generator Generator) TracerOption {
	return func(c *tracerConfig) {
		c.generator = generator
	}
}

type tracerConfig struct {
	generator Generator
}

func newTracerConfig(opts ...TracerOption) tracerConfig {
	var c tracerConfig
	defaultOpts := []TracerOption{
		TracerWithGenerator(NewCountGenerator()),
	}

	for _, opt := range append(defaultOpts, opts...) {
		opt(&c)
	}

	return c
}
