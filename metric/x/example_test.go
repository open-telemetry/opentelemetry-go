// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x_test

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/metric/x"
)

func ExampleFinisher() {
	var counter metric.Int64Counter = noop.Int64Counter{}
	ctx := context.Background()

	// Reuse the attributes for recording and finishing.
	attrs := []attribute.KeyValue{
		attribute.String("service.name", "checkout"),
		attribute.Int("service.port", 8080),
	}
	counter.Add(ctx, 1, metric.WithAttributes(attrs...))

	if finisher, ok := counter.(x.Finisher); ok {
		finisher.Finish(ctx, attrs...)
	}
}

func ExampleInt64CounterBinder() {
	var counter metric.Int64Counter = noop.Int64Counter{}
	attrs := []attribute.KeyValue{
		attribute.String("handler", "orders"),
	}

	var opts []metric.AddOption
	if b, ok := counter.(x.Int64CounterBinder); ok {
		counter = b.Bind(attrs...)
	} else {
		opts = append(opts, metric.WithAttributeSet(attribute.NewSet(attrs...)))
	}

	http.HandleFunc("/orders", func(_ http.ResponseWriter, r *http.Request) {
		counter.Add(r.Context(), 1, opts...)
	})
}
