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
	lock         sync.Mutex
	accumulator  *sdk.Accumulator
	provider     *registry.MeterProvider
	checkpointer export.Checkpointer
	exporter     export.Exporter
	wg           sync.WaitGroup
	stopCh       chan struct{}
	clock        controllerTime.Clock
	ticker       controllerTime.Ticker

	collectPeriod  time.Duration
	collectTimeout time.Duration
	exportTimeout  time.Duration

	// collectedTime is used only in configurations with no
	// exporter, when ticker != nil.
	collectedTime time.Time
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
	if c.ExportTimeout == 0 {
		c.ExportTimeout = c.CollectPeriod
	}

	impl := sdk.NewAccumulator(
		checkpointer,
		c.Resource,
	)
	return &Controller{
		provider:       registry.NewMeterProvider(impl),
		accumulator:    impl,
		checkpointer:   checkpointer,
		exporter:       c.Exporter,
		stopCh:         make(chan struct{}),
		collectPeriod:  c.CollectPeriod,
		collectTimeout: c.CollectTimeout,
		exportTimeout:  c.ExportTimeout,
		clock:          controllerTime.RealClock{},
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
		c.wg.Add(1)
		go c.runTicker(c.stopCh)
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
	if err := c.checkpoint(func() bool {
		return true
	}); err != nil {
		otel.Handle(err)
	}
	if c.exporter == nil {
		return
	}
	// This export has to complete before starting another
	// collection since it will hold a read lock and
	// checkpoint() must re-acquire a write lock.
	if err := c.export(); err != nil {
		otel.Handle(err)
	}
}

func (c *Controller) checkpoint(cond func() bool) error {
	ckpt := c.checkpointer.CheckpointSet()
	ckpt.Lock()
	defer ckpt.Unlock()

	if !cond() {
		return nil
	}
	c.checkpointer.StartCollection()

	ctx, cancel := context.WithTimeout(context.Background(), c.collectTimeout)
	defer cancel()

	c.accumulator.Collect(ctx)

	return c.checkpointer.FinishCollection()
}

func (c *Controller) export() error {
	ckpt := c.checkpointer.CheckpointSet()
	ckpt.RLock()
	defer ckpt.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), c.exportTimeout)
	defer cancel()

	return c.exporter.Export(ctx, ckpt)
}

// Foreach gives the caller read-locked access to the current
// export.CheckpointSet.
func (c *Controller) ForEach(ks export.ExportKindSelector, f func(export.Record) error) error {
	ckpt := c.checkpointer.CheckpointSet()
	ckpt.RLock()
	defer ckpt.RUnlock()

	return ckpt.ForEach(ks, f)
}

// IsRunning returns true if the controller was started via Start().
func (c *Controller) IsRunning() bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.ticker != nil
}

// Collect requests a collection.  The collection will be skipped if
// the last collection is aged less than the CachePeriod.
func (c *Controller) Collect(ctx context.Context) error {
	if c.IsRunning() {
		// When there's a non-nil ticker, there's a goroutine
		// computing checkpoints with the collection period.
		return nil
	}

	return c.checkpoint(func() bool {
		// This is called with the CheckpointSet exclusive
		// lock held.
		if c.collectPeriod == 0 {
			return true
		}
		now := c.clock.Now()
		if now.Sub(c.collectedTime) < c.collectPeriod {
			return false
		}
		c.collectedTime = now
		return true
	})
}
