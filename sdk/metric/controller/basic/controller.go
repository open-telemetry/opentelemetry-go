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

package basic // import "go.opentelemetry.io/otel/sdk/metric/controller/basic"

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/registry"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	controllerTime "go.opentelemetry.io/otel/sdk/metric/controller/time"
)

// DefaultPeriod is the minimum time between collections, maximum time
// for Export().
const DefaultPeriod = 10 * time.Second

// Controller organizes a periodic push of metric data.
type Controller struct {
	lock             sync.Mutex
	accumulator      *sdk.Accumulator
	provider         *registry.MeterProvider
	checkpointer     export.Checkpointer
	exporter         export.Exporter
	wg               sync.WaitGroup
	stopCh           chan struct{}
	exportRequestCh  chan struct{}
	exportResponseCh chan struct{}
	clock            controllerTime.Clock
	ticker           controllerTime.Ticker

	collectPeriod  time.Duration
	collectTimeout time.Duration
	exportTimeout  time.Duration
}

// New constructs a Controller, an implementation of MeterProvider, using the
// provided checkpointer, exporter, and options to configure an SDK with
// periodic collection.
func New(checkpointer export.Checkpointer, opts ...Option) *Controller {
	c := &Config{
		CollectPeriod: DefaultPeriod,
	}
	for _, opt := range opts {
		opt.Apply(c)
	}
	if c.CollectTimeout == 0 {
		c.CollectTimeout = c.CollectPeriod
	}

	var exportRequestCh chan struct{}
	var exportResponseCh chan struct{}
	if c.Exporter != nil {
		exportRequestCh = make(chan struct{}, 1)
		exportResponseCh = make(chan struct{}, 1)
	}
	if c.ExportTimeout == 0 {
		c.ExportTimeout = c.CollectPeriod
	}

	impl := sdk.NewAccumulator(
		checkpointer,
		c.Resource,
	)
	return &Controller{
		provider:         registry.NewMeterProvider(impl),
		accumulator:      impl,
		checkpointer:     checkpointer,
		exporter:         c.Exporter,
		stopCh:           make(chan struct{}),
		collectPeriod:    c.CollectPeriod,
		collectTimeout:   c.CollectTimeout,
		exportTimeout:    c.ExportTimeout,
		exportRequestCh:  exportRequestCh,
		exportResponseCh: exportResponseCh,
		clock:            controllerTime.RealClock{},
	}
}

// SetClock supports setting a mock clock for testing.  This must be
// called before Start().
func (c *Controller) SetClock(clock controllerTime.Clock) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.clock = clock
}

// MeterProvider returns a MeterProvider instance for this controller.
func (c *Controller) MeterProvider() metric.MeterProvider {
	return c.provider
}

// Start begins a ticker that periodically collects and exports
// metrics with the configured interval.
func (c *Controller) Start() {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.ticker != nil {
		return
	}

	c.wg.Add(1)
	c.ticker = c.clock.Ticker(c.collectPeriod)

	if c.exporter != nil {
		c.wg.Add(2)
		go c.runTicker(c.stopCh)
		go c.runExporter(c.stopCh)
	}
}

// Stop waits for the background goroutine to return and then collects
// and exports metrics one last time before returning.
func (c *Controller) Stop() {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.stopCh == nil {
		return
	}

	close(c.stopCh)
	c.stopCh = nil
	c.wg.Done()
	c.wg.Wait()
	c.ticker.Stop()

	c.collect(nil)
}

func (c *Controller) runTicker(stopCh chan struct{}) {
	defer c.wg.Done()
	for {
		select {
		case <-stopCh:
			return
		case <-c.ticker.C():
			c.collect(stopCh)
		}
	}
}

func (c *Controller) collect(stopCh chan struct{}) {
	if c.exporter != nil {
		// Wait for the previous export to finish or timeout.
		select {
		case <-c.exportResponseCh:
			// ok
		case <-stopCh:
			return
		}

	}
	if err := c.checkpoint(); err != nil {
		otel.Handle(err)
	}
	if c.exporter != nil {
		// Begin a new export.
		select {
		case c.exportRequestCh <- struct{}{}:
			// ok
		case <-stopCh:
			return
		}
	}
}

func (c *Controller) checkpoint() error {
	ctx, cancel := context.WithTimeout(context.Background(), c.collectTimeout)
	defer cancel()

	ckpt := c.checkpointer.CheckpointSet()
	ckpt.Lock()
	defer ckpt.Unlock()

	c.checkpointer.StartCollection()
	c.accumulator.Collect(ctx)
	return c.checkpointer.FinishCollection()
}

func (c *Controller) runExporter(stopCh chan struct{}) {
	defer c.wg.Done()
	c.exportResponseCh <- struct{}{}
	for {
		select {
		case <-c.exportRequestCh:
		case <-stopCh:
			return
		}
		c.export()
		select {
		case c.exportResponseCh <- struct{}{}:
		case <-stopCh:
			return
		}
	}
}

func (c *Controller) export() {
	ctx, cancel := context.WithTimeout(context.Background(), c.exportTimeout)
	defer cancel()

	ckpt := c.checkpointer.CheckpointSet()
	ckpt.RLock()
	defer ckpt.RUnlock()

	if err := c.exporter.Export(ctx, ckpt); err != nil {
		otel.Handle(err)
	}
}
