// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpprofilegrpc

import (
	"context"
	"sync"
	"sync/atomic"

	colprofilepb "go.opentelemetry.io/proto/otlp/collector/profiles/v1development"

	"go.opentelemetry.io/otel/exporters/otlp/otlpprofile/otlpprofilegrpc/internal/transform"
	"go.opentelemetry.io/otel/sdk/profile"
	"go.opentelemetry.io/otel/sdk/profile/profiledata"
)

type profileClient interface {
	UploadProfiles(ctx context.Context, req *colprofilepb.ExportProfilesServiceRequest) error
	Shutdown(context.Context) error
}

// Exporter is an OpenTelemetry profile exporter that transports OTLP Profiles
// v1development protobuf payloads over gRPC.
// All Exporter values must be created with [New].
type Exporter struct {
	clientMu sync.Mutex
	client   profileClient
	stopped  atomic.Bool
}

var _ profile.Exporter = (*Exporter)(nil)

// New returns a new [Exporter].
func New(_ context.Context, options ...Option) (*Exporter, error) {
	cfg := newConfig(options)
	c, err := newClient(cfg)
	if err != nil {
		return nil, err
	}
	return &Exporter{client: c}, nil
}

// Export transforms and transmits rp to an OTLP receiver.
// Returns [profile.ErrExporterShutdown] if called after Shutdown.
//
// This method is safe to call concurrently.
func (e *Exporter) Export(ctx context.Context, rp *profiledata.ResourceProfiles) error {
	if e.stopped.Load() {
		return profile.ErrExporterShutdown
	}
	req := transform.ExportRequest(rp)
	if req == nil {
		return nil
	}

	e.clientMu.Lock()
	defer e.clientMu.Unlock()

	if e.stopped.Load() {
		return profile.ErrExporterShutdown
	}
	return e.client.UploadProfiles(ctx, req)
}

// Shutdown shuts down the Exporter and closes any gRPC connection it created.
// Subsequent calls to Export return [profile.ErrExporterShutdown]. Calling
// Shutdown more than once is safe.
//
// This method is safe to call concurrently.
func (e *Exporter) Shutdown(ctx context.Context) error {
	if e.stopped.Swap(true) {
		return nil
	}

	e.clientMu.Lock()
	defer e.clientMu.Unlock()

	err := e.client.Shutdown(ctx)
	e.client = newNoopClient()
	return err
}

// ForceFlush does nothing because the Exporter holds no buffered state.
//
// This method is safe to call concurrently.
func (*Exporter) ForceFlush(context.Context) error {
	return nil
}
