// Code created by gotmpl. DO NOT MODIFY.
// source: internal/shared/otlp/otlplog/transform/log.go.tmpl

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package transform provides transformation functionality from the
// sdk/log data-types into OTLP data-types.
package transform // import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc/internal/transform"

import (
	"time"

	cpb "go.opentelemetry.io/proto/otlp/common/v1"
	lpb "go.opentelemetry.io/proto/otlp/logs/v1"
	rpb "go.opentelemetry.io/proto/otlp/resource/v1"

	"go.opentelemetry.io/otel/attribute"
	api "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/log"
)

// ResourceLogs returns an slice of OTLP ResourceLogs generated from records.
func ResourceLogs(records []log.Record) []*lpb.ResourceLogs {
	if len(records) == 0 {
		return nil
	}

	resMap := make(map[attribute.Distinct]*lpb.ResourceLogs)

	type key struct {
		r  attribute.Distinct
		is instrumentation.Scope
	}
	scopeMap := make(map[key]*lpb.ScopeLogs)

	var resources int
	for _, r := range records {
		res := r.Resource()
		rKey := res.Equivalent()
		scope := r.InstrumentationScope()
		k := key{
			r:  rKey,
			is: scope,
		}
		sl, iOk := scopeMap[k]
		if !iOk {
			sl = new(lpb.ScopeLogs)
			var emptyScope instrumentation.Scope
			if scope != emptyScope {
				sl.Scope = &cpb.InstrumentationScope{
					Name:       scope.Name,
					Version:    scope.Version,
					Attributes: AttrIter(scope.Attributes.Iter()),
				}
				sl.SchemaUrl = scope.SchemaURL
			}
			scopeMap[k] = sl
		}

		sl.LogRecords = append(sl.LogRecords, LogRecord(r))
		rl, rOk := resMap[rKey]
		if !rOk {
			resources++
			rl = new(lpb.ResourceLogs)
			if res.Len() > 0 {
				rl.Resource = &rpb.Resource{
					Attributes: AttrIter(res.Iter()),
				}
			}
			rl.SchemaUrl = res.SchemaURL()
			resMap[rKey] = rl
		}
		if !iOk {
			rl.ScopeLogs = append(rl.ScopeLogs, sl)
		}
	}

	// Transform the categorized map into a slice
	resLogs := make([]*lpb.ResourceLogs, 0, resources)
	for _, rl := range resMap {
		resLogs = append(resLogs, rl)
	}

	return resLogs
}

// LogRecord returns an OTLP LogRecord generated from record.
func LogRecord(record log.Record) *lpb.LogRecord {
	r := &lpb.LogRecord{
		TimeUnixNano:         timeUnixNano(record.Timestamp()),
		ObservedTimeUnixNano: timeUnixNano(record.ObservedTimestamp()),
		EventName:            record.EventName(),
		SeverityNumber:       SeverityNumber(record.Severity()),
		SeverityText:         record.SeverityText(),
		Body:                 LogAttrValue(record.Body()),
		Attributes:           make([]*cpb.KeyValue, 0, record.AttributesLen()),
		Flags:                uint32(record.TraceFlags()),
		// TODO: DroppedAttributesCount: /* ... */,
	}
	record.WalkAttributes(func(kv api.KeyValue) bool {
		r.Attributes = append(r.Attributes, LogAttr(kv))
		return true
	})
	if tID := record.TraceID(); tID.IsValid() {
		r.TraceId = tID[:]
	}
	if sID := record.SpanID(); sID.IsValid() {
		r.SpanId = sID[:]
	}
	return r
}

// timeUnixNano returns t as a Unix time, the number of nanoseconds elapsed
// since January 1, 1970 UTC as uint64. The result is undefined if the Unix
// time in nanoseconds cannot be represented by an int64 (a date before the
// year 1678 or after 2262). timeUnixNano on the zero Time returns 0. The
// result does not depend on the location associated with t.
func timeUnixNano(t time.Time) uint64 {
	nano := t.UnixNano()
	if nano < 0 {
		return 0
	}
	return uint64(nano) // nolint:gosec // Overflow checked.
}

// AttrIter transforms an [attribute.Iterator] into OTLP key-values.
func AttrIter(iter attribute.Iterator) []*cpb.KeyValue {
	l := iter.Len()
	if l == 0 {
		return nil
	}

	out := make([]*cpb.KeyValue, 0, l)
	for iter.Next() {
		out = append(out, Attr(iter.Attribute()))
	}
	return out
}

// Attrs transforms a slice of [attribute.KeyValue] into OTLP key-values.
func Attrs(attrs []attribute.KeyValue) []*cpb.KeyValue {
	if len(attrs) == 0 {
		return nil
	}

	out := make([]*cpb.KeyValue, 0, len(attrs))
	for _, kv := range attrs {
		out = append(out, Attr(kv))
	}
	return out
}

// Attr transforms an [attribute.KeyValue] into an OTLP key-value.
func Attr(kv attribute.KeyValue) *cpb.KeyValue {
	return &cpb.KeyValue{Key: string(kv.Key), Value: AttrValue(kv.Value)}
}

// AttrValue transforms an [attribute.Value] into an OTLP AnyValue.
func AttrValue(v attribute.Value) *cpb.AnyValue {
	av := new(cpb.AnyValue)
	switch v.Type() {
	case attribute.BOOL:
		av.Value = &cpb.AnyValue_BoolValue{
			BoolValue: v.AsBool(),
		}
	case attribute.BOOLSLICE:
		av.Value = &cpb.AnyValue_ArrayValue{
			ArrayValue: &cpb.ArrayValue{
				Values: boolSliceValues(v.AsBoolSlice()),
			},
		}
	case attribute.INT64:
		av.Value = &cpb.AnyValue_IntValue{
			IntValue: v.AsInt64(),
		}
	case attribute.INT64SLICE:
		av.Value = &cpb.AnyValue_ArrayValue{
			ArrayValue: &cpb.ArrayValue{
				Values: int64SliceValues(v.AsInt64Slice()),
			},
		}
	case attribute.FLOAT64:
		av.Value = &cpb.AnyValue_DoubleValue{
			DoubleValue: v.AsFloat64(),
		}
	case attribute.FLOAT64SLICE:
		av.Value = &cpb.AnyValue_ArrayValue{
			ArrayValue: &cpb.ArrayValue{
				Values: float64SliceValues(v.AsFloat64Slice()),
			},
		}
	case attribute.STRING:
		av.Value = &cpb.AnyValue_StringValue{
			StringValue: v.AsString(),
		}
	case attribute.STRINGSLICE:
		av.Value = &cpb.AnyValue_ArrayValue{
			ArrayValue: &cpb.ArrayValue{
				Values: stringSliceValues(v.AsStringSlice()),
			},
		}
	default:
		av.Value = &cpb.AnyValue_StringValue{
			StringValue: "INVALID",
		}
	}
	return av
}

func boolSliceValues(vals []bool) []*cpb.AnyValue {
	converted := make([]*cpb.AnyValue, len(vals))
	for i, v := range vals {
		converted[i] = &cpb.AnyValue{
			Value: &cpb.AnyValue_BoolValue{
				BoolValue: v,
			},
		}
	}
	return converted
}

func int64SliceValues(vals []int64) []*cpb.AnyValue {
	converted := make([]*cpb.AnyValue, len(vals))
	for i, v := range vals {
		converted[i] = &cpb.AnyValue{
			Value: &cpb.AnyValue_IntValue{
				IntValue: v,
			},
		}
	}
	return converted
}

func float64SliceValues(vals []float64) []*cpb.AnyValue {
	converted := make([]*cpb.AnyValue, len(vals))
	for i, v := range vals {
		converted[i] = &cpb.AnyValue{
			Value: &cpb.AnyValue_DoubleValue{
				DoubleValue: v,
			},
		}
	}
	return converted
}

func stringSliceValues(vals []string) []*cpb.AnyValue {
	converted := make([]*cpb.AnyValue, len(vals))
	for i, v := range vals {
		converted[i] = &cpb.AnyValue{
			Value: &cpb.AnyValue_StringValue{
				StringValue: v,
			},
		}
	}
	return converted
}

// LogAttrs transforms a slice of [api.KeyValue] into OTLP key-values.
func LogAttrs(attrs []api.KeyValue) []*cpb.KeyValue {
	if len(attrs) == 0 {
		return nil
	}

	out := make([]*cpb.KeyValue, 0, len(attrs))
	for _, kv := range attrs {
		out = append(out, LogAttr(kv))
	}
	return out
}

// LogAttr transforms an [api.KeyValue] into an OTLP key-value.
func LogAttr(attr api.KeyValue) *cpb.KeyValue {
	return &cpb.KeyValue{
		Key:   attr.Key,
		Value: LogAttrValue(attr.Value),
	}
}

// LogAttrValues transforms a slice of [api.Value] into an OTLP []AnyValue.
func LogAttrValues(vals []api.Value) []*cpb.AnyValue {
	if len(vals) == 0 {
		return nil
	}

	out := make([]*cpb.AnyValue, 0, len(vals))
	for _, v := range vals {
		out = append(out, LogAttrValue(v))
	}
	return out
}

// LogAttrValue transforms an [api.Value] into an OTLP AnyValue.
func LogAttrValue(v api.Value) *cpb.AnyValue {
	av := new(cpb.AnyValue)
	switch v.Kind() {
	case api.KindBool:
		av.Value = &cpb.AnyValue_BoolValue{
			BoolValue: v.AsBool(),
		}
	case api.KindInt64:
		av.Value = &cpb.AnyValue_IntValue{
			IntValue: v.AsInt64(),
		}
	case api.KindFloat64:
		av.Value = &cpb.AnyValue_DoubleValue{
			DoubleValue: v.AsFloat64(),
		}
	case api.KindString:
		av.Value = &cpb.AnyValue_StringValue{
			StringValue: v.AsString(),
		}
	case api.KindBytes:
		av.Value = &cpb.AnyValue_BytesValue{
			BytesValue: v.AsBytes(),
		}
	case api.KindSlice:
		av.Value = &cpb.AnyValue_ArrayValue{
			ArrayValue: &cpb.ArrayValue{
				Values: LogAttrValues(v.AsSlice()),
			},
		}
	case api.KindMap:
		av.Value = &cpb.AnyValue_KvlistValue{
			KvlistValue: &cpb.KeyValueList{
				Values: LogAttrs(v.AsMap()),
			},
		}
	default:
		av.Value = &cpb.AnyValue_StringValue{
			StringValue: "INVALID",
		}
	}
	return av
}

// SeverityNumber transforms a [log.Severity] into an OTLP SeverityNumber.
func SeverityNumber(s api.Severity) lpb.SeverityNumber {
	switch s {
	case api.SeverityTrace:
		return lpb.SeverityNumber_SEVERITY_NUMBER_TRACE
	case api.SeverityTrace2:
		return lpb.SeverityNumber_SEVERITY_NUMBER_TRACE2
	case api.SeverityTrace3:
		return lpb.SeverityNumber_SEVERITY_NUMBER_TRACE3
	case api.SeverityTrace4:
		return lpb.SeverityNumber_SEVERITY_NUMBER_TRACE4
	case api.SeverityDebug:
		return lpb.SeverityNumber_SEVERITY_NUMBER_DEBUG
	case api.SeverityDebug2:
		return lpb.SeverityNumber_SEVERITY_NUMBER_DEBUG2
	case api.SeverityDebug3:
		return lpb.SeverityNumber_SEVERITY_NUMBER_DEBUG3
	case api.SeverityDebug4:
		return lpb.SeverityNumber_SEVERITY_NUMBER_DEBUG4
	case api.SeverityInfo:
		return lpb.SeverityNumber_SEVERITY_NUMBER_INFO
	case api.SeverityInfo2:
		return lpb.SeverityNumber_SEVERITY_NUMBER_INFO2
	case api.SeverityInfo3:
		return lpb.SeverityNumber_SEVERITY_NUMBER_INFO3
	case api.SeverityInfo4:
		return lpb.SeverityNumber_SEVERITY_NUMBER_INFO4
	case api.SeverityWarn:
		return lpb.SeverityNumber_SEVERITY_NUMBER_WARN
	case api.SeverityWarn2:
		return lpb.SeverityNumber_SEVERITY_NUMBER_WARN2
	case api.SeverityWarn3:
		return lpb.SeverityNumber_SEVERITY_NUMBER_WARN3
	case api.SeverityWarn4:
		return lpb.SeverityNumber_SEVERITY_NUMBER_WARN4
	case api.SeverityError:
		return lpb.SeverityNumber_SEVERITY_NUMBER_ERROR
	case api.SeverityError2:
		return lpb.SeverityNumber_SEVERITY_NUMBER_ERROR2
	case api.SeverityError3:
		return lpb.SeverityNumber_SEVERITY_NUMBER_ERROR3
	case api.SeverityError4:
		return lpb.SeverityNumber_SEVERITY_NUMBER_ERROR4
	case api.SeverityFatal:
		return lpb.SeverityNumber_SEVERITY_NUMBER_FATAL
	case api.SeverityFatal2:
		return lpb.SeverityNumber_SEVERITY_NUMBER_FATAL2
	case api.SeverityFatal3:
		return lpb.SeverityNumber_SEVERITY_NUMBER_FATAL3
	case api.SeverityFatal4:
		return lpb.SeverityNumber_SEVERITY_NUMBER_FATAL4
	}
	return lpb.SeverityNumber_SEVERITY_NUMBER_UNSPECIFIED
}
