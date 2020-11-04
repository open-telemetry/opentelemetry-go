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

package oteltest // import "go.opentelemetry.io/otel/oteltest"

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"sync/atomic"

	"go.opentelemetry.io/otel"
	otelparent "go.opentelemetry.io/otel/internal/trace/parent"
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

var _ otel.Tracer = (*MockTracer)(nil)

// Start starts a MockSpan. It creates a new Span based on Parent SpanContext option.
// TraceID is used from Parent Span Context and SpanID is assigned.
// If Parent SpanContext option is not specified then random TraceID is used.
// No other options are supported.
func (mt *MockTracer) Start(ctx context.Context, name string, o ...otel.SpanOption) (context.Context, otel.Span) {
	config := otel.NewSpanConfig(o...)

	var span *MockSpan
	var sc otel.SpanContext

	parentSpanContext, _, _ := otelparent.GetSpanContextAndLinks(ctx, config.NewRoot)

	if !parentSpanContext.IsValid() {
		sc = otel.SpanContext{}
		_, _ = rand.Read(sc.TraceID[:])
		if mt.Sampled {
			sc.TraceFlags = otel.FlagsSampled
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

	return otel.ContextWithSpan(ctx, span), span
}
