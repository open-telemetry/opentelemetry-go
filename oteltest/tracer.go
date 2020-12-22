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
	"time"

	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
)

var _ trace.Tracer = (*Tracer)(nil)

// Tracer is an OpenTelemetry Tracer implementation used for testing.
type Tracer struct {
	// Name is the instrumentation name.
	Name string
	// Version is the instrumentation version.
	Version string

	config *config
}

// Start creates a span. If t is configured with a SpanRecorder its OnStart
// method will be called after the created Span has been initialized.
func (t *Tracer) Start(ctx context.Context, name string, opts ...trace.SpanOption) (context.Context, trace.Span) {
	c := trace.NewSpanConfig(opts...)
	startTime := time.Now()
	if st := c.Timestamp; !st.IsZero() {
		startTime = st
	}

	span := &Span{
		tracer:     t,
		startTime:  startTime,
		attributes: make(map[label.Key]label.Value),
		links:      []trace.Link{},
		spanKind:   c.SpanKind,
	}

	if c.NewRoot {
		span.spanContext = trace.SpanContext{}

		iodKey := label.Key("ignored-on-demand")
		if lsc := trace.SpanContextFromContext(ctx); lsc.IsValid() {
			span.links = append(span.links, trace.Link{
				SpanContext: lsc,
				Attributes:  []label.KeyValue{iodKey.String("current")},
			})
		}
		if rsc := trace.RemoteSpanContextFromContext(ctx); rsc.IsValid() {
			span.links = append(span.links, trace.Link{
				SpanContext: rsc,
				Attributes:  []label.KeyValue{iodKey.String("remote")},
			})
		}
	} else {
		span.spanContext = t.config.SpanContextFunc(ctx)
		if lsc := trace.SpanContextFromContext(ctx); lsc.IsValid() {
			span.spanContext.TraceID = lsc.TraceID
			span.parentSpanID = lsc.SpanID
		} else if rsc := trace.RemoteSpanContextFromContext(ctx); rsc.IsValid() {
			span.spanContext.TraceID = rsc.TraceID
			span.parentSpanID = rsc.SpanID
		}
	}

	for _, link := range c.Links {
		for i, sl := range span.links {
			if sl.SpanContext.SpanID == link.SpanContext.SpanID &&
				sl.SpanContext.TraceID == link.SpanContext.TraceID &&
				sl.SpanContext.TraceFlags == link.SpanContext.TraceFlags &&
				sl.SpanContext.TraceState.String() == link.SpanContext.TraceState.String() {
				span.links[i].Attributes = link.Attributes
				break
			}
		}
		span.links = append(span.links, link)
	}

	span.SetName(name)
	span.SetAttributes(c.Attributes...)

	if t.config.SpanRecorder != nil {
		t.config.SpanRecorder.OnStart(span)
	}
	return trace.ContextWithSpan(ctx, span), span
}
