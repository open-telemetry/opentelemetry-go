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
)

// DefaultCachePeriod determines how long a recently-computed result
// will be returned without gathering metric data again.
const DefaultCachePeriod = 10 * time.Second

type Controller struct {
	accumulator *sdk.Accumulator
	integrator  *integrator.Integrator
	provider    *registry.Provider
}

// TODO: Options: cached result period, resource, stateful, error handler.

func New(selector export.AggregationSelector, options ...Option) *Controller {
	config := &Config{
		ErrorHandler: sdk.DefaultErrorHandler,
		CachePeriod:  DefaultCachePeriod,
	}
	integrator := integrator.New(selector, config.Stateful)
	accum := sdk.NewAccumulator(integrator)
	return &Controller{
		accumulator: accum,
		integrator:  integrator,
		provider:    registry.NewProvider(accum),
	}
}

func (c *Controller) Provider() metric.Provider {
	return c.provider
}

func (c *Controller) ForEach(f func(export.Record) error) error {
	c.integrator.RLock()
	defer c.integrator.RUnlock()

	return c.integrator.CheckpointSet().ForEach(f)
}

func (c *Controller) Collect(ctx context.Context) {
	c.integrator.Lock()
	defer c.integrator.Unlock()

	c.accumulator.Collect(ctx)
}
