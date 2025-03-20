// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logtest // import "go.opentelemetry.io/otel/log/logtest"

import (
	"context"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"go.opentelemetry.io/otel/log"
)

// TestingT reports failure messages.
// [testing.T] implements this interface.
type TestingT interface {
	Errorf(format string, args ...any)
}

// AssertEqual asserts that the two concrete data-types from the logtest package are equal.
func AssertEqual[T Recording | Record](t TestingT, want, got T, opts ...AssertOption) bool {
	if h, ok := t.(interface{ Helper() }); ok {
		h.Helper()
	}

	cmpOpts := []cmp.Option{
		cmp.Comparer(func(x, y context.Context) bool { return x == y }),           // Compare context.
		cmpopts.SortSlices(func(a, b log.KeyValue) bool { return a.Key < b.Key }), // Unordered compare of the key values.
		cmpopts.EquateEmpty(), // Empty and nil collections are equal.
	}

	cfg := newAssertConfig(opts)
	if cfg.ignoreTimestamp {
		cmpOpts = append(cmpOpts, cmpopts.IgnoreTypes(time.Time{})) // Ignore Timestamps.
	}

	if diff := cmp.Diff(want, got, cmpOpts...); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
		return false
	}
	return true
}

type assertConfig struct {
	ignoreTimestamp bool
}

func newAssertConfig(opts []AssertOption) assertConfig {
	var cfg assertConfig
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}
	return cfg
}

// AssertOption allows for fine grain control over how AssertEqual operates.
type AssertOption interface {
	apply(cfg assertConfig) assertConfig
}

type fnOption func(cfg assertConfig) assertConfig

func (fn fnOption) apply(cfg assertConfig) assertConfig {
	return fn(cfg)
}

// IgnoreTimestamp disables checking if timestamps are different.
func IgnoreTimestamp() AssertOption {
	return fnOption(func(cfg assertConfig) assertConfig {
		cfg.ignoreTimestamp = true
		return cfg
	})
}
