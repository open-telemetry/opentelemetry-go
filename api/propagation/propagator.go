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

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/tag"
)

// TextFormatPropagator is an interface that specifies methods to create objects
// that encodes/decodes into/from text format representation of SpanContext and tag.Map.
//
// CarrierInjector method creates an Injector object and binds the carrier to the object.
// Injector object provides Inject method to inject SpanContext and tag.Map after serializing into
// a text format associated with the propagator.
//
// Similarly, CarrierExtractor method creates an Extractor object and binds the carrier to the
// object. Extractor object provides Extract method to extract text formatted de-serialized
// SpanContext and tag.Map
//
// Typically, a plugin for transport like HTTP uses this interface to allow user
// to configure appropriate text format propagators.
type TextFormatPropagator interface {
	// CarrierExtractor creates an object that implements Extractor interface.
	// It binds provided carrier to the object.
	CarrierExtractor(carrier interface{}) Extractor

	// CarrierInjector creates an object that implements Injector interface.
	// It binds provided carrier to the object.
	CarrierInjector(carrier interface{}) Injector
}

type Injector interface {
	// Inject encodes span context and tag.Map and inserts them in to
	// carrier associated with the injector. For example in case of http request,
	// span context could be added to the request (carrier) as W3C Trace context header.
	Inject(core.SpanContext, tag.Map)
}

type Extractor interface {
	// Extract decodes span context and tag.Map from a carrier associated with the
	// extractor. For example in case of http request, span context could be extracted
	// from the W3C Trace context header.
	Extract() (core.SpanContext, tag.Map)
}

// Propagator is an interface that specifies methods to inject and extract SpanContext
// into/from a carrier using Supplier interface.
// For example, HTTP Trace Context Propagator would encode SpanContext into W3C Trace
// Context Header and set the header into HttpRequest.
type Propagator interface {
	// Inject method retrieves current SpanContext from the ctx, encodes it into propagator
	// specific format and then injects the encoded SpanContext using supplier into a carrier
	// associated with the supplier.
	Inject(ctx context.Context, supplier Supplier)

	// Extract method retrieves SpanContext using supplier from the associated carrier with
	// the supplier. It decodes the SpanContext and it creates remote span using registered
	// tracer.
	Extract(ctx context.Context, supplier Supplier) context.Context

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
