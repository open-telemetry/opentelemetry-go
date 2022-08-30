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
Package instrumentation provides an instrumentation scope structure to be
passed to both the OpenTelemetry Tracer and Meter components.

For more information see
[this](https://github.com/open-telemetry/oteps/blob/main/text/0083-component.md).
*/
package instrumentation // import "go.opentelemetry.io/otel/sdk/instrumentation"

import "go.opentelemetry.io/otel/attribute"

// Scope represents the instrumentation source of OpenTelemetry data.
//
// Code that uses OpenTelemetry APIs or data-models to produce telemetry needs
// to be identifiable by the receiver of that data. A Scope is used for this
// purpose, it uniquely identifies that code as the source and the extent to
// which it is relevant.
type Scope struct {
	// Name is the name of the instrumentation scope. This should be the
	// Go package name of that scope.
	Name string
	// Version is the version of the instrumentation scope.
	Version string
	// SchemaURL of the telemetry emitted by the scope.
	SchemaURL string
	// Attributes describe the unique attributes of an instrumentation scope.
	//
	// These attributes are used to differentiate an instrumentation scope when
	// it emits data that belong to different domains. For example, if both
	// profiling data and client-side data are emitted as log records from the
	// same instrumentation library, they may need to be differentiated by a
	// telemetry receiver. In that case, these attributes are used to scope and
	// differentiate the data.
	Attributes attribute.Set
}
