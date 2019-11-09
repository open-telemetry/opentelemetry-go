// Copyright 2019, OpenTelemetry Authors
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

package push // import "go.opentelemetry.io/otel/sdk/metric/controller/push"

import (
	"context"
	"sync"
	"time"

	"github.com/benbjohnson/clock"

	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	sdk "go.opentelemetry.io/otel/sdk/metric"
)

// Controller organizes a periodic push of metric data.
type Controller struct {
	lock         sync.Mutex
	sdk          *sdk.SDK
	errorHandler sdk.ErrorHandler
	batcher      export.Batcher
	exporter     export.Exporter
	wg           sync.WaitGroup
	ch           chan struct{}
	period       time.Duration
	ticker       *clock.Ticker
	clock        clock.Clock
}

var _ metric.Provider = &Controller{}

// New constructs a Controller, an implementation of metric.Provider,
// using the provider batcher, exporter, period.  The batcher itself
// is configured with aggregation policy selection.
func New(batcher export.Batcher, exporter export.Exporter, period time.Duration) *Controller {
	lencoder, _ := exporter.(export.LabelEncoder)

	if lencoder == nil {
		lencoder = sdk.DefaultLabelEncoder()
	}

	return &Controller{
		sdk:          sdk.New(batcher, lencoder),
		errorHandler: sdk.DefaultErrorHandler,
		batcher:      batcher,
		exporter:     exporter,
		ch:           make(chan struct{}),
		period:       period,
		clock:        clock.New(),
	}
}

// SetClock supports setting a mock clock for testing.  This must be
// called before Start().
func (c *Controller) SetClock(clock clock.Clock) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.clock = clock
}

func (c *Controller) SetErrorHandler(errorHandler sdk.ErrorHandler) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.errorHandler = errorHandler
	c.sdk.SetErrorHandler(errorHandler)
}

// GetMeter returns a named Meter, satisifying the metric.Provider
// interface.
func (c *Controller) GetMeter(name string) metric.Meter {
	return c.sdk
}

// Start begins a ticker that periodically collects and exports
// metrics with the configured interval.
func (c *Controller) Start() {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.ticker != nil {
		return
	}

	c.ticker = c.clock.Ticker(c.period)
	c.wg.Add(1)
	go c.run(c.ch)
}

// Stop waits for the background goroutine to return and then collects
// and exports metrics one last time before returning.
func (c *Controller) Stop() {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.ch == nil {
		return
	}

	close(c.ch)
	c.ch = nil
	c.wg.Wait()
	c.ticker.Stop()

	c.tick()
}

func (c *Controller) run(ch chan struct{}) {
	for {
		select {
		case <-ch:
			c.wg.Done()
			return
		case <-c.ticker.C:
			c.tick()
		}
	}
}

func (c *Controller) tick() {
	ctx := context.Background()
	c.sdk.Collect(ctx)
	err := c.exporter.Export(ctx, c.batcher.ReadCheckpoint())

	if err != nil {
		c.errorHandler(err)
	}
}
