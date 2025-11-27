// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/trace"
)

// mockRuntimeTracer is a simple mock implementation of runtimeTracer for testing.
type mockRuntimeTracer struct {
	isEnabled        bool
	isEnabledCalls   atomic.Uint64
	newTaskCalls     atomic.Uint64
	startRegionCalls atomic.Uint64
	taskEndCalls     atomic.Uint64
	regionEndCalls   atomic.Uint64
}

func newMockRuntimeTracer(enabled bool) *mockRuntimeTracer {
	return &mockRuntimeTracer{
		isEnabled: enabled,
	}
}

func (m *mockRuntimeTracer) IsEnabled() bool {
	m.isEnabledCalls.Add(1)
	return m.isEnabled
}

func (m *mockRuntimeTracer) NewTask(ctx context.Context, _ string) (context.Context, runtimeTraceEndFn) {
	m.newTaskCalls.Add(1)
	endFunc := func() {
		m.taskEndCalls.Add(1)
	}
	return ctx, endFunc
}

func (m *mockRuntimeTracer) StartRegion(_ context.Context, _ string) runtimeTraceEndFn {
	m.startRegionCalls.Add(1)
	endFunc := func() {
		m.regionEndCalls.Add(1)
	}
	return endFunc
}

func assertCalls(t *testing.T, m *mockRuntimeTracer, isEnabled, newTask, startRegion, taskEnd, regionEnd int) {
	assert.Equal(t, uint64(isEnabled), m.isEnabledCalls.Load())
	assert.Equal(t, uint64(newTask), m.newTaskCalls.Load())
	assert.Equal(t, uint64(startRegion), m.startRegionCalls.Load())
	assert.Equal(t, uint64(taskEnd), m.taskEndCalls.Load())
	assert.Equal(t, uint64(regionEnd), m.regionEndCalls.Load())
}

func TestDefaultInstrumentation(t *testing.T) {
	// The cases in this test will be the default behavior for all OTel users that do not specify any runtime/trace
	// instrumentation options. The default behavior is to create a task for each root span and nothing else.

	originalRuntimeTracer := globalRuntimeTracer
	t.Cleanup(func() {
		globalRuntimeTracer = originalRuntimeTracer
	})

	tracerProvider := NewTracerProvider(WithSampler(AlwaysSample()))
	tracer := tracerProvider.Tracer("TestDefaultInstrumentation")

	mockRT := newMockRuntimeTracer(true)
	globalRuntimeTracer = mockRT

	ctx := t.Context()

	t.Run("local root span creates a task", func(t *testing.T) {
		_, span := tracer.Start(ctx, "root-span")
		assertCalls(t, mockRT, 1, 1, 0, 0, 0)
		profilingSpan, ok := span.(profilingSpan)
		require.True(t, ok)
		assert.True(t, profilingSpan.profilingStarted())
		span.End()
		assertCalls(t, mockRT, 1, 1, 0, 1, 0)
	})

	rootCtx, _ := tracer.Start(ctx, "root-span", trace.NoProfiling())
	t.Run("disable default instrumentation individually on root span", func(t *testing.T) {
		assertCalls(t, mockRT, 2, 1, 0, 1, 0) // no new task this time
	})

	t.Run("local child spans do nothing", func(t *testing.T) {
		_, childSpan := tracer.Start(rootCtx, "child-span")
		assertCalls(t, mockRT, 3, 1, 0, 1, 0)
		profilingSpan, ok := childSpan.(profilingSpan)
		require.True(t, ok)
		assert.False(t, profilingSpan.profilingStarted())
		childSpan.End()
		assertCalls(t, mockRT, 3, 1, 0, 1, 0)
	})

	t.Run("force manual task for child span", func(t *testing.T) {
		_, childSpan := tracer.Start(rootCtx, "child-span-task", trace.ProfileTask())
		assertCalls(t, mockRT, 4, 2, 0, 1, 0)
		profilingSpan, ok := childSpan.(profilingSpan)
		require.True(t, ok)
		assert.True(t, profilingSpan.profilingStarted())
		childSpan.End()
		assertCalls(t, mockRT, 4, 2, 0, 2, 0)
	})

	t.Run("force manual region for child span", func(t *testing.T) {
		_, childSpan := tracer.Start(rootCtx, "child-span-region", trace.ProfileRegion())
		assertCalls(t, mockRT, 5, 2, 1, 2, 0)
		profilingSpan, ok := childSpan.(profilingSpan)
		require.True(t, ok)
		assert.True(t, profilingSpan.profilingStarted())
		childSpan.End()
		assertCalls(t, mockRT, 5, 2, 1, 2, 1)
	})

	t.Run("force auto instrumentation for child span", func(t *testing.T) {
		// Not sure how useful this is for users, but trace.ProfilingAuto > trace.ProfilingDefault, so it works.
		_, childSpan := tracer.Start(
			rootCtx,
			"child-span-auto",
			trace.WithProfileTask(trace.ProfilingAuto),
			trace.WithProfileRegion(trace.ProfilingAuto),
		)
		assertCalls(t, mockRT, 6, 2, 2, 2, 1)
		profilingSpan, ok := childSpan.(profilingSpan)
		require.True(t, ok)
		assert.True(t, profilingSpan.profilingStarted())
		childSpan.End()
		assertCalls(t, mockRT, 6, 2, 2, 2, 2)
	})
}

func TestManualInstrumentation(t *testing.T) {
	// Notice that the tracer is created with trace.WithProfilingMode(trace.ProfilingManual), and all spans won't create
	// tasks or region unless they are explicitly tagged with trace.ProfileTask() or trace.ProfileRegion().

	originalRuntimeTracer := globalRuntimeTracer
	t.Cleanup(func() {
		globalRuntimeTracer = originalRuntimeTracer
	})

	tracerProvider := NewTracerProvider(WithSampler(AlwaysSample()))
	tracer := tracerProvider.Tracer("TestManualInstrumentation", trace.WithProfilingMode(trace.ProfilingManual))

	mockRT := newMockRuntimeTracer(true)
	globalRuntimeTracer = mockRT

	ctx := t.Context()

	t.Run("root span with task", func(t *testing.T) {
		_, rootSpanWithTask := tracer.Start(ctx, "root-span-with-task", trace.ProfileTask())
		assertCalls(t, mockRT, 1, 1, 0, 0, 0)
		profilingSpan, ok := rootSpanWithTask.(profilingSpan)
		require.True(t, ok)
		assert.True(t, profilingSpan.profilingStarted())
		rootSpanWithTask.End()
		assertCalls(t, mockRT, 1, 1, 0, 1, 0)
	})

	t.Run("root span with region", func(t *testing.T) {
		// This is a special case where the region does not have a task
		// associated but is allowed by the runtime/trace API. The region will
		// be associated with the background task.
		_, rootSpanWithRegion := tracer.Start(ctx, "root-span-with-region", trace.ProfileRegion())
		assertCalls(t, mockRT, 2, 1, 1, 1, 0)
		profilingSpan, ok := rootSpanWithRegion.(profilingSpan)
		require.True(t, ok)
		assert.True(t, profilingSpan.profilingStarted())
		rootSpanWithRegion.End()
		assertCalls(t, mockRT, 2, 1, 1, 1, 1)
	})

	t.Run("root span without profiling", func(t *testing.T) {
		_, rootSpanWithoutProfiling := tracer.Start(ctx, "root-span-without-profiling")
		assertCalls(t, mockRT, 3, 1, 1, 1, 1)
		profilingSpan, ok := rootSpanWithoutProfiling.(profilingSpan)
		require.True(t, ok)
		assert.False(t, profilingSpan.profilingStarted())
		rootSpanWithoutProfiling.End()
		assertCalls(t, mockRT, 3, 1, 1, 1, 1)
	})

	t.Run("root span with both task and region", func(t *testing.T) {
		// Not sure of the utility of this case, but our API is flexible/not opinionated.
		_, rootSpanWithBothTaskAndRegion := tracer.Start(
			ctx,
			"root-span-with-both-task-and-region",
			trace.ProfileTask(),
			trace.ProfileRegion(),
		)
		assertCalls(t, mockRT, 4, 2, 2, 1, 1)
		profilingSpan, ok := rootSpanWithBothTaskAndRegion.(profilingSpan)
		require.True(t, ok)
		assert.True(t, profilingSpan.profilingStarted())
		rootSpanWithBothTaskAndRegion.End()
		assertCalls(t, mockRT, 4, 2, 2, 2, 2)
	})

	rootCtx, _ := tracer.Start(ctx, "root-span")
	assertCalls(t, mockRT, 5, 2, 2, 2, 2)

	t.Run("child span with task", func(t *testing.T) {
		_, childSpanWithTask := tracer.Start(rootCtx, "child-span-with-task", trace.ProfileTask())
		assertCalls(t, mockRT, 6, 3, 2, 2, 2)
		profilingSpan, ok := childSpanWithTask.(profilingSpan)
		require.True(t, ok)
		assert.True(t, profilingSpan.profilingStarted())
		childSpanWithTask.End()
		assertCalls(t, mockRT, 6, 3, 2, 3, 2)
	})

	t.Run("child span with region", func(t *testing.T) {
		// Notice the parent span is not attached to a task, so this region is not associated with any task.
		// This is fine for runtime/trace API, the region will be associated with the background task.
		_, childSpanWithRegion := tracer.Start(rootCtx, "child-span-with-region", trace.ProfileRegion())
		assertCalls(t, mockRT, 7, 3, 3, 3, 2)
		profilingSpan, ok := childSpanWithRegion.(profilingSpan)
		require.True(t, ok)
		assert.True(t, profilingSpan.profilingStarted())
		childSpanWithRegion.End()
		assertCalls(t, mockRT, 7, 3, 3, 3, 3)
	})

	t.Run("child span without profiling", func(t *testing.T) {
		_, childSpanWithoutProfiling := tracer.Start(rootCtx, "child-span-without-profiling")
		assertCalls(t, mockRT, 8, 3, 3, 3, 3)
		profilingSpan, ok := childSpanWithoutProfiling.(profilingSpan)
		require.True(t, ok)
		assert.False(t, profilingSpan.profilingStarted())
		childSpanWithoutProfiling.End()
		assertCalls(t, mockRT, 8, 3, 3, 3, 3)
	})
}

func TestAutoInstrumentation(t *testing.T) {
	// Notice the usage of trace.AutoProfiling(), trace.AsyncEnd() and trace.NoProfiling() in this test.

	originalRuntimeTracer := globalRuntimeTracer
	t.Cleanup(func() {
		globalRuntimeTracer = originalRuntimeTracer
	})

	tracerProvider := NewTracerProvider(WithSampler(AlwaysSample()))
	tracer := tracerProvider.Tracer("TestAutoProfiling", trace.AutoProfiling())

	t.Run("local root creates a task and children create regions", func(t *testing.T) {
		mockRT := newMockRuntimeTracer(true)
		globalRuntimeTracer = mockRT

		ctx := t.Context()
		spanCtx, span := tracer.Start(ctx, "root-span")
		assertCalls(t, mockRT, 1, 1, 0, 0, 0)

		childSpanCtx, childSpan := tracer.Start(spanCtx, "child-span")
		assertCalls(t, mockRT, 2, 1, 1, 0, 0)

		_, grandchildSpan := tracer.Start(childSpanCtx, "grandchild-span")
		assertCalls(t, mockRT, 3, 1, 2, 0, 0)

		grandchildSpan.End()
		assertCalls(t, mockRT, 3, 1, 2, 0, 1)

		childSpan.End()
		assertCalls(t, mockRT, 3, 1, 2, 0, 2)

		span.End()
		assertCalls(t, mockRT, 3, 1, 2, 1, 2)
	})

	t.Run("async spans create tasks instead of regions", func(t *testing.T) {
		mockRT := newMockRuntimeTracer(true)
		globalRuntimeTracer = mockRT

		ctx := t.Context()
		rootCtx, rootSpan := tracer.Start(ctx, "root-span")
		assertCalls(t, mockRT, 1, 1, 0, 0, 0)

		_, childSpan := tracer.Start(rootCtx, "child-span", trace.AsyncEnd())
		assertCalls(t, mockRT, 2, 2, 0, 0, 0) // async span creates a task instead of a region

		childSpan.End()
		assertCalls(t, mockRT, 2, 2, 0, 1, 0)

		rootSpan.End()
		assertCalls(t, mockRT, 2, 2, 0, 2, 0)
	})

	t.Run("force manual instrumentation on individual span", func(t *testing.T) {
		mockRT := newMockRuntimeTracer(true)
		globalRuntimeTracer = mockRT

		ctx := t.Context()
		rootCtx, rootSpan := tracer.Start(ctx, "root-span")
		assertCalls(t, mockRT, 1, 1, 0, 0, 0)

		_, childSpan := tracer.Start(rootCtx, "child-span", trace.ProfileTask())
		// notice it also creates a region because it is part of the auto instrumentation and wasn't disabled explicitly
		assertCalls(t, mockRT, 2, 2, 1, 0, 0)

		childSpan.End()
		assertCalls(t, mockRT, 2, 2, 1, 1, 1)

		_, childSpan2 := tracer.Start(rootCtx, "child-span-2", trace.NoProfiling(), trace.ProfileTask())
		// no region this time
		assertCalls(t, mockRT, 3, 3, 1, 1, 1)

		childSpan2.End()
		assertCalls(t, mockRT, 3, 3, 1, 2, 1)

		rootSpan.End()
		assertCalls(t, mockRT, 3, 3, 1, 3, 1)
	})

	t.Run("disable profiling at span level", func(t *testing.T) {
		mockRT := newMockRuntimeTracer(true)
		globalRuntimeTracer = mockRT

		ctx := t.Context()
		_, span := tracer.Start(ctx, "root-span", trace.NoProfiling())
		assertCalls(t, mockRT, 1, 0, 0, 0, 0)

		span.End()
		assertCalls(t, mockRT, 1, 0, 0, 0, 0)
	})
}

func TestRuntimeTracerDisabled(t *testing.T) {
	originalRuntimeTracer := globalRuntimeTracer
	t.Cleanup(func() {
		globalRuntimeTracer = originalRuntimeTracer
	})

	mockRT := newMockRuntimeTracer(false)
	globalRuntimeTracer = mockRT

	tracerProvider := NewTracerProvider(WithSampler(AlwaysSample()))
	tracer := tracerProvider.Tracer("TestRuntimeTracerDisabled")
	// even manually tagged spans won't create tasks
	_, span := tracer.Start(t.Context(), "root-span", trace.ProfileTask())
	assertCalls(t, mockRT, 1, 0, 0, 0, 0)
	span.End()
	assertCalls(t, mockRT, 1, 0, 0, 0, 0)
}

func TestInstrumentationDisabled(t *testing.T) {
	originalRuntimeTracer := globalRuntimeTracer
	t.Cleanup(func() {
		globalRuntimeTracer = originalRuntimeTracer
	})

	mockRT := newMockRuntimeTracer(true)
	globalRuntimeTracer = mockRT

	tracerProvider := NewTracerProvider(WithSampler(AlwaysSample()))
	tracer := tracerProvider.Tracer("TestInstrumentationDisabled", trace.WithProfilingMode(trace.ProfilingDisabled))
	// even manually tagged spans won't create tasks
	_, span := tracer.Start(t.Context(), "root-span", trace.ProfileTask())
	assertCalls(t, mockRT, 1, 0, 0, 0, 0)
	span.End()
	assertCalls(t, mockRT, 1, 0, 0, 0, 0)
}
