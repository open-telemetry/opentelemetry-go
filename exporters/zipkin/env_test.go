// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package zipkin

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ottest "go.opentelemetry.io/otel/internal/internaltest"
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

	envStore := ottest.NewEnvStore()
	envStore.Record(envEndpoint)
	defer func() {
		require.NoError(t, envStore.Restore())
	}()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, os.Setenv(envEndpoint, tc.envEndpoint))

			endpoint := envOr(envEndpoint, tc.defaultCollectorEndpoint)

			assert.Equal(t, tc.expectedCollectorEndpoint, endpoint)
		})
	}
}
