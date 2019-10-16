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
	"regexp"
	"strconv"
	"strings"

	"go.opentelemetry.io/api/trace"

	"go.opentelemetry.io/api/core"
	dctx "go.opentelemetry.io/api/distributedcontext"
	apipropagation "go.opentelemetry.io/api/propagation"
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

var hexStr32ByteRegex = regexp.MustCompile("^[a-f0-9]{32}$")
var hexStr16ByteRegex = regexp.MustCompile("^[a-f0-9]{16}$")

func (b3 HTTPB3Propagator) Inject(ctx context.Context, supplier apipropagation.Supplier) {
	sc := trace.CurrentSpan(ctx).SpanContext()
	if sc.IsValid() {
		if b3.SingleHeader {
			sampled := sc.TraceFlags & core.TraceFlagsSampled
			supplier.Set(B3SingleHeader,
				fmt.Sprintf("%.16x%.16x-%.16x-%.1d", sc.TraceID.High, sc.TraceID.Low, sc.SpanID, sampled))
		} else {
			supplier.Set(B3TraceIDHeader,
				fmt.Sprintf("%.16x%.16x", sc.TraceID.High, sc.TraceID.Low))
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
	tid, ok := b3.extractTraceID(supplier.Get(B3TraceIDHeader))
	if !ok {
		return core.EmptySpanContext()
	}
	sid, ok := b3.extractSpanID(supplier.Get(B3SpanIDHeader))
	if !ok {
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

	var ok bool
	sc.TraceID, ok = b3.extractTraceID(parts[0])
	if !ok {
		return core.EmptySpanContext()
	}

	sc.SpanID, ok = b3.extractSpanID(parts[1])
	if !ok {
		return core.EmptySpanContext()
	}

	if l > 2 {
		sc.TraceFlags, ok = b3.extractSampledState(parts[2])
		if !ok {
			return core.EmptySpanContext()
		}
	}
	if l == 4 {
		_, ok = b3.extractSpanID(parts[3])
		if !ok {
			return core.EmptySpanContext()
		}
	}

	if !sc.IsValid() {
		return core.EmptySpanContext()
	}

	return sc
}

// extractTraceID parses the value of the X-B3-TraceId b3Header.
func (b3 HTTPB3Propagator) extractTraceID(tid string) (traceID core.TraceID, ok bool) {
	if hexStr32ByteRegex.MatchString(tid) {
		traceID.High, _ = strconv.ParseUint(tid[0:(16)], 16, 64)
		traceID.Low, _ = strconv.ParseUint(tid[(16):32], 16, 64)
		ok = true
	} else if b3.SingleHeader && hexStr16ByteRegex.MatchString(tid) {
		traceID.Low, _ = strconv.ParseUint(tid[:16], 16, 64)
		ok = true
	}
	return traceID, ok
}

// extractSpanID parses the value of the X-B3-SpanId or X-B3-ParentSpanId headers.
func (b3 HTTPB3Propagator) extractSpanID(sid string) (spanID uint64, ok bool) {
	if hexStr16ByteRegex.MatchString(sid) {
		spanID, _ = strconv.ParseUint(sid, 16, 64)
		ok = true
	}
	return spanID, ok
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
