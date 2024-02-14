package internal

import "go.opentelemetry.io/otel/attribute"

// This is a quick implementation of active entities for prototyping purposes.
// A proper implementation will use providers, exporters, etc. just like all other
// signals use.

var activeEntities map[attribute.Distinct]Entity = map[attribute.Distinct]Entity{}

func init() {
	go exportActive()
}

func exportActive() {

}
