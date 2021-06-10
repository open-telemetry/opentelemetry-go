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

package jaeger

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ottest "go.opentelemetry.io/otel/internal/internaltest"
)

func TestNewRawExporterWithDefault(t *testing.T) {
	const (
		collectorEndpoint = "http://localhost:14268/api/traces"
		username          = ""
		password          = ""
	)

	// Create Jaeger Exporter with default values
	exp, err := New(
		WithCollectorEndpoint(),
	)

	assert.NoError(t, err)

	require.IsType(t, &collectorUploader{}, exp.uploader)
	uploader := exp.uploader.(*collectorUploader)
	assert.Equal(t, collectorEndpoint, uploader.endpoint)
	assert.Equal(t, username, uploader.username)
	assert.Equal(t, password, uploader.password)
}

func TestNewRawExporterWithEnv(t *testing.T) {
	const (
		collectorEndpoint = "http://localhost"
		username          = "user"
		password          = "password"
	)

	envStore, err := ottest.SetEnvVariables(map[string]string{
		envEndpoint: collectorEndpoint,
		envUser:     username,
		envPassword: password,
	})
	require.NoError(t, err)
	defer func() {
		require.NoError(t, envStore.Restore())
	}()

	// Create Jaeger Exporter with environment variables
	exp, err := New(
		WithCollectorEndpoint(),
	)

	assert.NoError(t, err)

	require.IsType(t, &collectorUploader{}, exp.uploader)
	uploader := exp.uploader.(*collectorUploader)
	assert.Equal(t, collectorEndpoint, uploader.endpoint)
	assert.Equal(t, username, uploader.username)
	assert.Equal(t, password, uploader.password)
}

func TestNewRawExporterWithPassedOption(t *testing.T) {
	const (
		collectorEndpoint = "http://localhost"
		username          = "user"
		password          = "password"
		optionEndpoint    = "should not be overwritten"
	)

	envStore, err := ottest.SetEnvVariables(map[string]string{
		envEndpoint: collectorEndpoint,
		envUser:     username,
		envPassword: password,
	})
	require.NoError(t, err)
	defer func() {
		require.NoError(t, envStore.Restore())
	}()

	// Create Jaeger Exporter with passed endpoint option, should be used over envEndpoint
	exp, err := New(
		WithCollectorEndpoint(WithEndpoint(optionEndpoint)),
	)

	assert.NoError(t, err)

	require.IsType(t, &collectorUploader{}, exp.uploader)
	uploader := exp.uploader.(*collectorUploader)
	assert.Equal(t, optionEndpoint, uploader.endpoint)
	assert.Equal(t, username, uploader.username)
	assert.Equal(t, password, uploader.password)
}

func TestEnvOrWithAgentHostPortFromEnv(t *testing.T) {
	testCases := []struct {
		name         string
		envAgentHost string
		envAgentPort string
		defaultHost  string
		defaultPort  string
		expectedHost string
		expectedPort string
	}{
		{
			name:         "overrides default host/port values via environment variables",
			envAgentHost: "localhost",
			envAgentPort: "6832",
			defaultHost:  "hostNameToBeReplaced",
			defaultPort:  "8203",
			expectedHost: "localhost",
			expectedPort: "6832",
		},
		{
			name:         "envAgentHost is empty, will not overwrite default host value",
			envAgentHost: "",
			envAgentPort: "6832",
			defaultHost:  "hostNameNotToBeReplaced",
			defaultPort:  "8203",
			expectedHost: "hostNameNotToBeReplaced",
			expectedPort: "6832",
		},
		{
			name:         "envAgentPort is empty, will not overwrite default port value",
			envAgentHost: "localhost",
			envAgentPort: "",
			defaultHost:  "hostNameToBeReplaced",
			defaultPort:  "8203",
			expectedHost: "localhost",
			expectedPort: "8203",
		},
		{
			name:         "envAgentHost and envAgentPort are empty, will not overwrite default host/port values",
			envAgentHost: "",
			envAgentPort: "",
			defaultHost:  "hostNameNotToBeReplaced",
			defaultPort:  "8203",
			expectedHost: "hostNameNotToBeReplaced",
			expectedPort: "8203",
		},
	}

	envStore := ottest.NewEnvStore()
	envStore.Record(envAgentHost)
	envStore.Record(envAgentPort)
	defer func() {
		require.NoError(t, envStore.Restore())
	}()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, os.Setenv(envAgentHost, tc.envAgentHost))
			require.NoError(t, os.Setenv(envAgentPort, tc.envAgentPort))
			host := envOr(envAgentHost, tc.defaultHost)
			port := envOr(envAgentPort, tc.defaultPort)
			assert.Equal(t, tc.expectedHost, host)
			assert.Equal(t, tc.expectedPort, port)
		})
	}
}

func TestEnvOrWithCollectorEndpointOptionsFromEnv(t *testing.T) {
	testCases := []struct {
		name                             string
		envEndpoint                      string
		envUsername                      string
		envPassword                      string
		defaultCollectorEndpointOptions  collectorEndpointConfig
		expectedCollectorEndpointOptions collectorEndpointConfig
	}{
		{
			name:        "overrides value via environment variables",
			envEndpoint: "http://localhost:14252",
			envUsername: "username",
			envPassword: "password",
			defaultCollectorEndpointOptions: collectorEndpointConfig{
				endpoint: "endpoint not to be used",
				username: "foo",
				password: "bar",
			},
			expectedCollectorEndpointOptions: collectorEndpointConfig{
				endpoint: "http://localhost:14252",
				username: "username",
				password: "password",
			},
		},
		{
			name:        "environment variables is empty, will not overwrite value",
			envEndpoint: "",
			envUsername: "",
			envPassword: "",
			defaultCollectorEndpointOptions: collectorEndpointConfig{
				endpoint: "endpoint to be used",
				username: "foo",
				password: "bar",
			},
			expectedCollectorEndpointOptions: collectorEndpointConfig{
				endpoint: "endpoint to be used",
				username: "foo",
				password: "bar",
			},
		},
	}

	envStore := ottest.NewEnvStore()
	envStore.Record(envEndpoint)
	envStore.Record(envUser)
	envStore.Record(envPassword)
	defer func() {
		require.NoError(t, envStore.Restore())
	}()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, os.Setenv(envEndpoint, tc.envEndpoint))
			require.NoError(t, os.Setenv(envUser, tc.envUsername))
			require.NoError(t, os.Setenv(envPassword, tc.envPassword))

			endpoint := envOr(envEndpoint, tc.defaultCollectorEndpointOptions.endpoint)
			username := envOr(envUser, tc.defaultCollectorEndpointOptions.username)
			password := envOr(envPassword, tc.defaultCollectorEndpointOptions.password)

			assert.Equal(t, tc.expectedCollectorEndpointOptions.endpoint, endpoint)
			assert.Equal(t, tc.expectedCollectorEndpointOptions.username, username)
			assert.Equal(t, tc.expectedCollectorEndpointOptions.password, password)
		})
	}
}
