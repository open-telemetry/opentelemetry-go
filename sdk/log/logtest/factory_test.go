// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logtest

import (
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
)

func TestRecordFactoryEmpty(t *testing.T) {
	assert.Equal(t, sdklog.Record{}, RecordFactory{}.NewRecord())
}

func TestRecordFactory(t *testing.T) {
	now := time.Now()
	observed := now.Add(time.Second)
	eventName := "testing.name"
	severity := log.SeverityDebug
	severityText := "DBG"
	body := log.StringValue("Message")
	attrs := []log.KeyValue{
		log.Int("int", 1),
		log.String("str", "foo"),
		log.Float64("flt", 3.14),
	}
	traceID := trace.TraceID([16]byte{1})
	spanID := trace.SpanID([8]byte{2})
	traceFlags := trace.FlagsSampled
	dropped := 3
	scope := instrumentation.Scope{
		Name: t.Name(),
	}
	r := resource.NewSchemaless(attribute.Bool("works", true))

	got := RecordFactory{
		EventName:            eventName,
		Timestamp:            now,
		ObservedTimestamp:    observed,
		Severity:             severity,
		SeverityText:         severityText,
		Body:                 body,
		Attributes:           attrs,
		TraceID:              traceID,
		SpanID:               spanID,
		TraceFlags:           traceFlags,
		DroppedAttributes:    dropped,
		InstrumentationScope: &scope,
		Resource:             r,
	}.NewRecord()

	assert.Equal(t, eventName, got.EventName())
	assert.Equal(t, now, got.Timestamp())
	assert.Equal(t, observed, got.ObservedTimestamp())
	assert.Equal(t, severity, got.Severity())
	assert.Equal(t, severityText, got.SeverityText())
	assertBody(t, body, got)
	assertAttributes(t, attrs, got)
	assert.Equal(t, dropped, got.DroppedAttributes())
	assert.Equal(t, traceID, got.TraceID())
	assert.Equal(t, spanID, got.SpanID())
	assert.Equal(t, traceFlags, got.TraceFlags())
	assert.Equal(t, scope, got.InstrumentationScope())
	assert.Equal(t, r, got.Resource())
}

func TestRecordFactoryMultiple(t *testing.T) {
	now := time.Now()
	attrs := []log.KeyValue{
		log.Int("int", 1),
		log.String("str", "foo"),
		log.Float64("flt", 3.14),
	}
	scope := instrumentation.Scope{
		Name: t.Name(),
	}

	f := RecordFactory{
		Timestamp:            now,
		Attributes:           attrs,
		DroppedAttributes:    1,
		InstrumentationScope: &scope,
	}

	record1 := f.NewRecord()

	f.Attributes = append(f.Attributes, log.Bool("added", true))
	f.DroppedAttributes = 2
	record2 := f.NewRecord()

	assert.Equal(t, now, record2.Timestamp())
	assertAttributes(t, append(attrs, log.Bool("added", true)), record2)
	assert.Equal(t, 2, record2.DroppedAttributes())
	assert.Equal(t, scope, record2.InstrumentationScope())

	// Previously returned record is unharmed by the builder changes.
	assert.Equal(t, now, record1.Timestamp())
	assertAttributes(t, attrs, record1)
	assert.Equal(t, 1, record1.DroppedAttributes())
	assert.Equal(t, scope, record1.InstrumentationScope())
}

func assertBody(t *testing.T, want log.Value, r sdklog.Record) {
	t.Helper()
	got := r.Body()
	if !got.Equal(want) {
		t.Errorf("Body value is not equal:\nwant: %v\ngot:  %v", want, got)
	}
}

func assertAttributes(t *testing.T, want []log.KeyValue, r sdklog.Record) {
	t.Helper()
	var got []log.KeyValue
	r.WalkAttributes(func(kv log.KeyValue) bool {
		got = append(got, kv)
		return true
	})
	if !slices.EqualFunc(want, got, log.KeyValue.Equal) {
		t.Errorf("Attributes are not equal:\nwant: %v\ngot:  %v", want, got)
	}
}
