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

package resource_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	ottest "go.opentelemetry.io/otel/internal/testing"
	"go.opentelemetry.io/otel/label"
	opentelemetry "go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/resource"
)

const envVar = "OTEL_RESOURCE_ATTRIBUTES"

func TestDefaultConfig(t *testing.T) {
	store, err := ottest.SetEnvVariables(map[string]string{
		envVar: "",
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, store.Restore()) }()

	ctx := context.Background()
	res, err := resource.New(ctx)
	require.NoError(t, err)
	require.EqualValues(t, map[string]string{
		"host.name":              hostname(),
		"telemetry.sdk.name":     "opentelemetry-go",
		"telemetry.sdk.language": "go",
		"telemetry.sdk.version":  opentelemetry.Version(),
	}, toMap(res))
}

func TestDefaultConfigNoHost(t *testing.T) {
	store, err := ottest.SetEnvVariables(map[string]string{
		envVar: "",
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, store.Restore()) }()

	ctx := context.Background()
	res, err := resource.New(ctx, resource.WithHost(nil))
	require.NoError(t, err)
	require.EqualValues(t, map[string]string{
		"telemetry.sdk.name":     "opentelemetry-go",
		"telemetry.sdk.language": "go",
		"telemetry.sdk.version":  opentelemetry.Version(),
	}, toMap(res))
}

func TestDefaultConfigNoEnv(t *testing.T) {
	store, err := ottest.SetEnvVariables(map[string]string{
		envVar: "from=here",
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, store.Restore()) }()

	ctx := context.Background()
	res, err := resource.New(ctx, resource.WithFromEnv(nil))
	require.NoError(t, err)
	require.EqualValues(t, map[string]string{
		"host.name":              hostname(),
		"telemetry.sdk.name":     "opentelemetry-go",
		"telemetry.sdk.language": "go",
		"telemetry.sdk.version":  opentelemetry.Version(),
	}, toMap(res))
}

func TestDefaultConfigWithEnv(t *testing.T) {
	store, err := ottest.SetEnvVariables(map[string]string{
		envVar: "key=value,other=attr",
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, store.Restore()) }()

	ctx := context.Background()
	res, err := resource.New(ctx)
	require.NoError(t, err)
	require.EqualValues(t, map[string]string{
		"key":                    "value",
		"other":                  "attr",
		"host.name":              hostname(),
		"telemetry.sdk.name":     "opentelemetry-go",
		"telemetry.sdk.language": "go",
		"telemetry.sdk.version":  opentelemetry.Version(),
	}, toMap(res))
}

func TestWithoutBuiltin(t *testing.T) {
	store, err := ottest.SetEnvVariables(map[string]string{
		envVar: "key=value,other=attr",
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, store.Restore()) }()

	ctx := context.Background()
	res, err := resource.New(
		ctx,
		resource.WithoutBuiltin(),
		resource.WithAttributes(label.String("hello", "collector")),
	)
	require.NoError(t, err)
	require.EqualValues(t, map[string]string{
		"hello": "collector",
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
