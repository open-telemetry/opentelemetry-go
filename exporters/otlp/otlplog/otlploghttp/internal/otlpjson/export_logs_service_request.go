// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package otlpjson implements OTLP JSON Protobuf encoding for log data.
//
// The encoding conforms to the OTLP specs
// (https://opentelemetry.io/docs/specs/otlp/#json-protobuf-encoding):
//   - trace ID and span ID byte arrays are encoded as case-insensitive hex-encoded strings
//   - enum values encoded as integers
//   - field names in lowerCamelCase
//   - 64-bit integers encoded as quoted decimal strings (ProtoJSON specs)
//
// Common OTLP/JSON types and scalar codecs are generated from
// internal/shared/otlp/otlpjson.
package otlpjson

import (
	"encoding/json"

	collogpb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	logspb "go.opentelemetry.io/proto/otlp/logs/v1"
	resourcepb "go.opentelemetry.io/proto/otlp/resource/v1"
)

// ExportLogsServiceRequest corresponds to collogpb.ExportLogsServiceRequest.
type ExportLogsServiceRequest struct {
	ResourceLogs []*ResourceLogs `json:"resourceLogs,omitempty"`
}

// ResourceLogs corresponds to logspb.ResourceLogs.
type ResourceLogs struct {
	Resource  *Resource    `json:"resource,omitempty"`
	ScopeLogs []*ScopeLogs `json:"scopeLogs,omitempty"`
	SchemaURL string       `json:"schemaUrl,omitempty"`
}

// ScopeLogs corresponds to logspb.ScopeLogs.
type ScopeLogs struct {
	Scope      *InstrumentationScope `json:"scope,omitempty"`
	LogRecords []*LogRecord          `json:"logRecords,omitempty"`
	SchemaURL  string                `json:"schemaUrl,omitempty"`
}

// LogRecord corresponds to logspb.LogRecord.
type LogRecord struct {
	TimeUnixNano           Uint64      `json:"timeUnixNano,omitempty"`
	ObservedTimeUnixNano   Uint64      `json:"observedTimeUnixNano,omitempty"`
	SeverityNumber         int32       `json:"severityNumber,omitempty"`
	SeverityText           string      `json:"severityText,omitempty"`
	Body                   *AnyValue   `json:"body,omitempty"`
	Attributes             []*KeyValue `json:"attributes,omitempty"`
	DroppedAttributesCount uint32      `json:"droppedAttributesCount,omitempty"`
	Flags                  uint32      `json:"flags,omitempty"`
	TraceID                *TraceID    `json:"traceId,omitempty"`
	SpanID                 *SpanID     `json:"spanId,omitempty"`
	EventName              string      `json:"eventName,omitempty"`
}

// MarshalExportLogsServiceRequest encodes an ExportLogsServiceRequest as JSON Protobuf encoded bytes.
func MarshalExportLogsServiceRequest(req *collogpb.ExportLogsServiceRequest) ([]byte, error) {
	if req == nil {
		return []byte("{}"), nil
	}
	r := &ExportLogsServiceRequest{}
	for _, rl := range req.ResourceLogs {
		r.ResourceLogs = append(r.ResourceLogs, encodeResourceLogs(rl))
	}
	return json.Marshal(r)
}

func encodeResourceLogs(rl *logspb.ResourceLogs) *ResourceLogs {
	if rl == nil {
		return nil
	}
	out := &ResourceLogs{SchemaURL: rl.SchemaUrl}
	if rl.Resource != nil {
		out.Resource = &Resource{
			Attributes:             encodeKeyValues(rl.Resource.Attributes),
			DroppedAttributesCount: rl.Resource.DroppedAttributesCount,
			EntityRefs:             encodeEntityRefs(rl.Resource.EntityRefs),
		}
	}
	for _, sl := range rl.ScopeLogs {
		out.ScopeLogs = append(out.ScopeLogs, encodeScopeLogs(sl))
	}
	return out
}

func encodeScopeLogs(sl *logspb.ScopeLogs) *ScopeLogs {
	if sl == nil {
		return nil
	}
	out := &ScopeLogs{SchemaURL: sl.SchemaUrl}
	if sl.Scope != nil {
		out.Scope = &InstrumentationScope{
			Name:                   sl.Scope.Name,
			Version:                sl.Scope.Version,
			Attributes:             encodeKeyValues(sl.Scope.Attributes),
			DroppedAttributesCount: sl.Scope.DroppedAttributesCount,
		}
	}
	for _, lr := range sl.LogRecords {
		out.LogRecords = append(out.LogRecords, encodeLogRecord(lr))
	}
	return out
}

func encodeLogRecord(lr *logspb.LogRecord) *LogRecord {
	if lr == nil {
		return nil
	}
	out := &LogRecord{
		TimeUnixNano:           Uint64(lr.TimeUnixNano),
		ObservedTimeUnixNano:   Uint64(lr.ObservedTimeUnixNano),
		SeverityNumber:         int32(lr.SeverityNumber),
		SeverityText:           lr.SeverityText,
		Body:                   encodeAnyValue(lr.Body),
		Attributes:             encodeKeyValues(lr.Attributes),
		DroppedAttributesCount: lr.DroppedAttributesCount,
		Flags:                  lr.Flags,
		EventName:              lr.EventName,
	}
	if tid := encodeTraceID(lr.TraceId); tid != nil {
		out.TraceID = tid
	}
	if sid := encodeSpanID(lr.SpanId); sid != nil {
		out.SpanID = sid
	}
	return out
}

func encodeTraceID(b []byte) *TraceID {
	if len(b) == 0 {
		return nil
	}
	var tid TraceID
	copy(tid[:], b)
	if tid == (TraceID{}) {
		return nil
	}
	return &tid
}

func encodeSpanID(b []byte) *SpanID {
	if len(b) == 0 {
		return nil
	}
	var sid SpanID
	copy(sid[:], b)
	if sid == (SpanID{}) {
		return nil
	}
	return &sid
}

// UnmarshalExportLogsServiceRequest decodes JSON Protobuf encoded payload into an ExportLogsServiceRequest.
func UnmarshalExportLogsServiceRequest(data []byte, req *collogpb.ExportLogsServiceRequest) error {
	var jr ExportLogsServiceRequest
	if err := json.Unmarshal(data, &jr); err != nil {
		return err
	}

	for _, rl := range jr.ResourceLogs {
		req.ResourceLogs = append(req.ResourceLogs, decodeResourceLogs(rl))
	}
	return nil
}

func decodeResourceLogs(jrl *ResourceLogs) *logspb.ResourceLogs {
	rl := &logspb.ResourceLogs{SchemaUrl: jrl.SchemaURL}
	if jrl.Resource != nil {
		rl.Resource = &resourcepb.Resource{
			Attributes:             decodeKeyValues(jrl.Resource.Attributes),
			DroppedAttributesCount: jrl.Resource.DroppedAttributesCount,
			EntityRefs:             decodeEntityRefs(jrl.Resource.EntityRefs),
		}
	}
	for _, sl := range jrl.ScopeLogs {
		rl.ScopeLogs = append(rl.ScopeLogs, decodeScopeLogs(sl))
	}
	return rl
}

func decodeScopeLogs(jsl *ScopeLogs) *logspb.ScopeLogs {
	sl := &logspb.ScopeLogs{SchemaUrl: jsl.SchemaURL}
	if jsl.Scope != nil {
		sl.Scope = &commonpb.InstrumentationScope{
			Name:                   jsl.Scope.Name,
			Version:                jsl.Scope.Version,
			Attributes:             decodeKeyValues(jsl.Scope.Attributes),
			DroppedAttributesCount: jsl.Scope.DroppedAttributesCount,
		}
	}
	for _, lr := range jsl.LogRecords {
		sl.LogRecords = append(sl.LogRecords, decodeLogRecord(lr))
	}
	return sl
}

func decodeLogRecord(jlr *LogRecord) *logspb.LogRecord {
	lr := &logspb.LogRecord{
		TimeUnixNano:           uint64(jlr.TimeUnixNano),
		ObservedTimeUnixNano:   uint64(jlr.ObservedTimeUnixNano),
		SeverityNumber:         logspb.SeverityNumber(jlr.SeverityNumber),
		SeverityText:           jlr.SeverityText,
		Body:                   decodeAnyValue(jlr.Body),
		Attributes:             decodeKeyValues(jlr.Attributes),
		DroppedAttributesCount: jlr.DroppedAttributesCount,
		Flags:                  jlr.Flags,
		EventName:              jlr.EventName,
	}
	if jlr.TraceID != nil {
		lr.TraceId = jlr.TraceID[:]
	}
	if jlr.SpanID != nil {
		lr.SpanId = jlr.SpanID[:]
	}
	return lr
}
