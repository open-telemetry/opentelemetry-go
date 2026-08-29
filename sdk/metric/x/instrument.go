// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	metricx "go.opentelemetry.io/otel/metric/x"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

type bindingOption struct {
	sdkmetric.Option
}

func (bindingOption) WrapInt64CounterBinding(
	counter metric.Int64Counter,
	bind func(...attribute.KeyValue) metric.Int64Counter,
) metric.Int64Counter {
	return &bindingInt64Counter{Int64Counter: counter, bind: bind}
}

type finishingOption struct {
	sdkmetric.Option
}

func (finishingOption) WrapInt64CounterFinishing(
	counter metric.Int64Counter,
	finish func(context.Context, ...attribute.KeyValue),
) metric.Int64Counter {
	if binder, ok := counter.(metricx.Int64CounterBinder); ok {
		return &bindingFinishingInt64Counter{
			Int64Counter: counter,
			bind:         binder.Bind,
			finish:       finish,
		}
	}
	return &finishingInt64Counter{Int64Counter: counter, finish: finish}
}

type bindingInt64Counter struct {
	metric.Int64Counter
	bind func(...attribute.KeyValue) metric.Int64Counter
}

var (
	_ metric.Int64Counter        = (*bindingInt64Counter)(nil)
	_ metricx.Int64CounterBinder = (*bindingInt64Counter)(nil)
)

func (c *bindingInt64Counter) Bind(attrs ...attribute.KeyValue) metric.Int64Counter {
	return c.bind(attrs...)
}

type finishingInt64Counter struct {
	metric.Int64Counter
	finish func(context.Context, ...attribute.KeyValue)
}

var (
	_ metric.Int64Counter = (*finishingInt64Counter)(nil)
	_ metricx.Finisher    = (*finishingInt64Counter)(nil)
)

func (c *finishingInt64Counter) Finish(ctx context.Context, attrs ...attribute.KeyValue) {
	c.finish(ctx, attrs...)
}

type bindingFinishingInt64Counter struct {
	metric.Int64Counter
	bind   func(...attribute.KeyValue) metric.Int64Counter
	finish func(context.Context, ...attribute.KeyValue)
}

var (
	_ metric.Int64Counter        = (*bindingFinishingInt64Counter)(nil)
	_ metricx.Int64CounterBinder = (*bindingFinishingInt64Counter)(nil)
	_ metricx.Finisher           = (*bindingFinishingInt64Counter)(nil)
)

func (c *bindingFinishingInt64Counter) Bind(attrs ...attribute.KeyValue) metric.Int64Counter {
	return c.bind(attrs...)
}

func (c *bindingFinishingInt64Counter) Finish(ctx context.Context, attrs ...attribute.KeyValue) {
	c.finish(ctx, attrs...)
}
