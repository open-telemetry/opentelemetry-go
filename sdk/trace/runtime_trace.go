// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace // import "go.opentelemetry.io/otel/sdk/trace"

import (
	"context"

	runtimetrace "runtime/trace"

	"go.opentelemetry.io/otel/trace"
)

type runtimeTraceEndFn func()

// runtimeTracer abstracts the runtime/trace package so it can be mocked
// in tests. The runtime/trace API provides local, runtime-level
// instrumentation, recorded in Go execution profiles, similar to distributed
// tracing but confined to a single Go process. This interface allows
// integrating that instrumentation with distributed tracing without
// duplicating logic.
type runtimeTracer interface {
	IsEnabled() bool
	NewTask(ctx context.Context, name string) (context.Context, runtimeTraceEndFn)
	StartRegion(ctx context.Context, name string) runtimeTraceEndFn
}

// standardRuntimeTracer is the default implementation of runtimeTracer.
// It simply wraps the runtime/trace package.
type standardRuntimeTracer struct{}

func (standardRuntimeTracer) IsEnabled() bool {
	return runtimetrace.IsEnabled()
}

func (standardRuntimeTracer) NewTask(ctx context.Context, name string) (context.Context, runtimeTraceEndFn) {
	nctx, task := runtimetrace.NewTask(ctx, name)
	return nctx, task.End
}

func (standardRuntimeTracer) StartRegion(ctx context.Context, name string) runtimeTraceEndFn {
	region := runtimetrace.StartRegion(ctx, name)
	return region.End
}

// globalRuntimeTracer is the variable that holds the global runtimeTracer
// implementation. It defaults to the real implementation but can be swapped
// for testing.
var globalRuntimeTracer runtimeTracer = standardRuntimeTracer{}

// profilingSpan is an interface for spans that can integrate with
// runtimeTracer.
type profilingSpan interface {
	// startProfiling may start a "runtime/trace" Task (returning a new
	// context) or a Region (no context change). If tracing is disabled
	// (globally or for the span), it does nothing. Concrete implementations
	// may have their own defaults when config is not explicit.
	startProfiling(ctx context.Context, config *trace.SpanConfig, tracerSetting trace.ProfilingMode) context.Context
	endProfiling()
	profilingStarted() bool
}
