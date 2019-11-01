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
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/api/trace"

	"go.opentelemetry.io/otel/api/core"
	dctx "go.opentelemetry.io/otel/api/distributedcontext"
	apipropagation "go.opentelemetry.io/otel/api/propagation"
)

const (
	B3SingleHeader       = "X-B3"
	B3DebugFlagHeader    = "X-B3-Flags"
	B3TraceIDHeader      = "X-B3-TraceId"
	B3SpanIDHeader       = "X-B3-SpanId"
	B3SampledHeader      = "X-B3-Sampled"
	B3ParentSpanIDHeader = "X-B3-ParentSpanId"
)

// HTTPB3Propagator that facilitates core.SpanContext
// propagation using B3 Headers.
// This propagator supports both version of B3 headers,
//  1. Single Header :
//    X-B3: {TraceId}-{SpanId}-{SamplingState}-{ParentSpanId}
//  2. Multiple Headers:
//    X-B3-TraceId: {TraceId}
//    X-B3-ParentSpanId: {ParentSpanId}
//    X-B3-SpanId: {SpanId}
//    X-B3-Sampled: {SamplingState}
//    X-B3-Flags: {DebugFlag}
//
// If SingleHeader is set to true then X-B3 header is used to inject and extract. Otherwise,
// separate headers are used to inject and extract.
type HTTPB3Propagator struct {
	SingleHeader bool
}

var _ apipropagation.TextFormatPropagator = HTTPB3Propagator{}

func (b3 HTTPB3Propagator) Inject(ctx context.Context, supplier apipropagation.Supplier) {
	sc := trace.CurrentSpan(ctx).SpanContext()
	if sc.IsValid() {
		if b3.SingleHeader {
			sampled := sc.TraceFlags & core.TraceFlagsSampled
			supplier.Set(B3SingleHeader,
				fmt.Sprintf("%s-%.16x-%.1d", sc.TraceIDString(), sc.SpanID, sampled))
		} else {
			supplier.Set(B3TraceIDHeader, sc.TraceIDString())
			supplier.Set(B3SpanIDHeader,
				fmt.Sprintf("%.16x", sc.SpanID))

			var sampled string
			if sc.IsSampled() {
				sampled = "1"
			} else {
				sampled = "0"
			}
			supplier.Set(B3SampledHeader, sampled)
		}
	}
}

// Extract retrieves B3 Headers from the supplier
func (b3 HTTPB3Propagator) Extract(ctx context.Context, supplier apipropagation.Supplier) (core.SpanContext, dctx.Map) {
	if b3.SingleHeader {
		return b3.extractSingleHeader(supplier), dctx.NewEmptyMap()
	}
	return b3.extract(supplier), dctx.NewEmptyMap()
}

func (b3 HTTPB3Propagator) GetAllKeys() []string {
	if b3.SingleHeader {
		return []string{B3SingleHeader}
	}
	return []string{B3TraceIDHeader, B3SpanIDHeader, B3SampledHeader}
}

func (b3 HTTPB3Propagator) extract(supplier apipropagation.Supplier) core.SpanContext {
	tid, err := core.TraceIDFromHex(supplier.Get(B3TraceIDHeader))
	if err != nil {
		return core.EmptySpanContext()
	}
	sid, err := core.SpanIDFromHex(supplier.Get(B3SpanIDHeader))
	if err != nil {
		return core.EmptySpanContext()
	}
	sampled, ok := b3.extractSampledState(supplier.Get(B3SampledHeader))
	if !ok {
		return core.EmptySpanContext()
	}

	debug, ok := b3.extracDebugFlag(supplier.Get(B3DebugFlagHeader))
	if !ok {
		return core.EmptySpanContext()
	}
	if debug == core.TraceFlagsSampled {
		sampled = core.TraceFlagsSampled
	}

	sc := core.SpanContext{
		TraceID:    tid,
		SpanID:     sid,
		TraceFlags: sampled,
	}

	if !sc.IsValid() {
		return core.EmptySpanContext()
	}

	return sc
}

func (b3 HTTPB3Propagator) extractSingleHeader(supplier apipropagation.Supplier) core.SpanContext {
	h := supplier.Get(B3SingleHeader)
	if h == "" || h == "0" {
		core.EmptySpanContext()
	}
	sc := core.SpanContext{}
	parts := strings.Split(h, "-")
	l := len(parts)
	if l > 4 {
		return core.EmptySpanContext()
	}

	if l < 2 {
		return core.EmptySpanContext()
	}

	var err error
	sc.TraceID, err = core.TraceIDFromHex(parts[0])
	if err != nil {
		return core.EmptySpanContext()
	}

	sc.SpanID, err = core.SpanIDFromHex(parts[1])
	if err != nil {
		return core.EmptySpanContext()
	}

	if l > 2 {
		var ok bool
		sc.TraceFlags, ok = b3.extractSampledState(parts[2])
		if !ok {
			return core.EmptySpanContext()
		}
	}
	if l == 4 {
		_, err = core.SpanIDFromHex(parts[3])
		if err != nil {
			return core.EmptySpanContext()
		}
	}

	if !sc.IsValid() {
		return core.EmptySpanContext()
	}

	return sc
}

// extractSampledState parses the value of the X-B3-Sampled b3Header.
func (b3 HTTPB3Propagator) extractSampledState(sampled string) (flag byte, ok bool) {
	switch sampled {
	case "", "0":
		return 0, true
	case "1":
		return core.TraceFlagsSampled, true
	case "true":
		if !b3.SingleHeader {
			return core.TraceFlagsSampled, true
		}
	case "d":
		if b3.SingleHeader {
			return core.TraceFlagsSampled, true
		}
	}
	return 0, false
}

// extracDebugFlag parses the value of the X-B3-Sampled b3Header.
func (b3 HTTPB3Propagator) extracDebugFlag(debug string) (flag byte, ok bool) {
	switch debug {
	case "", "0":
		return 0, true
	case "1":
		return core.TraceFlagsSampled, true
	}
	return 0, false
}
