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
	"strings"

	label "go.opentelemetry.io/otel/attribute"
)

// envVar is the environment variable name OpenTelemetry Resource information can be assigned to.
const envVar = "OTEL_RESOURCE_ATTRIBUTES"

var (
	// errMissingValue is returned when a resource value is missing.
	errMissingValue = fmt.Errorf("%w: missing value", ErrPartialResource)
)

// FromEnv is a Detector that implements the Detector and collects
// resources from environment.  This Detector is included as a
// builtin.  If these resource attributes are not wanted, use the
// WithFromEnv(nil) or WithoutBuiltin() options to explicitly disable
// them.
type FromEnv struct{}

// compile time assertion that FromEnv implements Detector interface
var _ Detector = FromEnv{}

// Detect collects resources from environment
func (FromEnv) Detect(context.Context) (*Resource, error) {
	labels := strings.TrimSpace(os.Getenv(envVar))

	if labels == "" {
		return Empty(), nil
	}
	return constructOTResources(labels)
}

func constructOTResources(s string) (*Resource, error) {
	pairs := strings.Split(s, ",")
	labels := []label.KeyValue{}
	var invalid []string
	for _, p := range pairs {
		field := strings.SplitN(p, "=", 2)
		if len(field) != 2 {
			invalid = append(invalid, p)
			continue
		}
		k, v := strings.TrimSpace(field[0]), strings.TrimSpace(field[1])
		labels = append(labels, label.String(k, v))
	}
	var err error
	if len(invalid) > 0 {
		err = fmt.Errorf("%w: %v", errMissingValue, invalid)
	}
	return NewWithAttributes(labels...), err
}
