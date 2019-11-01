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
	"context"
	"encoding/json"
	"io"
	"os"
	"time"

	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/otel/api/core"
	apitrace "go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/sdk/export"
)

// Options are the options to be used when initializing a stdout export.
type Options struct {
	// PrettyPrint will pretty the json representation of the span,
	// making it print "pretty". Default is false.
	PrettyPrint bool
}

// Exporter is an implementation of trace.Exporter that writes spans to stdout.
type Exporter struct {
	pretty       bool
	outputWriter io.Writer
}

func NewExporter(o Options) (*Exporter, error) {
	return &Exporter{
		pretty:       o.PrettyPrint,
		outputWriter: os.Stdout,
	}, nil
}

type jsonValue struct {
	Type  string
	Value interface{}
}

type jsonKeyValue struct {
	Key   core.Key
	Value jsonValue
}

type jsonSpanData struct {
	SpanContext              core.SpanContext
	ParentSpanID             core.SpanID
	SpanKind                 apitrace.SpanKind
	Name                     string
	StartTime                time.Time
	EndTime                  time.Time
	Attributes               []jsonKeyValue
	MessageEvents            []export.Event
	Links                    []apitrace.Link
	Status                   codes.Code
	HasRemoteParent          bool
	DroppedAttributeCount    int
	DroppedMessageEventCount int
	DroppedLinkCount         int
	ChildSpanCount           int
}

func marshalSpanData(data *export.SpanData, pretty bool) ([]byte, error) {
	jsd := jsonSpanData{
		SpanContext:              data.SpanContext,
		ParentSpanID:             data.ParentSpanID,
		SpanKind:                 data.SpanKind,
		Name:                     data.Name,
		StartTime:                data.StartTime,
		EndTime:                  data.EndTime,
		Attributes:               toJSONAttributes(data.Attributes),
		MessageEvents:            data.MessageEvents,
		Links:                    data.Links,
		Status:                   data.Status,
		HasRemoteParent:          data.HasRemoteParent,
		DroppedAttributeCount:    data.DroppedAttributeCount,
		DroppedMessageEventCount: data.DroppedMessageEventCount,
		DroppedLinkCount:         data.DroppedLinkCount,
		ChildSpanCount:           data.ChildSpanCount,
	}

	if pretty {
		return json.MarshalIndent(jsd, "", "\t")
	}
	return json.Marshal(jsd)
}

func toJSONAttributes(attributes []core.KeyValue) []jsonKeyValue {
	jsonAttrs := make([]jsonKeyValue, len(attributes))
	for i := 0; i < len(attributes); i++ {
		jsonAttrs[i].Key = attributes[i].Key
		jsonAttrs[i].Value.Type = attributes[i].Value.Type().String()
		jsonAttrs[i].Value.Value = attributes[i].Value.AsInterface()
	}
	return jsonAttrs
}

// ExportSpan writes a SpanData in json format to stdout.
func (e *Exporter) ExportSpan(ctx context.Context, data *export.SpanData) {
	jsonSpan, err := marshalSpanData(data, e.pretty)
	if err != nil {
		// ignore writer failures for now
		_, _ = e.outputWriter.Write([]byte("Error converting spanData to json: " + err.Error()))
		return
	}
	// ignore writer failures for now
	_, _ = e.outputWriter.Write(append(jsonSpan, byte('\n')))
}
