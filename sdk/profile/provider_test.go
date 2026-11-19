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

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/profile"
	"go.opentelemetry.io/otel/sdk/profile/profiledata"
	"go.opentelemetry.io/otel/sdk/resource"
)

type stubProducer struct {
	mu       sync.Mutex
	sps      []profiledata.ScopeProfiles
	err      error
	calls    int
	shutdown bool
	onShut   func()
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
	p.mu.Unlock()
	if fn != nil {
		fn()
	}
	return nil
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
}

func (e *recordingExporter) Export(ctx context.Context, rp *profiledata.ResourceProfiles) error {
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

		var rp profiledata.ResourceProfiles
		require.NoError(t, mr.Collect(t.Context(), &rp))
		assert.Equal(t, res, rp.Resource)
		require.Len(t, rp.ScopeProfiles, 2)
		assert.Equal(t, "from-provider", rp.ScopeProfiles[0].Scope.Name)
		assert.Equal(t, "from-reader", rp.ScopeProfiles[1].Scope.Name)

		require.NoError(t, mr.ForceFlush(t.Context()))
		assert.NotNil(t, mr.MarshalLog())

		// Canceled context returns context error.
		canceledCtx, cancel := context.WithCancel(t.Context())
		cancel()
		assert.ErrorIs(t, mr.Collect(canceledCtx, &rp), context.Canceled)

		require.NoError(t, pp.Shutdown(t.Context()))
		assert.ErrorIs(t, mr.Collect(t.Context(), &rp), profile.ErrReaderShutdown)
		assert.ErrorIs(t, mr.ForceFlush(t.Context()), profile.ErrReaderShutdown)
		assert.ErrorIs(t, mr.Shutdown(t.Context()), profile.ErrReaderShutdown)
	})

	t.Run("producer_error_joined", func(t *testing.T) {
		wantErr := errors.New("boom")
		mr := profile.NewManualReader(profile.WithProducer(&stubProducer{err: wantErr}))
		pp := profile.NewProfilerProvider(profile.WithReader(mr))
		t.Cleanup(func() { _ = pp.Shutdown(t.Context()) })

		var rp profiledata.ResourceProfiles
		err := mr.Collect(t.Context(), &rp)
		assert.ErrorIs(t, err, wantErr)
	})
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

		assert.NotNil(t, pr.MarshalLog())
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

		// Nil exporter must not panic.
		nilPR := profile.NewPeriodicReader(nil, profile.WithInterval(time.Hour))
		nilPP := profile.NewProfilerProvider(profile.WithReader(nilPR))
		require.NoError(t, nilPP.ForceFlush(t.Context()))
		require.NoError(t, nilPP.Shutdown(t.Context()))
	})

	t.Run("collect_direct_and_concurrent", func(t *testing.T) {
		exp := &recordingExporter{}
		prod := &stubProducer{sps: sampleScopeProfiles("concurrent")}
		pr := profile.NewPeriodicReader(exp, profile.WithInterval(time.Hour))
		pp := profile.NewProfilerProvider(profile.WithProducer(prod), profile.WithReader(pr))
		t.Cleanup(func() { _ = pp.Shutdown(t.Context()) })

		assert.Error(t, pr.Collect(t.Context(), nil))

		var wg sync.WaitGroup
		var okCount atomic.Int32
		for range 8 {
			wg.Go(func() {
				var rp profiledata.ResourceProfiles
				if err := pr.Collect(t.Context(), &rp); err == nil && len(rp.ScopeProfiles) == 1 {
					okCount.Add(1)
				}
			})
		}
		wg.Wait()
		assert.Equal(t, int32(8), okCount.Load())
	})
}
