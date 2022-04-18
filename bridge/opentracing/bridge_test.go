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
	"reflect"
	"sort"
	"testing"

	ot "github.com/opentracing/opentracing-go"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/bridge/opentracing/internal"
	"go.opentelemetry.io/otel/propagation"
)

func TestTextMapAdapterKeys(t *testing.T) {
	carrier := ot.TextMapCarrier{
		"foo": "bar",
		"baz": "qux",
	}

	keys := textMapAdapter{r: carrier}.Keys()
	sort.Strings(keys)
	expected := []string{"baz", "foo"}
	if !reflect.DeepEqual(expected, keys) {
		t.Errorf("Keys do not match: %#v, %#v", expected, keys)
	}
	// Check what happens if we read from a write-capable adaptor.
	keys = textMapAdapter{w: carrier}.Keys()
	if keys != nil {
		t.Errorf("Keys should be nil: %#v", keys)
	}
}

func TestMapCarrier(t *testing.T) {
	carrier := propagation.MapCarrier{}
	testBridgeWithCarrier(t, carrier)
}

func TestHeaderCarrier(t *testing.T) {
	carrier := propagation.HeaderCarrier{}
	testBridgeWithCarrier(t, carrier)
}

func TestOTHTTPHeadersCarrier(t *testing.T) {
	carrier := ot.HTTPHeadersCarrier{}
	testBridgeWithCarrier(t, carrier)
}

func TestOTTextMapCarrier(t *testing.T) {
	carrier := ot.TextMapCarrier{}
	testBridgeWithCarrier(t, carrier)
}

func testBridgeWithCarrier(t *testing.T, carrier interface{}) {
	mockOtelTracer := internal.NewMockTracer()
	bridgeTracer, _ := NewTracerPair(mockOtelTracer)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	span := bridgeTracer.StartSpan("testSpan1")
	defer span.Finish()

	err := bridgeTracer.Inject(span.Context(), ot.HTTPHeaders, carrier)
	if err != nil {
		t.Errorf("Inject error: %s", err)
	}

	extractedSpanContext, err := bridgeTracer.Extract(ot.HTTPHeaders, carrier)
	if err != nil {
		t.Errorf("Extract error: %s", err)
	}

	// Make a copy of the SpanContext with remote set to true
	originalSpanContext := span.Context().(*bridgeSpanContext)
	copySpanContext := *originalSpanContext
	copySpanContext.otelSpanContext = copySpanContext.otelSpanContext.WithRemote(true)
	// Now the copy should be equal to the original we passed in.
	if !reflect.DeepEqual(&copySpanContext, extractedSpanContext) {
		t.Errorf("Extracted span context does not match: %#v, %#v", span.Context(), extractedSpanContext)
	}
}
