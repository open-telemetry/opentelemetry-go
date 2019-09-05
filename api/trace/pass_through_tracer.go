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

	"go.opentelemetry.io/api/core"
)

// PassThroughTracer is equivalent of noop tracer except that
// it facilitates forwarding incoming trace to downstream services.
// It does require to use appropriate propagators.
type PassThroughTracer struct{}

var _ Tracer = PassThroughTracer{}

// WithResources does nothing and returns noop implementation of Tracer.
func (t PassThroughTracer) WithResources(attributes ...core.KeyValue) Tracer {
	return t
}

// WithComponent does nothing and returns noop implementation of Tracer.
func (t PassThroughTracer) WithComponent(name string) Tracer {
	return t
}

// WithService does nothing and returns noop implementation of Tracer.
func (t PassThroughTracer) WithService(name string) Tracer {
	return t
}

// WithSpan wraps around execution of func with noop span.
func (t PassThroughTracer) WithSpan(ctx context.Context, name string, body func(context.Context) error) error {
	return body(ctx)
}

// Start starts a PassThroughSpan. It simply creates a copy of remote span.
// If RemoteSpanContext is not provided then it returns a NoopSpan.
func (PassThroughTracer) Start(ctx context.Context, name string, o ...SpanOption) (context.Context, Span) {
	var opts SpanOptions
	for _, op := range o {
		op(&opts)
	}
	if !opts.RemoteSpanContext.IsValid() {
		return ctx, NoopSpan{}
	}
	span := &PassThroughSpan{sc: opts.RemoteSpanContext}
	return SetCurrentSpan(ctx, span), span
}
