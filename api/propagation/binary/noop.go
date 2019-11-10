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

package binary

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

func (NoopSpanContextPropagator) ToBytes(core.SpanContext) []byte {
	return nil
}

func (NoopSpanContextPropagator) FromBytes([]byte) core.SpanContext {
	return core.EmptySpanContext()
}

func (NoopCorrelationsPropagator) ToBytes(dctx.Correlations) []byte {
	return nil
}

func (NoopCorrelationsPropagator) FromBytes([]byte) dctx.Correlations {
	return dctx.NewCorrelations()
}

func (NoopBaggagePropagator) ToBytes(dctx.Baggage) []byte {
	return nil
}

func (NoopBaggagePropagator) FromBytes([]byte) dctx.Baggage {
	return dctx.NewBaggage()
}

func (NoopPropagator) ToBytes(core.SpanContext, dctx.Correlations, dctx.Baggage) []byte {
	return nil
}

func (NoopPropagator) FromBytes([]byte) (core.SpanContext, dctx.Correlations, dctx.Baggage) {
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
