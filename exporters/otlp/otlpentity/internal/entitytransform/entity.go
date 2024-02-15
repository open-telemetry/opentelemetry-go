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

package entitytransform // import "go.opentelemetry.io/otel/exporters/otlp/otlpentity/internal/entitytransform"

import (
	enititespb "go.opentelemetry.io/proto/otlp/entities/v1"

	resourcesdk "go.opentelemetry.io/otel/sdk/resource"

	"go.opentelemetry.io/otel/sdk/instrumentation"
)

// Entities transforms a slice of OpenTelemetry spans into a slice of OTLP
// ResourceSpans.
func Entities(roEntities []resourcesdk.ReadOnlyEntity) (r []*enititespb.ScopeEntities) {
	if len(roEntities) == 0 {
		return nil
	}

	type key struct {
		is instrumentation.Scope
	}
	ssm := make(map[key]*enititespb.ScopeEntities)

	for _, roEntity := range roEntities {
		if roEntity == nil {
			continue
		}

		k := key{
			is: roEntity.InstrumentationScope(),
		}
		events, iOk := ssm[k]
		if !iOk {
			// Either the resource or instrumentation scope were unknown.
			events = &enititespb.ScopeEntities{
				Scope:        InstrumentationScope(roEntity.InstrumentationScope()),
				EntityEvents: []*enititespb.EntityEvent{},
				SchemaUrl:    roEntity.InstrumentationScope().SchemaURL,
			}
		}
		events.EntityEvents = append(events.EntityEvents, entityEvent(roEntity))
		ssm[k] = events
	}

	for _, v := range ssm {
		r = append(r, v)
	}

	return r
}

// entityEvent transforms a Span into an OTLP entityEvent.
func entityEvent(sd resourcesdk.ReadOnlyEntity) *enititespb.EntityEvent {
	if sd == nil {
		return nil
	}

	s := &enititespb.EntityEvent{
		TimeUnixNano: uint64(sd.StartTime().UnixNano()),
		EntityType:   sd.Type(),
		Id:           KeyValues(sd.Id()),
		Data: &enititespb.EntityEvent_EntityState{
			EntityState: &enititespb.EntityState{
				Attributes:             KeyValues(sd.Attributes()),
				DroppedAttributesCount: uint32(sd.DroppedAttributes()),
			},
		},
	}

	return s
}
