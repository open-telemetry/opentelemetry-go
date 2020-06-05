package jaeger

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/kv/value"
)

func Test_parseTags(t *testing.T) {
	require.NoError(t, os.Setenv("existing", "not-default"))

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
		disabled          = "disable"
		tags              = "key=value"
	)

	require.NoError(t, os.Setenv(envEndpoint, collectorEndpoint))
	require.NoError(t, os.Setenv(envUser, username))
	require.NoError(t, os.Setenv(envPassword, password))
	require.NoError(t, os.Setenv(envDisabled, disabled))
	require.NoError(t, os.Setenv(envServiceName, serviceName))
	require.NoError(t, os.Setenv(envTags, tags))
	defer func() {
		require.NoError(t, os.Unsetenv(envEndpoint))
		require.NoError(t, os.Unsetenv(envUser))
		require.NoError(t, os.Unsetenv(envPassword))
		require.NoError(t, os.Unsetenv(envDisabled))
		require.NoError(t, os.Unsetenv(envServiceName))
		require.NoError(t, os.Unsetenv(envTags))
	}()

	// Create Jaeger Exporter with environment variables
	exp, err := NewRawExporter(
		WithCollectorEndpoint(CollectorEndpointFromEnv(), WithCollectorEndpointOptionFromEnv()),
		WithDisabledFromEnv(),
		WithProcessFromEnv(),
	)

	assert.NoError(t, err)
	assert.Equal(t, exp.o.Disabled, false)
	assert.EqualValues(t, serviceName, exp.process.ServiceName)
	assert.Len(t, exp.process.Tags, 1)
}

func TestCollectorEndpointFromEnv(t *testing.T) {
	const (
		collectorEndpoint = "http://localhost"
	)

	require.NoError(t, os.Setenv(envEndpoint, collectorEndpoint))
	defer func() {
		require.NoError(t, os.Unsetenv(envEndpoint))
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, os.Setenv(envUser, tc.envUsername))
			require.NoError(t, os.Setenv(envPassword, tc.envPassword))

			f := WithCollectorEndpointOptionFromEnv()
			f(&tc.collectorEndpointOptions)

			assert.Equal(t, tc.expectedCollectorEndpointOptions, tc.collectorEndpointOptions)
		})
	}

	require.NoError(t, os.Unsetenv(envUser))
	require.NoError(t, os.Unsetenv(envPassword))
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, os.Setenv(envDisabled, tc.env))

			f := WithDisabledFromEnv()
			f(&tc.options)

			assert.Equal(t, tc.expectedOptions, tc.options)
		})
	}

	require.NoError(t, os.Unsetenv(envDisabled))
}

func TestProcessFromEnv(t *testing.T) {
	const (
		serviceName = "test-service"
		tags        = "key=value,key2=123"
	)

	require.NoError(t, os.Setenv(envServiceName, serviceName))
	require.NoError(t, os.Setenv(envTags, tags))
	defer func() {
		require.NoError(t, os.Unsetenv(envServiceName))
		require.NoError(t, os.Unsetenv(envTags))
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, os.Setenv(envServiceName, tc.envServiceName))
			require.NoError(t, os.Setenv(envTags, tc.envTags))

			f := WithProcessFromEnv()
			f(&tc.options)

			assert.Equal(t, tc.expectedOptions, tc.options)
		})
	}

	require.NoError(t, os.Unsetenv(envServiceName))
	require.NoError(t, os.Unsetenv(envTags))
}
