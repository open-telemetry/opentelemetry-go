// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package telemetry

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

var (
	traceID = [16]byte{0x1}

	spanID1 = [8]byte{0x1}
	spanID2 = [8]byte{0x2}

	now      = time.Now()
	nowPlus1 = now.Add(1 * time.Second)

	spanA = &Span{
		TraceID:      traceID,
		SpanID:       spanID2,
		ParentSpanID: spanID1,
		Flags:        1,
		Name:         "span-a",
		StartTime:    now,
		EndTime:      nowPlus1,
		Status: &Status{
			Message: "test status",
			Code:    StatusCodeOK,
		},
	}

	spanB = &Span{}

	scopeSpans = &ScopeSpans{
		Scope: &Scope{
			Name:    "TestTracer",
			Version: "v0.1.0",
		},
		SchemaURL: "http://go.opentelemetry.io/test",
		Spans:     []*Span{spanA, spanB},
	}
)

func BenchmarkJSONMarshalUnmarshal(b *testing.B) {
	var out ScopeSpans

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		var inBuf bytes.Buffer
		enc := json.NewEncoder(&inBuf)
		err := enc.Encode(scopeSpans)
		if err != nil {
			b.Fatal(err)
		}

		payload := inBuf.Bytes()

		dec := json.NewDecoder(bytes.NewReader(payload))
		err = dec.Decode(&out)
		if err != nil {
			b.Fatal(err)
		}
	}
	_ = out
}
