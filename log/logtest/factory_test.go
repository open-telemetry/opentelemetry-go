// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logtest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"

	"go.opentelemetry.io/otel/log"
)

func TestRecordFactory(t *testing.T) {
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

	got := RecordFactory{
		Timestamp:         now,
		ObservedTimestamp: observed,
		Severity:          severity,
		SeverityText:      severityText,
		Body:              body,
		Attributes:        attrs,
	}.NewRecord()

	assert.Equal(t, now, got.Timestamp())
	assert.Equal(t, observed, got.ObservedTimestamp())
	assert.Equal(t, severity, got.Severity())
	assert.Equal(t, severityText, got.SeverityText())
	assert.Equal(t, body, got.Body())
	assertAttributes(t, attrs, got)
}

func TestRecordFactoryMultiple(t *testing.T) {
	now := time.Now()
	attrs := []log.KeyValue{
		log.Int("int", 1),
		log.String("str", "foo"),
		log.Float64("flt", 3.14),
	}

	f := RecordFactory{
		Timestamp:  now,
		Attributes: attrs,
	}

	record1 := f.NewRecord()
	f.Attributes = append(f.Attributes, log.Bool("added", true))

	record2 := f.NewRecord()
	assert.Equal(t, now, record2.Timestamp())
	assertAttributes(t, append(attrs, log.Bool("added", true)), record2)

	// Previously returned record is unharmed by the builder changes.
	assert.Equal(t, now, record1.Timestamp())
	assertAttributes(t, attrs, record1)
}

func assertAttributes(t *testing.T, want []log.KeyValue, r log.Record) {
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
