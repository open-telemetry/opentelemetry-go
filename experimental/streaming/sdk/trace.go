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

package sdk

import (
	"context"
	"math/rand"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/key"
	"go.opentelemetry.io/api/trace"
	"go.opentelemetry.io/experimental/streaming/exporter"
)

// TODO These should move somewhere in the api, right?
var (
	ErrorKey   = key.New("error")
	SpanIDKey  = key.New("span_id")
	TraceIDKey = key.New("trace_id")
	MessageKey = key.New("message")
)

func (s *sdk) WithSpan(ctx context.Context, name string, body func(context.Context) error) error {
	// TODO: use runtime/trace.WithRegion for execution sdk support
	// TODO: use runtime/pprof.Do for profile tags support
	ctx, span := s.Start(ctx, name)
	defer span.End()

	if err := body(ctx); err != nil {
		span.SetAttribute(ErrorKey.Bool(true))
		span.AddEvent(ctx, "span error", MessageKey.String(err.Error()))
		return err
	}
	return nil
}

func (s *sdk) Start(ctx context.Context, name string, opts ...trace.SpanOption) (context.Context, trace.Span) {
	var child core.SpanContext

	child.SpanID = rand.Uint64()

	o := &trace.SpanOptions{}

	for _, opt := range opts {
		opt(o)
	}

	var parentScope exporter.ScopeID

	if o.Reference.HasTraceID() {
		parentScope.SpanContext = o.Reference.SpanContext
	} else {
		parentScope.SpanContext = trace.CurrentSpan(ctx).SpanContext()
	}

	if parentScope.HasTraceID() {
		parent := parentScope.SpanContext
		child.TraceID.High = parent.TraceID.High
		child.TraceID.Low = parent.TraceID.Low
	} else {
		child.TraceID.High = rand.Uint64()
		child.TraceID.Low = rand.Uint64()
	}

	childScope := exporter.ScopeID{
		SpanContext: child,
		EventID:     s.resources,
	}

	span := &span{
		sdk: s,
		initial: exporter.ScopeID{
			SpanContext: child,
			EventID: s.exporter.Record(exporter.Event{
				Time:    o.StartTime,
				Type:    exporter.START_SPAN,
				Scope:   s.exporter.NewScope(childScope, o.Attributes...),
				Context: ctx,
				Parent:  parentScope,
				String:  name,
			},
			),
		},
	}
	return trace.SetCurrentSpan(ctx, span), span
}
