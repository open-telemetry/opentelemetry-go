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

package binary_test

import (
	"go.opentelemetry.io/otel/api/core"
	dctx "go.opentelemetry.io/otel/api/distributedcontext"
	"go.opentelemetry.io/otel/api/propagation/binary"
)

type (
	testValuePropagator struct {
		binary.Propagators
	}

	testPointerPropagator struct {
		binary.Propagators
	}
)

var (
	_ binary.Propagator = testValuePropagator{}
	_ binary.Propagator = &testPointerPropagator{}
)

func (testValuePropagator) ToBytes(core.SpanContext, dctx.Correlations, dctx.Baggage) []byte {
	return nil
}

func (testValuePropagator) FromBytes([]byte) (core.SpanContext, dctx.Correlations, dctx.Baggage) {
	return core.EmptySpanContext(), dctx.NewCorrelations(), dctx.NewBaggage()
}

func (*testPointerPropagator) ToBytes(core.SpanContext, dctx.Correlations, dctx.Baggage) []byte {
	return nil
}

func (*testPointerPropagator) FromBytes([]byte) (core.SpanContext, dctx.Correlations, dctx.Baggage) {
	return core.EmptySpanContext(), dctx.NewCorrelations(), dctx.NewBaggage()
}
