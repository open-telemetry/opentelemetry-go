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

package metric // import "go.opentelemetry.io/otel/sdk/metric/reader"

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/sdk/metric/export"
)

// ManualReader is a a simple Reader that allows an application to
// read metrics on demand.  It simply stores the Producer interface
// provided through registration.  Flush and Shutdown are no-ops.
type ManualReader struct {
	producer
}

var _ Reader = &ManualReader{}

// NewManualReader returns an Reader that stores the Producer for
// manual use and returns a configurable `name` as its String(),
func NewManualReader() *ManualReader {
	return &ManualReader{}
}

// Register stores the Producer which enables the caller to read
// metrics on demand.
func (mr *ManualReader) register(p producer) {
	mr.producer = p
}

// ForceFlush is a no-op, always returns nil.
func (mr *ManualReader) ForceFlush(context.Context) error {
	return nil
}

// Shutdown is a no-op, always returns nil.
func (mr *ManualReader) Shutdown(context.Context) error {
	return nil
}

func (mr *ManualReader) Collect(ctx context.Context) (export.Metrics, error) {
	if mr.producer == nil {
		return export.Metrics{}, ErrReaderNotRegistered
	}
	return mr.produce(ctx), nil
}

var ErrReaderNotRegistered = fmt.Errorf("reader is not registered")
