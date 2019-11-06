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

package http

import (
	"go.opentelemetry.io/otel/api/core"
	dctx "go.opentelemetry.io/otel/api/distributedcontext"
)

// Supplier is an interface that specifies methods to retrieve and store
// value for a key to an associated carrier.
// Get method retrieves the value for a given key.
// Set method stores the value for a given key.
type Supplier interface {
	Get(key string) string
	Set(key, value string)
}

// SpanContextPropagator is an interface that specifies methods to
// convert SpanContext to/from byte array.
type SpanContextPropagator interface {
	Inject(core.SpanContext, Supplier)
	Extract(Supplier) core.SpanContext
}

type chainSpanContextPropagator struct {
	propagators []SpanContextPropagator
}

var _ SpanContextPropagator = chainSpanContextPropagator{}

func NewChainSpanContextPropagator(propagators ...SpanContextPropagator) SpanContextPropagator {
	return chainSpanContextPropagator{
		propagators: propagators,
	}
}

func (c chainSpanContextPropagator) Inject(sc core.SpanContext, supplier Supplier) {
	for _, p := range c.propagators {
		p.Inject(sc, supplier)
	}
}

func (c chainSpanContextPropagator) Extract(supplier Supplier) core.SpanContext {
	for _, p := range c.propagators {
		if sc := p.Extract(supplier); sc.IsValid() {
			return sc
		}
	}

	return core.EmptySpanContext()
}

type CorrelationsPropagator interface {
	Inject(dctx.Correlations, Supplier)
	Extract(Supplier) dctx.Correlations
}

type chainCorrelationsPropagator struct {
	propagators []CorrelationsPropagator
}

var _ CorrelationsPropagator = chainCorrelationsPropagator{}

func NewChainCorrelationsPropagator(propagators ...CorrelationsPropagator) CorrelationsPropagator {
	return chainCorrelationsPropagator{
		propagators: propagators,
	}
}

func (c chainCorrelationsPropagator) Inject(correlations dctx.Correlations, supplier Supplier) {
	for _, p := range c.propagators {
		p.Inject(correlations, supplier)
	}
}

func (c chainCorrelationsPropagator) Extract(supplier Supplier) dctx.Correlations {
	for _, p := range c.propagators {
		if correlations := p.Extract(supplier); correlations.Len() > 0 {
			return correlations
		}
	}

	return dctx.NewEmptyCorrelations()
}

type BaggagePropagator interface {
	Inject(dctx.Baggage, Supplier)
	Extract(Supplier) dctx.Baggage
}

type chainBaggagePropagator struct {
	propagators []BaggagePropagator
}

var _ BaggagePropagator = chainBaggagePropagator{}

func NewChainBaggagePropagator(propagators ...BaggagePropagator) BaggagePropagator {
	return chainBaggagePropagator{
		propagators: propagators,
	}
}

func (c chainBaggagePropagator) Inject(baggage dctx.Baggage, supplier Supplier) {
	for _, p := range c.propagators {
		p.Inject(baggage, supplier)
	}
}

func (c chainBaggagePropagator) Extract(supplier Supplier) dctx.Baggage {
	for _, p := range c.propagators {
		if baggage := p.Extract(supplier); baggage.Len() > 0 {
			return baggage
		}
	}

	return dctx.NewEmptyBaggage()
}

type Propagator interface {
	Inject(core.SpanContext, dctx.Correlations, dctx.Baggage, Supplier)
	Extract(Supplier) (core.SpanContext, dctx.Correlations, dctx.Baggage)
	SpanContextPropagator() SpanContextPropagator
	CorrelationsPropagator() CorrelationsPropagator
	BaggagePropagator() BaggagePropagator
}

type Propagators struct {
	Scp SpanContextPropagator
	Cp  CorrelationsPropagator
	Bp  BaggagePropagator
}

func (p Propagators) SpanContextPropagator() SpanContextPropagator {
	return p.Scp
}

func (p Propagators) CorrelationsPropagator() CorrelationsPropagator {
	return p.Cp
}

func (p Propagators) BaggagePropagator() BaggagePropagator {
	return p.Bp
}
