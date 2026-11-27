// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package profile_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/profile"
	"go.opentelemetry.io/otel/sdk/profile/profiledata"
	"go.opentelemetry.io/otel/sdk/resource"
)

type stubProducer struct {
	mu          sync.Mutex
	sps         []profiledata.ScopeProfiles
	err         error
	shutdownErr error
	calls       int
	shutdown    bool
	onShut      func()
}

func (p *stubProducer) Produce(ctx context.Context) ([]profiledata.ScopeProfiles, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.calls++
	if p.shutdown {
		return nil, errors.New("producer already shut down")
	}
	return p.sps, p.err
}

func (p *stubProducer) Shutdown(context.Context) error {
	p.mu.Lock()
	p.shutdown = true
	fn := p.onShut
	err := p.shutdownErr
	p.mu.Unlock()
	if fn != nil {
		fn()
	}
	return err
}

type recordingExporter struct {
	mu          sync.Mutex
	exported    []*profiledata.ResourceProfiles
	exportErr   error
	flushCalls  int
	flushErr    error
	shutdownErr error
	shutdown    bool
	onExport    func()
	blockExport func(context.Context) error
}

func (e *recordingExporter) Export(ctx context.Context, rp *profiledata.ResourceProfiles) error {
	if e.blockExport != nil {
		if err := e.blockExport(ctx); err != nil {
			return err
		}
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.shutdown {
		return profile.ErrExporterShutdown
	}
	if e.exportErr != nil {
		return e.exportErr
	}
	cp := *rp
	cp.ScopeProfiles = append([]profiledata.ScopeProfiles(nil), rp.ScopeProfiles...)
	e.exported = append(e.exported, &cp)
	if e.onExport != nil {
		e.onExport()
	}
	return nil
}

func (e *recordingExporter) ForceFlush(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.flushCalls++
	return e.flushErr
}

func (e *recordingExporter) Shutdown(context.Context) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.shutdown = true
	return e.shutdownErr
}

func (e *recordingExporter) count() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return len(e.exported)
}

func sampleScopeProfiles(name string) []profiledata.ScopeProfiles {
	return []profiledata.ScopeProfiles{
		{
			Scope: instrumentation.Scope{Name: name},
			Profiles: []profiledata.Profile{
				{
					SampleType: profiledata.ValueType{Type: "samples", Unit: "count"},
					Samples:    []profiledata.Sample{{Values: []int64{1}}},
				},
			},
		},
	}
}

func TestProfilerProviderLifecycle(t *testing.T) {
	res := resource.NewSchemaless(attribute.String("service.name", "test-svc"))
	var events []string
	var eventsMu sync.Mutex
	record := func(ev string) {
		eventsMu.Lock()
		events = append(events, ev)
		eventsMu.Unlock()
	}

	prod := &stubProducer{
		sps:    sampleScopeProfiles("provider-prod"),
		onShut: func() { record("producer_shutdown") },
	}
	exp := &recordingExporter{
		onExport: func() { record("export") },
	}

	pr := profile.NewPeriodicReader(exp, profile.WithInterval(time.Hour))
	pp := profile.NewProfilerProvider(
		profile.WithResource(res),
		profile.WithProducer(nil),
		profile.WithProducer(prod),
		profile.WithReader(nil),
		profile.WithReader(pr),
	)

	require.NoError(t, pp.ForceFlush(t.Context()))
	assert.Equal(t, 1, exp.count())
	assert.Equal(t, 1, exp.flushCalls)
	assert.Equal(t, res, exp.exported[0].Resource)

	// Shutdown must flush reader before shutting down provider producer.
	require.NoError(t, pp.Shutdown(t.Context()))
	assert.Equal(t, 2, exp.count())
	assert.Equal(t, []string{"export", "export", "producer_shutdown"}, events)

	// Subsequent ForceFlush and Shutdown calls are no-ops.
	require.NoError(t, pp.ForceFlush(t.Context()))
	require.NoError(t, pp.Shutdown(t.Context()))
}

func TestManualReader(t *testing.T) {
	t.Run("unregistered", func(t *testing.T) {
		mr := profile.NewManualReader()
		logData, ok := mr.MarshalLog().(struct {
			Type       string
			Registered bool
			Shutdown   bool
		})
		require.True(t, ok)
		assert.Equal(t, "ManualReader", logData.Type)
		assert.False(t, logData.Registered)
		assert.False(t, logData.Shutdown)

		var rp profiledata.ResourceProfiles
		assert.ErrorIs(t, mr.Collect(t.Context(), &rp), profile.ErrReaderNotRegistered)
		assert.Error(t, mr.Collect(t.Context(), nil))
	})

	t.Run("collect_and_shutdown", func(t *testing.T) {
		res := resource.NewSchemaless(attribute.String("service.name", "manual-svc"))
		provProd := &stubProducer{sps: sampleScopeProfiles("from-provider")}
		readerProd := &stubProducer{sps: sampleScopeProfiles("from-reader")}

		mr := profile.NewManualReader(profile.WithProducer(readerProd))
		pp := profile.NewProfilerProvider(
			profile.WithResource(res),
			profile.WithProducer(provProd),
			profile.WithReader(mr),
		)

		// Duplicate registration is ignored; original provider remains active.
		otherRes := resource.NewSchemaless(attribute.String("service.name", "other-svc"))
		_ = profile.NewProfilerProvider(profile.WithResource(otherRes), profile.WithReader(mr))

		var rp profiledata.ResourceProfiles
		require.NoError(t, mr.Collect(t.Context(), &rp))
		assert.Equal(t, res, rp.Resource)
		require.Len(t, rp.ScopeProfiles, 2)
		assert.Equal(t, "from-provider", rp.ScopeProfiles[0].Scope.Name)
		assert.Equal(t, "from-reader", rp.ScopeProfiles[1].Scope.Name)

		require.NoError(t, mr.ForceFlush(t.Context()))
		logData, ok := mr.MarshalLog().(struct {
			Type       string
			Registered bool
			Shutdown   bool
		})
		require.True(t, ok)
		assert.True(t, logData.Registered)
		assert.False(t, logData.Shutdown)

		// Canceled context returns context error on both Collect and ForceFlush.
		canceledCtx, cancel := context.WithCancel(t.Context())
		cancel()
		assert.ErrorIs(t, mr.Collect(canceledCtx, &rp), context.Canceled)
		assert.ErrorIs(t, mr.ForceFlush(canceledCtx), context.Canceled)

		require.NoError(t, pp.Shutdown(t.Context()))
		logAfter, ok := mr.MarshalLog().(struct {
			Type       string
			Registered bool
			Shutdown   bool
		})
		require.True(t, ok)
		assert.True(t, logAfter.Shutdown)
		assert.ErrorIs(t, mr.Collect(t.Context(), &rp), profile.ErrReaderShutdown)
		assert.ErrorIs(t, mr.ForceFlush(t.Context()), profile.ErrReaderShutdown)
		assert.ErrorIs(t, mr.Shutdown(t.Context()), profile.ErrReaderShutdown)
	})

	t.Run("producer_errors_joined_on_collect_and_shutdown", func(t *testing.T) {
		provErr := errors.New("provider producer boom")
		readerErr1 := errors.New("reader producer 1 boom")
		readerErr2 := errors.New("reader producer 2 boom")
		provShutErr := errors.New("provider shutdown boom")
		readerShutErr := errors.New("reader shutdown boom")

		mr := profile.NewManualReader(
			profile.WithProducer(&stubProducer{err: readerErr1, shutdownErr: readerShutErr}),
			profile.WithProducer(&stubProducer{sps: sampleScopeProfiles("partial-ok"), err: readerErr2}),
		)
		pp := profile.NewProfilerProvider(
			profile.WithProducer(&stubProducer{err: provErr, shutdownErr: provShutErr}),
			profile.WithReader(mr),
		)

		var rp profiledata.ResourceProfiles
		err := mr.Collect(t.Context(), &rp)
		assert.ErrorIs(t, err, provErr)
		assert.ErrorIs(t, err, readerErr1)
		assert.ErrorIs(t, err, readerErr2)
		require.Len(t, rp.ScopeProfiles, 1)
		assert.Equal(t, "partial-ok", rp.ScopeProfiles[0].Scope.Name)

		shutErr := pp.Shutdown(t.Context())
		assert.ErrorIs(t, shutErr, provShutErr)
		assert.ErrorIs(t, shutErr, readerShutErr)
	})
}

func TestWithResourceMergesEnvironment(t *testing.T) {
	t.Setenv("OTEL_RESOURCE_ATTRIBUTES", "env.only=from-env,service.name=from-env")
	res := resource.NewSchemaless(attribute.String("service.name", "explicit-svc"))

	mr := profile.NewManualReader()
	pp := profile.NewProfilerProvider(
		profile.WithResource(res),
		profile.WithReader(mr),
	)
	t.Cleanup(func() { _ = pp.Shutdown(t.Context()) })

	var rp profiledata.ResourceProfiles
	require.NoError(t, mr.Collect(t.Context(), &rp))
	require.NotNil(t, rp.Resource)

	val, ok := rp.Resource.Set().Value("env.only")
	require.True(t, ok, "expected OTEL_RESOURCE_ATTRIBUTES attribute to be merged")
	assert.Equal(t, "from-env", val.AsString())

	svc, ok := rp.Resource.Set().Value("service.name")
	require.True(t, ok)
	assert.Equal(t, "explicit-svc", svc.AsString(), "explicit resource should override env")
}

func TestPeriodicReader(t *testing.T) {
	t.Run("periodic_collection_and_shutdown", func(t *testing.T) {
		exported := make(chan struct{}, 10)
		exp := &recordingExporter{
			onExport: func() {
				select {
				case exported <- struct{}{}:
				default:
				}
			},
		}
		prod := &stubProducer{sps: sampleScopeProfiles("periodic")}
		pr := profile.NewPeriodicReader(
			exp,
			profile.WithInterval(10*time.Millisecond),
			profile.WithTimeout(5*time.Second),
			profile.WithInterval(-1), // Ignored.
			profile.WithTimeout(0),   // Ignored.
			profile.WithProducer(prod),
		)
		pp := profile.NewProfilerProvider(profile.WithReader(pr))

		select {
		case <-exported:
		case <-time.After(2 * time.Second):
			t.Fatal("timed out waiting for periodic profile export")
		}

		logData, ok := pr.MarshalLog().(struct {
			Type       string
			Exporter   profile.Exporter
			Registered bool
			Shutdown   bool
			Interval   time.Duration
			Timeout    time.Duration
		})
		require.True(t, ok)
		assert.Equal(t, "PeriodicReader", logData.Type)
		assert.Equal(t, exp, logData.Exporter)
		assert.True(t, logData.Registered)
		assert.False(t, logData.Shutdown)
		assert.Equal(t, 10*time.Millisecond, logData.Interval)
		assert.Equal(t, 5*time.Second, logData.Timeout)

		require.NoError(t, pp.Shutdown(t.Context()))
		assert.True(t, exp.shutdown)
		assert.True(t, prod.shutdown)
		assert.ErrorIs(t, pr.Shutdown(t.Context()), profile.ErrReaderShutdown)
		assert.ErrorIs(t, pr.ForceFlush(t.Context()), profile.ErrReaderShutdown)
	})

	t.Run("skips_export_when_empty_and_handles_nil_exporter", func(t *testing.T) {
		exp := &recordingExporter{}
		emptyProd := &stubProducer{}
		pr := profile.NewPeriodicReader(exp, profile.WithInterval(time.Hour))
		pp := profile.NewProfilerProvider(profile.WithProducer(emptyProd), profile.WithReader(pr))

		require.NoError(t, pp.ForceFlush(t.Context()))
		assert.Equal(t, 0, exp.count())
		require.NoError(t, pp.Shutdown(t.Context()))

		// Nil exporter must not panic even when non-empty profiles are exported.
		nilPR := profile.NewPeriodicReader(nil, profile.WithInterval(time.Hour))
		nilPP := profile.NewProfilerProvider(
			profile.WithProducer(&stubProducer{sps: sampleScopeProfiles("noop-export")}),
			profile.WithReader(nilPR),
		)
		require.NoError(t, nilPP.ForceFlush(t.Context()))
		require.NoError(t, nilPP.Shutdown(t.Context()))
	})

	t.Run("unregistered_duplicate_and_shutdown_paths", func(t *testing.T) {
		exp := &recordingExporter{}
		readerProd := &stubProducer{sps: sampleScopeProfiles("reader-prod")}
		pr := profile.NewPeriodicReader(exp, profile.WithInterval(time.Hour), profile.WithProducer(readerProd))

		var rp profiledata.ResourceProfiles
		assert.ErrorIs(t, pr.Collect(t.Context(), &rp), profile.ErrReaderNotRegistered)
		assert.Error(t, pr.Collect(t.Context(), nil))

		// Shutting down before registration shuts down exporter and reader producer.
		require.NoError(t, pr.Shutdown(t.Context()))
		assert.True(t, exp.shutdown)
		assert.True(t, readerProd.shutdown)
		assert.ErrorIs(t, pr.Collect(t.Context(), &rp), profile.ErrReaderShutdown)

		// Duplicate registration keeps the first provider active.
		res1 := resource.NewSchemaless(attribute.String("service.name", "first"))
		res2 := resource.NewSchemaless(attribute.String("service.name", "second"))
		pr2 := profile.NewPeriodicReader(&recordingExporter{}, profile.WithInterval(time.Hour))
		pp1 := profile.NewProfilerProvider(profile.WithResource(res1), profile.WithReader(pr2))
		_ = profile.NewProfilerProvider(profile.WithResource(res2), profile.WithReader(pr2))
		t.Cleanup(func() { _ = pp1.Shutdown(t.Context()) })

		require.NoError(t, pr2.Collect(t.Context(), &rp))
		assert.Equal(t, res1, rp.Resource)
	})

	t.Run("force_flush_errors_and_cancellation", func(t *testing.T) {
		exportErr := errors.New("export failed")
		flushErr := errors.New("flush failed")
		prodErr := errors.New("reader producer failed")

		exp := &recordingExporter{exportErr: exportErr}
		pr := profile.NewPeriodicReader(
			exp,
			profile.WithInterval(time.Hour),
			profile.WithProducer(&stubProducer{sps: sampleScopeProfiles("p"), err: prodErr}),
		)
		pp := profile.NewProfilerProvider(profile.WithReader(pr))
		t.Cleanup(func() { _ = pp.Shutdown(t.Context()) })

		err := pr.ForceFlush(t.Context())
		assert.ErrorIs(t, err, exportErr)
		assert.ErrorIs(t, err, prodErr)

		// When export succeeds, exporter.ForceFlush error is propagated.
		exp.mu.Lock()
		exp.exportErr = nil
		exp.flushErr = flushErr
		exp.mu.Unlock()
		assert.ErrorIs(t, pr.ForceFlush(t.Context()), prodErr)

		exp2 := &recordingExporter{flushErr: flushErr}
		pr2 := profile.NewPeriodicReader(exp2, profile.WithInterval(time.Hour))
		pp2 := profile.NewProfilerProvider(
			profile.WithProducer(&stubProducer{sps: sampleScopeProfiles("ok")}),
			profile.WithReader(pr2),
		)
		t.Cleanup(func() { _ = pp2.Shutdown(t.Context()) })
		assert.ErrorIs(t, pr2.ForceFlush(t.Context()), flushErr)

		// Blocked export honors ForceFlush context cancellation both while waiting
		// on errCh and while trying to send on flushCh.
		started := make(chan struct{})
		release := make(chan struct{})
		blockingExp := &recordingExporter{
			blockExport: func(ctx context.Context) error {
				select {
				case started <- struct{}{}:
				default:
				}
				select {
				case <-release:
					return nil
				case <-ctx.Done():
					return ctx.Err()
				}
			},
		}
		pr3 := profile.NewPeriodicReader(blockingExp, profile.WithInterval(time.Hour))
		pp3 := profile.NewProfilerProvider(
			profile.WithProducer(&stubProducer{sps: sampleScopeProfiles("block")}),
			profile.WithReader(pr3),
		)
		t.Cleanup(func() {
			close(release)
			_ = pp3.Shutdown(t.Context())
		})

		waitCtx, cancelWait := context.WithCancel(t.Context())
		errDone := make(chan error, 1)
		go func() { errDone <- pr3.ForceFlush(waitCtx) }()
		<-started

		// While the background loop is blocked in Export, a canceled ctx fails sending on flushCh.
		canceledCtx, cancelSend := context.WithCancel(t.Context())
		cancelSend()
		assert.ErrorIs(t, pr3.ForceFlush(canceledCtx), context.Canceled)

		// Canceling the waiting caller unblocks it from <-errCh.
		cancelWait()
		assert.ErrorIs(t, <-errDone, context.Canceled)
	})

	t.Run("periodic_tick_error_forwarded_to_otel_handler", func(t *testing.T) {
		wantErr := errors.New("tick export failed")
		gotErr := make(chan error, 1)
		prevHandler := otel.GetErrorHandler()
		otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
			select {
			case gotErr <- err:
			default:
			}
		}))
		t.Cleanup(func() { otel.SetErrorHandler(prevHandler) })

		exp := &recordingExporter{exportErr: wantErr}
		pr := profile.NewPeriodicReader(exp, profile.WithInterval(5*time.Millisecond))
		pp := profile.NewProfilerProvider(
			profile.WithProducer(&stubProducer{sps: sampleScopeProfiles("tick")}),
			profile.WithReader(pr),
		)
		t.Cleanup(func() { _ = pp.Shutdown(t.Context()) })

		select {
		case err := <-gotErr:
			assert.ErrorIs(t, err, wantErr)
		case <-time.After(2 * time.Second):
			t.Fatal("timed out waiting for periodic tick error")
		}
	})
}

func TestProfilerProviderConcurrentSafe(t *testing.T) {
	exp := &recordingExporter{}
	pr := profile.NewPeriodicReader(exp, profile.WithInterval(time.Hour))
	pp := profile.NewProfilerProvider(
		profile.WithProducer(&stubProducer{sps: sampleScopeProfiles("conc")}),
		profile.WithReader(pr),
	)

	var wg sync.WaitGroup
	for range 4 {
		wg.Go(func() {
			_ = pp.ForceFlush(t.Context())
		})
	}
	wg.Go(func() {
		_ = pp.Shutdown(t.Context())
	})
	wg.Wait()
	require.NoError(t, pp.Shutdown(t.Context()))
}

func TestManualReaderConcurrentSafe(t *testing.T) {
	mr := profile.NewManualReader(profile.WithProducer(&stubProducer{sps: sampleScopeProfiles("reader")}))
	pp := profile.NewProfilerProvider(
		profile.WithProducer(&stubProducer{sps: sampleScopeProfiles("provider")}),
		profile.WithReader(mr),
	)

	var wg sync.WaitGroup
	for range 4 {
		wg.Go(func() {
			var rp profiledata.ResourceProfiles
			_ = mr.Collect(t.Context(), &rp)
			_ = mr.ForceFlush(t.Context())
			_ = mr.MarshalLog()
		})
	}
	wg.Go(func() {
		_ = pp.Shutdown(t.Context())
	})
	wg.Wait()
	assert.ErrorIs(t, mr.Shutdown(t.Context()), profile.ErrReaderShutdown)
}

func TestPeriodicReaderConcurrentSafe(t *testing.T) {
	exp := &recordingExporter{}
	pr := profile.NewPeriodicReader(exp, profile.WithInterval(time.Hour))
	pp := profile.NewProfilerProvider(
		profile.WithProducer(&stubProducer{sps: sampleScopeProfiles("concurrent")}),
		profile.WithReader(pr),
	)

	var wg sync.WaitGroup
	var okCount atomic.Int32
	for range 6 {
		wg.Go(func() {
			var rp profiledata.ResourceProfiles
			if err := pr.Collect(t.Context(), &rp); err == nil && len(rp.ScopeProfiles) == 1 {
				okCount.Add(1)
			}
			_ = pr.ForceFlush(t.Context())
			_ = pr.MarshalLog()
		})
	}
	wg.Wait()
	assert.Equal(t, int32(6), okCount.Load())

	// Also race Collect/ForceFlush/MarshalLog with Shutdown.
	for range 4 {
		wg.Go(func() {
			var rp profiledata.ResourceProfiles
			_ = pr.Collect(t.Context(), &rp)
			_ = pr.ForceFlush(t.Context())
			_ = pr.MarshalLog()
		})
	}
	wg.Go(func() {
		_ = pp.Shutdown(t.Context())
	})
	wg.Wait()
	assert.ErrorIs(t, pr.Shutdown(t.Context()), profile.ErrReaderShutdown)
}
