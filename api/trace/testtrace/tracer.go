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

package testtrace

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel/api/kv/value"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/trace"

	"go.opentelemetry.io/otel/internal/trace/parent"
)

var _ trace.Tracer = (*Tracer)(nil)

// Tracer is a type of OpenTelemetry Tracer that tracks both active and ended spans,
// and which creates Spans that may be inspected to see what data has been set on them.
type Tracer struct {
	lock      *sync.RWMutex
	generator Generator
	spans     []*Span
}

func NewTracer(opts ...TracerOption) *Tracer {
	c := newTracerConfig(opts...)

	return &Tracer{
		lock:      &sync.RWMutex{},
		generator: c.generator,
	}
}

func (t *Tracer) Start(ctx context.Context, name string, opts ...trace.StartOption) (context.Context, trace.Span) {
	var c trace.StartConfig

	for _, opt := range opts {
		opt(&c)
	}

	var traceID trace.ID
	var parentSpanID trace.SpanID

	parentSpanContext, _, links := parent.GetSpanContextAndLinks(ctx, c.NewRoot)

	if parentSpanContext.IsValid() {
		traceID = parentSpanContext.TraceID
		parentSpanID = parentSpanContext.SpanID
	} else {
		traceID = t.generator.TraceID()
	}

	spanID := t.generator.SpanID()

	startTime := time.Now()

	if st := c.StartTime; !st.IsZero() {
		startTime = st
	}

	span := &Span{
		lock:      &sync.RWMutex{},
		tracer:    t,
		startTime: startTime,
		spanContext: trace.SpanContext{
			TraceID: traceID,
			SpanID:  spanID,
		},
		parentSpanID: parentSpanID,
		attributes:   make(map[kv.Key]value.Value),
		links:        make(map[trace.SpanContext][]kv.KeyValue),
	}

	span.SetName(name)
	span.SetAttributes(c.Attributes...)

	for _, link := range links {
		span.links[link.SpanContext] = link.Attributes
	}
	for _, link := range c.Links {
		span.links[link.SpanContext] = link.Attributes
	}

	t.lock.Lock()

	t.spans = append(t.spans, span)

	t.lock.Unlock()

	return trace.ContextWithSpan(ctx, span), span
}

func (t *Tracer) WithSpan(ctx context.Context, name string, body func(ctx context.Context) error, opts ...trace.StartOption) error {
	ctx, _ = t.Start(ctx, name, opts...)

	return body(ctx)
}

// Spans returns the list of current and ended Spans started via the Tracer.
func (t *Tracer) Spans() []*Span {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return append([]*Span{}, t.spans...)
}

// TracerOption enables configuration of a new Tracer.
type TracerOption func(*tracerConfig)

// TracerWithGenerator enables customization of the Generator that the Tracer will use
// to create new trace and span IDs.
// By default, new Tracers will use the CountGenerator.
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
