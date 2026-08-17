// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlploghttp

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	logpb "go.opentelemetry.io/proto/otlp/logs/v1"

	"go.opentelemetry.io/otel/sdk/log"
)

func TestExporterExportErrors(t *testing.T) {
	errUpload := errors.New("upload")
	c := &client{
		uploadLogs: func(context.Context, []*logpb.ResourceLogs) error {
			return errUpload
		},
	}

	e, err := newExporter(c, config{})
	require.NoError(t, err, "New")

	err = e.Export(t.Context(), make([]log.Record, 1))
	assert.ErrorIs(t, err, errUpload)
}

func TestExporterExport(t *testing.T) {
	var uploads int
	c := &client{
		uploadLogs: func(context.Context, []*logpb.ResourceLogs) error {
			uploads++
			return nil
		},
	}

	orig := transformResourceLogs
	var got []log.Record
	transformResourceLogs = func(r []log.Record) []*logpb.ResourceLogs {
		got = r
		return make([]*logpb.ResourceLogs, 1)
	}
	t.Cleanup(func() { transformResourceLogs = orig })

	e, err := newExporter(c, config{})
	require.NoError(t, err, "New")

	ctx := t.Context()
	want := make([]log.Record, 1)
	assert.NoError(t, e.Export(ctx, want))

	assert.Equal(t, 1, uploads, "client UploadLogs calls")
	assert.Equal(t, want, got, "transformed log records")
}

func TestExporterShutdown(t *testing.T) {
	ctx := t.Context()
	e, err := New(ctx)
	require.NoError(t, err, "New")
	assert.NoError(t, e.Shutdown(ctx), "Shutdown Exporter")

	// After Shutdown is called, calls to Export, Shutdown, or ForceFlush
	// should perform no operation and return nil error.
	r := make([]log.Record, 1)
	assert.NoError(t, e.Export(ctx, r), "Export on Shutdown Exporter")
	assert.NoError(t, e.ForceFlush(ctx), "ForceFlush on Shutdown Exporter")
	assert.NoError(t, e.Shutdown(ctx), "Shutdown on Shutdown Exporter")
}

func TestExporterForceFlush(t *testing.T) {
	ctx := t.Context()
	e, err := New(ctx)
	require.NoError(t, err, "New")

	assert.NoError(t, e.ForceFlush(ctx), "ForceFlush")
}

func BenchmarkExporterExport(b *testing.B) {
	c := &client{uploadLogs: func(context.Context, []*logpb.ResourceLogs) error { return nil }}
	e, err := newExporter(c, config{})
	require.NoError(b, err)
	records := make([]log.Record, 1)

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		err = e.Export(b.Context(), records)
	}
	_ = err
}
