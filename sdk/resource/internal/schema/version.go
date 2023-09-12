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

package schema // import "go.opentelemetry.io/otel/sdk/resource/internal/schema"

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Masterminds/semver/v3"
)

func GreatestVersion(schemaURLs ...string) (string, error) {
	greatest := struct {
		Version *semver.Version
		URL     string
	}{}
	var invalid []string
	for _, u := range schemaURLs {
		v, err := Version(u)
		if err != nil {
			invalid = append(invalid, u)
			continue
		}

		if greatest.Version == nil || greatest.Version.LessThan(v) {
			greatest.Version = v
			greatest.URL = u
		}
	}

	var err error
	if len(invalid) > 0 {
		err = fmt.Errorf("invalid schema URL: %v", invalid)
	}
	return greatest.URL, err
}

func Version(schemaURL string) (*semver.Version, error) {
	// https://github.com/open-telemetry/oteps/blob/main/text/0152-telemetry-schemas.md#schema-url
	u, err := url.Parse(schemaURL)
	if err != nil {
		return nil, err
	}

	return semver.NewVersion(u.Path[strings.LastIndex(u.Path, "/")+1:])
}
