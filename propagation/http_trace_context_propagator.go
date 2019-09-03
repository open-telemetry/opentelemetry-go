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
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"

	"go.opentelemetry.io/api/trace"

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
var _ apipropagation.Propagator = httpTraceContextPropagator{}

// CarrierExtractor implements TextFormatPropagator interface.
//
// It creates CarrierExtractor object and binds carrier to the object. The carrier
// is expected to be *http.Request. If the carrier is nil or its type is not *http.Request
// then a NoopExtractor is returned.
func (hp httpTraceContextPropagator) CarrierExtractor(carrier interface{}) apipropagation.Extractor {
	req, _ := carrier.(*http.Request)
	if req != nil {
		return traceContextExtractor{req: req}
	}
	return apipropagation.NoopExtractor{}
}

// CarrierInjector implements TextFormatPropagator interface.
//
// It creates CarrierInjector object and binds carrier to the object. The carrier
// is expected to be of type *http.Request. If the carrier is nil or its type is not *http.Request
// then a NoopInjector is returned.
func (hp httpTraceContextPropagator) CarrierInjector(carrier interface{}) apipropagation.Injector {
	req, _ := carrier.(*http.Request)
	if req != nil {
		return traceContextInjector{req: req}
	}
	return apipropagation.NoopInjector{}
}

func (hp httpTraceContextPropagator) Inject(ctx context.Context, supplier apipropagation.Supplier) {
	sc := trace.CurrentSpan(ctx).SpanContext()
	if sc.IsValid() {
		h := fmt.Sprintf("%.2x-%.16x%.16x-%.16x-%.2x",
			supportedVersion,
			sc.TraceID.High,
			sc.TraceID.Low,
			sc.SpanID,
			sc.TraceOptions)
		supplier.Set(traceparentHeader, h)
	}
}

func (hp httpTraceContextPropagator) Extract(ctx context.Context, supplier apipropagation.Supplier) context.Context {
	h := supplier.Get(traceparentHeader)
	if h == "" {
		return ctx
	}

	sections := strings.Split(h, "-")
	if len(sections) < 4 {
		return ctx
	}

	if len(sections[0]) != 2 {
		return ctx
	}
	ver, err := hex.DecodeString(sections[0])
	if err != nil {
		return ctx
	}
	version := int(ver[0])
	if version > maxVersion {
		return ctx
	}

	if version == 0 && len(sections) != 4 {
		return ctx
	}

	if len(sections[1]) != 32 {
		return ctx
	}

	result, err := strconv.ParseUint(sections[1][0:16], 16, 64)
	if err != nil {
		return ctx
	}
	var sc core.SpanContext

	sc.TraceID.High = result

	result, err = strconv.ParseUint(sections[1][16:32], 16, 64)
	if err != nil {
		return ctx
	}
	sc.TraceID.Low = result

	if len(sections[2]) != 16 {
		return ctx
	}
	result, err = strconv.ParseUint(sections[2][0:], 16, 64)
	if err != nil {
		return ctx
	}
	sc.SpanID = result

	opts, err := hex.DecodeString(sections[3])
	if err != nil || len(opts) < 1 {
		return ctx
	}
	sc.TraceOptions = opts[0]

	if !sc.IsValid() {
		return ctx
	}

	ctx, _ = trace.GlobalTracer().Start(ctx, "remote", trace.CopyOfRemote(sc))
	return ctx
}

func (hp httpTraceContextPropagator) GetAllKeys() []string {
	return nil
}

// HttpTraceContextPropagator creates a new text format propagator that propagates SpanContext
// in W3C TraceContext format.
//
// The propagator is then used to create CarrierInjector and CarrierExtractor associated with a
// specific request. Injectors and Extractors respectively provides method to
// inject and extract SpanContext into/from the http request. Inject method encodes
// SpanContext and tag.Map into W3C TraceContext Header and injects the header in the request.
// Extract method extracts the header and decodes SpanContext and tag.Map.
func HttpTraceContextPropagator() httpTraceContextPropagator {
	return httpTraceContextPropagator{}
}

type traceContextExtractor struct {
	req *http.Request
}

var _ apipropagation.Extractor = traceContextExtractor{}

// Extract implements Extractor interface.
//
// It extracts W3C TraceContext Header and then decodes SpanContext and tag.Map from the Header.
func (tce traceContextExtractor) Extract() (sc core.SpanContext, tm tag.Map) {
	if tce.req == nil {
		return noExtract()
	}
	h, ok := getRequestHeader(tce.req, traceparentHeader, false)
	if !ok {
		return noExtract()
	}
	sections := strings.Split(h, "-")
	if len(sections) < 4 {
		return noExtract()
	}

	if len(sections[0]) != 2 {
		return noExtract()
	}
	ver, err := hex.DecodeString(sections[0])
	if err != nil {
		return noExtract()
	}
	version := int(ver[0])
	if version > maxVersion {
		return noExtract()
	}

	if version == 0 && len(sections) != 4 {
		return noExtract()
	}

	if len(sections[1]) != 32 {
		return noExtract()
	}

	result, err := strconv.ParseUint(sections[1][0:16], 16, 64)
	if err != nil {
		return noExtract()
	}
	sc.TraceID.High = result

	result, err = strconv.ParseUint(sections[1][16:32], 16, 64)
	if err != nil {
		return noExtract()
	}
	sc.TraceID.Low = result

	if len(sections[2]) != 16 {
		return noExtract()
	}
	result, err = strconv.ParseUint(sections[2][0:], 16, 64)
	if err != nil {
		return noExtract()
	}
	sc.SpanID = result

	opts, err := hex.DecodeString(sections[3])
	if err != nil || len(opts) < 1 {
		return noExtract()
	}
	sc.TraceOptions = opts[0]

	if !sc.IsValid() {
		return noExtract()
	}

	// TODO: [rghetia] add tag.Map (distributed context) extraction
	return sc, tag.NewEmptyMap()
}

func noExtract() (core.SpanContext, tag.Map) {
	return core.EmptySpanContext(), tag.NewEmptyMap()
}

type traceContextInjector struct {
	req *http.Request
}

var _ apipropagation.Injector = traceContextInjector{}

// Inject implements Injector interface.
//
// It encodes SpanContext and tag.Map into W3C TraceContext Header and injects the header into
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
