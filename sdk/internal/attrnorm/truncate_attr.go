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
