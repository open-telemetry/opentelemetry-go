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
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func TestDetect(t *testing.T) {
	cases := []struct {
		name    string
		schema  []string
		wantErr error
	}{
		{
			name: "two different schema urls",
			schema: []string{
				"https://opentelemetry.io/schemas/1.3.0",
				"https://opentelemetry.io/schemas/1.4.0",
			},
			wantErr: resource.ErrSchemaURLConflict,
		},
		{
			name: "three different schema urls",
			schema: []string{
				"https://opentelemetry.io/schemas/1.3.0",
				"https://opentelemetry.io/schemas/1.4.0",
				"https://opentelemetry.io/schemas/1.5.0",
			},
			wantErr: resource.ErrSchemaURLConflict,
		},
		{
			name: "same schema url",
			schema: []string{
				"https://opentelemetry.io/schemas/1.4.0",
				"https://opentelemetry.io/schemas/1.4.0",
			},
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("case-%s", c.name), func(t *testing.T) {
			detectors := make([]resource.Detector, len(c.schema))
			for i, s := range c.schema {
				detectors[i] = resource.StringDetector(s, semconv.HostNameKey, os.Hostname)
			}
			r, err := resource.Detect(context.Background(), detectors...)
			assert.NotNil(t, r)
			if c.wantErr != nil {
				assert.ErrorIs(t, err, c.wantErr)
				if errors.Is(c.wantErr, resource.ErrSchemaURLConflict) {
					assert.Equal(t, "", r.SchemaURL())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
