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

package pull // import "go.opentelemetry.io/otel/sdk/metric/controller/pull"

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/metric/registry"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	controllerTime "go.opentelemetry.io/otel/sdk/metric/controller/time"
	integrator "go.opentelemetry.io/otel/sdk/metric/integrator/simple"
	"go.opentelemetry.io/otel/sdk/resource"
)

// DefaultCachePeriod determines how long a recently-computed result
// will be returned without gathering metric data again.
const DefaultCachePeriod time.Duration = 10 * time.Second

// Controller manages access to a *sdk.Accumulator and
// *simple.Integrator.  Use Provider() for obtaining Meters.  Use
// Foreach() for accessing current records.
type Controller struct {
	accumulator *sdk.Accumulator
	integrator  *integrator.Integrator
	provider    *registry.Provider
	period      time.Duration
	lastCollect time.Time
	clock       controllerTime.Clock
	checkpoint  export.CheckpointSet
}

// New returns a *Controller configured with an aggregation selector and options.
func New(selector export.AggregationSelector, options ...Option) *Controller {
	config := &Config{
		Resource:     resource.Empty(),
		ErrorHandler: sdk.DefaultErrorHandler,
		CachePeriod:  DefaultCachePeriod,
	}
	for _, opt := range options {
		opt.Apply(config)
	}
	integrator := integrator.New(selector, config.Stateful)
	accum := sdk.NewAccumulator(
		integrator,
		sdk.WithResource(config.Resource),
		sdk.WithErrorHandler(config.ErrorHandler),
	)
	return &Controller{
		accumulator: accum,
		integrator:  integrator,
		provider:    registry.NewProvider(accum),
		period:      config.CachePeriod,
		checkpoint:  integrator.CheckpointSet(),
		clock:       controllerTime.RealClock{},
	}
}

// SetClock sets the clock used for caching.  For testing purposes.
func (c *Controller) SetClock(clock controllerTime.Clock) {
	c.integrator.Lock()
	defer c.integrator.Unlock()
	c.clock = clock
}

// Provider returns a metric.Provider for the implementation managed
// by this controller.
func (c *Controller) Provider() metric.Provider {
	return c.provider
}

// Foreach gives the caller read-locked access to the current
// export.CheckpointSet.
func (c *Controller) ForEach(f func(export.Record) error) error {
	c.integrator.RLock()
	defer c.integrator.RUnlock()

	return c.checkpoint.ForEach(f)
}

// Collect requests a collection.  The collection will be skipped if
// the last collection is aged less than the CachePeriod.
func (c *Controller) Collect(ctx context.Context) {
	c.integrator.Lock()
	defer c.integrator.Unlock()

	if c.period > 0 {
		now := c.clock.Now()
		elapsed := now.Sub(c.lastCollect)

		if elapsed < c.period {
			return
		}
		c.lastCollect = now
	}

	c.accumulator.Collect(ctx)
	c.checkpoint = c.integrator.CheckpointSet()
}
