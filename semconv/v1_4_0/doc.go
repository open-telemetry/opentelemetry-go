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

// Package semconv implements OpenTelemetry semantic conventions start from specification
// version 1.4.0 and until the next version that introduced any changes.
//
// This package is currently in a pre-GA phase. Backwards incompatible changes
// may be introduced in subsequent minor version releases as we work to track
// the evolving OpenTelemetry specification and user feedback.
//
// OpenTelemetry semantic conventions are agreed standardized naming
// patterns for OpenTelemetry things. This package aims to be the
// centralized place to interact with these conventions.
//
// Example usage:
//
//  import semconv "go.opentelemetry.io/otel/semconv/v1_4_0"
//  ...
//  // Get a Tracer associcated with the Schema URL matching the semantic conventions
//  // that we want to use.
//	tracer := provider.Tracer("example_library", trace.WithSchemaURL(semconv.SchemaURL))
//	...
//  // Crate a Span.
//	var span trace.Span
//	_, span = tracer.Start(
//		context.Background(),
//		"operation",
//		// Note that here we use the semantic conventon for HostName that matches
//		// the Schema URL that we used in the Tracer() call.
//		trace.WithAttributes(semconv.HostNameKey.String("example.com")),
//	)
//	span.End()
//
package v1_4_0

// Note: semantic conventions are immutable. They can only change if a new specification
// version is introduced. If you need to change a semantic convention definition in this
// package it means there is a new specification version that introduced the change.
// In that case copy this package, give it a new name that corresponds to the new version
// number and modify the semantic convention definitions.
// Important: make sure to update SchemaURL string correspondingly.
