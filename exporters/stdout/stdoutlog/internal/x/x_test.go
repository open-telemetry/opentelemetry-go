// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelfObservabilityFeature(t *testing.T) {
	testCases := []struct {
		name     string
		envValue string
		enabled  bool
	}{
		{
			name:     "enabled_lowercase",
			envValue: "true",
			enabled:  true,
		},
		{
			name:     "enabled_uppercase",
			envValue: "TRUE",
			enabled:  true,
		},
		{
			name:     "enabled_mixed_case",
			envValue: "True",
			enabled:  true,
		},
		{
			name:     "disabled_false",
			envValue: "false",
			enabled:  false,
		},
		{
			name:     "disabled_invalid",
			envValue: "invalid",
			enabled:  false,
		},
		{
			name:     "disabled_empty",
			envValue: "",
			enabled:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.envValue != "" {
				t.Setenv(SelfObservability.Key(), tc.envValue)
			}

			assert.Equal(t, tc.enabled, SelfObservability.Enabled())

			value, ok := SelfObservability.Lookup()
			if tc.enabled {
				assert.True(t, ok)
				assert.Equal(t, tc.envValue, value)
			} else {
				assert.False(t, ok)
				assert.Empty(t, value)
			}
		})
	}
}
