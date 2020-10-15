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

package internal

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel"
)

// textMapPropagator is a default TextMapPropagator that delegates calls to a
// registered delegate if one is set, otherwise it defaults to delegating the
// calls to a the default no-op otel.TextMapPropagator.
type textMapPropagator struct {
	mtx      sync.Mutex
	once     sync.Once
	delegate otel.TextMapPropagator
	noop     otel.TextMapPropagator
}

// Compile-time guarantee that textMapPropagator implements the
// otel.TextMapPropagator interface.
var _ otel.TextMapPropagator = (*textMapPropagator)(nil)

func newTextMapPropagator() *textMapPropagator {
	return &textMapPropagator{
		noop: otel.NewCompositeTextMapPropagator(),
	}
}

// SetDelegate sets a delegate otel.TextMapPropagator that all calls are
// forwarded to. Delegation can only be performed once, all subsequent calls
// perform no delegation.
func (p *textMapPropagator) SetDelegate(delegate otel.TextMapPropagator) {
	if delegate == nil {
		return
	}

	p.mtx.Lock()
	p.once.Do(func() { p.delegate = delegate })
	p.mtx.Unlock()
}

// HasDelegate returns if a delegate is set for p.
func (p *textMapPropagator) HasDelegate() bool {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	return p.delegate != nil
}

// Inject set cross-cutting concerns from the Context into the carrier.
func (p *textMapPropagator) Inject(ctx context.Context, carrier otel.TextMapCarrier) {
	if p.HasDelegate() {
		p.delegate.Inject(ctx, carrier)
	}
	p.noop.Inject(ctx, carrier)
}

// Extract reads cross-cutting concerns from the carrier into a Context.
func (p *textMapPropagator) Extract(ctx context.Context, carrier otel.TextMapCarrier) context.Context {
	if p.HasDelegate() {
		return p.delegate.Extract(ctx, carrier)
	}
	return p.noop.Extract(ctx, carrier)
}

// Fields returns the keys who's values are set with Inject.
func (p *textMapPropagator) Fields() []string {
	if p.HasDelegate() {
		return p.delegate.Fields()
	}
	return p.noop.Fields()
}
