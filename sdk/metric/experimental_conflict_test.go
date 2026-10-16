// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
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

func TestValidateExperimental(t *testing.T) {
	a, b, c := expOptionA{}, expOptionB{}, expOptionC{}
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateExperimental(tt.opts)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateExperimentalErrorNamesOptions(t *testing.T) {
	err := validateExperimental([]experimentalOption{expOptionA{}, expOptionA{}, expOptionB{}, expOptionC{}})
	require.Error(t, err)
	assert.ErrorContains(t, err, "metric.expOptionA")
	assert.ErrorContains(t, err, "metric.expOptionB")
	assert.ErrorContains(t, err, "metric.expOptionC")
}

func TestNewMeterProviderDistinctExperimentalIsNoop(t *testing.T) {
	var logged string
	t.Cleanup(func(orig logr.Logger) func() {
		otel.SetLogger(funcr.New(func(_, args string) { logged += args }, funcr.Options{}))
		return func() { otel.SetLogger(orig) }
	}(logr.Discard()))

	rdr := NewManualReader()
	mp := NewMeterProvider(WithReader(rdr), testFinishOption(true), expOptionA{})

	assert.Equal(t, noop.Meter{}, mp.Meter("scope"), "distinct experimental options must disable the provider")
	assert.Nil(t, mp.meterFactory, "no experimental meter factory must be built")
	assert.Contains(t, logged, "metric.testFinishOption")
	assert.Contains(t, logged, "metric.expOptionA")

	// Readers keep their normal lifecycle.
	var rm metricdata.ResourceMetrics
	assert.NoError(t, rdr.Collect(t.Context(), &rm))
	assert.NoError(t, mp.ForceFlush(t.Context()))
	assert.NoError(t, mp.Shutdown(t.Context()))
	assert.ErrorIs(t, rdr.Collect(t.Context(), &rm), ErrReaderShutdown)
}

func TestNewMeterProviderRepeatedExperimental(t *testing.T) {
	mp := NewMeterProvider(testFinishOption(true), testFinishOption(true))
	t.Cleanup(func() { assert.NoError(t, mp.Shutdown(t.Context())) })

	counter, err := mp.Meter("scope").Int64Counter("requests")
	require.NoError(t, err)
	_, ok := counter.(metricx.Finisher)
	assert.True(t, ok, "repeating the same experimental option must keep its feature enabled")
}
