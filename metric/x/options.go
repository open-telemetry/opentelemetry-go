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

type unsafeAttributesOption struct {
	metric.MeasurementOption
	kvs []attribute.KeyValue
}

// Experimental prevents the API from panicking when the option is used.
func (*unsafeAttributesOption) Experimental() {}

// RawAttributes returns the raw key-values associated with the option.
func (o *unsafeAttributesOption) RawAttributes() []attribute.KeyValue {
	return o.kvs
}

// Settable implements the x.Settable interface.
// As with all Settable implementations, Set must not be called concurrently
// with the usage of the option in an instrument call (e.g., counter.Add)
// or by other goroutines.
func (o *unsafeAttributesOption) Set(kvs []attribute.KeyValue) {
	o.kvs = kvs
}

// WithUnsafeAttributes returns a metric.MeasurementOption that stores the raw attributes
// and associates them with a measurement without making a copy.
// The caller must not modify the attributes slice after passing it to this function.
// This is a work-in-progress, and does not yet have better performance than metric.WithAttributeSet.
func WithUnsafeAttributes(kvs ...attribute.KeyValue) metric.MeasurementOption {
	return &unsafeAttributesOption{kvs: kvs}
}

// Settable is an optional interface that Options can implement
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
//		if r, ok := opt.(x.Settable[attribute.Set]); ok {
//			r.Set(set)
//		} else {
//			opt = metric.WithAttributeSet(set)
//		}
//		counter.Add(ctx, 1, opt)
//	}
//
// WARNING: It is the user's responsibility to ensure that the option is not
// concurrently set while being passed to the API or used by another goroutine.
type Settable[T any] interface {
	Set(T)
}
