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

package jaeger

import (
	"sort"
	"testing"
	"time"

	apitrace "go.opentelemetry.io/api/trace"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/api/core"
	gen "go.opentelemetry.io/exporter/trace/jaeger/internal/gen-go/jaeger"
	"go.opentelemetry.io/sdk/export"
)

// TODO(rghetia): Test export.

func Test_spanDataToThrift(t *testing.T) {
	now := time.Now()
	traceID := core.TraceID{High: 0x0102030405060708, Low: 0x090a0b0c0d0e0f10}
	spanID := uint64(0x0102030405060708)

	linkTraceID := core.TraceID{High: 0x0102030405060709, Low: 0x090a0b0c0d0e0f11}
	linkSpanID := uint64(0x0102030405060709)

	keyValue := "value"
	statusCodeValue := int64(2)
	doubleValue := float64(123.456)
	boolTrue := true
	statusMessage := "Unknown"

	tests := []struct {
		name string
		data *export.SpanData
		want *gen.Span
	}{
		{
			name: "no parent",
			data: &export.SpanData{
				SpanContext: core.SpanContext{
					TraceID: traceID,
					SpanID:  spanID,
				},
				Name:      "/foo",
				StartTime: now,
				EndTime:   now,
				Links: []apitrace.Link{
					{
						SpanContext: core.SpanContext{
							TraceID: linkTraceID,
							SpanID:  linkSpanID,
						},
					},
				},
				Attributes: []core.KeyValue{
					{
						Key:   core.Key{Name: "key"},
						Value: core.Value{Type: core.STRING, String: keyValue},
					},
					{
						Key:   core.Key{Name: "double"},
						Value: core.Value{Type: core.FLOAT64, Float64: doubleValue},
					},
				},
				// TODO: [rghetia] add events test after event is concrete type.
				Status: codes.Unknown,
			},
			want: &gen.Span{
				TraceIdLow:    651345242494996240,
				TraceIdHigh:   72623859790382856,
				SpanId:        72623859790382856,
				OperationName: "/foo",
				StartTime:     now.UnixNano() / 1000,
				Duration:      0,
				Tags: []*gen.Tag{
					{Key: "double", VType: gen.TagType_DOUBLE, VDouble: &doubleValue},
					{Key: "key", VType: gen.TagType_STRING, VStr: &keyValue},
					{Key: "error", VType: gen.TagType_BOOL, VBool: &boolTrue},
					{Key: "status.code", VType: gen.TagType_LONG, VLong: &statusCodeValue},
					{Key: "status.message", VType: gen.TagType_STRING, VStr: &statusMessage},
				},
				References: []*gen.SpanRef{
					{
						RefType:     gen.SpanRefType_CHILD_OF,
						TraceIdLow:  int64(linkTraceID.Low),
						TraceIdHigh: int64(linkTraceID.High),
						SpanId:      int64(linkSpanID),
					},
				},
				// TODO [rghetia]: check Logs when event is added.
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := spanDataToThrift(tt.data)
			sort.Slice(got.Tags, func(i, j int) bool {
				return got.Tags[i].Key < got.Tags[j].Key
			})
			sort.Slice(tt.want.Tags, func(i, j int) bool {
				return tt.want.Tags[i].Key < tt.want.Tags[j].Key
			})
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Diff%v", diff)
			}
		})
	}
}
