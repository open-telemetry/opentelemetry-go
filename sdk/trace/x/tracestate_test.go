// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTracestateRandomness(t *testing.T) {
	const validRV = "0123456789abcd"
	const validRVValue uint64 = 0x0123456789abcd

	testCases := []struct {
		name       string
		otts       string
		wantRandom uint64
		wantHasRV  bool
	}{
		{"rv at beginning", "rv:" + validRV, validRVValue, true},
		{"rv at beginning with more keys", "rv:" + validRV + ";th:0;other:value", validRVValue, true},
		{"rv in middle", "th:0;rv:" + validRV + ";other:value", validRVValue, true},
		{"rv at end", "th:0;other:value;rv:" + validRV, validRVValue, true},
		{"rv with max 56-bit value", "rv:0fffffffffffff", 0x0fffffffffffff, true},
		{"no rv key", "th:0;other:value", 0, false},
		{"empty string", "", 0, false},
		{"rv value too short", "rv:0123456789abc", 0, false},
		{"rv value too long", "rv:0123456789abcdef", 0, false},
		{"rv with invalid hex", "rv:0123456789abcg", 0, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotRandom, gotHasRV := tracestateRandomness(tc.otts)
			assert.Equal(t, tc.wantHasRV, gotHasRV)
			if tc.wantHasRV {
				assert.Equal(t, tc.wantRandom, gotRandom)
			}
		})
	}
}

func TestEraseTraceStateThKeyValue(t *testing.T) {
	testCases := []struct {
		name string
		otts string
		want string
	}{
		{"empty string", "", ""},
		{"no th in existing", "rv:0123456789abcd;other:value", "rv:0123456789abcd;other:value"},
		{"only th returns empty", "th:0ad", ""},
		{"th at front", "th:0ad;rv:0123456789abcd", "rv:0123456789abcd"},
		{"th in middle", "rv:0123456789abcd;th:0ad;other:value", "rv:0123456789abcd;other:value"},
		{"th at end", "rv:0123456789abcd;th:0ad", "rv:0123456789abcd"},
		{
			"th substring in another key (path:0) is not erased",
			"path:0",
			"path:0",
		},
		{
			"erase real th only when another key ends with th before colon",
			"path:0;th:0ad",
			"path:0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := eraseTraceStateThKeyValue(tc.otts)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestInsertOrUpdateTraceStateThKeyValue(t *testing.T) {
	testCases := []struct {
		name         string
		existingOtts string
		thkv         string
		want         string
	}{
		{"empty existing adds th at front", "", "th:123456789abcd", "th:123456789abcd"},
		{
			"no th in existing adds th at front",
			"rv:0123456789abcd;other:value",
			"th:fedcba987654321",
			"th:fedcba987654321;rv:0123456789abcd;other:value",
		},
		{
			"existing th is replaced",
			"rv:0123456789abcd;th:0ad;other:value",
			"th:0e1",
			"th:0e1;rv:0123456789abcd;other:value",
		},
		{"th at front is replaced", "th:0ad;rv:0123456789abcd", "th:0e1", "th:0e1;rv:0123456789abcd"},
		{"only th in existing is replaced", "th:0ad", "th:0e1", "th:0e1"},
		{
			"th substring in another key (path:0) is left intact; th prepended",
			"path:0;rv:0123456789abcd",
			"th:fedcba987654321",
			"th:fedcba987654321;path:0;rv:0123456789abcd",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := InsertOrUpdateTraceStateThKeyValue(tc.existingOtts, tc.thkv)
			assert.Equal(t, tc.want, got)
		})
	}
}
