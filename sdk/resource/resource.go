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
	"encoding/json"
	"errors"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource/internal"
)

// Resource describes an entity about which identifying information
// and metadata is exposed.  Resource is an immutable object,
// equivalent to a map from key to unique value.
//
// Resources should be passed and stored as pointers
// (`*resource.Resource`).  The `nil` value is equivalent to an empty
// Resource.
type Resource struct {
	schemaURL string

	// Producing entity.
	entity internal.EntityData
}

var (
	defaultResource     *Resource
	defaultResourceOnce sync.Once
)

var errMergeConflictSchemaURL = errors.New("cannot merge resource due to conflicting Schema URL")

// New returns a Resource combined from the user-provided detectors.
func New(ctx context.Context, opts ...Option) (*Resource, error) {
	cfg := config{}
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}

	entityId, _ := attribute.NewSetWithFiltered(
		cfg.entityId, func(kv attribute.KeyValue) bool {
			return kv.Valid()
		},
	)

	r := &Resource{
		schemaURL: cfg.schemaURL,
		entity:    internal.EntityData{Type: cfg.entityType, Id: entityId},
	}
	return r, detect(ctx, r, cfg.detectors)
}

// NewWithAttributes creates a resource from attrs and associates the resource with a
// schema URL. If attrs contains duplicate keys, the last value will be used. If attrs
// contains any invalid items those items will be dropped. The attrs are assumed to be
// in a schema identified by schemaURL.
func NewWithAttributes(schemaURL string, attrs ...attribute.KeyValue) *Resource {
	resource := NewSchemaless(attrs...)
	resource.schemaURL = schemaURL
	return resource
}

// NewWithEntity creates a resource from entity and attrs and associates the resource with a
// schema URL. If attrs or entityId contains duplicate keys, the last value will be used. If attrs or entityId
// contains any invalid items those items will be dropped. The attrs and entityId are assumed to be
// in a schema identified by schemaURL.
func NewWithEntity(
	schemaURL string, entity *internal.EntityData,
) *Resource {
	resource := NewSchemaless(entity.Attrs.ToSlice()...)
	resource.schemaURL = schemaURL
	resource.entity = *entity
	//resource.entity.Type = entityType

	// Ensure attributes comply with the specification:
	// https://github.com/open-telemetry/opentelemetry-specification/blob/v1.20.0/specification/common/README.md#attribute
	//id, _ := attribute.NewSetWithFiltered(
	//	entityId, func(kv attribute.KeyValue) bool {
	//		return kv.Valid()
	//	},
	//)
	//
	//resource.entity.Id = id
	return resource
}

// NewSchemaless creates a resource from attrs. If attrs contains duplicate keys,
// the last value will be used. If attrs contains any invalid items those items will
// be dropped. The resource will not be associated with a schema URL. If the schema
// of the attrs is known use NewWithAttributes instead.
func NewSchemaless(attrs ...attribute.KeyValue) *Resource {
	if len(attrs) == 0 {
		return &Resource{}
	}

	// Ensure attributes comply with the specification:
	// https://github.com/open-telemetry/opentelemetry-specification/blob/v1.20.0/specification/common/README.md#attribute
	s, _ := attribute.NewSetWithFiltered(
		attrs, func(kv attribute.KeyValue) bool {
			return kv.Valid()
		},
	)

	// If attrs only contains invalid entries do not allocate a new resource.
	if s.Len() == 0 {
		return &Resource{}
	}

	return &Resource{entity: internal.EntityData{Id: attribute.NewSet(), Attrs: s}} //nolint
}

// String implements the Stringer interface and provides a
// human-readable form of the resource.
//
// Avoid using this representation as the key in a map of resources,
// use Equivalent() as the key instead.
func (r *Resource) String() string {
	if r == nil {
		return ""
	}
	return r.entity.Attrs.Encoded(attribute.DefaultEncoder())
}

// MarshalLog is the marshaling function used by the logging system to represent this Resource.
func (r *Resource) MarshalLog() interface{} {
	return struct {
		Attributes attribute.Set
		SchemaURL  string
	}{
		Attributes: r.entity.Attrs,
		SchemaURL:  r.schemaURL,
	}
}

// Attributes returns a copy of attributes from the resource in a sorted order.
// To avoid allocating a new slice, use an iterator.
func (r *Resource) Attributes() []attribute.KeyValue {
	if r == nil {
		r = Empty()
	}
	return r.entity.Attrs.ToSlice()
}

// SchemaURL returns the schema URL associated with Resource r.
func (r *Resource) SchemaURL() string {
	if r == nil {
		return ""
	}
	return r.schemaURL
}

func (r *Resource) EntityId() *attribute.Set {
	if r == nil {
		return attribute.EmptySet()
	}
	return &r.entity.Id
}

func (r *Resource) EntityType() string {
	if r == nil {
		return ""
	}
	return r.entity.Type
}

// Iter returns an iterator of the Resource attributes.
// This is ideal to use if you do not want a copy of the attributes.
func (r *Resource) Iter() attribute.Iterator {
	if r == nil {
		r = Empty()
	}
	return r.entity.Attrs.Iter()
}

// Equal returns true when a Resource is equivalent to this Resource.
func (r *Resource) Equal(eq *Resource) bool {
	if r == nil {
		r = Empty()
	}
	if eq == nil {
		eq = Empty()
	}
	return r.Equivalent() == eq.Equivalent()
}

// Merge creates a new resource by combining resource a and b.
//
// If there are common keys between resource a and b, then the value
// from resource b will overwrite the value from resource a, even
// if resource b's value is empty.
//
// The SchemaURL of the resources will be merged according to the spec rules:
// https://github.com/open-telemetry/opentelemetry-specification/blob/v1.20.0/specification/resource/sdk.md#merge
// If the resources have different non-empty schemaURL an empty resource and an error
// will be returned.
func Merge(a, b *Resource) (*Resource, error) {
	if a == nil && b == nil {
		return Empty(), nil
	}
	if a == nil {
		return b, nil
	}
	if b == nil {
		return a, nil
	}

	// Merge the schema URL.
	var schemaURL string
	switch true {
	case a.schemaURL == "":
		schemaURL = b.schemaURL
	case b.schemaURL == "":
		schemaURL = a.schemaURL
	case a.schemaURL == b.schemaURL:
		schemaURL = a.schemaURL
	default:
		return Empty(), errMergeConflictSchemaURL
	}

	mergedEntity := internal.MergeEntities(&a.entity, &b.entity)
	merged := NewWithEntity(schemaURL, mergedEntity)

	return merged, nil
}

// Empty returns an instance of Resource with no attributes. It is
// equivalent to a `nil` Resource.
func Empty() *Resource {
	return &Resource{}
}

// Default returns an instance of Resource with a default
// "service.name" and OpenTelemetrySDK attributes.
func Default() *Resource {
	defaultResourceOnce.Do(
		func() {
			var err error
			defaultResource, err = Detect(
				context.Background(),
				defaultServiceNameDetector{},
				fromEnv{},
				telemetrySDK{},
			)
			if err != nil {
				otel.Handle(err)
			}
			// If Detect did not return a valid resource, fall back to emptyResource.
			if defaultResource == nil {
				defaultResource = &Resource{}
			}
		},
	)
	return defaultResource
}

// Environment returns an instance of Resource with attributes
// extracted from the OTEL_RESOURCE_ATTRIBUTES environment variable.
func Environment() *Resource {
	detector := &fromEnv{}
	resource, err := detector.Detect(context.Background())
	if err != nil {
		otel.Handle(err)
	}
	return resource
}

// Equivalent returns an object that can be compared for equality
// between two resources. This value is suitable for use as a key in
// a map.
func (r *Resource) Equivalent() attribute.Distinct {
	return r.Set().Equivalent()
}

// Set returns the equivalent *attribute.Set of this resource's attributes.
func (r *Resource) Set() *attribute.Set {
	if r == nil {
		r = Empty()
	}
	return &r.entity.Attrs
}

// MarshalJSON encodes the resource attributes as a JSON list of { "Key":
// "...", "Value": ... } pairs in order sorted by key.
func (r *Resource) MarshalJSON() ([]byte, error) {
	if r == nil {
		r = Empty()
	}

	rjson := struct {
		Attributes any
		SchemaURL  string
		Entity     struct {
			Type string
			Id   any
		}
	}{
		Attributes: r.entity.Attrs.MarshalableToJSON(),
		SchemaURL:  r.schemaURL,
		Entity: struct {
			Type string
			Id   any
		}{
			Type: r.entity.Type,
			Id:   r.entity.Id.MarshalableToJSON(),
		},
	}

	return json.Marshal(rjson)
}

// Len returns the number of unique key-values in this Resource.
func (r *Resource) Len() int {
	if r == nil {
		return 0
	}
	return r.entity.Attrs.Len()
}

// Encoded returns an encoded representation of the resource.
func (r *Resource) Encoded(enc attribute.Encoder) string {
	if r == nil {
		return ""
	}
	return r.entity.Attrs.Encoded(enc)
}
