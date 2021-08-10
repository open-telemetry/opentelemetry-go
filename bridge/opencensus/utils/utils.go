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

// Package utils provides utilities for the OpenCensus bridge.
//
// Deprecated: Use the equivalent functions from the bridge/opencensus package
// instead.
package utils // import "go.opentelemetry.io/otel/bridge/opencensus/utils"

import (
	octrace "go.opencensus.io/trace"

	"go.opentelemetry.io/otel/bridge/opencensus/internal/oc2otel"
	"go.opentelemetry.io/otel/bridge/opencensus/internal/otel2oc"
	"go.opentelemetry.io/otel/trace"
)

// OTelSpanContextToOC converts from an OpenTelemetry SpanContext to an
// OpenCensus SpanContext, and handles any incompatibilities with the global
// error handler.
//
// Deprecated: Use OTelSpanContextToOC from bridge/opencensus instead.
func OTelSpanContextToOC(sc trace.SpanContext) octrace.SpanContext {
	return otel2oc.SpanContext(sc)
}

// OCSpanContextToOTel converts from an OpenCensus SpanContext to an
// OpenTelemetry SpanContext.
//
// Deprecated: Use OCSpanContextToOTel from bridge/opencensus instead.
func OCSpanContextToOTel(sc octrace.SpanContext) trace.SpanContext {
	return oc2otel.SpanContext(sc)
}
