package trace

import (
	"testing"
)

func BenchmarkDecodeHex(b *testing.B) {
	decodeHex1 := func(h string, b []byte) error {
		for i := 0; i < len(h); i++ {
			r := h[i]
			switch {
			case 'a' <= r && r <= 'f':
				continue
			case '0' <= r && r <= '9':
				continue
			default:
				return errInvalidHexID
			}
		}

		return nil
	}

	decodeHex2 := func(h string, b []byte) error {
		for _, r := range h {
			switch {
			case 'a' <= r && r <= 'f':
				continue
			case '0' <= r && r <= '9':
				continue
			default:
				return errInvalidHexID
			}
		}

		return nil
	}

	tests := []string{
		"0123456789abcdef0123456789abcdef",
		"abcdefabcdefabcdefabcdefabcdefab",
		"00000000000000000000000000000000",
		"ffffffffffffffffffffffffffffffff",
		"1234567890abcdef1234567890abcdef",
	}

	for _, h := range tests {
		b.Run("decodeHex1_"+h, func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()

			for range b.N {
				_ = decodeHex1(h, make([]byte, 16))
			}
		})

		b.Run("decodeHex2_"+h, func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()

			for range b.N {
				_ = decodeHex2(h, make([]byte, 16))
			}
		})
	}
}
