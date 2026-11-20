// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
)

// Finisher is implemented by synchronous instruments that support ending the
// lifetime of a series identified by exact attributes.
//
// The attributes passed to [Finisher.Finish] are matched as a complete
// collection, not as a subset. Attribute order does not affect identity, and
// duplicate keys use last-value-wins semantics. Calling Finish without
// attributes identifies only the series recorded without attributes; it does
// not select every series.
//
// Finishing an unknown, already-finished, dropped, or unsupported series is a
// no-op. Finishing attributes mapped to a shared cardinality-overflow series
// is also a no-op.
//
// Finisher is an optional interface. Callers should use a type assertion to
// determine whether an instrument supports it. Instruments obtained from the
// global MeterProvider before a delegate is registered cannot acquire this
// optional method later. Callers that require this capability should obtain
// their instruments directly from a configured MeterProvider.
type Finisher interface {
	// Finish marks the series identified by attrs as finished. Its final
	// aggregate is eligible for one collection before the series state is
	// released.
	//
	// Implementations must not retain or modify attrs.
	//
	// Finish is safe to call concurrently with measurement and collection.
	Finish(ctx context.Context, attrs ...attribute.KeyValue)
}
