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

package otel // import "go.opentelemetry.io/otel"

import (
	"go.opentelemetry.io/otel/entity"
	"go.opentelemetry.io/otel/internal/global"
)

// EntityEmitter creates a named entityEmitter that implements EntityEmitter interface.
// If the name is an empty string then provider uses default name.
//
// This is short for GetEntityEmitterProvider().EntityEmitter(name, opts...)
func EntityEmitter(name string, opts ...entity.EntityEmitterOption) entity.EntityEmitter {
	return GetEntityEmitterProvider().EntityEmitter(name, opts...)
}

// GetEntityEmitterProvider returns the registered global entity provider.
// If none is registered then an instance of NoopEntityEmitterProvider is returned.
//
// Use the entity provider to create a named entityEmitter. E.g.
//
//	entityEmitter := otel.GetEntityEmitterProvider().EntityEmitter("example.com/foo")
//
// or
//
//	entityEmitter := otel.EntityEmitter("example.com/foo")
func GetEntityEmitterProvider() entity.EntityEmitterProvider {
	return global.EntityEmitterProvider()
}

// SetEntityEmitterProvider registers `tp` as the global entity provider.
func SetEntityEmitterProvider(tp entity.EntityEmitterProvider) {
	global.SetEntityEmitterProvider(tp)
}
