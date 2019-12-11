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

package propagation

import (
	"context"

	"go.opentelemetry.io/otel/api/core"
	dctx "go.opentelemetry.io/otel/api/distributedcontext"
)

// HTTPSupplier is implemented by http.Headers.
type HTTPSupplier interface {
	Get(key string) string
	Set(key string, value string)
}

type HTTPExtractor interface {
	// Extract method retrieves encoded SpanContext using supplier
	// from the associated carrier.  It decodes the SpanContext
	// and returns it and a dctx of correlated context.  If no
	// SpanContext was retrieved OR if the retrieved SpanContext
	// is invalid then an empty SpanContext is returned.
	Extract(context.Context, HTTPSupplier) (core.SpanContext, dctx.Map)
}

type HTTPInjector interface {
	// Inject method retrieves current SpanContext from the ctx,
	// encodes it into propagator specific format and then injects
	// the encoded SpanContext using supplier into a carrier
	// associated with the supplier. It also takes a
	// correlationCtx whose values will be injected into a carrier
	// using the supplier.
	Inject(context.Context, HTTPSupplier) context.Context
}

type Config struct {
	httpEx []HTTPExtractor
	httpIn []HTTPInjector
}

type Propagators interface {
	HTTPExtractors() []HTTPExtractor
	HTTPInjectors() []HTTPInjector
}

type Option func(*Config)

type propagators struct {
	config Config
}

func New(options ...Option) Propagators {
	config := Config{}
	for _, opt := range options {
		opt(&config)
	}
	return &propagators{
		config: config,
	}
}

func WithInjectors(inj ...HTTPInjector) Option {
	return func(config *Config) {
		config.httpIn = append(config.httpIn, inj...)
	}
}

func WithExtractors(ext ...HTTPExtractor) Option {
	return func(config *Config) {
		config.httpEx = append(config.httpEx, ext...)
	}
}

func (p *propagators) HTTPExtractors() []HTTPExtractor {
	return p.config.httpEx
}

func (p *propagators) HTTPInjectors() []HTTPInjector {
	return p.config.httpIn
}

type NoopPropagators struct{}

func (NoopPropagators) HTTPExtractors() []HTTPExtractor { return nil }
func (NoopPropagators) HTTPInjectors() []HTTPInjector   { return nil }
