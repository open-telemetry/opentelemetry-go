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

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"context"
)

// ManualReader is a a simple Reader that allows an application to
// read metrics on demand.  It simply stores the Producer interface
// provided through registration.  Flush and Shutdown are no-ops.
type ManualReader struct {
	Name string
	Producer
}

var _ Reader = &ManualReader{}

// NewManualReader returns an Reader that stores the Producer for
// manual use and returns a configurable `name` as its String(),
func NewManualReader(name string) *ManualReader {
	return &ManualReader{
		Name: name,
	}
}

// String returns the name of this ManualReader.
func (mr *ManualReader) String() string {
	return mr.Name
}

// Register stores the Producer which enables the caller to read
// metrics on demand.
func (mr *ManualReader) Register(p Producer) {
	mr.Producer = p
}

// ForceFlush is a no-op, always returns nil.
func (mr *ManualReader) ForceFlush(context.Context) error {
	return nil
}

// Shutdown is a no-op, always returns nil.
func (mr *ManualReader) Shutdown(context.Context) error {
	return nil
}
