// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package x documents experimental features for [go.opentelemetry.io/otel/sdk/log].
package x // import "go.opentelemetry.io/otel/sdk/log/internal/x"

import (
	"os"
	"strings"
)

// Observability is an experimental feature flag that determines if SDK
// observability metrics are enabled.
//
// To enable this feature set the OTEL_GO_X_OBSERVABILITY environment variable
// to the case-insensitive string value of "true" (i.e. "True" and "TRUE"
// will also enable this).
var Observability = newFeature(
	[]string{"OBSERVABILITY", "SELF_OBSERVABILITY"},
	func(v string) (string, bool) {
		if strings.EqualFold(v, "true") {
			return v, true
		}
		return "", false
	},
)

// Feature is an experimental feature control flag. It provides a uniform way
// to interact with these feature flags and parse their values.
type Feature[T any] struct {
	keys  []string
	parse func(v string) (T, bool)
}

func newFeature[T any](suffix []string, parse func(string) (T, bool)) Feature[T] {
	const envKeyRoot = "OTEL_GO_X_"
	keys := make([]string, 0, len(suffix))
	for _, s := range suffix {
		keys = append(keys, envKeyRoot+s)
	}
	return Feature[T]{
		keys:  keys,
		parse: parse,
	}
}

// Keys returns the environment variable keys that can be set to enable the
// feature.
func (f Feature[T]) Keys() []string { return f.keys }

// Lookup returns the user configured value for the feature and true if the
// user has enabled the feature. Otherwise, if the feature is not enabled, a
// zero-value and false are returned.
func (f Feature[T]) Lookup() (v T, ok bool) {
	// https://github.com/open-telemetry/opentelemetry-specification/blob/62effed618589a0bec416a87e559c0a9d96289bb/specification/configuration/sdk-environment-variables.md#parsing-empty-value
	//
	// > The SDK MUST interpret an empty value of an environment variable the
	// > same way as when the variable is unset.
	for _, key := range f.keys {
		vRaw := os.Getenv(key)
		if vRaw != "" {
			return f.parse(vRaw)
		}
	}
	return v, ok
}

// Enabled reports whether the feature is enabled.
func (f Feature[T]) Enabled() bool {
	_, ok := f.Lookup()
	return ok
}
