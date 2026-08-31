// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
)

// Finisher is implemented by synchronous instruments that support ending the
// lifetime of a series identified by an exact attribute set.
//
// The set passed to [Finisher.Finish] is the same input set used when recording
// measurements, before any SDK View attribute filtering. An empty set
// identifies only the series recorded without attributes; it does not select
// every series.
//
// Finishing an unknown, already-finished, dropped, or unsupported series is a
// no-op. A set mapped to a shared cardinality-overflow series cannot be
// individually finished and is also a no-op.
//
// Finisher is an optional interface. Callers should use a type assertion to
// determine whether an instrument supports it. Instruments obtained from the
// global MeterProvider before a delegate is registered cannot acquire this
// optional method later. Callers that require this capability should obtain
// their instruments directly from a configured MeterProvider.
type Finisher interface {
	// Finish marks the series identified by set as finished. Its final
	// aggregate is eligible for one collection before the series state is
	// released. A measurement made before that collection cancels the pending
	// finish and continues the existing series lifetime.
	//
	// Finish is safe to call concurrently with measurement and collection. It
	// completes its in-process state transition before returning.
	Finish(ctx context.Context, set attribute.Set)
}
