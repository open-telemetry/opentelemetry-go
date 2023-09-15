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
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/Masterminds/semver/v3"
)

var errInvalid = errors.New("invalid schema URL")

func Version(schemaURL string) (*semver.Version, error) {
	if schemaURL == "" {
		return nil, fmt.Errorf("%w: empty", errInvalid)
	}

	// https://github.com/open-telemetry/oteps/blob/main/text/0152-telemetry-schemas.md#schema-url
	u, err := url.Parse(schemaURL)
	if err != nil {
		return nil, err
	}

	if h := u.Hostname(); h != "opentelemetry.io" {
		return nil, fmt.Errorf("%w: host not \"opentelemetry.io\": %s", errInvalid, h)
	}

	i := strings.LastIndex(u.Path, "/") + 1
	if p := u.Path[:i]; p != "/schemas" {
		return nil, fmt.Errorf("%w: path not \"/schemas\": %s", errInvalid, p)
	}

	return semver.NewVersion(u.Path[i:])
}
