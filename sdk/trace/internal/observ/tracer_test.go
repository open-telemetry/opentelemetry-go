// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package observ_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/internal/observ"
	"go.opentelemetry.io/otel/semconv/v1.39.0/otelconv"
	tapi "go.opentelemetry.io/otel/trace"
)

func live(dPts ...metricdata.DataPoint[int64]) metricdata.Metrics {
	return metricdata.Metrics{
		Name:        otelconv.SDKSpanLive{}.Name(),
		Description: otelconv.SDKSpanLive{}.Description(),
		Unit:        otelconv.SDKSpanLive{}.Unit(),
		Data: metricdata.Sum[int64]{
			Temporality: metricdata.CumulativeTemporality,
			DataPoints:  dPts,
		},
	}
}

func sampledLive(value int64) metricdata.Metrics {
	set := attribute.NewSet(
		otelconv.SDKSpanLive{}.AttrSpanSamplingResult(
			otelconv.SpanSamplingResultRecordAndSample,
		),
	)
	return live(dPt(set, value))
}

func started(dPts ...metricdata.DataPoint[int64]) metricdata.Metrics {
	return metricdata.Metrics{
		Name:        otelconv.SDKSpanStarted{}.Name(),
		Description: otelconv.SDKSpanStarted{}.Description(),
		Unit:        otelconv.SDKSpanStarted{}.Unit(),
		Data: metricdata.Sum[int64]{
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: true,
			DataPoints:  dPts,
		},
	}
}

func sampledStarted(value int64) metricdata.Metrics {
	set := attribute.NewSet(
		otelconv.SDKSpanStarted{}.AttrSpanParentOrigin(
			otelconv.SpanParentOriginNone,
		),
		otelconv.SDKSpanStarted{}.AttrSpanSamplingResult(
			otelconv.SpanSamplingResultRecordAndSample,
		),
	)
	return started(dPt(set, value))
}

func TestTracer(t *testing.T) {
	collect := setup(t)
	tracer := trace.NewTracerProvider().Tracer(t.Name())

	_, span := tracer.Start(t.Context(), "span")
	check(t, collect(), sampledLive(1), sampledStarted(1))

	span.End()
	check(t, collect(), sampledLive(0), sampledStarted(1))
}

func dropStarted(value int64) metricdata.Metrics {
	set := attribute.NewSet(
		otelconv.SDKSpanStarted{}.AttrSpanParentOrigin(
			otelconv.SpanParentOriginNone,
		),
		otelconv.SDKSpanStarted{}.AttrSpanSamplingResult(
			otelconv.SpanSamplingResultDrop,
		),
	)
	return started(dPt(set, value))
}

func TestTracerNonRecording(t *testing.T) {
	collect := setup(t)
	tracer := trace.NewTracerProvider(
		trace.WithSampler(trace.NeverSample()),
	).Tracer(t.Name())

	_, _ = tracer.Start(t.Context(), "span")
	check(t, collect(), dropStarted(1))
}

func recLive(value int64) metricdata.Metrics {
	set := attribute.NewSet(
		otelconv.SDKSpanLive{}.AttrSpanSamplingResult(
			otelconv.SpanSamplingResultRecordOnly,
		),
	)
	return live(dPt(set, value))
}

func recStarted(value int64) metricdata.Metrics {
	set := attribute.NewSet(
		otelconv.SDKSpanStarted{}.AttrSpanParentOrigin(
			otelconv.SpanParentOriginNone,
		),
		otelconv.SDKSpanStarted{}.AttrSpanSamplingResult(
			otelconv.SpanSamplingResultRecordOnly,
		),
	)
	return started(dPt(set, value))
}

type recOnly struct{}

func (recOnly) ShouldSample(p trace.SamplingParameters) trace.SamplingResult {
	psc := tapi.SpanContextFromContext(p.ParentContext)
	return trace.SamplingResult{
		Decision:   trace.RecordOnly,
		Tracestate: psc.TraceState(),
	}
}

func (recOnly) Description() string { return "RecordingOnly" }

func TestTracerRecordOnly(t *testing.T) {
	collect := setup(t)
	tracer := trace.NewTracerProvider(
		trace.WithSampler(recOnly{}),
	).Tracer(t.Name())

	_, _ = tracer.Start(t.Context(), "span")
	check(t, collect(), recLive(1), recStarted(1))
}

func remoteStarted(value int64) metricdata.Metrics {
	set := attribute.NewSet(
		otelconv.SDKSpanStarted{}.AttrSpanParentOrigin(
			otelconv.SpanParentOriginRemote,
		),
		otelconv.SDKSpanStarted{}.AttrSpanSamplingResult(
			otelconv.SpanSamplingResultRecordAndSample,
		),
	)
	return started(dPt(set, value))
}

func TestTracerRemoteParent(t *testing.T) {
	collect := setup(t)
	tracer := trace.NewTracerProvider().Tracer(t.Name())

	ctx := tapi.ContextWithRemoteSpanContext(
		t.Context(),
		tapi.NewSpanContext(tapi.SpanContextConfig{
			TraceID:    tapi.TraceID{0x01},
			SpanID:     tapi.SpanID{0x01},
			TraceFlags: 0x1,
			Remote:     true,
		}))

	_, _ = tracer.Start(ctx, "span")
	check(t, collect(), sampledLive(1), remoteStarted(1))
}

func chainStarted(parent, child int64) metricdata.Metrics {
	noParentSet := attribute.NewSet(
		otelconv.SDKSpanStarted{}.AttrSpanParentOrigin(
			otelconv.SpanParentOriginNone,
		),
		otelconv.SDKSpanStarted{}.AttrSpanSamplingResult(
			otelconv.SpanSamplingResultRecordAndSample,
		),
	)
	localSet := attribute.NewSet(
		otelconv.SDKSpanStarted{}.AttrSpanParentOrigin(
			otelconv.SpanParentOriginLocal,
		),
		otelconv.SDKSpanStarted{}.AttrSpanSamplingResult(
			otelconv.SpanSamplingResultRecordAndSample,
		),
	)
	return started(dPt(noParentSet, parent), dPt(localSet, child))
}

func TestTracerLocalParent(t *testing.T) {
	collect := setup(t)
	tracer := trace.NewTracerProvider().Tracer(t.Name())

	ctx, parent := tracer.Start(t.Context(), "parent")
	_, child := tracer.Start(ctx, "child")

	check(t, collect(), sampledLive(2), chainStarted(1, 1))

	child.End()
	parent.End()

	check(t, collect(), sampledLive(0), chainStarted(1, 1))
}

func TestNewTracerObservabilityDisabled(t *testing.T) {
	// Do not set OTEL_GO_X_OBSERVABILITY
	tracer, err := observ.NewTracer()
	assert.NoError(t, err)
	assert.False(t, tracer.Enabled())
}

func TestNewTracerErrors(t *testing.T) {
	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	orig := otel.GetMeterProvider()
	t.Cleanup(func() { otel.SetMeterProvider(orig) })

	mp := &errMeterProvider{err: assert.AnError}
	otel.SetMeterProvider(mp)

	_, err := observ.NewTracer()
	require.ErrorIs(t, err, assert.AnError, "new instrument errors")

	assert.ErrorContains(t, err, "span live metric")
	assert.ErrorContains(t, err, "span started metric")
}

func BenchmarkTracer(b *testing.B) {
	b.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	orig := otel.GetMeterProvider()
	b.Cleanup(func() { otel.SetMeterProvider(orig) })

	// Ensure deterministic benchmark by using noop meter.
	otel.SetMeterProvider(noop.NewMeterProvider())

	tracer, err := observ.NewTracer()
	require.NoError(b, err)
	require.True(b, tracer.Enabled())

	t := otel.GetTracerProvider().Tracer(b.Name())
	ctx, span := t.Start(b.Context(), "parent")
	psc := span.SpanContext()
	span.End()

	b.Run("SpanStarted", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				tracer.SpanStarted(ctx, psc, span)
			}
		})
	})

	b.Run("SpanLive", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				tracer.SpanLive(ctx, span)
			}
		})
	})

	b.Run("SpanEnded", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				tracer.SpanEnded(ctx, span)
			}
		})
	})
}

func BenchmarkNewTracer(b *testing.B) {
	b.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	orig := otel.GetMeterProvider()
	b.Cleanup(func() { otel.SetMeterProvider(orig) })

	// Ensure deterministic benchmark by using noop meter.
	otel.SetMeterProvider(noop.NewMeterProvider())

	var tracer observ.Tracer

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		tracer, _ = observ.NewTracer()
	}

	_ = tracer
}
