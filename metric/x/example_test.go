// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x_test

import (
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/metric/x"
)

func ExampleInt64CounterBinder() {
	var counter metric.Int64Counter = noop.Int64Counter{}
	attrs := []attribute.KeyValue{
		attribute.String("handler", "orders"),
	}

	isBound := false
	if b, ok := counter.(x.Int64CounterBinder); ok {
		counter = b.Bind(attrs...)
		isBound = true
	}

	http.HandleFunc("/orders", func(_ http.ResponseWriter, r *http.Request) {
		if isBound {
			counter.Add(r.Context(), 1)
		} else {
			counter.Add(r.Context(), 1, metric.WithAttributes(attrs...))
		}
	})
}
