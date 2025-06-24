// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
)

// func TestBatchSpanProcessorMetricsDisabled(t *testing.T) {
// 	// TODO: test queueCapacityUpDownCounter
// 	// TODO: test queueSizeUpDownCounter
// 	// TODO: test spansProcessedCounter
// }

// pausedExporter waits until shutdown to export spans.
type pausedExporter struct {
	stop chan (struct{})
}

func newPausedExporter(t *testing.T) pausedExporter {
	e := pausedExporter{stop: make(chan struct{})}
	return e
}

func (e pausedExporter) Shutdown(context.Context) error {
	return nil
}

func (e pausedExporter) ExportSpans(ctx context.Context, _ []ReadOnlySpan) error {
	<-e.stop
	return ctx.Err()
}

func TestBatchSpanProcessorQueueMetrics(t *testing.T) {
	t.Setenv("OTEL_GO_X_SELF_OBSERVABILITY", "true")
	tp := basicTracerProvider(t)
	reader := sdkmetric.NewManualReader()
	meterProvider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	pe := newPausedExporter(t)
	bsp := NewBatchSpanProcessor(
		pe,
		WithBlocking(),
		WithMaxExportBatchSize(1),
		withMeterProvider(meterProvider),
	)
	tp.RegisterSpanProcessor(bsp)
	t.Cleanup(func() {
		tp.UnregisterSpanProcessor(bsp)
	})
	t.Cleanup(func() {
		close(pe.stop)
	})

	tr := tp.Tracer("TestBatchSpanProcessorQueueMetrics")
	generateSpan(t, tr, testOption{genNumSpans: 1})
	gotMetrics := new(metricdata.ResourceMetrics)
	reader.Collect(context.Background(), gotMetrics)
	require.Len(t, gotMetrics.ScopeMetrics, 1)
	assertScopeMetrics(t, gotMetrics.ScopeMetrics[0], expectMetrics{queueCapacity: 2048, queueSize: 1})
	// TODO: test queueCapacityUpDownCounter
	// TODO: test queueSizeUpDownCounter
}

// func TestBatchSpanProcessorDropOnQueueFullMetrics(t *testing.T) {
// 	// TODO: test spansProcessedCounter with queueFullAttributes (non-blocking, queue full)
// 	// TODO: test spansProcessedCounter with queueFullAttributes (blocking, ctx expires)
// 	// TODO: test spansProcessedCounter with successAttributes
// 	// TODO: test spansProcessedCounter with alreadyShutdownAttributes
// }

type expectMetrics struct {
	queueCapacity            int64
	queueSize                int64
	successProcessed         int64
	alreadyShutdownProcessed int64
	queueFullProcessed       int64
}

func assertScopeMetrics(t *testing.T, sm metricdata.ScopeMetrics, expectation expectMetrics) {
	assert.Equal(t, instrumentation.Scope{
		Name:    "go.opentelemetry.io/otel/sdk/trace",
		Version: version(),
	}, sm.Scope)
	wantProcessedDataPoints := []metricdata.DataPoint[int64]{}
	if expectation.successProcessed > 0 {
		// TODO attrs
		wantProcessedDataPoints = append(wantProcessedDataPoints, metricdata.DataPoint[int64]{Value: expectation.successProcessed})
	}
	if expectation.alreadyShutdownProcessed > 0 {
		// TODO attrs
		wantProcessedDataPoints = append(wantProcessedDataPoints, metricdata.DataPoint[int64]{Value: expectation.alreadyShutdownProcessed})
	}
	if expectation.queueFullProcessed > 0 {
		// TODO attrs
		wantProcessedDataPoints = append(wantProcessedDataPoints, metricdata.DataPoint[int64]{Value: expectation.queueFullProcessed})
	}

	if len(wantProcessedDataPoints) > 0 {
		require.Len(t, sm.Metrics, 3)
	} else {
		require.Len(t, sm.Metrics, 2)
	}

	componentTypeAttr := componentTypeKey.String("batching_span_processor")
	componentNameAttr := componentNameKey.String(fmt.Sprintf("batching_span_processor/%d", int64(processorID.Load()-1)))

	baseAttrs := attribute.NewSet(
		componentTypeAttr,
		componentNameAttr,
	)

	want := metricdata.Metrics{
		Name:        queueCapacityMetricName,
		Description: queueCapacityMetricDescription,
		Unit:        spanCountUnit,
		Data: metricdata.Sum[int64]{
			DataPoints:  []metricdata.DataPoint[int64]{{Attributes: baseAttrs, Value: expectation.queueCapacity}},
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: false,
		},
	}
	metricdatatest.AssertEqual(t, want, sm.Metrics[0], metricdatatest.IgnoreTimestamp())

	want = metricdata.Metrics{
		Name:        queueSizeMetricName,
		Description: queueSizeMetricDescription,
		Unit:        spanCountUnit,
		Data: metricdata.Sum[int64]{
			DataPoints:  []metricdata.DataPoint[int64]{{Attributes: baseAttrs, Value: expectation.queueSize}},
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: false,
		},
	}
	metricdatatest.AssertEqual(t, want, sm.Metrics[1], metricdatatest.IgnoreTimestamp())

	if len(wantProcessedDataPoints) > 0 {
		want = metricdata.Metrics{
			Name:        spansProcessedMetricName,
			Description: spansProcessedMetricDescription,
			Unit:        spanCountUnit,
			Data: metricdata.Sum[int64]{
				DataPoints:  wantProcessedDataPoints,
				Temporality: metricdata.CumulativeTemporality,
				IsMonotonic: true,
			},
		}
		metricdatatest.AssertEqual(t, want, sm.Metrics[2], metricdatatest.IgnoreTimestamp(), metricdatatest.IgnoreValue())
	}
}
