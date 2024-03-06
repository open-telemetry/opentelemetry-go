// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resource_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	k, err := registry.OpenKey(
		registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)

	require.NoError(t, err, "should open Windows CurrentVersion registry key")

	defer k.Close()

	assert.NotEmpty(t, resource.ReadProductName(k), "should read ProductName")
	assert.NotEmpty(t, resource.ReadReleaseID(k), "should read ReleaseID")
	assert.NotEmpty(t, resource.ReadCurrentMajorVersionNumber(k), "should read CurrentMajorVersionNumber")
	assert.NotEmpty(t, resource.ReadCurrentMinorVersionNumber(k), "should read CurrentMinorVersionNumber")
	assert.NotEmpty(t, resource.ReadCurrentBuildNumber(k), "should read CurrentBuildNumber")
	assert.NotEmpty(t, resource.ReadUBR(k), "should read UBR")
	assert.NotPanics(t, func() { resource.ReadDisplayVersion(k) }, "should not panic when reading DisplayVersion")
}
