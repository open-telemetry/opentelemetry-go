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

package resource

import (
	"context"
	"errors"
	"os"
	"strings"

	"go.opentelemetry.io/otel/api/kv"
)

// envVar is the environment variable name OpenTelemetry Resource information can be assigned to.
const envVar = "OTEL_RESOURCE_LABELS"

var (
	//ErrMissingValue is returned when a resource value is missing.
	ErrMissingValue = errors.New("missing value")

	//ErrUnescape happens when '%' is not followed by two hexadecimal digits
	ErrUnescape = errors.New("invalid resource format in attribute")
)

// FromEnv is a detector that implements the Detector and collects resources
// from environment
type FromEnv struct{}

// compile time assertion that FromEnv implements Detector interface
var _ Detector = (*FromEnv)(nil)

// Detect collects resources from environment
func (d *FromEnv) Detect(context.Context) (*Resource, error) {
	labels := strings.TrimSpace(os.Getenv(envVar))

	if labels == "" {
		return Empty(), nil
	}
	return constructOTResources(labels)
}

func constructOTResources(s string) (*Resource, error) {
	pairs := strings.Split(s, ",")
	labels := make([]kv.KeyValue, len(pairs))
	for i, p := range pairs {
		field := strings.SplitN(p, "=", 2)
		if len(field) != 2 {
			return Empty(), ErrMissingValue
		}
		k, v := strings.TrimSpace(field[0]), strings.TrimSpace(field[1])

		labels[i] = kv.String(k, v)
	}

	return New(labels...), nil
}
