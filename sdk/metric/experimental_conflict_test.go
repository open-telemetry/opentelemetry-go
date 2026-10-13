// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"reflect"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/funcr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// featureAOption, featureBOption and featureCOption identify stub experimental
// features, as finishOption identifies the Finish feature.
type (
	featureAOption interface{ featureA() }
	featureBOption interface{ featureB() }
	featureCOption interface{ featureC() }
)

type expMeterOptionA struct{}

func (expMeterOptionA) apply(config) config {
	panic("experimental options must not be applied as stable options")
}
func (expMeterOptionA) Experimental() {}
func (expMeterOptionA) featureA()     {}

type expMeterOptionB struct{}

func (expMeterOptionB) apply(config) config {
	panic("experimental options must not be applied as stable options")
}
func (expMeterOptionB) Experimental() {}
func (expMeterOptionB) featureB()     {}

type expMeterOptionC struct{}

func (expMeterOptionC) apply(config) config {
	panic("experimental options must not be applied as stable options")
}
func (expMeterOptionC) Experimental() {}
func (expMeterOptionC) featureC()     {}

var (
	conflictAB = experimentalConflict{reflect.TypeFor[featureAOption](), reflect.TypeFor[featureBOption]()}
	conflictAC = experimentalConflict{reflect.TypeFor[featureAOption](), reflect.TypeFor[featureCOption]()}
)

func TestValidateExperimental(t *testing.T) {
	a, b, c := expMeterOptionA{}, expMeterOptionB{}, expMeterOptionC{}
	tests := []struct {
		name      string
		opts      []experimentalOption
		conflicts []experimentalConflict
		wantErrs  int
	}{
		{"NoOptions", nil, []experimentalConflict{conflictAB}, 0},
		{"OnlyA", []experimentalOption{a}, []experimentalConflict{conflictAB}, 0},
		{"OnlyB", []experimentalOption{b}, []experimentalConflict{conflictAB}, 0},
		{"NoRules", []experimentalOption{a, b}, nil, 0},
		{"AB", []experimentalOption{a, b}, []experimentalConflict{conflictAB}, 1},
		{"BA", []experimentalOption{b, a}, []experimentalConflict{conflictAB}, 1},
		{"DuplicateA", []experimentalOption{a, a, b}, []experimentalConflict{conflictAB}, 1},
		{"BCNotTransitive", []experimentalOption{b, c}, []experimentalConflict{conflictAB, conflictAC}, 0},
		{"TwoPairs", []experimentalOption{a, b, c}, []experimentalConflict{conflictAB, conflictAC}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateExperimental(tt.opts, tt.conflicts)
			if tt.wantErrs == 0 {
				assert.NoError(t, err)
				return
			}
			require.Error(t, err)
			joined, ok := err.(interface{ Unwrap() []error })
			require.True(t, ok, "conflicts must be returned joined")
			assert.Len(t, joined.Unwrap(), tt.wantErrs)
		})
	}
}

func TestValidateExperimentalErrorNamesOptions(t *testing.T) {
	err := validateExperimental(
		[]experimentalOption{expMeterOptionA{}, expMeterOptionB{}},
		[]experimentalConflict{conflictAB},
	)
	require.Error(t, err)
	assert.ErrorContains(t, err, "metric.expMeterOptionA")
	assert.ErrorContains(t, err, "metric.expMeterOptionB")
}

func TestNewMeterProviderConflictingExperimentalIsNoop(t *testing.T) {
	var logged string
	t.Cleanup(func(orig logr.Logger) func() {
		otel.SetLogger(funcr.New(func(_, args string) { logged += args }, funcr.Options{}))
		return func() { otel.SetLogger(orig) }
	}(logr.Discard()))

	rdr := NewManualReader()
	mp := newMeterProvider(
		[]experimentalConflict{{reflect.TypeFor[finishOption](), reflect.TypeFor[featureAOption]()}},
		WithReader(rdr), testFinishOption(true), expMeterOptionA{},
	)

	assert.Equal(t, noop.Meter{}, mp.Meter("scope"), "a provider with conflicting options must return no-op meters")
	assert.Nil(t, mp.meterFactory, "no experimental meter factory must be built")
	assert.Contains(t, logged, "metric.testFinishOption")
	assert.Contains(t, logged, "metric.expMeterOptionA")

	// Readers keep their normal lifecycle.
	var rm metricdata.ResourceMetrics
	assert.NoError(t, rdr.Collect(t.Context(), &rm))
	assert.NoError(t, mp.ForceFlush(t.Context()))
	assert.NoError(t, mp.Shutdown(t.Context()))
	assert.ErrorIs(t, rdr.Collect(t.Context(), &rm), ErrReaderShutdown)
}

func TestNewMeterProviderNonConflictingExperimental(t *testing.T) {
	mp := newMeterProvider([]experimentalConflict{conflictAB}, expMeterOptionA{})
	assert.NotEqual(t, noop.Meter{}, mp.Meter("scope"), "options without a conflict must not disable the provider")
}

func TestNewMeterProviderDefaultConflicts(t *testing.T) {
	// The stub features are not production conflicts.
	mp := NewMeterProvider(expMeterOptionA{}, expMeterOptionB{})
	assert.NotEqual(t, noop.Meter{}, mp.Meter("scope"))
}
