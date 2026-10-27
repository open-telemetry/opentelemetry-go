// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x_test

import (
	"slices"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/x"
)

func TestWithOptIn(t *testing.T) {
	opt := x.WithOptIn()
	if opt == nil {
		t.Fatal("WithOptIn() returned nil")
	}

	exp, ok := opt.(interface{ Experimental() })
	if !ok {
		t.Fatal("WithOptIn() does not implement Experimental()")
	}
	exp.Experimental()

	optIn, ok := opt.(interface{ OptIn() bool })
	if !ok {
		t.Fatal("WithOptIn() does not implement OptIn() bool")
	}
	if !optIn.OptIn() {
		t.Fatal("OptIn() returned false, want true")
	}

	// Ensure all instrument config constructors accept WithOptIn without panicking.
	_ = metric.NewInt64CounterConfig(opt)
	_ = metric.NewInt64UpDownCounterConfig(opt)
	_ = metric.NewInt64HistogramConfig(opt)
	_ = metric.NewInt64GaugeConfig(opt)
	_ = metric.NewInt64ObservableCounterConfig(opt)
	_ = metric.NewInt64ObservableUpDownCounterConfig(opt)
	_ = metric.NewInt64ObservableGaugeConfig(opt)
	_ = metric.NewFloat64CounterConfig(opt)
	_ = metric.NewFloat64UpDownCounterConfig(opt)
	_ = metric.NewFloat64HistogramConfig(opt)
	_ = metric.NewFloat64GaugeConfig(opt)
	_ = metric.NewFloat64ObservableCounterConfig(opt)
	_ = metric.NewFloat64ObservableUpDownCounterConfig(opt)
	_ = metric.NewFloat64ObservableGaugeConfig(opt)
}

func TestWithDefaultAttributes(t *testing.T) {
	k := attribute.Key("k")
	opt := x.WithDefaultAttributes(k)
	if opt == nil {
		t.Fatal("WithDefaultAttributes() returned nil")
	}

	allowed, ok := opt.(interface{ AllowedKeys() []attribute.Key })
	if !ok {
		t.Fatal("WithDefaultAttributes() does not implement AllowedKeys() []attribute.Key")
	}
	if got, want := allowed.AllowedKeys(), []attribute.Key{k}; !slices.Equal(got, want) {
		t.Fatalf("AllowedKeys() = %v, want %v", got, want)
	}
}
