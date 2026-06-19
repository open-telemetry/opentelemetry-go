// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutlog // import "go.opentelemetry.io/otel/exporters/stdout/stdoutlog"

import (
	"encoding/json"
	"errors"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
)

func newValue(v attribute.Value) value {
	return value{Value: v}
}

type value struct {
	attribute.Value
}

// MarshalJSON implements a custom marshal function to encode attribute.Value.
func (v value) MarshalJSON() ([]byte, error) {
	var jsonVal struct {
		Type  string
		Value any
	}
	jsonVal.Type = valueType(v.Value)

	switch v.Type() {
	case attribute.STRING:
		jsonVal.Value = v.AsString()
	case attribute.INT64:
		jsonVal.Value = v.AsInt64()
	case attribute.FLOAT64:
		jsonVal.Value = v.AsFloat64()
	case attribute.BOOL:
		jsonVal.Value = v.AsBool()
	case attribute.BYTESLICE:
		jsonVal.Value = v.AsByteSlice()
	case attribute.BOOLSLICE:
		jsonVal.Value = v.AsBoolSlice()
	case attribute.INT64SLICE:
		jsonVal.Value = v.AsInt64Slice()
	case attribute.FLOAT64SLICE:
		jsonVal.Value = v.AsFloat64Slice()
	case attribute.STRINGSLICE:
		jsonVal.Value = v.AsStringSlice()
	case attribute.MAP:
		m := v.AsMap()
		values := make([]keyValue, 0, len(m))
		for _, kv := range m {
			values = append(values, keyValue{
				Key:   string(kv.Key),
				Value: newValue(kv.Value),
			})
		}

		jsonVal.Value = values
	case attribute.SLICE:
		s := v.AsSlice()
		values := make([]value, 0, len(s))
		for _, e := range s {
			values = append(values, newValue(e))
		}

		jsonVal.Value = values
	case attribute.EMPTY:
		jsonVal.Value = nil
	default:
		return nil, errors.New("invalid attribute.Type")
	}

	return json.Marshal(jsonVal)
}

func valueType(v attribute.Value) string {
	switch v.Type() {
	case attribute.EMPTY:
		return "Empty"
	case attribute.BOOL:
		return "Bool"
	case attribute.INT64:
		return "Int64"
	case attribute.FLOAT64:
		return "Float64"
	case attribute.STRING:
		return "String"
	case attribute.BYTESLICE:
		return "Bytes"
	case attribute.SLICE:
		return "Slice"
	case attribute.MAP:
		return "Map"
	case attribute.BOOLSLICE:
		return "BoolSlice"
	case attribute.INT64SLICE:
		return "Int64Slice"
	case attribute.FLOAT64SLICE:
		return "Float64Slice"
	case attribute.STRINGSLICE:
		return "StringSlice"
	default:
		return "Invalid"
	}
}

type keyValue struct {
	Key   string
	Value value
}

// recordJSON is a JSON-serializable representation of a Record.
type recordJSON struct {
	Timestamp         *time.Time `json:",omitempty"`
	ObservedTimestamp *time.Time `json:",omitempty"`
	EventName         string     `json:",omitempty"`
	Severity          log.Severity
	SeverityText      string
	Body              value
	Attributes        []keyValue
	TraceID           trace.TraceID
	SpanID            trace.SpanID
	TraceFlags        trace.TraceFlags
	Resource          *resource.Resource
	Scope             instrumentation.Scope
	DroppedAttributes int
}

func (e *Exporter) newRecordJSON(r sdklog.Record) recordJSON {
	res := r.Resource()
	newRecord := recordJSON{
		EventName:    r.EventName(),
		Severity:     r.Severity(),
		SeverityText: r.SeverityText(),
		Body:         newValue(r.Body()),

		TraceID:    r.TraceID(),
		SpanID:     r.SpanID(),
		TraceFlags: r.TraceFlags(),

		Attributes: make([]keyValue, 0, r.AttributesLen()),

		Resource: res,
		Scope:    r.InstrumentationScope(),

		DroppedAttributes: r.DroppedAttributes(),
	}

	r.WalkAttributes(func(kv attribute.KeyValue) bool {
		newRecord.Attributes = append(newRecord.Attributes, keyValue{
			Key:   string(kv.Key),
			Value: newValue(kv.Value),
		})
		return true
	})

	if e.timestamps {
		timestamp := r.Timestamp()
		newRecord.Timestamp = &timestamp

		observedTimestamp := r.ObservedTimestamp()
		newRecord.ObservedTimestamp = &observedTimestamp
	}

	return newRecord
}
