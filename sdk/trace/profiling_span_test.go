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

func TestProfilingSpan(t *testing.T) {
	originalRuntimeTracer := globalRuntimeTracer
	t.Cleanup(func() {
		globalRuntimeTracer = originalRuntimeTracer
	})

	tracerProvider := NewTracerProvider(WithSampler(AlwaysSample()))
	tracer := tracerProvider.Tracer("TestProfilingSpan")

	t.Run("local root span creates Task by default", func(t *testing.T) {
		mockRT := newMockRuntimeTracer(true)
		globalRuntimeTracer = mockRT

		ctx := t.Context()
		_, span := tracer.Start(ctx, "root-span")
		assertCalls(t, mockRT, 1, 1, 0, 0, 0)

		profilingSpan, ok := span.(profilingSpan)
		require.True(t, ok)
		assert.True(t, profilingSpan.profilingStarted())
	})

	t.Run("local root span WithProfileTask(true) creates Task", func(t *testing.T) {
		mockRT := newMockRuntimeTracer(true)
		globalRuntimeTracer = mockRT

		ctx := t.Context()
		_, span := tracer.Start(ctx, "root-span", trace.WithProfileTask(true))
		assertCalls(t, mockRT, 1, 1, 0, 0, 0)

		profilingSpan, ok := span.(profilingSpan)
		require.True(t, ok)
		assert.True(t, profilingSpan.profilingStarted())
	})

	t.Run("local root span WithProfileTask(false) does not create Task", func(t *testing.T) {
		mockRT := newMockRuntimeTracer(true)
		globalRuntimeTracer = mockRT

		ctx := t.Context()
		_, span := tracer.Start(ctx, "root-span", trace.WithProfileTask(false))
		assertCalls(t, mockRT, 1, 0, 0, 0, 0)

		profilingSpan, ok := span.(profilingSpan)
		require.True(t, ok)
		assert.False(t, profilingSpan.profilingStarted())
	})

	t.Run("local child span WithProfileRegion(true) creates Region", func(t *testing.T) {
		mockRT := newMockRuntimeTracer(true)
		globalRuntimeTracer = mockRT

		ctx := t.Context()
		rootCtx, _ := tracer.Start(ctx, "root-span")
		assertCalls(t, mockRT, 1, 1, 0, 0, 0) // root span should create a task

		_, childSpan := tracer.Start(rootCtx, "child-span", trace.WithProfileRegion(true))
		assertCalls(t, mockRT, 2, 1, 1, 0, 0) // child span should create a region

		profilingSpan, ok := childSpan.(profilingSpan)
		require.True(t, ok)
		assert.True(t, profilingSpan.profilingStarted())
	})

	t.Run("local child span WithProfileTask(true) creates Task", func(t *testing.T) {
		mockRT := newMockRuntimeTracer(true)
		globalRuntimeTracer = mockRT

		ctx := t.Context()
		rootCtx, _ := tracer.Start(ctx, "root-span")
		assertCalls(t, mockRT, 1, 1, 0, 0, 0) // root span should create a task

		_, childSpan := tracer.Start(rootCtx, "child-span", trace.WithProfileTask(true))
		assertCalls(t, mockRT, 2, 2, 0, 0, 0) // child span should create another task

		profilingSpan, ok := childSpan.(profilingSpan)
		require.True(t, ok)
		assert.True(t, profilingSpan.profilingStarted())
	})

	t.Run("local child span without profiling options does not create Task or Region", func(t *testing.T) {
		mockRT := newMockRuntimeTracer(true)
		globalRuntimeTracer = mockRT

		ctx := t.Context()
		rootCtx, _ := tracer.Start(ctx, "root-span")
		assertCalls(t, mockRT, 1, 1, 0, 0, 0) // root span should create a task

		_, childSpan := tracer.Start(rootCtx, "child-span")
		assertCalls(t, mockRT, 2, 1, 0, 0, 0)

		profilingChildSpan, ok := childSpan.(profilingSpan)
		require.True(t, ok)
		assert.False(t, profilingChildSpan.profilingStarted())
	})

	t.Run("profiling options ignored when runtime trace is not enabled", func(t *testing.T) {
		mockRT := newMockRuntimeTracer(false)
		globalRuntimeTracer = mockRT

		ctx := t.Context()
		_, span := tracer.Start(ctx, "root-span", trace.WithProfileTask(true))
		assertCalls(t, mockRT, 1, 0, 0, 0, 0)

		profilingSpan, ok := span.(profilingSpan)
		require.True(t, ok)
		assert.False(t, profilingSpan.profilingStarted())
	})

	t.Run("WithProfileRegion(true) is converted to task for root local span", func(t *testing.T) {
		mockRT := newMockRuntimeTracer(true)
		globalRuntimeTracer = mockRT

		ctx := t.Context()
		_, span := tracer.Start(ctx, "root-span", trace.WithProfileRegion(true))
		assertCalls(t, mockRT, 1, 1, 0, 0, 0)

		profilingSpan, ok := span.(profilingSpan)
		require.True(t, ok)
		assert.True(t, profilingSpan.profilingStarted())
	})

	t.Run("special case: creating a region without an associated task", func(t *testing.T) {
		// This is allowed by runtime/trace. The region will be associated with the background task. Given our
		// implementation, the option WithProfileTask(false) is required to achieve this.

		mockRT := newMockRuntimeTracer(true)
		globalRuntimeTracer = mockRT

		ctx := t.Context()
		_, span := tracer.Start(ctx, "root-span", trace.WithProfileTask(false), trace.WithProfileRegion(true))
		assertCalls(t, mockRT, 1, 0, 1, 0, 0)

		profilingSpan, ok := span.(profilingSpan)
		require.True(t, ok)
		assert.True(t, profilingSpan.profilingStarted())
	})

	t.Run("WithProfileTask takes precedence over WithProfileRegion regardless of option order", func(t *testing.T) {
		t.Run("WithProfileTask first", func(t *testing.T) {
			mockRT := newMockRuntimeTracer(true)
			globalRuntimeTracer = mockRT

			ctx := t.Context()
			rootCtx, _ := tracer.Start(ctx, "root-span")
			assertCalls(t, mockRT, 1, 1, 0, 0, 0) // root span should create a task

			_, _ = tracer.Start(rootCtx, "child-span", trace.WithProfileTask(true), trace.WithProfileRegion(true))
			assertCalls(t, mockRT, 2, 2, 0, 0, 0)
		})
		t.Run("WithProfileRegion first", func(t *testing.T) {
			mockRT := newMockRuntimeTracer(true)
			globalRuntimeTracer = mockRT

			ctx := t.Context()
			rootCtx, _ := tracer.Start(ctx, "root-span")
			assertCalls(t, mockRT, 1, 1, 0, 0, 0) // root span should create a task

			_, _ = tracer.Start(rootCtx, "child-span", trace.WithProfileRegion(true), trace.WithProfileTask(true))
			assertCalls(t, mockRT, 2, 2, 0, 0, 0)
		})
	})

	t.Run("profiling task ends when span ends", func(t *testing.T) {
		mockRT := newMockRuntimeTracer(true)
		globalRuntimeTracer = mockRT

		ctx := t.Context()
		_, span := tracer.Start(ctx, "root-span")
		assertCalls(t, mockRT, 1, 1, 0, 0, 0)

		span.End()
		assertCalls(t, mockRT, 1, 1, 0, 1, 0)
	})

	t.Run("profiling region ends when span ends", func(t *testing.T) {
		mockRT := newMockRuntimeTracer(true)
		globalRuntimeTracer = mockRT

		ctx := t.Context()
		rootCtx, span := tracer.Start(ctx, "root-span")
		assertCalls(t, mockRT, 1, 1, 0, 0, 0)

		_, childSpan := tracer.Start(rootCtx, "child-span", trace.WithProfileRegion(true))
		assertCalls(t, mockRT, 2, 1, 1, 0, 0)

		childSpan.End()
		assertCalls(t, mockRT, 2, 1, 1, 0, 1)

		childProfilingSpan, ok := childSpan.(profilingSpan)
		require.True(t, ok)
		assert.True(t, childProfilingSpan.profilingStarted()) // should return true even after it ended

		span.End()
		assertCalls(t, mockRT, 2, 1, 1, 1, 1)

		rootProfilingSpan, ok := span.(profilingSpan)
		require.True(t, ok)
		assert.True(t, rootProfilingSpan.profilingStarted()) // should return true even after it ended
	})
}
