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
	"time"

	"go.opentelemetry.io/otel/sdk/metric/data"
)

// PushExporter is an interface for push-based exporters.
type PushExporter interface {
	String() string
	ExportMetrics(context.Context, data.Metrics) error
	ShutdownMetrics(context.Context, data.Metrics) error
	ForceFlushMetrics(context.Context, data.Metrics) error
}

// PeriodicReader is an implementation of Reader that manages periodic
// exporter, flush, and shutdown.  This implementation re-uses data
// from one collection to the next, to lower memory costs.
type PeriodicReader struct {
	lock     sync.Mutex
	data     data.Metrics
	interval time.Duration
	exporter PushExporter
	producer Producer
	stop     context.CancelFunc
	wait     sync.WaitGroup
}

// NewPeriodicReader constructs a PeriodicReader from a push-based
// exporter given an interval (TODO: and options).
func NewPeriodicReader(exporter PushExporter, interval time.Duration /* opts ...Option*/) Reader {
	return &PeriodicReader{
		interval: interval,
		exporter: exporter,
	}
}

// String returns the exporter name and the configured interval.
func (pr *PeriodicReader) String() string {
	return fmt.Sprintf("%v interval %v", pr.exporter.String(), pr.interval)
}

// Register starts the periodic export loop.
func (pr *PeriodicReader) Register(producer Producer) {
	ctx, cancel := context.WithCancel(context.Background())

	pr.producer = producer
	pr.stop = cancel
	pr.wait.Add(1)

	go pr.start(ctx)
}

// start runs the export loop.
func (pr *PeriodicReader) start(ctx context.Context) {
	defer pr.wait.Done()
	ticker := time.NewTicker(pr.interval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			pr.collect(ctx, pr.exporter.ExportMetrics)
		}
	}
}

// Shutdown stops the stops the export loop, canceling its Context,
// and waits for it to return.  Then it issues a ShutdownMetrics with
// final data.
func (pr *PeriodicReader) Shutdown(ctx context.Context) error {
	pr.stop()
	pr.wait.Wait()

	return pr.collect(ctx, pr.exporter.ShutdownMetrics)
}

// ForceFlush immediately waits for an existing collection, otherwise
// immediately begins collection without regards to timing and calls
// ForceFlush with current data.
func (pr *PeriodicReader) ForceFlush(ctx context.Context) error {
	return pr.collect(ctx, pr.exporter.ForceFlushMetrics)
}

// collect serializes access to re-usable metrics data, in each case
// calling through to an underlying PushExporter method with current
// data.
func (pr *PeriodicReader) collect(ctx context.Context, method func(context.Context, data.Metrics) error) error {
	pr.lock.Lock()
	defer pr.lock.Unlock()

	// The lock ensures that re-use of `pr.data` is successful, it
	// means that shutdown, flush, and ordinary collection are
	// exclusive.  Note that shutdown will cancel an concurrent
	// (ordinary) export, while flush will wait on for a
	// concurrent export.
	pr.data = pr.producer.Produce(&pr.data)

	return method(ctx, pr.data)
}
