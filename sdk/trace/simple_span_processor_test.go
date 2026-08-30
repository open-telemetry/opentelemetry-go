// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
	"go.opentelemetry.io/otel/semconv/v1.43.0/otelconv"
)

type simpleTestExporter struct {
	spans    []ReadOnlySpan
	shutdown bool
}

func (t *simpleTestExporter) ExportSpans(_ context.Context, spans []ReadOnlySpan) error {
	t.spans = append(t.spans, spans...)
	return nil
}

func (t *simpleTestExporter) Shutdown(ctx context.Context) error {
	t.shutdown = true
	select {
	case <-ctx.Done():
		// Ensure context deadline tests receive the expected error.
		return ctx.Err()
	default:
		return nil
	}
}

var _ SpanExporter = (*failingTestExporter)(nil)

type failingTestExporter struct {
	simpleTestExporter
}

func (f *failingTestExporter) ExportSpans(ctx context.Context, spans []ReadOnlySpan) error {
	_ = f.simpleTestExporter.ExportSpans(ctx, spans)
	return errors.New("failed to export spans")
}

// callbackTestExporter invokes onExport at the start of ExportSpans, allowing a
// test to observe state at the moment a span is submitted to the exporter.
type callbackTestExporter struct {
	simpleTestExporter
	onExport func()
}

func (c *callbackTestExporter) ExportSpans(ctx context.Context, spans []ReadOnlySpan) error {
	if c.onExport != nil {
		c.onExport()
	}
	return c.simpleTestExporter.ExportSpans(ctx, spans)
}

var _ SpanExporter = (*simpleTestExporter)(nil)

func TestNewSimpleSpanProcessor(t *testing.T) {
	if ssp := NewSimpleSpanProcessor(&simpleTestExporter{}); ssp == nil {
		t.Error("failed to create new SimpleSpanProcessor")
	}
}

func TestNewSimpleSpanProcessorWithNilExporter(t *testing.T) {
	ssp := NewSimpleSpanProcessor(nil)
	if ssp == nil {
		t.Error("failed to create new SimpleSpanProcessor with nil exporter")
	}
	if err := ssp.Shutdown(t.Context()); err != nil {
		t.Errorf("Shutdown with nil exporter should return nil, got %v", err)
	}
}

func TestSimpleSpanProcessorOnEnd(t *testing.T) {
	tp := basicTracerProvider(t)
	te := simpleTestExporter{}
	ssp := NewSimpleSpanProcessor(&te)

	tp.RegisterSpanProcessor(ssp)
	startSpan(tp, "TestSimpleSpanProcessorOnEnd").End()

	wantTraceID := tid
	gotTraceID := te.spans[0].SpanContext().TraceID()
	if wantTraceID != gotTraceID {
		t.Errorf("SimplerSpanProcessor OnEnd() check: got %+v, want %+v\n", gotTraceID, wantTraceID)
	}
}

func TestSimpleSpanProcessorShutdown(t *testing.T) {
	exporter := &simpleTestExporter{}
	ssp := NewSimpleSpanProcessor(exporter)

	// Ensure we can export a span before we test we cannot after shutdown.
	tp := basicTracerProvider(t)
	tp.RegisterSpanProcessor(ssp)
	startSpan(tp, "TestSimpleSpanProcessorShutdown").End()
	nExported := len(exporter.spans)
	if nExported != 1 {
		t.Error("failed to verify span export")
	}

	if err := ssp.Shutdown(t.Context()); err != nil {
		t.Errorf("shutting the SimpleSpanProcessor down: %v", err)
	}
	if !exporter.shutdown {
		t.Error("SimpleSpanProcessor.Shutdown did not shut down exporter")
	}

	startSpan(tp, "TestSimpleSpanProcessorShutdown").End()
	if len(exporter.spans) > nExported {
		t.Error("exported span to shutdown exporter")
	}
}

func TestSimpleSpanProcessorShutdownOnEndConcurrentSafe(t *testing.T) {
	exporter := &simpleTestExporter{}
	ssp := NewSimpleSpanProcessor(exporter)
	tp := basicTracerProvider(t)
	tp.RegisterSpanProcessor(ssp)

	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer func() {
			done <- struct{}{}
		}()
		for {
			select {
			case <-stop:
				return
			default:
				startSpan(tp, "TestSimpleSpanProcessorShutdownOnEndConcurrentSafe").End()
			}
		}
	}()

	if err := ssp.Shutdown(t.Context()); err != nil {
		t.Errorf("shutting the SimpleSpanProcessor down: %v", err)
	}
	if !exporter.shutdown {
		t.Error("SimpleSpanProcessor.Shutdown did not shut down exporter")
	}

	stop <- struct{}{}
	<-done
}

func TestSimpleSpanProcessorShutdownOnEndConcurrentSafe2(t *testing.T) {
	exporter := &simpleTestExporter{}
	ssp := NewSimpleSpanProcessor(exporter)
	tp := basicTracerProvider(t)
	tp.RegisterSpanProcessor(ssp)

	var wg sync.WaitGroup
	wg.Add(2)

	span := func(spanName string) {
		assert.NotPanics(t, func() {
			defer wg.Done()
			_, span := tp.Tracer("test").Start(t.Context(), spanName)
			span.End()
		})
	}

	go span("test-span-1")
	go span("test-span-2")

	wg.Wait()

	assert.NoError(t, ssp.Shutdown(t.Context()))
	assert.True(t, exporter.shutdown, "exporter shutdown")
}

func TestSimpleSpanProcessorShutdownHonorsContextDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), time.Nanosecond)
	defer cancel()
	<-ctx.Done()

	ssp := NewSimpleSpanProcessor(&simpleTestExporter{})
	if got, want := ssp.Shutdown(ctx), context.DeadlineExceeded; !errors.Is(got, want) {
		t.Errorf("SimpleSpanProcessor.Shutdown did not return %v, got %v", want, got)
	}
}

func TestSimpleSpanProcessorShutdownHonorsContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	ssp := NewSimpleSpanProcessor(&simpleTestExporter{})
	if got, want := ssp.Shutdown(ctx), context.Canceled; !errors.Is(got, want) {
		t.Errorf("SimpleSpanProcessor.Shutdown did not return %v, got %v", want, got)
	}
}

func TestSimpleSpanProcessorObservability(t *testing.T) {
	tests := []struct {
		name          string
		enabled       bool
		exporter      SpanExporter
		assertMetrics func(t *testing.T, rm metricdata.ResourceMetrics)
	}{
		{
			name:     "Disabled",
			enabled:  false,
			exporter: &simpleTestExporter{},
			assertMetrics: func(t *testing.T, rm metricdata.ResourceMetrics) {
				assert.Empty(t, rm.ScopeMetrics)
			},
		},
		{
			name:     "Enabled",
			enabled:  true,
			exporter: &simpleTestExporter{},
			assertMetrics: func(t *testing.T, rm metricdata.ResourceMetrics) {
				assert.Len(t, rm.ScopeMetrics, 1)
				sm := rm.ScopeMetrics[0]

				want := metricdata.ScopeMetrics{
					Scope: instrumentation.Scope{
						Name:      "go.opentelemetry.io/otel/sdk/trace/internal/observ",
						Version:   sdk.Version(),
						SchemaURL: semconv.SchemaURL,
					},
					Metrics: []metricdata.Metrics{
						{
							Name:        otelconv.SDKProcessorSpanProcessed{}.Name(),
							Description: otelconv.SDKProcessorSpanProcessed{}.Description(),
							Unit:        otelconv.SDKProcessorSpanProcessed{}.Unit(),
							Data: metricdata.Sum[int64]{
								DataPoints: []metricdata.DataPoint[int64]{
									{
										Value: 1,
										Attributes: attribute.NewSet(
											semconv.OTelComponentName("simple_span_processor/0"),
											semconv.OTelComponentTypeKey.String("simple_span_processor"),
										),
									},
								},
								Temporality: metricdata.CumulativeTemporality,
								IsMonotonic: true,
							},
						},
					},
				}

				metricdatatest.AssertEqual(
					t,
					want,
					sm,
					metricdatatest.IgnoreTimestamp(),
					metricdatatest.IgnoreExemplars(),
				)
			},
		},
		{
			name:    "Enabled, Exporter error",
			enabled: true,
			exporter: &failingTestExporter{
				simpleTestExporter: simpleTestExporter{},
			},
			assertMetrics: func(t *testing.T, rm metricdata.ResourceMetrics) {
				assert.Len(t, rm.ScopeMetrics, 1)
				sm := rm.ScopeMetrics[0]

				want := metricdata.ScopeMetrics{
					Scope: instrumentation.Scope{
						Name:      "go.opentelemetry.io/otel/sdk/trace/internal/observ",
						Version:   sdk.Version(),
						SchemaURL: semconv.SchemaURL,
					},
					Metrics: []metricdata.Metrics{
						{
							Name:        otelconv.SDKProcessorSpanProcessed{}.Name(),
							Description: otelconv.SDKProcessorSpanProcessed{}.Description(),
							Unit:        otelconv.SDKProcessorSpanProcessed{}.Unit(),
							Data: metricdata.Sum[int64]{
								DataPoints: []metricdata.DataPoint[int64]{
									{
										Value: 1,
										Attributes: attribute.NewSet(
											semconv.OTelComponentName("simple_span_processor/0"),
											semconv.OTelComponentTypeKey.String("simple_span_processor"),
										),
									},
								},
								Temporality: metricdata.CumulativeTemporality,
								IsMonotonic: true,
							},
						},
					},
				}

				metricdatatest.AssertEqual(
					t,
					want,
					sm,
					metricdatatest.IgnoreTimestamp(),
					metricdatatest.IgnoreExemplars(),
				)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("OTEL_GO_X_OBSERVABILITY", strconv.FormatBool(test.enabled))

			original := otel.GetMeterProvider()
			t.Cleanup(func() { otel.SetMeterProvider(original) })

			r := metric.NewManualReader()
			mp := metric.NewMeterProvider(
				metric.WithReader(r),
				metric.WithView(dropSpanMetricsView),
			)
			otel.SetMeterProvider(mp)

			ssp := NewSimpleSpanProcessor(test.exporter)
			tp := basicTracerProvider(t)
			tp.RegisterSpanProcessor(ssp)
			startSpan(tp, test.name).End()

			var rm metricdata.ResourceMetrics
			require.NoError(t, r.Collect(t.Context(), &rm))
			test.assertMetrics(t, rm)
			simpleProcessorIDCounter.Store(0) // reset simpleProcessorIDCounter
		})
	}
}

// processedSpanCount returns the current total value of the
// otel.sdk.processor.span.processed metric recorded by r.
func processedSpanCount(t *testing.T, r metric.Reader) int64 {
	t.Helper()
	var rm metricdata.ResourceMetrics
	require.NoError(t, r.Collect(t.Context(), &rm))
	name := otelconv.SDKProcessorSpanProcessed{}.Name()
	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			if m.Name != name {
				continue
			}
			var total int64
			for _, dp := range m.Data.(metricdata.Sum[int64]).DataPoints {
				total += dp.Value
			}
			return total
		}
	}
	return 0
}

// TestSimpleSpanProcessorProcessedBeforeExport asserts the span is counted as
// processed at the point it is submitted to the exporter, not after the export
// completes (semantic-conventions requirement).
func TestSimpleSpanProcessorProcessedBeforeExport(t *testing.T) {
	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	original := otel.GetMeterProvider()
	t.Cleanup(func() { otel.SetMeterProvider(original) })
	t.Cleanup(func() { simpleProcessorIDCounter.Store(0) })

	r := metric.NewManualReader()
	mp := metric.NewMeterProvider(metric.WithReader(r), metric.WithView(dropSpanMetricsView))
	otel.SetMeterProvider(mp)

	var duringExport int64
	exp := &callbackTestExporter{
		onExport: func() { duringExport = processedSpanCount(t, r) },
	}

	ssp := NewSimpleSpanProcessor(exp)
	tp := basicTracerProvider(t)
	tp.RegisterSpanProcessor(ssp)
	startSpan(tp, "test").End()

	assert.Equal(t, int64(1), duringExport,
		"processed must be recorded before the span is submitted to the exporter")
}
