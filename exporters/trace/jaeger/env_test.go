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

func TestAgentEndpointFromEnv(t *testing.T) {
	testCases := []struct {
		name             string
		envAgentHost     string
		envAgentPort     string
		hostPort         string
		expectedHostPort string
	}{
		{
			name:             "overrides HostPort value via environment variables",
			envAgentHost:     "localhost",
			envAgentPort:     "6832",
			hostPort:         "hostNameToBeReplaced:8203",
			expectedHostPort: "localhost:6832",
		},
		{
			name:             "envAgentHost is empty, will not overwrite HostPort value",
			envAgentHost:     "",
			envAgentPort:     "6832",
			hostPort:         "hostNameNotToBeReplaced:8203",
			expectedHostPort: "hostNameNotToBeReplaced:8203",
		},
		{
			name:             "envAgentPort is empty, will not overwrite HostPort value",
			envAgentHost:     "localhost",
			envAgentPort:     "",
			hostPort:         "hostNameNotToBeReplaced:8203",
			expectedHostPort: "hostNameNotToBeReplaced:8203",
		},
		{
			name:             "envAgentHost and envAgentPort are empty, will not overwrite HostPort value",
			envAgentHost:     "",
			envAgentPort:     "",
			hostPort:         "hostNameNotToBeReplaced:8203",
			expectedHostPort: "hostNameNotToBeReplaced:8203",
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
			AgentEndpointFromEnv(&tc.hostPort)
			assert.Equal(t, tc.hostPort, tc.expectedHostPort)
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
