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

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/kv/value"
	ottest "go.opentelemetry.io/otel/internal/testing"
)

func Test_parseTags(t *testing.T) {
	envStore, err := ottest.SetEnvVariables(map[string]string{
		"existing": "not-default",
	})
	require.NoError(t, err)
	defer func() {
		require.NoError(t, envStore.Restore())
	}()

	tags := "key=value,k1=${nonExisting:default}, k2=${withSpace:default},k3=${existing:default},k4=true,k5=42,k6=-1.2"
	ts := parseTags(tags)
	assert.Equal(t, 7, len(ts))

	assert.Equal(t, kv.Key("key"), ts[0].Key)
	assert.Equal(t, value.String("value"), ts[0].Value)

	assert.Equal(t, kv.Key("k1"), ts[1].Key)
	assert.Equal(t, value.String("default"), ts[1].Value)

	assert.Equal(t, kv.Key("k2"), ts[2].Key)
	assert.Equal(t, value.String("default"), ts[2].Value)

	assert.Equal(t, kv.Key("k3"), ts[3].Key)
	assert.Equal(t, value.String("not-default"), ts[3].Value)

	assert.Equal(t, kv.Key("k4"), ts[4].Key)
	assert.Equal(t, value.Bool(true), ts[4].Value)

	assert.Equal(t, kv.Key("k5"), ts[5].Key)
	assert.Equal(t, value.Int64(42), ts[5].Value)

	assert.Equal(t, kv.Key("k6"), ts[6].Key)
	assert.Equal(t, value.Float64(-1.2), ts[6].Value)

	require.NoError(t, os.Unsetenv("existing"))
}

func Test_parseValue(t *testing.T) {
	testCases := []struct {
		name     string
		str      string
		expected value.Value
	}{
		{
			name:     "bool: true",
			str:      "true",
			expected: value.Bool(true),
		},
		{
			name:     "bool: false",
			str:      "false",
			expected: value.Bool(false),
		},
		{
			name:     "int64: 012340",
			str:      "012340",
			expected: value.Int64(12340),
		},
		{
			name:     "int64: -012340",
			str:      "-012340",
			expected: value.Int64(-12340),
		},
		{
			name:     "int64: 0",
			str:      "0",
			expected: value.Int64(0),
		},
		{
			name:     "float64: -0.1",
			str:      "-0.1",
			expected: value.Float64(-0.1),
		},
		{
			name:     "float64: 00.001",
			str:      "00.001",
			expected: value.Float64(0.001),
		},
		{
			name:     "float64: 1E23",
			str:      "1E23",
			expected: value.Float64(1e23),
		},
		{
			name:     "string: foo",
			str:      "foo",
			expected: value.String("foo"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := parseValue(tc.str)
			assert.Equal(t, tc.expected, v)
		})
	}
}

func TestNewRawExporterWithEnv(t *testing.T) {
	const (
		collectorEndpoint = "http://localhost"
		username          = "user"
		password          = "password"
		serviceName       = "test-service"
		disabled          = "false"
		tags              = "key=value"
	)

	envStore, err := ottest.SetEnvVariables(map[string]string{
		envEndpoint:    collectorEndpoint,
		envUser:        username,
		envPassword:    password,
		envDisabled:    disabled,
		envServiceName: serviceName,
		envTags:        tags,
	})
	require.NoError(t, err)
	defer func() {
		require.NoError(t, envStore.Restore())
	}()

	// Create Jaeger Exporter with environment variables
	exp, err := NewRawExporter(
		WithCollectorEndpoint(CollectorEndpointFromEnv(), WithCollectorEndpointOptionFromEnv()),
		WithDisabled(true),
		WithDisabledFromEnv(),
		WithProcessFromEnv(),
	)

	assert.NoError(t, err)
	assert.Equal(t, false, exp.o.Disabled)
	assert.EqualValues(t, serviceName, exp.process.ServiceName)
	assert.Len(t, exp.process.Tags, 1)

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
		serviceName       = "test-service"
		disabled          = "false"
		tags              = "key=value"
	)

	envStore, err := ottest.SetEnvVariables(map[string]string{
		envEndpoint:    collectorEndpoint,
		envUser:        username,
		envPassword:    password,
		envDisabled:    disabled,
		envServiceName: serviceName,
		envTags:        tags,
	})
	require.NoError(t, err)
	defer func() {
		require.NoError(t, envStore.Restore())
	}()

	// Create Jaeger Exporter with environment variables
	exp, err := NewRawExporter(
		WithCollectorEndpoint("should be overwritten"),
		WithDisabled(true),
	)

	assert.NoError(t, err)
	assert.Equal(t, false, exp.o.Disabled)
	assert.EqualValues(t, serviceName, exp.process.ServiceName)
	assert.Len(t, exp.process.Tags, 1)

	require.IsType(t, &collectorUploader{}, exp.uploader)
	uploader := exp.uploader.(*collectorUploader)
	assert.Equal(t, collectorEndpoint, uploader.endpoint)
	assert.Equal(t, username, uploader.username)
	assert.Equal(t, password, uploader.password)
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

func TestWithDisabledFromEnv(t *testing.T) {
	testCases := []struct {
		name            string
		env             string
		options         options
		expectedOptions options
	}{
		{
			name:            "overwriting",
			env:             "true",
			options:         options{},
			expectedOptions: options{Disabled: true},
		},
		{
			name:            "no overwriting",
			env:             "",
			options:         options{Disabled: true},
			expectedOptions: options{Disabled: true},
		},
	}

	envStore := ottest.NewEnvStore()
	envStore.Record(envDisabled)
	defer func() {
		require.NoError(t, envStore.Restore())
	}()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, os.Setenv(envDisabled, tc.env))

			f := WithDisabledFromEnv()
			f(&tc.options)

			assert.Equal(t, tc.expectedOptions, tc.options)
		})
	}
}

func TestProcessFromEnv(t *testing.T) {
	const (
		serviceName = "test-service"
		tags        = "key=value,key2=123"
	)

	envStore, err := ottest.SetEnvVariables(map[string]string{
		envServiceName: serviceName,
		envTags:        tags,
	})
	require.NoError(t, err)
	defer func() {
		require.NoError(t, envStore.Restore())
	}()

	p := ProcessFromEnv()

	assert.Equal(t, Process{
		ServiceName: serviceName,
		Tags: []kv.KeyValue{
			kv.String("key", "value"),
			kv.Int64("key2", 123),
		},
	}, p)
}

func TestWithProcessFromEnv(t *testing.T) {
	testCases := []struct {
		name            string
		envServiceName  string
		envTags         string
		options         options
		expectedOptions options
	}{
		{
			name:           "overwriting",
			envServiceName: "service-name",
			envTags:        "key=value",
			options: options{
				Process: Process{
					ServiceName: "old-name",
					Tags: []kv.KeyValue{
						kv.String("old-key", "old-value"),
					},
				},
			},
			expectedOptions: options{
				Process: Process{
					ServiceName: "service-name",
					Tags: []kv.KeyValue{
						kv.String("key", "value"),
					},
				},
			},
		},
		{
			name:           "no overwriting",
			envServiceName: "",
			envTags:        "",
			options: options{
				Process: Process{
					ServiceName: "old-name",
					Tags: []kv.KeyValue{
						kv.String("old-key", "old-value"),
					},
				},
			},
			expectedOptions: options{
				Process: Process{
					ServiceName: "old-name",
					Tags: []kv.KeyValue{
						kv.String("old-key", "old-value"),
					},
				},
			},
		},
	}

	envStore := ottest.NewEnvStore()
	envStore.Record(envServiceName)
	envStore.Record(envTags)
	defer func() {
		require.NoError(t, envStore.Restore())
	}()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, os.Setenv(envServiceName, tc.envServiceName))
			require.NoError(t, os.Setenv(envTags, tc.envTags))

			f := WithProcessFromEnv()
			f(&tc.options)

			assert.Equal(t, tc.expectedOptions, tc.options)
		})
	}
}
