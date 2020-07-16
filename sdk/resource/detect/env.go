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

package detect

import (
	"context"
	"errors"
	"net/url"
	"os"
	"strings"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/sdk/resource"
)

// Environment variable used by "env" to decode a resource.
const envVar = "OTEL_RESOURCE_LABELS"

var (
	//ErrMissingPairValue happens when a pair does not include value
	ErrMissingPairValue = errors.New("missing pair value")

	//ErrUnescape happens when '%' is not followed by two hexadecimal digits
	ErrUnescape = errors.New("invalid resource format in attribute")
)

// FromEnv is a detector that implments the ResourceDetector and collects resources
// from environment
type FromEnv struct{}

// Detect is a function that collects resources from environment
func (d *FromEnv) Detect(context.Context) (*resource.Resource, error) {
	labels := strings.TrimSpace(os.Getenv(envVar))

	if labels == "" {
		return resource.Empty(), nil
	}
	return constructOTResources(labels)
}

func constructOTResources(s string) (*resource.Resource, error) {
	pairs := strings.Split(s, ",")
	labels := make([]kv.KeyValue, len(pairs))
	for i, p := range pairs {
		field := strings.SplitN(p, "=", 2)
		if len(field) != 2 {
			return resource.Empty(), ErrMissingPairValue
		}
		k, v := strings.TrimSpace(field[0]), strings.TrimSpace(field[1])

		var err error
		if v, err = url.QueryUnescape(v); err != nil {
			return resource.Empty(), ErrUnescape
		}

		labels[i] = kv.String(k, v)
	}

	return resource.New(labels...), nil
}
