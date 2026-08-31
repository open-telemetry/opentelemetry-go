// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x_test

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/metric/x"
)

func ExampleFinisher() {
	var counter metric.Int64Counter = noop.Int64Counter{}
	ctx := context.Background()

	// Construct the set once and reuse it for recording and finishing.
	attrs := attribute.NewSet(
		attribute.String("service.name", "checkout"),
		attribute.Int("service.port", 8080),
	)
	counter.Add(ctx, 1, metric.WithAttributeSet(attrs))

	if finisher, ok := counter.(x.Finisher); ok {
		finisher.Finish(ctx, attrs)
	}
}
