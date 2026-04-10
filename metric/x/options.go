// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package x contains experimental metric options.
package x // import "go.opentelemetry.io/otel/metric/x"

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type defaultAttributesOption struct {
	metric.InstrumentOption
	keys []attribute.Key
}

// Experimental prevents the API from panicking when the option is used.
func (defaultAttributesOption) Experimental() {}

func (o defaultAttributesOption) AllowedKeys() []attribute.Key {
	return o.keys
}

// WithDefaultAttributes returns a metric.InstrumentOption that specifies default attribute keys.
// The implementation should treat keys that are not included in the passed keys as opt-in, and they should be filtered out by default.
// Users of [go.opentelemetry.io/otel/sdk/metric] can enable these attributes by using the AttributeFilter of a View to include them.
func WithDefaultAttributes(keys ...attribute.Key) metric.InstrumentOption {
	return defaultAttributesOption{keys: keys}
}
