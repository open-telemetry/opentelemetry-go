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

// Package opencensus provides a migration bridge from OpenCensus to
// OpenTelemetry. The NewTracer function should be used to create an
// OpenCensus Tracer from an OpenTelemetry Tracer. This Tracer can be use in
// place of any existing OpenCensus Tracer and will generate OpenTelemetry
// spans for traces. These spans will be exported by the OpenTelemetry
// TracerProvider the original OpenTelemetry Tracer came from.
//
// There are known limitations to this bridge:
//
// - The AddLink method for OpenCensus Spans is not compatible with the
// OpenTelemetry Span. No link can be added to an OpenTelemetry Span once it
// is started. Any calls to this method for the OpenCensus Span will result
// in an error being sent to the OpenTelemetry default ErrorHandler.
//
// - The NewContext method of the OpenCensus Tracer cannot embed an OpenCensus
// Span in a context unless that Span was created by that Tracer.
//
// - Conversion of custom OpenCensus Samplers to OpenTelemetry is not
// implemented. An error will be sent to the OpenTelemetry default
// ErrorHandler if this is attempted.
package opencensus // import "go.opentelemetry.io/otel/bridge/opencensus"
