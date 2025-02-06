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
	"fmt"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// Resource describes an entity about which identifying information
// and metadata is exposed.  Resource is an immutable object,
// equivalent to a map from key to unique value.
//
// Resources should be passed and stored as pointers
// (`*resource.Resource`).  The `nil` value is equivalent to an empty
// Resource.
type Resource struct {
	impl *resourceImpl
}

// Compile-time check that the Resource remains comparable.
var _ map[Resource]struct{} = nil

type resourceImpl struct {
	attrs map[attribute.Key]attribute.Value

	// attrSet is cached attribute.Set representation of attrs.
	attrSet attribute.Set

	schemaURL  string
	entityRefs []resourceEntityRef
}

var (
	defaultResource     *Resource
	defaultResourceOnce sync.Once
)

var errMergeConflictSchemaURL = errors.New("cannot merge resource due to conflicting Schema URL")

func NewResource() *Resource {
	return &Resource{impl: &resourceImpl{}}
}

// New returns a Resource combined from the user-provided detectors.
func New(ctx context.Context, opts ...Option) (*Resource, error) {
	cfg := config{}
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}

	r := &Resource{impl: &resourceImpl{schemaURL: cfg.schemaURL}}
	return r, detect(ctx, r, cfg.detectors)
}

// NewWithAttributes creates a resource from attrs and associates the resource with a
// schema URL. If attrs contains duplicate keys, the last value will be used. If attrs
// contains any invalid items those items will be dropped. The attrs are assumed to be
// in a schema identified by schemaURL.
func NewWithAttributes(schemaURL string, attrs ...attribute.KeyValue) *Resource {
	resource := NewSchemaless(attrs...)
	resource.impl.schemaURL = schemaURL
	return resource
}

// NewWithEntities creates a resource from entity and attrs and associates the resource with a
// schema URL. If attrs or entityId contains duplicate keys, the last value will be used. If attrs or entityId
// contains any invalid items those items will be dropped. The attrs and entityId are assumed to be
// in a schema identified by schemaURL.
func NewWithEntities(
	entities []Entity,
) (*Resource, error) {
	resource := &Resource{impl: &resourceImpl{}}

	for _, entity := range entities {
		b := &Resource{
			impl: &resourceImpl{
				schemaURL:  entity.SchemaURL,
				attrs:      map[attribute.Key]attribute.Value{},
				entityRefs: []resourceEntityRef{{}},
			},
		}

		entityRef := &b.impl.entityRefs[0]
		entityRef.typ = entity.Type
		entityRef.id = map[attribute.Key]bool{}
		entityRef.attrs = map[attribute.Key]bool{}
		entityRef.schemaUrl = entity.SchemaURL

		ids := entity.Id.Iter()
		for ids.Next() {
			attr := ids.Attribute()
			if !attr.Valid() {
				continue
			}
			entityRef.id[attr.Key] = true
			b.impl.attrs[attr.Key] = attr.Value
		}
		attrs := entity.Attrs.Iter()
		for attrs.Next() {
			attr := attrs.Attribute()
			if !attr.Valid() {
				continue
			}
			if _, exists := b.impl.attrs[attr.Key]; exists {
				return nil, fmt.Errorf("invalid Entity, key %q is both an id and Attr", attr.Key)
			}
			entityRef.attrs[attr.Key] = true
			b.impl.attrs[attr.Key] = attr.Value
		}
		entityRef.updateCache()

		var err error
		resource, err = Merge(resource, b)
		if err != nil {
			return nil, err
		}
	}

	return resource, nil
}

// NewSchemaless creates a resource from attrs. If attrs contains duplicate keys,
// the last value will be used. If attrs contains any invalid items those items will
// be dropped. The resource will not be associated with a schema URL. If the schema
// of the attrs is known use NewWithAttributes instead.
func NewSchemaless(attrs ...attribute.KeyValue) *Resource {
	if len(attrs) == 0 {
		return &Resource{impl: &resourceImpl{}}
	}

	m := map[attribute.Key]attribute.Value{}
	for _, attr := range attrs {
		// Ensure attributes comply with the specification:
		// https://github.com/open-telemetry/opentelemetry-specification/blob/v1.20.0/specification/common/README.md#attribute
		if attr.Valid() {
			m[attr.Key] = attr.Value
		}
	}

	// If attrs only contains invalid entries do not allocate a new resource.
	if len(m) == 0 {
		return &Resource{impl: &resourceImpl{}}
	}

	r := &Resource{impl: &resourceImpl{attrs: m}}
	r.updateCache()
	return r
}

func mapAttrsToSlice(m map[attribute.Key]attribute.Value) []attribute.KeyValue {
	var kv []attribute.KeyValue
	for k, v := range m {
		kv = append(
			kv, attribute.KeyValue{
				Key:   k,
				Value: v,
			},
		)
	}
	return kv
}

func mapAttrsToSet(m map[attribute.Key]attribute.Value) attribute.Set {
	return attribute.NewSet(mapAttrsToSlice(m)...)
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
	s := mapAttrsToSet(r.impl.attrs)
	return s.Encoded(attribute.DefaultEncoder())
}

// MarshalLog is the marshaling function used by the logging system to represent this Resource.
func (r *Resource) MarshalLog() interface{} {
	return struct {
		Attributes attribute.Set
		SchemaURL  string
	}{
		Attributes: mapAttrsToSet(r.impl.attrs),
		SchemaURL:  r.impl.schemaURL,
	}
}

// Attributes returns a copy of attributes from the resource in a sorted order.
// To avoid allocating a new slice, use an iterator.
func (r *Resource) Attributes() []attribute.KeyValue {
	if r == nil {
		r = Empty()
	}
	return mapAttrsToSlice(r.impl.attrs)
}

// SchemaURL returns the schema URL associated with Resource r.
func (r *Resource) SchemaURL() string {
	if r == nil {
		return ""
	}
	return r.impl.schemaURL
}

// Iter returns an iterator of the Resource attributes.
// This is ideal to use if you do not want a copy of the attributes.
func (r *Resource) Iter() attribute.Iterator {
	if r == nil {
		r = Empty()
	}
	return r.impl.attrSet.Iter()
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
	return merge(a, b, mergeOptions{})
}

type mergeOptions struct {
	allowMultipleOfSameType bool
}

func merge(a, b *Resource, options mergeOptions) (*Resource, error) {
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
	case a.impl.schemaURL == "":
		schemaURL = b.impl.schemaURL
	case b.impl.schemaURL == "":
		schemaURL = a.impl.schemaURL
	case a.impl.schemaURL == b.impl.schemaURL:
		schemaURL = a.impl.schemaURL
	default:
		return Empty(), errMergeConflictSchemaURL
	}

	merged := &Resource{
		impl: &resourceImpl{
			attrs:     cloneAttrs(a.impl.attrs),
			schemaURL: schemaURL,
		},
	}
	if len(a.impl.entityRefs) > 0 {
		merged.impl.entityRefs = make([]resourceEntityRef, len(a.impl.entityRefs))
	}

	for k, v := range b.impl.attrs {
		merged.impl.attrs[k] = v
	}

	copy(merged.impl.entityRefs, a.impl.entityRefs)

	entityTypes := map[string]resourceEntityRef{}
	for _, er := range a.impl.entityRefs {
		entityTypes[er.typ] = er
	}

	for _, er := range b.impl.entityRefs {
		if existingEr, exists := entityTypes[er.typ]; !exists {
			merged.impl.entityRefs = append(merged.impl.entityRefs, er)
			entityTypes[er.typ] = er

			for k := range er.id {
				merged.impl.attrs[k] = b.impl.attrs[k]
			}
			for k := range er.attrs {
				merged.impl.attrs[k] = b.impl.attrs[k]
			}

		} else {
			err := mergeEntity(merged, existingEr, b, er, options)
			if err != nil {
				return nil, err
			}
		}
	}

	merged.updateCache()

	return merged, nil
}

func cloneAttrs(attrs map[attribute.Key]attribute.Value) map[attribute.Key]attribute.Value {
	m := map[attribute.Key]attribute.Value{}
	for k, v := range attrs {
		m[k] = v
	}
	return m
}

func mergeEntity(
	intoRes *Resource, intoEnt resourceEntityRef, fromRes *Resource, fromEnt resourceEntityRef,
	options mergeOptions,
) error {
	intoId, err := intoRes.getEntityId(intoEnt)
	if err != nil {
		return err
	}

	fromId, err := fromRes.getEntityId(fromEnt)
	if err != nil {
		return err
	}

	if options.allowMultipleOfSameType || intoId == fromId {
		// id is the same or allowMultipleOfSameType is set.
		if intoEnt.schemaUrl == fromEnt.schemaUrl {
			// SchemaURL is the same too.
			// Merge descriptive attributes.
			attrs, err := fromRes.getEntityDescr(fromEnt)
			if err != nil {
				return err
			}

			err = intoRes.mergeEntity(intoEnt, fromId, attrs)
			if err != nil {
				return err
			}
		} else {
			// SchemaURL is different.
			// Overwrite entity
			err = intoRes.overwriteEntity(intoEnt, fromRes, fromEnt)
			if err != nil {
				return err
			}
		}
	} else {
		// id is different
		// Overwrite entity
		err = intoRes.overwriteEntity(intoEnt, fromRes, fromEnt)
		if err != nil {
			return err
		}
	}
	return nil
}

// Empty returns an instance of Resource with no attributes. It is
// equivalent to a `nil` Resource.
func Empty() *Resource {
	return &Resource{impl: &resourceImpl{}}
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
				defaultResource = NewResource()
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
	return &r.impl.attrSet
}

// MarshalJSON encodes the resource attributes as a JSON list of { "Key":
// "...", "Value": ... } pairs in order sorted by key.
func (r *Resource) MarshalJSON() ([]byte, error) {
	if r == nil {
		r = Empty()
	}

	type entityRef struct {
		Type      string
		Id        any
		Attrs     any
		SchemaURL string
	}

	rjson := struct {
		Attributes any
		SchemaURL  string
		EntityRefs []entityRef
	}{
		Attributes: r.impl.attrSet.MarshalableToJSON(),
		SchemaURL:  r.impl.schemaURL,
	}
	for _, er := range r.impl.entityRefs {
		rjson.EntityRefs = append(
			rjson.EntityRefs, entityRef{
				Type:      er.typ,
				Id:        er.idAsSlice,
				Attrs:     er.attrsAsSlice,
				SchemaURL: er.schemaUrl,
			},
		)
	}

	return json.Marshal(rjson)
}

// Len returns the number of unique key-values in this Resource.
func (r *Resource) Len() int {
	if r == nil {
		return 0
	}
	return len(r.impl.attrs)
}

// Encoded returns an encoded representation of the resource.
func (r *Resource) Encoded(enc attribute.Encoder) string {
	if r == nil {
		return ""
	}
	return r.impl.attrSet.Encoded(enc)
}

func (r *Resource) getEntityId(entity resourceEntityRef) (attribute.Set, error) {
	return r.getAttrsByKeys(entity.id)
}

func (r *Resource) getEntityDescr(entity resourceEntityRef) (attribute.Set, error) {
	return r.getAttrsByKeys(entity.attrs)
}

func (r *Resource) getAttrsByKeys(keys map[attribute.Key]bool) (attribute.Set, error) {
	var id []attribute.KeyValue
	for key := range keys {
		val, exists := r.impl.attrs[key]
		if !exists {
			return attribute.NewSet(), fmt.Errorf(
				"invalid resourceEntityRef, key %s not found in Resource attrs", key,
			)
		}
		id = append(id, attribute.KeyValue{Key: attribute.Key(key), Value: val})
	}
	return attribute.NewSet(id...), nil
}

func (r *Resource) mergeEntity(entity resourceEntityRef, id, attrs attribute.Set) error {
	idx := r.findEntity(entity)
	if idx < 0 {
		return errors.New("invalid resourceEntityRef")
	}
	updateEnt := &r.impl.entityRefs[idx]

	iter := id.Iter()
	for iter.Next() {
		attr := iter.Attribute()
		r.impl.attrs[attr.Key] = attr.Value
		updateEnt.id[attr.Key] = true
	}

	iter = attrs.Iter()
	for iter.Next() {
		attr := iter.Attribute()
		r.impl.attrs[attr.Key] = attr.Value
		updateEnt.attrs[attr.Key] = true
	}
	return nil
}

func (r *Resource) mergeEntityDescr(entity resourceEntityRef, attrs attribute.Set) error {
	idx := r.findEntity(entity)
	if idx < 0 {
		return errors.New("invalid resourceEntityRef")
	}
	updateEnt := &r.impl.entityRefs[idx]

	iter := attrs.Iter()
	for iter.Next() {
		idAttr := iter.Attribute()
		r.impl.attrs[idAttr.Key] = idAttr.Value
		updateEnt.attrs[idAttr.Key] = true
	}
	return nil
}

func (r *Resource) findEntity(entity resourceEntityRef) int {
	for i, e := range r.impl.entityRefs {
		if e.typ == entity.typ && equalEntityIdKeys(e.id, entity.id) {
			return i
		}
	}
	return -1
}

func (r *Resource) overwriteEntity(
	intoEnt resourceEntityRef, fromRes *Resource, fromEnt resourceEntityRef,
) error {
	idx := r.findEntity(intoEnt)
	if idx < 0 {
		return errors.New("invalid resourceEntityRef")
	}
	updateEnt := &r.impl.entityRefs[idx]

	updateEnt.typ = fromEnt.typ
	updateEnt.schemaUrl = fromEnt.schemaUrl

	id, err := fromRes.getEntityId(fromEnt)
	if err != nil {
		return err
	}

	r.setEntityId(updateEnt, id)
	attrs, err := fromRes.getEntityDescr(fromEnt)
	if err != nil {
		return err
	}
	r.setEntityDescr(updateEnt, attrs)

	return nil
}

func (r *Resource) setEntityId(ent *resourceEntityRef, id attribute.Set) {
	iter := id.Iter()
	ent.id = map[attribute.Key]bool{}
	for iter.Next() {
		attr := iter.Attribute()
		ent.id[attr.Key] = true
		r.impl.attrs[attr.Key] = attr.Value
	}
}

func (r *Resource) setEntityDescr(ent *resourceEntityRef, attrs attribute.Set) {
	iter := attrs.Iter()
	ent.attrs = map[attribute.Key]bool{}
	for iter.Next() {
		attr := iter.Attribute()
		ent.attrs[attr.Key] = true
		r.impl.attrs[attr.Key] = attr.Value
	}

}

func (r *Resource) updateCache() {
	r.impl.attrSet = mapAttrsToSet(r.impl.attrs)
	for i := range r.impl.entityRefs {
		r.impl.entityRefs[i].updateCache()
	}
}

func (r *Resource) EntityRefs() []resourceEntityRef {
	return r.impl.entityRefs
}

func equalEntityIdKeys(id1 map[attribute.Key]bool, id2 map[attribute.Key]bool) bool {
	if len(id1) != len(id2) {
		return false
	}
	for k := range id1 {
		if !id2[k] {
			return false
		}
	}
	return true
}
