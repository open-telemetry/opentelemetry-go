// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

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
					assert.Empty(t, r.SchemaURL())
				}
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, c.want.SchemaURL(), r.SchemaURL())
			assert.ElementsMatch(t, c.want.Attributes(), r.Attributes())
		})
	}
}
