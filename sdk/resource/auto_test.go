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
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

func TestDetect(t *testing.T) {

	cases := []struct {
		name             string
		schema1, schema2 string
		isErr            bool
	}{
		{
			name:    "different schema urls",
			schema1: "https://opentelemetry.io/schemas/1.3.0",
			schema2: "https://opentelemetry.io/schemas/1.4.0",
			isErr:   true,
		},
		{
			name:    "same schema url",
			schema1: "https://opentelemetry.io/schemas/1.4.0",
			schema2: "https://opentelemetry.io/schemas/1.4.0",
			isErr:   false,
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("case-%s", c.name), func(t *testing.T) {
			d1 := resource.StringDetector(c.schema1, semconv.HostNameKey, os.Hostname)
			d2 := resource.StringDetector(c.schema2, semconv.HostNameKey, os.Hostname)
			r, err := resource.Detect(context.Background(), d1, d2)
			assert.NotNil(t, r)
			if c.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
