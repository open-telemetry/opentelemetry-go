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

package sdk

import (
	"context"

	"go.opentelemetry.io/otel/api/trace"
)

// Tracer is a simple tracer used for testing purpose only.
// It only supports ChildOf option. SpanId is atomically increased every time a
// new span is created.
type Tracer struct {
	// Name is the instrumentation name.
	Name string
	// Version is the instrumentation version.
	Version string

	Config *TracingConfig
}

var _ trace.Tracer = (*Tracer)(nil)

// WithSpan does nothing except executing the body.
func (t *Tracer) WithSpan(ctx context.Context, name string, body func(context.Context) error, opts ...trace.StartOption) error {
	ctx, span := t.Start(ctx, name, opts...)
	defer span.End()

	return body(ctx)
}

// Start starts a Span.
func (t *Tracer) Start(ctx context.Context, name string, opts ...trace.StartOption) (context.Context, trace.Span) {
	var conf trace.StartConfig
	for _, opt := range opts {
		opt(&conf)
	}
	span := &Span{
		sc:         t.Config.SpanContextFunc(ctx),
		tracer:     t,
		Name:       name,
		Attributes: conf.Attributes,
	}
	if t.Config.SpanRecorder != nil {
		t.Config.SpanRecorder.OnStart(span)
	}
	return trace.ContextWithSpan(ctx, span), span
}
