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

package httptrace

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/api/core"
	dctx "go.opentelemetry.io/otel/api/distributedcontext"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/trace"
	propagation "go.opentelemetry.io/otel/propagation/http"
)

const (
	Vendor = "ot"
)

var (
	HostKey = key.New("http.host")
	URLKey  = key.New("http.url")

	scPropagator = propagation.TraceContextPropagator{}
	cPropagator  = propagation.CorrelationContextPropagator{}
)

// Returns the Attributes, Context Entries, and SpanContext that were encoded by Inject.
func Extract(ctx context.Context, req *http.Request) ([]core.KeyValue, []dctx.Correlation, core.SpanContext) {
	sc := scPropagator.Extract(req.Header)
	correlationCtx := cPropagator.Extract(req.Header)

	attrs := []core.KeyValue{
		URLKey.String(req.URL.String()),
		// Etc.
	}

	return attrs, correlationCtx.Correlations(), sc
}

func Inject(ctx context.Context, req *http.Request) {
	scPropagator.Inject(trace.CurrentSpan(ctx).SpanContext(), req.Header)
	cPropagator.Inject(dctx.CorrelationsFromContext(ctx), req.Header)
}
