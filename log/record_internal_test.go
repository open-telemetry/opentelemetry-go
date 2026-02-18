// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log

import "testing"

type namedErr struct{}

func (namedErr) Error() string { return "named" }

func TestTypeStrBuiltin(t *testing.T) {
	got := typeStr(struct{}{})
	if got != "struct {}" {
		t.Fatalf("typeStr(struct{}{}) = %q, want %q", got, "struct {}")
	}
}

func TestTypeStrNamed(t *testing.T) {
	got := typeStr(namedErr{})
	want := "go.opentelemetry.io/otel/log.namedErr"
	if got != want {
		t.Fatalf("typeStr(namedErr{}) = %q, want %q", got, want)
	}
}
