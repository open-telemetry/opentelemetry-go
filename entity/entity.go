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

package entity // import "go.opentelemetry.io/otel/sdk/resource"

import (
	"context"
	"errors"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// Entity describes an entity about which identifying information
// and metadata is exposed.  Entity is an immutable object,
// equivalent to a map from key to unique value.
//
// Resources should be passed and stored as pointers
// (`*resource.Entity`).  The `nil` value is equivalent to an empty
// Entity.
type Entity struct {
	attrs     attribute.Set
	schemaURL string
}

var (
	defaultResource     *Entity
	defaultResourceOnce sync.Once
)

var errMergeConflictSchemaURL = errors.New("cannot merge resource due to conflicting Schema URL")

// New returns a Entity combined from the user-provided detectors.
func New(ctx context.Context, opts ...Option) (*Entity, error) {
	cfg := config{}
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}

	r := &Entity{schemaURL: cfg.schemaURL}
	return r, detect(ctx, r, cfg.detectors)
}

// NewWithAttributes creates a resource from attrs and associates the resource with a
// schema URL. If attrs contains duplicate keys, the last value will be used. If attrs
// contains any invalid items those items will be dropped. The attrs are assumed to be
// in a schema identified by schemaURL.
func NewWithAttributes(schemaURL string, attrs ...attribute.KeyValue) *Entity {
	resource := NewSchemaless(attrs...)
	resource.schemaURL = schemaURL
	return resource
}

// String implements the Stringer interface and provides a
// human-readable form of the resource.
//
// Avoid using this representation as the key in a map of resources,
// use Equivalent() as the key instead.
func (r *Entity) String() string {
	if r == nil {
		return ""
	}
	return r.attrs.Encoded(attribute.DefaultEncoder())
}

// MarshalLog is the marshaling function used by the logging system to represent this Entity.
func (r *Entity) MarshalLog() interface{} {
	return struct {
		Attributes attribute.Set
		SchemaURL  string
	}{
		Attributes: r.attrs,
		SchemaURL:  r.schemaURL,
	}
}

// Attributes returns a copy of attributes from the resource in a sorted order.
// To avoid allocating a new slice, use an iterator.
func (r *Entity) Attributes() []attribute.KeyValue {
	if r == nil {
		r = Empty()
	}
	return r.attrs.ToSlice()
}

// SchemaURL returns the schema URL associated with Entity r.
func (r *Entity) SchemaURL() string {
	if r == nil {
		return ""
	}
	return r.schemaURL
}

// Equal returns true when a Entity is equivalent to this Entity.
func (r *Entity) Equal(eq *Entity) bool {
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
func Merge(a, b *Entity) (*Entity, error) {
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

	// Note: 'b' attributes will overwrite 'a' with last-value-wins in attribute.Key()
	// Meaning this is equivalent to: append(a.Attributes(), b.Attributes()...)
	mi := attribute.NewMergeIterator(b.Set(), a.Set())
	combine := make([]attribute.KeyValue, 0, a.Len()+b.Len())
	for mi.Next() {
		combine = append(combine, mi.Attribute())
	}
	merged := NewWithAttributes(schemaURL, combine...)
	return merged, nil
}

// Empty returns an instance of Entity with no attributes. It is
// equivalent to a `nil` Entity.
func Empty() *Entity {
	return &Entity{}
}

// Default returns an instance of Entity with a default
// "service.name" and OpenTelemetrySDK attributes.
func Default() *Entity {
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
				defaultResource = &Entity{}
			}
		},
	)
	return defaultResource
}

// Environment returns an instance of Entity with attributes
// extracted from the OTEL_RESOURCE_ATTRIBUTES environment variable.
func Environment() *Entity {
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
func (r *Entity) Equivalent() attribute.Distinct {
	return r.Set().Equivalent()
}

// MarshalJSON encodes the resource attributes as a JSON list of { "Key":
// "...", "Value": ... } pairs in order sorted by key.
func (r *Entity) MarshalJSON() ([]byte, error) {
	if r == nil {
		r = Empty()
	}
	return r.attrs.MarshalJSON()
}

// Encoded returns an encoded representation of the resource.
func (r *Entity) Encoded(enc attribute.Encoder) string {
	if r == nil {
		return ""
	}
	return r.attrs.Encoded(enc)
}
