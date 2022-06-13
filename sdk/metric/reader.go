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

	"go.opentelemetry.io/otel/sdk/metric/export"
)

// errDuplicateRegister is logged by a Reader when an attempt to registered it
// more than once occurs.
var errDuplicateRegister = fmt.Errorf("duplicate reader registration")

// ErrReaderNotRegistered is returned if Collect or Shutdown are called before
// the reader is registered with a MeterProvider.
var ErrReaderNotRegistered = fmt.Errorf("reader is not registered")

// ErrReaderShutdown is returned if Collect or Shutdown are called after a
// reader has been Shutdown once.
var ErrReaderShutdown = fmt.Errorf("reader is shutdown")

// Reader is the interface used between the SDK and an
// exporter.  Control flow is bi-directional through the
// Reader, since the SDK initiates ForceFlush and Shutdown
// while the initiates collection.  The Register() method here
// informs the Reader that it can begin reading, signaling the
// start of bi-directional control flow.
//
// Typically, push-based exporters that are periodic will
// implement PeroidicExporter themselves and construct a
// PeriodicReader to satisfy this interface.
//
// Pull-based exporters will typically implement Register
// themselves, since they read on demand.
type Reader interface {
	// register registers a Reader with a MeterProvider.
	// The producer argument allows the Reader to signal the sdk to collect
	// and send aggregated metric measurements.
	register(producer)

	// temporality reports the Temporality for the instrument kind provided.
	temporality(InstrumentKind) Temporality

	// Collect gathers and returns all metric data related to the Reader from
	// the SDK. An error is returned if this is called after Shutdown.
	Collect(context.Context) (export.Metrics, error)

	// ForceFlush flushes all metric measurements held in an export pipeline.
	//
	// This deadline or cancellation of the passed context are honored. An appropriate
	// error will be returned in these situations. There is no guaranteed that all
	// telemetry be flushed or all resources have been released in these
	// situations.
	ForceFlush(context.Context) error

	// Shutdown flushes all metric measurements held in an export pipeline and releases any
	// held computational resources.
	//
	// This deadline or cancellation of the passed context are honored. An appropriate
	// error will be returned in these situations. There is no guaranteed that all
	// telemetry be flushed or all resources have been released in these
	// situations.
	//
	// After Shutdown is called, calls to Collect will perform no operation and instead will return
	// an error indicating the shutdown state.
	Shutdown(context.Context) error
}

//  producer produces metrics for a Reader.
type producer interface {
	// produce returns aggregated metrics from a single collection.
	//
	// This method is safe to call concurrently.
	produce(context.Context) (export.Metrics, error)
}

// produceHolder is used as an atomic.Value to wrap the non-concrete producer
// type.
type produceHolder struct {
	produce func(context.Context) (export.Metrics, error)
}

// shutdownProducer produces an ErrReaderShutdown error always.
type shutdownProducer struct{}

// produce returns an ErrReaderShutdown error.
func (p shutdownProducer) produce(context.Context) (export.Metrics, error) {
	return export.Metrics{}, ErrReaderShutdown
}

// ReaderOption applies a configuration option value to either a ManualReader or
// a PeriodicReader.
type ReaderOption interface {
	ManualReaderOption
	PeriodicReaderOption
}

// WithTemporality uses the selector to determine the Temporality measurements
// from instrument should be recorded with.
func WithTemporality(selector func(instrument InstrumentKind) Temporality) ReaderOption {
	return temporalitySelectorOption{selector: selector}
}

type temporalitySelectorOption struct {
	selector func(instrument InstrumentKind) Temporality
}

// applyManual returns a manualReaderConfig with option applied.
func (t temporalitySelectorOption) applyManual(mrc manualReaderConfig) manualReaderConfig {
	mrc.temporalitySelector = t.selector
	return mrc
}

// applyPeriodic returns a periodicReaderConfig with option applied.
func (t temporalitySelectorOption) applyPeriodic(prc periodicReaderConfig) periodicReaderConfig {
	prc.temporalitySelector = t.selector
	return prc
}

// defaultTemporalitySelector returns the default Temporality measurements
// from instrument should be recorded with: cumulative.
func defaultTemporalitySelector(InstrumentKind) Temporality {
	return CumulativeTemporality
}
