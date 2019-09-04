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
	"strconv"
	"strings"

	"go.opentelemetry.io/api/trace"

	"go.opentelemetry.io/api/core"
	apipropagation "go.opentelemetry.io/api/propagation"
)

const (
	supportedVersion  = 0
	maxVersion        = 254
	traceparentHeader = "traceparent"
)

type httpTraceContextPropagator struct{}

var _ apipropagation.Propagator = httpTraceContextPropagator{}

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
