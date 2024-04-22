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
)

func TestRecordBuilder(t *testing.T) {
	now := time.Now()
	observed := now.Add(time.Second)
	severity := log.SeverityDebug
	severityText := "DBG"
	body := log.StringValue("Message")
	attrs := []log.KeyValue{
		log.Int("int", 1),
		log.String("str", "foo"),
		log.Float64("flt", 3.14),
	}
	dropped := 3
	scope := instrumentation.Scope{
		Name: t.Name(),
	}
	r := resource.NewSchemaless(attribute.Bool("works", true))

	b := RecordBuilder{}.
		SetTimestamp(now).
		SetObservedTimestamp(observed).
		SetSeverity(severity).
		SetSeverityText(severityText).
		SetBody(body).
		SetAttributes(attrs...).
		SetDroppedAttributes(dropped).
		SetInstrumentationScope(scope).
		SetResource(r)
	got := b.Record()

	assert.Equal(t, now, got.Timestamp())
	assert.Equal(t, observed, got.ObservedTimestamp())
	assert.Equal(t, severity, got.Severity())
	assert.Equal(t, severityText, got.SeverityText())
	assertBody(t, body, got)
	assertAttributes(t, attrs, got)
	assert.Equal(t, dropped, got.DroppedAttributes())
	assert.Equal(t, scope, got.InstrumentationScope())
	assert.Equal(t, *r, got.Resource())
}

func TestRecordBuilderMultiple(t *testing.T) {
	now := time.Now()
	attrs := []log.KeyValue{
		log.Int("int", 1),
		log.String("str", "foo"),
		log.Float64("flt", 3.14),
	}
	scope := instrumentation.Scope{
		Name: t.Name(),
	}

	b := RecordBuilder{}.
		SetTimestamp(now).
		AddAttributes(attrs...).
		SetDroppedAttributes(1).
		SetInstrumentationScope(scope)

	record1 := b.Record()

	record2 := b.
		AddAttributes(log.Bool("added", true)).
		SetDroppedAttributes(2).
		Record()

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
