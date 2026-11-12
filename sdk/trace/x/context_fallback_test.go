// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/sdk/trace/x"
	"go.opentelemetry.io/otel/trace"
)

const fallbackTraceParent = "00-0123456789abcdef0123456789abcdef-0123456789abcdef-01"

func fallbackPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
}

func fallbackCarrier() propagation.MapCarrier {
	return propagation.MapCarrier{
		"traceparent": fallbackTraceParent,
		"tracestate":  "vendor=value",
		"baggage":     "startup=incoming,shared=incoming;property=value",
	}
}

func fallbackProvider(t *testing.T, opts ...sdktrace.TracerProviderOption) *sdktrace.TracerProvider {
	t.Helper()
	provider := sdktrace.NewTracerProvider(opts...)
	t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.WithoutCancel(t.Context()))) })
	return provider
}

type fallbackSampler func(sdktrace.SamplingParameters) sdktrace.SamplingResult

func (s fallbackSampler) ShouldSample(p sdktrace.SamplingParameters) sdktrace.SamplingResult {
	return s(p)
}

func (fallbackSampler) Description() string { return "context fallback test sampler" }

func baggageContext(t *testing.T, ctx context.Context, entries ...string) context.Context {
	t.Helper()
	members := make([]baggage.Member, 0, len(entries)/2)
	for i := 0; i < len(entries); i += 2 {
		member, err := baggage.NewMember(entries[i], entries[i+1])
		require.NoError(t, err)
		members = append(members, member)
	}
	b, err := baggage.New(members...)
	require.NoError(t, err)
	return baggage.ContextWithBaggage(ctx, b)
}

func baggageValues(ctx context.Context) map[string]string {
	values := make(map[string]string)
	for _, member := range baggage.FromContext(ctx).Members() {
		values[member.Key()] = member.Value()
	}
	return values
}

func otherSpanContext(t *testing.T, remote bool) trace.SpanContext {
	t.Helper()
	tid, err := trace.TraceIDFromHex("fedcba9876543210fedcba9876543210")
	require.NoError(t, err)
	sid, err := trace.SpanIDFromHex("fedcba9876543210")
	require.NoError(t, err)
	state, err := trace.ParseTraceState("other=parent")
	require.NoError(t, err)
	return trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: tid, SpanID: sid, TraceFlags: trace.FlagsSampled, TraceState: state, Remote: remote,
	})
}

func TestContextFallbackParent(t *testing.T) {
	carrier := propagation.MapCarrier{
		"traceparent": "00-0123456789abcdef0123456789abcdef-0123456789abcdef-01",
		"tracestate":  "vendor=value",
	}
	want := trace.SpanContextFromContext(propagation.TraceContext{}.Extract(t.Context(), carrier))
	recorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(
		x.WithContextFallback(propagation.TraceContext{}, carrier),
		sdktrace.WithSpanProcessor(recorder),
	)
	t.Cleanup(func() { require.NoError(t, provider.Shutdown(context.WithoutCancel(t.Context()))) })
	_, span := provider.Tracer("test").Start(t.Context(), "job")
	span.End()
	require.Len(t, recorder.Ended(), 1)
	assert.True(t, want.Equal(recorder.Ended()[0].Parent()))
	assert.Equal(t, want.TraceID(), span.SpanContext().TraceID())
	assert.Equal(t, "vendor=value", span.SpanContext().TraceState().String())
}

func TestContextFallbackBeforeSampling(t *testing.T) {
	carrier := fallbackCarrier()
	want := trace.SpanContextFromContext(fallbackPropagator().Extract(t.Context(), carrier))
	var sampled bool
	processor := &fallbackContextProcessor{SpanRecorder: tracetest.NewSpanRecorder()}
	sampler := fallbackSampler(func(p sdktrace.SamplingParameters) sdktrace.SamplingResult {
		sampled = true
		assert.True(t, want.Equal(trace.SpanContextFromContext(p.ParentContext)))
		assert.Equal(t, want.TraceID(), p.TraceID)
		assert.Equal(t, "incoming", baggage.FromContext(p.ParentContext).Member("startup").Value())
		assert.Equal(t, "job", p.Name)
		assert.Equal(t, trace.SpanKindConsumer, p.Kind)
		assert.Equal(t, []attribute.KeyValue{attribute.String("job", "work")}, p.Attributes)
		return sdktrace.SamplingResult{Decision: sdktrace.RecordAndSample, Tracestate: want.TraceState()}
	})
	provider := fallbackProvider(t,
		x.WithContextFallback(fallbackPropagator(), carrier),
		sdktrace.WithSampler(sampler), sdktrace.WithSpanProcessor(processor),
	)
	ctx, span := provider.Tracer("test").Start(t.Context(), "job",
		trace.WithSpanKind(trace.SpanKindConsumer), trace.WithAttributes(attribute.String("job", "work")),
	)
	span.End()
	assert.True(t, sampled)
	assert.Equal(t, "incoming", baggage.FromContext(ctx).Member("startup").Value())
	require.Len(t, processor.contexts, 1)
	assert.True(t, want.Equal(trace.SpanContextFromContext(processor.contexts[0])))
	assert.Equal(t, baggageValues(ctx), baggageValues(processor.contexts[0]))
	require.Len(t, processor.Ended(), 1)
	assert.True(t, want.Equal(processor.Ended()[0].Parent()))
	assert.Equal(t, "vendor=value", span.SpanContext().TraceState().String())
}

type fallbackIDGenerator struct {
	newIDs    func(context.Context) (trace.TraceID, trace.SpanID)
	newSpanID func(context.Context, trace.TraceID) trace.SpanID
}

func (g fallbackIDGenerator) NewIDs(ctx context.Context) (trace.TraceID, trace.SpanID) {
	return g.newIDs(ctx)
}

func (g fallbackIDGenerator) NewSpanID(ctx context.Context, id trace.TraceID) trace.SpanID {
	return g.newSpanID(ctx, id)
}

func TestContextFallbackBeforeIDGeneration(t *testing.T) {
	for _, tc := range []struct {
		name string
		mode x.ContextFallbackMode
	}{{"parent", x.ContextFallbackParent}, {"link", x.ContextFallbackLink}} {
		t.Run(tc.name, func(t *testing.T) {
			incoming := trace.SpanContextFromContext(fallbackPropagator().Extract(t.Context(), fallbackCarrier()))
			called := false
			checkContext := func(ctx context.Context) {
				called = true
				assert.Equal(t, "incoming", baggage.FromContext(ctx).Member("startup").Value())
			}
			generator := fallbackIDGenerator{
				newIDs: func(ctx context.Context) (trace.TraceID, trace.SpanID) {
					checkContext(ctx)
					assert.Equal(t, x.ContextFallbackLink, tc.mode)
					assert.True(t, trace.SpanContextFromContext(ctx).Equal(trace.SpanContext{}))
					return trace.TraceID{3}, trace.SpanID{4}
				},
				newSpanID: func(ctx context.Context, id trace.TraceID) trace.SpanID {
					checkContext(ctx)
					assert.Equal(t, x.ContextFallbackParent, tc.mode)
					assert.Equal(t, incoming.TraceID(), id)
					assert.True(t, incoming.Equal(trace.SpanContextFromContext(ctx)))
					return trace.SpanID{4}
				},
			}
			provider := fallbackProvider(t,
				x.WithContextFallback(fallbackPropagator(), fallbackCarrier(), x.WithContextFallbackMode(tc.mode)),
				sdktrace.WithIDGenerator(generator), sdktrace.WithSampler(sdktrace.NeverSample()),
			)
			_, span := provider.Tracer("test").Start(t.Context(), "job")
			span.End()
			assert.True(t, called)
		})
	}
}

type fallbackContextProcessor struct {
	*tracetest.SpanRecorder
	contexts []context.Context
}

func (p *fallbackContextProcessor) OnStart(ctx context.Context, span sdktrace.ReadWriteSpan) {
	p.contexts = append(p.contexts, ctx)
	p.SpanRecorder.OnStart(ctx, span)
}

func TestContextFallbackParentBasedSampling(t *testing.T) {
	for _, tc := range []struct {
		name        string
		traceparent string
		sampled     bool
	}{
		{"sampled", fallbackTraceParent, true},
		{"not sampled", "00-0123456789abcdef0123456789abcdef-0123456789abcdef-00", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			carrier := propagation.MapCarrier{"traceparent": tc.traceparent}
			incoming := trace.SpanContextFromContext(propagation.TraceContext{}.Extract(t.Context(), carrier))
			provider := fallbackProvider(t, x.WithContextFallback(propagation.TraceContext{}, carrier),
				sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.NeverSample())))
			ctx, span := provider.Tracer("test").Start(t.Context(), "job")
			assert.Equal(t, incoming.TraceID(), span.SpanContext().TraceID())
			assert.Equal(t, tc.sampled, span.IsRecording())
			assert.Equal(t, tc.sampled, span.SpanContext().IsSampled())
			_, child := provider.Tracer("test").Start(ctx, "child")
			assert.Equal(t, span.SpanContext().TraceID(), child.SpanContext().TraceID())
			assert.Equal(t, tc.sampled, child.IsRecording())
			child.End()
			span.End()
		})
	}
}

func TestContextFallbackPrecedence(t *testing.T) {
	for _, mode := range []x.ContextFallbackMode{x.ContextFallbackParent, x.ContextFallbackLink} {
		for _, tc := range []struct {
			name    string
			parent  trace.SpanContext
			newRoot bool
		}{
			{"local parent", otherSpanContext(t, false), false},
			{"remote parent", otherSpanContext(t, true), false},
			{"explicit root", trace.SpanContext{}, true},
			{"explicit root with parent", otherSpanContext(t, true), true},
			{"partial trace context", trace.NewSpanContext(trace.SpanContextConfig{TraceID: otherSpanContext(t, false).TraceID()}), false},
			{"partial span context", trace.NewSpanContext(trace.SpanContextConfig{SpanID: otherSpanContext(t, false).SpanID()}), false},
			{"trace state only", trace.NewSpanContext(trace.SpanContextConfig{TraceState: otherSpanContext(t, false).TraceState()}), false},
			{"trace flags only", trace.NewSpanContext(trace.SpanContextConfig{TraceFlags: trace.FlagsSampled}), false},
			{"remote flag only", trace.NewSpanContext(trace.SpanContextConfig{Remote: true}), false},
		} {
			t.Run(fmt.Sprintf("%s/%d", tc.name, mode), func(t *testing.T) {
				var seen sdktrace.SamplingParameters
				sampler := fallbackSampler(func(p sdktrace.SamplingParameters) sdktrace.SamplingResult {
					seen = p
					return sdktrace.SamplingResult{Decision: sdktrace.RecordAndSample}
				})
				provider := fallbackProvider(t,
					x.WithContextFallback(fallbackPropagator(), fallbackCarrier(), x.WithContextFallbackMode(mode)),
					sdktrace.WithSampler(sampler),
				)
				parent := baggageContext(t, trace.ContextWithSpanContext(t.Context(), tc.parent), "application", "kept")
				var options []trace.SpanStartOption
				if tc.newRoot {
					options = append(options, trace.WithNewRoot())
				}
				ctx, span := provider.Tracer("test").Start(parent, "job", options...)
				span.End()
				assert.Equal(t, map[string]string{"application": "kept"}, baggageValues(ctx))
				assert.Empty(t, seen.Links)
				if tc.newRoot {
					assert.False(t, trace.SpanContextFromContext(seen.ParentContext).IsValid())
					assert.NotEqual(t, tc.parent.TraceID(), span.SpanContext().TraceID())
				} else {
					assert.True(t, tc.parent.Equal(trace.SpanContextFromContext(seen.ParentContext)))
					if tc.parent.TraceID().IsValid() {
						assert.Equal(t, tc.parent.TraceID(), span.SpanContext().TraceID())
					}
				}
			})
		}
	}
}

func TestContextFallbackLinks(t *testing.T) {
	carrier := fallbackCarrier()
	incoming := trace.SpanContextFromContext(fallbackPropagator().Extract(t.Context(), carrier))
	for _, duplicate := range []bool{false, true} {
		t.Run(fmt.Sprintf("duplicate=%t", duplicate), func(t *testing.T) {
			first := trace.Link{
				SpanContext: otherSpanContext(t, false),
				Attributes:  []attribute.KeyValue{attribute.Bool("existing", true)},
			}
			original := []trace.Link{first}
			if duplicate {
				// A link to the same span is sufficient even if its metadata differs.
				original = append(original, trace.Link{
					SpanContext: incoming.WithRemote(false).WithTraceFlags(0),
					Attributes:  []attribute.KeyValue{attribute.String("keep", "metadata")},
				})
			}
			want := append([]trace.Link(nil), original...)
			if !duplicate {
				want = append(want, trace.Link{SpanContext: incoming})
			}
			var seenLinks []trace.Link
			sampler := fallbackSampler(func(p sdktrace.SamplingParameters) sdktrace.SamplingResult {
				assert.False(t, trace.SpanContextFromContext(p.ParentContext).IsValid())
				assert.NotEqual(t, incoming.TraceID(), p.TraceID)
				seenLinks = p.Links
				return sdktrace.SamplingResult{Decision: sdktrace.RecordAndSample}
			})
			recorder := tracetest.NewSpanRecorder()
			provider := fallbackProvider(t,
				x.WithContextFallback(fallbackPropagator(), carrier, x.WithContextFallbackMode(x.ContextFallbackLink)),
				sdktrace.WithSampler(sampler), sdktrace.WithSpanProcessor(recorder),
			)
			option := trace.WithLinks(original...)
			for range 2 {
				ctx, span := provider.Tracer("test").Start(t.Context(), "job", option)
				span.End()
				assert.Equal(t, want, seenLinks)
				assert.Equal(t, "incoming", baggage.FromContext(ctx).Member("startup").Value())
			}
			config := trace.NewSpanStartConfig(option)
			assert.Equal(t, original, config.Links())
			require.Len(t, recorder.Ended(), 2)
			for _, span := range recorder.Ended() {
				require.Len(t, span.Links(), len(want))
				for i, link := range span.Links() {
					assert.True(t, want[i].SpanContext.Equal(link.SpanContext))
					assert.Equal(t, want[i].Attributes, link.Attributes)
				}
				assert.False(t, span.Parent().IsValid())
			}
		})
	}
}

func TestContextFallbackLinksBeforeLimits(t *testing.T) {
	carrier := fallbackCarrier()
	incoming := trace.SpanContextFromContext(fallbackPropagator().Extract(t.Context(), carrier))
	var seen []trace.Link
	sampler := fallbackSampler(func(p sdktrace.SamplingParameters) sdktrace.SamplingResult {
		seen = p.Links
		return sdktrace.SamplingResult{Decision: sdktrace.RecordAndSample}
	})
	limits := sdktrace.NewSpanLimits()
	limits.LinkCountLimit = 1
	recorder := tracetest.NewSpanRecorder()
	provider := fallbackProvider(t,
		x.WithContextFallback(fallbackPropagator(), carrier, x.WithContextFallbackMode(x.ContextFallbackLink)),
		sdktrace.WithSampler(sampler), sdktrace.WithRawSpanLimits(limits), sdktrace.WithSpanProcessor(recorder),
	)
	_, span := provider.Tracer("test").
		Start(t.Context(), "job", trace.WithLinks(trace.Link{SpanContext: otherSpanContext(t, false)}))
	span.End()
	require.Len(t, seen, 2)
	assert.True(t, incoming.Equal(seen[1].SpanContext))
	require.Len(t, recorder.Ended(), 1)
	assert.Len(t, recorder.Ended()[0].Links(), 1)
	assert.Equal(t, 1, recorder.Ended()[0].DroppedLinks())
}

func TestContextFallbackBaggageInheritance(t *testing.T) {
	for _, mode := range []x.ContextFallbackMode{x.ContextFallbackParent, x.ContextFallbackLink} {
		t.Run(fmt.Sprint(mode), func(t *testing.T) {
			provider := fallbackProvider(
				t,
				x.WithContextFallback(fallbackPropagator(), fallbackCarrier(), x.WithContextFallbackMode(mode)),
			)
			tracer := provider.Tracer("test")
			original := baggageContext(t, t.Context(), "shared", "application", "application", "kept")
			ctx, span := tracer.Start(original, "job")
			defer span.End()
			assert.Equal(
				t,
				map[string]string{"startup": "incoming", "shared": "application", "application": "kept"},
				baggageValues(ctx),
			)
			assert.Equal(t, map[string]string{"shared": "application", "application": "kept"}, baggageValues(original))

			modified := baggage.FromContext(ctx).DeleteMember("startup")
			member, err := baggage.NewMember("shared", "changed")
			require.NoError(t, err)
			modified, err = modified.SetMember(member)
			require.NoError(t, err)
			modifiedCtx := baggage.ContextWithBaggage(ctx, modified)
			childCtx, child := tracer.Start(modifiedCtx, "child")
			child.End()
			assert.Equal(t, map[string]string{"shared": "changed", "application": "kept"}, baggageValues(childCtx))
			clearedCtx, cleared := tracer.Start(baggage.ContextWithoutBaggage(ctx), "cleared")
			cleared.End()
			assert.Empty(t, baggageValues(clearedCtx), "descendants must not resurrect removed startup baggage")

			repeatedCtx, repeated := tracer.Start(original, "another root")
			repeated.End()
			assert.Equal(
				t,
				baggageValues(ctx),
				baggageValues(repeatedCtx),
				"independent roots receive configured defaults again",
			)
		})
	}
}

func TestContextFallbackBaggageOnly(t *testing.T) {
	for _, value := range []string{"", "malformed", "00-00000000000000000000000000000000-0123456789abcdef-01"} {
		t.Run(value, func(t *testing.T) {
			carrier := fallbackCarrier()
			carrier["traceparent"] = value
			provider := fallbackProvider(t, x.WithContextFallback(fallbackPropagator(), carrier))
			ctx, span := provider.Tracer("test").Start(t.Context(), "job")
			span.End()
			assert.True(t, span.SpanContext().IsValid())
			assert.Equal(t, "incoming", baggage.FromContext(ctx).Member("startup").Value())
			member := baggage.FromContext(ctx).Member("shared")
			require.Len(t, member.Properties(), 1)
			assert.Equal(t, "property", member.Properties()[0].Key())
			_, child := provider.Tracer("test").Start(ctx, "child")
			child.End()
			assert.Equal(t, span.SpanContext().TraceID(), child.SpanContext().TraceID())
		})
	}
}

func TestContextFallbackBaggageLimitsPreserveApplication(t *testing.T) {
	for _, entries := range [][]string{
		{"a", strings.Repeat("a", 4090), "b", strings.Repeat("b", 4090)},
		func() []string {
			entries := make([]string, 0, 128)
			for i := range 64 {
				entries = append(entries, fmt.Sprintf("application%d", i), "kept")
			}
			return entries
		}(),
	} {
		t.Run(fmt.Sprint(len(entries)), func(t *testing.T) {
			original := baggageContext(t, t.Context(), entries...)
			provider := fallbackProvider(t, x.WithContextFallback(fallbackPropagator(), fallbackCarrier()))
			ctx, span := provider.Tracer("test").Start(original, "job")
			span.End()
			assert.Equal(t, baggageValues(original), baggageValues(ctx))
		})
	}
}

func TestContextFallbackPreservesCallerContext(t *testing.T) {
	type key struct{}
	for _, mode := range []x.ContextFallbackMode{x.ContextFallbackParent, x.ContextFallbackLink} {
		t.Run(fmt.Sprint(mode), func(t *testing.T) {
			parent, cancel := context.WithDeadline(
				context.WithValue(t.Context(), key{}, "application"),
				time.Now().Add(time.Hour),
			)
			defer cancel()
			provider := fallbackProvider(
				t,
				x.WithContextFallback(fallbackPropagator(), fallbackCarrier(), x.WithContextFallbackMode(mode)),
			)
			ctx, span := provider.Tracer("test").Start(parent, "job")
			defer span.End()
			assert.Equal(t, "application", ctx.Value(key{}))
			wantDeadline, wantOK := parent.Deadline()
			deadline, ok := ctx.Deadline()
			assert.Equal(t, wantDeadline, deadline)
			assert.Equal(t, wantOK, ok)
			cancel()
			assert.ErrorIs(t, ctx.Err(), context.Canceled)
			select {
			case <-ctx.Done():
			default:
				t.Fatal("fallback context did not retain cancellation")
			}
		})
	}
}

type fallbackTestPropagator struct {
	extract func(context.Context, propagation.TextMapCarrier) context.Context
	calls   atomic.Int64
}

func (*fallbackTestPropagator) Inject(context.Context, propagation.TextMapCarrier) {}

func (p *fallbackTestPropagator) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	p.calls.Add(1)
	return p.extract(ctx, carrier)
}

func (*fallbackTestPropagator) Fields() []string { return []string{"custom-context"} }

func TestContextFallbackExtractsOnceUsingConfiguredPropagator(t *testing.T) {
	type extractedKey struct{}
	want := otherSpanContext(t, false)
	carrier := propagation.MapCarrier{"custom-context": "first"}
	propagator := &fallbackTestPropagator{
		extract: func(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
			assert.Nil(t, ctx.Done())
			assert.Equal(t, "first", carrier.Get("custom-context"))
			ctx = trace.ContextWithSpanContext(ctx, want)
			ctx = baggageContext(t, ctx, "custom", carrier.Get("custom-context"))
			return context.WithValue(ctx, extractedKey{}, "not retained")
		},
	}
	option := x.WithContextFallback(propagator, carrier)
	assert.EqualValues(t, 1, propagator.calls.Load(), "extraction happens at option construction")
	carrier["custom-context"] = "changed"
	for range 2 {
		recorder := tracetest.NewSpanRecorder()
		provider := fallbackProvider(t, option, sdktrace.WithSpanProcessor(recorder))
		for range 2 {
			ctx, span := provider.Tracer("test").Start(t.Context(), "job")
			span.End()
			assert.Equal(t, want.TraceID(), span.SpanContext().TraceID())
			assert.Equal(t, "first", baggage.FromContext(ctx).Member("custom").Value())
			assert.Nil(t, ctx.Value(extractedKey{}), "retain propagation values, not arbitrary extracted context")
		}
		for _, span := range recorder.Ended() {
			assert.True(t, want.Equal(span.Parent()), "custom propagator's remote flag is preserved")
		}
	}
	assert.EqualValues(t, 1, propagator.calls.Load())
}

func TestContextFallbackDisabledAndInvalid(t *testing.T) {
	previous := otel.GetErrorHandler()
	var errorsHandled atomic.Int64
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(error) { errorsHandled.Add(1) }))
	t.Cleanup(func() { otel.SetErrorHandler(previous) })
	for _, tc := range []struct {
		name       string
		propagator propagation.TextMapPropagator
		carrier    propagation.TextMapCarrier
		options    []x.ContextFallbackOption
		wantErrors int64
	}{
		{"none", fallbackPropagator(), fallbackCarrier(), []x.ContextFallbackOption{x.WithContextFallbackMode(x.ContextFallbackNone)}, 0},
		{"none with nil inputs", nil, nil, []x.ContextFallbackOption{x.WithContextFallbackMode(x.ContextFallbackNone)}, 0},
		{"nil propagator", nil, fallbackCarrier(), nil, 1},
		{"nil carrier", fallbackPropagator(), nil, nil, 1},
		{"negative mode", fallbackPropagator(), fallbackCarrier(), []x.ContextFallbackOption{x.WithContextFallbackMode(-1)}, 1},
		{"unknown mode", fallbackPropagator(), fallbackCarrier(), []x.ContextFallbackOption{x.WithContextFallbackMode(42)}, 1},
		{"nil extraction context", &fallbackTestPropagator{extract: func(context.Context, propagation.TextMapCarrier) context.Context { return nil }}, fallbackCarrier(), nil, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			before := errorsHandled.Load()
			provider := fallbackProvider(t, x.WithContextFallback(tc.propagator, tc.carrier, tc.options...))
			ctx, span := provider.Tracer("test").Start(t.Context(), "job")
			span.End()
			assert.Empty(t, baggageValues(ctx))
			incoming := trace.SpanContextFromContext(fallbackPropagator().Extract(t.Context(), fallbackCarrier()))
			assert.NotEqual(t, incoming.TraceID(), span.SpanContext().TraceID())
			assert.Equal(t, tc.wantErrors, errorsHandled.Load()-before)
		})
	}
	propagator := &fallbackTestPropagator{extract: func(context.Context, propagation.TextMapCarrier) context.Context {
		t.Fatal("none must not invoke the propagator")
		return nil
	}}
	x.WithContextFallback(propagator, fallbackCarrier(), x.WithContextFallbackMode(x.ContextFallbackNone))
	assert.Zero(t, propagator.calls.Load())
}

func TestContextFallbackLastOptionWins(t *testing.T) {
	carrier := fallbackCarrier()
	incoming := trace.SpanContextFromContext(fallbackPropagator().Extract(t.Context(), carrier))
	for _, tc := range []struct {
		name      string
		options   []sdktrace.TracerProviderOption
		wantTrace bool
		wantLinks int
	}{
		{"absent", nil, false, 0},
		{"enabled then none", []sdktrace.TracerProviderOption{
			x.WithContextFallback(fallbackPropagator(), carrier),
			x.WithContextFallback(nil, nil, x.WithContextFallbackMode(x.ContextFallbackNone)),
		}, false, 0},
		{"none then enabled", []sdktrace.TracerProviderOption{
			x.WithContextFallback(nil, nil, x.WithContextFallbackMode(x.ContextFallbackNone)),
			x.WithContextFallback(fallbackPropagator(), carrier),
		}, true, 0},
		{"enabled then empty", []sdktrace.TracerProviderOption{
			x.WithContextFallback(fallbackPropagator(), carrier),
			x.WithContextFallback(fallbackPropagator(), propagation.MapCarrier{}),
		}, false, 0},
		{"mode option order", []sdktrace.TracerProviderOption{
			x.WithContextFallback(fallbackPropagator(), carrier,
				x.WithContextFallbackMode(x.ContextFallbackNone), x.WithContextFallbackMode(x.ContextFallbackLink)),
		}, false, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			recorder := tracetest.NewSpanRecorder()
			options := append([]sdktrace.TracerProviderOption{sdktrace.WithSpanProcessor(recorder)}, tc.options...)
			provider := fallbackProvider(t, options...)
			_, span := provider.Tracer("test").Start(t.Context(), "job")
			span.End()
			assert.Equal(t, tc.wantTrace, span.SpanContext().TraceID() == incoming.TraceID())
			require.Len(t, recorder.Ended(), 1)
			assert.Len(t, recorder.Ended()[0].Links(), tc.wantLinks)
		})
	}
}

func TestContextFallbackUnrelatedRequests(t *testing.T) {
	carrier := fallbackCarrier()
	incomingContext := fallbackPropagator().Extract(t.Context(), carrier)
	incoming := trace.SpanContextFromContext(incomingContext)
	for _, persistent := range []bool{false, true} {
		t.Run(fmt.Sprintf("provider fallback=%t", persistent), func(t *testing.T) {
			var options []sdktrace.TracerProviderOption
			if persistent {
				options = append(options, x.WithContextFallback(fallbackPropagator(), carrier))
			}
			provider := fallbackProvider(t, options...)
			tracer := provider.Tracer("test")
			invocation := t.Context()
			if !persistent {
				invocation = incomingContext
			}
			_, job := tracer.Start(invocation, "job")
			job.End()
			assert.Equal(t, incoming.TraceID(), job.SpanContext().TraceID())

			// Empty request headers are also an empty parent context. Provider-wide
			// fallback cannot determine that this is a different invocation.
			request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.test/", nil)
			requestContext := fallbackPropagator().Extract(request.Context(), propagation.HeaderCarrier(request.Header))
			_, server := tracer.Start(requestContext, "unrelated request", trace.WithSpanKind(trace.SpanKindServer))
			server.End()
			assert.Equal(t, persistent, server.SpanContext().TraceID() == incoming.TraceID())
			_, independent := tracer.Start(requestContext, "independent request", trace.WithNewRoot())
			independent.End()
			assert.NotEqual(t, incoming.TraceID(), independent.SpanContext().TraceID())
		})
	}
}

func TestContextFallbackConcurrentStarts(t *testing.T) {
	for _, mode := range []x.ContextFallbackMode{x.ContextFallbackParent, x.ContextFallbackLink} {
		t.Run(fmt.Sprint(mode), func(t *testing.T) {
			recorder := tracetest.NewSpanRecorder()
			carrier := fallbackCarrier()
			incoming := trace.SpanContextFromContext(fallbackPropagator().Extract(t.Context(), carrier))
			provider := fallbackProvider(t,
				x.WithContextFallback(fallbackPropagator(), carrier, x.WithContextFallbackMode(mode)),
				sdktrace.WithSpanProcessor(recorder),
			)
			tracer := provider.Tracer("test")
			parent := baggageContext(t, t.Context(), "shared", "application")
			links := []trace.Link{{SpanContext: otherSpanContext(t, false)}}
			option := trace.WithLinks(links...)
			const workers = 64
			var wg sync.WaitGroup
			for range workers {
				wg.Go(func() {
					ctx, span := tracer.Start(parent, "job", option)
					assert.Equal(
						t,
						map[string]string{"startup": "incoming", "shared": "application"},
						baggageValues(ctx),
					)
					span.End()
				})
			}
			wg.Wait()
			assert.Equal(t, map[string]string{"shared": "application"}, baggageValues(parent))
			config := trace.NewSpanStartConfig(option)
			assert.Equal(t, links, config.Links())
			require.Len(t, recorder.Ended(), workers)
			for _, span := range recorder.Ended() {
				if mode == x.ContextFallbackParent {
					assert.True(t, incoming.Equal(span.Parent()))
					assert.Len(t, span.Links(), 1)
				} else {
					assert.False(t, span.Parent().IsValid())
					require.Len(t, span.Links(), 2)
					assert.True(t, incoming.Equal(span.Links()[1].SpanContext))
				}
			}
		})
	}
}

func TestContextFallbackBoundsCustomBaggage(t *testing.T) {
	for _, tc := range []struct {
		name       string
		members    int
		valueBytes int
	}{
		{"member count", 100, 1},
		{"total bytes", 3, 4090},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var incoming baggage.Baggage
			for i := range tc.members {
				member, err := baggage.NewMember(fmt.Sprintf("k%d", i), strings.Repeat("a", tc.valueBytes))
				require.NoError(t, err)
				incoming, err = incoming.SetMember(member)
				require.NoError(t, err)
			}
			propagator := &fallbackTestPropagator{
				extract: func(ctx context.Context, _ propagation.TextMapCarrier) context.Context {
					return baggage.ContextWithBaggage(ctx, incoming)
				},
			}
			provider := fallbackProvider(t, x.WithContextFallback(propagator, propagation.MapCarrier{}))
			for _, parent := range []context.Context{
				t.Context(), baggageContext(t, t.Context(), "application", "kept"),
			} {
				ctx, span := provider.Tracer("test").Start(parent, "job")
				span.End()
				got := baggage.FromContext(ctx)
				assert.LessOrEqual(t, got.Len(), 64)
				assert.LessOrEqual(t, len(got.String()), 8192)
				assert.Positive(t, got.Len())
				assert.Equal(
					t,
					baggage.FromContext(parent).Member("application").Value(),
					got.Member("application").Value(),
				)
			}
		})
	}
}
