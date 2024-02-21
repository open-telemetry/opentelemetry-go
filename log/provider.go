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

package log // import "go.opentelemetry.io/otel/log"

import "go.opentelemetry.io/otel/log/embedded"

// LoggerProvider provides access to [Logger].
//
// Warning: Methods may be added to this interface in minor releases. See
// package documentation on API implementation for information on how to set
// default behavior for unimplemented methods.
type LoggerProvider interface {
	// Users of the interface can ignore this. This embedded type is only used
	// by implementations of this interface. See the "API Implementations"
	// section of the package documentation for more information.
	embedded.LoggerProvider

	// Logger returns a new [Logger] with the provided name and configuration.
	//
	// If name is empty, implementations need to provide a default name.
	//
	// Implementations of this method need to be safe for a user to call
	// concurrently.
	Logger(name string, options ...LoggerOption) Logger
}
