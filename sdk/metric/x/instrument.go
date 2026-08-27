// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x // import "go.opentelemetry.io/otel/sdk/metric/x"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	metricx "go.opentelemetry.io/otel/metric/x"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

type boundCounterOption struct {
	sdkmetric.Option
}

func (boundCounterOption) Experimental() {}

func (boundCounterOption) WrapInt64Counter(
	counter metric.Int64Counter,
	bind func(...attribute.KeyValue) metric.Int64Counter,
	finish func(context.Context, ...attribute.KeyValue),
) metric.Int64Counter {
	return &int64Counter{Int64Counter: counter, bind: bind, finish: finish}
}

type int64Counter struct {
	metric.Int64Counter
	bind   func(...attribute.KeyValue) metric.Int64Counter
	finish func(context.Context, ...attribute.KeyValue)
}

var (
	_ metric.Int64Counter        = (*int64Counter)(nil)
	_ metricx.Int64CounterBinder = (*int64Counter)(nil)
	_ metricx.Finisher           = (*int64Counter)(nil)
)

func (c *int64Counter) Bind(attrs ...attribute.KeyValue) metric.Int64Counter {
	return c.bind(attrs...)
}

func (c *int64Counter) Finish(ctx context.Context, attrs ...attribute.KeyValue) {
	c.finish(ctx, attrs...)
}
