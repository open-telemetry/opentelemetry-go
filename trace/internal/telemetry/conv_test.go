// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package telemetry

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const schema100 = "http://go.opentelemetry.io/schema/v1.0.0"

var y2k = time.Unix(0, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano()) // No location.

func runJSONEncodingTests[T any](decoded T, encoded []byte) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		t.Run("Unmarshal", runJSONUnmarshalTest(decoded, encoded))
		t.Run("Marshal", runJSONMarshalTest(decoded, encoded))
	}
}

func runJSONMarshalTest[T any](decoded T, encoded []byte) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		got, err := json.Marshal(decoded)
		require.NoError(t, err)

		var want bytes.Buffer
		require.NoError(t, json.Compact(&want, encoded))
		assert.Equal(t, want.String(), string(got))
	}
}

func runJSONUnmarshalTest[T any](decoded T, encoded []byte) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		var got T
		require.NoError(t, json.Unmarshal(encoded, &got))
		assert.Equal(t, decoded, got)
	}
}
