package trace

import (
	"regexp"
	"strings"
	"testing"
)

func FuzzTraceIDFromHex(f *testing.F) {
	// Seed corpus with valid and edge-case examples
	f.Add("00000000000000000000000000000001") // lowest valid (non-zero)
	f.Add("0123456789abcdef0123456789abcdef")
	f.Add("ffffffffffffffffffffffffffffffff")
	f.Add("0123456789abcdefabcdefabcdefabcd")
	f.Add("invalidhexstringnot32chars")

	f.Fuzz(func(t *testing.T, s string) {
		id, err := TraceIDFromHex(s)

		// OTel-valid TraceIDs: 32 lowercase hex chars, not all zeros.
		isValidTraceID := regexp.MustCompile(`^[0-9a-f]{32}$`).MatchString(s) && !strings.EqualFold(s, "00000000000000000000000000000000")

		if isValidTraceID {
			if err != nil {
				t.Errorf("expected no error for valid hex input: %q, got err: %v", s, err)
				return
			}
		}

		// Invalid input is fine for fuzzing — skip further checks.
		if err != nil {
			return
		}

		got := id.String()

		// TraceIDFromHex normalizes input to lowercase.
		if got != strings.ToLower(s) {
			t.Errorf("roundtrip mismatch: in=%q out=%q", s, got)
		}
	})
}

func FuzzSpanIDFromHex(f *testing.F) {
	// Seed corpus with valid and edge-case examples
	f.Add("0000000000000001") // lowest valid (non-zero)
	f.Add("0123456789abcdef")
	f.Add("ffffffffffffffff")
	f.Add("abcdefabcdefabcd")
	f.Add("invalidhex")

	f.Fuzz(func(t *testing.T, s string) {
		id, err := SpanIDFromHex(s)

		// OTel-valid SpanIDs: 16 lowercase hex chars, not all zeros.
		isValidSpanID := regexp.MustCompile(`^[0-9a-f]{16}$`).MatchString(s) && !strings.EqualFold(s, "0000000000000000")

		if isValidSpanID {
			if err != nil {
				t.Errorf("expected no error for valid hex input: %q, got err: %v", s, err)
				return
			}
		}

		// Invalid input is fine for fuzzing — skip further checks.
		if err != nil {
			return
		}

		got := id.String()

		// SpanIDFromHex normalizes input to lowercase.
		if got != strings.ToLower(s) {
			t.Errorf("roundtrip mismatch: in=%q out=%q", s, got)
		}
	})
}
