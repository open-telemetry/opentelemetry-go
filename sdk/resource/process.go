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

package resource // import "go.opentelemetry.io/otel/sdk/resource"

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"go.opentelemetry.io/otel/semconv"
)

type pidProvider func() int
type executableNameProvider func() string
type executablePathProvider func() (string, error)
type commandArgsProvider func() []string
type ownerProvider func() (string, error)
type runtimeNameProvider func() string
type runtimeVersionProvider func() string
type runtimeDescriptionProvider func() string

var (
	pid            pidProvider            = os.Getpid
	executableName executableNameProvider = func() string { return filepath.Base(os.Args[0]) }
	executablePath executablePathProvider = os.Executable
	commandArgs    commandArgsProvider    = func() []string { return os.Args }
	owner          ownerProvider          = func() (string, error) {
		user, err := user.Current()
		if err != nil {
			return "", err
		}

		return user.Username, nil
	}
	runtimeName        runtimeNameProvider        = func() string { return runtime.Compiler }
	runtimeVersion     runtimeVersionProvider     = runtime.Version
	runtimeDescription runtimeDescriptionProvider = func() string {
		return fmt.Sprintf("go version %s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	}
)

func setProcessAttributesProviders(
	pidProvider pidProvider,
	executableNameProvider executableNameProvider,
	executablePathProvider executablePathProvider,
	commandArgsProvider commandArgsProvider,
	ownerProvider ownerProvider,
	runtimeNameProvider runtimeNameProvider,
	runtimeVersionProvider runtimeVersionProvider,
	runtimeDescriptionProvider runtimeDescriptionProvider,
) {
	pid = pidProvider
	executableName = executableNameProvider
	executablePath = executablePathProvider
	commandArgs = commandArgsProvider
	owner = ownerProvider
	runtimeName = runtimeNameProvider
	runtimeVersion = runtimeVersionProvider
	runtimeDescription = runtimeDescriptionProvider
}

type processPIDDetector struct{}
type processExecutableNameDetector struct{}
type processExecutablePathDetector struct{}
type processCommandArgsDetector struct{}
type processOwnerDetector struct{}
type processRuntimeNameDetector struct{}
type processRuntimeVersionDetector struct{}
type processRuntimeDescriptionDetector struct{}

// Detect returns a *Resource that describes the process identifier (PID) of the
// executing process.
func (processPIDDetector) Detect(ctx context.Context) (*Resource, error) {
	return NewWithAttributes(semconv.ProcessPIDKey.Int(pid())), nil
}

// Detect returns a *Resource that describes the name of the process executable.
func (processExecutableNameDetector) Detect(ctx context.Context) (*Resource, error) {
	return NewWithAttributes(semconv.ProcessExecutableNameKey.String(executableName())), nil
}

// Detect returns a *Resource that describes the full path of the process executable.
func (processExecutablePathDetector) Detect(ctx context.Context) (*Resource, error) {
	executablePath, err := executablePath()
	if err != nil {
		return nil, err
	}

	return NewWithAttributes(semconv.ProcessExecutablePathKey.String(executablePath)), nil
}

// Detect returns a *Resource that describes all the command arguments as received
// by the process.
func (processCommandArgsDetector) Detect(ctx context.Context) (*Resource, error) {
	return NewWithAttributes(semconv.ProcessCommandArgsKey.Array(commandArgs())), nil
}

// Detect returns a *Resource that describes the username of the user that owns the
// process.
func (processOwnerDetector) Detect(ctx context.Context) (*Resource, error) {
	owner, err := owner()
	if err != nil {
		return nil, err
	}

	return NewWithAttributes(semconv.ProcessOwnerKey.String(owner)), nil
}

// Detect returns a *Resource that describes the name of the compiler used to compile
// this process image.
func (processRuntimeNameDetector) Detect(ctx context.Context) (*Resource, error) {
	return NewWithAttributes(semconv.ProcessRuntimeNameKey.String(runtimeName())), nil
}

// Detect returns a *Resource that describes the version of the runtime of this process.
func (processRuntimeVersionDetector) Detect(ctx context.Context) (*Resource, error) {
	return NewWithAttributes(semconv.ProcessRuntimeVersionKey.String(runtimeVersion())), nil
}

// Detect returns a *Resource that describes the runtime of this process.
func (processRuntimeDescriptionDetector) Detect(ctx context.Context) (*Resource, error) {
	return NewWithAttributes(
		semconv.ProcessRuntimeDescriptionKey.String(runtimeDescription()),
	), nil
}

// WithProcessPID adds an attribute with the process identifier (PID) to the
// configured Resource.
func WithProcessPID() Option {
	return WithDetectors(processPIDDetector{})
}

// WithProcessExecutableName adds an attribute with the name of the process
// executable to the configured Resource.
func WithProcessExecutableName() Option {
	return WithDetectors(processExecutableNameDetector{})
}

// WithProcessExecutablePath adds an attribute with the full path to the process
// executable to the configured Resource.
func WithProcessExecutablePath() Option {
	return WithDetectors(processExecutablePathDetector{})
}

// WithProcessCommandArgs adds an attribute with all the command arguments (including
// the command/executable itself) as received by the process the configured Resource.
func WithProcessCommandArgs() Option {
	return WithDetectors(processCommandArgsDetector{})
}

// WithProcessOwner adds an attribute the username of the user that owns the process
// to the configured Resource.
func WithProcessOwner() Option {
	return WithDetectors(processOwnerDetector{})
}

// WithProcessRuntimeName adds an attribute with the name of the runtime of this
// process to the configured Resource.
func WithProcessRuntimeName() Option {
	return WithDetectors(processRuntimeNameDetector{})
}

// WithProcessRuntimeVersion adds an attribute with the version of the runtime of
// this process to the configured Resource.
func WithProcessRuntimeVersion() Option {
	return WithDetectors(processRuntimeVersionDetector{})
}

// WithProcessRuntimeDescription adds an attribute with an additional description
// about the runtime of the process to the configured Resource.
func WithProcessRuntimeDescription() Option {
	return WithDetectors(processRuntimeDescriptionDetector{})
}

// WithProcess adds all the Process attributes to the configured Resource.
// See individual WithProcess* functions to configure specific attributes.
func WithProcess() Option {
	return WithDetectors(
		processPIDDetector{},
		processExecutableNameDetector{},
		processExecutablePathDetector{},
		processCommandArgsDetector{},
		processOwnerDetector{},
		processRuntimeNameDetector{},
		processRuntimeVersionDetector{},
		processRuntimeDescriptionDetector{},
	)
}
