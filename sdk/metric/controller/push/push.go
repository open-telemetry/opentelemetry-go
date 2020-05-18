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

package push // import "go.opentelemetry.io/otel/sdk/metric/controller/push"

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/metric/registry"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	controllerTime "go.opentelemetry.io/otel/sdk/metric/controller/time"
	"go.opentelemetry.io/otel/sdk/metric/integrator/simple"
	"go.opentelemetry.io/otel/sdk/resource"
)

const DefaultPushPeriod = 10 * time.Second

// Controller organizes a periodic push of metric data.
type Controller struct {
	lock         sync.Mutex
	accumulator  *sdk.Accumulator
	resource     *resource.Resource
	uniq         metric.MeterImpl
	named        map[string]metric.Meter
	errorHandler sdk.ErrorHandler
	integrator   *simple.Integrator
	exporter     export.Exporter
	wg           sync.WaitGroup
	ch           chan struct{}
	period       time.Duration
	ticker       controllerTime.Ticker
	clock        controllerTime.Clock
}

var _ metric.Provider = &Controller{}

// New constructs a Controller, an implementation of metric.Provider,
// using the provided exporter and options to configure an SDK with
// periodic collection.
func New(selector export.AggregationSelector, exporter export.Exporter, opts ...Option) *Controller {
	c := &Config{
		ErrorHandler: sdk.DefaultErrorHandler,
		Period:       DefaultPushPeriod,
	}
	for _, opt := range opts {
		opt.Apply(c)
	}

	integrator := simple.New(selector, c.Stateful)
	impl := sdk.NewAccumulator(integrator, sdk.WithErrorHandler(c.ErrorHandler))
	return &Controller{
		accumulator:  impl,
		resource:     c.Resource,
		uniq:         registry.NewUniqueInstrumentMeterImpl(impl),
		named:        map[string]metric.Meter{},
		errorHandler: c.ErrorHandler,
		integrator:   integrator,
		exporter:     exporter,
		ch:           make(chan struct{}),
		period:       c.Period,
		clock:        controllerTime.RealClock{},
	}
}

// SetClock supports setting a mock clock for testing.  This must be
// called before Start().
func (c *Controller) SetClock(clock controllerTime.Clock) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.clock = clock
}

func (c *Controller) SetErrorHandler(errorHandler sdk.ErrorHandler) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.errorHandler = errorHandler
	c.accumulator.SetErrorHandler(errorHandler)
}

// Meter returns a named Meter, satisifying the metric.Provider
// interface.
func (c *Controller) Meter(name string) metric.Meter {
	c.lock.Lock()
	defer c.lock.Unlock()

	if meter, ok := c.named[name]; ok {
		return meter
	}

	meter := metric.WrapMeterImpl(c.uniq, name)
	c.named[name] = meter
	return meter
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
	c.integrator.Lock()
	defer c.integrator.Unlock()

	c.accumulator.Collect(ctx)

	err := c.exporter.Export(ctx, c.resource, c.integrator.CheckpointSet())
	c.integrator.FinishedCollection()

	if err != nil {
		c.errorHandler(err)
	}
}
