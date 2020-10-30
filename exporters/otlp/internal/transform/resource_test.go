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

package transform

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/sdk/resource"
)

func TestNilResource(t *testing.T) {
	assert.Empty(t, Resource(nil))
}

func TestEmptyResource(t *testing.T) {
	assert.Empty(t, Resource(&resource.Resource{}))
}

/*
* This does not include any testing on the ordering of Resource Attributes.
* They are stored as a map internally to the Resource and their order is not
* guaranteed.
 */

func TestResourceAttributes(t *testing.T) {
	attrs := []label.KeyValue{label.Int("one", 1), label.Int("two", 2)}

	got := Resource(resource.NewWithAttributes(attrs...)).GetAttributes()
	if !assert.Len(t, attrs, 2) {
		return
	}
	assert.ElementsMatch(t, Attributes(attrs), got)
}
