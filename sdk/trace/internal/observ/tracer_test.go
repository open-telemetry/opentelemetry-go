// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package observ_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	mapi "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/internal/observ"
	"go.opentelemetry.io/otel/semconv/v1.37.0/otelconv"
	tapi "go.opentelemetry.io/otel/trace"
)

func setup(t *testing.T) func() metricdata.ScopeMetrics {
	t.Setenv("OTEL_GO_X_SELF_OBSERVABILITY", "true")

	orig := otel.GetMeterProvider()
	t.Cleanup(func() { otel.SetMeterProvider(orig) })

	reader := metric.NewManualReader()
	mp := metric.NewMeterProvider(metric.WithReader(reader))
	otel.SetMeterProvider(mp)

	return func() metricdata.ScopeMetrics {
		var got metricdata.ResourceMetrics
		require.NoError(t, reader.Collect(context.Background(), &got))
		if len(got.ScopeMetrics) != 1 {
			return metricdata.ScopeMetrics{}
		}
		return got.ScopeMetrics[0]
	}
}

func scopeMetrics(metrics ...metricdata.Metrics) metricdata.ScopeMetrics {
	return metricdata.ScopeMetrics{
		Scope: instrumentation.Scope{
			Name:      observ.ScopeName,
			Version:   sdk.Version(),
			SchemaURL: observ.SchemaURL,
		},
		Metrics: metrics,
	}
}

func check(t *testing.T, got metricdata.ScopeMetrics, want ...metricdata.Metrics) {
	o := []metricdatatest.Option{
		metricdatatest.IgnoreTimestamp(),
		metricdatatest.IgnoreExemplars(),
	}
	metricdatatest.AssertEqual(t, scopeMetrics(want...), got, o...)
}

func dPt(set attribute.Set, value int64) metricdata.DataPoint[int64] {
	return metricdata.DataPoint[int64]{Attributes: set, Value: value}
}

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

func sampledStarted() metricdata.Metrics {
	set := attribute.NewSet(
		otelconv.SDKSpanStarted{}.AttrSpanParentOrigin(
			otelconv.SpanParentOriginNone,
		),
		otelconv.SDKSpanStarted{}.AttrSpanSamplingResult(
			otelconv.SpanSamplingResultRecordAndSample,
		),
	)
	return started(dPt(set, 1))
}

func TestTracer(t *testing.T) {
	collect := setup(t)
	tracer := trace.NewTracerProvider().Tracer(t.Name())

	_, span := tracer.Start(context.Background(), "span")
	check(t, collect(), sampledLive(1), sampledStarted())

	span.End()
	check(t, collect(), sampledLive(0), sampledStarted())
}

func dropStarted() metricdata.Metrics {
	set := attribute.NewSet(
		otelconv.SDKSpanStarted{}.AttrSpanParentOrigin(
			otelconv.SpanParentOriginNone,
		),
		otelconv.SDKSpanStarted{}.AttrSpanSamplingResult(
			otelconv.SpanSamplingResultDrop,
		),
	)
	return started(dPt(set, 1))
}

func TestTracerNonRecording(t *testing.T) {
	collect := setup(t)
	tracer := trace.NewTracerProvider(
		trace.WithSampler(trace.NeverSample()),
	).Tracer(t.Name())

	_, _ = tracer.Start(context.Background(), "span")
	check(t, collect(), dropStarted())
}

func recLive(value int64) metricdata.Metrics {
	set := attribute.NewSet(
		otelconv.SDKSpanLive{}.AttrSpanSamplingResult(
			otelconv.SpanSamplingResultRecordOnly,
		),
	)
	return live(dPt(set, value))
}

func recStarted() metricdata.Metrics {
	set := attribute.NewSet(
		otelconv.SDKSpanStarted{}.AttrSpanParentOrigin(
			otelconv.SpanParentOriginNone,
		),
		otelconv.SDKSpanStarted{}.AttrSpanSamplingResult(
			otelconv.SpanSamplingResultRecordOnly,
		),
	)
	return started(dPt(set, 1))
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

func remoteStarted() metricdata.Metrics {
	set := attribute.NewSet(
		otelconv.SDKSpanStarted{}.AttrSpanParentOrigin(
			otelconv.SpanParentOriginRemote,
		),
		otelconv.SDKSpanStarted{}.AttrSpanSamplingResult(
			otelconv.SpanSamplingResultRecordAndSample,
		),
	)
	return started(dPt(set, 1))
}

func TestTracerRecordOnly(t *testing.T) {
	collect := setup(t)
	tracer := trace.NewTracerProvider(
		trace.WithSampler(recOnly{}),
	).Tracer(t.Name())

	_, _ = tracer.Start(context.Background(), "span")
	check(t, collect(), recLive(1), recStarted())
}

func TestTracerRemoteParent(t *testing.T) {
	collect := setup(t)
	tracer := trace.NewTracerProvider().Tracer(t.Name())

	ctx := tapi.ContextWithRemoteSpanContext(
		context.Background(),
		tapi.NewSpanContext(tapi.SpanContextConfig{
			TraceID:    tapi.TraceID{0x01},
			SpanID:     tapi.SpanID{0x01},
			TraceFlags: 0x1,
			Remote:     true,
		}))

	_, _ = tracer.Start(ctx, "span")
	check(t, collect(), sampledLive(1), remoteStarted())
}

func chainStarted() metricdata.Metrics {
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
	return started(dPt(noParentSet, 1), dPt(localSet, 1))
}

func TestTracerLocalParent(t *testing.T) {
	collect := setup(t)
	tracer := trace.NewTracerProvider().Tracer(t.Name())

	ctx, parent := tracer.Start(context.Background(), "parent")
	_, child := tracer.Start(ctx, "child")

	check(t, collect(), sampledLive(2), chainStarted())

	child.End()
	parent.End()

	check(t, collect(), sampledLive(0), chainStarted())
}

func TestNewTracerObservabilityDisabled(t *testing.T) {
	// Do not set OTEL_GO_X_SELF_OBSERVABILITY
	tracer, err := observ.NewTracer()
	assert.NoError(t, err)
	assert.Nil(t, tracer)
}

type errMeterProvider struct {
	mapi.MeterProvider

	err error
}

func (m *errMeterProvider) Meter(string, ...mapi.MeterOption) mapi.Meter {
	return &errMeter{err: m.err}
}

type errMeter struct {
	mapi.Meter

	err error
}

func (m *errMeter) Int64UpDownCounter(string, ...mapi.Int64UpDownCounterOption) (mapi.Int64UpDownCounter, error) {
	return nil, m.err
}

func (m *errMeter) Int64Counter(string, ...mapi.Int64CounterOption) (mapi.Int64Counter, error) {
	return nil, m.err
}

func (m *errMeter) Int64ObservableUpDownCounter(
	string,
	...mapi.Int64ObservableUpDownCounterOption,
) (mapi.Int64ObservableUpDownCounter, error) {
	return nil, m.err
}

func (m *errMeter) RegisterCallback(mapi.Callback, ...mapi.Observable) (mapi.Registration, error) {
	return nil, m.err
}

func TestNewTracerErrors(t *testing.T) {
	t.Setenv("OTEL_GO_X_SELF_OBSERVABILITY", "true")

	orig := otel.GetMeterProvider()
	t.Cleanup(func() { otel.SetMeterProvider(orig) })

	mp := &errMeterProvider{err: assert.AnError}
	otel.SetMeterProvider(mp)

	_, err := observ.NewTracer()
	require.ErrorIs(t, err, assert.AnError, "new instrument errors")

	assert.ErrorContains(t, err, "span live metric")
	assert.ErrorContains(t, err, "span started metric")
}
