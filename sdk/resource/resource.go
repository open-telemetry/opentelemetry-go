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

var emptyResource Resource

// New creates a resource from a set of attributes.
// If there are duplicate keys then the last value of the key is preserved.
func New(kvs ...core.KeyValue) *Resource {
	return &Resource{
		labels: label.NewSet(kvs...),
	}
}

func Empty() *Resource {
	return &emptyResource
}

// @@@ Note this allocates a copy
func (r *Resource) Attributes() []core.KeyValue {
	if r == nil {
		r = Empty()
	}
	return r.labels.ToSlice()
}

func (r *Resource) Equal(eq *Resource) bool {
	if r == nil {
		r = Empty()
	}
	if eq == nil {
		eq = Empty()
	}
	return r.Equivalent() == eq.Equivalent()
}

func (r *Resource) Len() int {
	if r == nil {
		r = Empty()
	}
	return r.labels.Len()
}

func (r *Resource) Equivalent() label.Distinct {
	if r == nil {
		r = Empty()
	}
	return r.labels.Equivalent()
}

func (r *Resource) MarshalJSON() ([]byte, error) {
	if r == nil {
		r = Empty()
	}
	return r.labels.MarshalJSON()
}

func (r *Resource) String() string {
	if r == nil {
		r = Empty()
	}
	return r.labels.Encoded(label.DefaultEncoder())
}

func (r *Resource) Iter() label.Iterator {
	if r == nil {
		r = Empty()
	}
	return r.labels.Iter()
}

func Merge(a, b *Resource) *Resource {
	if a == nil {
		a = Empty()
	}
	if b == nil {
		b = Empty()
	}
	// Note: 'b' is listed first so that 'a' will overwrite with
	// last-value-wins in label.New()
	combine := append(b.Attributes(), a.Attributes()...)
	return New(combine...)
}
