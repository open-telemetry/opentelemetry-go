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
	"os/user"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/sdk/resource"
)

func TestWithProcessPID(t *testing.T) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithoutBuiltin(),
		resource.WithProcessPID(),
	)

	require.NoError(t, err)
	require.EqualValues(t, map[string]string{
		"process.pid": fmt.Sprint(pid()),
	}, toMap(res))
}

func TestWithProcessExecutableName(t *testing.T) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithoutBuiltin(),
		resource.WithProcessExecutableName(),
	)

	require.NoError(t, err)
	require.EqualValues(t, map[string]string{
		"process.executable.name": executableName(),
	}, toMap(res))
}

func TestWithProcessExecutablePath(t *testing.T) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithoutBuiltin(),
		resource.WithProcessExecutablePath(),
	)

	require.NoError(t, err)
	require.EqualValues(t, map[string]string{
		"process.executable.path": executablePath(),
	}, toMap(res))
}

func TestWithProcessCommandArgs(t *testing.T) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithoutBuiltin(),
		resource.WithProcessCommandArgs(),
	)

	require.NoError(t, err)
	require.EqualValues(t, map[string]string{
		"process.command_args": fmt.Sprint(commandArgs()),
	}, toMap(res))
}

func TestWithProcessOwner(t *testing.T) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithoutBuiltin(),
		resource.WithProcessOwner(),
	)

	require.NoError(t, err)
	require.EqualValues(t, map[string]string{
		"process.owner": owner(),
	}, toMap(res))
}

func TestWithProcessRuntimeName(t *testing.T) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithoutBuiltin(),
		resource.WithProcessRuntimeName(),
	)

	require.NoError(t, err)
	require.EqualValues(t, map[string]string{
		"process.runtime.name": runtimeName(),
	}, toMap(res))
}

func TestWithProcessRuntimeVersion(t *testing.T) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithoutBuiltin(),
		resource.WithProcessRuntimeVersion(),
	)

	require.NoError(t, err)
	require.EqualValues(t, map[string]string{
		"process.runtime.version": runtimeVersion(),
	}, toMap(res))
}

func TestWithProcessRuntimeDescription(t *testing.T) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithoutBuiltin(),
		resource.WithProcessRuntimeDescription(),
	)

	require.NoError(t, err)
	require.EqualValues(t, map[string]string{
		"process.runtime.description": runtimeDescription(),
	}, toMap(res))
}

func TestWithProcess(t *testing.T) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithoutBuiltin(),
		resource.WithProcess(),
	)

	require.NoError(t, err)
	require.EqualValues(t, map[string]string{
		"process.pid":                 fmt.Sprint(pid()),
		"process.executable.name":     executableName(),
		"process.executable.path":     executablePath(),
		"process.command_args":        fmt.Sprint(commandArgs()),
		"process.owner":               owner(),
		"process.runtime.name":        runtimeName(),
		"process.runtime.version":     runtimeVersion(),
		"process.runtime.description": runtimeDescription(),
	}, toMap(res))
}

func pid() int {
	return os.Getpid()
}

func executableName() string {
	return filepath.Base(os.Args[0])
}

func executablePath() string {
	executable, _ := os.Executable()

	return executable
}

func commandArgs() []string {
	return os.Args
}

func owner() string {
	user, _ := user.Current()

	return user.Username
}

func runtimeName() string {
	return runtime.Compiler
}

func runtimeVersion() string {
	return runtime.Version()
}

func runtimeDescription() string {
	return fmt.Sprintf("go version %s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
