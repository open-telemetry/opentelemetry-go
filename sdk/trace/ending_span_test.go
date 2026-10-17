// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
)

// endedRecordingSpan starts and ends a span, returning the underlying
// *recordingSpan. The span is already ended when returned.
func endedRecordingSpan(t *testing.T) *recordingSpan {
	t.Helper()
	te := NewTestExporter()
	tp := NewTracerProvider(WithSyncer(te), WithResource(resource.Empty()))
	_, span := tp.Tracer("test").Start(t.Context(), "test")
	span.End()
	rs, ok := span.(*recordingSpan)
	require.True(t, ok)
	return rs
}

func TestEndingSpanIsRecording(t *testing.T) {
	rs := endedRecordingSpan(t)
	assert.False(t, rs.IsRecording(), "ended span should not be recording")
	assert.True(t, endingSpan{recordingSpan: rs}.IsRecording(), "endingSpan should always report recording")
}

func TestEndingSpanEndIsNoop(t *testing.T) {
	rs := endedRecordingSpan(t)
	endTime := rs.EndTime()
	endingSpan{recordingSpan: rs}.End(trace.WithTimestamp(time.Now().Add(time.Hour)))
	assert.Equal(t, endTime, rs.EndTime(), "endingSpan.End should not change end time")
}

func TestEndingSpanSetStatus(t *testing.T) {
	rs := endedRecordingSpan(t)
	// Direct mutation is a no-op because the span is ended.
	rs.SetStatus(codes.Error, "via rs")
	assert.Equal(t, codes.Unset, rs.Status().Code, "SetStatus on ended span should be no-op")

	// endingSpan bypasses that guard.
	endingSpan{recordingSpan: rs}.SetStatus(codes.Error, "via endingSpan")
	assert.Equal(t, codes.Error, rs.Status().Code)
	assert.Equal(t, "via endingSpan", rs.Status().Description)
}

func TestEndingSpanSetAttributes(t *testing.T) {
	rs := endedRecordingSpan(t)
	rs.SetAttributes(attribute.String("key", "via rs"))
	assert.Empty(t, rs.Attributes(), "SetAttributes on ended span should be no-op")

	endingSpan{recordingSpan: rs}.SetAttributes(attribute.String("key", "via endingSpan"))
	attrs := rs.Attributes()
	require.Len(t, attrs, 1)
	assert.Equal(t, "via endingSpan", attrs[0].Value.AsString())
}

func TestEndingSpanSetName(t *testing.T) {
	rs := endedRecordingSpan(t)
	rs.SetName("via rs")
	assert.Equal(t, "test", rs.Name(), "SetName on ended span should be no-op")

	endingSpan{recordingSpan: rs}.SetName("via endingSpan")
	assert.Equal(t, "via endingSpan", rs.Name())
}

func TestEndingSpanAddEvent(t *testing.T) {
	rs := endedRecordingSpan(t)
	rs.AddEvent("via rs")
	assert.Empty(t, rs.Events(), "AddEvent on ended span should be no-op")

	endingSpan{recordingSpan: rs}.AddEvent("via endingSpan")
	evts := rs.Events()
	require.Len(t, evts, 1)
	assert.Equal(t, "via endingSpan", evts[0].Name)
}

func TestEndingSpanAddLink(t *testing.T) {
	rs := endedRecordingSpan(t)
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: [16]byte{1},
		SpanID:  [8]byte{1},
	})
	rs.AddLink(trace.Link{SpanContext: sc})
	assert.Empty(t, rs.Links(), "AddLink on ended span should be no-op")

	endingSpan{recordingSpan: rs}.AddLink(trace.Link{SpanContext: sc})
	links := rs.Links()
	require.Len(t, links, 1)
	assert.Equal(t, sc, links[0].SpanContext)
}

func TestEndingSpanRecordError(t *testing.T) {
	rs := endedRecordingSpan(t)
	testErr := errors.New("test error")
	rs.RecordError(testErr)
	assert.Empty(t, rs.Events(), "RecordError on ended span should be no-op")

	endingSpan{recordingSpan: rs}.RecordError(testErr)
	evts := rs.Events()
	require.Len(t, evts, 1)
	assert.Equal(t, "exception", evts[0].Name)
}

func TestEndingSpanRecordErrorNil(t *testing.T) {
	rs := endedRecordingSpan(t)
	endingSpan{recordingSpan: rs}.RecordError(nil)
	assert.Empty(t, rs.Events())
}
