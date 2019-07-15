// Copyright 2019, OpenTelemetry Authors
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

package event

import (
	"go.opentelemetry.io/api/core"
)

// Event interface provides methods to retrieve Event properties.
type Event interface {

	// Message interface retrieves message string of the Event.
	Message() string

	// Attributes interface returns a copy of attributes associated with the Event.
	Attributes() []core.KeyValue
}
