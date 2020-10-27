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

package oteltest

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

// Start starts a MockSpan. It creates a new Span based on Parent SpanReference option.
// TraceID is used from Parent Span Context and SpanID is assigned.
// If Parent SpanReference option is not specified then random TraceID is used.
// No other options are supported.
func (mt *MockTracer) Start(ctx context.Context, name string, o ...otel.SpanOption) (context.Context, otel.Span) {
	config := otel.NewSpanConfig(o...)

	var span *MockSpan
	var sr otel.SpanReference

	parentSpanReference, _, _ := otelparent.GetSpanReferenceAndLinks(ctx, config.NewRoot)

	if !parentSpanReference.IsValid() {
		sr = otel.SpanReference{}
		_, _ = rand.Read(sr.TraceID[:])
		if mt.Sampled {
			sr.TraceFlags = otel.FlagsSampled
		}
	} else {
		sr = parentSpanReference
	}

	binary.BigEndian.PutUint64(sr.SpanID[:], atomic.AddUint64(mt.StartSpanID, 1))
	span = &MockSpan{
		sr:     sr,
		tracer: mt,
		Name:   name,
	}
	if mt.OnSpanStarted != nil {
		mt.OnSpanStarted(span)
	}

	return otel.ContextWithSpan(ctx, span), span
}
