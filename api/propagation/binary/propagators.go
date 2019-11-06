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

// SpanContextPropagator is an interface that specifies methods to
// convert SpanContext to/from byte array.
type SpanContextPropagator interface {
	// ToBytes serializes span context into a byte array and
	// returns the array.
	ToBytes(core.SpanContext) []byte

	// FromBytes de-serializes byte array into span context and
	// returns the span context.
	FromBytes([]byte) core.SpanContext
}

type CorrelationsPropagator interface {
	ToBytes(dctx.Correlations) []byte
	FromBytes([]byte) dctx.Correlations
}

type BaggagePropagator interface {
	ToBytes(dctx.Baggage) []byte
	FromBytes([]byte) dctx.Baggage
}

type Propagator interface {
	ToBytes(core.SpanContext, dctx.Correlations, dctx.Baggage) []byte
	FromBytes([]byte) (core.SpanContext, dctx.Correlations, dctx.Baggage)
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
