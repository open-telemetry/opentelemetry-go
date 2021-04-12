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
	exp, err := NewRawExporter(
		WithCollectorEndpoint(CollectorEndpointFromEnv(), WithCollectorEndpointOptionFromEnv()),
	)

	assert.NoError(t, err)

	require.IsType(t, &collectorUploader{}, exp.uploader)
	uploader := exp.uploader.(*collectorUploader)
	assert.Equal(t, collectorEndpoint, uploader.endpoint)
	assert.Equal(t, username, uploader.username)
	assert.Equal(t, password, uploader.password)
}

func TestNewRawExporterWithEnvImplicitly(t *testing.T) {
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
	exp, err := NewRawExporter(
		WithCollectorEndpoint("should be overwritten"),
	)

	assert.NoError(t, err)

	require.IsType(t, &collectorUploader{}, exp.uploader)
	uploader := exp.uploader.(*collectorUploader)
	assert.Equal(t, collectorEndpoint, uploader.endpoint)
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

func TestCollectorEndpointFromEnv(t *testing.T) {
	const (
		collectorEndpoint = "http://localhost"
	)

	envStore, err := ottest.SetEnvVariables(map[string]string{
		envEndpoint: collectorEndpoint,
	})
	require.NoError(t, err)
	defer func() {
		require.NoError(t, envStore.Restore())
	}()

	assert.Equal(t, collectorEndpoint, CollectorEndpointFromEnv())
}

func TestWithCollectorEndpointOptionFromEnv(t *testing.T) {
	testCases := []struct {
		name                             string
		envUsername                      string
		envPassword                      string
		collectorEndpointOptions         CollectorEndpointOptions
		expectedCollectorEndpointOptions CollectorEndpointOptions
	}{
		{
			name:        "overrides value via environment variables",
			envUsername: "username",
			envPassword: "password",
			collectorEndpointOptions: CollectorEndpointOptions{
				username: "foo",
				password: "bar",
			},
			expectedCollectorEndpointOptions: CollectorEndpointOptions{
				username: "username",
				password: "password",
			},
		},
		{
			name:        "environment variables is empty, will not overwrite value",
			envUsername: "",
			envPassword: "",
			collectorEndpointOptions: CollectorEndpointOptions{
				username: "foo",
				password: "bar",
			},
			expectedCollectorEndpointOptions: CollectorEndpointOptions{
				username: "foo",
				password: "bar",
			},
		},
	}

	envStore := ottest.NewEnvStore()
	envStore.Record(envUser)
	envStore.Record(envPassword)
	defer func() {
		require.NoError(t, envStore.Restore())
	}()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, os.Setenv(envUser, tc.envUsername))
			require.NoError(t, os.Setenv(envPassword, tc.envPassword))

			f := WithCollectorEndpointOptionFromEnv()
			f(&tc.collectorEndpointOptions)

			assert.Equal(t, tc.expectedCollectorEndpointOptions, tc.collectorEndpointOptions)
		})
	}
}
