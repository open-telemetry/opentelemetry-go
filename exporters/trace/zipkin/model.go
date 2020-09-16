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

package zipkin

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	zkmodel "github.com/openzipkin/zipkin-go/model"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	export "go.opentelemetry.io/otel/sdk/export/trace"
)

const (
	keyInstrumentationLibraryName    = "otel.instrumentation_library.name"
	keyInstrumentationLibraryVersion = "otel.instrumentation_library.version"
)

func toZipkinSpanModels(batch []*export.SpanData, serviceName string) []zkmodel.SpanModel {
	models := make([]zkmodel.SpanModel, 0, len(batch))
	for _, data := range batch {
		models = append(models, toZipkinSpanModel(data, serviceName))
	}
	return models
}

func toZipkinSpanModel(data *export.SpanData, serviceName string) zkmodel.SpanModel {
	return zkmodel.SpanModel{
		SpanContext: toZipkinSpanContext(data),
		Name:        data.Name,
		Kind:        toZipkinKind(data.SpanKind),
		Timestamp:   data.StartTime,
		Duration:    data.EndTime.Sub(data.StartTime),
		Shared:      false,
		LocalEndpoint: &zkmodel.Endpoint{
			ServiceName: serviceName,
		},
		RemoteEndpoint: nil, // *Endpoint
		Annotations:    toZipkinAnnotations(data.MessageEvents),
		Tags:           toZipkinTags(data),
	}
}

func toZipkinSpanContext(data *export.SpanData) zkmodel.SpanContext {
	return zkmodel.SpanContext{
		TraceID:  toZipkinTraceID(data.SpanContext.TraceID),
		ID:       toZipkinID(data.SpanContext.SpanID),
		ParentID: toZipkinParentID(data.ParentSpanID),
		Debug:    false,
		Sampled:  nil,
		Err:      nil,
	}
}

func toZipkinTraceID(traceID otel.ID) zkmodel.TraceID {
	return zkmodel.TraceID{
		High: binary.BigEndian.Uint64(traceID[:8]),
		Low:  binary.BigEndian.Uint64(traceID[8:]),
	}
}

func toZipkinID(spanID otel.SpanID) zkmodel.ID {
	return zkmodel.ID(binary.BigEndian.Uint64(spanID[:]))
}

func toZipkinParentID(spanID otel.SpanID) *zkmodel.ID {
	if spanID.IsValid() {
		id := toZipkinID(spanID)
		return &id
	}
	return nil
}

func toZipkinKind(kind otel.SpanKind) zkmodel.Kind {
	switch kind {
	case otel.SpanKindUnspecified:
		return zkmodel.Undetermined
	case otel.SpanKindInternal:
		// The spec says we should set the kind to nil, but
		// the model does not allow that.
		return zkmodel.Undetermined
	case otel.SpanKindServer:
		return zkmodel.Server
	case otel.SpanKindClient:
		return zkmodel.Client
	case otel.SpanKindProducer:
		return zkmodel.Producer
	case otel.SpanKindConsumer:
		return zkmodel.Consumer
	}
	return zkmodel.Undetermined
}

func toZipkinAnnotations(events []export.Event) []zkmodel.Annotation {
	if len(events) == 0 {
		return nil
	}
	annotations := make([]zkmodel.Annotation, 0, len(events))
	for _, event := range events {
		value := event.Name
		if len(event.Attributes) > 0 {
			jsonString := attributesToJSONMapString(event.Attributes)
			if jsonString != "" {
				value = fmt.Sprintf("%s: %s", event.Name, jsonString)
			}
		}
		annotations = append(annotations, zkmodel.Annotation{
			Timestamp: event.Time,
			Value:     value,
		})
	}
	return annotations
}

func attributesToJSONMapString(attributes []label.KeyValue) string {
	m := make(map[string]interface{}, len(attributes))
	for _, attribute := range attributes {
		m[(string)(attribute.Key)] = attribute.Value.AsInterface()
	}
	// if an error happens, the result will be an empty string
	jsonBytes, _ := json.Marshal(m)
	return (string)(jsonBytes)
}

// extraZipkinTags are those that may be added to every outgoing span
var extraZipkinTags = []string{
	"ot.status_code",
	"ot.status_description",
	keyInstrumentationLibraryName,
	keyInstrumentationLibraryVersion,
}

func toZipkinTags(data *export.SpanData) map[string]string {
	m := make(map[string]string, len(data.Attributes)+len(extraZipkinTags))
	for _, kv := range data.Attributes {
		m[(string)(kv.Key)] = kv.Value.Emit()
	}
	if v, ok := m["error"]; ok && v == "false" {
		delete(m, "error")
	}
	m["ot.status_code"] = data.StatusCode.String()
	m["ot.status_description"] = data.StatusMessage

	if il := data.InstrumentationLibrary; il.Name != "" {
		m[keyInstrumentationLibraryName] = il.Name
		if il.Version != "" {
			m[keyInstrumentationLibraryVersion] = il.Version
		}
	}
	return m
}
