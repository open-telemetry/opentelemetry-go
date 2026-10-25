// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// stableMetricFilter mirrors the internal interface the SDK checks for so we
// can verify the experimental option satisfies it structurally without the x
// package importing the SDK's internal types.
type stableMetricFilter interface {
	Experimental()
	TestMetric(instrumentationScope instrumentation.Scope, name string, kind metricdata.Aggregation, unit string) int
	TestAttributes(
		instrumentationScope instrumentation.Scope,
		name string,
		kind metricdata.Aggregation,
		unit string,
		attributes attribute.Set,
	) int
}

type testFilter struct {
	metricResult     Result
	attributesResult AttributesFilterResult
	metricParams     *TestMetricParameters
	attributesParams *TestAttributesParameters
}

func (f testFilter) TestMetric(params TestMetricParameters) Result {
	if f.metricParams != nil {
		*f.metricParams = params
	}
	return f.metricResult
}

func (f testFilter) TestAttributes(params TestAttributesParameters) AttributesFilterResult {
	if f.attributesParams != nil {
		*f.attributesParams = params
	}
	return f.attributesResult
}

func TestWithMetricFilterOptionSatisfiesStableInterface(t *testing.T) {
	opt := WithMetricFilter(testFilter{
		metricResult:     Accept,
		attributesResult: AttrAccept,
	})

	require.NotNil(t, opt)

	mf, ok := opt.(stableMetricFilter)
	require.True(t, ok, "option must satisfy the stable metricFilter interface")

	assert.NotPanics(t, mf.Experimental)
}

func TestWithMetricFilterNil(t *testing.T) {
	opt := WithMetricFilter(nil)

	_, ok := opt.(stableMetricFilter)
	assert.False(t, ok)

	_, experimental := opt.(interface{ Experimental() })
	assert.True(t, experimental)
}

func TestWithMetricFilterEnumMapping(t *testing.T) {
	tests := []struct {
		name     string
		filter   testFilter
		wantM    int
		wantAttr int
	}{
		{
			name:     "Accept",
			filter:   testFilter{metricResult: Accept, attributesResult: AttrAccept},
			wantM:    0,
			wantAttr: 0,
		},
		{
			name:     "Drop",
			filter:   testFilter{metricResult: Drop, attributesResult: AttrDrop},
			wantM:    1,
			wantAttr: 1,
		},
		{
			name:     "AcceptPartial",
			filter:   testFilter{metricResult: AcceptPartial, attributesResult: AttrAccept},
			wantM:    2,
			wantAttr: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithMetricFilter(tt.filter)
			mf := opt.(stableMetricFilter)

			scope := instrumentation.Scope{Name: "test"}
			kind := metricdata.Sum[int64]{}
			attrs := attribute.NewSet(attribute.String("k", "v"))

			assert.Equal(t, tt.wantM, mf.TestMetric(scope, "name", kind, "1"))
			assert.Equal(t, tt.wantAttr, mf.TestAttributes(scope, "name", kind, "1", attrs))
		})
	}
}
