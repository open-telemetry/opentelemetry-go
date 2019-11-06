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

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

const (
	Vendor = "ot"
)

var (
	HostKey = otel.Key("http.host")
	URLKey  = otel.Key("http.url")

	propagator = propagation.HTTPTraceContextPropagator{}
)

// Returns the Attributes, Context Entries, and SpanContext that were encoded by Inject.
func Extract(ctx context.Context, req *http.Request) ([]otel.KeyValue, []otel.KeyValue, otel.SpanContext) {
	sc, correlationCtx := propagator.Extract(ctx, req.Header)

	attrs := []otel.KeyValue{
		URLKey.String(req.URL.String()),
		// Etc.
	}

	var correlationCtxKVs []otel.KeyValue
	correlationCtx.Foreach(func(kv otel.KeyValue) bool {
		correlationCtxKVs = append(correlationCtxKVs, kv)
		return true
	})

	return attrs, correlationCtxKVs, sc
}

func Inject(ctx context.Context, req *http.Request) {
	propagator.Inject(ctx, req.Header)
}
