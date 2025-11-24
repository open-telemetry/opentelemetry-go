// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace // import "go.opentelemetry.io/otel/sdk/trace"

import (
	"context"

	profile "runtime/trace"

	"go.opentelemetry.io/otel/trace"
)

type endFunc func()

// runtimeTraceAPI abstracts the runtime/trace package so it can be mocked
// in tests. The runtime/trace API provides local, runtime-level
// instrumentation, recorded in Go execution profiles, similar to distributed
// tracing but confined to a single Go process. This interface allows
// integrating that instrumentation with distributed tracing without
// duplicating logic.
type runtimeTraceAPI interface {
	IsEnabled() bool
	NewTask(ctx context.Context, name string) (context.Context, endFunc)
	StartRegion(ctx context.Context, name string) endFunc
}

type standardRuntimeTraceWrapper struct{}

func (r standardRuntimeTraceWrapper) IsEnabled() bool {
	return profile.IsEnabled()
}

func (r standardRuntimeTraceWrapper) NewTask(ctx context.Context, name string) (context.Context, endFunc) {
	nctx, task := profile.NewTask(ctx, name)
	return nctx, task.End
}

func (r standardRuntimeTraceWrapper) StartRegion(ctx context.Context, name string) endFunc {
	region := profile.StartRegion(ctx, name)
	return region.End
}

// globalRuntimeTracer is the variable that holds the global runtime tracer implementation.
// It defaults to the real implementation but can be swapped for testing.
var globalRuntimeTracer runtimeTraceAPI = standardRuntimeTraceWrapper{}

// profilingSpan is an interface for spans that can integrate with runtimeTraceAPI
type profilingSpan interface {
	// startProfiling may start a "runtime/trace" Task (returning a new context)
	// or a Region (no context change). If tracing is disabled (globally or
	// for the span), it does nothing.
	startProfiling(ctx context.Context, config *trace.SpanConfig) context.Context
	endProfiling()
	profilingStarted() bool
}
