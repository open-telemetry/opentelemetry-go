// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logtest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/cmplxattr"
	"go.opentelemetry.io/otel/log"
)

func TestRecordFactory(t *testing.T) {
	now := time.Now()
	observed := now.Add(time.Second)
	severity := log.SeverityDebug
	severityText := "DBG"
	body := cmplxattr.StringValue("Message")
	attrs := []cmplxattr.KeyValue{
		cmplxattr.Int("int", 1),
		cmplxattr.String("str", "foo"),
		cmplxattr.Float64("flt", 3.14),
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
	assertBody(t, body, got)
	assertAttributes(t, attrs, got)
}

func TestRecordFactoryMultiple(t *testing.T) {
	now := time.Now()
	attrs := []cmplxattr.KeyValue{
		cmplxattr.Int("int", 1),
		cmplxattr.String("str", "foo"),
		cmplxattr.Float64("flt", 3.14),
	}

	f := RecordFactory{
		Timestamp:  now,
		Attributes: attrs,
	}

	record1 := f.NewRecord()
	f.Attributes = append(f.Attributes, cmplxattr.Bool("added", true))

	record2 := f.NewRecord()
	assert.Equal(t, now, record2.Timestamp())
	assertAttributes(t, append(attrs, cmplxattr.Bool("added", true)), record2)

	// Previously returned record is unharmed by the builder changes.
	assert.Equal(t, now, record1.Timestamp())
	assertAttributes(t, attrs, record1)
}
