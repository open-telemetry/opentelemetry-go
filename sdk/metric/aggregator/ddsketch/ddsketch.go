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
package ddsketch // import "go.opentelemetry.io/otel/sdk/metric/aggregator/ddsketch"

import (
	"context"
	"math"
	"sync"

	sdk "github.com/DataDog/sketches-go/ddsketch"

	"go.opentelemetry.io/otel/api/metric"

	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
)

// Config is an alias for the underlying DDSketch config object.
type Config = sdk.Config

// Aggregator aggregates events into a distribution.
type Aggregator struct {
	lock       sync.Mutex
	cfg        *Config
	kind       metric.NumberKind
	current    *sdk.DDSketch
	checkpoint *sdk.DDSketch
}

var _ export.Aggregator = &Aggregator{}
var _ aggregator.MinMaxSumCount = &Aggregator{}
var _ aggregator.Distribution = &Aggregator{}

// New returns a new DDSketch aggregator.
func New(cfg *Config, desc *metric.Descriptor) *Aggregator {
	return &Aggregator{
		cfg:        cfg,
		kind:       desc.NumberKind(),
		current:    sdk.NewDDSketch(cfg),
		checkpoint: sdk.NewDDSketch(cfg),
	}
}

// NewDefaultConfig returns a new, default DDSketch config.
//
// TODO: Should the Config constructor set minValue to -Inf to
// when the descriptor has absolute=false?  This requires providing
// values for alpha and maxNumBins, apparently.
func NewDefaultConfig() *Config {
	return sdk.NewDefaultConfig()
}

// Sum returns the sum of values in the checkpoint.
func (c *Aggregator) Sum() (metric.Number, error) {
	return c.toNumber(c.checkpoint.Sum()), nil
}

// Count returns the number of values in the checkpoint.
func (c *Aggregator) Count() (int64, error) {
	return c.checkpoint.Count(), nil
}

// Max returns the maximum value in the checkpoint.
func (c *Aggregator) Max() (metric.Number, error) {
	return c.Quantile(1)
}

// Min returns the minimum value in the checkpoint.
func (c *Aggregator) Min() (metric.Number, error) {
	return c.Quantile(0)
}

// Quantile returns the estimated quantile of data in the checkpoint.
// It is an error if `q` is less than 0 or greated than 1.
func (c *Aggregator) Quantile(q float64) (metric.Number, error) {
	if c.checkpoint.Count() == 0 {
		return metric.Number(0), aggregator.ErrNoData
	}
	f := c.checkpoint.Quantile(q)
	if math.IsNaN(f) {
		return metric.Number(0), aggregator.ErrInvalidQuantile
	}
	return c.toNumber(f), nil
}

func (c *Aggregator) toNumber(f float64) metric.Number {
	if c.kind == metric.Float64NumberKind {
		return metric.NewFloat64Number(f)
	}
	return metric.NewInt64Number(int64(f))
}

// Checkpoint saves the current state and resets the current state to
// the empty set, taking a lock to prevent concurrent Update() calls.
func (c *Aggregator) Checkpoint(ctx context.Context, _ *metric.Descriptor) {
	replace := sdk.NewDDSketch(c.cfg)

	c.lock.Lock()
	c.checkpoint = c.current
	c.current = replace
	c.lock.Unlock()
}

// Update adds the recorded measurement to the current data set.
// Update takes a lock to prevent concurrent Update() and Checkpoint()
// calls.
func (c *Aggregator) Update(_ context.Context, number metric.Number, desc *metric.Descriptor) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.current.Add(number.CoerceToFloat64(desc.NumberKind()))
	return nil
}

// Merge combines two sketches into one.
func (c *Aggregator) Merge(oa export.Aggregator, d *metric.Descriptor) error {
	o, _ := oa.(*Aggregator)
	if o == nil {
		return aggregator.NewInconsistentMergeError(c, oa)
	}

	c.checkpoint.Merge(o.checkpoint)
	return nil
}
