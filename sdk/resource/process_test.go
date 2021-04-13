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
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/sdk/resource"
)

var (
	fakePID                = 123
	fakeExecutableName     = "mock"
	fakeExecutablePath     = "/fake/path"
	fakeCommandArgs        = []string{"-t", "30"}
	fakeOwner              = "gopher"
	fakeRuntimeName        = "gcmock"
	fakeRuntimeVersion     = "go1.2.3"
	fakeRuntimeDescription = "go version go1.2.3 mock/mock"
)

func mockProcessAttributesProviders() {
	resource.SetProcessAttributesProviders(
		func() int { return fakePID },
		func() string { return fakeExecutableName },
		func() (string, error) { return fakeExecutablePath, nil },
		func() []string { return fakeCommandArgs },
		func() (string, error) { return fakeOwner, nil },
		func() string { return fakeRuntimeName },
		func() string { return fakeRuntimeVersion },
		func() string { return fakeRuntimeDescription },
	)
}

func TestWithProcessFuncs(t *testing.T) {
	mockProcessAttributesProviders()

	t.Run("WithPID", testWithProcessPID)
	t.Run("WithExecutableName", testWithProcessExecutableName)
	t.Run("WithExecutablePath", testWithProcessExecutablePath)
	t.Run("WithCommandArgs", testWithProcessCommandArgs)
	t.Run("WithOwner", testWithProcessOwner)
	t.Run("WithRuntimeName", testWithProcessRuntimeName)
	t.Run("WithRuntimeVersion", testWithProcessRuntimeVersion)
	t.Run("WithRuntimeDescription", testWithProcessRuntimeDescription)
	t.Run("WithProcess", testWithProcess)

	// TODO: tear down (revert to original functions?)
}

func testWithProcessPID(t *testing.T) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithoutBuiltin(),
		resource.WithProcessPID(),
	)

	require.NoError(t, err)
	require.EqualValues(t, map[string]string{
		"process.pid": fmt.Sprint(fakePID),
	}, toMap(res))
}

func testWithProcessExecutableName(t *testing.T) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithoutBuiltin(),
		resource.WithProcessExecutableName(),
	)

	require.NoError(t, err)
	require.EqualValues(t, map[string]string{
		"process.executable.name": fakeExecutableName,
	}, toMap(res))
}

func testWithProcessExecutablePath(t *testing.T) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithoutBuiltin(),
		resource.WithProcessExecutablePath(),
	)

	require.NoError(t, err)
	require.EqualValues(t, map[string]string{
		"process.executable.path": fakeExecutablePath,
	}, toMap(res))
}

func testWithProcessCommandArgs(t *testing.T) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithoutBuiltin(),
		resource.WithProcessCommandArgs(),
	)

	require.NoError(t, err)
	require.EqualValues(t, map[string]string{
		"process.command_args": fmt.Sprint(fakeCommandArgs),
	}, toMap(res))
}

func testWithProcessOwner(t *testing.T) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithoutBuiltin(),
		resource.WithProcessOwner(),
	)

	require.NoError(t, err)
	require.EqualValues(t, map[string]string{
		"process.owner": fakeOwner,
	}, toMap(res))
}

func testWithProcessRuntimeName(t *testing.T) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithoutBuiltin(),
		resource.WithProcessRuntimeName(),
	)

	require.NoError(t, err)
	require.EqualValues(t, map[string]string{
		"process.runtime.name": fakeRuntimeName,
	}, toMap(res))
}

func testWithProcessRuntimeVersion(t *testing.T) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithoutBuiltin(),
		resource.WithProcessRuntimeVersion(),
	)

	require.NoError(t, err)
	require.EqualValues(t, map[string]string{
		"process.runtime.version": fakeRuntimeVersion,
	}, toMap(res))
}

func testWithProcessRuntimeDescription(t *testing.T) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithoutBuiltin(),
		resource.WithProcessRuntimeDescription(),
	)

	require.NoError(t, err)
	require.EqualValues(t, map[string]string{
		"process.runtime.description": fakeRuntimeDescription,
	}, toMap(res))
}

func testWithProcess(t *testing.T) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithoutBuiltin(),
		resource.WithProcess(),
	)

	require.NoError(t, err)
	require.EqualValues(t, map[string]string{
		"process.pid":                 fmt.Sprint(fakePID),
		"process.executable.name":     fakeExecutableName,
		"process.executable.path":     fakeExecutablePath,
		"process.command_args":        fmt.Sprint(fakeCommandArgs),
		"process.owner":               fakeOwner,
		"process.runtime.name":        fakeRuntimeName,
		"process.runtime.version":     fakeRuntimeVersion,
		"process.runtime.description": fakeRuntimeDescription,
	}, toMap(res))
}
