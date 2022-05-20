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

// manualReader is a a simple Reader that allows an application to
// read metrics on demand.
type manualReader struct {
	lock     sync.Mutex
	producer producer
	shutdown bool
}

// Compile time check the manualReader implements Reader.
var _ Reader = &manualReader{}

// NewManualReader returns a Reader which is directly called to collect metrics.
func NewManualReader() Reader {
	return &manualReader{}
}

// register stores the Producer which enables the caller to read
// metrics on demand.
func (mr *manualReader) register(p producer) {
	mr.lock.Lock()
	defer mr.lock.Unlock()
	mr.producer = p
}

// ForceFlush is a no-op, it always returns nil.
func (mr *manualReader) ForceFlush(context.Context) error {
	return nil
}

// Shutdown closes any connections and frees any resources used by the reader.
func (mr *manualReader) Shutdown(context.Context) error {
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
func (mr *manualReader) Collect(ctx context.Context) (export.Metrics, error) {
	mr.lock.Lock()
	defer mr.lock.Unlock()
	if mr.producer == nil {
		return export.Metrics{}, ErrReaderNotRegistered
	}
	if mr.shutdown {
		return export.Metrics{}, ErrReaderShutdown
	}
	return mr.producer.produce(ctx)
}

// ErrReaderNotRegistered is returned if Collect or Shutdown are called before
// the reader is registered with a MeterProvider.
var ErrReaderNotRegistered = fmt.Errorf("reader is not registered")

// ErrReaderShutdown is returned if Collect or Shutdown are called after a
// reader has been Shutdown once.
var ErrReaderShutdown = fmt.Errorf("reader is shutdown")
