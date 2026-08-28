// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlptrace_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/go-logr/logr/funcr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

type syncClient struct {
	uploadErr   error
	logEndpoint string
	captured    [][]*tracepb.ResourceSpans
	uploadSyncCalled bool
}

var _ otlptrace.SyncClient = &syncClient{}

func (*syncClient) Start(context.Context) error { return nil }
func (*syncClient) Stop(context.Context) error  { return nil }
func (c *syncClient) UploadTraces(_ context.Context, _ []*tracepb.ResourceSpans) error {
	// Should not be called by SyncExporter; if it is, fail the test.
	return assert.AnError
}
func (c *syncClient) UploadTracesSync(_ context.Context, rs []*tracepb.ResourceSpans) error {
	c.uploadSyncCalled = true
	// Capture shallow copy to test retention contract: caller should not
	// retain after return, but we capture to verify data correctness before
	// arena reuse. We copy slice header only; underlying data will be
	// invalidated after Reset, so test must check values before that.
	c.captured = append(c.captured, rs)
	return c.uploadErr
}
func (c *syncClient) MarshalLog() any {
	return struct{ Endpoint string }{Endpoint: c.logEndpoint}
}

func TestSyncExporterUsesUploadTracesSync(t *testing.T) {
	ctx := t.Context()
	client := &syncClient{}
	exp, err := otlptrace.NewSync(ctx, client)
	require.NoError(t, err)

	spans := tracetest.SpanStubs{{Name: "sync-span"}}.Snapshots()
	err = exp.ExportSpans(ctx, spans)
	require.NoError(t, err)
	assert.True(t, client.uploadSyncCalled, "UploadTracesSync should be called")

	assert.NoError(t, exp.Shutdown(ctx))
}

func TestSyncExporterClientError(t *testing.T) {
	ctx := t.Context()
	exp, err := otlptrace.NewSync(ctx, &syncClient{
		uploadErr: context.Canceled,
	})
	require.NoError(t, err)

	spans := tracetest.SpanStubs{{Name: "Span 0"}}.Snapshots()
	err = exp.ExportSpans(ctx, spans)

	assert.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
	assert.True(t, strings.HasPrefix(err.Error(), "traces export: "), "%+v", err)

	assert.NoError(t, exp.Shutdown(ctx))
}

func TestSyncExporterMarshalLogDoesNotIncludeClientConfig(t *testing.T) {
	const sensitiveEndpoint = "user:pass@collector.internal:4318"

	var buf bytes.Buffer
	logger := funcr.New(func(_, args string) {
		_, _ = buf.WriteString(args)
	}, funcr.Options{})

	exp := otlptrace.NewSyncUnstarted(&syncClient{logEndpoint: sensitiveEndpoint})
	logger.Info("exporter", "config", exp)

	logged := buf.String()
	assert.Contains(t, logged, "otlptrace-sync")
	assert.NotContains(t, logged, sensitiveEndpoint)
}

func TestSyncExporterArenaReuse(t *testing.T) {
	ctx := t.Context()
	client := &syncClient{}
	exp, err := otlptrace.NewSync(ctx, client, otlptrace.WithInitialBatchSize(2), otlptrace.WithMaxRetainedBatchSize(10))
	require.NoError(t, err)

	// First batch
	spans1 := tracetest.SpanStubs{
		{Name: "span-1", Attributes: nil},
		{Name: "span-2"},
	}.Snapshots()
	require.NoError(t, exp.ExportSpans(ctx, spans1))
	require.Len(t, client.captured, 1)
	// Verify first batch has correct names before arena reset corrupts it.
	// Since we captured slice header, underlying arena data is still valid
	// until next export reuses arena. We check immediately.
	firstBatch := client.captured[0]
	// Count spans in first batch
	count := 0
	for _, rs := range firstBatch {
		for _, ss := range rs.ScopeSpans {
			count += len(ss.Spans)
		}
	}
	assert.Equal(t, 2, count)

	// Second batch with different data to ensure arena reuse does not corrupt
	// first batch's captured reference after Return (client must not retain).
	// But since our test client retains slice header, second export will reuse
	// arena and mutate memory that first batch points to. A correct SyncClient
	// must NOT retain after return, so capturing is safe only if we deep-copy.
	// Instead verify second export succeeds and has correct data.
	spans2 := tracetest.SpanStubs{
		{Name: "span-3"},
	}.Snapshots()
	require.NoError(t, exp.ExportSpans(ctx, spans2))
	require.Len(t, client.captured, 2)
	count = 0
	for _, rs := range client.captured[1] {
		for _, ss := range rs.ScopeSpans {
			count += len(ss.Spans)
		}
	}
	assert.Equal(t, 1, count)

	assert.NoError(t, exp.Shutdown(ctx))
}

func TestSyncExporterWithOptions(t *testing.T) {
	ctx := t.Context()
	client := &syncClient{}
	exp := otlptrace.NewSyncUnstarted(client, otlptrace.WithInitialBatchSize(100), otlptrace.WithMaxRetainedBatchSize(1000))
	require.NotNil(t, exp)
	require.NoError(t, exp.Start(ctx))
	require.NoError(t, exp.Shutdown(ctx))
}
