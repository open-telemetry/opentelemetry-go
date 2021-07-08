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

package oteltest // import "go.opentelemetry.io/otel/oteltest"

import (
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// NewTracerProvider returns a *TracerProvider configured with options.
func NewTracerProvider(options ...Option) trace.TracerProvider {
	cfg := newConfig(options)
	return sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(cfg.SpanRecorder),
		sdktrace.WithSampler(sdktrace.AlwaysSample()))
}

// DefaultTracer returns a default tracer for testing purposes.
func DefaultTracer() trace.Tracer {
	return NewTracerProvider().Tracer("")
}
