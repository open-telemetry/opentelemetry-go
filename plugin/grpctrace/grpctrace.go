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

package grpctrace

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/propagators"
)

const (
	// Vendor is the integration provider
	Vendor = "ot"
)

var (
	propagator = propagators.TraceContext{}
)

type metadataSupplier struct {
	metadata *metadata.MD
}

func (s *metadataSupplier) Get(key string) string {
	values := s.metadata.Get(key)
	return strings.Join(values, ",")
}

func (s *metadataSupplier) Set(key string, value string) {
	s.metadata.Append(key, value)
}

// Inject injects the gRPC call metadata into the Span
func Inject(ctx context.Context, metadata *metadata.MD) {
	propagator.Inject(ctx, &metadataSupplier{
		metadata: metadata,
	})
}

// Extract returns the Context Entries and SpanContext that were encoded by Inject.
func Extract(ctx context.Context, metadata *metadata.MD) ([]core.KeyValue, core.SpanContext) {
	spanContext, correlationCtx := propagator.Extract(ctx, &metadataSupplier{
		metadata: metadata,
	})

	var correlationCtxKVs []core.KeyValue
	correlationCtx.Foreach(func(kv core.KeyValue) bool {
		correlationCtxKVs = append(correlationCtxKVs, kv)
		return true
	})

	return correlationCtxKVs, spanContext
}
