// Copyright The OpenTelemetry Authors
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
	"encoding/json"
	"testing"
	"time"

	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/trace"
	export "go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/resource"
)

func TestExporter_ExportSpan(t *testing.T) {
	// write to buffer for testing
	var b bytes.Buffer
	ex, err := NewExporter(Options{Writer: &b})
	if err != nil {
		t.Errorf("Error constructing stdout exporter %s", err)
	}

	// setup test span
	now := time.Now()
	traceID, _ := trace.IDFromHex("0102030405060708090a0b0c0d0e0f10")
	spanID, _ := trace.SpanIDFromHex("0102030405060708")
	keyValue := "value"
	doubleValue := 123.456
	resource := resource.New(kv.String("rk1", "rv11"))

	testSpan := &export.SpanData{
		SpanContext: trace.SpanContext{
			TraceID: traceID,
			SpanID:  spanID,
		},
		Name:      "/foo",
		StartTime: now,
		EndTime:   now,
		Attributes: []kv.KeyValue{
			kv.String("key", keyValue),
			kv.Float64("double", doubleValue),
		},
		MessageEvents: []export.Event{
			{Name: "foo", Attributes: []kv.KeyValue{kv.String("key", keyValue)}, Time: now},
			{Name: "bar", Attributes: []kv.KeyValue{kv.Float64("double", doubleValue)}, Time: now},
		},
		SpanKind:      trace.SpanKindInternal,
		StatusCode:    codes.Unknown,
		StatusMessage: "interesting",
		Resource:      resource,
	}
	ex.ExportSpan(context.Background(), testSpan)

	expectedSerializedNow, _ := json.Marshal(now)

	got := b.String()
	expectedOutput := `{"SpanContext":{` +
		`"TraceID":"0102030405060708090a0b0c0d0e0f10",` +
		`"SpanID":"0102030405060708","TraceFlags":0},` +
		`"ParentSpanID":"0000000000000000",` +
		`"SpanKind":1,` +
		`"Name":"/foo",` +
		`"StartTime":` + string(expectedSerializedNow) + "," +
		`"EndTime":` + string(expectedSerializedNow) + "," +
		`"Attributes":[` +
		`{` +
		`"Key":"key",` +
		`"Value":{"Type":"STRING","Value":"value"}` +
		`},` +
		`{` +
		`"Key":"double",` +
		`"Value":{"Type":"FLOAT64","Value":123.456}` +
		`}],` +
		`"MessageEvents":[` +
		`{` +
		`"Name":"foo",` +
		`"Attributes":[` +
		`{` +
		`"Key":"key",` +
		`"Value":{"Type":"STRING","Value":"value"}` +
		`}` +
		`],` +
		`"Time":` + string(expectedSerializedNow) +
		`},` +
		`{` +
		`"Name":"bar",` +
		`"Attributes":[` +
		`{` +
		`"Key":"double",` +
		`"Value":{"Type":"FLOAT64","Value":123.456}` +
		`}` +
		`],` +
		`"Time":` + string(expectedSerializedNow) +
		`}` +
		`],` +
		`"Links":null,` +
		`"StatusCode":2,` +
		`"StatusMessage":"interesting",` +
		`"HasRemoteParent":false,` +
		`"DroppedAttributeCount":0,` +
		`"DroppedMessageEventCount":0,` +
		`"DroppedLinkCount":0,` +
		`"ChildSpanCount":0,` +
		`"Resource":[` +
		`{` +
		`"Key":"rk1",` +
		`"Value":{"Type":"STRING","Value":"rv11"}` +
		`}]}` + "\n"

	if got != expectedOutput {
		t.Errorf("Want: %v but got: %v", expectedOutput, got)
	}
}
