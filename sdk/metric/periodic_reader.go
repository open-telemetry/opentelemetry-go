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
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric/export"
)

// Default periodic reader timing.
const (
	defaultTimeout  = time.Millisecond * 30000
	defaultInterval = time.Millisecond * 60000
)

// periodicReaderConfig contains configuration options for a PeriodicReader.
type periodicReaderConfig struct {
	interval time.Duration
	timeout  time.Duration
}

// newPeriodicReaderConfig returns a periodicReaderConfig configured with
// options.
func newPeriodicReaderConfig(options []PeriodicReaderOption) periodicReaderConfig {
	c := periodicReaderConfig{
		interval: defaultInterval,
		timeout:  defaultTimeout,
	}
	for _, o := range options {
		c = o.apply(c)
	}
	return c
}

// PeriodicReaderOption applies a configuration option value to a PeriodicReader.
type PeriodicReaderOption interface {
	apply(periodicReaderConfig) periodicReaderConfig
}

// periodicReaderOptionFunc applies a set of options to a periodicReaderConfig.
type periodicReaderOptionFunc func(periodicReaderConfig) periodicReaderConfig

// apply returns a periodicReaderConfig with option(s) applied.
func (o periodicReaderOptionFunc) apply(conf periodicReaderConfig) periodicReaderConfig {
	return o(conf)
}

// WithTimeout configures the time a PeriodicReader waits for an export to
// complete before canceling it.
//
// If this option is not used or d is less than or equal to zero, 30 seconds
// is used as the default.
func WithTimeout(d time.Duration) PeriodicReaderOption {
	return periodicReaderOptionFunc(func(conf periodicReaderConfig) periodicReaderConfig {
		if d <= 0 {
			return conf
		}
		conf.timeout = d
		return conf
	})
}

// WithInterval configures the intervening time between exports for a
// PeriodicReader.
//
// If this option is not used or d is less than or equal to zero, 60 seconds
// is used as the default.
func WithInterval(d time.Duration) PeriodicReaderOption {
	return periodicReaderOptionFunc(func(conf periodicReaderConfig) periodicReaderConfig {
		if d <= 0 {
			return conf
		}
		conf.interval = d
		return conf
	})
}

// NewPeriodicReader returns a Reader that collects and exports metric data to
// the exporter at a defined interval. By default, the returned Reader will
// collect and export data every 60 seconds, and will cancel export attempts
// that exceed 30 seconds. The export time is not counted towards the interval
// between attempts.
//
// The Collect method of the returned Reader continues to gather and return
// metric data to the user. It will not automatically send that data to the
// exporter. That is left to the user to accomplish.
func NewPeriodicReader(exporter Exporter, options ...PeriodicReaderOption) Reader {
	conf := newPeriodicReaderConfig(options)
	ctx, cancel := context.WithCancel(context.Background())
	r := &periodicReader{
		timeout:  conf.timeout,
		exporter: exporter,
		cancel:   cancel,
	}

	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		r.run(ctx, conf.interval)
	}()

	return r
}

// periodicReader is a Reader that continuously collects and exports metric
// data at a set interval.
type periodicReader struct {
	producer atomic.Value

	timeout  time.Duration
	exporter Exporter

	wg           sync.WaitGroup
	cancel       context.CancelFunc
	shutdownOnce sync.Once
}

// newTicker allows testing override.
var newTicker = time.NewTicker

// run continuously collects and exports metric data at the specified
// interval. This will run until ctx is canceled or times out.
func (r *periodicReader) run(ctx context.Context, interval time.Duration) {
	ticker := newTicker(interval)
	defer func() { ticker.Stop() }()

	for {
		select {
		case <-ticker.C:
			m, err := r.Collect(ctx)
			if err == nil {
				c, cancel := context.WithTimeout(ctx, r.timeout)
				err = r.exporter.Export(c, m)
				cancel()
			}
			if err != nil {
				otel.Handle(err)
			}
		case <-ctx.Done():
			return
		}
	}
}

// register registers p as the producer of this reader.
func (r *periodicReader) register(p producer) {
	r.producer.Store(produceHolder{produce: p.produce})
}

// Collect gathers and returns all metric data related to the Reader from
// the SDK. The returned metric data is not exported to the configured
// exporter, it is left to the caller to handle that if desired.
//
// An error is returned if this is called after Shutdown.
func (r *periodicReader) Collect(ctx context.Context) (export.Metrics, error) {
	p := r.producer.Load()
	if p == nil {
		return export.Metrics{}, ErrReaderNotRegistered
	}

	ph, ok := p.(produceHolder)
	if !ok {
		// The atomic.Value is entirely in the periodicReader's control so
		// this should never happen. In the unforeseen case that this does
		// happen, return an error instead of panicking so a users code does
		// not halt in the processes.
		err := fmt.Errorf("periodic reader: invalid producer: %T", p)
		return export.Metrics{}, err
	}
	return ph.produce(ctx)
}

// ForceFlush flushes the Exporter.
func (r *periodicReader) ForceFlush(ctx context.Context) error {
	return r.exporter.ForceFlush(ctx)
}

// Shutdown stops the export pipeline.
func (r *periodicReader) Shutdown(ctx context.Context) error {
	err := ErrReaderShutdown
	r.shutdownOnce.Do(func() {
		// Stop the run loop.
		r.cancel()
		r.wg.Wait()

		// Any future call to Collect will now return ErrReaderShutdown.
		r.producer.Store(produceHolder{
			produce: shutdownProducer{}.produce,
		})

		err = r.exporter.Shutdown(ctx)
	})
	return err
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
