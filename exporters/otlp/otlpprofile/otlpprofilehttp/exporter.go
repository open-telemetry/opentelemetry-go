// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpprofilehttp

import (
	"context"
	"sync/atomic"

	"go.opentelemetry.io/otel/exporters/otlp/otlpprofile/otlpprofilehttp/internal/transform"
	"go.opentelemetry.io/otel/sdk/profile"
	"go.opentelemetry.io/otel/sdk/profile/profiledata"
)

// Exporter is an OpenTelemetry profile exporter that transports OTLP Profiles
// v1development protobuf payloads over HTTP.
// Exporter must be created with [New].
type Exporter struct {
	client  atomic.Pointer[client]
	stopped atomic.Bool
}

var _ profile.Exporter = (*Exporter)(nil)

// New returns a new [Exporter].
func New(ctx context.Context, options ...Option) (*Exporter, error) {
	cfg := newConfig(options)
	c, err := newHTTPClient(ctx, cfg)
	if err != nil {
		return nil, err
	}
	e := &Exporter{}
	e.client.Store(c)
	return e, nil
}

// Export transforms and transmits rp to an OTLP receiver. It returns
// [profile.ErrExporterShutdown] if called after Shutdown.
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

	c := e.client.Load()
	if e.stopped.Load() {
		return profile.ErrExporterShutdown
	}
	return c.UploadProfiles(ctx, req)
}

// Shutdown shuts down the Exporter. Subsequent calls to Export return
// [profile.ErrExporterShutdown]. Calling Shutdown more than once is safe.
//
// This method is safe to call concurrently.
func (e *Exporter) Shutdown(context.Context) error {
	if e.stopped.Swap(true) {
		return nil
	}
	if prev := e.client.Swap(newNoopClient()); prev != nil && prev.stop != nil {
		prev.stop()
	}
	return nil
}

// ForceFlush does nothing because the Exporter holds no buffered state.
//
// This method is safe to call concurrently.
func (*Exporter) ForceFlush(context.Context) error {
	return nil
}
