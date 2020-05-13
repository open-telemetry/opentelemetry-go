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

package opentracing

import (
	"context"
	"fmt"
	"testing"

	ot "github.com/opentracing/opentracing-go"

	otelcorrelation "go.opentelemetry.io/otel/api/correlation"
	otelglobal "go.opentelemetry.io/otel/api/global"
	otelcore "go.opentelemetry.io/otel/api/kv"
	oteltrace "go.opentelemetry.io/otel/api/trace"

	"go.opentelemetry.io/otel/bridge/opentracing/internal"
)

type mixedAPIsTestCase struct {
	desc string

	setup func(*testing.T, *internal.MockTracer)
	run   func(*testing.T, context.Context)
	check func(*testing.T, *internal.MockTracer)
}

func getMixedAPIsTestCases() []mixedAPIsTestCase {
	st := newSimpleTest()
	cast := newCurrentActiveSpanTest()
	coin := newContextIntactTest()
	bip := newBaggageItemsPreservationTest()
	tm := newTracerMessTest()
	bio := newBaggageInteroperationTest()

	return []mixedAPIsTestCase{
		{
			desc:  "simple otel -> ot -> otel",
			setup: st.setup,
			run:   st.runOtelOTOtel,
			check: st.check,
		},
		{
			desc:  "simple ot -> otel -> ot",
			setup: st.setup,
			run:   st.runOTOtelOT,
			check: st.check,
		},
		{
			desc:  "current/active span otel -> ot -> otel",
			setup: cast.setup,
			run:   cast.runOtelOTOtel,
			check: cast.check,
		},
		{
			desc:  "current/active span ot -> otel -> ot",
			setup: cast.setup,
			run:   cast.runOTOtelOT,
			check: cast.check,
		},
		{
			desc:  "context intact otel -> ot -> otel",
			setup: coin.setup,
			run:   coin.runOtelOTOtel,
			check: coin.check,
		},
		{
			desc:  "context intact ot -> otel -> ot",
			setup: coin.setup,
			run:   coin.runOTOtelOT,
			check: coin.check,
		},
		{
			desc:  "baggage items preservation across layers otel -> ot -> otel",
			setup: bip.setup,
			run:   bip.runOtelOTOtel,
			check: bip.check,
		},
		{
			desc:  "baggage items preservation across layers ot -> otel -> ot",
			setup: bip.setup,
			run:   bip.runOTOtelOT,
			check: bip.check,
		},
		{
			desc:  "consistent tracers otel -> ot -> otel",
			setup: tm.setup,
			run:   tm.runOtelOTOtel,
			check: tm.check,
		},
		{
			desc:  "consistent tracers ot -> otel -> ot",
			setup: tm.setup,
			run:   tm.runOTOtelOT,
			check: tm.check,
		},
		{
			desc:  "baggage items interoperation across layers ot -> otel -> ot",
			setup: bio.setup,
			run:   bio.runOTOtelOT,
			check: bio.check,
		},
		{
			desc:  "baggage items interoperation across layers otel -> ot -> otel",
			setup: bio.setup,
			run:   bio.runOtelOTOtel,
			check: bio.check,
		},
	}
}

func TestMixedAPIs(t *testing.T) {
	for idx, tc := range getMixedAPIsTestCases() {
		t.Logf("Running test case %d: %s", idx, tc.desc)
		mockOtelTracer := internal.NewMockTracer()
		ctx, otTracer, otelProvider := NewTracerPairWithContext(context.Background(), mockOtelTracer)
		otTracer.SetWarningHandler(func(msg string) {
			t.Log(msg)
		})

		otelglobal.SetTraceProvider(otelProvider)
		ot.SetGlobalTracer(otTracer)

		tc.setup(t, mockOtelTracer)
		tc.run(t, ctx)
		tc.check(t, mockOtelTracer)
	}
}

// simple test

type simpleTest struct {
	traceID oteltrace.ID
	spanIDs []oteltrace.SpanID
}

func newSimpleTest() *simpleTest {
	return &simpleTest{
		traceID: simpleTraceID(),
		spanIDs: simpleSpanIDs(3),
	}
}

func (st *simpleTest) setup(t *testing.T, tracer *internal.MockTracer) {
	tracer.SpareTraceIDs = append(tracer.SpareTraceIDs, st.traceID)
	tracer.SpareSpanIDs = append(tracer.SpareSpanIDs, st.spanIDs...)
}

func (st *simpleTest) check(t *testing.T, tracer *internal.MockTracer) {
	checkTraceAndSpans(t, tracer, st.traceID, st.spanIDs)
}

func (st *simpleTest) runOtelOTOtel(t *testing.T, ctx context.Context) {
	runOtelOTOtel(t, ctx, "simple", st.noop)
}

func (st *simpleTest) runOTOtelOT(t *testing.T, ctx context.Context) {
	runOTOtelOT(t, ctx, "simple", st.noop)
}

func (st *simpleTest) noop(t *testing.T, ctx context.Context) context.Context {
	return ctx
}

// current/active span test

type currentActiveSpanTest struct {
	traceID oteltrace.ID
	spanIDs []oteltrace.SpanID

	recordedCurrentOtelSpanIDs []oteltrace.SpanID
	recordedActiveOTSpanIDs    []oteltrace.SpanID
}

func newCurrentActiveSpanTest() *currentActiveSpanTest {
	return &currentActiveSpanTest{
		traceID: simpleTraceID(),
		spanIDs: simpleSpanIDs(3),
	}
}

func (cast *currentActiveSpanTest) setup(t *testing.T, tracer *internal.MockTracer) {
	tracer.SpareTraceIDs = append(tracer.SpareTraceIDs, cast.traceID)
	tracer.SpareSpanIDs = append(tracer.SpareSpanIDs, cast.spanIDs...)

	cast.recordedCurrentOtelSpanIDs = nil
	cast.recordedActiveOTSpanIDs = nil
}

func (cast *currentActiveSpanTest) check(t *testing.T, tracer *internal.MockTracer) {
	checkTraceAndSpans(t, tracer, cast.traceID, cast.spanIDs)
	if len(cast.recordedCurrentOtelSpanIDs) != len(cast.spanIDs) {
		t.Errorf("Expected to have %d recorded Otel current spans, got %d", len(cast.spanIDs), len(cast.recordedCurrentOtelSpanIDs))
	}
	if len(cast.recordedActiveOTSpanIDs) != len(cast.spanIDs) {
		t.Errorf("Expected to have %d recorded OT active spans, got %d", len(cast.spanIDs), len(cast.recordedActiveOTSpanIDs))
	}

	minLen := min(len(cast.recordedCurrentOtelSpanIDs), len(cast.spanIDs))
	minLen = min(minLen, len(cast.recordedActiveOTSpanIDs))
	for i := 0; i < minLen; i++ {
		if cast.recordedCurrentOtelSpanIDs[i] != cast.spanIDs[i] {
			t.Errorf("Expected span idx %d (%d) to be recorded as current span in Otel, got %d", i, cast.spanIDs[i], cast.recordedCurrentOtelSpanIDs[i])
		}
		if cast.recordedActiveOTSpanIDs[i] != cast.spanIDs[i] {
			t.Errorf("Expected span idx %d (%d) to be recorded as active span in OT, got %d", i, cast.spanIDs[i], cast.recordedActiveOTSpanIDs[i])
		}
	}
}

func (cast *currentActiveSpanTest) runOtelOTOtel(t *testing.T, ctx context.Context) {
	runOtelOTOtel(t, ctx, "cast", cast.recordSpans)
}

func (cast *currentActiveSpanTest) runOTOtelOT(t *testing.T, ctx context.Context) {
	runOTOtelOT(t, ctx, "cast", cast.recordSpans)
}

func (cast *currentActiveSpanTest) recordSpans(t *testing.T, ctx context.Context) context.Context {
	spanID := oteltrace.SpanFromContext(ctx).SpanContext().SpanID
	cast.recordedCurrentOtelSpanIDs = append(cast.recordedCurrentOtelSpanIDs, spanID)

	spanID = oteltrace.SpanID{}
	if bridgeSpan, ok := ot.SpanFromContext(ctx).(*bridgeSpan); ok {
		spanID = bridgeSpan.otelSpan.SpanContext().SpanID
	}
	cast.recordedActiveOTSpanIDs = append(cast.recordedActiveOTSpanIDs, spanID)
	return ctx
}

// context intact test

type contextIntactTest struct {
	contextKeyValues []internal.MockContextKeyValue

	recordedContextValues []interface{}
	recordIdx             int
}

type coin1Key struct{}

type coin1Value struct{}

type coin2Key struct{}

type coin2Value struct{}

type coin3Key struct{}

type coin3Value struct{}

func newContextIntactTest() *contextIntactTest {
	return &contextIntactTest{
		contextKeyValues: []internal.MockContextKeyValue{
			{
				Key:   coin1Key{},
				Value: coin1Value{},
			},
			{
				Key:   coin2Key{},
				Value: coin2Value{},
			},
			{
				Key:   coin3Key{},
				Value: coin3Value{},
			},
		},
	}
}

func (coin *contextIntactTest) setup(t *testing.T, tracer *internal.MockTracer) {
	tracer.SpareContextKeyValues = append(tracer.SpareContextKeyValues, coin.contextKeyValues...)

	coin.recordedContextValues = nil
	coin.recordIdx = 0
}

func (coin *contextIntactTest) check(t *testing.T, tracer *internal.MockTracer) {
	if len(coin.recordedContextValues) != len(coin.contextKeyValues) {
		t.Errorf("Expected to have %d recorded context values, got %d", len(coin.contextKeyValues), len(coin.recordedContextValues))
	}

	minLen := min(len(coin.recordedContextValues), len(coin.contextKeyValues))
	for i := 0; i < minLen; i++ {
		key := coin.contextKeyValues[i].Key
		value := coin.contextKeyValues[i].Value
		gotValue := coin.recordedContextValues[i]
		if value != gotValue {
			t.Errorf("Expected value %#v for key %#v, got %#v", value, key, gotValue)
		}
	}
}

func (coin *contextIntactTest) runOtelOTOtel(t *testing.T, ctx context.Context) {
	runOtelOTOtel(t, ctx, "coin", coin.recordValue)
}

func (coin *contextIntactTest) runOTOtelOT(t *testing.T, ctx context.Context) {
	runOTOtelOT(t, ctx, "coin", coin.recordValue)
}

func (coin *contextIntactTest) recordValue(t *testing.T, ctx context.Context) context.Context {
	if coin.recordIdx >= len(coin.contextKeyValues) {
		t.Errorf("Too many steps?")
		return ctx
	}
	key := coin.contextKeyValues[coin.recordIdx].Key
	coin.recordIdx++
	coin.recordedContextValues = append(coin.recordedContextValues, ctx.Value(key))
	return ctx
}

// baggage items preservation test

type bipBaggage struct {
	key   string
	value string
}

type baggageItemsPreservationTest struct {
	baggageItems []bipBaggage

	step            int
	recordedBaggage []map[string]string
}

func newBaggageItemsPreservationTest() *baggageItemsPreservationTest {
	return &baggageItemsPreservationTest{
		baggageItems: []bipBaggage{
			{
				key:   "First",
				value: "one",
			},
			{
				key:   "Second",
				value: "two",
			},
			{
				key:   "Third",
				value: "three",
			},
		},
	}
}

func (bip *baggageItemsPreservationTest) setup(t *testing.T, tracer *internal.MockTracer) {
	bip.step = 0
	bip.recordedBaggage = nil
}

func (bip *baggageItemsPreservationTest) check(t *testing.T, tracer *internal.MockTracer) {
	if len(bip.recordedBaggage) != len(bip.baggageItems) {
		t.Errorf("Expected %d recordings, got %d", len(bip.baggageItems), len(bip.recordedBaggage))
	}
	minLen := min(len(bip.recordedBaggage), len(bip.baggageItems))

	for i := 0; i < minLen; i++ {
		recordedItems := bip.recordedBaggage[i]
		if len(recordedItems) != i+1 {
			t.Errorf("Expected %d recorded baggage items in recording %d, got %d", i+1, i+1, len(bip.recordedBaggage[i]))
		}
		minItemLen := min(len(bip.baggageItems), i+1)
		for j := 0; j < minItemLen; j++ {
			expectedItem := bip.baggageItems[j]
			if gotValue, ok := recordedItems[expectedItem.key]; !ok {
				t.Errorf("Missing baggage item %q in recording %d", expectedItem.key, i+1)
			} else if gotValue != expectedItem.value {
				t.Errorf("Expected recorded baggage item %q in recording %d + 1to be %q, got %q", expectedItem.key, i, expectedItem.value, gotValue)
			} else {
				delete(recordedItems, expectedItem.key)
			}
		}
		for key, value := range recordedItems {
			t.Errorf("Unexpected baggage item in recording %d: %q -> %q", i+1, key, value)
		}
	}
}

func (bip *baggageItemsPreservationTest) runOtelOTOtel(t *testing.T, ctx context.Context) {
	runOtelOTOtel(t, ctx, "bip", bip.addAndRecordBaggage)
}

func (bip *baggageItemsPreservationTest) runOTOtelOT(t *testing.T, ctx context.Context) {
	runOTOtelOT(t, ctx, "bip", bip.addAndRecordBaggage)
}

func (bip *baggageItemsPreservationTest) addAndRecordBaggage(t *testing.T, ctx context.Context) context.Context {
	if bip.step >= len(bip.baggageItems) {
		t.Errorf("Too many steps?")
		return ctx
	}
	span := ot.SpanFromContext(ctx)
	if span == nil {
		t.Errorf("No active OpenTracing span")
		return ctx
	}
	idx := bip.step
	bip.step++
	span.SetBaggageItem(bip.baggageItems[idx].key, bip.baggageItems[idx].value)
	sctx := span.Context()
	recording := make(map[string]string)
	sctx.ForeachBaggageItem(func(key, value string) bool {
		recording[key] = value
		return true
	})
	bip.recordedBaggage = append(bip.recordedBaggage, recording)
	return ctx
}

// tracer mess test

type tracerMessTest struct {
	recordedOTSpanTracers   []ot.Tracer
	recordedOtelSpanTracers []oteltrace.Tracer
}

func newTracerMessTest() *tracerMessTest {
	return &tracerMessTest{
		recordedOTSpanTracers:   nil,
		recordedOtelSpanTracers: nil,
	}
}

func (tm *tracerMessTest) setup(t *testing.T, tracer *internal.MockTracer) {
	tm.recordedOTSpanTracers = nil
	tm.recordedOtelSpanTracers = nil
}

func (tm *tracerMessTest) check(t *testing.T, tracer *internal.MockTracer) {
	globalOtTracer := ot.GlobalTracer()
	globalOtelTracer := otelglobal.Tracer("")
	if len(tm.recordedOTSpanTracers) != 3 {
		t.Errorf("Expected 3 recorded OpenTracing tracers from spans, got %d", len(tm.recordedOTSpanTracers))
	}
	if len(tm.recordedOtelSpanTracers) != 3 {
		t.Errorf("Expected 3 recorded OpenTelemetry tracers from spans, got %d", len(tm.recordedOtelSpanTracers))
	}
	for idx, tracer := range tm.recordedOTSpanTracers {
		if tracer != globalOtTracer {
			t.Errorf("Expected OpenTracing tracer %d to be the same as global tracer (%#v), but got %#v", idx, globalOtTracer, tracer)
		}
	}
	for idx, tracer := range tm.recordedOtelSpanTracers {
		if tracer != globalOtelTracer {
			t.Errorf("Expected OpenTelemetry tracer %d to be the same as global tracer (%#v), but got %#v", idx, globalOtelTracer, tracer)
		}
	}
}

func (tm *tracerMessTest) runOtelOTOtel(t *testing.T, ctx context.Context) {
	runOtelOTOtel(t, ctx, "tm", tm.recordTracers)
}

func (tm *tracerMessTest) runOTOtelOT(t *testing.T, ctx context.Context) {
	runOTOtelOT(t, ctx, "tm", tm.recordTracers)
}

func (tm *tracerMessTest) recordTracers(t *testing.T, ctx context.Context) context.Context {
	otSpan := ot.SpanFromContext(ctx)
	if otSpan == nil {
		t.Errorf("No current OpenTracing span?")
	} else {
		tm.recordedOTSpanTracers = append(tm.recordedOTSpanTracers, otSpan.Tracer())
	}

	otelSpan := oteltrace.SpanFromContext(ctx)
	tm.recordedOtelSpanTracers = append(tm.recordedOtelSpanTracers, otelSpan.Tracer())
	return ctx
}

// baggage interoperation test

type baggageInteroperationTest struct {
	baggageItems []bipBaggage

	step                int
	recordedOTBaggage   []map[string]string
	recordedOtelBaggage []map[string]string
}

func newBaggageInteroperationTest() *baggageInteroperationTest {
	return &baggageInteroperationTest{
		baggageItems: []bipBaggage{
			{
				key:   "First",
				value: "one",
			},
			{
				key:   "Second",
				value: "two",
			},
			{
				key:   "Third",
				value: "three",
			},
		},
	}
}

func (bio *baggageInteroperationTest) setup(t *testing.T, tracer *internal.MockTracer) {
	bio.step = 0
	bio.recordedOTBaggage = nil
	bio.recordedOtelBaggage = nil
}

func (bio *baggageInteroperationTest) check(t *testing.T, tracer *internal.MockTracer) {
	checkBIORecording(t, "OT", bio.baggageItems, bio.recordedOTBaggage)
	checkBIORecording(t, "Otel", bio.baggageItems, bio.recordedOtelBaggage)
}

func checkBIORecording(t *testing.T, apiDesc string, initialItems []bipBaggage, recordings []map[string]string) {
	// expect recordings count to equal the number of initial
	// items

	// each recording should have a duplicated item from initial
	// items, one with OT suffix, another one with Otel suffix

	// expect each subsequent recording to have two more items, up
	// to double of the count of the initial items

	if len(initialItems) != len(recordings) {
		t.Errorf("Expected %d recordings from %s, got %d", len(initialItems), apiDesc, len(recordings))
	}
	minRecLen := min(len(initialItems), len(recordings))
	for i := 0; i < minRecLen; i++ {
		recordedItems := recordings[i]
		expectedItemsInStep := (i + 1) * 2
		if expectedItemsInStep != len(recordedItems) {
			t.Errorf("Expected %d recorded items in recording %d from %s, got %d", expectedItemsInStep, i, apiDesc, len(recordedItems))
		}
		recordedItemsCopy := make(map[string]string, len(recordedItems))
		for k, v := range recordedItems {
			recordedItemsCopy[k] = v
		}
		for j := 0; j < i+1; j++ {
			otKey, otelKey := generateBaggageKeys(initialItems[j].key)
			value := initialItems[j].value
			for _, k := range []string{otKey, otelKey} {
				if v, ok := recordedItemsCopy[k]; ok {
					if value != v {
						t.Errorf("Expected value %s under key %s in recording %d from %s, got %s", value, k, i, apiDesc, v)
					}
					delete(recordedItemsCopy, k)
				} else {
					t.Errorf("Missing key %s in recording %d from %s", k, i, apiDesc)
				}
			}
		}
		for k, v := range recordedItemsCopy {
			t.Errorf("Unexpected key-value pair %s = %s in recording %d from %s", k, v, i, apiDesc)
		}
	}
}

func (bio *baggageInteroperationTest) runOtelOTOtel(t *testing.T, ctx context.Context) {
	runOtelOTOtel(t, ctx, "bio", bio.addAndRecordBaggage)
}

func (bio *baggageInteroperationTest) runOTOtelOT(t *testing.T, ctx context.Context) {
	runOTOtelOT(t, ctx, "bio", bio.addAndRecordBaggage)
}

func (bio *baggageInteroperationTest) addAndRecordBaggage(t *testing.T, ctx context.Context) context.Context {
	if bio.step >= len(bio.baggageItems) {
		t.Errorf("Too many steps?")
		return ctx
	}
	otSpan := ot.SpanFromContext(ctx)
	if otSpan == nil {
		t.Errorf("No active OpenTracing span")
		return ctx
	}
	idx := bio.step
	bio.step++
	key := bio.baggageItems[idx].key
	otKey, otelKey := generateBaggageKeys(key)
	value := bio.baggageItems[idx].value

	otSpan.SetBaggageItem(otKey, value)
	ctx = otelcorrelation.NewContext(ctx, otelcore.String(otelKey, value))

	otRecording := make(map[string]string)
	otSpan.Context().ForeachBaggageItem(func(key, value string) bool {
		otRecording[key] = value
		return true
	})
	otelRecording := make(map[string]string)
	otelcorrelation.MapFromContext(ctx).Foreach(func(kv otelcore.KeyValue) bool {
		otelRecording[string(kv.Key)] = kv.Value.Emit()
		return true
	})
	bio.recordedOTBaggage = append(bio.recordedOTBaggage, otRecording)
	bio.recordedOtelBaggage = append(bio.recordedOtelBaggage, otelRecording)
	return ctx
}

func generateBaggageKeys(key string) (otKey, otelKey string) {
	otKey, otelKey = key+"-Ot", key+"-Otel"
	return
}

// helpers

func checkTraceAndSpans(t *testing.T, tracer *internal.MockTracer, expectedTraceID oteltrace.ID, expectedSpanIDs []oteltrace.SpanID) {
	expectedSpanCount := len(expectedSpanIDs)

	// reverse spanIDs, since first span ID belongs to root, that
	// finishes last
	spanIDs := make([]oteltrace.SpanID, len(expectedSpanIDs))
	copy(spanIDs, expectedSpanIDs)
	reverse(len(spanIDs), func(i, j int) {
		spanIDs[i], spanIDs[j] = spanIDs[j], spanIDs[i]
	})
	// the last finished span has no parent
	parentSpanIDs := append(spanIDs[1:], oteltrace.SpanID{})

	sks := map[oteltrace.SpanID]oteltrace.SpanKind{
		{125}: oteltrace.SpanKindProducer,
		{124}: oteltrace.SpanKindInternal,
		{123}: oteltrace.SpanKindClient,
	}

	if len(tracer.FinishedSpans) != expectedSpanCount {
		t.Errorf("Expected %d finished spans, got %d", expectedSpanCount, len(tracer.FinishedSpans))
	}
	for idx, span := range tracer.FinishedSpans {
		sctx := span.SpanContext()
		if sctx.TraceID != expectedTraceID {
			t.Errorf("Expected trace ID %v in span %d (%d), got %v", expectedTraceID, idx, sctx.SpanID, sctx.TraceID)
		}
		expectedSpanID := spanIDs[idx]
		expectedParentSpanID := parentSpanIDs[idx]
		if sctx.SpanID != expectedSpanID {
			t.Errorf("Expected finished span %d to have span ID %d, but got %d", idx, expectedSpanID, sctx.SpanID)
		}
		if span.ParentSpanID != expectedParentSpanID {
			t.Errorf("Expected finished span %d (span ID: %d) to have parent span ID %d, but got %d", idx, sctx.SpanID, expectedParentSpanID, span.ParentSpanID)
		}
		if span.SpanKind != sks[span.SpanContext().SpanID] {
			t.Errorf("Expected finished span %d (span ID: %d) to have span.kind to be '%v' but was '%v'", idx, sctx.SpanID, sks[span.SpanContext().SpanID], span.SpanKind)
		}
	}
}

func reverse(length int, swap func(i, j int)) {
	for left, right := 0, length-1; left < right; left, right = left+1, right-1 {
		swap(left, right)
	}
}

func simpleTraceID() oteltrace.ID {
	return [16]byte{123, 42}
}

func simpleSpanIDs(count int) []oteltrace.SpanID {
	base := []oteltrace.SpanID{
		{123},
		{124},
		{125},
		{126},
		{127},
		{128},
	}
	return base[:count]
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func runOtelOTOtel(t *testing.T, ctx context.Context, name string, callback func(*testing.T, context.Context) context.Context) {
	tr := otelglobal.Tracer("")
	ctx, span := tr.Start(ctx, fmt.Sprintf("%s_Otel_OTOtel", name), oteltrace.WithSpanKind(oteltrace.SpanKindClient))
	defer span.End()
	ctx = callback(t, ctx)
	func(ctx2 context.Context) {
		span, ctx2 := ot.StartSpanFromContext(ctx2, fmt.Sprintf("%sOtel_OT_Otel", name))
		defer span.Finish()
		ctx2 = callback(t, ctx2)
		func(ctx3 context.Context) {
			ctx3, span := tr.Start(ctx3, fmt.Sprintf("%sOtelOT_Otel_", name), oteltrace.WithSpanKind(oteltrace.SpanKindProducer))
			defer span.End()
			_ = callback(t, ctx3)
		}(ctx2)
	}(ctx)
}

func runOTOtelOT(t *testing.T, ctx context.Context, name string, callback func(*testing.T, context.Context) context.Context) {
	tr := otelglobal.Tracer("")
	span, ctx := ot.StartSpanFromContext(ctx, fmt.Sprintf("%s_OT_OtelOT", name), ot.Tag{Key: "span.kind", Value: "client"})
	defer span.Finish()
	ctx = callback(t, ctx)
	func(ctx2 context.Context) {
		ctx2, span := tr.Start(ctx2, fmt.Sprintf("%sOT_Otel_OT", name))
		defer span.End()
		ctx2 = callback(t, ctx2)
		func(ctx3 context.Context) {
			span, ctx3 := ot.StartSpanFromContext(ctx3, fmt.Sprintf("%sOTOtel_OT_", name), ot.Tag{Key: "span.kind", Value: "producer"})
			defer span.Finish()
			_ = callback(t, ctx3)
		}(ctx2)
	}(ctx)
}
