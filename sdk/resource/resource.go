// Copyright 2020, OpenTelemetry Authors
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
	"reflect"

	"go.opentelemetry.io/otel/api/core"
)

// Resource describes an entity about which identifying information and metadata is exposed.
type Resource struct {
	labels map[core.Key]core.Value
}

// New creates a resource from a set of attributes.
// If there are duplicates keys then the first value of the key is preserved.
func New(kvs ...core.KeyValue) *Resource {
	res := &Resource{
		labels: map[core.Key]core.Value{},
	}
	for _, kv := range kvs {
		if _, ok := res.labels[kv.Key]; !ok {
			res.labels[kv.Key] = kv.Value
		}
	}
	return res
}

// Merge creates a new resource by combining resource a and b.
// If there are common key between resource a and b then value from resource a is preserved.
// If one of the resources is nil then the other resource is returned without creating a new one.
func Merge(a, b *Resource) *Resource {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	res := &Resource{
		labels: map[core.Key]core.Value{},
	}
	for k, v := range b.labels {
		res.labels[k] = v
	}
	// labels from resource a overwrite labels from resource b.
	for k, v := range a.labels {
		res.labels[k] = v
	}
	return res
}

// Attributes returns a copy of attributes from the resource.
func (r Resource) Attributes() []core.KeyValue {
	attrs := make([]core.KeyValue, 0, len(r.labels))
	for k, v := range r.labels {
		attrs = append(attrs, core.KeyValue{Key: k, Value: v})
	}
	return attrs
}

// Equal returns true if other Resource is the equal to r.
func (r Resource) Equal(other Resource) bool {
	return reflect.DeepEqual(r.labels, other.labels)
}
