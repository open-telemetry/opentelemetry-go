// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resource_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var (
	kv11 = attribute.String("k1", "v11")
	kv12 = attribute.String("k1", "v12")
	kv21 = attribute.String("k2", "v21")
	kv31 = attribute.String("k3", "v31")
	kv41 = attribute.String("k4", "v41")
	kv42 = attribute.String("k4", "")
)

func TestNewWithAttributes(t *testing.T) {
	cases := []struct {
		name string
		in   []attribute.KeyValue
		want []attribute.KeyValue
	}{
		{
			name: "Key with common key order1",
			in:   []attribute.KeyValue{kv12, kv11, kv21},
			want: []attribute.KeyValue{kv11, kv21},
		},
		{
			name: "Key with common key order2",
			in:   []attribute.KeyValue{kv11, kv12, kv21},
			want: []attribute.KeyValue{kv12, kv21},
		},
		{
			name: "Key with nil",
			in:   nil,
			want: nil,
		},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("case-%s", c.name), func(t *testing.T) {
			res := resource.NewSchemaless(c.in...)
			if diff := cmp.Diff(
				res.Attributes(),
				c.want,
				cmp.AllowUnexported(attribute.Value{})); diff != "" {
				t.Fatalf("unwanted result: diff %+v,", diff)
			}
		})
	}
}

func TestMerge(t *testing.T) {
	cases := []struct {
		name      string
		a, b      *resource.Resource
		want      []attribute.KeyValue
		isErr     bool
		schemaURL string
	}{
		{
			name: "Merge 2 nils",
			a:    nil,
			b:    nil,
			want: nil,
		},
		{
			name: "Merge with no overlap, no nil",
			a:    resource.NewSchemaless(kv11, kv31),
			b:    resource.NewSchemaless(kv21, kv41),
			want: []attribute.KeyValue{kv11, kv21, kv31, kv41},
		},
		{
			name: "Merge with no overlap, no nil, not interleaved",
			a:    resource.NewSchemaless(kv11, kv21),
			b:    resource.NewSchemaless(kv31, kv41),
			want: []attribute.KeyValue{kv11, kv21, kv31, kv41},
		},
		{
			name: "Merge with common key order1",
			a:    resource.NewSchemaless(kv11),
			b:    resource.NewSchemaless(kv12, kv21),
			want: []attribute.KeyValue{kv12, kv21},
		},
		{
			name: "Merge with common key order2",
			a:    resource.NewSchemaless(kv12, kv21),
			b:    resource.NewSchemaless(kv11),
			want: []attribute.KeyValue{kv11, kv21},
		},
		{
			name: "Merge with common key order4",
			a:    resource.NewSchemaless(kv11, kv21, kv41),
			b:    resource.NewSchemaless(kv31, kv41),
			want: []attribute.KeyValue{kv11, kv21, kv31, kv41},
		},
		{
			name: "Merge with no keys",
			a:    resource.NewSchemaless(),
			b:    resource.NewSchemaless(),
			want: nil,
		},
		{
			name: "Merge with first resource no keys",
			a:    resource.NewSchemaless(),
			b:    resource.NewSchemaless(kv21),
			want: []attribute.KeyValue{kv21},
		},
		{
			name: "Merge with second resource no keys",
			a:    resource.NewSchemaless(kv11),
			b:    resource.NewSchemaless(),
			want: []attribute.KeyValue{kv11},
		},
		{
			name: "Merge with first resource nil",
			a:    nil,
			b:    resource.NewSchemaless(kv21),
			want: []attribute.KeyValue{kv21},
		},
		{
			name: "Merge with second resource nil",
			a:    resource.NewSchemaless(kv11),
			b:    nil,
			want: []attribute.KeyValue{kv11},
		},
		{
			name: "Merge with first resource value empty string",
			a:    resource.NewSchemaless(kv42),
			b:    resource.NewSchemaless(kv41),
			want: []attribute.KeyValue{kv41},
		},
		{
			name: "Merge with second resource value empty string",
			a:    resource.NewSchemaless(kv41),
			b:    resource.NewSchemaless(kv42),
			want: []attribute.KeyValue{kv42},
		},
		{
			name:      "Merge with first resource with schema",
			a:         resource.NewWithAttributes("https://opentelemetry.io/schemas/1.4.0", kv41),
			b:         resource.NewSchemaless(kv42),
			want:      []attribute.KeyValue{kv42},
			schemaURL: "https://opentelemetry.io/schemas/1.4.0",
		},
		{
			name:      "Merge with second resource with schema",
			a:         resource.NewSchemaless(kv41),
			b:         resource.NewWithAttributes("https://opentelemetry.io/schemas/1.4.0", kv42),
			want:      []attribute.KeyValue{kv42},
			schemaURL: "https://opentelemetry.io/schemas/1.4.0",
		},
		{
			name:  "Merge with different schemas",
			a:     resource.NewWithAttributes("https://opentelemetry.io/schemas/1.4.0", kv41),
			b:     resource.NewWithAttributes("https://opentelemetry.io/schemas/1.3.0", kv42),
			want:  []attribute.KeyValue{kv42},
			isErr: true,
		},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("case-%s", c.name), func(t *testing.T) {
			res, err := resource.Merge(c.a, c.b)
			if c.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, c.schemaURL, res.SchemaURL())
			if diff := cmp.Diff(
				res.Attributes(),
				c.want,
				cmp.AllowUnexported(attribute.Value{})); diff != "" {
				t.Fatalf("unwanted result: diff %+v,", diff)
			}
		})
	}
}

func TestEmpty(t *testing.T) {
	var res *resource.Resource
	assert.Empty(t, res.SchemaURL())
	assert.Empty(t, res.String())
	assert.Equal(t, []attribute.KeyValue(nil), res.Attributes())

	it := res.Iter()
	assert.Equal(t, 0, it.Len())
	assert.True(t, res.Equal(res))
}

func TestDefault(t *testing.T) {
	res := resource.Default()
	require.False(t, res.Equal(resource.Empty()))
	require.True(t, res.Set().HasValue(semconv.ServiceNameKey))

	serviceName, _ := res.Set().Value(semconv.ServiceNameKey)
	require.True(t, strings.HasPrefix(serviceName.AsString(), "unknown_service:"))
	require.Greaterf(t, len(serviceName.AsString()), len("unknown_service:"),
		"default service.name should include executable name")

	require.Contains(t, res.Attributes(), semconv.TelemetrySDKLanguageGo)
	require.Contains(t, res.Attributes(), semconv.TelemetrySDKVersion(sdk.Version()))
	require.Contains(t, res.Attributes(), semconv.TelemetrySDKName("opentelemetry"))
}

func TestString(t *testing.T) {
	for _, test := range []struct {
		kvs  []attribute.KeyValue
		want string
	}{
		{
			kvs:  nil,
			want: "",
		},
		{
			kvs:  []attribute.KeyValue{},
			want: "",
		},
		{
			kvs:  []attribute.KeyValue{kv11},
			want: "k1=v11",
		},
		{
			kvs:  []attribute.KeyValue{kv11, kv12},
			want: "k1=v12",
		},
		{
			kvs:  []attribute.KeyValue{kv11, kv21},
			want: "k1=v11,k2=v21",
		},
		{
			kvs:  []attribute.KeyValue{kv21, kv11},
			want: "k1=v11,k2=v21",
		},
		{
			kvs:  []attribute.KeyValue{kv11, kv21, kv31},
			want: "k1=v11,k2=v21,k3=v31",
		},
		{
			kvs:  []attribute.KeyValue{kv31, kv11, kv21},
			want: "k1=v11,k2=v21,k3=v31",
		},
		{
			kvs:  []attribute.KeyValue{attribute.String("A", "a"), attribute.String("B", "b")},
			want: "A=a,B=b",
		},
		{
			kvs:  []attribute.KeyValue{attribute.String("A", "a,B=b")},
			want: `A=a\,B\=b`,
		},
		{
			kvs:  []attribute.KeyValue{attribute.String("A", `a,B\=b`)},
			want: `A=a\,B\\\=b`,
		},
		{
			kvs:  []attribute.KeyValue{attribute.String("A=a,B", `b`)},
			want: `A\=a\,B=b`,
		},
		{
			kvs:  []attribute.KeyValue{attribute.String(`A=a\,B`, `b`)},
			want: `A\=a\\\,B=b`,
		},
		{
			kvs:  []attribute.KeyValue{attribute.String("", "invalid")},
			want: "",
		},
		{
			kvs:  []attribute.KeyValue{attribute.String("", "invalid"), attribute.String("B", "b")},
			want: "B=b",
		},
	} {
		if got := resource.NewSchemaless(test.kvs...).String(); got != test.want {
			t.Errorf("Resource(%v).String() = %q, want %q", test.kvs, got, test.want)
		}
	}
}

const envVar = "OTEL_RESOURCE_ATTRIBUTES"

func TestMarshalJSON(t *testing.T) {
	r := resource.NewSchemaless(attribute.Int64("A", 1), attribute.String("C", "D"))
	data, err := json.Marshal(r)
	require.NoError(t, err)
	require.JSONEq(t,
		`[{"Key":"A","Value":{"Type":"INT64","Value":1}},{"Key":"C","Value":{"Type":"STRING","Value":"D"}}]`,
		string(data))
}

func TestNew(t *testing.T) {
	tc := []struct {
		name      string
		envars    string
		detectors []resource.Detector
		options   []resource.Option

		resourceValues map[string]string
		schemaURL      string
		wantErr        error
	}{
		{
			name:           "No Options returns empty resource",
			envars:         "key=value,other=attr",
			options:        nil,
			resourceValues: map[string]string{},
		},
		{
			name:   "Nil Detectors works",
			envars: "key=value,other=attr",
			options: []resource.Option{
				resource.WithDetectors(),
			},
			resourceValues: map[string]string{},
		},
		{
			name:   "Only Host",
			envars: "from=here",
			options: []resource.Option{
				resource.WithHost(),
			},
			resourceValues: map[string]string{
				"host.name": hostname(),
			},
			schemaURL: semconv.SchemaURL,
		},
		{
			name:   "Only Env",
			envars: "key=value,other=attr",
			options: []resource.Option{
				resource.WithFromEnv(),
			},
			resourceValues: map[string]string{
				"key":   "value",
				"other": "attr",
			},
		},
		{
			name:   "Only TelemetrySDK",
			envars: "",
			options: []resource.Option{
				resource.WithTelemetrySDK(),
			},
			resourceValues: map[string]string{
				"telemetry.sdk.name":     "opentelemetry",
				"telemetry.sdk.language": "go",
				"telemetry.sdk.version":  sdk.Version(),
			},
			schemaURL: semconv.SchemaURL,
		},
		{
			name:   "WithAttributes",
			envars: "key=value,other=attr",
			options: []resource.Option{
				resource.WithAttributes(attribute.String("A", "B")),
			},
			resourceValues: map[string]string{
				"A": "B",
			},
		},
		{
			name:   "With schema url",
			envars: "",
			options: []resource.Option{
				resource.WithAttributes(attribute.String("A", "B")),
				resource.WithSchemaURL("https://opentelemetry.io/schemas/1.0.0"),
			},
			resourceValues: map[string]string{
				"A": "B",
			},
			schemaURL: "https://opentelemetry.io/schemas/1.0.0",
		},
		{
			name:   "With conflicting schema urls",
			envars: "",
			options: []resource.Option{
				resource.WithDetectors(
					resource.StringDetector("https://opentelemetry.io/schemas/1.0.0", semconv.HostNameKey, os.Hostname),
				),
				resource.WithSchemaURL("https://opentelemetry.io/schemas/1.1.0"),
			},
			resourceValues: map[string]string{
				string(semconv.HostNameKey): func() (hostname string) {
					hostname, _ = os.Hostname()
					return hostname
				}(),
			},
			schemaURL: "",
			wantErr:   resource.ErrSchemaURLConflict,
		},
		{
			name:   "With conflicting detector schema urls",
			envars: "",
			options: []resource.Option{
				resource.WithDetectors(
					resource.StringDetector("https://opentelemetry.io/schemas/1.0.0", semconv.HostNameKey, os.Hostname),
					resource.StringDetector(
						"https://opentelemetry.io/schemas/1.1.0",
						semconv.HostNameKey,
						func() (string, error) { return "", errors.New("fail") },
					),
				),
				resource.WithSchemaURL("https://opentelemetry.io/schemas/1.2.0"),
			},
			resourceValues: map[string]string{
				string(semconv.HostNameKey): func() (hostname string) {
					hostname, _ = os.Hostname()
					return hostname
				}(),
			},
			schemaURL: "",
			wantErr:   resource.ErrSchemaURLConflict,
		},
	}
	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(envVar, tt.envars)
			ctx := context.Background()
			res, err := resource.New(ctx, tt.options...)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			require.Equal(t, tt.resourceValues, toMap(res))

			// TODO: do we need to ensure that resource is never nil and eliminate the
			// following if?
			if res != nil {
				assert.Equal(t, tt.schemaURL, res.SchemaURL())
			}
		})
	}
}

func TestNewWrappedError(t *testing.T) {
	localErr := errors.New("local error")
	_, err := resource.New(
		context.Background(),
		resource.WithDetectors(
			resource.StringDetector("", "", func() (string, error) {
				return "", localErr
			}),
			resource.StringDetector("", "", func() (string, error) {
				return "", assert.AnError
			}),
		),
	)

	assert.ErrorIs(t, err, localErr)
	assert.ErrorIs(t, err, assert.AnError)
	assert.NotErrorIs(t, err, errors.New("false positive error"))
}

func TestWithHostID(t *testing.T) {
	mockHostIDProvider()
	t.Cleanup(restoreHostIDProvider)

	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithHostID(),
	)

	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"host.id": "f2c668b579780554f70f72a063dc0864",
	}, toMap(res))
}

func TestWithHostIDError(t *testing.T) {
	mockHostIDProviderWithError()
	t.Cleanup(restoreHostIDProvider)

	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithHostID(),
	)

	assert.ErrorIs(t, err, assert.AnError)
	require.Equal(t, map[string]string{}, toMap(res))
}

func TestWithOSType(t *testing.T) {
	mockRuntimeProviders()
	t.Cleanup(restoreAttributesProviders)

	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithOSType(),
	)

	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"os.type": "linux",
	}, toMap(res))
}

func TestWithOSDescription(t *testing.T) {
	mockRuntimeProviders()
	t.Cleanup(restoreAttributesProviders)

	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithOSDescription(),
	)

	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"os.description": "Test",
	}, toMap(res))
}

func TestWithOS(t *testing.T) {
	mockRuntimeProviders()
	t.Cleanup(restoreAttributesProviders)

	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithOS(),
	)

	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"os.type":        "linux",
		"os.description": "Test",
	}, toMap(res))
}

func TestWithProcessPID(t *testing.T) {
	mockProcessAttributesProvidersWithErrors()
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithProcessPID(),
	)

	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"process.pid": fmt.Sprint(fakePID),
	}, toMap(res))
}

func TestWithProcessExecutableName(t *testing.T) {
	mockProcessAttributesProvidersWithErrors()
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithProcessExecutableName(),
	)

	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"process.executable.name": fakeExecutableName,
	}, toMap(res))
}

func TestWithProcessExecutablePath(t *testing.T) {
	mockProcessAttributesProviders()
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithProcessExecutablePath(),
	)

	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"process.executable.path": fakeExecutablePath,
	}, toMap(res))
}

func TestWithProcessCommandArgs(t *testing.T) {
	mockProcessAttributesProvidersWithErrors()
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithProcessCommandArgs(),
	)

	require.NoError(t, err)
	jsonCommandArgs, _ := json.Marshal(fakeCommandArgs)
	require.Equal(t, map[string]string{
		"process.command_args": string(jsonCommandArgs),
	}, toMap(res))
}

func TestWithProcessOwner(t *testing.T) {
	mockProcessAttributesProviders()
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithProcessOwner(),
	)

	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"process.owner": fakeOwner,
	}, toMap(res))
}

func TestWithProcessRuntimeName(t *testing.T) {
	mockProcessAttributesProvidersWithErrors()
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithProcessRuntimeName(),
	)

	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"process.runtime.name": fakeRuntimeName,
	}, toMap(res))
}

func TestWithProcessRuntimeVersion(t *testing.T) {
	mockProcessAttributesProvidersWithErrors()
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithProcessRuntimeVersion(),
	)

	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"process.runtime.version": fakeRuntimeVersion,
	}, toMap(res))
}

func TestWithProcessRuntimeDescription(t *testing.T) {
	mockProcessAttributesProvidersWithErrors()
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithProcessRuntimeDescription(),
	)

	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"process.runtime.description": fakeRuntimeDescription,
	}, toMap(res))
}

func TestWithProcess(t *testing.T) {
	mockProcessAttributesProviders()
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithProcess(),
	)

	require.NoError(t, err)
	jsonCommandArgs, _ := json.Marshal(fakeCommandArgs)
	require.Equal(t, map[string]string{
		"process.pid":                 fmt.Sprint(fakePID),
		"process.executable.name":     fakeExecutableName,
		"process.executable.path":     fakeExecutablePath,
		"process.command_args":        string(jsonCommandArgs),
		"process.owner":               fakeOwner,
		"process.runtime.name":        fakeRuntimeName,
		"process.runtime.version":     fakeRuntimeVersion,
		"process.runtime.description": fakeRuntimeDescription,
	}, toMap(res))
}

func toMap(res *resource.Resource) map[string]string {
	m := map[string]string{}
	for _, attr := range res.Attributes() {
		m[string(attr.Key)] = attr.Value.Emit()
	}
	return m
}

func hostname() string {
	hn, err := os.Hostname()
	if err != nil {
		return fmt.Sprintf("hostname(%s)", err)
	}
	return hn
}

func TestWithContainerID(t *testing.T) {
	t.Cleanup(restoreAttributesProviders)

	fakeContainerID := "fake-container-id"

	testCases := []struct {
		name                string
		containerIDProvider func() (string, error)
		expectedResource    map[string]string
		expectedErr         bool
	}{
		{
			name: "get container id",
			containerIDProvider: func() (string, error) {
				return fakeContainerID, nil
			},
			expectedResource: map[string]string{
				string(semconv.ContainerIDKey): fakeContainerID,
			},
		},
		{
			name: "no container id found",
			containerIDProvider: func() (string, error) {
				return "", nil
			},
			expectedResource: map[string]string{},
		},
		{
			name: "error",
			containerIDProvider: func() (string, error) {
				return "", fmt.Errorf("unable to get container id")
			},
			expectedResource: map[string]string{},
			expectedErr:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resource.SetContainerProviders(tc.containerIDProvider)

			res, err := resource.New(context.Background(),
				resource.WithContainerID(),
			)

			if tc.expectedErr {
				assert.Error(t, err)
			}
			assert.Equal(t, tc.expectedResource, toMap(res))
		})
	}
}

func TestWithContainer(t *testing.T) {
	t.Cleanup(restoreAttributesProviders)

	fakeContainerID := "fake-container-id"
	resource.SetContainerProviders(func() (string, error) {
		return fakeContainerID, nil
	})

	res, err := resource.New(context.Background(),
		resource.WithContainer(),
	)

	assert.NoError(t, err)
	assert.Equal(t, map[string]string{
		string(semconv.ContainerIDKey): fakeContainerID,
	}, toMap(res))
}

func TestResourceConcurrentSafe(t *testing.T) {
	// Creating Resources should also be free of any data races,
	// because Resources are immutable.
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			d := &fakeDetector{}
			_, err := resource.Detect(context.Background(), d)
			assert.NoError(t, err)
		}()
	}
	wg.Wait()
}

type fakeDetector struct{}

func (f fakeDetector) Detect(_ context.Context) (*resource.Resource, error) {
	// A bit pedantic, but resource.NewWithAttributes returns an empty Resource when
	// no attributes specified. We want to make sure that this is concurrent-safe.
	return resource.NewWithAttributes("https://opentelemetry.io/schemas/1.3.0"), nil
}

var _ resource.Detector = &fakeDetector{}
