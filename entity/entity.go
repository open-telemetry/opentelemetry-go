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

package entity // import "go.opentelemetry.io/otel/entity"

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/entity/embedded"
)

type Entity interface {
	embedded.Entity

	// IsRecording returns the recording state of the Span. It will return
	// true if the Span is active and events can be recorded.
	IsRecording() bool

	// SetAttributes sets kv as attributes of the Span. If a key from kv
	// already exists for an attribute of the Span it will be overwritten with
	// the value contained in kv.
	SetAttributes(kv ...attribute.KeyValue)

	// TracerProvider returns a TracerProvider that can be used to generate
	// additional Spans on the same telemetry pipeline as the current Span.
	EntityEmitterProvider() EntityEmitterProvider
}

// EntityEmitter is the creator of Spans.
//
// Warning: Methods may be added to this interface in minor releases. See
// package documentation on API implementation for information on how to set
// default behavior for unimplemented methods.
type EntityEmitter interface {
	// Users of the interface can ignore this. This embedded type is only used
	// by implementations of this interface. See the "API Implementations"
	// section of the package documentation for more information.
	embedded.EntityEmitter

	// Start creates a span and a context.Context containing the newly-created span.
	//
	// If the context.Context provided in `ctx` contains a Span then the newly-created
	// Span will be a child of that span, otherwise it will be a root span. This behavior
	// can be overridden by providing `WithNewRoot()` as a SpanOption, causing the
	// newly-created Span to be a root span even if `ctx` contains a Span.
	//
	// When creating a Span it is recommended to provide all known span attributes using
	// the `WithAttributes()` SpanOption as samplers will only have access to the
	// attributes provided when a Span is created.
	//
	// Any Span that is created MUST also be ended. This is the responsibility of the user.
	// Implementations of this API may leak memory or other resources if Spans are not ended.
	//Start(ctx context.Context, spanName string, opts ...SpanStartOption) (context.Context, Span)
}

// EntityEmitterProvider provides EntityEmitters that are used by instrumentation code to
// entity computational workflows.
//
// A EntityEmitterProvider is the collection destination of all Spans from EntityEmitters it
// provides, it represents a unique telemetry collection pipeline. How that
// pipeline is defined, meaning how those Spans are collected, processed, and
// where they are exported, depends on its implementation. Instrumentation
// authors do not need to define this implementation, rather just use the
// provided EntityEmitters to instrument code.
//
// Commonly, instrumentation code will accept a EntityEmitterProvider implementation
// at runtime from its users or it can simply use the globally registered one
// (see https://pkg.go.dev/go.opentelemetry.io/otel#GetEntityEmitterProvider).
//
// Warning: Methods may be added to this interface in minor releases. See
// package documentation on API implementation for information on how to set
// default behavior for unimplemented methods.
type EntityEmitterProvider interface {
	// Users of the interface can ignore this. This embedded type is only used
	// by implementations of this interface. See the "API Implementations"
	// section of the package documentation for more information.
	embedded.EntityEmitterProvider

	// EntityEmitter returns a unique EntityEmitter scoped to be used by instrumentation code
	// to entity computational workflows. The scope and identity of that
	// instrumentation code is uniquely defined by the name and options passed.
	//
	// The passed name needs to uniquely identify instrumentation code.
	// Therefore, it is recommended that name is the Go package name of the
	// library providing instrumentation (note: not the code being
	// instrumented). Instrumentation libraries can have multiple versions,
	// therefore, the WithInstrumentationVersion option should be used to
	// distinguish these different codebases. Additionally, instrumentation
	// libraries may sometimes use entitys to communicate different domains of
	// workflow data (i.e. using spans to communicate workflow events only). If
	// this is the case, the WithScopeAttributes option should be used to
	// uniquely identify EntityEmitters that handle the different domains of workflow
	// data.
	//
	// If the same name and options are passed multiple times, the same EntityEmitter
	// will be returned (it is up to the implementation if this will be the
	// same underlying instance of that EntityEmitter or not). It is not necessary to
	// call this multiple times with the same name and options to get an
	// up-to-date EntityEmitter. All implementations will ensure any EntityEmitterProvider
	// configuration changes are propagated to all provided EntityEmitters.
	//
	// If name is empty, then an implementation defined default name will be
	// used instead.
	//
	// This method is safe to call concurrently.
	EntityEmitter(name string, options ...EntityEmitterOption) EntityEmitter
}
