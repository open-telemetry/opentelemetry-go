// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logtest // import "go.opentelemetry.io/otel/log/logtest"

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"go.opentelemetry.io/otel/log"
)

// AssertEqual asserts that the two concrete data-types from the logtest package are equal.
func AssertEqual[T Recording | Record](t *testing.T, want, got T, opts ...AssertOption) bool {
	t.Helper()
	return assertEqual(t, want, got, opts...)
}

// testingT reports failure messages.
// *testing.T implements this interface.
type testingT interface {
	Errorf(format string, args ...any)
}

func assertEqual[T Recording | Record](t testingT, want, got T, opts ...AssertOption) bool {
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
		t.Errorf("mismatch (-want +got):\n%s", diff)
		return false
	}
	return true
}

type assertConfig struct {
	cmpOpts []cmp.Option
}

// AssertOption allows for fine grain control over how AssertEqual operates.
type AssertOption interface {
	apply(cfg assertConfig) assertConfig
}

type fnOption func(cfg assertConfig) assertConfig

func (fn fnOption) apply(cfg assertConfig) assertConfig {
	return fn(cfg)
}

// WithCmpOpts adds additional [cmp.Option] used when comparing values.
func WithCmpOpts(opts ...cmp.Option) AssertOption {
	return fnOption(func(cfg assertConfig) assertConfig {
		cfg.cmpOpts = append(cfg.cmpOpts, opts...)
		return cfg
	})
}
