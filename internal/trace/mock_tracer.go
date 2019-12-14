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
	"crypto/rand"
	"encoding/binary"
	"sync/atomic"

	"go.opentelemetry.io/otel/api/core"
	apitrace "go.opentelemetry.io/otel/api/trace"
)

// MockTracer is a simple tracer used for testing purpose only.
// It only supports ChildOf option. SpanId is atomically increased every time a
// new span is created.
type MockTracer struct {
	// Sampled specifies if the new span should be sampled or not.
	Sampled bool

	// StartSpanID is used to initialize spanId. It is incremented by one
	// every time a new span is created.
	StartSpanID *uint64
}

var _ apitrace.Tracer = (*MockTracer)(nil)

// WithSpan does nothing except executing the body.
func (mt *MockTracer) WithSpan(ctx context.Context, name string, body func(context.Context) error) error {
	return body(ctx)
}

// Start starts a MockSpan. It creates a new Span based on Relation SpanContext option.
// TracdID is used from Relation Span Context and SpanID is assigned.
// If Relation SpanContext option is not specified then random TraceID is used.
// No other options are supported.
func (mt *MockTracer) Start(ctx context.Context, name string, o ...apitrace.StartOption) (context.Context, apitrace.Span) {
	var opts apitrace.StartConfig
	for _, op := range o {
		op(&opts)
	}
	var span *MockSpan
	var sc core.SpanContext
	if !opts.Relation.SpanContext.IsValid() {
		sc = core.SpanContext{}
		_, _ = rand.Read(sc.TraceID[:])
		if mt.Sampled {
			sc.TraceFlags = core.TraceFlagsSampled
		}
	} else {
		sc = opts.Relation.SpanContext
	}

	binary.BigEndian.PutUint64(sc.SpanID[:], atomic.AddUint64(mt.StartSpanID, 1))
	span = &MockSpan{
		sc:     sc,
		tracer: mt,
	}

	return apitrace.ContextWithSpan(ctx, span), span
}
