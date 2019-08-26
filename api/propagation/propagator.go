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

// Package tracecontext contains HTTP propagator for TraceContext standard.
// See https://github.com/w3c/distributed-tracing for more information.
package propagation // import "go.opentelemetry.io/api/propagation"

import (
	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/tag"
)

// TextFormatPropagator is an interface that specifies methods to create CarrierInjector
// and CarrierExtractor objects. These methods bind given carrier to the created object.
// CarrierInjector object implements Inject method to inject
// SpanContext and tag.Map as a text format into carrier like HTTP request.
// Similarly, CarrierExtractor object implements Extract method to extract SpanContext
// encoded in text format from a carrier like HTTP request.
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
	// Inject serializes span context and tag.Map and inserts them in to
	// carrier associated with the injector. For example in case of http request,
	// span context could be added to the request (carrier) as W3C Trace context header.
	Inject(core.SpanContext, tag.Map)
}

type Extractor interface {
	// Extract de-serializes span context and tag.Map from a carrier associated with the
	// extractor. For example in case of http request, span context could be extracted
	// from the W3C Trace context header.
	Extract() (core.SpanContext, tag.Map)
}
