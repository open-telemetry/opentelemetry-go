// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package profile

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/sdk/profile/profiledata"
)

// ErrExporterShutdown is returned if Export is called after an Exporter has
// been shut down.
var ErrExporterShutdown = errors.New("exporter is shutdown")

// Exporter handles the delivery of profile data to external receivers.
type Exporter interface {
	// Export serializes and transmits profile data to a receiver.
	//
	// Implementations must honor cancellation and deadlines on ctx and must
	// not retain rp after Export returns.
	Export(ctx context.Context, rp *profiledata.ResourceProfiles) error

	// ForceFlush flushes any profile data held by the exporter.
	//
	// This method must be safe to call concurrently.
	ForceFlush(ctx context.Context) error

	// Shutdown flushes any remaining profile data and releases resources held
	// by the exporter. Subsequent calls to Export should return
	// [ErrExporterShutdown].
	//
	// This method must be safe to call concurrently.
	Shutdown(ctx context.Context) error
}

type noopExporter struct{}

func (noopExporter) Export(context.Context, *profiledata.ResourceProfiles) error { return nil }
func (noopExporter) ForceFlush(context.Context) error                            { return nil }
func (noopExporter) Shutdown(context.Context) error                              { return nil }
