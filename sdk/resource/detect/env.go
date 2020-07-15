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

// Package detect provides a detector that loads resource information from
// the OTEL_RESOURCE environment variable. A list of labels of the form
// `<key1>=<value1>,<key2>=<value2>,...` is accepted. Domain names and
// paths are accepted as label keys.
package detect

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/sdk/resource"
)

// Environment variable used by "env" to decode a resource.
const envVar = "OTEL_RESOURCE_LABELS"

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
			return resource.Empty(), fmt.Errorf("%v missing tag value", p)
		}
		k, v := strings.TrimSpace(field[0]), strings.TrimSpace(field[1])

		labels[i] = kv.String(k, v)
	}

	return resource.New(labels...), nil
}
