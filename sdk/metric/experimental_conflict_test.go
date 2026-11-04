// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"context"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/funcr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric/noop"
	metricx "go.opentelemetry.io/otel/metric/x"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// expOptionA, expOptionB and expOptionC are stub experimental options that
// enable no feature. Only their distinct types matter.
type (
	expOptionA struct{}
	expOptionB struct{}
	expOptionC struct{}
)

// expViewOptionA and expViewOptionB are stub experimental options that only
// contribute views.
type (
	expViewOptionA struct{}
	expViewOptionB struct{}
)

func (expOptionA) apply(config) config {
	panic("experimental options must not be applied as stable options")
}
func (expOptionA) Experimental() {}

func (expOptionB) apply(config) config {
	panic("experimental options must not be applied as stable options")
}
func (expOptionB) Experimental() {}

func (expOptionC) apply(config) config {
	panic("experimental options must not be applied as stable options")
}
func (expOptionC) Experimental() {}

func (expViewOptionA) apply(config) config {
	panic("experimental options must not be applied as stable options")
}
func (expViewOptionA) Experimental() {}
func (expViewOptionA) ExperimentalViews() []func(Instrument) (Stream, *bool, bool) {
	return nil
}

func (expViewOptionB) apply(config) config {
	panic("experimental options must not be applied as stable options")
}
func (expViewOptionB) Experimental() {}
func (expViewOptionB) ExperimentalViews() []func(Instrument) (Stream, *bool, bool) {
	return nil
}

func TestRejectDistinctExperimental(t *testing.T) {
	a, b, c := expOptionA{}, expOptionB{}, expOptionC{}
	va, vb := expViewOptionA{}, expViewOptionB{}
	tests := []struct {
		name    string
		opts    []experimentalOption
		wantErr bool
	}{
		{"None", nil, false},
		{"Single", []experimentalOption{a}, false},
		{"Repeated", []experimentalOption{a, a}, false},
		{"TwoDistinct", []experimentalOption{a, b}, true},
		{"TwoDistinctReversed", []experimentalOption{b, a}, true},
		{"ThreeDistinct", []experimentalOption{a, b, c}, true},
		{"RepeatedAndDistinct", []experimentalOption{a, a, b}, true},
		{"FinishAndOther", []experimentalOption{testFinishOption(true), a}, true},
		{"DisabledFinishAndOther", []experimentalOption{testFinishOption(false), a}, true},
		{"ViewOnly", []experimentalOption{va}, false},
		{"RepeatedView", []experimentalOption{va, va}, false},
		{"ViewAndMeter", []experimentalOption{va, a}, false},
		{"MeterAndView", []experimentalOption{a, va}, false},
		{"ViewAndFinish", []experimentalOption{va, testFinishOption(true), va}, false},
		{"RepeatedViewAndMeter", []experimentalOption{va, a, va, a}, false},
		{"TwoDistinctViews", []experimentalOption{va, vb}, true},
		{"TwoDistinctViewsAndMeter", []experimentalOption{va, a, vb}, true},
		{"TwoDistinctMetersAndView", []experimentalOption{a, va, b}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rejectDistinctExperimental(tt.opts)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRejectDistinctExperimentalErrorNamesOptions(t *testing.T) {
	err := rejectDistinctExperimental([]experimentalOption{expOptionA{}, expOptionA{}, expOptionB{}, expOptionC{}})
	require.Error(t, err)
	assert.ErrorContains(t, err, "metric.expOptionA")
	assert.ErrorContains(t, err, "metric.expOptionB")
	assert.ErrorContains(t, err, "metric.expOptionC")
}

func TestRejectDistinctExperimentalErrorNamesViewOptions(t *testing.T) {
	err := rejectDistinctExperimental([]experimentalOption{expViewOptionA{}, expOptionA{}, expViewOptionB{}})
	require.Error(t, err)
	assert.ErrorContains(t, err, "metric.expViewOptionA")
	assert.ErrorContains(t, err, "metric.expViewOptionB")
	assert.NotContains(t, err.Error(), "metric.expOptionA")
}

func TestNewMeterProviderDistinctExperimentalIsNoop(t *testing.T) {
	var logged string
	t.Cleanup(func(orig logr.Logger) func() {
		otel.SetLogger(funcr.New(func(_, args string) { logged += args }, funcr.Options{}))
		return func() { otel.SetLogger(orig) }
	}(logr.Discard()))

	rdr := NewManualReader()
	// ManualReader has no ForceFlush, so flushRdr verifies flushes still
	// reach the Readers.
	flushRdr := &reader{
		forceFlushFunc: func(context.Context) error { return assert.AnError },
		shutdownFunc:   func(context.Context) error { return nil },
	}
	mp := NewMeterProvider(WithReader(rdr), WithReader(flushRdr), testFinishOption(true), expOptionA{})

	assert.Equal(t, noop.Meter{}, mp.Meter("scope"), "distinct experimental options must disable the provider")
	assert.Nil(t, mp.meterFactory, "no experimental meter factory must be built")
	assert.Contains(t, logged, "metric.testFinishOption")
	assert.Contains(t, logged, "metric.expOptionA")

	// Readers keep their normal lifecycle.
	var rm metricdata.ResourceMetrics
	assert.NoError(t, rdr.Collect(t.Context(), &rm))
	assert.ErrorIs(t, mp.ForceFlush(t.Context()), assert.AnError)
	assert.NoError(t, mp.Shutdown(t.Context()))
	assert.ErrorIs(t, rdr.Collect(t.Context(), &rm), ErrReaderShutdown)
}

func TestNewMeterProviderExperimentalViewWithMeterOption(t *testing.T) {
	mp := NewMeterProvider(testFinishOption(true), expViewOptionA{})
	t.Cleanup(func() { assert.NoError(t, mp.Shutdown(t.Context())) })

	counter, err := mp.Meter("scope").Int64Counter("requests")
	require.NoError(t, err)
	_, ok := counter.(metricx.Finisher)
	assert.True(t, ok, "an experimental view option must not disable a meter option")
}

func TestNewMeterProviderRepeatedExperimental(t *testing.T) {
	mp := NewMeterProvider(testFinishOption(true), testFinishOption(true))
	t.Cleanup(func() { assert.NoError(t, mp.Shutdown(t.Context())) })

	counter, err := mp.Meter("scope").Int64Counter("requests")
	require.NoError(t, err)
	_, ok := counter.(metricx.Finisher)
	assert.True(t, ok, "repeating the same experimental option must keep its feature enabled")
}
