// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logtest // import "go.opentelemetry.io/otel/log/logtest"

import (
	"context"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"go.opentelemetry.io/otel/log"
)

// TestingT reports failure messages.
// *testing.T implements this interface.
type TestingT interface {
	Errorf(format string, args ...any)
}

// AssertEqual asserts that the two concrete data-types from the logtest package are equal.
func AssertEqual[T Recording | Record](t TestingT, want, got T, opts ...AssertOption) bool {
	if h, ok := t.(interface{ Helper() }); ok {
		h.Helper()
	}

	var cfg assertConfig
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}

	cmpOpts := []cmp.Option{
		cmp.Comparer(func(x, y context.Context) bool { return x == y }), // Compare context.
		cmpopts.SortSlices(
			func(a, b log.KeyValue) bool { return a.Key < b.Key },
		), // Unordered compare of the key values.
		cmpopts.EquateEmpty(), // Empty and nil collections are equal.
	}
	cmpOpts = append(cmpOpts, cfg.cmpOpts...)

	if diff := cmp.Diff(want, got, cmpOpts...); diff != "" {
		msg := "mismatch (-want +got):\n%s"
		if cfg.msg != "" {
			msg = cfg.msg + "\n" + msg
		}

		args := make([]any, 0, len(cfg.args)+1)
		args = append(args, cfg.args...)
		args = append(args, diff)

		t.Errorf(msg, args...)
		return false
	}
	return true
}

type assertConfig struct {
	cmpOpts []cmp.Option
	msg     string
	args    []any
}

// AssertOption allows for fine grain control over how AssertEqual operates.
type AssertOption interface {
	apply(cfg assertConfig) assertConfig
}

type fnOption func(cfg assertConfig) assertConfig

func (fn fnOption) apply(cfg assertConfig) assertConfig {
	return fn(cfg)
}

// Transform applies a transformation f function that
// converts values of a certain type into that of another.
// f must not mutate A in any way.
func Transform[A, B any](f func(A) B) AssertOption {
	return fnOption(func(cfg assertConfig) assertConfig {
		cfg.cmpOpts = append(cfg.cmpOpts, cmp.Transformer("", f))
		return cfg
	})
}

// Desc prepends the given text to an assertion failure message.
// The text is formatted with the args using fmt.Sprintf.
func Desc(text string, args ...any) AssertOption {
	return fnOption(func(cfg assertConfig) assertConfig {
		cfg.msg = text
		cfg.args = args
		return cfg
	})
}
