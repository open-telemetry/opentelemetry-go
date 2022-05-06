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
	"fmt"
	"sync"

	"go.opentelemetry.io/otel/sdk/metric/export"
)

// ManualReader is a a simple Reader that allows an application to
// read metrics on demand.  It simply stores the Producer interface
// provided through registration.  Flush and Shutdown are no-ops.
type ManualReader struct {
	lock     sync.Mutex
	producer producer
	shutdown bool
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
	mr.lock.Lock()
	defer mr.lock.Unlock()
	mr.producer = p
}

// ForceFlush is a no-op, always returns nil.
func (mr *ManualReader) ForceFlush(context.Context) error {
	return nil
}

// Shutdown closes any connections and frees any resources used by the reader.
func (mr *ManualReader) Shutdown(context.Context) error {
	mr.lock.Lock()
	defer mr.lock.Unlock()
	if mr.shutdown {
		return ErrReaderShutdown
	}
	mr.shutdown = true
	return nil
}

// Collect gathers all metrics from the SDK, calling any callbacks necessary.
// Collect will return an error if called after shutdown.
func (mr *ManualReader) Collect(ctx context.Context) (export.Metrics, error) {
	mr.lock.Lock()
	defer mr.lock.Unlock()
	if mr.producer == nil {
		return export.Metrics{}, ErrReaderNotRegistered
	}
	if mr.shutdown {
		return export.Metrics{}, ErrReaderShutdown
	}
	return mr.producer.produce(ctx), nil
}

var ErrReaderNotRegistered = fmt.Errorf("reader is not registered")
var ErrReaderShutdown = fmt.Errorf("reader is shutdown")
