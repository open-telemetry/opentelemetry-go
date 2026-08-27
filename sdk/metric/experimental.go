// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	api "go.opentelemetry.io/otel/metric"
)

// int64CounterWrapper is implemented by an experimental Option that needs to
// wrap counters without adding experimental methods to stable instruments.
// All operations passed to the wrapper are backed by this provider's own
// pipelines and aggregation state.
type int64CounterWrapper interface {
	experimentalOption
	WrapInt64Counter(
		api.Int64Counter,
		func(...attribute.KeyValue) api.Int64Counter,
		func(context.Context, ...attribute.KeyValue),
	) api.Int64Counter
}
