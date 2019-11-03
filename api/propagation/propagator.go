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

package propagation

import (
	"context"

	"go.opentelemetry.io/otel/api/core"
	dctx "go.opentelemetry.io/otel/api/distributedcontext"
)

// TextFormatPropagator is an interface that specifies methods to inject and extract SpanContext
// and distributed context into/from a carrier using Supplier interface.
// For example, HTTP Trace Context propagator would encode SpanContext into W3C Trace
// Context Header and set the header into HttpRequest.
type TextFormatPropagator interface {
	// Inject method retrieves current SpanContext from the ctx, encodes it into propagator
	// specific format and then injects the encoded SpanContext using supplier into a carrier
	// associated with the supplier. It also takes a correlationCtx whose values will be
	// injected into a carrier using the supplier.
	Inject(ctx context.Context, supplier Supplier)

	// Extract method retrieves encoded SpanContext using supplier from the associated carrier.
	// It decodes the SpanContext and returns it and a dctx of correlated context.
	// If no SpanContext was retrieved OR if the retrieved SpanContext is invalid then
	// an empty SpanContext is returned.
	Extract(ctx context.Context, supplier Supplier) (core.SpanContext, dctx.Map)

	// GetAllKeys returns all the keys that this propagator injects/extracts into/from a
	// carrier. The use cases for this are
	// * allow pre-allocation of fields, especially in systems like gRPC Metadata
	// * allow a single-pass over an iterator (ex OpenTracing has no getter in TextMap)
	GetAllKeys() []string
}

// Supplier is an interface that specifies methods to retrieve and store
// value for a key to an associated carrier.
// Get method retrieves the value for a given key.
// Set method stores the value for a given key.
type Supplier interface {
	Get(key string) string
	Set(key string, value string)
}

// BinaryFormatPropagator is an interface that specifies methods to convert SpanContext
// to/from byte array.
type BinaryFormatPropagator interface {
	// ToBytes serializes span context into a byte array and returns the array.
	ToBytes(sc core.SpanContext) []byte

	// FromBytes de-serializes byte array into span context and returns the span context.
	FromBytes([]byte) core.SpanContext
}
