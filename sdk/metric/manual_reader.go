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

//go:build go1.18
// +build go1.18

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/sdk/metric/export"
)

// manualReader is a a simple Reader that allows an application to
// read metrics on demand.
type manualReader struct {
	producer     atomic.Value
	shutdownOnce sync.Once

	temporalitySelector func(InstrumentKind) Temporality
}

// Compile time check the manualReader implements Reader.
var _ Reader = &manualReader{}

// NewManualReader returns a Reader which is directly called to collect metrics.
func NewManualReader(opts ...ManualReaderOption) Reader {
	cfg := newManualReaderConfig(opts)
	return &manualReader{
		temporalitySelector: cfg.temporalitySelector,
	}
}

// register stores the Producer which enables the caller to read
// metrics on demand.
func (mr *manualReader) register(p producer) {
	// Only register once. If producer is already set, do nothing.
	if !mr.producer.CompareAndSwap(nil, produceHolder{produce: p.produce}) {
		msg := "did not register manual reader"
		global.Error(errDuplicateRegister, msg)
	}
}

// temporality reports the Temporality for the instrument kind provided.
func (mr *manualReader) temporality(kind InstrumentKind) Temporality {
	return mr.temporalitySelector(kind)
}

// ForceFlush is a no-op, it always returns nil.
func (mr *manualReader) ForceFlush(context.Context) error {
	return nil
}

// Shutdown closes any connections and frees any resources used by the reader.
func (mr *manualReader) Shutdown(context.Context) error {
	err := ErrReaderShutdown
	mr.shutdownOnce.Do(func() {
		// Any future call to Collect will now return ErrReaderShutdown.
		mr.producer.Store(produceHolder{
			produce: shutdownProducer{}.produce,
		})
		err = nil
	})
	return err
}

// Collect gathers all metrics from the SDK, calling any callbacks necessary.
// Collect will return an error if called after shutdown.
func (mr *manualReader) Collect(ctx context.Context) (export.Metrics, error) {
	p := mr.producer.Load()
	if p == nil {
		return export.Metrics{}, ErrReaderNotRegistered
	}

	ph, ok := p.(produceHolder)
	if !ok {
		// The atomic.Value is entirely in the periodicReader's control so
		// this should never happen. In the unforeseen case that this does
		// happen, return an error instead of panicking so a users code does
		// not halt in the processes.
		err := fmt.Errorf("manual reader: invalid producer: %T", p)
		return export.Metrics{}, err
	}
	return ph.produce(ctx)
}

// manualReaderConfig contains configuration options for a ManualReader.
type manualReaderConfig struct {
	temporalitySelector func(InstrumentKind) Temporality
}

// newManualReaderConfig returns a manualReaderConfig configured with options.
func newManualReaderConfig(opts []ManualReaderOption) manualReaderConfig {
	cfg := manualReaderConfig{
		temporalitySelector: defaultTemporalitySelector,
	}
	for _, opt := range opts {
		cfg = opt.applyManual(cfg)
	}
	return cfg
}

// ManualReaderOption applies a configuration option value to a ManualReader.
type ManualReaderOption interface {
	applyManual(manualReaderConfig) manualReaderConfig
}
