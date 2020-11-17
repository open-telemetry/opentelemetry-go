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

package binary // import "go.opentelemetry.io/otel/bridge/opencensus/binary"

import (
	"context"

	ocpropagation "go.opencensus.io/trace/propagation"

	"go.opentelemetry.io/otel/bridge/opencensus/utils"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type key uint

const binaryKey key = 0

// binaryHeader is the same as traceContextKey is in opencensus:
// https://github.com/census-instrumentation/opencensus-go/blob/3fb168f674736c026e623310bfccb0691e6dec8a/plugin/ocgrpc/trace_common.go#L30
const binaryHeader = "grpc-trace-bin"

// Binary is an OpenTelemetry implementation of the OpenCensus grpc binary format.
// Binary propagation was temporarily removed from opentelemetry.  See
// https://github.com/open-telemetry/opentelemetry-specification/issues/437
type Binary struct{}

var _ propagation.TextMapPropagator = Binary{}

// Inject injects context into the TextMapCarrier
func (b Binary) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	binaryContext := ctx.Value(binaryKey)
	if state, ok := binaryContext.(string); binaryContext != nil && ok {
		carrier.Set(binaryHeader, state)
	}

	sc := trace.SpanContextFromContext(ctx)
	if !sc.IsValid() {
		return
	}
	h := ocpropagation.Binary(utils.OTelSpanContextToOC(sc))
	carrier.Set(binaryHeader, string(h))
}

// Extract extracts the SpanContext from the TextMapCarrier
func (b Binary) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	state := carrier.Get(binaryHeader)
	if state != "" {
		ctx = context.WithValue(ctx, binaryKey, state)
	}

	sc := b.extract(carrier)
	if !sc.IsValid() {
		return ctx
	}
	return trace.ContextWithRemoteSpanContext(ctx, sc)
}

func (b Binary) extract(carrier propagation.TextMapCarrier) trace.SpanContext {
	h := carrier.Get(binaryHeader)
	if h == "" {
		return trace.SpanContext{}
	}
	ocContext, ok := ocpropagation.FromBinary([]byte(h))
	if !ok {
		return trace.SpanContext{}
	}
	return utils.OCSpanContextToOTel(ocContext)
}

// Fields returns the fields that this propagator modifies.
func (b Binary) Fields() []string {
	return []string{binaryHeader}
}
