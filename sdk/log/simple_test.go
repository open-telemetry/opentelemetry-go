// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log_test

import (
	"context"
	"errors"
	"io"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	apilog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/log/internal/observ"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	semconv "go.opentelemetry.io/otel/semconv/v1.42.0"
	"go.opentelemetry.io/otel/semconv/v1.42.0/otelconv"
)

type exporter struct {
	records []log.Record
	calls   []string

	exportCalled     bool
	shutdownCalled   bool
	forceFlushCalled bool
	shutdownErr      error
	forceFlushErr    error
}

func (e *exporter) Export(_ context.Context, r []log.Record) error {
	e.records = slices.Clone(r) // don't retain the original slice
	e.exportCalled = true
	return nil
}

func (e *exporter) Shutdown(context.Context) error {
	e.calls = append(e.calls, "Shutdown")
	e.shutdownCalled = true
	return e.shutdownErr
}

func (e *exporter) ForceFlush(context.Context) error {
	e.calls = append(e.calls, "ForceFlush")
	e.forceFlushCalled = true
	return e.forceFlushErr
}

var _ log.Exporter = (*failingTestExporter)(nil)

type failingTestExporter struct {
	exporter
}

func (f *failingTestExporter) Export(ctx context.Context, r []log.Record) error {
	_ = f.exporter.Export(ctx, r)
	return assert.AnError
}

func TestSimpleProcessorEnabled(t *testing.T) {
	e := log.NewSimpleProcessor(nil)
	enabled := e.Enabled(t.Context(), log.EnabledParameters{})
	assert.True(t, enabled, "Enabled should return true")
}

func TestSimpleProcessorOnEmit(t *testing.T) {
	e := new(exporter)
	s := log.NewSimpleProcessor(e)

	r := new(log.Record)
	r.SetSeverityText("test")
	_ = s.OnEmit(t.Context(), r)

	require.True(t, e.exportCalled, "exporter Export not called")
	assert.Equal(t, []log.Record{*r}, e.records)
}

func TestSimpleProcessorShutdown(t *testing.T) {
	t.Run("FlushesBeforeShutdown", func(t *testing.T) {
		e := new(exporter)
		provider := log.NewLoggerProvider(
			log.WithProcessor(log.NewSimpleProcessor(e)),
		)

		require.NoError(t, provider.Shutdown(t.Context()))
		require.NoError(t, provider.Shutdown(t.Context()))
		assert.Equal(t, []string{"ForceFlush", "Shutdown"}, e.calls)
	})

	t.Run("Error", func(t *testing.T) {
		forceFlushErr := errors.New("force flush")
		shutdownErr := errors.New("shutdown")
		e := &exporter{
			forceFlushErr: forceFlushErr,
			shutdownErr:   shutdownErr,
		}
		provider := log.NewLoggerProvider(
			log.WithProcessor(log.NewSimpleProcessor(e)),
		)

		err := provider.Shutdown(t.Context())
		assert.ErrorIs(t, err, forceFlushErr)
		assert.ErrorIs(t, err, shutdownErr)
		assert.Equal(t, []string{"ForceFlush", "Shutdown"}, e.calls)
	})
}

func TestSimpleProcessorForceFlush(t *testing.T) {
	e := new(exporter)
	s := log.NewSimpleProcessor(e)
	_ = s.ForceFlush(t.Context())
	require.True(t, e.forceFlushCalled, "exporter ForceFlush not called")
}

type writerExporter struct {
	io.Writer
}

func (e *writerExporter) Export(_ context.Context, records []log.Record) error {
	for _, r := range records {
		_, _ = io.WriteString(e.Writer, r.Body().String())
	}
	return nil
}

func (*writerExporter) Shutdown(context.Context) error {
	return nil
}

func (*writerExporter) ForceFlush(context.Context) error {
	return nil
}

func TestSimpleProcessorEmpty(t *testing.T) {
	assert.NotPanics(t, func() {
		var s log.SimpleProcessor
		ctx := t.Context()
		record := new(log.Record)
		assert.NoError(t, s.OnEmit(ctx, record), "OnEmit")
		assert.NoError(t, s.ForceFlush(ctx), "ForceFlush")
		assert.NoError(t, s.Shutdown(ctx), "Shutdown")
	})
}

func TestSimpleProcessorConcurrentSafe(t *testing.T) {
	const goRoutineN = 10

	r := apilog.Record{}
	r.SetBody(attribute.StringValue("test"))
	ctx := t.Context()
	buf := new(strings.Builder)
	e := &writerExporter{buf}
	provider := log.NewLoggerProvider(
		log.WithProcessor(log.NewSimpleProcessor(e)),
	)
	logger := provider.Logger("test")

	var wg sync.WaitGroup
	flushErrs := make(chan error, goRoutineN)
	for range goRoutineN {
		wg.Go(func() {
			logger.Emit(ctx, r)
			flushErrs <- provider.ForceFlush(ctx)
		})
	}
	wg.Wait()
	close(flushErrs)
	for err := range flushErrs {
		require.NoError(t, err)
	}

	require.NoError(t, provider.Shutdown(ctx))
	assert.Equal(t, strings.Repeat("test", goRoutineN), buf.String())
}

func BenchmarkSimpleProcessorOnEmit(b *testing.B) {
	r := new(log.Record)
	r.SetSeverityText("test")
	ctx := b.Context()
	s := log.NewSimpleProcessor(nil)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var out error

		for pb.Next() {
			out = s.OnEmit(ctx, r)
		}

		_ = out
	})
}

func BenchmarkSimpleProcessorObservability(b *testing.B) {
	run := func(b *testing.B) {
		slp := log.NewSimpleProcessor(&failingTestExporter{exporter: exporter{}})
		record := new(log.Record)
		record.SetSeverityText("test")

		ctx := b.Context()
		b.ReportAllocs()
		b.ResetTimer()

		var err error
		for b.Loop() {
			err = slp.OnEmit(ctx, record)
		}
		_ = err
	}

	b.Run("Observability", func(b *testing.B) {
		b.Setenv("OTEL_GO_X_OBSERVABILITY", "true")
		run(b)
	})
	b.Run("NoObservability", run)
}

func TestSimpleLogProcessorObservability(t *testing.T) {
	testcases := []struct {
		name          string
		enabled       bool
		exporter      log.Exporter
		wantErr       error
		assertMetrics func(t *testing.T, rm metricdata.ResourceMetrics)
	}{
		{
			name:     "disabled",
			enabled:  false,
			exporter: new(exporter),
			assertMetrics: func(t *testing.T, rm metricdata.ResourceMetrics) {
				assert.Empty(t, rm.ScopeMetrics)
			},
		},
		{
			name:     "enabled",
			enabled:  true,
			exporter: new(exporter),
			assertMetrics: func(t *testing.T, rm metricdata.ResourceMetrics) {
				assert.Len(t, rm.ScopeMetrics, 1)
				sm := rm.ScopeMetrics[0]

				p := otelconv.SDKProcessorLogProcessed{}

				want := metricdata.ScopeMetrics{
					Scope: instrumentation.Scope{
						Name:      observ.ScopeName,
						Version:   sdk.Version(),
						SchemaURL: semconv.SchemaURL,
					},
					Metrics: []metricdata.Metrics{
						{
							Name:        p.Name(),
							Description: p.Description(),
							Unit:        p.Unit(),
							Data: metricdata.Sum[int64]{
								DataPoints: []metricdata.DataPoint[int64]{
									{
										Value: 1,
										Attributes: attribute.NewSet(
											observ.GetSLPComponentName(0),
											semconv.OTelComponentTypeKey.String(
												string(otelconv.ComponentTypeSimpleLogProcessor),
											),
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
					metricdatatest.IgnoreExemplars(),
					metricdatatest.IgnoreTimestamp(),
				)
			},
		},
		{
			name:    "Enable Exporter error",
			enabled: true,
			wantErr: assert.AnError,
			exporter: &failingTestExporter{
				exporter: exporter{},
			},
			assertMetrics: func(t *testing.T, rm metricdata.ResourceMetrics) {
				assert.Len(t, rm.ScopeMetrics, 1)
				sm := rm.ScopeMetrics[0]
				p := otelconv.SDKProcessorLogProcessed{}

				want := metricdata.ScopeMetrics{
					Scope: instrumentation.Scope{
						Name:      "go.opentelemetry.io/otel/sdk/log/internal/observ",
						Version:   sdk.Version(),
						SchemaURL: semconv.SchemaURL,
					},
					Metrics: []metricdata.Metrics{
						{
							Name:        p.Name(),
							Description: p.Description(),
							Unit:        p.Unit(),
							Data: metricdata.Sum[int64]{
								DataPoints: []metricdata.DataPoint[int64]{
									{
										Value: 1,
										Attributes: attribute.NewSet(
											observ.GetSLPComponentName(0),
											semconv.OTelComponentTypeKey.String(
												string(otelconv.ComponentTypeSimpleLogProcessor),
											),
											semconv.ErrorTypeKey.String("*errors.errorString"),
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

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("OTEL_GO_X_OBSERVABILITY", strconv.FormatBool(tc.enabled))

			original := otel.GetMeterProvider()
			t.Cleanup(func() {
				otel.SetMeterProvider(original)
			})

			r := metric.NewManualReader()
			mp := metric.NewMeterProvider(metric.WithReader(r))
			otel.SetMeterProvider(mp)

			slp := log.NewSimpleProcessor(tc.exporter)
			record := new(log.Record)
			record.SetSeverityText("test")
			err := slp.OnEmit(t.Context(), record)
			require.ErrorIs(t, err, tc.wantErr)
			var rm metricdata.ResourceMetrics
			require.NoError(t, r.Collect(t.Context(), &rm))
			tc.assertMetrics(t, rm)
			observ.SetSimpleProcessorID(0)
		})
	}
}
