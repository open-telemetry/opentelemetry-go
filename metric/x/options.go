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

// Resettable is an optional interface that Options can implement
// to allow reuse without additional allocations.
//
// Example usage with sync.Pool:
//
//	var optionPool = sync.Pool{
//		New: func() any {
//			return metric.WithAttributeSet(*attribute.EmptySet())
//		},
//	}
//
//	func record(ctx context.Context, counter metric.Int64Counter, set attribute.Set) {
//		opt := optionPool.Get().(metric.MeasurementOption)
//		defer optionPool.Put(opt)
//
//		if r, ok := opt.(x.Resettable[attribute.Set]); ok {
//			r.Reset(set)
//		} else {
//			opt = metric.WithAttributeSet(set)
//		}
//		counter.Add(ctx, 1, opt)
//	}
//
// WARNING: It is the user's responsibility to ensure that the option is not
// concurrently reset while being passed to the API or used by another goroutine.
type Resettable[T any] interface {
	Reset(T)
}
