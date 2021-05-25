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
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/sys/windows/registry"

	"go.opentelemetry.io/otel/sdk/resource"
)

func TestPlatformOSDescription(t *testing.T) {
	osDescription, err := resource.PlatformOSDescription()

	require.NoError(t, err)
	require.Regexp(t, `^(\w+\s)+\(\d+\)\s\[Version\s\d+(\.\d+){3}\]$`, osDescription)
}

func TestReadRegistryValues(t *testing.T) {
	k, _ := registry.OpenKey(
		registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)

	defer k.Close()

	t.Run("ReadProductName", testReadProductName(t, k))
	t.Run("ReadDisplayVersion", testReadDisplayVersion(t, k))
	t.Run("ReadReleaseID", testReadReleaseID(t, k))
	t.Run("ReadCurrentMajorVersionNumber", testReadCurrentMajorVersionNumber(t, k))
	t.Run("ReadCurrentMinorVersionNumber", testReadCurrentMinorVersionNumber(t, k))
	t.Run("ReadCurrentBuildNumber", testReadCurrentBuildNumber(t, k))
	t.Run("ReadUBR", testReadUBR(t, k))
}

func testReadProductName(t *testing.T, key registry.Key) func(*testing.T) {
	return func(t *testing.T) {
		value := resource.ReadProductName(key)

		require.NotEmpty(t, value)
	}
}

func testReadDisplayVersion(t *testing.T, key registry.Key) func(*testing.T) {
	return func(t *testing.T) {
		value := resource.ReadDisplayVersion(key)

		require.NotEmpty(t, value)
	}
}

func testReadReleaseID(t *testing.T, key registry.Key) func(*testing.T) {
	return func(t *testing.T) {
		value := resource.ReadReleaseID(key)

		require.NotEmpty(t, value)
	}
}

func testReadCurrentMajorVersionNumber(t *testing.T, key registry.Key) func(*testing.T) {
	return func(t *testing.T) {
		value := resource.ReadCurrentMajorVersionNumber(key)

		require.NotEmpty(t, value)
	}
}

func testReadCurrentMinorVersionNumber(t *testing.T, key registry.Key) func(*testing.T) {
	return func(t *testing.T) {
		value := resource.ReadCurrentMinorVersionNumber(key)

		require.NotEmpty(t, value)
	}
}

func testReadCurrentBuildNumber(t *testing.T, key registry.Key) func(*testing.T) {
	return func(t *testing.T) {
		value := resource.ReadCurrentBuildNumber(key)

		require.NotEmpty(t, value)
	}
}

func testReadUBR(t *testing.T, key registry.Key) func(*testing.T) {
	return func(t *testing.T) {
		value := resource.ReadUBR(key)

		require.NotEmpty(t, value)
	}
}
