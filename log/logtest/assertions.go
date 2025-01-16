// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logtest // import "go.opentelemetry.io/otel/log/logtest"

import (
	"slices"
	"testing"

	"go.opentelemetry.io/otel/cmplxattr"
	"go.opentelemetry.io/otel/log"
)

// AssertRecordEqual compares two log records, and fails the test if they are
// not equal.
func AssertRecordEqual(t testing.TB, want, got log.Record) bool {
	t.Helper()

	if !want.Timestamp().Equal(got.Timestamp()) {
		t.Errorf("Timestamp value is not equal:\nwant: %v\ngot:  %v", want.Timestamp(), got.Timestamp())
		return false
	}
	if !want.ObservedTimestamp().Equal(got.ObservedTimestamp()) {
		t.Errorf("ObservedTimestamp value is not equal:\nwant: %v\ngot:  %v", want.ObservedTimestamp(), got.ObservedTimestamp())
		return false
	}
	if want.Severity() != got.Severity() {
		t.Errorf("Severity value is not equal:\nwant: %v\ngot:  %v", want.Severity(), got.Severity())
		return false
	}
	if want.SeverityText() != got.SeverityText() {
		t.Errorf("SeverityText value is not equal:\nwant: %v\ngot:  %v", want.SeverityText(), got.SeverityText())
		return false
	}
	if !assertBody(t, want.Body(), got) {
		return false
	}

	var attrs []cmplxattr.KeyValue
	want.WalkAttributes(func(kv cmplxattr.KeyValue) bool {
		attrs = append(attrs, kv)
		return true
	})
	return assertAttributes(t, attrs, got)
}

func assertBody(t testing.TB, want cmplxattr.Value, r log.Record) bool {
	t.Helper()
	got := r.Body()
	if !got.Equal(want) {
		t.Errorf("Body value is not equal:\nwant: %v\ngot:  %v", want, got)
		return false
	}

	return true
}

func assertAttributes(t testing.TB, want []cmplxattr.KeyValue, r log.Record) bool {
	t.Helper()
	var got []cmplxattr.KeyValue
	r.WalkAttributes(func(kv cmplxattr.KeyValue) bool {
		got = append(got, kv)
		return true
	})
	if !slices.EqualFunc(want, got, cmplxattr.KeyValue.Equal) {
		t.Errorf("Attributes are not equal:\nwant: %v\ngot:  %v", want, got)
		return false
	}

	return true
}
