// Copyright 2019, OpenTelemetry Authors
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

package propagation

import (
	"fmt"
	"net/http"
	"strconv"

	"go.opentelemetry.io/api/core"
	apipropagation "go.opentelemetry.io/api/propagation"
	"go.opentelemetry.io/api/tag"
)

const (
	B3TraceIDHeader = "X-B3-TraceId"
	B3SpanIDHeader  = "X-B3-SpanId"
	B3SampledHeader = "X-B3-Sampled"
)

type httpB3Propagator struct{}

var _ apipropagation.TextFormatPropagator = httpB3Propagator{}

// CarrierExtractor implements CarrierExtractor method of TextFormatPropagator interface.
//
// It creates CarrierExtractor object and binds carrier to the object. The carrier
// is expected to be *http.Request. If the carrier type is not *http.Request
// then an empty extractor. Extract method on empty extractor does nothing.
func (t httpB3Propagator) CarrierExtractor(carrier interface{}) apipropagation.Extractor {
	req, ok := carrier.(*http.Request)
	if ok {
		return b3Extractor{req: req}
	}
	return b3Extractor{}
}

// CarrierInjector implements CarrierInjector method of TextFormatPropagator interface.
//
// It creates CarrierInjector object and binds carrier to the object. The carrier
// is expected to be of type *http.Request. If the carrier type is not *http.Request
// then an empty injector is returned. Inject method on empty injector does nothing.
func (t httpB3Propagator) CarrierInjector(carrier interface{}) apipropagation.Injector {
	req, ok := carrier.(*http.Request)
	if ok {
		return b3Injector{req: req}
	}
	return b3Injector{}
}

// HttpB3Propagator creates a new propagator. The propagator is then used
// to create Injector and Extrator associated with a specific request. Injectors
// and Extractors respectively provides method to inject and extract SpanContext
// into/from the http request. These methods encode/decode SpanContext to/from B3
// Headers.
func HttpB3Propagator() httpB3Propagator {
	return httpB3Propagator{}
}

type b3Extractor struct {
	req *http.Request
}

var _ apipropagation.Extractor = b3Extractor{}

// Extract implements Extract method of trace.Extractor interface. It extracts
// B3 Trace Headers and decodes SpanContext from these headers. These headers are
// X-B3-TraceId
// X-B3-SpanId
// X-B3-Sampled
// It skips the X-B3-ParentId and X-B3-Flags headers as they are not supported
// by core.SpanContext
func (be b3Extractor) Extract() (sc core.SpanContext, tm tag.Map) {
	if be.req == nil {
		return core.EmptySpanContext(), tag.NewEmptyMap()
	}
	tid := extractTraceID(be.req.Header.Get(B3TraceIDHeader))
	sid := extractSpanID(be.req.Header.Get(B3SpanIDHeader))
	sampled := extractTraceOption(be.req.Header.Get(B3SampledHeader))
	sc = core.SpanContext{
		TraceID:      tid,
		SpanID:       sid,
		TraceOptions: sampled,
	}

	if !sc.IsValid() {
		return core.EmptySpanContext(), tag.NewEmptyMap()
	}

	return sc, tag.NewEmptyMap()
}

// extractTraceID parses the value of the X-B3-TraceId header.
func extractTraceID(tid string) (traceID core.TraceID) {
	if tid == "" {
		return traceID
	}
	l := len(tid)
	if l > 32 {
		return core.TraceID{}
	} else if l > 16 {
		traceID.High, _ = strconv.ParseUint(tid[0:(l-16)], 16, 64)
		traceID.Low, _ = strconv.ParseUint(tid[(l-16):l], 16, 64)
	} else {
		traceID.Low, _ = strconv.ParseUint(tid[:l], 16, 64)
	}
	return traceID
}

// extractSpanID parses the value of the X-B3-SpanId or X-B3-ParentSpanId headers.
func extractSpanID(sid string) (spanID uint64) {
	if sid == "" {
		return spanID
	}
	spanID, _ = strconv.ParseUint(sid, 16, 64)
	return spanID
}

// extractTraceOption parses the value of the X-B3-Sampled header.
func extractTraceOption(sampled string) byte {
	switch sampled {
	case "true", "1":
		return core.TraceOptionSampled
	default:
		return 0
	}
}

type b3Injector struct {
	req *http.Request
}

var _ apipropagation.Injector = b3Injector{}

// Inject implements Inject method of trace.Injector interface. It encodes
// SpanContext into W3C TraceContext Header and injects the header into
// the associated request.
func (b3i b3Injector) Inject(sc core.SpanContext, tm tag.Map) {
	if b3i.req == nil {
		return
	}
	if sc.IsValid() {
		b3i.req.Header.Set(B3TraceIDHeader,
			fmt.Sprintf("%.16x%.16x", sc.TraceID.High, sc.TraceID.Low))
		b3i.req.Header.Set(B3SpanIDHeader,
			fmt.Sprintf("%.16x", sc.SpanID))

		var sampled string
		if sc.IsSampled() {
			sampled = "1"
		} else {
			sampled = "0"
		}
		b3i.req.Header.Set(B3SampledHeader, sampled)
	}
}
