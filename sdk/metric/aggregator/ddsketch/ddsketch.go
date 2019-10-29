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
package ddsketch

import (
	"context"
	"sync"

	sdk "github.com/DataDog/sketches-go/ddsketch"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/sdk/export"
)

type (
	// Aggregator aggregates measure events.
	Aggregator struct {
		lock sync.Mutex
		cfg  *sdk.Config
		kind core.NumberKind
		live *sdk.DDSketch
		save *sdk.DDSketch
	}
)

var _ export.MetricAggregator = &Aggregator{}

// New returns a new DDSketch aggregator.
func New(cfg *sdk.Config, desc *export.Descriptor) *Aggregator {
	return &Aggregator{
		cfg:  cfg,
		kind: desc.NumberKind(),
		live: sdk.NewDDSketch(cfg),
	}
}

func NewDefaultConfig() *sdk.Config {
	return sdk.NewDefaultConfig()
}

func (c *Aggregator) Sum() float64 {
	return c.save.Sum()
}

func (c *Aggregator) Count() int64 {
	return c.save.Count()
}

func (c *Aggregator) Max() float64 {
	return c.save.Quantile(1)
}

func (c *Aggregator) Min() float64 {
	return c.save.Quantile(0)
}

func (c *Aggregator) Quantile(q float64) float64 {
	return c.save.Quantile(q)
}

// Collect saves the current value (atomically) and exports it.
func (c *Aggregator) Collect(ctx context.Context, rec export.MetricRecord, exp export.MetricBatcher) {
	replace := sdk.NewDDSketch(c.cfg)

	c.lock.Lock()
	c.save = c.live
	c.live = replace
	c.lock.Unlock()

	if c.save.Count() != 0 {
		exp.Export(ctx, rec, c)
	}
}

// Collect updates the current value (atomically) for later export.
func (c *Aggregator) Update(_ context.Context, number core.Number, rec export.MetricRecord) {
	desc := rec.Descriptor()
	kind := desc.NumberKind()

	if !desc.Alternate() && number.IsNegative(kind) {
		// TODO warn
		return
	}

	c.lock.Lock()
	defer c.lock.Unlock()
	c.live.Add(number.CoerceToFloat64(kind))
}
