// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package attrnorm

import (
	"slices"

	"go.opentelemetry.io/otel/attribute"
)

// NeedsTruncation reports whether v would be modified by TruncateValue for
// the given limit.
func NeedsTruncation(limit int, v attribute.Value) bool {
	switch v.Type() {
	case attribute.STRING:
		return StringNeedsTruncation(limit, v.AsString())
	case attribute.BYTESLICE:
		// len(v.AsString()) is identical to len(v.AsByteSlice()) but
		// avoids memory allocation.
		if limit >= 0 && len(v.AsString()) > limit {
			return true
		}
	case attribute.STRINGSLICE:
		return StringSliceNeedsTruncation(limit, v)
	case attribute.SLICE:
		return slices.ContainsFunc(v.AsSlice(), func(e attribute.Value) bool { return NeedsTruncation(limit, e) })
	case attribute.MAP:
		return slices.ContainsFunc(
			v.AsMap(),
			func(kv attribute.KeyValue) bool { return NeedsTruncation(limit, kv.Value) },
		)
	}
	return false
}

// TruncateAttr returns a truncated version of attr. Only string, string
// slice, byte slice, slice, and map attribute values are truncated. String
// values are truncated to at most a length of limit. Each string slice value
// is truncated in this fashion (the slice length itself is unaffected), and
// byte slice values are truncated to at most limit bytes. For slice and map
// attribute values, the limit is applied recursively to contained values.
//
// No truncation is performed for a negative limit.
func TruncateAttr(limit int, attr attribute.KeyValue) attribute.KeyValue {
	if limit < 0 {
		return attr
	}
	switch attr.Value.Type() {
	case attribute.STRING:
		return attr.Key.String(Truncate(limit, attr.Value.AsString()))
	case attribute.STRINGSLICE:
		if !StringSliceNeedsTruncation(limit, attr.Value) {
			return attr
		}
		v := attr.Value.AsStringSlice()
		for i := range v {
			v[i] = Truncate(limit, v[i])
		}
		return attr.Key.StringSlice(v)
	case attribute.BYTESLICE:
		v := attr.Value.AsString()
		if len(v) > limit {
			return attr.Key.ByteSlice([]byte(v[:limit]))
		}
		return attr
	case attribute.SLICE:
		v := attr.Value.AsSlice()
		if !slices.ContainsFunc(v, func(e attribute.Value) bool { return NeedsTruncation(limit, e) }) {
			return attr
		}
		newV := make([]attribute.Value, len(v))
		for i, elem := range v {
			newV[i] = TruncateValue(limit, elem)
		}
		return attr.Key.Slice(newV...)
	case attribute.MAP:
		v := attr.Value.AsMap()
		if !slices.ContainsFunc(v, func(kv attribute.KeyValue) bool { return NeedsTruncation(limit, kv.Value) }) {
			return attr
		}
		newV := make([]attribute.KeyValue, len(v))
		for i, elem := range v {
			elem.Value = TruncateValue(limit, elem.Value)
			newV[i] = elem
		}
		return attr.Key.Map(newV...)
	}
	return attr
}

// TruncateValue returns a truncated version of v. Only string, string slice,
// byte slice, and (recursively) slice and map values are modified.
//
// No truncation is performed for a negative limit.
func TruncateValue(limit int, v attribute.Value) attribute.Value {
	switch v.Type() {
	case attribute.STRING:
		return attribute.StringValue(Truncate(limit, v.AsString()))
	case attribute.STRINGSLICE:
		if !StringSliceNeedsTruncation(limit, v) {
			return v
		}
		ss := v.AsStringSlice()
		for i := range ss {
			ss[i] = Truncate(limit, ss[i])
		}
		return attribute.StringSliceValue(ss)
	case attribute.BYTESLICE:
		// len(v.AsString()) is identical to len(v.AsByteSlice()) but
		// avoids allocating the full slice before truncation.
		s := v.AsString()
		if limit >= 0 && len(s) > limit {
			return attribute.ByteSliceValue([]byte(s[:limit]))
		}
	case attribute.SLICE:
		sl := v.AsSlice()
		if !slices.ContainsFunc(sl, func(e attribute.Value) bool { return NeedsTruncation(limit, e) }) {
			return v
		}
		newSl := make([]attribute.Value, len(sl))
		for i, elem := range sl {
			newSl[i] = TruncateValue(limit, elem)
		}
		return attribute.SliceValue(newSl...)
	case attribute.MAP:
		m := v.AsMap()
		if !slices.ContainsFunc(m, func(kv attribute.KeyValue) bool { return NeedsTruncation(limit, kv.Value) }) {
			return v
		}
		newM := make([]attribute.KeyValue, len(m))
		for i, elem := range m {
			elem.Value = TruncateValue(limit, elem.Value)
			newM[i] = elem
		}
		return attribute.MapValue(newM...)
	}
	return v
}
