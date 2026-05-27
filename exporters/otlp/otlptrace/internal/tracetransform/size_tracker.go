// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package tracetransform // import "go.opentelemetry.io/otel/exporters/otlp/otlptrace/internal/tracetransform"

import (
	math_bits "math/bits"

	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
	"google.golang.org/protobuf/proto"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

type resourceScopeKey struct {
	r  attribute.Distinct
	is instrumentation.Scope
}

type resourceSizeState struct {
	size int
}

type scopeSizeState struct {
	size     int
	resource *resourceSizeState
}

// SizeTracker incrementally tracks OTLP trace request size.
type SizeTracker struct {
	total     int
	resources map[attribute.Distinct]*resourceSizeState
	scopes    map[resourceScopeKey]*scopeSizeState
}

// NewSizeTracker returns a trace request size tracker.
func NewSizeTracker() *SizeTracker {
	return &SizeTracker{
		resources: make(map[attribute.Distinct]*resourceSizeState),
		scopes:    make(map[resourceScopeKey]*scopeSizeState),
	}
}

// Add includes sd in the tracked request and returns the new request size.
func (t *SizeTracker) Add(sd tracesdk.ReadOnlySpan) int {
	if sd == nil {
		return t.total
	}

	spanMsg := span(sd)
	spanSize := proto.Size(spanMsg)

	rKey := sd.Resource().Equivalent()
	sKey := resourceScopeKey{r: rKey, is: sd.InstrumentationScope()}

	if state, ok := t.scopes[sKey]; ok {
		oldScopeSize := state.size
		newScopeSize := oldScopeSize + repeatedFieldSize(spanSize)
		state.size = newScopeSize

		oldResourceSize := state.resource.size
		newResourceSize := oldResourceSize + repeatedFieldSize(newScopeSize) - repeatedFieldSize(oldScopeSize)
		state.resource.size = newResourceSize

		t.total += repeatedFieldSize(newResourceSize) - repeatedFieldSize(oldResourceSize)
		return t.total
	}

	scopeMsg := &tracepb.ScopeSpans{
		Scope:     InstrumentationScope(sd.InstrumentationScope()),
		Spans:     []*tracepb.Span{spanMsg},
		SchemaUrl: sd.InstrumentationScope().SchemaURL,
	}
	scopeSize := proto.Size(scopeMsg)

	resourceState, ok := t.resources[rKey]
	if !ok {
		resourceMsg := &tracepb.ResourceSpans{
			Resource:   Resource(sd.Resource()),
			ScopeSpans: []*tracepb.ScopeSpans{scopeMsg},
			SchemaUrl:  sd.Resource().SchemaURL(),
		}
		resourceState = &resourceSizeState{size: proto.Size(resourceMsg)}
		t.resources[rKey] = resourceState
		t.total += repeatedFieldSize(resourceState.size)
		t.scopes[sKey] = &scopeSizeState{size: scopeSize, resource: resourceState}
		return t.total
	}

	oldResourceSize := resourceState.size
	resourceState.size += repeatedFieldSize(scopeSize)
	t.total += repeatedFieldSize(resourceState.size) - repeatedFieldSize(oldResourceSize)
	t.scopes[sKey] = &scopeSizeState{size: scopeSize, resource: resourceState}
	return t.total
}

func repeatedFieldSize(size int) int {
	return 1 + size + sov(size)
}

func sov(x int) int {
	return (math_bits.Len(uint(x)|1) + 6) / 7
}
