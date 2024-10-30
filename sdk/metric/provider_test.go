// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/go-logr/logr/funcr"
	"github.com/go-logr/logr/testr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestMeterConcurrentSafe(t *testing.T) {
	const name = "TestMeterConcurrentSafe meter"
	mp := NewMeterProvider()

	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = mp.Meter(name)
	}()

	_ = mp.Meter(name)
	<-done
}

func TestForceFlushConcurrentSafe(t *testing.T) {
	mp := NewMeterProvider()

	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = mp.ForceFlush(context.Background())
	}()

	_ = mp.ForceFlush(context.Background())
	<-done
}

func TestShutdownConcurrentSafe(t *testing.T) {
	mp := NewMeterProvider()

	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = mp.Shutdown(context.Background())
	}()

	_ = mp.Shutdown(context.Background())
	<-done
}

func TestMeterAndShutdownConcurrentSafe(t *testing.T) {
	const name = "TestMeterAndShutdownConcurrentSafe meter"
	mp := NewMeterProvider()

	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = mp.Shutdown(context.Background())
	}()

	_ = mp.Meter(name)
	<-done
}

func TestMeterDoesNotPanicForEmptyMeterProvider(t *testing.T) {
	mp := MeterProvider{}
	assert.NotPanics(t, func() { _ = mp.Meter("") })
}

func TestForceFlushDoesNotPanicForEmptyMeterProvider(t *testing.T) {
	mp := MeterProvider{}
	assert.NotPanics(t, func() { _ = mp.ForceFlush(context.Background()) })
}

func TestShutdownDoesNotPanicForEmptyMeterProvider(t *testing.T) {
	mp := MeterProvider{}
	assert.NotPanics(t, func() { _ = mp.Shutdown(context.Background()) })
}

func TestMeterProviderReturnsSameMeter(t *testing.T) {
	mp := MeterProvider{}
	mtr := mp.Meter("")

	assert.Same(t, mtr, mp.Meter(""))
	assert.NotSame(t, mtr, mp.Meter("diff"))
	assert.NotSame(t, mtr, mp.Meter("", api.WithInstrumentationAttributes(attribute.String("k", "v"))))
}

func TestEmptyMeterName(t *testing.T) {
	var buf strings.Builder
	warnLevel := 1
	l := funcr.New(func(prefix, args string) {
		_, _ = buf.WriteString(fmt.Sprint(prefix, args))
	}, funcr.Options{Verbosity: warnLevel})
	otel.SetLogger(l)
	mp := NewMeterProvider()

	mp.Meter("")

	assert.Contains(t, buf.String(), `"level"=1 "msg"="Invalid Meter name." "name"=""`)
}

func TestMeterProviderReturnsNoopMeterAfterShutdown(t *testing.T) {
	mp := NewMeterProvider()

	m := mp.Meter("")
	_, ok := m.(noop.Meter)
	assert.False(t, ok, "Meter from running MeterProvider is NoOp")

	require.NoError(t, mp.Shutdown(context.Background()))

	m = mp.Meter("")
	_, ok = m.(noop.Meter)
	assert.Truef(t, ok, "Meter from shutdown MeterProvider is not NoOp: %T", m)
}

func TestMeterProviderMixingOnRegisterErrors(t *testing.T) {
	otel.SetLogger(testr.New(t))

	rdr0 := NewManualReader()
	mp0 := NewMeterProvider(WithReader(rdr0))

	rdr1 := NewManualReader()
	mp1 := NewMeterProvider(WithReader(rdr1))

	// Meters with the same scope but different MeterProviders.
	m0 := mp0.Meter("TestMeterProviderMixingOnRegisterErrors")
	m1 := mp1.Meter("TestMeterProviderMixingOnRegisterErrors")

	m0Gauge, err := m0.Float64ObservableGauge("float64Gauge")
	require.NoError(t, err)

	m1Gauge, err := m1.Int64ObservableGauge("int64Gauge")
	require.NoError(t, err)

	_, err = m0.RegisterCallback(
		func(_ context.Context, o api.Observer) error {
			o.ObserveFloat64(m0Gauge, 2)
			// Observe an instrument from a different MeterProvider.
			o.ObserveInt64(m1Gauge, 1)

			return nil
		},
		m0Gauge, m1Gauge,
	)
	assert.Error(
		t,
		err,
		"Instrument registered with Meter from different MeterProvider",
	)

	var data metricdata.ResourceMetrics
	_ = rdr0.Collect(context.Background(), &data)
	// Only the metrics from mp0 should be produced.
	assert.Len(t, data.ScopeMetrics, 1)

	err = rdr1.Collect(context.Background(), &data)
	assert.NoError(t, err, "Errored when collect should be a noop")
	assert.Empty(
		t, data.ScopeMetrics,
		"Metrics produced for instrument collected by different MeterProvider",
	)
}
