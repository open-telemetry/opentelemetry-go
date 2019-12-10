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
	"sync/atomic"

	"go.opentelemetry.io/otel/api/core"
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

type Propagators interface {
	// HTTP propagation
	SetHTTPExtractors(...HTTPExtractor)
	SetHTTPInjectors(...HTTPInjector)
	HTTPExtractors() []HTTPExtractor
	HTTPInjectors() []HTTPInjector

	// Binary propagation
	// TODO
}

type propagators struct {
	// TODO Nah.  Use an options pattern to avoid mutation.

	httpEx atomic.Value // []HTTPExtractor
	httpIn atomic.Value // []HTTPInjector
}

func New() Propagators {
	p := &propagators{}

	p.httpIn.Store([]HTTPInjector{})
	p.httpEs.Store([]HTTPExtractor{})

	return p
}

func (p *Propagators) SetHTTPExtractors(ex ...HTTPExtractor) {
	if ex == nil {
		ex = []HTTPExtractor{}
	}
	p.httpEx.Store(ex)
}

func (p *Propagators) SetHTTPInjectors(in ...HTTPInjector) {
	if in == nil {
		in = []HTTPInjector{}
	}
	p.httpIn.Store(in)
}

func (p *Propagators) HTTPExtractors() []HTTPExtractor {
	return p.httpEx.Load().([]HTTPExtractor)
}

func (p *Propagators) HTTPInjectors() []HTTPInjector {
	return p.httpIn.Load().([]HTTPInjector)
}

type NoopPropagators struct{}

func (NoopPropagators) SetHTTPExtractors(ex ...HTTPExtractor) {}
func (NoopPropagators) SetHTTPInjectors(in ...HTTPInjector)   {}
func (NoopPropagators) HTTPExtractors() []HTTPExtractor       { return nil }
func (NoopPropagators) HTTPInjectors() []HTTPInjector         { return nil }
