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

var (
	// runtimeRatesMu serializes global Go runtime profiling rate transitions
	// across concurrent RuntimeProducer instances.
	runtimeRatesMu     sync.Mutex
	mutexActiveCount   int
	savedMutexFraction int
	blockActiveCount   int
)

type runtimeProducerConfig struct {
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

func defaultRuntimeProducerConfig() runtimeProducerConfig {
	return runtimeProducerConfig{
		heap:          true,
		allocs:        true,
		goroutine:     true,
		deltaProfiles: true,
	}
}

// RuntimeProducerOption configures a [RuntimeProducer].
type RuntimeProducerOption interface {
	applyRuntime(runtimeProducerConfig) runtimeProducerConfig
}

type runtimeProducerOptionFunc func(runtimeProducerConfig) runtimeProducerConfig

func (f runtimeProducerOptionFunc) applyRuntime(c runtimeProducerConfig) runtimeProducerConfig {
	return f(c)
}

// WithHeapProfile enables live heap (inuse_space and inuse_objects) profile
// collection. Heap profiling is enabled by default.
func WithHeapProfile() RuntimeProducerOption {
	return runtimeProducerOptionFunc(func(c runtimeProducerConfig) runtimeProducerConfig {
		c.heap = true
		return c
	})
}

// WithoutHeapProfile disables live heap profile collection.
func WithoutHeapProfile() RuntimeProducerOption {
	return runtimeProducerOptionFunc(func(c runtimeProducerConfig) runtimeProducerConfig {
		c.heap = false
		return c
	})
}

// WithAllocsProfile enables allocation (alloc_space and alloc_objects) profile
// collection. Allocation profiling is enabled by default.
func WithAllocsProfile() RuntimeProducerOption {
	return runtimeProducerOptionFunc(func(c runtimeProducerConfig) runtimeProducerConfig {
		c.allocs = true
		return c
	})
}

// WithoutAllocsProfile disables allocation profile collection.
func WithoutAllocsProfile() RuntimeProducerOption {
	return runtimeProducerOptionFunc(func(c runtimeProducerConfig) runtimeProducerConfig {
		c.allocs = false
		return c
	})
}

// WithGoroutineProfile enables goroutine stack snapshot collection. Goroutine
// profiling is enabled by default.
func WithGoroutineProfile() RuntimeProducerOption {
	return runtimeProducerOptionFunc(func(c runtimeProducerConfig) runtimeProducerConfig {
		c.goroutine = true
		return c
	})
}

// WithoutGoroutineProfile disables goroutine stack snapshot collection.
func WithoutGoroutineProfile() RuntimeProducerOption {
	return runtimeProducerOptionFunc(func(c runtimeProducerConfig) runtimeProducerConfig {
		c.goroutine = false
		return c
	})
}

// WithMutexProfile enables mutex contention profile collection. If runtime
// mutex profiling is currently disabled, [DefaultMutexProfileFraction] is set
// and restored to the previous fraction on Shutdown.
func WithMutexProfile() RuntimeProducerOption {
	return runtimeProducerOptionFunc(func(c runtimeProducerConfig) runtimeProducerConfig {
		c.mutex = true
		return c
	})
}

// WithMutexProfileFraction enables mutex contention profile collection and sets
// [runtime.SetMutexProfileFraction] to fraction when fraction is positive. The
// previous fraction is restored on Shutdown. Non-positive fractions are ignored.
func WithMutexProfileFraction(fraction int) RuntimeProducerOption {
	return runtimeProducerOptionFunc(func(c runtimeProducerConfig) runtimeProducerConfig {
		if fraction > 0 {
			c.mutex = true
			c.mutexFraction = fraction
		}
		return c
	})
}

// WithBlockProfile enables blocking profile collection using
// [DefaultBlockProfileRate]. Block profiling is reset to 0 on Shutdown.
func WithBlockProfile() RuntimeProducerOption {
	return runtimeProducerOptionFunc(func(c runtimeProducerConfig) runtimeProducerConfig {
		c.block = true
		return c
	})
}

// WithBlockProfileRate enables blocking profile collection and sets
// [runtime.SetBlockProfileRate] to rate (in nanoseconds) when rate is positive.
// Block profiling is reset to 0 on Shutdown. Non-positive rates are ignored.
func WithBlockProfileRate(rate int) RuntimeProducerOption {
	return runtimeProducerOptionFunc(func(c runtimeProducerConfig) runtimeProducerConfig {
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
func WithCPUProfile() RuntimeProducerOption {
	return runtimeProducerOptionFunc(func(c runtimeProducerConfig) runtimeProducerConfig {
		c.cpu = true
		return c
	})
}

// WithoutDeltaProfiles disables interval delta subtraction for cumulative Go
// runtime profiles (allocs, mutex, and block), emitting raw cumulative counts
// since process start instead.
func WithoutDeltaProfiles() RuntimeProducerOption {
	return runtimeProducerOptionFunc(func(c runtimeProducerConfig) runtimeProducerConfig {
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
// [RuntimeProducerOption] helpers. Samples carrying valid "otel.trace_id" and
// "otel.span_id" pprof string labels are converted into [profiledata.Link].
//
// Call [RuntimeProducer.Shutdown] (or register it with a [ProfilerProvider] or
// [Reader] that shuts it down) to stop CPU profiling and restore any modified
// runtime mutex or block profiling rates.
type RuntimeProducer struct {
	cfg      runtimeProducerConfig
	buildIDs *pprof.BuildIDCache

	mu               sync.Mutex
	shutdown         bool
	buf              bytes.Buffer
	cpuBufs          [2]bytes.Buffer
	cpuBufIdx        int
	cpuActive        bool
	mutexRateChanged bool
	blockRateChanged bool
	prevAllocs       *pprof.Profile
	prevMutex        *pprof.Profile
	prevBlock        *pprof.Profile
}

var (
	_ Producer           = (*RuntimeProducer)(nil)
	_ producerShutdowner = (*RuntimeProducer)(nil)
)

// NewRuntimeProducer creates and starts a [RuntimeProducer] configured with
// opts.
func NewRuntimeProducer(opts ...RuntimeProducerOption) *RuntimeProducer {
	cfg := defaultRuntimeProducerConfig()
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
		rp.mutexRateChanged = acquireMutexProfileRate(cfg.mutexFraction)
		if cfg.deltaProfiles {
			if p, err := rp.snapshotLookup("mutex"); err == nil {
				rp.prevMutex = p
			}
		}
	}

	if cfg.block {
		acquireBlockProfileRate(cfg.blockRate)
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
		if err := runtimepprof.StartCPUProfile(&rp.cpuBufs[rp.cpuBufIdx]); err != nil {
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
	if err := ctx.Err(); err != nil {
		return nil, err
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
		if e := ctx.Err(); e != nil {
			return out, joinErrors(err, e)
		}
		sps, e := rp.collectHeapAndAllocs()
		if e != nil {
			err = joinErrors(err, e)
		}
		out = append(out, sps...)
	}

	if rp.cfg.goroutine {
		if e := ctx.Err(); e != nil {
			return out, joinErrors(err, e)
		}
		sp, e := rp.collectSnapshot("goroutine", true, "")
		if e != nil {
			err = joinErrors(err, e)
		}
		if len(sp.Profiles) > 0 {
			out = append(out, sp)
		}
	}

	if rp.cfg.mutex {
		if e := ctx.Err(); e != nil {
			return out, joinErrors(err, e)
		}
		sp, e := rp.collectCumulative("mutex", "delay", &rp.prevMutex)
		if e != nil {
			err = joinErrors(err, e)
		}
		if len(sp.Profiles) > 0 {
			out = append(out, sp)
		}
	}

	if rp.cfg.block {
		if e := ctx.Err(); e != nil {
			return out, joinErrors(err, e)
		}
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
	}
	if rp.mutexRateChanged {
		releaseMutexProfileRate()
		rp.mutexRateChanged = false
	}
	if rp.blockRateChanged {
		releaseBlockProfileRate()
		rp.blockRateChanged = false
	}
	rp.buf = bytes.Buffer{}
	rp.cpuBufs = [2]bytes.Buffer{}
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
		rp.cpuBufs[rp.cpuBufIdx].Reset()
		if err := runtimepprof.StartCPUProfile(&rp.cpuBufs[rp.cpuBufIdx]); err != nil {
			return profiledata.ScopeProfiles{}, fmt.Errorf("sdk/profile: start CPU profile: %w", err)
		}
		rp.cpuActive = true
		return profiledata.ScopeProfiles{}, nil
	}

	runtimepprof.StopCPUProfile()
	rp.cpuActive = false

	// Restart CPU profiling into the alternate buffer before parsing so CPU
	// profiling stays active during gzip/protobuf decoding.
	finished := &rp.cpuBufs[rp.cpuBufIdx]
	rp.cpuBufIdx ^= 1
	next := &rp.cpuBufs[rp.cpuBufIdx]
	next.Reset()

	var restartErr error
	if err := runtimepprof.StartCPUProfile(next); err != nil {
		restartErr = fmt.Errorf("sdk/profile: restart CPU profile: %w", err)
	} else {
		rp.cpuActive = true
	}

	parsed, parseErr := pprof.Parse(finished.Bytes())
	finished.Reset()

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

func acquireMutexProfileRate(fraction int) bool {
	runtimeRatesMu.Lock()
	defer runtimeRatesMu.Unlock()

	if mutexActiveCount == 0 {
		prev := runtime.SetMutexProfileFraction(-1)
		switch {
		case fraction > 0:
			savedMutexFraction = runtime.SetMutexProfileFraction(fraction)
		case prev == 0:
			savedMutexFraction = runtime.SetMutexProfileFraction(DefaultMutexProfileFraction)
		default:
			// Respect an existing externally configured rate.
			return false
		}
		mutexActiveCount = 1
		return true
	}

	if fraction > 0 {
		runtime.SetMutexProfileFraction(fraction)
	}
	mutexActiveCount++
	return true
}

func releaseMutexProfileRate() {
	runtimeRatesMu.Lock()
	defer runtimeRatesMu.Unlock()

	if mutexActiveCount == 0 {
		return
	}
	mutexActiveCount--
	if mutexActiveCount == 0 {
		runtime.SetMutexProfileFraction(savedMutexFraction)
	}
}

func acquireBlockProfileRate(rate int) {
	if rate <= 0 {
		rate = DefaultBlockProfileRate
	}
	runtimeRatesMu.Lock()
	defer runtimeRatesMu.Unlock()

	runtime.SetBlockProfileRate(rate)
	blockActiveCount++
}

func releaseBlockProfileRate() {
	runtimeRatesMu.Lock()
	defer runtimeRatesMu.Unlock()

	if blockActiveCount == 0 {
		return
	}
	blockActiveCount--
	if blockActiveCount == 0 {
		runtime.SetBlockProfileRate(0)
	}
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
			Version:    version(),
			SchemaURL:  semconv.SchemaURL,
			Attributes: attribute.NewSet(scopeAttrs...),
		},
		Profiles: profs,
	}
}
