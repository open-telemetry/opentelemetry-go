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

// testRuntimeTraceAPI is a simple mock implementation of runtimeTraceAPI for testing
type testRuntimeTraceAPI struct {
	isEnabled        bool
	isEnabledCalls   atomic.Int64
	newTaskCalls     atomic.Int64
	startRegionCalls atomic.Int64
	taskEndCalls     atomic.Int64
	regionEndCalls   atomic.Int64
}

func newTestRuntimeTraceAPI(enabled bool) *testRuntimeTraceAPI {
	return &testRuntimeTraceAPI{
		isEnabled: enabled,
	}
}

func (m *testRuntimeTraceAPI) IsEnabled() bool {
	m.isEnabledCalls.Add(1)
	return m.isEnabled
}

func (m *testRuntimeTraceAPI) NewTask(ctx context.Context, name string) (context.Context, endFunc) {
	m.newTaskCalls.Add(1)
	endFunc := func() {
		m.taskEndCalls.Add(1)
	}
	return ctx, endFunc
}

func (m *testRuntimeTraceAPI) StartRegion(ctx context.Context, name string) endFunc {
	m.startRegionCalls.Add(1)
	endFunc := func() {
		m.regionEndCalls.Add(1)
	}
	return endFunc
}

func TestProfilingSpan_Start(t *testing.T) {
	originalRuntimeTracer := globalRuntimeTracer
	t.Cleanup(func() {
		globalRuntimeTracer = originalRuntimeTracer
	})

	tracerProvider := NewTracerProvider(WithSampler(AlwaysSample()))
	tracer := tracerProvider.Tracer("TestProfilingSpan_StartProfiling")

	t.Run("local root span creates Task by default", func(t *testing.T) {
		mockRT := newTestRuntimeTraceAPI(true)
		globalRuntimeTracer = mockRT

		ctx := t.Context()
		_, span := tracer.Start(ctx, "root-span")

		profilingSpan, ok := span.(profilingSpan)
		require.True(t, ok)

		assert.Equal(t, int64(1), mockRT.newTaskCalls.Load())
		assert.Equal(t, int64(0), mockRT.startRegionCalls.Load())
		assert.Equal(t, int64(0), mockRT.taskEndCalls.Load())
		assert.Equal(t, int64(0), mockRT.regionEndCalls.Load())
		assert.True(t, profilingSpan.profilingStarted())
	})

	t.Run("local root span WithProfileTask(true) creates Task", func(t *testing.T) {
		mockRT := newTestRuntimeTraceAPI(true)
		globalRuntimeTracer = mockRT

		ctx := t.Context()
		_, span := tracer.Start(ctx, "root-span", trace.WithProfileTask(true))

		profilingSpan, ok := span.(profilingSpan)
		require.True(t, ok)

		assert.Equal(t, int64(1), mockRT.newTaskCalls.Load())
		assert.Equal(t, int64(0), mockRT.startRegionCalls.Load())
		assert.Equal(t, int64(0), mockRT.taskEndCalls.Load())
		assert.Equal(t, int64(0), mockRT.regionEndCalls.Load())
		assert.True(t, profilingSpan.profilingStarted())
	})

	t.Run("local root span WithProfileTask(false) does not create Task", func(t *testing.T) {
		mockRT := newTestRuntimeTraceAPI(true)
		globalRuntimeTracer = mockRT

		ctx := t.Context()
		_, span := tracer.Start(ctx, "root-span", trace.WithProfileTask(false))

		profilingSpan, ok := span.(profilingSpan)
		require.True(t, ok)

		assert.Equal(t, int64(0), mockRT.newTaskCalls.Load())
		assert.Equal(t, int64(0), mockRT.startRegionCalls.Load())
		assert.Equal(t, int64(0), mockRT.taskEndCalls.Load())
		assert.Equal(t, int64(0), mockRT.regionEndCalls.Load())
		assert.False(t, profilingSpan.profilingStarted())
	})

	t.Run("local child span WithProfileRegion(true) creates Region", func(t *testing.T) {
		mockRT := newTestRuntimeTraceAPI(true)
		globalRuntimeTracer = mockRT

		ctx := t.Context()
		rootCtx, _ := tracer.Start(ctx, "root-span")

		require.Equal(t, int64(1), mockRT.newTaskCalls.Load()) // root span should have created a task
		assert.Equal(t, int64(0), mockRT.startRegionCalls.Load())
		assert.Equal(t, int64(0), mockRT.taskEndCalls.Load())
		assert.Equal(t, int64(0), mockRT.regionEndCalls.Load())

		_, childSpan := tracer.Start(rootCtx, "child-span", trace.WithProfileRegion(true))

		profilingSpan, ok := childSpan.(profilingSpan)
		require.True(t, ok)
		assert.Equal(t, int64(1), mockRT.newTaskCalls.Load()) // child span should have created a task
		assert.Equal(t, int64(1), mockRT.startRegionCalls.Load())
		assert.Equal(t, int64(0), mockRT.taskEndCalls.Load())
		assert.Equal(t, int64(0), mockRT.regionEndCalls.Load())
		assert.True(t, profilingSpan.profilingStarted())
	})

	t.Run("local child span without profiling options does not create Task or Region", func(t *testing.T) {
		mockRT := newTestRuntimeTraceAPI(true)
		globalRuntimeTracer = mockRT

		ctx := t.Context()
		rootCtx, _ := tracer.Start(ctx, "root-span")

		assert.Equal(t, int64(1), mockRT.newTaskCalls.Load()) // root span should have created a task
		assert.Equal(t, int64(0), mockRT.startRegionCalls.Load())
		assert.Equal(t, int64(0), mockRT.taskEndCalls.Load())
		assert.Equal(t, int64(0), mockRT.regionEndCalls.Load())

		_, childSpan := tracer.Start(rootCtx, "child-span")

		profilingChildSpan, ok := childSpan.(profilingSpan)
		require.True(t, ok)
		assert.Equal(t, int64(1), mockRT.newTaskCalls.Load())
		assert.Equal(t, int64(0), mockRT.startRegionCalls.Load())
		assert.Equal(t, int64(0), mockRT.taskEndCalls.Load())
		assert.Equal(t, int64(0), mockRT.regionEndCalls.Load())
		assert.False(t, profilingChildSpan.profilingStarted())
	})

	t.Run("profiling disabled when runtime trace is not enabled", func(t *testing.T) {
		mockRT := newTestRuntimeTraceAPI(false)
		globalRuntimeTracer = mockRT

		ctx := t.Context()
		_, span := tracer.Start(ctx, "root-span")
		profilingSpan, ok := span.(profilingSpan)
		require.True(t, ok)
		assert.Equal(t, int64(0), mockRT.newTaskCalls.Load())
		assert.Equal(t, int64(0), mockRT.startRegionCalls.Load())
		assert.Equal(t, int64(0), mockRT.taskEndCalls.Load())
		assert.Equal(t, int64(0), mockRT.regionEndCalls.Load())
		assert.False(t, profilingSpan.profilingStarted())
	})

	t.Run("end span with task", func(t *testing.T) {
		mockRT := newTestRuntimeTraceAPI(true)
		globalRuntimeTracer = mockRT

		ctx := t.Context()
		_, span := tracer.Start(ctx, "root-span")

		assert.Equal(t, int64(1), mockRT.newTaskCalls.Load())
		assert.Equal(t, int64(0), mockRT.startRegionCalls.Load())
		assert.Equal(t, int64(0), mockRT.taskEndCalls.Load())
		assert.Equal(t, int64(0), mockRT.regionEndCalls.Load())

		span.End()

		assert.Equal(t, int64(1), mockRT.newTaskCalls.Load())
		assert.Equal(t, int64(0), mockRT.startRegionCalls.Load())
		assert.Equal(t, int64(1), mockRT.taskEndCalls.Load())
		assert.Equal(t, int64(0), mockRT.regionEndCalls.Load())
	})

	t.Run("end span with region", func(t *testing.T) {
		mockRT := newTestRuntimeTraceAPI(true)
		globalRuntimeTracer = mockRT

		ctx := t.Context()
		rootCtx, span := tracer.Start(ctx, "root-span")

		assert.Equal(t, int64(1), mockRT.newTaskCalls.Load())
		assert.Equal(t, int64(0), mockRT.startRegionCalls.Load())
		assert.Equal(t, int64(0), mockRT.taskEndCalls.Load())
		assert.Equal(t, int64(0), mockRT.regionEndCalls.Load())

		_, childSpan := tracer.Start(rootCtx, "child-span", trace.WithProfileRegion(true))

		assert.Equal(t, int64(1), mockRT.newTaskCalls.Load())
		assert.Equal(t, int64(1), mockRT.startRegionCalls.Load())
		assert.Equal(t, int64(0), mockRT.taskEndCalls.Load())
		assert.Equal(t, int64(0), mockRT.regionEndCalls.Load())

		childSpan.End()

		assert.Equal(t, int64(1), mockRT.newTaskCalls.Load())
		assert.Equal(t, int64(1), mockRT.startRegionCalls.Load())
		assert.Equal(t, int64(0), mockRT.taskEndCalls.Load())
		assert.Equal(t, int64(1), mockRT.regionEndCalls.Load())
		childProfilingSpan, ok := childSpan.(profilingSpan)
		require.True(t, ok)
		assert.True(t, childProfilingSpan.profilingStarted()) // should return true even after it ended

		span.End()

		assert.Equal(t, int64(1), mockRT.newTaskCalls.Load())
		assert.Equal(t, int64(1), mockRT.startRegionCalls.Load())
		assert.Equal(t, int64(1), mockRT.taskEndCalls.Load())
		assert.Equal(t, int64(1), mockRT.regionEndCalls.Load())
		rootProfilingSpan, ok := span.(profilingSpan)
		require.True(t, ok)
		assert.True(t, rootProfilingSpan.profilingStarted()) // should return true even after it ended
	})
}
