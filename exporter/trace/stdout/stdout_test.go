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

package stdout

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/sdk/export"
)

func traceIDFromString(s string) core.TraceID {
	b, _ := hex.DecodeString(s)
	t := core.TraceID{}
	copy(t[:], b)
	return t
}

func TestExporter_ExportSpan(t *testing.T) {
	ex, err := NewExporter(Options{})
	if err != nil {
		t.Errorf("Error constructing stdout exporter %s", err)
	}

	// override output writer for testing
	var b bytes.Buffer
	ex.outputWriter = &b

	// setup test span
	now := time.Now()
	traceID := traceIDFromString("0102030405060708090a0b0c0d0e0f10")
	spanID := uint64(0x0102030405060708)
	keyValue := "value"
	doubleValue := float64(123.456)

	testSpan := &export.SpanData{
		SpanContext: core.SpanContext{
			TraceID: traceID,
			SpanID:  spanID,
		},
		Name:      "/foo",
		StartTime: now,
		EndTime:   now,
		Attributes: []core.KeyValue{
			{
				Key:   core.Key("key"),
				Value: core.Value{Type: core.STRING, String: keyValue},
			},
			{
				Key:   core.Key("double"),
				Value: core.Value{Type: core.FLOAT64, Float64: doubleValue},
			},
		},
		Status: codes.Unknown,
	}
	ex.ExportSpan(context.Background(), testSpan)

	expectedSerializedNow, _ := json.Marshal(now)

	got := b.String()
	expectedOutput := `{"SpanContext":{` +
		//`"TraceID":{"High":72623859790382856,"Low":651345242494996240},` +
		// FIXME should this printed as an array of bytes?
		`"TraceID":[1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16],` +
		`"SpanID":72623859790382856,"TraceFlags":0},` +
		`"ParentSpanID":0,` +
		`"SpanKind":0,` +
		`"Name":"/foo",` +
		`"StartTime":` + string(expectedSerializedNow) + "," +
		`"EndTime":` + string(expectedSerializedNow) + "," +
		`"Attributes":[` +
		`{` +
		`"Key":"key",` +
		`"Value":{"Type":8,"Bool":false,"Int64":0,"Uint64":0,"Float64":0,"String":"value","Bytes":null}` +
		`},` +
		`{` +
		`"Key":"double",` +
		`"Value":{"Type":7,"Bool":false,"Int64":0,"Uint64":0,"Float64":123.456,"String":"","Bytes":null}` +
		`}` +
		`],` +
		`"MessageEvents":null,` +
		`"Links":null,` +
		`"Status":2,` +
		`"HasRemoteParent":false,` +
		`"DroppedAttributeCount":0,` +
		`"DroppedMessageEventCount":0,` +
		`"DroppedLinkCount":0,` +
		`"ChildSpanCount":0}` + "\n"

	if got != expectedOutput {
		t.Errorf("Want: %v but got: %v", expectedOutput, got)
	}
}
