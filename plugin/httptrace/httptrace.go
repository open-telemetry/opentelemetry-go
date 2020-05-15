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

package httptrace

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/api/correlation"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/propagation"
	"go.opentelemetry.io/otel/api/trace"
)

var (
	HostKey = kv.Key("http.host")
	URLKey  = kv.Key("http.url")
)

// Returns the Attributes, Context Entries, and SpanContext that were encoded by Inject.
func Extract(ctx context.Context, req *http.Request) ([]kv.KeyValue, []kv.KeyValue, trace.SpanContext) {
	ctx = propagation.ExtractHTTP(ctx, global.Propagators(), req.Header)

	attrs := []kv.KeyValue{
		URLKey.String(req.URL.String()),
		// Etc.
	}

	var correlationCtxKVs []kv.KeyValue
	correlation.MapFromContext(ctx).Foreach(func(kv kv.KeyValue) bool {
		correlationCtxKVs = append(correlationCtxKVs, kv)
		return true
	})

	return attrs, correlationCtxKVs, trace.RemoteSpanContextFromContext(ctx)
}

func Inject(ctx context.Context, req *http.Request) {
	propagation.InjectHTTP(ctx, global.Propagators(), req.Header)
}
