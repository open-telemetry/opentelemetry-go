package internal

import "go.opentelemetry.io/otel/attribute"

type EntityData struct {
	// Defines the producing entity type of this resource, e.g "service", "k8s.pod", etc.
	// Empty for legacy Resources that are not entity-aware.
	Type string
	// Set of attributes that identify the entity.
	// Note that a copy of identifying attributes will be also recorded in the Attrs field.
	Id attribute.Set

	Attrs attribute.Set
}

func MergeEntities(a, b *EntityData) *EntityData {
	// Note: 'b' attributes will overwrite 'a' with last-value-wins in attribute.Key()
	// Meaning this is equivalent to: append(a.Attributes(), b.Attributes()...)
	mergedAttrs := mergeAttrs(&b.Attrs, &a.Attrs)

	var mergedType string
	var mergedId attribute.Set

	if a.Type == b.Type {
		mergedType = a.Type
		mergedId = mergeAttrs(&b.Id, &a.Id)
	} else {
		if a.Type == "" {
			mergedType = b.Type
			mergedId = b.Id
		} else if b.Type == "" {
			mergedType = a.Type
			mergedId = a.Id
		} else {
			// Different non-empty entities.
			mergedId = a.Id
			// TODO: merge the id of the updating Entity into the non-identifying
			// attributes of the old Resource, attributes from the updating Entity
			// take precedence.
			panic("not implemented")
		}
	}

	return &EntityData{
		Type:  mergedType,
		Id:    mergedId,
		Attrs: mergedAttrs,
	}
}

func mergeAttrs(a, b *attribute.Set) attribute.Set {
	if a.Len()+b.Len() == 0 {
		return *attribute.EmptySet()
	}

	mi := attribute.NewMergeIterator(a, b)
	combine := make([]attribute.KeyValue, 0, a.Len()+b.Len())
	for mi.Next() {
		combine = append(combine, mi.Attribute())
	}
	return attribute.NewSet(combine...)
}
