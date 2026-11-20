// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdkmetricx "go.opentelemetry.io/otel/sdk/metric/x"
)

type dropNamedMetricFilter struct {
	name string
}

func (f dropNamedMetricFilter) TestMetric(params sdkmetricx.TestMetricParameters) sdkmetricx.Result {
	if params.Name == f.name {
		return sdkmetricx.Drop
	}
	return sdkmetricx.Accept
}

func (dropNamedMetricFilter) TestAttributes(sdkmetricx.TestAttributesParameters) sdkmetricx.AttributesFilterResult {
	return sdkmetricx.AttrAccept
}

func TestWithMetricFilterCollect(t *testing.T) {
	reader := sdkmetric.NewManualReader(
		sdkmetricx.WithMetricFilter(dropNamedMetricFilter{name: "drop"}),
	)
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	t.Cleanup(func() {
		require.NoError(t, provider.Shutdown(t.Context()))
	})

	meter := provider.Meter("test")
	keep, err := meter.Int64Counter("keep")
	require.NoError(t, err)
	drop, err := meter.Int64Counter("drop")
	require.NoError(t, err)

	keep.Add(t.Context(), 3)
	drop.Add(t.Context(), 7)

	var collected metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &collected))

	require.Len(t, collected.ScopeMetrics, 1)
	metrics := collected.ScopeMetrics[0].Metrics
	require.Len(t, metrics, 1)
	assert.Equal(t, "keep", metrics[0].Name)

	sum, ok := metrics[0].Data.(metricdata.Sum[int64])
	require.True(t, ok)
	require.Len(t, sum.DataPoints, 1)
	assert.Equal(t, int64(3), sum.DataPoints[0].Value)
}
