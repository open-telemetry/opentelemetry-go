package internal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Tests that GetNormalBase2 returns the base-2 exponent as documented, unlike
// math.Frexp.
func TestGetNormalBase2(t *testing.T) {
	require.Equal(t, int32(-1022), MinNormalExponent)
	require.Equal(t, int32(+1023), MaxNormalExponent)

	require.Equal(t, MaxNormalExponent, GetNormalBase2(0x1p+1023))
	require.Equal(t, int32(1022), GetNormalBase2(0x1p+1022))

	require.Equal(t, int32(0), GetNormalBase2(1))

	require.Equal(t, int32(-1021), GetNormalBase2(0x1p-1021))
	require.Equal(t, int32(-1022), GetNormalBase2(0x1p-1022))

	// Subnormals below this point
	require.Equal(t, int32(-1023), GetNormalBase2(0x1p-1023))
	require.Equal(t, int32(-1023), GetNormalBase2(0x1p-1024))
	require.Equal(t, int32(-1023), GetNormalBase2(0x1p-1025))
	require.Equal(t, int32(-1023), GetNormalBase2(0x1p-1074))
}

func TestGetSignificand(t *testing.T) {
	// The number 1.5 has a single most-significant bit set, i.e., 1<<51.
	require.Equal(t, int64(1)<<(SignificandWidth-1), GetSignificand(1.5))
}
