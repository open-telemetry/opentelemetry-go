// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resource // import "go.opentelemetry.io/otel/sdk/resource"

import "go.opentelemetry.io/otel/attribute"

type Entity struct {
	Type        string
	Id          attribute.Set
	Description attribute.Set
}

// EntitySet is based on attribute.Set. Pretend this is implemented for the PoC.
type EntitySet struct {
	hash uint64
	data any
}

func (e EntitySet) Distinct() EntityDistinct {
	return EntityDistinct{hash: e.hash}
}

// Distinct is an identifier of a Set which is very likely to be unique.
//
// Distinct should be used as a map key instead of a Set for to provide better
// performance for map operations.
type EntityDistinct struct {
	hash uint64
}

// NewEntitySet behaves similarly to attribute.NewSet.  Pretend this is implemented for the PoC.
func NewEntitySet(entities ...Entity) EntitySet {
	return EntitySet{}
}

func MergeEntities(res *Resource, entities EntitySet) *Resource {
	// Pretend that this merges entities into a resource for PoC purposes.
	return res
}
