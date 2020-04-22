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

// Package resource provides functionality for resource, which capture
// identifying information about the entities for which signals are exported.
package resource

import (
	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/label"
)

// Resource describes an entity about which identifying information
// and metadata is exposed.
type Resource struct {
	labels label.Set
}

// New creates a resource from a set of attributes.
// If there are duplicates keys then the first value of the key is preserved.
func New(kvs ...core.KeyValue) *Resource {
	return &Resource{
		labels: label.NewSet(kvs...),
	}
}

// @@@ Note this allocates a copy
func (r *Resource) Attributes() []core.KeyValue {
	return r.labels.ToSlice()
}

func (r *Resource) MarshalJSON() ([]byte, error) {
	return r.labels.MarshalJSON()
}
