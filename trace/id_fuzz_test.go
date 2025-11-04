package trace

import (
	"strings"
	"testing"
)

func FuzzTraceIDFromHex(f *testing.F) {
	// seed with valid-looking values
	f.Add("00000000000000000000000000000000")
	f.Add("0123456789abcdef0123456789abcdef")
	f.Add("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")

	f.Fuzz(func(t *testing.T, s string) {
		id, err := TraceIDFromHex(s)
		if err != nil {
			// invalid input is fine, that’s exactly why we fuzz
			return
		}

		// for valid-length inputs, String() should give a 32-char lowercase hex
		got := id.String()
		if len(s) == 32 {
			if len(got) != 32 {
				t.Fatalf("TraceIDFromHex(%q) -> String() len = %d, want 32", s, len(got))
			}
			if got != strings.ToLower(s) {
				t.Fatalf("roundtrip mismatch: in=%q out=%q", s, got)
			}
		}
	})
}

func FuzzSpanIDFromHex(f *testing.F) {
	// span IDs are 8 bytes => 16 hex chars
	f.Add("0000000000000000")
	f.Add("0123456789abcdef")
	f.Add("FFFFFFFFFFFFFFFF")

	f.Fuzz(func(t *testing.T, s string) {
		id, err := SpanIDFromHex(s)
		if err != nil {
			return
		}

		got := id.String()
		if len(s) == 16 {
			if len(got) != 16 {
				t.Fatalf("SpanIDFromHex(%q) -> String() len = %d, want 16", s, len(got))
			}
			if got != strings.ToLower(s) {
				t.Fatalf("roundtrip mismatch: in=%q out=%q", s, got)
			}
		}
	})
}
