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
	"encoding/hex"
	"fmt"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"

	"go.opentelemetry.io/api/core"
	apipropagation "go.opentelemetry.io/api/propagation"
	"go.opentelemetry.io/api/tag"
)

const (
	supportedVersion  = 0
	maxVersion        = 254
	traceparentHeader = "traceparent"
)

type httpTraceContextPropagator struct{}

var _ apipropagation.TextFormatPropagator = httpTraceContextPropagator{}

// CarrierExtractor implements CarrierExtractor method of TextFormatPropagator interface.
//
// It creates CarrierExtractor object and binds carrier to the object. The carrier
// is expected to be *http.Request. If the carrier type is not *http.Request
// then an empty extractor. Extract method on empty extractor does nothing.
func (hp httpTraceContextPropagator) CarrierExtractor(carrier interface{}) apipropagation.Extractor {
	req, ok := carrier.(*http.Request)
	if ok {
		return traceContextExtractor{req: req}
	}
	return traceContextExtractor{}
}

// CarrierInjector implements CarrierInjector method of TextFormatPropagator interface.
//
// It creates CarrierInjector object and binds carrier to the object. The carrier
// is expected to be of type *http.Request. If the carrier type is not *http.Request
// then an empty injector is returned. Inject method on empty injector does nothing.
func (hp httpTraceContextPropagator) CarrierInjector(carrier interface{}) apipropagation.Injector {
	req, ok := carrier.(*http.Request)
	if ok {
		return traceContextInjector{req: req}
	}
	return traceContextInjector{}
}

// HttpTraceContextPropagator creates a new propagator that propagates SpanContext
// in W3C TraceContext format.
//
// The propagator is then used to create CarrierInjector and CarrierExtractor associated with a
// specific request. Injectors and Extractors respectively provides method to
// inject and extract SpanContext into/from the http request. Inject method encodes
// SpanContext into W3C TraceContext Header and injects the header in the request.
// Extract method extracts the header and decodes SpanContext.
func HttpTraceContextPropagator() httpTraceContextPropagator {
	return httpTraceContextPropagator{}
}

type traceContextExtractor struct {
	req *http.Request
}

var _ apipropagation.Extractor = traceContextExtractor{}

// Extract implements Extract method of propagation.Extractor interface. It extracts
// W3C TraceContext Header and decodes SpanContext from the Header.
func (tce traceContextExtractor) Extract() (sc core.SpanContext, tm tag.Map) {
	if tce.req == nil {
		return core.EmptySpanContext(), tag.NewEmptyMap()
	}
	h, ok := getRequestHeader(tce.req, traceparentHeader, false)
	if !ok {
		return core.EmptySpanContext(), tag.NewEmptyMap()
	}
	sections := strings.Split(h, "-")
	if len(sections) < 4 {
		return core.EmptySpanContext(), tag.NewEmptyMap()
	}

	if len(sections[0]) != 2 {
		return core.EmptySpanContext(), tag.NewEmptyMap()
	}
	ver, err := hex.DecodeString(sections[0])
	if err != nil {
		return core.EmptySpanContext(), tag.NewEmptyMap()
	}
	version := int(ver[0])
	if version > maxVersion {
		return core.EmptySpanContext(), tag.NewEmptyMap()
	}

	if version == 0 && len(sections) != 4 {
		return core.EmptySpanContext(), tag.NewEmptyMap()
	}

	if len(sections[1]) != 32 {
		return core.EmptySpanContext(), tag.NewEmptyMap()
	}

	result, err := strconv.ParseUint(sections[1][0:16], 16, 64)
	if err != nil {
		return core.EmptySpanContext(), tag.NewEmptyMap()
	}
	sc.TraceID.High = result

	result, err = strconv.ParseUint(sections[1][16:32], 16, 64)
	if err != nil {
		return core.EmptySpanContext(), tag.NewEmptyMap()
	}
	sc.TraceID.Low = result

	if len(sections[2]) != 16 {
		return core.EmptySpanContext(), tag.NewEmptyMap()
	}
	result, err = strconv.ParseUint(sections[2][0:], 16, 64)
	if err != nil {
		return core.EmptySpanContext(), tag.NewEmptyMap()
	}
	sc.SpanID = result

	opts, err := hex.DecodeString(sections[3])
	if err != nil || len(opts) < 1 {
		return core.EmptySpanContext(), tag.NewEmptyMap()
	}
	sc.TraceOptions = opts[0]

	if !sc.IsValid() {
		return core.EmptySpanContext(), tag.NewEmptyMap()
	}

	// TODO: [rghetia] add tag.Map (distributed context) extraction
	return sc, tag.NewEmptyMap()
}

type traceContextInjector struct {
	req *http.Request
}

var _ apipropagation.Injector = traceContextInjector{}

// Inject implements Inject method of propagation.Injector interface. It encodes
// SpanContext into W3C TraceContext Header and injects the header into
// the associated request.
func (tci traceContextInjector) Inject(sc core.SpanContext, tm tag.Map) {
	if tci.req == nil {
		return
	}
	if sc.IsValid() {
		h := fmt.Sprintf("%.2x-%.16x%.16x-%.16x-%.2x",
			supportedVersion,
			sc.TraceID.High,
			sc.TraceID.Low,
			sc.SpanID,
			sc.TraceOptions)
		tci.req.Header.Set(traceparentHeader, h)
	}
	// TODO: [rghetia] add tag.Map (distributed context) injection
}

// getRequestHeader returns a combined header field according to RFC7230 section 3.2.2.
// If commaSeparated is true, multiple header fields with the same field name using be
// combined using ",".
// If no header was found using the given name, "ok" would be false.
// If more than one headers was found using the given name, while commaSeparated is false,
// "ok" would be false.
func getRequestHeader(req *http.Request, name string, commaSeparated bool) (hdr string, ok bool) {
	v := req.Header[textproto.CanonicalMIMEHeaderKey(name)]
	switch len(v) {
	case 0:
		return "", false
	case 1:
		return v[0], true
	default:
		return strings.Join(v, ","), commaSeparated
	}
}
