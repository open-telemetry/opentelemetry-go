// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package profile

import (
	"bytes"
	"context"
	"fmt"
	"runtime"
	runtimepprof "runtime/pprof"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/profile/internal/pprof"
	"go.opentelemetry.io/otel/sdk/profile/profiledata"
	semconv "go.opentelemetry.io/otel/semconv/v1.45.0"
)

const (
	// DefaultMutexProfileFraction is the sampling fraction configured when
	// [WithMutexProfile] is enabled and mutex profiling is not already active.
	DefaultMutexProfileFraction = 5
	// DefaultBlockProfileRate is the block profile sampling rate in nanoseconds
	// configured when [WithBlockProfile] is enabled without an explicit rate.
	DefaultBlockProfileRate = 10_000

	runtimeScopePrefix = "go.opentelemetry.io/otel/sdk/profile/runtime/"
)

type runtimeConfig struct {
	heap          bool
	allocs        bool
	goroutine     bool
	mutex         bool
	mutexFraction int
	block         bool
	blockRate     int
	cpu           bool
	deltaProfiles bool
}

func defaultRuntimeConfig() runtimeConfig {
	return runtimeConfig{
		heap:          true,
		allocs:        true,
		goroutine:     true,
		deltaProfiles: true,
	}
}

// RuntimeOption configures a [RuntimeProducer].
type RuntimeOption interface {
	applyRuntime(runtimeConfig) runtimeConfig
}

type runtimeOptionFunc func(runtimeConfig) runtimeConfig

func (f runtimeOptionFunc) applyRuntime(c runtimeConfig) runtimeConfig { return f(c) }

// WithHeapProfile enables live heap (inuse_space and inuse_objects) profile
// collection. Heap profiling is enabled by default.
func WithHeapProfile() RuntimeOption {
	return runtimeOptionFunc(func(c runtimeConfig) runtimeConfig {
		c.heap = true
		return c
	})
}

// WithoutHeapProfile disables live heap profile collection.
func WithoutHeapProfile() RuntimeOption {
	return runtimeOptionFunc(func(c runtimeConfig) runtimeConfig {
		c.heap = false
		return c
	})
}

// WithAllocsProfile enables allocation (alloc_space and alloc_objects) profile
// collection. Allocation profiling is enabled by default.
func WithAllocsProfile() RuntimeOption {
	return runtimeOptionFunc(func(c runtimeConfig) runtimeConfig {
		c.allocs = true
		return c
	})
}

// WithoutAllocsProfile disables allocation profile collection.
func WithoutAllocsProfile() RuntimeOption {
	return runtimeOptionFunc(func(c runtimeConfig) runtimeConfig {
		c.allocs = false
		return c
	})
}

// WithGoroutineProfile enables goroutine stack snapshot collection. Goroutine
// profiling is enabled by default.
func WithGoroutineProfile() RuntimeOption {
	return runtimeOptionFunc(func(c runtimeConfig) runtimeConfig {
		c.goroutine = true
		return c
	})
}

// WithoutGoroutineProfile disables goroutine stack snapshot collection.
func WithoutGoroutineProfile() RuntimeOption {
	return runtimeOptionFunc(func(c runtimeConfig) runtimeConfig {
		c.goroutine = false
		return c
	})
}

// WithMutexProfile enables mutex contention profile collection. If runtime
// mutex profiling is currently disabled, [DefaultMutexProfileFraction] is set
// and restored to the previous fraction on Shutdown.
func WithMutexProfile() RuntimeOption {
	return runtimeOptionFunc(func(c runtimeConfig) runtimeConfig {
		c.mutex = true
		return c
	})
}

// WithMutexProfileFraction enables mutex contention profile collection and sets
// [runtime.SetMutexProfileFraction] to fraction when fraction is positive. The
// previous fraction is restored on Shutdown. Non-positive fractions are ignored.
func WithMutexProfileFraction(fraction int) RuntimeOption {
	return runtimeOptionFunc(func(c runtimeConfig) runtimeConfig {
		if fraction > 0 {
			c.mutex = true
			c.mutexFraction = fraction
		}
		return c
	})
}

// WithBlockProfile enables blocking profile collection using
// [DefaultBlockProfileRate]. Block profiling is reset to 0 on Shutdown.
func WithBlockProfile() RuntimeOption {
	return runtimeOptionFunc(func(c runtimeConfig) runtimeConfig {
		c.block = true
		return c
	})
}

// WithBlockProfileRate enables blocking profile collection and sets
// [runtime.SetBlockProfileRate] to rate (in nanoseconds) when rate is positive.
// Block profiling is reset to 0 on Shutdown. Non-positive rates are ignored.
func WithBlockProfileRate(rate int) RuntimeOption {
	return runtimeOptionFunc(func(c runtimeConfig) runtimeConfig {
		if rate > 0 {
			c.block = true
			c.blockRate = rate
		}
		return c
	})
}

// WithCPUProfile enables continuous in-process CPU profiling via
// [runtime/pprof.StartCPUProfile]. Because the Go runtime permits only one
// active CPU profiler per process, CPU profiling is disabled by default.
func WithCPUProfile() RuntimeOption {
	return runtimeOptionFunc(func(c runtimeConfig) runtimeConfig {
		c.cpu = true
		return c
	})
}

// WithoutDeltaProfiles disables interval delta subtraction for cumulative Go
// runtime profiles (allocs, mutex, and block), emitting raw cumulative counts
// since process start instead.
func WithoutDeltaProfiles() RuntimeOption {
	return runtimeOptionFunc(func(c runtimeConfig) runtimeConfig {
		c.deltaProfiles = false
		return c
	})
}

// RuntimeProducer collects profiles from the Go runtime via [runtime/pprof] and
// converts them into [profiledata.ScopeProfiles].
//
// By default, RuntimeProducer collects live heap, allocation (delta), and
// goroutine profiles without mutating global Go runtime profiling rates or
// starting CPU profiling. Mutex, block, and CPU profiling can be enabled via
// [RuntimeOption] helpers.
//
// Call [RuntimeProducer.Shutdown] (or register it with a [ProfilerProvider] or
// [Reader] that shuts it down) to stop CPU profiling and restore any modified
// runtime mutex or block profiling rates.
type RuntimeProducer struct {
	cfg      runtimeConfig
	buildIDs *pprof.BuildIDCache

	mu                sync.Mutex
	shutdown          bool
	buf               bytes.Buffer
	cpuBuf            bytes.Buffer
	cpuActive         bool
	mutexRateChanged  bool
	prevMutexFraction int
	blockRateChanged  bool
	prevAllocs        *pprof.Profile
	prevMutex         *pprof.Profile
	prevBlock         *pprof.Profile
}

var _ Producer = (*RuntimeProducer)(nil)

// NewRuntimeProducer creates and starts a [RuntimeProducer] configured with
// opts.
func NewRuntimeProducer(opts ...RuntimeOption) *RuntimeProducer {
	cfg := defaultRuntimeConfig()
	for _, opt := range opts {
		if opt != nil {
			cfg = opt.applyRuntime(cfg)
		}
	}

	rp := &RuntimeProducer{
		cfg:      cfg,
		buildIDs: pprof.NewBuildIDCache(),
	}

	if cfg.mutex {
		prev := runtime.SetMutexProfileFraction(-1)
		switch {
		case cfg.mutexFraction > 0:
			rp.prevMutexFraction = runtime.SetMutexProfileFraction(cfg.mutexFraction)
			rp.mutexRateChanged = true
		case prev == 0:
			rp.prevMutexFraction = runtime.SetMutexProfileFraction(DefaultMutexProfileFraction)
			rp.mutexRateChanged = true
		}
		if cfg.deltaProfiles {
			if p, err := rp.snapshotLookup("mutex"); err == nil {
				rp.prevMutex = p
			}
		}
	}

	if cfg.block {
		rate := cfg.blockRate
		if rate <= 0 {
			rate = DefaultBlockProfileRate
		}
		runtime.SetBlockProfileRate(rate)
		rp.blockRateChanged = true
		if cfg.deltaProfiles {
			if p, err := rp.snapshotLookup("block"); err == nil {
				rp.prevBlock = p
			}
		}
	}

	if cfg.allocs && cfg.deltaProfiles {
		if p, err := rp.snapshotLookup("heap"); err == nil {
			rp.prevAllocs = p
		}
	}

	if cfg.cpu {
		if err := runtimepprof.StartCPUProfile(&rp.cpuBuf); err != nil {
			otel.Handle(fmt.Errorf("sdk/profile: start CPU profile: %w", err))
		} else {
			rp.cpuActive = true
		}
	}

	return rp
}

// Produce collects all enabled Go runtime profiles and returns them as
// [profiledata.ScopeProfiles], one per non-empty profile kind.
//
// This method is safe to call concurrently.
func (rp *RuntimeProducer) Produce(ctx context.Context) ([]profiledata.ScopeProfiles, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	rp.mu.Lock()
	defer rp.mu.Unlock()

	if rp.shutdown {
		return nil, nil
	}

	var out []profiledata.ScopeProfiles
	var err error

	if rp.cfg.cpu {
		sp, e := rp.collectCPU()
		if e != nil {
			err = joinErrors(err, e)
		}
		if len(sp.Profiles) > 0 {
			out = append(out, sp)
		}
	}

	if rp.cfg.heap || rp.cfg.allocs {
		sps, e := rp.collectHeapAndAllocs()
		if e != nil {
			err = joinErrors(err, e)
		}
		out = append(out, sps...)
	}

	if rp.cfg.goroutine {
		sp, e := rp.collectSnapshot("goroutine", true, "")
		if e != nil {
			err = joinErrors(err, e)
		}
		if len(sp.Profiles) > 0 {
			out = append(out, sp)
		}
	}

	if rp.cfg.mutex {
		sp, e := rp.collectCumulative("mutex", "delay", &rp.prevMutex)
		if e != nil {
			err = joinErrors(err, e)
		}
		if len(sp.Profiles) > 0 {
			out = append(out, sp)
		}
	}

	if rp.cfg.block {
		sp, e := rp.collectCumulative("block", "delay", &rp.prevBlock)
		if e != nil {
			err = joinErrors(err, e)
		}
		if len(sp.Profiles) > 0 {
			out = append(out, sp)
		}
	}

	return out, err
}

// Shutdown stops CPU profiling (if started by this RuntimeProducer) and
// restores any runtime mutex or block profiling rates modified at construction.
//
// Calling Shutdown more than once is safe and performs no action.
func (rp *RuntimeProducer) Shutdown(context.Context) error {
	rp.mu.Lock()
	defer rp.mu.Unlock()

	if rp.shutdown {
		return nil
	}
	rp.shutdown = true

	if rp.cpuActive {
		runtimepprof.StopCPUProfile()
		rp.cpuActive = false
		rp.cpuBuf.Reset()
	}
	if rp.mutexRateChanged {
		runtime.SetMutexProfileFraction(rp.prevMutexFraction)
		rp.mutexRateChanged = false
	}
	if rp.blockRateChanged {
		runtime.SetBlockProfileRate(0)
		rp.blockRateChanged = false
	}
	rp.buf.Reset()
	rp.prevAllocs = nil
	rp.prevMutex = nil
	rp.prevBlock = nil
	return nil
}

func (rp *RuntimeProducer) snapshotLookup(name string) (*pprof.Profile, error) {
	prof := runtimepprof.Lookup(name)
	if prof == nil {
		return nil, fmt.Errorf("sdk/profile: unknown runtime pprof profile %q", name)
	}
	rp.buf.Reset()
	if err := prof.WriteTo(&rp.buf, 0); err != nil {
		return nil, fmt.Errorf("sdk/profile: write %s profile: %w", name, err)
	}
	parsed, err := pprof.Parse(rp.buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("sdk/profile: parse %s profile: %w", name, err)
	}
	return parsed, nil
}

func (rp *RuntimeProducer) collectCPU() (profiledata.ScopeProfiles, error) {
	if !rp.cpuActive {
		// Retry starting CPU profiling if a previous external profiler finished.
		if err := runtimepprof.StartCPUProfile(&rp.cpuBuf); err != nil {
			return profiledata.ScopeProfiles{}, fmt.Errorf("sdk/profile: start CPU profile: %w", err)
		}
		rp.cpuActive = true
		return profiledata.ScopeProfiles{}, nil
	}

	runtimepprof.StopCPUProfile()
	rp.cpuActive = false

	parsed, parseErr := pprof.Parse(rp.cpuBuf.Bytes())
	rp.cpuBuf.Reset()

	var restartErr error
	if err := runtimepprof.StartCPUProfile(&rp.cpuBuf); err != nil {
		restartErr = fmt.Errorf("sdk/profile: restart CPU profile: %w", err)
	} else {
		rp.cpuActive = true
	}

	if parseErr != nil {
		return profiledata.ScopeProfiles{}, joinErrors(
			fmt.Errorf("sdk/profile: parse cpu profile: %w", parseErr),
			restartErr,
		)
	}

	profs, scopeAttrs := pprof.Convert(parsed, pprof.ConvertOptions{
		DefaultSampleType: "cpu",
		FrameType:         semconv.ProfileFrameTypeGo,
		BuildIDs:          rp.buildIDs,
	})
	return makeRuntimeScopeProfiles("cpu", profs, scopeAttrs), restartErr
}

func (rp *RuntimeProducer) collectHeapAndAllocs() ([]profiledata.ScopeProfiles, error) {
	cur, err := rp.snapshotLookup("heap")
	if err != nil {
		return nil, err
	}

	var out []profiledata.ScopeProfiles
	if rp.cfg.heap {
		profs, scopeAttrs := pprof.Convert(cur, pprof.ConvertOptions{
			SampleTypes:       []string{"inuse_space", "inuse_objects"},
			DefaultSampleType: "inuse_space",
			FrameType:         semconv.ProfileFrameTypeGo,
			BuildIDs:          rp.buildIDs,
			Instantaneous:     true,
		})
		if len(profs) > 0 {
			out = append(out, makeRuntimeScopeProfiles("heap", profs, scopeAttrs))
		}
	}

	if rp.cfg.allocs {
		src := cur
		if rp.cfg.deltaProfiles {
			src = pprof.Delta(rp.prevAllocs, cur)
			rp.prevAllocs = cur
		}
		profs, scopeAttrs := pprof.Convert(src, pprof.ConvertOptions{
			SampleTypes:       []string{"alloc_space", "alloc_objects"},
			DefaultSampleType: "alloc_space",
			FrameType:         semconv.ProfileFrameTypeGo,
			BuildIDs:          rp.buildIDs,
		})
		if len(profs) > 0 {
			out = append(out, makeRuntimeScopeProfiles("allocs", profs, scopeAttrs))
		}
	}

	return out, nil
}

func (rp *RuntimeProducer) collectSnapshot(
	kind string,
	instantaneous bool,
	defaultType string,
) (profiledata.ScopeProfiles, error) {
	cur, err := rp.snapshotLookup(kind)
	if err != nil {
		return profiledata.ScopeProfiles{}, err
	}
	profs, scopeAttrs := pprof.Convert(cur, pprof.ConvertOptions{
		DefaultSampleType: defaultType,
		FrameType:         semconv.ProfileFrameTypeGo,
		BuildIDs:          rp.buildIDs,
		Instantaneous:     instantaneous,
	})
	return makeRuntimeScopeProfiles(kind, profs, scopeAttrs), nil
}

func (rp *RuntimeProducer) collectCumulative(
	kind string,
	defaultType string,
	prev **pprof.Profile,
) (profiledata.ScopeProfiles, error) {
	cur, err := rp.snapshotLookup(kind)
	if err != nil {
		return profiledata.ScopeProfiles{}, err
	}
	src := cur
	if rp.cfg.deltaProfiles {
		src = pprof.Delta(*prev, cur)
		*prev = cur
	}
	profs, scopeAttrs := pprof.Convert(src, pprof.ConvertOptions{
		DefaultSampleType: defaultType,
		FrameType:         semconv.ProfileFrameTypeGo,
		BuildIDs:          rp.buildIDs,
	})
	return makeRuntimeScopeProfiles(kind, profs, scopeAttrs), nil
}

func makeRuntimeScopeProfiles(
	kind string,
	profs []profiledata.Profile,
	scopeAttrs []attribute.KeyValue,
) profiledata.ScopeProfiles {
	if len(profs) == 0 {
		return profiledata.ScopeProfiles{}
	}
	return profiledata.ScopeProfiles{
		Scope: instrumentation.Scope{
			Name:       runtimeScopePrefix + kind,
			Version:    Version(),
			SchemaURL:  semconv.SchemaURL,
			Attributes: attribute.NewSet(scopeAttrs...),
		},
		Profiles: profs,
	}
}
