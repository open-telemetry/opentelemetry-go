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
package otlpjson

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

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

// Resource corresponds to resourcepb.Resource.
type Resource struct {
	Attributes             []*KeyValue  `json:"attributes,omitempty"`
	DroppedAttributesCount uint32       `json:"droppedAttributesCount,omitempty"`
	EntityRefs             []*EntityRef `json:"entityRefs,omitempty"`
}

// ScopeLogs corresponds to logspb.ScopeLogs.
type ScopeLogs struct {
	Scope      *InstrumentationScope `json:"scope,omitempty"`
	LogRecords []*LogRecord          `json:"logRecords,omitempty"`
	SchemaURL  string                `json:"schemaUrl,omitempty"`
}

// InstrumentationScope corresponds to commonpb.InstrumentationScope.
type InstrumentationScope struct {
	Name                   string      `json:"name,omitempty"`
	Version                string      `json:"version,omitempty"`
	Attributes             []*KeyValue `json:"attributes,omitempty"`
	DroppedAttributesCount uint32      `json:"droppedAttributesCount,omitempty"`
}

// EntityRef corresponds to commonpb.EntityRef.
type EntityRef struct {
	SchemaURL       string   `json:"schemaUrl,omitempty"`
	Type            string   `json:"type,omitempty"`
	IdKeys          []string `json:"idKeys,omitempty"`
	DescriptionKeys []string `json:"descriptionKeys,omitempty"`
}

// KeyValue corresponds to commonpb.KeyValue.
type KeyValue struct {
	Key   string    `json:"key"`
	Value *AnyValue `json:"value,omitempty"`
}

// AnyValue corresponds to commonpb.AnyValue.
type AnyValue struct {
	StringValue *string      `json:"stringValue,omitempty"`
	BoolValue   *bool        `json:"boolValue,omitempty"`
	IntValue    *Int64       `json:"intValue,omitempty"`
	DoubleValue *float64     `json:"doubleValue,omitempty"`
	ArrayValue  *ArrayValue  `json:"arrayValue,omitempty"`
	KvlistValue *KvlistValue `json:"kvlistValue,omitempty"`
	BytesValue  []byte       `json:"bytesValue,omitempty"`
}

// ArrayValue corresponds to commonpb.ArrayValue.
type ArrayValue struct {
	Values []*AnyValue `json:"values,omitempty"`
}

// KvlistValue corresponds to commonpb.KeyValueList.
type KvlistValue struct {
	Values []*KeyValue `json:"values,omitempty"`
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

func encodeEntityRefs(ers []*commonpb.EntityRef) []*EntityRef {
	if len(ers) == 0 {
		return nil
	}
	out := make([]*EntityRef, len(ers))
	for i, er := range ers {
		if er == nil {
			continue
		}
		out[i] = &EntityRef{
			SchemaURL:       er.SchemaUrl,
			Type:            er.Type,
			IdKeys:          er.IdKeys,
			DescriptionKeys: er.DescriptionKeys,
		}
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

func encodeKeyValues(kvs []*commonpb.KeyValue) []*KeyValue {
	if len(kvs) == 0 {
		return nil
	}
	out := make([]*KeyValue, len(kvs))
	for i, kv := range kvs {
		if kv == nil {
			continue
		}
		out[i] = &KeyValue{
			Key:   kv.Key,
			Value: encodeAnyValue(kv.Value),
		}
	}
	return out
}

func encodeAnyValue(av *commonpb.AnyValue) *AnyValue {
	if av == nil {
		return nil
	}
	out := &AnyValue{}
	switch v := av.Value.(type) {
	case *commonpb.AnyValue_StringValue:
		out.StringValue = &v.StringValue
	case *commonpb.AnyValue_BoolValue:
		out.BoolValue = &v.BoolValue
	case *commonpb.AnyValue_IntValue:
		iv := Int64(v.IntValue)
		out.IntValue = &iv
	case *commonpb.AnyValue_DoubleValue:
		out.DoubleValue = &v.DoubleValue
	case *commonpb.AnyValue_ArrayValue:
		if v.ArrayValue != nil {
			arr := &ArrayValue{}
			for _, val := range v.ArrayValue.Values {
				arr.Values = append(arr.Values, encodeAnyValue(val))
			}
			out.ArrayValue = arr
		}
	case *commonpb.AnyValue_KvlistValue:
		if v.KvlistValue != nil {
			out.KvlistValue = &KvlistValue{
				Values: encodeKeyValues(v.KvlistValue.Values),
			}
		}
	case *commonpb.AnyValue_BytesValue:
		out.BytesValue = v.BytesValue
	}
	return out
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

func decodeEntityRefs(jers []*EntityRef) []*commonpb.EntityRef {
	if len(jers) == 0 {
		return nil
	}
	ers := make([]*commonpb.EntityRef, len(jers))
	for i, jer := range jers {
		if jer == nil {
			continue
		}
		ers[i] = &commonpb.EntityRef{
			SchemaUrl:       jer.SchemaURL,
			Type:            jer.Type,
			IdKeys:          jer.IdKeys,
			DescriptionKeys: jer.DescriptionKeys,
		}
	}
	return ers
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

func decodeKeyValues(jkvs []*KeyValue) []*commonpb.KeyValue {
	if len(jkvs) == 0 {
		return nil
	}
	kvs := make([]*commonpb.KeyValue, len(jkvs))
	for i, jkv := range jkvs {
		kvs[i] = &commonpb.KeyValue{
			Key:   jkv.Key,
			Value: decodeAnyValue(jkv.Value),
		}
	}
	return kvs
}

func decodeAnyValue(jav *AnyValue) *commonpb.AnyValue {
	if jav == nil {
		return nil
	}
	av := &commonpb.AnyValue{}
	switch {
	case jav.StringValue != nil:
		av.Value = &commonpb.AnyValue_StringValue{StringValue: *jav.StringValue}
	case jav.BoolValue != nil:
		av.Value = &commonpb.AnyValue_BoolValue{BoolValue: *jav.BoolValue}
	case jav.IntValue != nil:
		av.Value = &commonpb.AnyValue_IntValue{IntValue: int64(*jav.IntValue)}
	case jav.DoubleValue != nil:
		av.Value = &commonpb.AnyValue_DoubleValue{DoubleValue: *jav.DoubleValue}
	case jav.ArrayValue != nil:
		arr := &commonpb.ArrayValue{}
		for _, v := range jav.ArrayValue.Values {
			arr.Values = append(arr.Values, decodeAnyValue(v))
		}
		av.Value = &commonpb.AnyValue_ArrayValue{ArrayValue: arr}
	case jav.KvlistValue != nil:
		av.Value = &commonpb.AnyValue_KvlistValue{
			KvlistValue: &commonpb.KeyValueList{
				Values: decodeKeyValues(jav.KvlistValue.Values),
			},
		}
	case jav.BytesValue != nil:
		av.Value = &commonpb.AnyValue_BytesValue{BytesValue: jav.BytesValue}
	}
	return av
}

// Int64 encodes int64 as a quoted decimal string per ProtoJSON specs.
type Int64 int64

func (i Int64) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strconv.FormatInt(int64(i), 10) + `"`), nil
}

func (i *Int64) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		*i = Int64(v)
		return nil
	}
	var v int64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*i = Int64(v)
	return nil
}

// Uint64 encodes uint64 as a quoted decimal string per ProtoJSON specs.
type Uint64 uint64

func (i Uint64) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strconv.FormatUint(uint64(i), 10) + `"`), nil
}

func (i *Uint64) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		*i = Uint64(v)
		return nil
	}
	var v uint64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*i = Uint64(v)
	return nil
}

const base16Alphabets = "0123456789ABCDEF"

// TraceID encodes a 16-byte trace ID as a case-insensitive hex-encoded string.
type TraceID [16]byte

func (t TraceID) MarshalJSON() ([]byte, error) {
	var b [34]byte
	b[0] = '"'
	for i, v := range t {
		b[1+i*2] = base16Alphabets[v>>4]
		b[2+i*2] = base16Alphabets[v&0x0f]
	}
	b[33] = '"'
	return b[:], nil
}

func (t *TraceID) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	b, err := hex.DecodeString(str)
	if err != nil {
		return err
	}
	if len(b) != len(t) {
		return fmt.Errorf("invalid trace ID length: got %d, want %d", len(b), len(t))
	}
	copy(t[:], b)
	return nil
}

// SpanID encodes an 8-byte span ID as a case-insensitive hex-encoded string.
type SpanID [8]byte

func (s SpanID) MarshalJSON() ([]byte, error) {
	var b [18]byte
	b[0] = '"'
	for i, v := range s {
		b[1+i*2] = base16Alphabets[v>>4]
		b[2+i*2] = base16Alphabets[v&0x0f]
	}
	b[17] = '"'
	return b[:], nil
}

func (s *SpanID) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	b, err := hex.DecodeString(str)
	if err != nil {
		return err
	}
	if len(b) != len(s) {
		return fmt.Errorf("invalid span ID length: got %d, want %d", len(b), len(s))
	}
	copy(s[:], b)
	return nil
}
