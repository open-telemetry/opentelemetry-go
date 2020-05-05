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
	"crypto/rand"
	"encoding/binary"
	"sync/atomic"

	apitrace "go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/internal/trace/parent"
)

// MockTracer is a simple tracer used for testing purpose only.
// It only supports ChildOf option. SpanId is atomically increased every time a
// new span is created.
type MockTracer struct {
	// StartSpanID is used to initialize spanId. It is incremented by one
	// every time a new span is created.
	//
	// StartSpanID has to be aligned for 64-bit atomic operations.
	StartSpanID *uint64

	// Sampled specifies if the new span should be sampled or not.
	Sampled bool

	// OnSpanStarted is called every time a new trace span is started
	OnSpanStarted func(span *MockSpan)
}

var _ apitrace.Tracer = (*MockTracer)(nil)

// WithSpan does nothing except executing the body.
func (mt *MockTracer) WithSpan(ctx context.Context, name string, body func(context.Context) error, opts ...apitrace.StartOption) error {
	ctx, span := mt.Start(ctx, name, opts...)
	defer span.End()

	return body(ctx)
}

// Start starts a MockSpan. It creates a new Span based on Parent SpanContext option.
// TracdID is used from Parent Span Context and SpanID is assigned.
// If Parent SpanContext option is not specified then random TraceID is used.
// No other options are supported.
func (mt *MockTracer) Start(ctx context.Context, name string, o ...apitrace.StartOption) (context.Context, apitrace.Span) {
	var opts apitrace.StartConfig
	for _, op := range o {
		op(&opts)
	}
	var span *MockSpan
	var sc apitrace.SpanContext

	parentSpanContext, _, _ := parent.GetSpanContextAndLinks(ctx, opts.NewRoot)

	if !parentSpanContext.IsValid() {
		sc = apitrace.SpanContext{}
		_, _ = rand.Read(sc.TraceID[:])
		if mt.Sampled {
			sc.TraceFlags = apitrace.FlagsSampled
		}
	} else {
		sc = parentSpanContext
	}

	binary.BigEndian.PutUint64(sc.SpanID[:], atomic.AddUint64(mt.StartSpanID, 1))
	span = &MockSpan{
		sc:     sc,
		tracer: mt,
		Name:   name,
	}
	if mt.OnSpanStarted != nil {
		mt.OnSpanStarted(span)
	}

	return apitrace.ContextWithSpan(ctx, span), span
}
