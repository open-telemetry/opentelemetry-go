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
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/sdk/resource"
)

func mockRuntimeProviders() {
	resource.SetRuntimeProviders(
		fakeRuntimeNameProvider,
		fakeRuntimeVersionProvider,
		func() string { return "LINUX" },
		fakeRuntimeArchProvider,
	)
}

func TestWithOSType(t *testing.T) {
	mockRuntimeProviders()

	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithoutBuiltin(),
		resource.WithOSType(),
	)

	require.NoError(t, err)
	require.EqualValues(t, map[string]string{
		"os.type": "linux",
	}, toMap(res))

	restoreProcessAttributesProviders()
}

func TestWithOSDescription(t *testing.T) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithoutBuiltin(),
		resource.WithOSDescription(),
	)

	require.NoError(t, err)
	require.EqualValues(t, map[string]string{
		"os.description": osDescription(),
	}, toMap(res))
}

func TestWithOS(t *testing.T) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithoutBuiltin(),
		resource.WithOS(),
	)

	require.NoError(t, err)
	require.EqualValues(t, map[string]string{
		"os.type":        osType(),
		"os.description": osDescription(),
	}, toMap(res))
}

func osType() string {
	return strings.ToUpper(runtime.GOOS)
}

func osDescription() string {
	description, _ := resource.OSDescription()

	return description
}
