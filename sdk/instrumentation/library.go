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

/*
Package instrumentation provides an instrumentation library structure to be
passed to both the OpenTelemetry Tracer and Meter components.

For more information see
[this](https://github.com/open-telemetry/oteps/blob/main/text/0083-component.md).
*/
package instrumentation // import "go.opentelemetry.io/otel/sdk/instrumentation"
import "go.opentelemetry.io/otel/attribute"

// Library represents the instrumentation library.
type Library struct {
	// Name is the name of the instrumentation library. This should be the
	// Go package name of that library.
	Name string
	// Version is the version of the instrumentation library.
	Version string
	// SchemaURL of the telemetry emitted by the library.
	SchemaURL string
	// Scope attributes.
	Attrs attribute.Set
}

type LibraryDistinct struct {
	name      string
	version   string
	schemaURL string
	attrs     attribute.Distinct
}

// Equivalent returns an object that can be compared for equality
// between two libraries. This value is suitable for use as a key in
// a map.
func (l Library) Equivalent() LibraryDistinct {
	return LibraryDistinct{
		name:      l.Name,
		version:   l.Version,
		schemaURL: l.SchemaURL,
		attrs:     l.Attrs.Equivalent(),
	}
}
