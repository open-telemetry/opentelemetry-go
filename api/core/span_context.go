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

package core

import (
	"fmt"
)

type TraceID struct {
	High uint64
	Low  uint64
}

const (
	traceOptionBitMaskSampled = byte(0x01)
	traceOptionBitMaskUnused  = byte(0xFE)

	// TraceOptionSampled is a byte with sampled bit set. It is a convenient value initialize
	// SpanContext when a trace is sampled.
	TraceOptionSampled = traceOptionBitMaskSampled
)

type SpanContext struct {
	TraceID      TraceID
	SpanID       uint64
	TraceOptions byte
}

var (
	// INVALID_SPAN_CONTEXT is meant for internal use to return invalid span context during error
	// conditions.
	INVALID_SPAN_CONTEXT = SpanContext{}
)

func (sc SpanContext) HasTraceID() bool {
	return sc.TraceID.High != 0 || sc.TraceID.Low != 0
}

func (sc SpanContext) HasSpanID() bool {
	return sc.SpanID != 0
}

func (sc SpanContext) SpanIDString() string {
	p := fmt.Sprintf("%.16x", sc.SpanID)
	return p[0:3] + ".." + p[13:16]
}

func (sc SpanContext) TraceIDString() string {
	p1 := fmt.Sprintf("%.16x", sc.TraceID.High)
	p2 := fmt.Sprintf("%.16x", sc.TraceID.Low)
	return p1[0:3] + ".." + p2[13:16]
}

func (sc SpanContext) IsSampled() bool {
	return sc.TraceOptions&traceOptionBitMaskSampled == traceOptionBitMaskSampled
}
