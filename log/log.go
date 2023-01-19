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

package log // import "go.opentelemetry.io/otel/log"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
)

// LogRecord is the individual component of a trace. It represents a single named
// and timed operation of a workflow that is traced. A Logger is used to
// create a LogRecord and it is then up to the operation the LogRecord represents to
// properly end the LogRecord when the operation itself ends.
//
// Warning: methods may be added to this interface in minor releases.
type LogRecord interface {
	// IsRecording returns the recording state of the LogRecord. It will return
	// true if the LogRecord is active and events can be recorded.
	IsRecording() bool

	// SetAttributes sets kv as attributes of the LogRecord. If a key from kv
	// already exists for an attribute of the LogRecord it will be overwritten with
	// the value contained in kv.
	SetAttributes(kv ...attribute.KeyValue)

	// LoggerProvider returns a LoggerProvider that can be used to generate
	// additional Spans on the same telemetry pipeline as the current LogRecord.
	LoggerProvider() LoggerProvider
}

// Logger is the creator of Spans.
//
// Warning: methods may be added to this interface in minor releases.
type Logger interface {
	// Emit creates a span and a context.Context containing the newly-created span.
	//
	// If the context.Context provided in `ctx` contains a LogRecord then the newly-created
	// LogRecord will be a child of that span, otherwise it will be a root span. This behavior
	// can be overridden by providing `WithNewRoot()` as a SpanOption, causing the
	// newly-created LogRecord to be a root span even if `ctx` contains a LogRecord.
	//
	// When creating a LogRecord it is recommended to provide all known span attributes using
	// the `WithAttributes()` SpanOption as samplers will only have access to the
	// attributes provided when a LogRecord is created.
	//
	// Any LogRecord that is created MUST also be ended. This is the responsibility of the user.
	// Implementations of this API may leak memory or other resources if Spans are not ended.
	Emit(ctx context.Context, opts ...LogRecordOption)
}

// LoggerProvider provides Tracers that are used by instrumentation code to
// trace computational workflows.
//
// A LoggerProvider is the collection destination of all Spans from Tracers it
// provides, it represents a unique telemetry collection pipeline. How that
// pipeline is defined, meaning how those Spans are collected, processed, and
// where they are exported, depends on its implementation. Instrumentation
// authors do not need to define this implementation, rather just use the
// provided Tracers to instrument code.
//
// Commonly, instrumentation code will accept a LoggerProvider implementation
// at runtime from its users or it can simply use the globally registered one
// (see https://pkg.go.dev/go.opentelemetry.io/otel#GetTracerProvider).
//
// Warning: methods may be added to this interface in minor releases.
type LoggerProvider interface {
	// Logger returns a unique Logger scoped to be used by instrumentation code
	// to trace computational workflows. The scope and identity of that
	// instrumentation code is uniquely defined by the name and options passed.
	//
	// The passed name needs to uniquely identify instrumentation code.
	// Therefore, it is recommended that name is the Go package name of the
	// library providing instrumentation (note: not the code being
	// instrumented). Instrumentation libraries can have multiple versions,
	// therefore, the WithInstrumentationVersion option should be used to
	// distinguish these different codebases. Additionally, instrumentation
	// libraries may sometimes use traces to communicate different domains of
	// workflow data (i.e. using spans to communicate workflow events only). If
	// this is the case, the WithScopeAttributes option should be used to
	// uniquely identify Tracers that handle the different domains of workflow
	// data.
	//
	// If the same name and options are passed multiple times, the same Logger
	// will be returned (it is up to the implementation if this will be the
	// same underlying instance of that Logger or not). It is not necessary to
	// call this multiple times with the same name and options to get an
	// up-to-date Logger. All implementations will ensure any LoggerProvider
	// configuration changes are propagated to all provided Tracers.
	//
	// If name is empty, then an implementation defined default name will be
	// used instead.
	//
	// This method is safe to call concurrently.
	Logger(name string, options ...LoggerOption) Logger
}
