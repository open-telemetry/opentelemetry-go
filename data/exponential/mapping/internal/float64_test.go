// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
