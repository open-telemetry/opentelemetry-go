// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resource // import "go.opentelemetry.io/otel/sdk/resource"

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/internal/attrdedup"
)

// Entity represents an object of interest associated with produced telemetry:
// traces, metrics, or logs.
//
// For example, telemetry produced using the OpenTelemetry SDK is normally
// associated with a Service entity. Similarly, OpenTelemetry defines system
// metrics for a host. The Host is the entity we want to associate metrics
// with in this case.
//
// Entity is an immutable object.
type Entity struct {
	typ         string
	id          attribute.Set
	description attribute.Set
	schemaURL   string
}

type entityConfig struct {
	description []attribute.KeyValue
	schemaURL   string
}

// EntityOption applies a configuration option to an [Entity].
type EntityOption interface {
	applyEntity(entityConfig) entityConfig
}

type entityDescriptionOption []attribute.KeyValue

func (o entityDescriptionOption) applyEntity(cfg entityConfig) entityConfig {
	cfg.description = append(cfg.description, []attribute.KeyValue(o)...)
	return cfg
}

// WithEntityDescription sets descriptive (non-identifying) attributes for an [Entity].
// If duplicate keys are provided, the last value is used. Any identifying attribute
// key already present in the Entity ID will take precedence over descriptive attributes
// with the same key.
func WithEntityDescription(attrs ...attribute.KeyValue) EntityOption {
	return entityDescriptionOption(attrs)
}

type entitySchemaURLOption string

func (o entitySchemaURLOption) applyEntity(cfg entityConfig) entityConfig {
	cfg.schemaURL = string(o)
	return cfg
}

// WithEntitySchemaURL sets the OpenTelemetry schema URL for an [Entity].
func WithEntitySchemaURL(schemaURL string) EntityOption {
	return entitySchemaURLOption(schemaURL)
}

// NewEntity returns an [Entity] with type typ and identifying attributes id.
// Additional configuration such as descriptive attributes and schema URL can be
// provided via opts.
//
// If id contains duplicate keys, the last value is used. Invalid attributes
// are dropped. Any descriptive attribute key in opts that conflicts with a key
// in id is dropped so that identifying attributes always take precedence.
func NewEntity(typ string, id []attribute.KeyValue, opts ...EntityOption) *Entity {
	cfg := entityConfig{}
	for _, opt := range opts {
		cfg = opt.applyEntity(cfg)
	}

	idAttrs, _ := attrdedup.KeyValues(id)
	idSet, _ := attribute.NewSetWithFiltered(idAttrs, func(kv attribute.KeyValue) bool {
		return kv.Valid()
	})

	descAttrs, _ := attrdedup.KeyValues(cfg.description)
	descSet, _ := attribute.NewSetWithFiltered(descAttrs, func(kv attribute.KeyValue) bool {
		if !kv.Valid() {
			return false
		}
		return !idSet.HasValue(kv.Key)
	})

	return &Entity{
		typ:         typ,
		id:          idSet,
		description: descSet,
		schemaURL:   cfg.schemaURL,
	}
}

// Type returns the entity type string (e.g., "service", "host", "process").
func (e *Entity) Type() string {
	if e == nil {
		return ""
	}
	return e.typ
}

// ID returns a copy of the identifying attributes of the entity.
func (e *Entity) ID() []attribute.KeyValue {
	if e == nil {
		return nil
	}
	return e.id.ToSlice()
}

// Description returns a copy of the descriptive (non-identifying) attributes of the entity.
func (e *Entity) Description() []attribute.KeyValue {
	if e == nil {
		return nil
	}
	return e.description.ToSlice()
}

// SchemaURL returns the OpenTelemetry schema URL associated with the entity.
func (e *Entity) SchemaURL() string {
	if e == nil {
		return ""
	}
	return e.schemaURL
}

// Attributes returns all attributes of the entity (ID and Description combined).
func (e *Entity) Attributes() []attribute.KeyValue {
	if e == nil {
		return nil
	}
	mi := attribute.NewMergeIterator(&e.description, &e.id)
	out := make([]attribute.KeyValue, 0, e.id.Len()+e.description.Len())
	for mi.Next() {
		out = append(out, mi.Attribute())
	}
	return out
}

// Equal reports whether e and o represent the same entity.
// Two entities are equal if they have the same type, equivalent identifying
// attributes, equivalent descriptive attributes, and the same schema URL.
func (e *Entity) Equal(o *Entity) bool {
	if e == nil && o == nil {
		return true
	}
	if e == nil || o == nil {
		return false
	}
	return e.typ == o.typ &&
		e.schemaURL == o.schemaURL &&
		e.id.Equivalent() == o.id.Equivalent() &&
		e.description.Equivalent() == o.description.Equivalent()
}

// mergeEntities merges a base slice of entities with a next slice of entities
// following the OpenTelemetry Entity specification merge algorithm:
// https://github.com/open-telemetry/opentelemetry-specification/blob/main/oteps/entities/0264-resource-and-entities.md#entity-merging-and-resource
func mergeEntities(base, next []*Entity) []*Entity {
	if len(base) == 0 {
		return copyEntities(next)
	}
	if len(next) == 0 {
		return copyEntities(base)
	}

	byType := make(map[string]*Entity, len(base)+len(next))
	order := make([]string, 0, len(base)+len(next))

	for _, e := range base {
		if e == nil {
			continue
		}
		if _, exists := byType[e.typ]; !exists {
			order = append(order, e.typ)
		}
		byType[e.typ] = e
	}

	for _, e := range next {
		if e == nil {
			continue
		}
		old, exists := byType[e.typ]
		if !exists {
			byType[e.typ] = e
			order = append(order, e.typ)
			continue
		}

		// If existing entity and new entity have different schema URLs, drop new entity.
		if old.schemaURL != e.schemaURL {
			continue
		}
		// If identifying attributes differ, drop new entity.
		if old.id.Equivalent() != e.id.Equivalent() {
			continue
		}

		// Merge descriptive attributes: existing (base/old) descriptive attributes
		// take precedence over new (next/e) descriptive attributes for keys present in both.
		mergedDesc := append(e.Description(), old.Description()...)
		byType[e.typ] = NewEntity(
			old.typ, old.ID(),
			WithEntityDescription(mergedDesc...),
			WithEntitySchemaURL(old.schemaURL),
		)
	}

	result := make([]*Entity, 0, len(order))
	for _, typ := range order {
		if e, ok := byType[typ]; ok && e != nil {
			result = append(result, e)
		}
	}
	return result
}

func copyEntities(in []*Entity) []*Entity {
	if len(in) == 0 {
		return nil
	}
	out := make([]*Entity, 0, len(in))
	for _, e := range in {
		if e != nil {
			out = append(out, e)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func findEntityHoldingKey(entities []*Entity, key attribute.Key) (*Entity, int) {
	for i, e := range entities {
		if e == nil {
			continue
		}
		if e.id.HasValue(key) || e.description.HasValue(key) {
			return e, i
		}
	}
	return nil, -1
}

func mergeRawAttributesAndConflicts(
	baseRaw, nextRaw []attribute.KeyValue,
	entities []*Entity,
) ([]attribute.KeyValue, []*Entity) {
	raw := make([]attribute.KeyValue, 0, len(baseRaw)+len(nextRaw))
	for _, kv := range baseRaw {
		if e, _ := findEntityHoldingKey(entities, kv.Key); e == nil {
			raw = append(raw, kv)
		}
	}

	for _, kv := range nextRaw {
		for {
			e, idx := findEntityHoldingKey(entities, kv.Key)
			if e == nil {
				break
			}
			raw = append(raw, e.Attributes()...)
			entities = append(entities[:idx], entities[idx+1:]...)
		}
		raw = append(raw, kv)
	}

	deduped, _ := attrdedup.KeyValues(raw)
	return deduped, entities
}

func mergeEntitySchemaURL(entities []*Entity, baseSchema, nextSchema string) (string, error) {
	var entitySchema string
	hasEntitySchema := false

	for _, e := range entities {
		if e == nil || e.schemaURL == "" {
			continue
		}
		if !hasEntitySchema {
			entitySchema = e.schemaURL
			hasEntitySchema = true
		} else if entitySchema != e.schemaURL {
			return "", ErrSchemaURLConflict
		}
	}

	var target string
	if hasEntitySchema {
		target = entitySchema
	} else {
		switch {
		case baseSchema == "":
			target = nextSchema
		case nextSchema == "":
			target = baseSchema
		case baseSchema == nextSchema:
			target = baseSchema
		default:
			return "", ErrSchemaURLConflict
		}
	}

	if target == "" {
		return "", nil
	}

	if baseSchema != "" && baseSchema != target {
		return "", ErrSchemaURLConflict
	}
	if nextSchema != "" && nextSchema != target {
		return "", ErrSchemaURLConflict
	}
	return target, nil
}
