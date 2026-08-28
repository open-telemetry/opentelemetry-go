// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	api "go.opentelemetry.io/otel/metric"
)

type int64CounterBindingWrapper func(
	api.Int64Counter,
	func(...attribute.KeyValue) api.Int64Counter,
) api.Int64Counter

type int64CounterFinishingWrapper func(
	api.Int64Counter,
	func(context.Context, ...attribute.KeyValue),
) api.Int64Counter
