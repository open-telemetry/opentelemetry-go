// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package opentracing

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"

	ot "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/bridge/opentracing/internal"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type testOnlyTextMapReader struct{}

func newTestOnlyTextMapReader() *testOnlyTextMapReader {
	return &testOnlyTextMapReader{}
}

func (t *testOnlyTextMapReader) ForeachKey(handler func(key string, val string) error) error {
	_ = handler("key1", "val1")
	_ = handler("key2", "val2")

	return nil
}

type testOnlyTextMapWriter struct {
	m map[string]string
}

func newTestOnlyTextMapWriter() *testOnlyTextMapWriter {
	return &testOnlyTextMapWriter{m: map[string]string{}}
}

func (t *testOnlyTextMapWriter) Set(key, val string) {
	t.m[key] = val
}

type testTextMapReaderAndWriter struct {
	*testOnlyTextMapReader
	*testOnlyTextMapWriter
}

func newTestTextMapReaderAndWriter() *testTextMapReaderAndWriter {
	return &testTextMapReaderAndWriter{
		testOnlyTextMapReader: newTestOnlyTextMapReader(),
		testOnlyTextMapWriter: newTestOnlyTextMapWriter(),
	}
}

func TestTextMapWrapper_New(t *testing.T) {
	_, err := newTextMapWrapperForExtract(newTestOnlyTextMapReader())
	assert.NoError(t, err)

	_, err = newTextMapWrapperForExtract(newTestOnlyTextMapWriter())
	assert.ErrorIs(t, err, ot.ErrInvalidCarrier)

	_, err = newTextMapWrapperForExtract(newTestTextMapReaderAndWriter())
	assert.NoError(t, err)

	_, err = newTextMapWrapperForInject(newTestOnlyTextMapWriter())
	assert.NoError(t, err)

	_, err = newTextMapWrapperForInject(newTestOnlyTextMapReader())
	assert.ErrorIs(t, err, ot.ErrInvalidCarrier)

	_, err = newTextMapWrapperForInject(newTestTextMapReaderAndWriter())
	assert.NoError(t, err)
}

func TestTextMapWrapper_action(t *testing.T) {
	testExtractFunc := func(carrier propagation.TextMapCarrier) {
		str := carrier.Keys()
		assert.Len(t, str, 2)
		assert.Contains(t, str, "key1", "key2")

		assert.Equal(t, "val1", carrier.Get("key1"))
		assert.Equal(t, "val2", carrier.Get("key2"))
	}

	testInjectFunc := func(carrier propagation.TextMapCarrier) {
		carrier.Set("key1", "val1")
		carrier.Set("key2", "val2")

		wrap, ok := carrier.(*textMapWrapper)
		assert.True(t, ok)

		writer, ok := wrap.TextMapWriter.(*testOnlyTextMapWriter)
		if ok {
			assert.Contains(t, writer.m, "key1", "key2", "val1", "val2")
			return
		}

		writer2, ok := wrap.TextMapWriter.(*testTextMapReaderAndWriter)
		assert.True(t, ok)
		assert.Contains(t, writer2.m, "key1", "key2", "val1", "val2")
	}

	onlyWriter, err := newTextMapWrapperForExtract(newTestOnlyTextMapReader())
	assert.NoError(t, err)
	testExtractFunc(onlyWriter)

	onlyReader, err := newTextMapWrapperForInject(&testOnlyTextMapWriter{m: map[string]string{}})
	assert.NoError(t, err)
	testInjectFunc(onlyReader)

	both, err := newTextMapWrapperForExtract(newTestTextMapReaderAndWriter())
	assert.NoError(t, err)
	testExtractFunc(both)

	both, err = newTextMapWrapperForInject(newTestTextMapReaderAndWriter())
	assert.NoError(t, err)
	testInjectFunc(both)
}

var (
	testHeader               = "test-trace-id"
	traceID    trace.TraceID = [16]byte{byte(10)}
	spanID     trace.SpanID  = [8]byte{byte(11)}
)

type testTextMapPropagator struct{}

func (t testTextMapPropagator) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	carrier.Set(testHeader, strings.Join([]string{traceID.String(), spanID.String()}, ":"))

	// Test for panic
	_ = carrier.Get("test")
	_ = carrier.Keys()
}

func (t testTextMapPropagator) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	traces := carrier.Get(testHeader)

	str := strings.Split(traces, ":")
	if len(str) != 2 {
		return ctx
	}

	exist := false

	for _, key := range carrier.Keys() {
		if strings.EqualFold(testHeader, key) {
			exist = true

			break
		}
	}

	if !exist {
		return ctx
	}

	var (
		traceID, _ = trace.TraceIDFromHex(str[0])
		spanID, _  = trace.SpanIDFromHex(str[1])
		sc         = trace.NewSpanContext(trace.SpanContextConfig{
			TraceID: traceID,
			SpanID:  spanID,
		})
	)

	// Test for panic
	carrier.Set("key", "val")

	return trace.ContextWithRemoteSpanContext(ctx, sc)
}

func (t testTextMapPropagator) Fields() []string {
	return []string{"test"}
}

// textMapCarrier  Implemented propagation.TextMapCarrier interface.
type textMapCarrier struct {
	m map[string]string
}

var _ propagation.TextMapCarrier = (*textMapCarrier)(nil)

func newTextCarrier() *textMapCarrier {
	return &textMapCarrier{m: map[string]string{}}
}

func (t *textMapCarrier) Get(key string) string {
	return t.m[key]
}

func (t *textMapCarrier) Set(key string, value string) {
	t.m[key] = value
}

func (t *textMapCarrier) Keys() []string {
	str := make([]string, 0, len(t.m))

	for key := range t.m {
		str = append(str, key)
	}

	return str
}

// testTextMapReader only implemented opentracing.TextMapReader interface.
type testTextMapReader struct {
	m *map[string]string
}

func newTestTextMapReader(m *map[string]string) *testTextMapReader {
	return &testTextMapReader{m: m}
}

func (t *testTextMapReader) ForeachKey(handler func(key string, val string) error) error {
	for key, val := range *t.m {
		if err := handler(key, val); err != nil {
			return err
		}
	}

	return nil
}

// testTextMapWriter only implemented opentracing.TextMapWriter interface.
type testTextMapWriter struct {
	m *map[string]string
}

func newTestTextMapWriter(m *map[string]string) *testTextMapWriter {
	return &testTextMapWriter{m: m}
}

func (t *testTextMapWriter) Set(key, val string) {
	(*t.m)[key] = val
}

type samplable interface {
	IsSampled() bool
}

func TestBridgeTracer_ExtractAndInject(t *testing.T) {
	bridge := NewBridgeTracer()
	bridge.SetTextMapPropagator(new(testTextMapPropagator))

	tmc := newTextCarrier()
	shareMap := map[string]string{}
	otTextMap := ot.TextMapCarrier{}
	httpHeader := ot.HTTPHeadersCarrier(http.Header{})

	testCases := []struct {
		name               string
		injectCarrierType  ot.BuiltinFormat
		extractCarrierType ot.BuiltinFormat
		extractCarrier     interface{}
		injectCarrier      interface{}
		extractErr         error
		injectErr          error
	}{
		{
			name:               "support for propagation.TextMapCarrier",
			injectCarrierType:  ot.TextMap,
			injectCarrier:      tmc,
			extractCarrierType: ot.TextMap,
			extractCarrier:     tmc,
		},
		{
			name:               "support for opentracing.TextMapReader and opentracing.TextMapWriter",
			injectCarrierType:  ot.TextMap,
			injectCarrier:      otTextMap,
			extractCarrierType: ot.TextMap,
			extractCarrier:     otTextMap,
		},
		{
			name:               "support for HTTPHeaders",
			injectCarrierType:  ot.HTTPHeaders,
			injectCarrier:      httpHeader,
			extractCarrierType: ot.HTTPHeaders,
			extractCarrier:     httpHeader,
		},
		{
			name:               "support for opentracing.TextMapReader and opentracing.TextMapWriter,non-same instance",
			injectCarrierType:  ot.TextMap,
			injectCarrier:      newTestTextMapWriter(&shareMap),
			extractCarrierType: ot.TextMap,
			extractCarrier:     newTestTextMapReader(&shareMap),
		},
		{
			name:              "inject: format type is HTTPHeaders, but carrier is not HTTPHeadersCarrier",
			injectCarrierType: ot.HTTPHeaders,
			injectCarrier:     struct{}{},
			injectErr:         ot.ErrInvalidCarrier,
		},
		{
			name:               "extract: format type is HTTPHeaders, but carrier is not HTTPHeadersCarrier",
			injectCarrierType:  ot.HTTPHeaders,
			injectCarrier:      httpHeader,
			extractCarrierType: ot.HTTPHeaders,
			extractCarrier:     struct{}{},
			extractErr:         ot.ErrInvalidCarrier,
		},
		{
			name:              "inject: format type is TextMap, but carrier is cannot be wrapped into propagation.TextMapCarrier",
			injectCarrierType: ot.TextMap,
			injectCarrier:     struct{}{},
			injectErr:         ot.ErrInvalidCarrier,
		},
		{
			name:               "extract: format type is TextMap, but carrier is cannot be wrapped into propagation.TextMapCarrier",
			injectCarrierType:  ot.TextMap,
			injectCarrier:      otTextMap,
			extractCarrierType: ot.TextMap,
			extractCarrier:     struct{}{},
			extractErr:         ot.ErrInvalidCarrier,
		},
		{
			name:              "inject: unsupported format type",
			injectCarrierType: ot.Binary,
			injectErr:         ot.ErrUnsupportedFormat,
		},
		{
			name:               "extract: unsupported format type",
			injectCarrierType:  ot.TextMap,
			injectCarrier:      otTextMap,
			extractCarrierType: ot.Binary,
			extractCarrier:     struct{}{},
			extractErr:         ot.ErrUnsupportedFormat,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := bridge.Inject(newBridgeSpanContext(trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: [16]byte{byte(1)},
				SpanID:  [8]byte{byte(2)},
			}), nil), tc.injectCarrierType, tc.injectCarrier)
			assert.Equal(t, tc.injectErr, err)

			if tc.injectErr == nil {
				spanContext, err := bridge.Extract(tc.extractCarrierType, tc.extractCarrier)
				assert.Equal(t, tc.extractErr, err)

				if tc.extractErr == nil {
					bsc, ok := spanContext.(*bridgeSpanContext)
					assert.True(t, ok)
					require.NotNil(t, bsc)
					require.NotNil(t, bsc.SpanContext)
					require.NotNil(t, bsc.SpanID())
					require.NotNil(t, bsc.TraceID())

					assert.Equal(t, spanID.String(), bsc.SpanID().String())
					assert.Equal(t, traceID.String(), bsc.TraceID().String())
				}
			}
		})
	}
}

type nonDeferWrapperTracer struct {
	*WrapperTracer
}

func (t *nonDeferWrapperTracer) Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	// Run start on the parent wrapper with a brand new context
	// so `WithDeferredSetup` hasn't been called, and the OpenTracing context is injected.
	return t.WrapperTracer.Start(context.Background(), name, opts...)
}

func TestBridgeTracer_StartSpan(t *testing.T) {
	testCases := []struct {
		name           string
		before         func(*testing.T, *BridgeTracer)
		expectWarnings []string
	}{
		{
			name: "with no option set",
			expectWarnings: []string{
				"The OpenTelemetry tracer is not set, default no-op tracer is used! Call SetOpenTelemetryTracer to set it up.\n",
			},
		},
		{
			name: "with wrapper tracer set",
			before: func(t *testing.T, bridge *BridgeTracer) {
				wTracer := NewWrapperTracer(bridge, otel.Tracer("test"))
				bridge.SetOpenTelemetryTracer(wTracer)
			},
			expectWarnings: []string(nil),
		},
		{
			name: "with a non-deferred wrapper tracer",
			before: func(t *testing.T, bridge *BridgeTracer) {
				wTracer := &nonDeferWrapperTracer{
					NewWrapperTracer(bridge, otel.Tracer("test")),
				}
				bridge.SetOpenTelemetryTracer(wTracer)
			},
			expectWarnings: []string{
				"SDK should have deferred the context setup, see the documentation of go.opentelemetry.io/otel/bridge/opentracing/migration\n",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var warningMessages []string
			bridge := NewBridgeTracer()
			bridge.SetWarningHandler(func(msg string) {
				warningMessages = append(warningMessages, msg)
			})

			if tc.before != nil {
				tc.before(t, bridge)
			}

			span := bridge.StartSpan("test")
			assert.NotNil(t, span)

			assert.Equal(t, tc.expectWarnings, warningMessages)
		})
	}
}

func Test_otTagToOTelAttr(t *testing.T) {
	key := attribute.Key("test")
	testCases := []struct {
		value    interface{}
		expected attribute.KeyValue
	}{
		{
			value:    int8(12),
			expected: key.Int64(int64(12)),
		},
		{
			value:    uint8(12),
			expected: key.Int64(int64(12)),
		},
		{
			value:    int16(12),
			expected: key.Int64(int64(12)),
		},
		{
			value:    uint16(12),
			expected: key.Int64(int64(12)),
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s %v", reflect.TypeOf(tc.value), tc.value), func(t *testing.T) {
			att := otTagToOTelAttr(string(key), tc.value)
			assert.Equal(t, tc.expected, att)
		})
	}
}

func Test_otTagsToOTelAttributesKindAndError(t *testing.T) {
	tracer := internal.NewMockTracer()
	sc := &bridgeSpanContext{}

	testCases := []struct {
		name     string
		opt      []ot.StartSpanOption
		expected trace.SpanKind
	}{
		{
			name:     "client",
			opt:      []ot.StartSpanOption{ext.SpanKindRPCClient},
			expected: trace.SpanKindClient,
		},
		{
			name:     "server",
			opt:      []ot.StartSpanOption{ext.RPCServerOption(sc)},
			expected: trace.SpanKindServer,
		},
		{
			name:     "client string",
			opt:      []ot.StartSpanOption{ot.Tag{Key: "span.kind", Value: "client"}},
			expected: trace.SpanKindClient,
		},
		{
			name:     "server string",
			opt:      []ot.StartSpanOption{ot.Tag{Key: "span.kind", Value: "server"}},
			expected: trace.SpanKindServer,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, _ := NewTracerPair(tracer)

			s := b.StartSpan(tc.name, tc.opt...)
			assert.Equal(t, tc.expected, s.(*bridgeSpan).otelSpan.(*internal.MockSpan).SpanKind)
		})
	}
}

func TestBridge_SpanContext_IsSampled(t *testing.T) {
	testCases := []struct {
		name     string
		flags    trace.TraceFlags
		expected bool
	}{
		{
			name:     "not sampled",
			flags:    0,
			expected: false,
		},
		{
			name:     "sampled",
			flags:    trace.FlagsSampled,
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tracer := internal.NewMockTracer()
			tracer.TraceFlags = tc.flags

			b, _ := NewTracerPair(tracer)
			s := b.StartSpan("abc")
			sc := s.Context()

			assert.Equal(t, tc.expected, sc.(samplable).IsSampled())
		})
	}
}

func TestBridgeSpanContextPromotedMethods(t *testing.T) {
	bridge := NewBridgeTracer()
	bridge.SetTextMapPropagator(new(testTextMapPropagator))

	tmc := newTextCarrier()

	type spanContextProvider interface {
		HasTraceID() bool
		TraceID() trace.TraceID
		HasSpanID() bool
		SpanID() trace.SpanID
	}

	err := bridge.Inject(newBridgeSpanContext(trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: [16]byte{byte(1)},
		SpanID:  [8]byte{byte(2)},
	}), nil), ot.TextMap, tmc)
	assert.NoError(t, err)

	spanContext, err := bridge.Extract(ot.TextMap, tmc)
	assert.NoError(t, err)

	assert.NotPanics(t, func() {
		assert.Equal(t, spanID.String(), spanContext.(spanContextProvider).SpanID().String())
		assert.Equal(t, traceID.String(), spanContext.(spanContextProvider).TraceID().String())
		assert.True(t, spanContext.(spanContextProvider).HasSpanID())
		assert.True(t, spanContext.(spanContextProvider).HasTraceID())
	})
}

func TestBridgeCarrierBaggagePropagation(t *testing.T) {
	carriers := []struct {
		name    string
		factory func() interface{}
		format  ot.BuiltinFormat
	}{
		{
			name:    "TextMapCarrier",
			factory: func() interface{} { return ot.TextMapCarrier{} },
			format:  ot.TextMap,
		},
		{
			name:    "HTTPHeadersCarrier",
			factory: func() interface{} { return ot.HTTPHeadersCarrier{} },
			format:  ot.HTTPHeaders,
		},
	}

	testCases := []struct {
		name         string
		baggageItems []bipBaggage
	}{
		{
			name: "single baggage item",
			baggageItems: []bipBaggage{
				{
					key:   "foo",
					value: "bar",
				},
			},
		},
		{
			name: "multiple baggage items",
			baggageItems: []bipBaggage{
				{
					key:   "foo",
					value: "bar",
				},
				{
					key:   "foo2",
					value: "bar2",
				},
			},
		},
		{
			name: "with characters escaped by baggage propagator",
			baggageItems: []bipBaggage{
				{
					key:   "space",
					value: "Hello world!",
				},
				{
					key:   "utf8",
					value: "Świat",
				},
			},
		},
	}

	for _, c := range carriers {
		for _, tc := range testCases {
			t.Run(fmt.Sprintf("%s %s", c.name, tc.name), func(t *testing.T) {
				mockOtelTracer := internal.NewMockTracer()
				b, _ := NewTracerPair(mockOtelTracer)
				b.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
					propagation.TraceContext{},
					propagation.Baggage{}), // Required for baggage propagation.
				)

				// Set baggage items.
				span := b.StartSpan("test")
				for _, bi := range tc.baggageItems {
					span.SetBaggageItem(bi.key, bi.value)
				}
				defer span.Finish()

				carrier := c.factory()
				err := b.Inject(span.Context(), c.format, carrier)
				assert.NoError(t, err)

				spanContext, err := b.Extract(c.format, carrier)
				assert.NoError(t, err)

				// Check baggage items.
				bsc, ok := spanContext.(*bridgeSpanContext)
				assert.True(t, ok)

				var got []bipBaggage
				for _, m := range bsc.bag.Members() {
					got = append(got, bipBaggage{m.Key(), m.Value()})
				}

				assert.ElementsMatch(t, tc.baggageItems, got)
			})
		}
	}
}
