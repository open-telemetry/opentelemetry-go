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
)

// HTTPSupplier is implemented by http.Headers.
type HTTPSupplier interface {
	Get(key string) string
	Set(key string, value string)
}

// HTTPExtractor extracts from a HTTPSupplier.
type HTTPExtractor interface {
	// Extract method retrieves encoded SpanContext using supplier
	// from the associated carrier.  It decodes the SpanContext
	// and returns it and a dctx of correlated context.  If no
	// SpanContext was retrieved OR if the retrieved SpanContext
	// is invalid then an empty SpanContext is returned.
	Extract(context.Context, HTTPSupplier) context.Context
}

// HTTPExtractor injects into a HTTPSupplier.
type HTTPInjector interface {
	// Inject method retrieves current SpanContext from the ctx,
	// encodes it into propagator specific format and then injects
	// the encoded SpanContext using supplier into a carrier
	// associated with the supplier. It also takes a
	// correlationCtx whose values will be injected into a carrier
	// using the supplier.
	Inject(context.Context, HTTPSupplier)
}

// Config contains the current set of extractors and injectors.
type Config struct {
	httpEx []HTTPExtractor
	httpIn []HTTPInjector
}

// Propagators is the interface to a set of injectors and extractors
// for all supported carrier formats.
type Propagators interface {
	// HTTPExtractors returns the configured extractors.
	HTTPExtractors() []HTTPExtractor

	// HTTPInjectors returns the configured injectors.
	HTTPInjectors() []HTTPInjector
}

// HTTPPropagator is the interface to inject and extract to and from
// http.Headers.
type HTTPPropagator interface {
	HTTPInjector
	HTTPExtractor

	// GetAllKeys returns the HTTP header names used.
	GetAllKeys() []string
}

// Option support passing configuration parameters to New().
type Option func(*Config)

// propagators is the default Propagators implementation.
type propagators struct {
	config Config
}

// New returns a standard Propagators implementation.
func New(options ...Option) Propagators {
	config := Config{}
	for _, opt := range options {
		opt(&config)
	}
	return &propagators{
		config: config,
	}
}

// WithInjectors appends to the optional injector set.
func WithInjectors(inj ...HTTPInjector) Option {
	return func(config *Config) {
		config.httpIn = append(config.httpIn, inj...)
	}
}

// WithExtractors appends to the optional extractor set.
func WithExtractors(ext ...HTTPExtractor) Option {
	return func(config *Config) {
		config.httpEx = append(config.httpEx, ext...)
	}
}

// HTTPExtractors implements Propagators.
func (p *propagators) HTTPExtractors() []HTTPExtractor {
	return p.config.httpEx
}

// HTTPExtractors implements Propagators.
func (p *propagators) HTTPInjectors() []HTTPInjector {
	return p.config.httpIn
}

// ExtractHTTP applies props.HTTPExtractors() to the supplier and
// returns the combined result Context.
func ExtractHTTP(ctx context.Context, props Propagators, supplier HTTPSupplier) context.Context {
	for _, ex := range props.HTTPExtractors() {
		ctx = ex.Extract(ctx, supplier)
	}
	return ctx
}

// ExtractHTTP applies props.HTTPInjectors() and returns the combined
// to the supplier.
func InjectHTTP(ctx context.Context, props Propagators, supplier HTTPSupplier) {
	for _, in := range props.HTTPInjectors() {
		in.Inject(ctx, supplier)
	}
}
