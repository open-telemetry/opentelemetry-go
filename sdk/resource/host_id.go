// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resource // import "go.opentelemetry.io/otel/sdk/resource"

import (
	"context"
	"os"
	"os/exec"

	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type hostIDProvider func() (string, error)

var defaultHostIDProvider hostIDProvider = platformHostIDReader.read

var hostID = defaultHostIDProvider

func setDefaultHostIDProvider() {
	setHostIDProvider(defaultHostIDProvider)
}

func setHostIDProvider(hostIDProvider hostIDProvider) {
	hostID = hostIDProvider
}

type hostIDReader interface {
	read() (string, error)
}

type fileReader func(string) (string, error)

// nolint:unused // commandExecutor is used on darwin and BSD, but excluded on linux where the lint job runs
type commandExecutor func(string, ...string) (string, error)

func readFile(filename string) (string, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return "", nil
	}

	return string(b), nil
}

// nolint:unused // execCommand is used on darwin and BSD, but excluded on linux where the lint job runs
func execCommand(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	b, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(b), nil
}

type hostIDDetector struct{}

// Detect returns a *Resource containing the platform specific host id.
func (hostIDDetector) Detect(ctx context.Context) (*Resource, error) {
	hostID, err := hostID()
	if err != nil {
		return nil, err
	}

	return NewWithAttributes(
		semconv.SchemaURL,
		semconv.HostID(hostID),
	), nil
}
