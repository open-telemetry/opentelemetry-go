// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpprofilegrpc

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	colprofilepb "go.opentelemetry.io/proto/otlp/collector/profiles/v1development"

	"go.opentelemetry.io/otel/sdk/profile"
	"go.opentelemetry.io/otel/sdk/profile/profiledata"
)

type stubProfileClient struct {
	uploadFn   func(context.Context, *colprofilepb.ExportProfilesServiceRequest) error
	shutdownFn func(context.Context) error
	uploads    atomic.Int32
}

func (s *stubProfileClient) UploadProfiles(
	ctx context.Context,
	req *colprofilepb.ExportProfilesServiceRequest,
) error {
	s.uploads.Add(1)
	if s.uploadFn != nil {
		return s.uploadFn(ctx, req)
	}
	return nil
}

func (s *stubProfileClient) Shutdown(ctx context.Context) error {
	if s.shutdownFn != nil {
		return s.shutdownFn(ctx)
	}
	return nil
}

func nonEmptyResourceProfiles() *profiledata.ResourceProfiles {
	return &profiledata.ResourceProfiles{
		ScopeProfiles: []profiledata.ScopeProfiles{
			{
				Profiles: []profiledata.Profile{
					{
						SampleType: profiledata.ValueType{Type: "samples", Unit: "count"},
						Samples:    []profiledata.Sample{{Values: []int64{1}}},
					},
				},
			},
		},
	}
}

func TestExporterConcurrentSafe(t *testing.T) {
	e := &Exporter{client: &stubProfileClient{}}

	const goroutines = 10

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(t.Context())
	var runs atomic.Uint64
	rp := nonEmptyResourceProfiles()
	for range goroutines {
		wg.Go(func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					_ = e.Export(ctx, rp)
					_ = e.ForceFlush(ctx)
					runs.Add(1)
				}
			}
		})
	}

	for runs.Load() == 0 {
		runtime.Gosched()
	}

	assert.NoError(t, e.Shutdown(ctx))
	cancel()
	wg.Wait()
}

func TestExporterExportConcurrentShutdown(t *testing.T) {
	orig := transformExportRequest
	started := make(chan struct{})
	resume := make(chan struct{})
	transformExportRequest = func(*profiledata.ResourceProfiles) *colprofilepb.ExportProfilesServiceRequest {
		close(started)
		<-resume
		return &colprofilepb.ExportProfilesServiceRequest{}
	}
	t.Cleanup(func() { transformExportRequest = orig })

	stub := &stubProfileClient{}
	e := &Exporter{client: stub}
	ctx := t.Context()

	errCh := make(chan error, 1)
	go func() {
		errCh <- e.Export(ctx, nonEmptyResourceProfiles())
	}()

	<-started
	require.NoError(t, e.Shutdown(ctx))
	close(resume)

	assert.ErrorIs(t, <-errCh, profile.ErrExporterShutdown)
	assert.Zero(t, stub.uploads.Load())
	assert.NoError(t, e.ForceFlush(ctx))
	assert.ErrorIs(t, e.Export(ctx, nil), profile.ErrExporterShutdown)
}

func TestExporterConcurrentShutdownWaitsForTeardown(t *testing.T) {
	inShutdown := make(chan struct{})
	releaseShutdown := make(chan struct{})
	var finished atomic.Bool

	stub := &stubProfileClient{
		shutdownFn: func(context.Context) error {
			close(inShutdown)
			<-releaseShutdown
			finished.Store(true)
			return nil
		},
	}
	e := &Exporter{client: stub}

	var wg sync.WaitGroup
	wg.Go(func() {
		assert.NoError(t, e.Shutdown(t.Context()))
	})

	<-inShutdown

	secondDone := make(chan struct{})
	go func() {
		assert.NoError(t, e.Shutdown(t.Context()))
		assert.True(t, finished.Load(), "second Shutdown returned before client teardown completed")
		close(secondDone)
	}()

	select {
	case <-secondDone:
		t.Fatal("second Shutdown returned while first Shutdown was still running")
	default:
	}

	close(releaseShutdown)
	<-secondDone
	wg.Wait()
}
