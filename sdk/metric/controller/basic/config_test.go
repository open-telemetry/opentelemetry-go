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

package basic

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
)

func TestWithResource(t *testing.T) {
	r := resource.NewWithAttributes(attribute.String("A", "a"))

	c := &Config{}
	WithResource(r).Apply(c)
	assert.Equal(t, r.Equivalent(), c.Resource.Equivalent())

	// Ensure overwriting works.
	c = &Config{Resource: &resource.Resource{}}
	WithResource(r).Apply(c)
	assert.Equal(t, r.Equivalent(), c.Resource.Equivalent())
}
