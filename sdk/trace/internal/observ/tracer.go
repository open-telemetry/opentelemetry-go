// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package observ

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/trace/internal/x"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/semconv/v1.37.0/otelconv"
	"go.opentelemetry.io/otel/trace"
)

const (
	// ScopeName is the name of the instrumentation scope.
	ScopeName = "go.opentelemetry.io/otel/sdk/trace/internal/observ"

	// SchemaURL is the schema URL of the instrumentation.
	SchemaURL = semconv.SchemaURL
)

// Tracer is instrumentation for an OTel SDK Tracer.
type Tracer struct {
	spanLiveMetric    otelconv.SDKSpanLive
	spanStartedMetric otelconv.SDKSpanStarted
}

func NewTracer() (*Tracer, error) {
	if !x.SelfObservability.Enabled() {
		return nil, nil
	}
	meter := otel.GetMeterProvider().Meter(
		ScopeName,
		metric.WithInstrumentationVersion(sdk.Version()),
		metric.WithSchemaURL(SchemaURL),
	)

	var err error
	spanLiveMetric, e := otelconv.NewSDKSpanLive(meter)
	if e != nil {
		e = fmt.Errorf("failed to create span live metric: %w", e)
		err = errors.Join(err, e)
	}

	spanStartedMetric, e := otelconv.NewSDKSpanStarted(meter)
	if e != nil {
		e = fmt.Errorf("failed to create span started metric: %w", e)
		err = errors.Join(err, e)
	}
	return &Tracer{
		spanLiveMetric:    spanLiveMetric,
		spanStartedMetric: spanStartedMetric,
	}, err
}

func (t *Tracer) SpanStarted(ctx context.Context, psc trace.SpanContext, span trace.Span) {
	set := spanStartedSet(psc, span)
	t.spanStartedMetric.AddSet(ctx, 1, set)
}

func (t *Tracer) SpanLive(ctx context.Context, span trace.Span) {
	set := spanLiveSet(span.SpanContext().IsSampled())
	t.spanLiveMetric.AddSet(ctx, 1, set)
}

func (t *Tracer) SpanEnded(span trace.Span) {
	// Add the span to the context to ensure the metric is recorded
	// with the correct span context.
	ctx := trace.ContextWithSpan(context.Background(), span)
	set := spanLiveSet(span.SpanContext().IsSampled())
	t.spanLiveMetric.AddSet(ctx, -1, set)
}

type parentState int

const (
	parentStateNoParent parentState = iota
	parentStateLocalParent
	parentStateRemoteParent
)

type samplingState int

const (
	samplingStateDrop samplingState = iota
	samplingStateRecordOnly
	samplingStateRecordAndSample
)

type spanStartedSetKey struct {
	parent   parentState
	sampling samplingState
}

var spanStartedSetCache = map[spanStartedSetKey]attribute.Set{
	{parentStateNoParent, samplingStateDrop}: attribute.NewSet(
		otelconv.SDKSpanStarted{}.AttrSpanParentOrigin(otelconv.SpanParentOriginNone),
		otelconv.SDKSpanStarted{}.AttrSpanSamplingResult(otelconv.SpanSamplingResultDrop),
	),
	{parentStateLocalParent, samplingStateDrop}: attribute.NewSet(
		otelconv.SDKSpanStarted{}.AttrSpanParentOrigin(otelconv.SpanParentOriginLocal),
		otelconv.SDKSpanStarted{}.AttrSpanSamplingResult(otelconv.SpanSamplingResultDrop),
	),
	{parentStateRemoteParent, samplingStateDrop}: attribute.NewSet(
		otelconv.SDKSpanStarted{}.AttrSpanParentOrigin(otelconv.SpanParentOriginRemote),
		otelconv.SDKSpanStarted{}.AttrSpanSamplingResult(otelconv.SpanSamplingResultDrop),
	),

	{parentStateNoParent, samplingStateRecordOnly}: attribute.NewSet(
		otelconv.SDKSpanStarted{}.AttrSpanParentOrigin(otelconv.SpanParentOriginNone),
		otelconv.SDKSpanStarted{}.AttrSpanSamplingResult(otelconv.SpanSamplingResultRecordOnly),
	),
	{parentStateLocalParent, samplingStateRecordOnly}: attribute.NewSet(
		otelconv.SDKSpanStarted{}.AttrSpanParentOrigin(otelconv.SpanParentOriginLocal),
		otelconv.SDKSpanStarted{}.AttrSpanSamplingResult(otelconv.SpanSamplingResultRecordOnly),
	),
	{parentStateRemoteParent, samplingStateRecordOnly}: attribute.NewSet(
		otelconv.SDKSpanStarted{}.AttrSpanParentOrigin(otelconv.SpanParentOriginRemote),
		otelconv.SDKSpanStarted{}.AttrSpanSamplingResult(otelconv.SpanSamplingResultRecordOnly),
	),

	{parentStateNoParent, samplingStateRecordAndSample}: attribute.NewSet(
		otelconv.SDKSpanStarted{}.AttrSpanParentOrigin(otelconv.SpanParentOriginNone),
		otelconv.SDKSpanStarted{}.AttrSpanSamplingResult(otelconv.SpanSamplingResultRecordAndSample),
	),
	{parentStateLocalParent, samplingStateRecordAndSample}: attribute.NewSet(
		otelconv.SDKSpanStarted{}.AttrSpanParentOrigin(otelconv.SpanParentOriginLocal),
		otelconv.SDKSpanStarted{}.AttrSpanSamplingResult(otelconv.SpanSamplingResultRecordAndSample),
	),
	{parentStateRemoteParent, samplingStateRecordAndSample}: attribute.NewSet(
		otelconv.SDKSpanStarted{}.AttrSpanParentOrigin(otelconv.SpanParentOriginRemote),
		otelconv.SDKSpanStarted{}.AttrSpanSamplingResult(otelconv.SpanSamplingResultRecordAndSample),
	),
}

func spanStartedSet(psc trace.SpanContext, span trace.Span) attribute.Set {
	key := spanStartedSetKey{
		parent:   parentStateNoParent,
		sampling: samplingStateDrop,
	}

	if psc.IsValid() {
		if psc.IsRemote() {
			key.parent = parentStateRemoteParent
		} else {
			key.parent = parentStateLocalParent
		}
	}

	if span.IsRecording() {
		if span.SpanContext().IsSampled() {
			key.sampling = samplingStateRecordAndSample
		} else {
			key.sampling = samplingStateRecordOnly
		}
	}

	return spanStartedSetCache[key]
}

type spanLiveSetKey struct {
	sampled bool
}

var spanLiveSetCache = map[spanLiveSetKey]attribute.Set{
	{true}: attribute.NewSet(
		otelconv.SDKSpanLive{}.AttrSpanSamplingResult(
			otelconv.SpanSamplingResultRecordAndSample,
		),
	),
	{false}: attribute.NewSet(
		otelconv.SDKSpanLive{}.AttrSpanSamplingResult(
			otelconv.SpanSamplingResultRecordOnly,
		),
	),
}

func spanLiveSet(sampled bool) attribute.Set {
	key := spanLiveSetKey{sampled: sampled}
	return spanLiveSetCache[key]
}
