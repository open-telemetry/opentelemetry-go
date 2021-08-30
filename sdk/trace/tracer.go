// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package trace // import "go.opentelemetry.io/otel/sdk/trace"

import (
	"context"
	rt "runtime/trace"
	"time"

	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel/sdk/instrumentation"
)

type tracer struct {
	provider               *TracerProvider
	instrumentationLibrary instrumentation.Library
}

var _ trace.Tracer = &tracer{}

// Start starts a Span and returns it along with a context containing it.
//
// The Span is created with the provided name and as a child of any existing
// span context found in the passed context. The created Span will be
// configured appropriately by any SpanOption passed.
func (tr *tracer) Start(ctx context.Context, name string, options ...trace.SpanStartOption) (context.Context, trace.Span) {
	config := trace.NewSpanStartConfig(options...)

	// For local spans created by this SDK, track child span count.
	if p := trace.SpanFromContext(ctx); p != nil {
		if sdkSpan, ok := p.(*span); ok {
			sdkSpan.addChild()
		}
	}

	s := tr.newSpan(ctx, name, config)
	if s.IsRecording() {
		sps, _ := tr.provider.spanProcessors.Load().(spanProcessorStates)
		for _, sp := range sps {
			sp.sp.OnStart(ctx, s)
		}
	}

	ctx, s.executionTracerTaskEnd = newRuntimeTask(ctx, name)

	return trace.ContextWithSpan(ctx, s), s
}

// newSpan returns a new configured span.
func (tr *tracer) newSpan(ctx context.Context, name string, config *trace.SpanConfig) *span {
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

	s := new(span)
	s.attributes = newAttributesMap(tr.provider.spanLimits.AttributeCountLimit)
	s.events = newEvictedQueue(tr.provider.spanLimits.EventCountLimit)
	s.links = newEvictedQueue(tr.provider.spanLimits.LinkCountLimit)
	s.spanLimits = tr.provider.spanLimits

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
	s.spanContext = trace.NewSpanContext(scc)

	if !isRecording(samplingResult) {
		return s
	}

	startTime := config.Timestamp()
	if startTime.IsZero() {
		startTime = time.Now()
	}
	s.startTime = startTime

	s.spanKind = trace.ValidateSpanKind(config.SpanKind())
	s.name = name
	s.parent = psc
	s.resource = tr.provider.resource
	s.instrumentationLibrary = tr.instrumentationLibrary
	s.tracer = tr

	s.SetAttributes(samplingResult.Attributes...)
	s.SetAttributes(config.Attributes()...)

	for _, l := range config.Links() {
		s.addLink(l)
	}

	return s
}

// newRuntimeTask starts a runtime.Task with the passed name and returns both
// a context containing the task and an end function for the task.
//
// If the runtime tracing is not enabled the original context is returned with
// an empty function to call for end.
func newRuntimeTask(ctx context.Context, name string) (context.Context, func()) {
	if !rt.IsEnabled() {
		// Avoid additional overhead if
		// runtime/trace is not enabled.
		return ctx, func() {}
	}
	nctx, task := rt.NewTask(ctx, name)
	return nctx, task.End

}
