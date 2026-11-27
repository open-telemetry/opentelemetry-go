// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package profile_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"runtime"
	runtimepprof "runtime/pprof"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
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
	assert.Contains(t, scopes, "go.opentelemetry.io/otel/sdk/profile/runtime/allocs")
	assert.Contains(t, scopes, "go.opentelemetry.io/otel/sdk/profile/runtime/goroutine")
	assert.NotContains(t, scopes, "go.opentelemetry.io/otel/sdk/profile/runtime/cpu")
	assert.NotContains(t, scopes, "go.opentelemetry.io/otel/sdk/profile/runtime/mutex")
	assert.NotContains(t, scopes, "go.opentelemetry.io/otel/sdk/profile/runtime/block")

	heapSP := scopes["go.opentelemetry.io/otel/sdk/profile/runtime/heap"]
	require.NotEmpty(t, heapSP.Profiles)
	assert.Equal(t, profiledata.ValueType{Type: "inuse_space", Unit: "bytes"}, heapSP.Profiles[0].SampleType)
	assert.Zero(t, heapSP.Profiles[0].Duration)

	allocsSP := scopes["go.opentelemetry.io/otel/sdk/profile/runtime/allocs"]
	require.NotEmpty(t, allocsSP.Profiles)
	assert.Equal(t, profiledata.ValueType{Type: "alloc_space", Unit: "bytes"}, allocsSP.Profiles[0].SampleType)

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
		profile.WithMutexProfileFraction(0), // Ignored.
		profile.WithBlockProfileRate(-1),    // Ignored.
	)

	assert.Equal(t, 1, runtime.SetMutexProfileFraction(-1))

	// Generate contention in a bounded retry loop so slow/raced scheduling never flakes.
	contend := func() {
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
	}

	scopes := make(map[string]profiledata.ScopeProfiles)
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		contend()
		sps, err := rp.Produce(t.Context())
		require.NoError(t, err)
		for _, sp := range sps {
			scopes[sp.Scope.Name] = sp
		}
		if _, hasM := scopes["go.opentelemetry.io/otel/sdk/profile/runtime/mutex"]; hasM {
			if _, hasB := scopes["go.opentelemetry.io/otel/sdk/profile/runtime/block"]; hasB {
				break
			}
		}
	}
	assert.Contains(t, scopes, "go.opentelemetry.io/otel/sdk/profile/runtime/mutex")
	assert.Contains(t, scopes, "go.opentelemetry.io/otel/sdk/profile/runtime/block")

	// A follow-up Produce with no new mutex contention omits empty delta scopes.
	emptySps, err := rp.Produce(t.Context())
	require.NoError(t, err)
	for _, sp := range emptySps {
		assert.NotEmpty(t, sp.Profiles)
	}

	require.NoError(t, rp.Shutdown(t.Context()))
	assert.Equal(t, 7, runtime.SetMutexProfileFraction(-1), "mutex fraction must be restored on Shutdown")

	// Second Shutdown must be a no-op and must not overwrite subsequent external changes.
	runtime.SetMutexProfileFraction(11)
	require.NoError(t, rp.Shutdown(t.Context()))
	assert.Equal(t, 11, runtime.SetMutexProfileFraction(-1))

	// Subsequent Produce after Shutdown returns nil, nil.
	after, err := rp.Produce(t.Context())
	require.NoError(t, err)
	assert.Nil(t, after)

	t.Run("default_preserved_and_refcounted_rates", func(t *testing.T) {
		// Default rate when disabled (0) sets DefaultMutexProfileFraction and restores 0.
		runtime.SetMutexProfileFraction(0)
		pDefault := profile.NewRuntimeProducer(
			profile.WithoutHeapProfile(),
			profile.WithoutAllocsProfile(),
			profile.WithoutGoroutineProfile(),
			profile.WithMutexProfile(),
			profile.WithBlockProfile(),
		)
		assert.Equal(t, profile.DefaultMutexProfileFraction, runtime.SetMutexProfileFraction(-1))
		require.NoError(t, pDefault.Shutdown(t.Context()))
		assert.Equal(t, 0, runtime.SetMutexProfileFraction(-1))

		// Existing non-zero rate is preserved by WithMutexProfile().
		runtime.SetMutexProfileFraction(9)
		pPreserve := profile.NewRuntimeProducer(
			profile.WithoutHeapProfile(),
			profile.WithoutAllocsProfile(),
			profile.WithoutGoroutineProfile(),
			profile.WithMutexProfile(),
		)
		assert.Equal(t, 9, runtime.SetMutexProfileFraction(-1))
		require.NoError(t, pPreserve.Shutdown(t.Context()))
		assert.Equal(t, 9, runtime.SetMutexProfileFraction(-1))

		// Overlapping producers refcount global rates so shutting down A does not
		// disable profiling for B, and shutting down B restores the pre-SDK fraction.
		runtime.SetMutexProfileFraction(0)
		pA := profile.NewRuntimeProducer(
			profile.WithoutHeapProfile(),
			profile.WithoutAllocsProfile(),
			profile.WithoutGoroutineProfile(),
			profile.WithMutexProfileFraction(2),
			profile.WithBlockProfileRate(100),
		)
		pB := profile.NewRuntimeProducer(
			profile.WithoutHeapProfile(),
			profile.WithoutAllocsProfile(),
			profile.WithoutGoroutineProfile(),
			profile.WithMutexProfileFraction(3),
			profile.WithBlockProfileRate(200),
		)
		assert.Equal(t, 3, runtime.SetMutexProfileFraction(-1))
		require.NoError(t, pA.Shutdown(t.Context()))
		assert.Equal(t, 3, runtime.SetMutexProfileFraction(-1), "mutex profiling must stay active while pB is live")
		require.NoError(t, pB.Shutdown(t.Context()))
		assert.Equal(
			t,
			0,
			runtime.SetMutexProfileFraction(-1),
			"pre-SDK fraction must be restored after last producer shuts down",
		)
	})
}

func burnCPU(d time.Duration) {
	deadline := time.Now().Add(d)
	var buf [32]byte
	for time.Now().Before(deadline) {
		buf = sha256.Sum256(buf[:])
	}
	runtime.KeepAlive(buf)
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

		var sps []profiledata.ScopeProfiles
		for range 8 {
			burnCPU(30 * time.Millisecond)
			var err error
			sps, err = rp.Produce(t.Context())
			require.NoError(t, err)
			if len(sps) > 0 {
				break
			}
		}
		require.Len(t, sps, 1, "expected CPU profile scope after sustained CPU work")
		assert.Equal(t, "go.opentelemetry.io/otel/sdk/profile/runtime/cpu", sps[0].Scope.Name)
		require.NotEmpty(t, sps[0].Profiles)
		assert.Equal(t, profiledata.ValueType{Type: "cpu", Unit: "nanoseconds"}, sps[0].Profiles[0].SampleType)
	})

	t.Run("handles_already_active_cpu_profiler_and_recovers", func(t *testing.T) {
		var handled []error
		var handledMu sync.Mutex
		prevHandler := otel.GetErrorHandler()
		otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
			handledMu.Lock()
			handled = append(handled, err)
			handledMu.Unlock()
		}))
		t.Cleanup(func() { otel.SetErrorHandler(prevHandler) })

		var external bytes.Buffer
		require.NoError(t, runtimepprof.StartCPUProfile(&external))
		externalStopped := false
		t.Cleanup(func() {
			if !externalStopped {
				runtimepprof.StopCPUProfile()
			}
		})

		rp := profile.NewRuntimeProducer(
			profile.WithoutHeapProfile(),
			profile.WithoutAllocsProfile(),
			profile.WithoutGoroutineProfile(),
			profile.WithCPUProfile(),
		)
		t.Cleanup(func() { require.NoError(t, rp.Shutdown(t.Context())) })

		handledMu.Lock()
		assert.NotEmpty(t, handled, "expected initial StartCPUProfile failure to be reported via otel.Handle")
		handledMu.Unlock()

		_, err := rp.Produce(t.Context())
		assert.Error(t, err)

		// Once the external profiler stops, Produce re-arms CPU profiling and subsequent calls collect samples.
		runtimepprof.StopCPUProfile()
		externalStopped = true

		sps, err := rp.Produce(t.Context())
		require.NoError(t, err)
		assert.Empty(t, sps)

		for range 8 {
			burnCPU(30 * time.Millisecond)
			sps, err = rp.Produce(t.Context())
			require.NoError(t, err)
			if len(sps) > 0 {
				break
			}
		}
		require.Len(t, sps, 1)
		assert.Equal(t, "go.opentelemetry.io/otel/sdk/profile/runtime/cpu", sps[0].Scope.Name)
	})
}

func TestRuntimeProducerConcurrentSafe(t *testing.T) {
	rp := profile.NewRuntimeProducer(
		profile.WithoutHeapProfile(),
		profile.WithoutAllocsProfile(),
		profile.WithGoroutineProfile(),
	)

	var wg sync.WaitGroup
	for range 4 {
		wg.Go(func() {
			_, _ = rp.Produce(t.Context())
		})
	}
	wg.Go(func() {
		_ = rp.Shutdown(t.Context())
	})
	wg.Wait()
	require.NoError(t, rp.Shutdown(t.Context()))
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
