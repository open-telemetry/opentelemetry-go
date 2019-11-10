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

type (
	NoopSpanContextPropagator  struct{}
	NoopCorrelationsPropagator struct{}
	NoopBaggagePropagator      struct{}
	NoopPropagator             struct{}
)

var (
	_ SpanContextPropagator  = NoopSpanContextPropagator{}
	_ CorrelationsPropagator = NoopCorrelationsPropagator{}
	_ BaggagePropagator      = NoopBaggagePropagator{}
	_ Propagator             = NoopPropagator{}
)

func (NoopSpanContextPropagator) Inject(core.SpanContext, Supplier) {
}

func (NoopSpanContextPropagator) Extract(Supplier) core.SpanContext {
	return core.EmptySpanContext()
}

func (NoopCorrelationsPropagator) Inject(dctx.Correlations, Supplier) {
}

func (NoopCorrelationsPropagator) Extract(Supplier) dctx.Correlations {
	return dctx.NewCorrelations()
}

func (NoopBaggagePropagator) Inject(dctx.Baggage, Supplier) {
}

func (NoopBaggagePropagator) Extract(Supplier) dctx.Baggage {
	return dctx.NewBaggage()
}

func (NoopPropagator) Inject(core.SpanContext, dctx.Correlations, dctx.Baggage, Supplier) {
}

func (NoopPropagator) Extract(Supplier) (core.SpanContext, dctx.Correlations, dctx.Baggage) {
	return core.EmptySpanContext(), dctx.NewCorrelations(), dctx.NewBaggage()
}

func (NoopPropagator) SpanContextPropagator() SpanContextPropagator {
	return NoopSpanContextPropagator{}
}

func (NoopPropagator) CorrelationsPropagator() CorrelationsPropagator {
	return NoopCorrelationsPropagator{}
}

func (NoopPropagator) BaggagePropagator() BaggagePropagator {
	return NoopBaggagePropagator{}
}
