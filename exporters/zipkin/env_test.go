// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package zipkin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvOrWithCollectorEndpointOptionsFromEnv(t *testing.T) {
	testCases := []struct {
		name                      string
		envEndpoint               string
		defaultCollectorEndpoint  string
		expectedCollectorEndpoint string
	}{
		{
			name:                      "overrides value via environment variables",
			envEndpoint:               "http://localhost:19411/foo",
			defaultCollectorEndpoint:  defaultCollectorURL,
			expectedCollectorEndpoint: "http://localhost:19411/foo",
		},
		{
			name:                      "environment variables is empty, will not overwrite value",
			envEndpoint:               "",
			defaultCollectorEndpoint:  defaultCollectorURL,
			expectedCollectorEndpoint: defaultCollectorURL,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv(envEndpoint, tc.envEndpoint)

			endpoint := envOr(envEndpoint, tc.defaultCollectorEndpoint)

			assert.Equal(t, tc.expectedCollectorEndpoint, endpoint)
		})
	}
}
