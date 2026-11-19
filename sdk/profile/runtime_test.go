// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package profile_test

import (
	"bytes"
	"context"
	"runtime"
	runtimepprof "runtime/pprof"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/sdk/profile"
	"go.opentelemetry.io/otel/sdk/profile/profiledata"
)

func TestRuntimeProducerDefaultProfiles(t *testing.T) {
	rp := profile.NewRuntimeProducer()
	t.Cleanup(func() { require.NoError(t, rp.Shutdown(t.Context())) })

	// Retain allocations across a GC cycle so both heap inuse and allocs
	// profiles always have observable samples.
	sink := make([][]byte, 0, 64)
	for range 64 {
		sink = append(sink, make([]byte, 32*1024))
	}
	runtime.GC()

	sps, err := rp.Produce(t.Context())
	runtime.KeepAlive(sink)
	require.NoError(t, err)

	scopes := make(map[string]profiledata.ScopeProfiles, len(sps))
	for _, sp := range sps {
		scopes[sp.Scope.Name] = sp
	}

	assert.Contains(t, scopes, "go.opentelemetry.io/otel/sdk/profile/runtime/heap")
	assert.Contains(t, scopes, "go.opentelemetry.io/otel/sdk/profile/runtime/goroutine")
	assert.NotContains(t, scopes, "go.opentelemetry.io/otel/sdk/profile/runtime/cpu")
	assert.NotContains(t, scopes, "go.opentelemetry.io/otel/sdk/profile/runtime/mutex")
	assert.NotContains(t, scopes, "go.opentelemetry.io/otel/sdk/profile/runtime/block")

	heapSP := scopes["go.opentelemetry.io/otel/sdk/profile/runtime/heap"]
	require.NotEmpty(t, heapSP.Profiles)
	assert.Equal(t, profiledata.ValueType{Type: "inuse_space", Unit: "bytes"}, heapSP.Profiles[0].SampleType)
	assert.Zero(t, heapSP.Profiles[0].Duration)

	goroutineSP := scopes["go.opentelemetry.io/otel/sdk/profile/runtime/goroutine"]
	require.Len(t, goroutineSP.Profiles, 1)
	assert.Equal(t, profiledata.ValueType{Type: "goroutine", Unit: "count"}, goroutineSP.Profiles[0].SampleType)
	assert.Zero(t, goroutineSP.Profiles[0].Duration)
}

func TestRuntimeProducerMutexAndBlockRatesRestored(t *testing.T) {
	origMutex := runtime.SetMutexProfileFraction(-1)
	t.Cleanup(func() { runtime.SetMutexProfileFraction(origMutex) })
	t.Cleanup(func() { runtime.SetBlockProfileRate(0) })

	runtime.SetMutexProfileFraction(7)

	rp := profile.NewRuntimeProducer(
		profile.WithoutHeapProfile(),
		profile.WithoutAllocsProfile(),
		profile.WithoutGoroutineProfile(),
		profile.WithMutexProfileFraction(1),
		profile.WithBlockProfileRate(1),
	)

	assert.Equal(t, 1, runtime.SetMutexProfileFraction(-1))

	// Generate contention on a mutex so both mutex and block profiles record samples.
	var mu sync.Mutex
	var wg sync.WaitGroup
	ready := make(chan struct{})
	mu.Lock()
	wg.Go(func() {
		close(ready)
		mu.Lock()
		defer mu.Unlock()
	})
	<-ready
	time.Sleep(5 * time.Millisecond)
	mu.Unlock()
	wg.Wait()

	sps, err := rp.Produce(t.Context())
	require.NoError(t, err)

	scopes := make(map[string]profiledata.ScopeProfiles, len(sps))
	for _, sp := range sps {
		scopes[sp.Scope.Name] = sp
	}
	assert.Contains(t, scopes, "go.opentelemetry.io/otel/sdk/profile/runtime/mutex")
	assert.Contains(t, scopes, "go.opentelemetry.io/otel/sdk/profile/runtime/block")

	require.NoError(t, rp.Shutdown(t.Context()))
	assert.Equal(t, 7, runtime.SetMutexProfileFraction(-1), "mutex fraction must be restored on Shutdown")

	// Subsequent Produce after Shutdown returns nil, nil.
	after, err := rp.Produce(t.Context())
	require.NoError(t, err)
	assert.Nil(t, after)
}

func TestRuntimeProducerCPUProfile(t *testing.T) {
	t.Run("collects_cpu_profile", func(t *testing.T) {
		rp := profile.NewRuntimeProducer(
			profile.WithoutHeapProfile(),
			profile.WithoutAllocsProfile(),
			profile.WithoutGoroutineProfile(),
			profile.WithCPUProfile(),
		)
		t.Cleanup(func() { require.NoError(t, rp.Shutdown(t.Context())) })

		deadline := time.Now().Add(40 * time.Millisecond)
		x := 0
		for time.Now().Before(deadline) {
			x++
		}
		_ = x

		sps, err := rp.Produce(t.Context())
		require.NoError(t, err)
		for _, sp := range sps {
			assert.Equal(t, "go.opentelemetry.io/otel/sdk/profile/runtime/cpu", sp.Scope.Name)
			require.NotEmpty(t, sp.Profiles)
			assert.Equal(t, profiledata.ValueType{Type: "cpu", Unit: "nanoseconds"}, sp.Profiles[0].SampleType)
		}
	})

	t.Run("handles_already_active_cpu_profiler", func(t *testing.T) {
		var external bytes.Buffer
		require.NoError(t, runtimepprof.StartCPUProfile(&external))
		defer runtimepprof.StopCPUProfile()

		rp := profile.NewRuntimeProducer(
			profile.WithoutHeapProfile(),
			profile.WithoutAllocsProfile(),
			profile.WithoutGoroutineProfile(),
			profile.WithCPUProfile(),
		)
		t.Cleanup(func() { require.NoError(t, rp.Shutdown(t.Context())) })

		_, err := rp.Produce(t.Context())
		assert.Error(t, err)
	})
}

func TestRuntimeProducerCumulativeModeAndCancellation(t *testing.T) {
	origMutex := runtime.SetMutexProfileFraction(-1)
	t.Cleanup(func() { runtime.SetMutexProfileFraction(origMutex) })
	t.Cleanup(func() { runtime.SetBlockProfileRate(0) })

	// Exercise explicit opt-in helpers.
	all := profile.NewRuntimeProducer(
		profile.WithHeapProfile(),
		profile.WithGoroutineProfile(),
		profile.WithMutexProfile(),
		profile.WithBlockProfile(),
	)
	require.NoError(t, all.Shutdown(t.Context()))

	rp := profile.NewRuntimeProducer(
		profile.WithoutHeapProfile(),
		profile.WithoutGoroutineProfile(),
		profile.WithAllocsProfile(),
		profile.WithoutDeltaProfiles(),
	)
	t.Cleanup(func() { require.NoError(t, rp.Shutdown(t.Context())) })

	sink := make([][]byte, 64)
	for i := range sink {
		sink[i] = make([]byte, 32*1024)
	}
	runtime.GC()

	sps, err := rp.Produce(t.Context())
	runtime.KeepAlive(sink)
	require.NoError(t, err)
	require.Len(t, sps, 1)
	assert.Equal(t, "go.opentelemetry.io/otel/sdk/profile/runtime/allocs", sps[0].Scope.Name)
	assert.Equal(t, profiledata.ValueType{Type: "alloc_space", Unit: "bytes"}, sps[0].Profiles[0].SampleType)

	canceledCtx, cancel := context.WithCancel(t.Context())
	cancel()
	_, err = rp.Produce(canceledCtx)
	assert.ErrorIs(t, err, context.Canceled)
}
