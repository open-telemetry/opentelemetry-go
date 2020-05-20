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
	integrator "go.opentelemetry.io/otel/sdk/metric/integrator/simple"
	"go.opentelemetry.io/otel/sdk/resource"
)

// DefaultCachePeriod determines how long a recently-computed result
// will be returned without gathering metric data again.  If the period
// is zero, caching of the result is disabled.
const DefaultCachePeriod time.Duration = 0

type Controller struct {
	accumulator *sdk.Accumulator
	integrator  *integrator.Integrator
	provider    *registry.Provider
	period      time.Duration
	lastCollect time.Time
	checkpoint  export.CheckpointSet
}

func New(selector export.AggregationSelector, options ...Option) *Controller {
	config := &Config{
		Resource:     resource.Empty(),
		ErrorHandler: sdk.DefaultErrorHandler,
		CachePeriod:  DefaultCachePeriod,
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
	}
}

func (c *Controller) Stop() error {
	*c = Controller{}
	return nil
}

func (c *Controller) Provider() metric.Provider {
	return c.provider
}

func (c *Controller) ForEach(f func(export.Record) error) error {
	c.integrator.RLock()
	defer c.integrator.RUnlock()

	return c.checkpoint.ForEach(f)
}

func (c *Controller) Collect(ctx context.Context) {
	c.integrator.Lock()
	defer c.integrator.Unlock()

	if c.period > 0 {
		now := time.Now()
		elapsed := now.Sub(c.lastCollect)

		if elapsed < c.period {
			return
		}
		c.lastCollect = now
	}

	c.accumulator.Collect(ctx)
	c.checkpoint = c.integrator.CheckpointSet()
}
