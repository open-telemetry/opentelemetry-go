// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"regexp"
	"testing"
)

func FuzzTraceIDFromHex(f *testing.F) {
	// Seed corpus with valid and edge-case examples.
	f.Add("00000000000000000000000000000001") // Lowest valid (non-zero).
	f.Add("0123456789abcdef0123456789abcdef")
	f.Add("ffffffffffffffffffffffffffffffff") // Highest valid.
	f.Add("0123456789abcdefabcdefabcdefabcd")
	f.Add("invalidhexstringnot32chars") // Invalid.

	validTraceIDRe := regexp.MustCompile(`^[0-9a-f]{32}$`)

	f.Fuzz(func(t *testing.T, s string) {
		id, err := TraceIDFromHex(s)

		// OTel-valid TraceIDs: 32 lowercase hex chars, not all zeros.
		validTraceID := validTraceIDRe.MatchString(s) && s != "00000000000000000000000000000000"
		if validTraceID && err != nil {
			t.Fatalf("expected no error for valid hex input: %q, got err: %v", s, err)
		}
		if !validTraceID && err == nil {
			t.Fatalf("expected error for invalid input: %q", s)
		}
		if err != nil {
			return
		}

		got := id.String()
		if got != s {
			t.Errorf("roundtrip mismatch: in=%q out=%q", s, got)
		}
	})
}

func FuzzSpanIDFromHex(f *testing.F) {
	// Seed corpus with valid and edge-case examples.
	f.Add("0000000000000001") // Lowest valid (non-zero).
	f.Add("0123456789abcdef")
	f.Add("ffffffffffffffff") // Highest valid.
	f.Add("abcdefabcdefabcd")
	f.Add("invalidhex") // Invalid.

	validSpanIDRe := regexp.MustCompile(`^[0-9a-f]{16}$`)

	f.Fuzz(func(t *testing.T, s string) {
		id, err := SpanIDFromHex(s)

		// OTel-valid SpanIDs: 16 lowercase hex chars, not all zeros.
		validSpanID := validSpanIDRe.MatchString(s) && s != "0000000000000000"
		if validSpanID && err != nil {
			t.Fatalf("expected no error for valid hex input: %q, got err: %v", s, err)
		}
		if !validSpanID && err == nil {
			t.Fatalf("expected error for invalid input: %q", s)
		}
		if err != nil {
			return
		}

		got := id.String()
		if got != s {
			t.Errorf("roundtrip mismatch: in=%q out=%q", s, got)
		}
	})
}
