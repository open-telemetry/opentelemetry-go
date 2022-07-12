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
	"errors"
	"net/http"
	"strings"
	"testing"

	ot "github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type testOnlyTextMapReader struct {
}

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
	assert.True(t, errors.Is(err, ot.ErrInvalidCarrier))

	_, err = newTextMapWrapperForExtract(newTestTextMapReaderAndWriter())
	assert.NoError(t, err)

	_, err = newTextMapWrapperForInject(newTestOnlyTextMapWriter())
	assert.NoError(t, err)

	_, err = newTextMapWrapperForInject(newTestOnlyTextMapReader())
	assert.True(t, errors.Is(err, ot.ErrInvalidCarrier))

	_, err = newTextMapWrapperForInject(newTestTextMapReaderAndWriter())
	assert.NoError(t, err)
}

func TestTextMapWrapper_action(t *testing.T) {
	testExtractFunc := func(carrier propagation.TextMapCarrier) {
		str := carrier.Keys()
		assert.Len(t, str, 2)
		assert.Contains(t, str, "key1", "key2")

		assert.Equal(t, carrier.Get("key1"), "val1")
		assert.Equal(t, carrier.Get("key2"), "val2")
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

type testTextMapPropagator struct {
}

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

	var exist = false

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

					assert.Equal(t, spanID.String(), bsc.otelSpanContext.SpanID().String())
					assert.Equal(t, traceID.String(), bsc.otelSpanContext.TraceID().String())
				}
			}
		})
	}
}
