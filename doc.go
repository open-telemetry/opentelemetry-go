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

/*
Package otel provides an implementation of the OpenTelemetry API.

The provided API is used to instrument code and measure data about that code's
performance and operation. The measured data, by default, is not processed or
transmitted anywhere. An implementation of the OpenTelemetry SDK, like the
default SDK implementation (go.opentelemetry.io/otel/sdk), and associated
exporters are used to process and transport this data.

Tracing

To participate in distributed traces a Span needs to be created for the
operation being performed as part of a traced workflow. It its simplest form:

	var tracer otel.Tracer

	func init() {
		tracer := global.Tracer("instrumentation/package/name")
	}

	func operation(ctx context.Context) {
		var span trace.Span
		ctx, span = tracer.Start(ctx, "operation")
		defer span.End()
		// ...
	}

A Tracer is unique to the instrumentation and is used to create Spans.
Instrumentation should be designed to accept a TracerProvider from which it
can create its own unique Tracer. Alternatively, the registered global
TracerProvider from the go.opentelemetry.io/otel/api/global package can be
used as a default.

	const (
		name    = "instrumentation/package/name"
		version = "0.1.0"
	)

	type Instrumentation struct {
		tracer otel.Tracer
	}

	func NewInstrumentation(tp otel.TracerProvider) *Instrumentation {
		if tp == nil {
			tp := global.TracerProvider()
		}
		return &Instrumentation{
			tracer: tp.Tracer(name, otel.WithTracerVersion(version)),
		}
	}

	func operation(ctx context.Context, inst *Instrumentation) {
		var span trace.Span
		ctx, span = inst.tracer.Start(ctx, "operation")
		defer span.End()
		// ...
	}
*/
package otel // import "go.opentelemetry.io/otel"
