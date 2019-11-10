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

package http_test

import (
	"go.opentelemetry.io/otel/api/core"
	dctx "go.opentelemetry.io/otel/api/distributedcontext"
	"go.opentelemetry.io/otel/api/propagation/http"
)

type (
	// This is to test that embedding binary.Propagators helps to
	// implement the Propagator interface with the implementation
	// having a value receiver.
	testValuePropagator struct {
		http.Propagators
	}

	// This is to test that embedding binary.Propagators helps to
	// implement the Propagator interface with the implementation
	// having a pointer receiver.
	testPointerPropagator struct {
		http.Propagators
	}
)

var (
	_ http.Propagator = testValuePropagator{}
	_ http.Propagator = &testPointerPropagator{}
)

func (testValuePropagator) Inject(core.SpanContext, dctx.Correlations, dctx.Baggage, http.Supplier) {
}

func (testValuePropagator) Extract(http.Supplier) (core.SpanContext, dctx.Correlations, dctx.Baggage) {
	return core.EmptySpanContext(), dctx.NewCorrelations(), dctx.NewBaggage()
}

func (*testPointerPropagator) Inject(core.SpanContext, dctx.Correlations, dctx.Baggage, http.Supplier) {
}

func (*testPointerPropagator) Extract(http.Supplier) (core.SpanContext, dctx.Correlations, dctx.Baggage) {
	return core.EmptySpanContext(), dctx.NewCorrelations(), dctx.NewBaggage()
}
