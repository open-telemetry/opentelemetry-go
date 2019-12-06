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
	ticker       Ticker
	clock        Clock
}

var _ metric.Provider = &Controller{}

// Several types below are created to match "github.com/benbjohnson/clock"
// so that it remains a test-only dependency.

type Clock interface {
	Now() time.Time
	Ticker(time.Duration) Ticker
}

type Ticker interface {
	Stop()
	C() <-chan time.Time
}

type realClock struct {
}

type realTicker struct {
	ticker *time.Ticker
}

var _ Clock = realClock{}
var _ Ticker = realTicker{}

// New constructs a Controller, an implementation of metric.Provider,
// using the provided batcher, exporter, and collection period to
// configure an SDK with periodic collection.  The batcher itself is
// configured with the aggregation selector policy.
//
// If the Exporter implements the export.LabelEncoder interface, the
// exporter will be used as the label encoder for the SDK itself,
// otherwise the SDK will be configured with the default label
// encoder.
func New(batcher export.Batcher, exporter export.Exporter, period time.Duration) *Controller {
	lencoder, _ := exporter.(export.LabelEncoder)

	if lencoder == nil {
		lencoder = sdk.NewDefaultLabelEncoder()
	}

	return &Controller{
		sdk:          sdk.New(batcher, lencoder),
		errorHandler: sdk.DefaultErrorHandler,
		batcher:      batcher,
		exporter:     exporter,
		ch:           make(chan struct{}),
		period:       period,
		clock:        realClock{},
	}
}

// SetClock supports setting a mock clock for testing.  This must be
// called before Start().
func (c *Controller) SetClock(clock Clock) {
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

// Meter returns a named Meter, satisifying the metric.Provider
// interface.
func (c *Controller) Meter(name string) metric.Meter {
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
		case <-c.ticker.C():
			c.tick()
		}
	}
}

func (c *Controller) tick() {
	// TODO: either remove the context argument from Export() or
	// configure a timeout here?
	ctx := context.Background()
	c.sdk.Collect(ctx)
	err := c.exporter.Export(ctx, c.batcher.CheckpointSet())
	c.batcher.FinishedCollection()

	if err != nil {
		c.errorHandler(err)
	}
}

func (realClock) Now() time.Time {
	return time.Now()
}

func (realClock) Ticker(period time.Duration) Ticker {
	return realTicker{time.NewTicker(period)}
}

func (t realTicker) Stop() {
	t.ticker.Stop()
}

func (t realTicker) C() <-chan time.Time {
	return t.ticker.C
}
