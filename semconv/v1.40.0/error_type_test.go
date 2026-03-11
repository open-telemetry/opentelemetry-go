// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package semconv // import "go.opentelemetry.io/otel/semconv/v1.40.0"

import (
	"errors"
	"testing"
)

const pkg = "go.opentelemetry.io/otel/semconv/v1.40.0"

func TestErrorType(t *testing.T) {
	check(t, nil, ErrorTypeOther.Value.AsString())
	check(t, errors.New("msg"), "*errors.errorString")
	check(t, custom("aborted"), "aborted")
	check(t, errors.Join(custom("left"), custom("right")), "left") // first errors.As match is used.
	check(t, custom(""), pkg+".ErrCustomType")                     // empty ErrorType, use concrete type.
	check(t, wrapped(custom("wrapped-aborted")), "wrapped-aborted")
	check(t, wrapped(custom("")), pkg+".wrappedErr") // empty ErrorType in chain, use concrete top-level type.
}

func check(t *testing.T, err error, want string) {
	t.Helper()
	got := ErrorType(err)
	if got.Key != ErrorTypeKey {
		t.Errorf("ErrorType(%v) key = %v, want %v", err, got.Key, ErrorTypeKey)
	}
	if got.Value.AsString() != want {
		t.Errorf("ErrorType(%v) value = %v, want %v", err, got.Value.AsString(), want)
	}
}

func custom(typ string) error {
	return ErrCustomType{Type: typ}
}

func wrapped(err error) error {
	return wrappedErr{err: err}
}

type ErrCustomType struct {
	Type string
}

func (e ErrCustomType) Error() string {
	return "custom: " + e.Type
}

func (e ErrCustomType) ErrorType() string {
	return e.Type
}

type wrappedErr struct {
	err error
}

func (e wrappedErr) Error() string {
	return "wrapped: " + e.err.Error()
}

func (e wrappedErr) Unwrap() error {
	return e.err
}
