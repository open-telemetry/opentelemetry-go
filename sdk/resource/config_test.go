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

	"go.opentelemetry.io/otel"
	ottest "go.opentelemetry.io/otel/internal/internaltest"
	"go.opentelemetry.io/otel/sdk/resource"
)

const envVar = "OTEL_RESOURCE_ATTRIBUTES"

func TestConfig(t *testing.T) {
	tc := []struct {
		name      string
		envars    string
		detectors []resource.Detector

		resouceValues map[string]string
	}{
		{
			name:          "No detectors disables detection",
			envars:        "key=value,other=attr",
			detectors:     []resource.Detector{},
			resouceValues: map[string]string{},
		},
		{
			name:          "Nil detectors disables detection",
			envars:        "key=value,other=attr",
			detectors:     nil,
			resouceValues: map[string]string{},
		},
		{
			name:   "Only Host",
			envars: "from=here",
			detectors: []resource.Detector{
				resource.Host{},
			},
			resouceValues: map[string]string{
				"host.name": hostname(),
			},
		},
		{
			name:   "Only Env",
			envars: "key=value,other=attr",
			detectors: []resource.Detector{
				resource.FromEnv{},
			},
			resouceValues: map[string]string{
				"key":   "value",
				"other": "attr",
			},
		},
		{
			name:   "Only TelemetrySDK",
			envars: "",
			detectors: []resource.Detector{
				resource.TelemetrySDK{},
			},
			resouceValues: map[string]string{
				"telemetry.sdk.name":     "opentelemetry",
				"telemetry.sdk.language": "go",
				"telemetry.sdk.version":  otel.Version(),
			},
		},
		{
			name:      "Builtins",
			envars:    "key=value,other=attr",
			detectors: resource.BuiltinDetectors,
			resouceValues: map[string]string{
				"host.name":              hostname(),
				"telemetry.sdk.name":     "opentelemetry",
				"telemetry.sdk.language": "go",
				"telemetry.sdk.version":  otel.Version(),
				"key":                    "value",
				"other":                  "attr",
			},
		},
	}
	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			store, err := ottest.SetEnvVariables(map[string]string{
				envVar: tt.envars,
			})
			require.NoError(t, err)
			defer func() { require.NoError(t, store.Restore()) }()

			ctx := context.Background()
			res, err := resource.New(ctx,
				resource.WithDetectors(tt.detectors...),
			)

			require.NoError(t, err)
			require.EqualValues(t, tt.resouceValues, toMap(res))
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	resouceValues := map[string]string{
		"host.name":              hostname(),
		"telemetry.sdk.name":     "opentelemetry",
		"telemetry.sdk.language": "go",
		"telemetry.sdk.version":  otel.Version(),
	}
	store, err := ottest.SetEnvVariables(map[string]string{
		envVar: "",
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, store.Restore()) }()

	ctx := context.Background()

	res, err := resource.New(ctx)

	require.NoError(t, err)
	require.EqualValues(t, resouceValues, toMap(res))
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
