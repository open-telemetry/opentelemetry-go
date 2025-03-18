// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resource_test

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/sdk/resource"
)

var (
	fakePID            = 123
	fakeExecutablePath = "/fake/path/mock"
	fakeCommandArgs    = []string{"mock", "-t", "30"}
	fakeOwner          = "gopher"
	fakeRuntimeName    = "gcmock"
	fakeRuntimeVersion = "go1.2.3"
	fakeRuntimeOS      = "linux"
	fakeRuntimeArch    = "amd64"
)

var (
	fakeExecutableName     = "mock"
	fakeRuntimeDescription = "go version go1.2.3 linux/amd64"
)

var (
	fakePidProvider            = func() int { return fakePID }
	fakeExecutablePathProvider = func() (string, error) { return fakeExecutablePath, nil }
	fakeCommandArgsProvider    = func() []string { return fakeCommandArgs }
	fakeOwnerProvider          = func() (*user.User, error) { return &user.User{Username: fakeOwner}, nil }
	fakeRuntimeNameProvider    = func() string { return fakeRuntimeName }
	fakeRuntimeVersionProvider = func() string { return fakeRuntimeVersion }
	fakeRuntimeOSProvider      = func() string { return fakeRuntimeOS }
	fakeRuntimeArchProvider    = func() string { return fakeRuntimeArch }
)

var (
	fakeExecutablePathProviderWithError = func() (string, error) {
		return "", fmt.Errorf("unable to get process executable")
	}
	fakeOwnerProviderWithError = func() (*user.User, error) {
		return nil, fmt.Errorf("unable to get process user")
	}
)

func mockProcessAttributesProviders() {
	resource.SetOSProviders(
		fakePidProvider,
		fakeExecutablePathProvider,
		fakeCommandArgsProvider,
	)
	resource.SetRuntimeProviders(
		fakeRuntimeNameProvider,
		fakeRuntimeVersionProvider,
		fakeRuntimeOSProvider,
		fakeRuntimeArchProvider,
	)
	resource.SetUserProviders(
		fakeOwnerProvider,
	)
}

func mockProcessAttributesProvidersWithErrors() {
	resource.SetOSProviders(
		fakePidProvider,
		fakeExecutablePathProviderWithError,
		fakeCommandArgsProvider,
	)
	resource.SetRuntimeProviders(
		fakeRuntimeNameProvider,
		fakeRuntimeVersionProvider,
		fakeRuntimeOSProvider,
		fakeRuntimeArchProvider,
	)
	resource.SetUserProviders(
		fakeOwnerProviderWithError,
	)
}

func restoreAttributesProviders() {
	resource.SetDefaultOSProviders()
	resource.SetDefaultRuntimeProviders()
	resource.SetDefaultUserProviders()
	resource.SetDefaultOSDescriptionProvider()
	resource.SetDefaultContainerProviders()
}

func TestWithProcessFuncsErrors(t *testing.T) {
	mockProcessAttributesProvidersWithErrors()

	t.Run("WithExecutablePath", testWithProcessExecutablePathError)
	t.Run("WithOwner", testWithProcessOwnerError)

	restoreAttributesProviders()
}

func TestCommandArgs(t *testing.T) {
	require.Equal(t, os.Args, resource.CommandArgs())
}

func TestRuntimeName(t *testing.T) {
	if runtime.Compiler == "gc" {
		require.Equal(t, "go", resource.RuntimeName())
	} else {
		require.Equal(t, runtime.Compiler, resource.RuntimeName())
	}
}

func TestRuntimeOS(t *testing.T) {
	require.Equal(t, runtime.GOOS, resource.RuntimeOS())
}

func TestRuntimeArch(t *testing.T) {
	require.Equal(t, runtime.GOARCH, resource.RuntimeArch())
}

func testWithProcessExecutablePathError(t *testing.T) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithProcessExecutablePath(),
	)

	require.Error(t, err)
	require.Equal(t, map[string]string{}, toMap(res))
}

func testWithProcessOwnerError(t *testing.T) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithProcessOwner(),
	)

	require.Error(t, err)
	require.Equal(t, map[string]string{}, toMap(res))
}
