package resource

import "go.opentelemetry.io/otel/attribute"

type resourceEntityRef struct {
	schemaUrl string

	// Defines the entity type, e.g "service", "k8s.pod", etc.
	typ string

	// Set of Resource attribute keys that identify the entity.
	id map[attribute.Key]bool

	// Set of Resource attribute keys that describe the entity.
	attrs map[attribute.Key]bool

	// id and attrs are cached as slices for faster exporting.
	idAsSlice    []string
	attrsAsSlice []string
}

// updateCache must be called after id or attrs are modified to make
// sure the cache is up to date.
func (r *resourceEntityRef) updateCache() {
	r.idAsSlice = r.idAsSlice[:0]
	for k := range r.id {
		r.idAsSlice = append(r.idAsSlice, string(k))
	}
	r.attrsAsSlice = r.attrsAsSlice[:0]
	for k := range r.attrs {
		r.attrsAsSlice = append(r.attrsAsSlice, string(k))
	}
}

func (r *resourceEntityRef) SchemaUrl() string {
	return r.schemaUrl
}

func (r *resourceEntityRef) Type() string {
	return r.typ
}

func (r *resourceEntityRef) Id() []string {
	return r.idAsSlice
}

func (r *resourceEntityRef) Attrs() []string {
	return r.attrsAsSlice
}
