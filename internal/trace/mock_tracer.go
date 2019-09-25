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

package trace

import (
	"context"
	"math/rand"
	"sync/atomic"

	"go.opentelemetry.io/api/core"
	apitrace "go.opentelemetry.io/api/trace"
)

// MockTracer is a simple tracer used for testing purpose only.
// It only supports ChildOf option. SpanId is atomically increased every time a
// new span is created.
type MockTracer struct {
	// Sampled specifies if the new span should be sampled or not.
	Sampled bool

	// StartSpanId is used to initialize spanId. It is incremented by one
	// every time a new span is created.
	StartSpanId *uint64
}

var _ apitrace.Tracer = (*MockTracer)(nil)

// WithResources does nothing and returns MockTracer implementation of Tracer.
func (mt *MockTracer) WithResources(attributes ...core.KeyValue) apitrace.Tracer {
	return mt
}

// WithComponent does nothing and returns MockTracer implementation of Tracer.
func (mt *MockTracer) WithComponent(name string) apitrace.Tracer {
	return mt
}

// WithService does nothing and returns MockTracer implementation of Tracer.
func (mt *MockTracer) WithService(name string) apitrace.Tracer {
	return mt
}

// WithSpan does nothing except executing the body.
func (mt *MockTracer) WithSpan(ctx context.Context, name string, body func(context.Context) error) error {
	return body(ctx)
}

// Start starts a MockSpan. It creates a new Span based on Reference SpanContext option.
// TracdID is used from Reference Span Context and SpanID is assigned.
// If Reference SpanContext option is not specified then random TraceID is used.
// No other options are supported.
func (mt *MockTracer) Start(ctx context.Context, name string, o ...apitrace.SpanOption) (context.Context, apitrace.Span) {
	var opts apitrace.SpanOptions
	for _, op := range o {
		op(&opts)
	}
	var span *MockSpan
	var sc core.SpanContext
	if !opts.Reference.SpanContext.IsValid() {
		sc = core.SpanContext{
			TraceID: core.TraceID{
				High: rand.Uint64(),
				Low:  rand.Uint64(),
			},
		}
		if mt.Sampled {
			sc.TraceFlags = core.TraceFlagsSampled
		}
	} else {
		sc = opts.Reference.SpanContext
	}
	sc.SpanID = atomic.AddUint64(mt.StartSpanId, 1)
	span = &MockSpan{
		sc:     sc,
		tracer: mt,
	}

	return apitrace.SetCurrentSpan(ctx, span), span
}
