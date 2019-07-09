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

	"github.com/open-telemetry/opentelemetry-go/api/core"
)

type noopTracer struct{}

var _ Tracer = noopTracer{}

// WithResources does nothing and returns noop implementation of Tracer.
func (t noopTracer) WithResources(attributes ...core.KeyValue) Tracer {
	return t
}

// WithComponent does nothing and returns noop implementation of Tracer.
func (t noopTracer) WithComponent(name string) Tracer {
	return t
}

// WithService does nothing and returns noop implementation of Tracer.
func (t noopTracer) WithService(name string) Tracer {
	return t
}

// WithSpan wraps around execution of func with noop span.
func (t noopTracer) WithSpan(ctx context.Context, name string, body func(context.Context) error) error {
	return body(ctx)
}

// Start starts a noop span.
func (noopTracer) Start(ctx context.Context, name string, opts ...SpanOption) (context.Context, Span) {
	span := noopSpan{}
	return SetCurrentSpan(ctx, span), span
}

// Inject does nothing.
func (noopTracer) Inject(ctx context.Context, span Span, injector Injector) {
}
