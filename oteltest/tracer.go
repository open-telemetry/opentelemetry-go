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
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
)

var _ otel.Tracer = (*Tracer)(nil)

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
func (t *Tracer) Start(ctx context.Context, name string, opts ...otel.SpanOption) (context.Context, otel.Span) {
	c := otel.NewSpanConfig(opts...)
	startTime := time.Now()
	if st := c.Timestamp; !st.IsZero() {
		startTime = st
	}

	span := &Span{
		tracer:     t,
		startTime:  startTime,
		attributes: make(map[label.Key]label.Value),
		links:      make(map[otel.SpanReference][]label.KeyValue),
		spanKind:   c.SpanKind,
	}

	if c.NewRoot {
		span.spanReference = otel.SpanReference{}

		iodKey := label.Key("ignored-on-demand")
		if lsr := otel.SpanFromContext(ctx).SpanReference(); lsr.IsValid() {
			span.links[lsr] = []label.KeyValue{iodKey.String("current")}
		}
		if rsr := otel.RemoteSpanReferenceFromContext(ctx); rsr.IsValid() {
			span.links[rsr] = []label.KeyValue{iodKey.String("remote")}
		}
	} else {
		span.spanReference = t.config.SpanReferenceFunc(ctx)
		if lsr := otel.SpanFromContext(ctx).SpanReference(); lsr.IsValid() {
			span.spanReference.TraceID = lsr.TraceID
			span.parentSpanID = lsr.SpanID
		} else if rsr := otel.RemoteSpanReferenceFromContext(ctx); rsr.IsValid() {
			span.spanReference.TraceID = rsr.TraceID
			span.parentSpanID = rsr.SpanID
		}
	}

	for _, link := range c.Links {
		span.links[link.SpanReference] = link.Attributes
	}

	span.SetName(name)
	span.SetAttributes(c.Attributes...)

	if t.config.SpanRecorder != nil {
		t.config.SpanRecorder.OnStart(span)
	}
	return otel.ContextWithSpan(ctx, span), span
}
