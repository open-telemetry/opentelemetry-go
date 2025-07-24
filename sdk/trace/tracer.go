// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace // import "go.opentelemetry.io/otel/sdk/trace"

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/trace/internal/x"
	semconv "go.opentelemetry.io/otel/semconv/v1.36.0"
	"go.opentelemetry.io/otel/semconv/v1.36.0/otelconv"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/embedded"
)

type tracer struct {
	embedded.Tracer

	provider             *TracerProvider
	instrumentationScope instrumentation.Scope

	selfObservabilityEnabled bool
	spanLiveMetric           otelconv.SDKSpanLive
	spanStartedMetric        otelconv.SDKSpanStarted
}

var _ trace.Tracer = &tracer{}

func (tr *tracer) initSelfObservability() {
	if !x.SelfObservability.Enabled() {
		return
	}

	tr.selfObservabilityEnabled = true
	mp := otel.GetMeterProvider()
	m := mp.Meter("go.opentelemetry.io/otel/sdk/trace",
		metric.WithInstrumentationVersion(sdk.Version()),
		metric.WithSchemaURL(semconv.SchemaURL))

	var err error
	if tr.spanLiveMetric, err = otelconv.NewSDKSpanLive(m); err != nil {
		otel.Handle(err)
	}
	if tr.spanStartedMetric, err = otelconv.NewSDKSpanStarted(m); err != nil {
		otel.Handle(err)
	}
}

// Start starts a Span and returns it along with a context containing it.
//
// The Span is created with the provided name and as a child of any existing
// span context found in the passed context. The created Span will be
// configured appropriately by any SpanOption passed.
func (tr *tracer) Start(
	ctx context.Context,
	name string,
	options ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	config := trace.NewSpanStartConfig(options...)

	if ctx == nil {
		// Prevent trace.ContextWithSpan from panicking.
		ctx = context.Background()
	}

	// For local spans created by this SDK, track child span count.
	if p := trace.SpanFromContext(ctx); p != nil {
		if sdkSpan, ok := p.(*recordingSpan); ok {
			sdkSpan.addChild()
		}
	}

	s := tr.newSpan(ctx, name, &config)
	if tr.selfObservabilityEnabled {
		// Check if the span has a parent span and set the origin attribute accordingly.
		var attrParentOrigin attribute.KeyValue
		if psc := trace.SpanContextFromContext(ctx); psc.IsValid() {
			if psc.IsRemote() {
				attrParentOrigin = tr.spanStartedMetric.AttrSpanParentOrigin(otelconv.SpanParentOriginRemote)
			} else {
				attrParentOrigin = tr.spanStartedMetric.AttrSpanParentOrigin(otelconv.SpanParentOriginLocal)
			}
		} else {
			attrParentOrigin = tr.spanStartedMetric.AttrSpanParentOrigin(otelconv.SpanParentOriginNone)
		}

		// Determine the sampling result and create the corresponding attribute.
		var attrSamplingResult attribute.KeyValue
		if s.SpanContext().IsSampled() && s.IsRecording() {
			attrSamplingResult = tr.spanStartedMetric.AttrSpanSamplingResult(
				otelconv.SpanSamplingResultRecordAndSample,
			)
		} else if s.IsRecording() {
			attrSamplingResult = tr.spanStartedMetric.AttrSpanSamplingResult(otelconv.SpanSamplingResultRecordOnly)
		} else {
			attrSamplingResult = tr.spanStartedMetric.AttrSpanSamplingResult(otelconv.SpanSamplingResultDrop)
		}

		tr.spanStartedMetric.Add(context.Background(), 1, attrParentOrigin, attrSamplingResult)
	}

	if rw, ok := s.(ReadWriteSpan); ok && s.IsRecording() {
		sps := tr.provider.getSpanProcessors()
		for _, sp := range sps {
			sp.sp.OnStart(ctx, rw)
		}
	}
	if rtt, ok := s.(runtimeTracer); ok {
		ctx = rtt.runtimeTrace(ctx)
	}

	return trace.ContextWithSpan(ctx, s), s
}

type runtimeTracer interface {
	// runtimeTrace starts a "runtime/trace".Task for the span and
	// returns a context containing the task.
	runtimeTrace(ctx context.Context) context.Context
}

// newSpan returns a new configured span.
func (tr *tracer) newSpan(ctx context.Context, name string, config *trace.SpanConfig) trace.Span {
	// If told explicitly to make this a new root use a zero value SpanContext
	// as a parent which contains an invalid trace ID and is not remote.
	var psc trace.SpanContext
	if config.NewRoot() {
		ctx = trace.ContextWithSpanContext(ctx, psc)
	} else {
		psc = trace.SpanContextFromContext(ctx)
	}

	// If there is a valid parent trace ID, use it to ensure the continuity of
	// the trace. Always generate a new span ID so other components can rely
	// on a unique span ID, even if the Span is non-recording.
	var tid trace.TraceID
	var sid trace.SpanID
	if !psc.TraceID().IsValid() {
		tid, sid = tr.provider.idGenerator.NewIDs(ctx)
	} else {
		tid = psc.TraceID()
		sid = tr.provider.idGenerator.NewSpanID(ctx, tid)
	}

	samplingResult := tr.provider.sampler.ShouldSample(SamplingParameters{
		ParentContext: ctx,
		TraceID:       tid,
		Name:          name,
		Kind:          config.SpanKind(),
		Attributes:    config.Attributes(),
		Links:         config.Links(),
	})

	scc := trace.SpanContextConfig{
		TraceID:    tid,
		SpanID:     sid,
		TraceState: samplingResult.Tracestate,
	}
	if isSampled(samplingResult) {
		scc.TraceFlags = psc.TraceFlags() | trace.FlagsSampled
	} else {
		scc.TraceFlags = psc.TraceFlags() &^ trace.FlagsSampled
	}
	sc := trace.NewSpanContext(scc)

	if !isRecording(samplingResult) {
		return tr.newNonRecordingSpan(sc)
	}
	return tr.newRecordingSpan(psc, sc, name, samplingResult, config)
}

// newRecordingSpan returns a new configured recordingSpan.
func (tr *tracer) newRecordingSpan(
	psc, sc trace.SpanContext,
	name string,
	sr SamplingResult,
	config *trace.SpanConfig,
) *recordingSpan {
	startTime := config.Timestamp()
	if startTime.IsZero() {
		startTime = time.Now()
	}

	s := &recordingSpan{
		// Do not pre-allocate the attributes slice here! Doing so will
		// allocate memory that is likely never going to be used, or if used,
		// will be over-sized. The default Go compiler has been tested to
		// dynamically allocate needed space very well. Benchmarking has shown
		// it to be more performant than what we can predetermine here,
		// especially for the common use case of few to no added
		// attributes.

		parent:      psc,
		spanContext: sc,
		spanKind:    trace.ValidateSpanKind(config.SpanKind()),
		name:        name,
		startTime:   startTime,
		events:      newEvictedQueueEvent(tr.provider.spanLimits.EventCountLimit),
		links:       newEvictedQueueLink(tr.provider.spanLimits.LinkCountLimit),
		tracer:      tr,
	}

	for _, l := range config.Links() {
		s.AddLink(l)
	}

	s.SetAttributes(sr.Attributes...)
	s.SetAttributes(config.Attributes()...)

	if tr.selfObservabilityEnabled {
		// Determine the sampling result and create the corresponding attribute.
		var attrSamplingResult attribute.KeyValue
		if s.spanContext.IsSampled() {
			attrSamplingResult = tr.spanLiveMetric.AttrSpanSamplingResult(
				otelconv.SpanSamplingResultRecordAndSample,
			)
		} else {
			attrSamplingResult = tr.spanLiveMetric.AttrSpanSamplingResult(otelconv.SpanSamplingResultRecordOnly)
		}

		tr.spanLiveMetric.Add(context.Background(), 1, attrSamplingResult)
	}

	return s
}

// newNonRecordingSpan returns a new configured nonRecordingSpan.
func (tr *tracer) newNonRecordingSpan(sc trace.SpanContext) nonRecordingSpan {
	return nonRecordingSpan{tracer: tr, sc: sc}
}
