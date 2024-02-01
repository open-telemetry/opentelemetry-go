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
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
)

type detector struct {
	SchemaURL  string
	Attributes []attribute.KeyValue
}

func newDetector(schemaURL string, attrs ...attribute.KeyValue) resource.Detector {
	return detector{schemaURL, attrs}
}

func (d detector) Detect(context.Context) (*resource.Resource, error) {
	return resource.NewWithAttributes(d.SchemaURL, d.Attributes...), nil
}

func TestDetect(t *testing.T) {
	v130 := "https://opentelemetry.io/schemas/1.3.0"
	v140 := "https://opentelemetry.io/schemas/1.4.0"
	v150 := "https://opentelemetry.io/schemas/1.5.0"

	alice := attribute.String("name", "Alice")
	bob := attribute.String("name", "Bob")
	carol := attribute.String("name", "Carol")

	admin := attribute.Bool("admin", true)
	user := attribute.Bool("admin", false)

	cases := []struct {
		name      string
		detectors []resource.Detector
		want      *resource.Resource
		wantErr   error
	}{
		{
			name: "two different schema urls",
			detectors: []resource.Detector{
				newDetector(v130, alice, admin),
				newDetector(v140, bob, user),
			},
			want:    resource.NewSchemaless(bob, user),
			wantErr: resource.ErrSchemaURLConflict,
		},
		{
			name: "three different schema urls",
			detectors: []resource.Detector{
				newDetector(v130, alice, admin),
				newDetector(v140, bob, user),
				newDetector(v150, carol),
			},
			want:    resource.NewSchemaless(carol, user),
			wantErr: resource.ErrSchemaURLConflict,
		},
		{
			name: "same schema url",
			detectors: []resource.Detector{
				newDetector(v140, alice, admin),
				newDetector(v140, bob, user),
			},
			want: resource.NewWithAttributes(v140, bob, user),
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("case-%s", c.name), func(t *testing.T) {
			r, err := resource.Detect(context.Background(), c.detectors...)
			if c.wantErr != nil {
				assert.ErrorIs(t, err, c.wantErr)
				if errors.Is(c.wantErr, resource.ErrSchemaURLConflict) {
					assert.Zero(t, r.SchemaURL())
				}
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, c.want.SchemaURL(), r.SchemaURL())
			assert.ElementsMatch(t, c.want.Attributes(), r.Attributes())
		})
	}
}
