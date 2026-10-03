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

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

type syncClient struct {
	uploadErr        error
	logEndpoint      string
	captured         [][]*tracepb.ResourceSpans
	uploadSyncCalled bool
	onUpload         func([]*tracepb.ResourceSpans)
}

var _ otlptrace.SyncClient = &syncClient{}

func (*syncClient) Start(context.Context) error { return nil }
func (*syncClient) Stop(context.Context) error  { return nil }
func (*syncClient) UploadTraces(_ context.Context, _ []*tracepb.ResourceSpans) error {
	// Should not be called by SyncExporter; if it is, fail the test.
	return assert.AnError
}

func (c *syncClient) UploadTracesSync(_ context.Context, rs []*tracepb.ResourceSpans) error {
	c.uploadSyncCalled = true
	// Verify synchronously before return: arena is Reset after Upload returns,
	// so arena-owned attributes/events must be checked here, not after.
	if c.onUpload != nil {
		c.onUpload(rs)
	}
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
	exp, err := otlptrace.NewSync(
		ctx,
		client,
		otlptrace.WithInitialBatchSize(2),
		otlptrace.WithMaxRetainedBatchSize(10),
	)
	require.NoError(t, err)

	// First batch with arena-owned data (attributes + events are arena-backed,
	// span structs and names are heap-allocated). Values must be checked
	// synchronously in onUpload: ExportSpans Resets the arena after Upload
	// returns, so captured references are invalid afterwards.
	client.onUpload = func(rs []*tracepb.ResourceSpans) {
		assert.Equal(t, "v1", spanAttrString(rs, "span-1", "k"))
		assert.Equal(t, int64(7), spanAttrInt(rs, "span-2", "n"))
		assert.Equal(t, "ev1", spanEventAttrString(rs, "span-1", "e1", "ek"))
	}
	spans1 := tracetest.SpanStubs{
		{
			Name:       "span-1",
			Attributes: []attribute.KeyValue{attribute.String("k", "v1")},
			Events: []tracesdk.Event{
				{Name: "e1", Attributes: []attribute.KeyValue{attribute.String("ek", "ev1")}},
			},
		},
		{Name: "span-2", Attributes: []attribute.KeyValue{attribute.Int("n", 7)}},
	}.Snapshots()
	require.NoError(t, exp.ExportSpans(ctx, spans1))
	require.Len(t, client.captured, 1)

	// Second batch reuses the arena with different values.
	client.onUpload = func(rs []*tracepb.ResourceSpans) {
		assert.Equal(t, "v2", spanAttrString(rs, "span-3", "k"))
		assert.Equal(t, "ev2", spanEventAttrString(rs, "span-3", "e1", "ek"))
	}
	spans2 := tracetest.SpanStubs{
		{
			Name:       "span-3",
			Attributes: []attribute.KeyValue{attribute.String("k", "v2")},
			Events: []tracesdk.Event{
				{Name: "e1", Attributes: []attribute.KeyValue{attribute.String("ek", "ev2")}},
			},
		},
	}.Snapshots()
	require.NoError(t, exp.ExportSpans(ctx, spans2))
	require.Len(t, client.captured, 2)

	assert.NoError(t, exp.Shutdown(ctx))
}

func findPBSpan(
	rs []*tracepb.ResourceSpans,
	name string,
) *tracepb.Span {
	for _, r := range rs {
		for _, ss := range r.ScopeSpans {
			for _, s := range ss.Spans {
				if s.Name == name {
					return s
				}
			}
		}
	}
	return nil
}

func spanAttrString(rs []*tracepb.ResourceSpans, spanName, key string) string {
	s := findPBSpan(rs, spanName)
	if s == nil {
		return ""
	}
	for _, a := range s.Attributes {
		if a.Key == key {
			return a.Value.GetStringValue()
		}
	}
	return ""
}

func spanAttrInt(rs []*tracepb.ResourceSpans, spanName, key string) int64 {
	s := findPBSpan(rs, spanName)
	if s == nil {
		return 0
	}
	for _, a := range s.Attributes {
		if a.Key == key {
			return a.Value.GetIntValue()
		}
	}
	return 0
}

func spanEventAttrString(rs []*tracepb.ResourceSpans, spanName, eventName, key string) string {
	s := findPBSpan(rs, spanName)
	if s == nil {
		return ""
	}
	for _, e := range s.Events {
		if e.Name == eventName {
			for _, a := range e.Attributes {
				if a.Key == key {
					return a.Value.GetStringValue()
				}
			}
		}
	}
	return ""
}

func TestSyncExporterWithOptions(t *testing.T) {
	ctx := t.Context()
	client := &syncClient{}
	exp := otlptrace.NewSyncUnstarted(
		client,
		otlptrace.WithInitialBatchSize(100),
		otlptrace.WithMaxRetainedBatchSize(1000),
	)
	require.NotNil(t, exp)
	require.NoError(t, exp.Start(ctx))
	require.NoError(t, exp.Shutdown(ctx))
}
