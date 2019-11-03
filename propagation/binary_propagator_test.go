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

package propagation_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/propagation"
)

func TestExtractSpanContextFromBytes(t *testing.T) {
	traceID, _ := core.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	spanID, _ := core.SpanIDFromHex("00f067aa0ba902b7")

	propagator := propagation.BinaryPropagator()
	tests := []struct {
		name   string
		bytes  []byte
		wantSc core.SpanContext
	}{
		{
			name: "future version of the proto",
			bytes: []byte{
				0x02, 0x00, 0x4b, 0xf9, 0x2f, 0x35, 0x77, 0xb3, 0x4d, 0xa6, 0xa3, 0xce, 0x92, 0x9d, 0x0e, 0x0e, 0x47, 0x36,
				0x01, 0x00, 0xf0, 0x67, 0xaa, 0x0b, 0xa9, 0x02, 0xb7,
				0x02, 0x01,
			},
			wantSc: core.SpanContext{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: core.TraceFlagsSampled,
			},
		},
		{
			name: "current version with valid SpanContext and with Sampled bit set",
			bytes: []byte{
				0x00, 0x00, 0x4b, 0xf9, 0x2f, 0x35, 0x77, 0xb3, 0x4d, 0xa6, 0xa3, 0xce, 0x92, 0x9d, 0x0e, 0x0e, 0x47, 0x36,
				0x01, 0x00, 0xf0, 0x67, 0xaa, 0x0b, 0xa9, 0x02, 0xb7,
				0x02, 0x01,
			},
			wantSc: core.SpanContext{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: core.TraceFlagsSampled,
			},
		},
		{
			name: "valid SpanContext without option",
			bytes: []byte{
				0x00, 0x00, 0x4b, 0xf9, 0x2f, 0x35, 0x77, 0xb3, 0x4d, 0xa6, 0xa3, 0xce, 0x92, 0x9d, 0x0e, 0x0e, 0x47, 0x36,
				0x01, 0x00, 0xf0, 0x67, 0xaa, 0x0b, 0xa9, 0x02, 0xb7,
			},
			wantSc: core.SpanContext{
				TraceID: traceID,
				SpanID:  spanID,
			},
		},
		{
			name: "zero trace ID",
			bytes: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x02, 0x01,
			},
			wantSc: core.EmptySpanContext(),
		},
		{
			name: "zero span ID",
			bytes: []byte{
				0x00, 0x00, 0x4b, 0xf9, 0x2f, 0x35, 0x77, 0xb3, 0x4d, 0xa6, 0xa3, 0xce, 0x92, 0x9d, 0x0e, 0x0e, 0x47, 0x36,
				0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x02, 0x01,
			},
			wantSc: core.EmptySpanContext(),
		},
		{
			name: "wrong trace ID field number",
			bytes: []byte{
				0x00, 0x01, 0x4b, 0xf9, 0x2f, 0x35, 0x77, 0xb3, 0x4d, 0xa6, 0xa3, 0xce, 0x92, 0x9d, 0x0e, 0x0e, 0x47, 0x36,
				0x01, 0x00, 0xf0, 0x67, 0xaa, 0x0b, 0xa9, 0x02, 0xb7,
			},
			wantSc: core.EmptySpanContext(),
		},
		{
			name: "short byte array",
			bytes: []byte{
				0x00, 0x00, 0x4b, 0xf9, 0x2f, 0x35, 0x77, 0xb3, 0x4d,
			},
			wantSc: core.EmptySpanContext(),
		},
		{
			name:   "nil byte array",
			wantSc: core.EmptySpanContext(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSc := propagator.FromBytes(tt.bytes)
			if diff := cmp.Diff(gotSc, tt.wantSc); diff != "" {
				t.Errorf("Deserialize SpanContext from byte array: %s: -got +want %s", tt.name, diff)
			}
		})
	}
}

func TestConvertSpanContextToBytes(t *testing.T) {
	traceID, _ := core.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	spanID, _ := core.SpanIDFromHex("00f067aa0ba902b7")

	propagator := propagation.BinaryPropagator()
	tests := []struct {
		name  string
		sc    core.SpanContext
		bytes []byte
	}{
		{
			name: "valid SpanContext, with sampling bit set",
			sc: core.SpanContext{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: core.TraceFlagsSampled,
			},
			bytes: []byte{
				0x00, 0x00, 0x4b, 0xf9, 0x2f, 0x35, 0x77, 0xb3, 0x4d, 0xa6, 0xa3, 0xce, 0x92, 0x9d, 0x0e, 0x0e, 0x47, 0x36,
				0x01, 0x00, 0xf0, 0x67, 0xaa, 0x0b, 0xa9, 0x02, 0xb7,
				0x02, 0x01,
			},
		},
		{
			name: "valid SpanContext, with sampling bit cleared",
			sc: core.SpanContext{
				TraceID: traceID,
				SpanID:  spanID,
			},
			bytes: []byte{
				0x00, 0x00, 0x4b, 0xf9, 0x2f, 0x35, 0x77, 0xb3, 0x4d, 0xa6, 0xa3, 0xce, 0x92, 0x9d, 0x0e, 0x0e, 0x47, 0x36,
				0x01, 0x00, 0xf0, 0x67, 0xaa, 0x0b, 0xa9, 0x02, 0xb7,
				0x02, 0x00,
			},
		},
		{
			name: "invalid spancontext",
			sc:   core.EmptySpanContext(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBytes := propagator.ToBytes(tt.sc)
			if diff := cmp.Diff(gotBytes, tt.bytes); diff != "" {
				t.Errorf("Serialize SpanContext to byte array: %s: -got +want %s", tt.name, diff)
			}
		})
	}
}
