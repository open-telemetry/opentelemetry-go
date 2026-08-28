// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x_test

import (
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	metricx "go.opentelemetry.io/otel/metric/x"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdkmetricx "go.opentelemetry.io/otel/sdk/metric/x"
)

func TestFeatureOptionsIndependent(t *testing.T) {
	tests := []struct {
		name       string
		options    []sdkmetric.Option
		wantBinder bool
		wantFinish bool
	}{
		{
			name:       "binding",
			options:    []sdkmetric.Option{sdkmetricx.WithBinding()},
			wantBinder: true,
		},
		{
			name:       "finish",
			options:    []sdkmetric.Option{sdkmetricx.WithFinish()},
			wantFinish: true,
		},
		{
			name:       "binding then finish",
			options:    []sdkmetric.Option{sdkmetricx.WithBinding(), sdkmetricx.WithFinish()},
			wantBinder: true,
			wantFinish: true,
		},
		{
			name:       "finish then binding",
			options:    []sdkmetric.Option{sdkmetricx.WithFinish(), sdkmetricx.WithBinding()},
			wantBinder: true,
			wantFinish: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			provider := sdkmetric.NewMeterProvider(test.options...)
			counter, err := provider.Meter("test").Int64Counter("requests")
			if err != nil {
				t.Fatal(err)
			}

			_, gotBinder := counter.(metricx.Int64CounterBinder)
			if gotBinder != test.wantBinder {
				t.Errorf("Int64CounterBinder implemented = %t, want %t", gotBinder, test.wantBinder)
			}
			_, gotFinish := counter.(metricx.Finisher)
			if gotFinish != test.wantFinish {
				t.Errorf("Finisher implemented = %t, want %t", gotFinish, test.wantFinish)
			}
		})
	}
}

func TestFinishWithoutBinding(t *testing.T) {
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader),
		sdkmetricx.WithFinish(),
	)
	counter, err := provider.Meter("test").Int64Counter("requests")
	if err != nil {
		t.Fatal(err)
	}
	attrs := []attribute.KeyValue{attribute.String("route", "/orders")}
	counter.Add(t.Context(), 3, metric.WithAttributes(attrs...))
	counter.(metricx.Finisher).Finish(t.Context(), attrs...)

	assertSum(
		t,
		collect(t, reader),
		metricdata.CumulativeTemporality,
		"/orders",
		3,
	)
}
