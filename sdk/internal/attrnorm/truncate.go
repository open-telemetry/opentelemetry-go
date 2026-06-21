// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package attrnorm

import (
	"reflect"
	"strings"
	"unicode/utf8"
	"unsafe"

	"go.opentelemetry.io/otel/attribute"
)

// StringNeedsTruncation reports whether s would be modified by Truncate for
// the given limit.
func StringNeedsTruncation(limit int, s string) bool {
	if limit < 0 || len(s) <= limit {
		return false
	}
	return utf8.RuneCountInString(s) > limit || !utf8.ValidString(s)
}

// Truncate returns a truncated version of s such that it contains less than
// the limit number of characters. Truncation is applied by returning the limit
// number of valid characters contained in s.
//
// If limit is negative, it returns the original string.
//
// UTF-8 is supported. When truncating, all invalid characters are dropped
// before applying truncation.
//
// If s already contains less than the limit number of bytes, it is returned
// unchanged. No invalid characters are removed.
func Truncate(limit int, s string) string {
	// This prioritize performance in the following order based on the most
	// common expected use-cases.
	//
	//  - Short values less than the default limit (128).
	//  - Strings with valid encodings that exceed the limit.
	//  - No limit.
	//  - Strings with invalid encodings that exceed the limit.
	if limit < 0 || len(s) <= limit {
		return s
	}

	// Optimistically, assume all valid UTF-8.
	var b strings.Builder
	count := 0
	for i, c := range s {
		if c != utf8.RuneError {
			count++
			if count > limit {
				return s[:i]
			}
			continue
		}

		_, size := utf8.DecodeRuneInString(s[i:])
		if size == 1 {
			// Invalid encoding.
			b.Grow(len(s) - 1)
			_, _ = b.WriteString(s[:i])
			s = s[i:]
			break
		}
	}

	// Fast-path, no invalid input.
	if b.Cap() == 0 {
		return s
	}

	// Truncate while validating UTF-8.
	for i := 0; i < len(s) && count < limit; {
		c := s[i]
		if c < utf8.RuneSelf {
			// Optimization for single byte runes (common case).
			_ = b.WriteByte(c)
			i++
			count++
			continue
		}

		_, size := utf8.DecodeRuneInString(s[i:])
		if size == 1 {
			// We checked for all 1-byte runes above, this is a RuneError.
			i++
			continue
		}

		_, _ = b.WriteString(s[i : i+size])
		i += size
		count++
	}

	return b.String()
}

// rawAttrValue mirrors the internal layout of attribute.Value. It is used
// only to read the immutable backing storage of STRINGSLICE values directly,
// avoiding the allocation incurred by returning a copy.
type rawAttrValue struct {
	vtype    attribute.Type
	numeric  uint64
	stringly string
	slice    any
}

// attrValueSlice returns the slice backing storage of v directly from memory.
func attrValueSlice(v attribute.Value) any {
	return (*rawAttrValue)(
		unsafe.Pointer(&v),
	).slice //nolint:gosec // Read-only mirror of attribute.Value; only the immutable backing storage is read.
}

// StringSliceNeedsTruncation reports whether any element in the STRINGSLICE
// value v would be modified by Truncate for the given limit.
//
// It reads the backing storage of v directly to avoid the copy allocation
// that any public accessor incurs on the no-op path.
func StringSliceNeedsTruncation(limit int, v attribute.Value) bool {
	switch ss := attrValueSlice(v).(type) {
	case [0]string:
		return false
	case [1]string:
		return StringNeedsTruncation(limit, ss[0])
	case [2]string:
		return StringNeedsTruncation(limit, ss[0]) || StringNeedsTruncation(limit, ss[1])
	case [3]string:
		return StringNeedsTruncation(limit, ss[0]) || StringNeedsTruncation(limit, ss[1]) ||
			StringNeedsTruncation(limit, ss[2])
	default:
		// 4+ elements are stored as a reflect-allocated [N]string array.
		// rv.Index(i).String() reads each string directly without allocating.
		rv := reflect.ValueOf(attrValueSlice(v))
		if !rv.IsValid() || rv.Kind() != reflect.Array {
			return false
		}
		for i := range rv.Len() {
			if StringNeedsTruncation(limit, rv.Index(i).String()) {
				return true
			}
		}
		return false
	}
}
