package trace

import (
	"regexp"
	"strings"
	"testing"
)

func FuzzTraceIDFromHex(f *testing.F) {
	// Seed corpus with valid and edge-case examples.
	f.Add("00000000000000000000000000000001") // Lowest valid (non-zero).
	f.Add("0123456789abcdef0123456789abcdef")
	f.Add("ffffffffffffffffffffffffffffffff") // Highest valid.
	f.Add("0123456789abcdefabcdefabcdefabcd")
	f.Add("invalidhexstringnot32chars") // Invalid.

	// Precompile regex for efficiency.
	validTraceIDRe := regexp.MustCompile(`^[0-9a-f]{32}$`)

	f.Fuzz(func(t *testing.T, s string) {
		id, err := TraceIDFromHex(s)

		// OTel-valid TraceIDs: 32 lowercase hex chars, not all zeros.
		isValidTraceID := validTraceIDRe.MatchString(s) && !strings.EqualFold(s, "00000000000000000000000000000000")

		if isValidTraceID {
			if err != nil {
				t.Errorf("expected no error for valid hex input: %q, got err: %v", s, err)
				t.Fatalf("expected no error for valid hex input: %q, got err: %v", s, err)
				return
			}
		}

		if !isValidTraceID && err == nil {
			t.Errorf("expected error for invalid input: %q", s)
			return
		}

		if err != nil {
			return
		}

		got := id.String()

		// TraceIDFromHex normalizes input to lowercase.
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

	// Precompile regex for efficiency.
	validSpanIDRe := regexp.MustCompile(`^[0-9a-f]{16}$`)

	f.Fuzz(func(t *testing.T, s string) {
		id, err := SpanIDFromHex(s)

		// OTel-valid SpanIDs: 16 lowercase hex chars, not all zeros.
		isValidSpanID := validSpanIDRe.MatchString(s) && !strings.EqualFold(s, "0000000000000000")

		if isValidSpanID {
			if err != nil {
				t.Errorf("expected no error for valid hex input: %q, got err: %v", s, err)
				t.Fatalf("expected no error for valid hex input: %q, got err: %v", s, err)
				return
			}
		}

		if !isValidSpanID && err == nil {
			t.Errorf("expected error for invalid input: %q", s)
			return
		}

		if err != nil {
			return
		}

		got := id.String()

		// SpanIDFromHex normalizes input to lowercase.
		if got != s {
			t.Errorf("roundtrip mismatch: in=%q out=%q", s, got)
		}
	})
}
