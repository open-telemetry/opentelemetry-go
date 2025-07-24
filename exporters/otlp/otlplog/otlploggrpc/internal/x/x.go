// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package x contains support for OTel SDK experimental features.
//
// This package should only be used for features defined in the specification.
// It should not be used for experiments or new project ideas.

package x

import (
	"os"
	"strings"
)

type Feature[T any] struct {
	key   string
	parse func(v string) (T, bool)
}

func newFeature[T any](suffix string, parse func(string) (T, bool)) Feature[T] {
	const envKeyRoot = "OTEL_GO_X_"
	return Feature[T]{
		key:   envKeyRoot + suffix,
		parse: parse,
	}
}

var SelfObservability = newFeature("SELF_OBSERVABILITY", func(v string) (string, bool) {
	if strings.ToLower(v) == "true" {
		return v, true
	}
	return "", false
})

func (f Feature[T]) Key() string {
	return f.key
}

func (f Feature[T]) Lookup() (v T, ok bool) {
	vRaw := os.Getenv(f.key)
	if vRaw == "" {
		return v, ok
	}
	return f.parse(vRaw)
}

func (f Feature[T]) Enable() bool {
	_, ok := f.Lookup()
	return ok
}
