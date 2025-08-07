// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelfObservability(t *testing.T) {
	assert.Equal(t, "OTEL_GO_X_SELF_OBSERVABILITY", SelfObservability.Key())

	testCases := []struct {
		name     string
		envValue string
		enabled  bool
		value    string
	}{
		{
			name:     "unset",
			envValue: "",
			enabled:  false,
			value:    "",
		},
		{
			name:     "lowercase true",
			envValue: "true",
			enabled:  true,
			value:    "true",
		},
		{
			name:     "uppercase true",
			envValue: "TRUE",
			enabled:  true,
			value:    "TRUE",
		},
		{
			name:     "mixed case true",
			envValue: "True",
			enabled:  true,
			value:    "True",
		},
		{
			name:     "false value",
			envValue: "false",
			enabled:  false,
			value:    "",
		},
		{
			name:     "invalid value",
			envValue: "invalid",
			enabled:  false,
			value:    "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv(SelfObservability.Key(), tc.envValue)

			assert.Equal(t, tc.enabled, SelfObservability.Enabled())

			value, ok := SelfObservability.Lookup()
			assert.Equal(t, tc.enabled, ok)
			assert.Equal(t, tc.value, value)
		})
	}
}
