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
	"encoding/json"
	"errors"
	"fmt"
)

const (
	traceFlagsBitMaskSampled = byte(0x01)
	traceFlagsBitMaskUnused  = byte(0xFE)

	// TraceFlagsSampled is a byte with sampled bit set. It is a convenient value initializer
	// for SpanContext TraceFlags field when a trace is sampled.
	TraceFlagsSampled = traceFlagsBitMaskSampled
	TraceFlagsUnused  = traceFlagsBitMaskUnused
)

type TraceID [16]byte

var nilTraceID TraceID

func (t TraceID) isValid() bool {
	return !bytes.Equal(t[:], nilTraceID[:])
}

// TraceIDFromHex returns a TraceID from a hex string if it is compliant
// with the w3c trace-context specification.
// See more at https://www.w3.org/TR/trace-context/#trace-id
func TraceIDFromHex(h string) (TraceID, error) {
	t := TraceID{}
	if len(h) != 32 {
		return t, errors.New("hex encoded trace-id must have length equals to 32")
	}

	for _, r := range h {
		switch {
		case 'a' <= r && r <= 'f':
			continue
		case '0' <= r && r <= '9':
			continue
		default:
			return t, errors.New("trace-id can only contain [0-9a-f] characters, all lowercase")
		}
	}

	b, err := hex.DecodeString(h)
	if err != nil {
		return t, err
	}
	copy(t[:], b)

	if !t.isValid() {
		return t, errors.New("trace-id can't be all zero")
	}
	return t, nil
}

type SpanContext struct {
	TraceID    TraceID
	SpanID     uint64
	TraceFlags byte
}

var _ json.Marshaler = (*SpanContext)(nil)

// EmptySpanContext is meant for internal use to return invalid span context during error conditions.
func EmptySpanContext() SpanContext {
	return SpanContext{}
}

// MarshalJSON implements a custom marshal function to encode SpanContext
// in a human readable format with hex encoded TraceID and SpanID.
func (sc SpanContext) MarshalJSON() ([]byte, error) {
	type JSONSpanContext struct {
		TraceID    string
		SpanID     string
		TraceFlags byte
	}

	return json.Marshal(JSONSpanContext{
		TraceID:    sc.TraceIDString(),
		SpanID:     sc.SpanIDString(),
		TraceFlags: sc.TraceFlags,
	})
}

func (sc SpanContext) IsValid() bool {
	return sc.HasTraceID() && sc.HasSpanID()
}

func (sc SpanContext) HasTraceID() bool {
	return sc.TraceID.isValid()
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
