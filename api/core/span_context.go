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
	"bytes"
	"encoding/hex"
	"fmt"
)

type TraceID [16]byte

var nilTraceID TraceID

const (
	traceFlagsBitMaskSampled = byte(0x01)
	traceFlagsBitMaskUnused  = byte(0xFE)

	// TraceFlagsSampled is a byte with sampled bit set. It is a convenient value initializer
	// for SpanContext TraceFlags field when a trace is sampled.
	TraceFlagsSampled = traceFlagsBitMaskSampled
	TraceFlagsUnused  = traceFlagsBitMaskUnused
)

type SpanContext struct {
	TraceID    TraceID
	SpanID     uint64
	TraceFlags byte
}

// EmptySpanContext is meant for internal use to return invalid span context during error conditions.
func EmptySpanContext() SpanContext {
	return SpanContext{}
}

func (sc SpanContext) IsValid() bool {
	return sc.HasTraceID() && sc.HasSpanID()
}

func (sc SpanContext) HasTraceID() bool {
	return !bytes.Equal(sc.TraceID[:], nilTraceID[:])
}

func (sc SpanContext) HasSpanID() bool {
	return sc.SpanID != 0
}

func (sc SpanContext) SpanIDString() string {
	return fmt.Sprintf("%.16x", sc.SpanID)
}

func (sc SpanContext) TraceIDString() string {
	return hex.EncodeToString(sc.TraceID[:])
}

func (sc SpanContext) IsSampled() bool {
	return sc.TraceFlags&traceFlagsBitMaskSampled == traceFlagsBitMaskSampled
}
