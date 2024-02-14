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

// Package embedded provides interfaces embedded within the [OpenTelemetry
// entity API].
//
// Implementers of the [OpenTelemetry entity API] can embed the relevant type
// from this package into their implementation directly. Doing so will result
// in a compilation error for users when the [OpenTelemetry entity API] is
// extended (which is something that can happen without a major version bump of
// the API package).
//
// [OpenTelemetry entity API]: https://pkg.go.dev/go.opentelemetry.io/otel/entity
package embedded // import "go.opentelemetry.io/otel/entity/embedded"

// EntityEmitterProvider is embedded in
// [go.opentelemetry.io/otel/entity.EntityEmitterProvider].
//
// Embed this interface in your implementation of the
// [go.opentelemetry.io/otel/entity.EntityEmitterProvider] if you want users to
// experience a compilation error, signaling they need to update to your latest
// implementation, when the [go.opentelemetry.io/otel/entity.EntityEmitterProvider]
// interface is extended (which is something that can happen without a major
// version bump of the API package).
type EntityEmitterProvider interface{ entityEmitterProvider() }

// EntityEmitter is embedded in [go.opentelemetry.io/otel/entity.EntityEmitter].
//
// Embed this interface in your implementation of the
// [go.opentelemetry.io/otel/entity.EntityEmitter] if you want users to experience a
// compilation error, signaling they need to update to your latest
// implementation, when the [go.opentelemetry.io/otel/entity.EntityEmitter] interface
// is extended (which is something that can happen without a major version bump
// of the API package).
type EntityEmitter interface{ entityEmitter() }

// Entity is embedded in [go.opentelemetry.io/otel/entity.Entity].
//
// Embed this interface in your implementation of the
// [go.opentelemetry.io/otel/entity.Entity] if you want users to experience a
// compilation error, signaling they need to update to your latest
// implementation, when the [go.opentelemetry.io/otel/entity.Entity] interface is
// extended (which is something that can happen without a major version bump of
// the API package).
type Entity interface{ entity() }
