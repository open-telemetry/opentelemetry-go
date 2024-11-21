// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package opentracing

import (
	"context"
	"fmt"
	"testing"

	ot "github.com/opentracing/opentracing-go"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/bridge/opentracing/internal"
	"go.opentelemetry.io/otel/trace"
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

		otel.SetTracerProvider(otelProvider)
		ot.SetGlobalTracer(otTracer)

		tc.setup(t, mockOtelTracer)
		tc.run(t, ctx)
		tc.check(t, mockOtelTracer)
	}
}

// simple test

type simpleTest struct {
	traceID trace.TraceID
	spanIDs []trace.SpanID
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
	traceID trace.TraceID
	spanIDs []trace.SpanID

	recordedCurrentOtelSpanIDs []trace.SpanID
	recordedActiveOTSpanIDs    []trace.SpanID
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
	spanID := trace.SpanContextFromContext(ctx).SpanID()
	cast.recordedCurrentOtelSpanIDs = append(cast.recordedCurrentOtelSpanIDs, spanID)

	spanID = trace.SpanID{}
	if bridgeSpan, ok := ot.SpanFromContext(ctx).(*bridgeSpan); ok {
		spanID = bridgeSpan.otelSpan.SpanContext().SpanID()
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
				key:   "third",
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
				key:   "third",
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

	m, err := baggage.NewMemberRaw(otelKey, value)
	if err != nil {
		t.Error(err)
		return ctx
	}
	b, err := baggage.FromContext(ctx).SetMember(m)
	if err != nil {
		t.Error(err)
		return ctx
	}
	ctx = baggage.ContextWithBaggage(ctx, b)

	otRecording := make(map[string]string)
	otSpan.Context().ForeachBaggageItem(func(key, value string) bool {
		otRecording[key] = value
		return true
	})
	otelRecording := make(map[string]string)
	for _, m := range baggage.FromContext(ctx).Members() {
		otelRecording[m.Key()] = m.Value()
	}
	bio.recordedOTBaggage = append(bio.recordedOTBaggage, otRecording)
	bio.recordedOtelBaggage = append(bio.recordedOtelBaggage, otelRecording)
	return ctx
}

func generateBaggageKeys(key string) (otKey, otelKey string) {
	otKey, otelKey = key+"-Ot", key+"-Otel"
	return
}

// helpers

func checkTraceAndSpans(t *testing.T, tracer *internal.MockTracer, expectedTraceID trace.TraceID, expectedSpanIDs []trace.SpanID) {
	expectedSpanCount := len(expectedSpanIDs)

	// reverse spanIDs, since first span ID belongs to root, that
	// finishes last
	spanIDs := make([]trace.SpanID, len(expectedSpanIDs))
	copy(spanIDs, expectedSpanIDs)
	reverse(len(spanIDs), func(i, j int) {
		spanIDs[i], spanIDs[j] = spanIDs[j], spanIDs[i]
	})
	// the last finished span has no parent
	parentSpanIDs := append(spanIDs[1:], trace.SpanID{})

	sks := map[trace.SpanID]trace.SpanKind{
		{125}: trace.SpanKindProducer,
		{124}: trace.SpanKindInternal,
		{123}: trace.SpanKindClient,
	}

	if len(tracer.FinishedSpans) != expectedSpanCount {
		t.Errorf("Expected %d finished spans, got %d", expectedSpanCount, len(tracer.FinishedSpans))
	}
	for idx, span := range tracer.FinishedSpans {
		sctx := span.SpanContext()
		if sctx.TraceID() != expectedTraceID {
			t.Errorf("Expected trace ID %v in span %d (%d), got %v", expectedTraceID, idx, sctx.SpanID(), sctx.TraceID())
		}
		expectedSpanID := spanIDs[idx]
		expectedParentSpanID := parentSpanIDs[idx]
		if sctx.SpanID() != expectedSpanID {
			t.Errorf("Expected finished span %d to have span ID %d, but got %d", idx, expectedSpanID, sctx.SpanID())
		}
		if span.ParentSpanID != expectedParentSpanID {
			t.Errorf("Expected finished span %d (span ID: %d) to have parent span ID %d, but got %d", idx, sctx.SpanID(), expectedParentSpanID, span.ParentSpanID)
		}
		if span.SpanKind != sks[span.SpanContext().SpanID()] {
			t.Errorf("Expected finished span %d (span ID: %d) to have span.kind to be '%v' but was '%v'", idx, sctx.SpanID(), sks[span.SpanContext().SpanID()], span.SpanKind)
		}
	}
}

func reverse(length int, swap func(i, j int)) {
	for left, right := 0, length-1; left < right; left, right = left+1, right-1 {
		swap(left, right)
	}
}

func simpleTraceID() trace.TraceID {
	return [16]byte{123, 42}
}

func simpleSpanIDs(count int) []trace.SpanID {
	base := []trace.SpanID{
		{123},
		{124},
		{125},
		{126},
		{127},
		{128},
	}
	return base[:count]
}

func runOtelOTOtel(t *testing.T, ctx context.Context, name string, callback func(*testing.T, context.Context) context.Context) {
	tr := otel.Tracer("")
	ctx, span := tr.Start(ctx, fmt.Sprintf("%s_Otel_OTOtel", name), trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()
	ctx = callback(t, ctx)
	func(ctx2 context.Context) {
		span, ctx2 := ot.StartSpanFromContext(ctx2, fmt.Sprintf("%sOtel_OT_Otel", name))
		defer span.Finish()
		ctx2 = callback(t, ctx2)
		func(ctx3 context.Context) {
			ctx3, span := tr.Start(ctx3, fmt.Sprintf("%sOtelOT_Otel_", name), trace.WithSpanKind(trace.SpanKindProducer))
			defer span.End()
			_ = callback(t, ctx3)
		}(ctx2)
	}(ctx)
}

func runOTOtelOT(t *testing.T, ctx context.Context, name string, callback func(*testing.T, context.Context) context.Context) {
	tr := otel.Tracer("")
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

func TestOtTagToOTelAttrCheckTypeConversions(t *testing.T) {
	tableTest := []struct {
		key               string
		value             interface{}
		expectedValueType attribute.Type
	}{
		{
			key:               "bool to bool",
			value:             true,
			expectedValueType: attribute.BOOL,
		},
		{
			key:               "int to int64",
			value:             123,
			expectedValueType: attribute.INT64,
		},
		{
			key:               "uint to string",
			value:             uint(1234),
			expectedValueType: attribute.STRING,
		},
		{
			key:               "int32 to int64",
			value:             int32(12345),
			expectedValueType: attribute.INT64,
		},
		{
			key:               "uint32 to int64",
			value:             uint32(123456),
			expectedValueType: attribute.INT64,
		},
		{
			key:               "int64 to int64",
			value:             int64(1234567),
			expectedValueType: attribute.INT64,
		},
		{
			key:               "uint64 to string",
			value:             uint64(12345678),
			expectedValueType: attribute.STRING,
		},
		{
			key:               "float32 to float64",
			value:             float32(3.14),
			expectedValueType: attribute.FLOAT64,
		},
		{
			key:               "float64 to float64",
			value:             float64(3.14),
			expectedValueType: attribute.FLOAT64,
		},
		{
			key:               "string to string",
			value:             "string_value",
			expectedValueType: attribute.STRING,
		},
		{
			key:               "unexpected type to string",
			value:             struct{}{},
			expectedValueType: attribute.STRING,
		},
	}

	for _, test := range tableTest {
		got := otTagToOTelAttr(test.key, test.value)
		if test.expectedValueType != got.Value.Type() {
			t.Errorf("Expected type %s, but got %s after conversion '%v' value",
				test.expectedValueType,
				got.Value.Type(),
				test.value)
		}
	}
}
