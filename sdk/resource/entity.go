package resource

import "go.opentelemetry.io/otel/attribute"

type Entity struct {
	Type      string
	Id        attribute.Set
	Attrs     attribute.Set
	SchemaURL string
}
