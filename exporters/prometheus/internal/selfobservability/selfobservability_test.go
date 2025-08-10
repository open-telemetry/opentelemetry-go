package selfobservability

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestExporterMetrics(t *testing.T) {
	// Set up test meter provider
	reader := metric.NewManualReader()
	mp := metric.NewMeterProvider(metric.WithReader(reader))
	prev := otel.GetMeterProvider()
	otel.SetMeterProvider(mp)
	t.Cleanup(func() { otel.SetMeterProvider(prev) })

	em := NewExporterMetrics("test-component")
	ctx := context.Background()

	// Test AddInflight and AddExported
	em.AddInflight(ctx, 3)
	em.AddInflight(ctx, -1)
	em.AddExported(ctx, 5)

	// Test TrackCollectionDuration
	endCollect := em.TrackCollectionDuration(ctx)
	endCollect(nil)

	// Test TrackOperationDuration
	endOp := em.TrackOperationDuration(ctx)
	endOp(nil)

	// Collect metrics
	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(ctx, &rm))

	// Helper to find metrics
	findMetric := func(name string) *metricdata.Metrics {
		for _, sm := range rm.ScopeMetrics {
			for i := range sm.Metrics {
				if sm.Metrics[i].Name == name {
					return &sm.Metrics[i]
				}
			}
		}
		return nil
	}

	// Verify exported metric
	exported := findMetric("otel.sdk.exporter.metric_data_point.exported")
	require.NotNil(t, exported)
	switch data := exported.Data.(type) {
	case metricdata.Sum[int64]:
		assert.Equal(t, int64(5), data.DataPoints[0].Value)
	case metricdata.Sum[float64]:
		assert.InDelta(t, 5.0, data.DataPoints[0].Value, 0.001)
	}

	// Verify inflight metric (3 - 1 = 2)
	inflight := findMetric("otel.sdk.exporter.metric_data_point.inflight")
	require.NotNil(t, inflight)
	switch data := inflight.Data.(type) {
	case metricdata.Sum[int64]:
		assert.Equal(t, int64(2), data.DataPoints[0].Value)
	case metricdata.Sum[float64]:
		assert.InDelta(t, 2.0, data.DataPoints[0].Value, 0.001)
	}

	// Verify duration metrics exist
	collectionDuration := findMetric("otel.sdk.metric_reader.collection.duration")
	require.NotNil(t, collectionDuration)

	operationDuration := findMetric("otel.sdk.exporter.operation.duration")
	require.NotNil(t, operationDuration)
}
