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
	"go.opentelemetry.io/otel/sdk/metric/aggregator"

	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
)

// Config is an alias for the underlying DDSketch config object.
type Config = sdk.Config

type sketchValue struct {
	sketch *sdk.DDSketch
	self   *Aggregator
}

// Aggregator aggregates events into a distribution.
type Aggregator struct {
	lock       sync.Mutex
	cfg        *Config
	kind       metric.NumberKind
	current    sketchValue
	checkpoint sketchValue
}

var _ export.Aggregator = &Aggregator{}
var _ aggregation.MinMaxSumCount = &sketchValue{}
var _ aggregation.Distribution = &sketchValue{}

// New returns a new DDSketch aggregation.
func New(desc *metric.Descriptor, cfg *Config) *Aggregator {
	if cfg == nil {
		cfg = NewDefaultConfig()
	}
	agg := &Aggregator{
		cfg:  cfg,
		kind: desc.NumberKind(),
	}
	agg.current = sketchValue{
		sketch: sdk.NewDDSketch(cfg),
		self:   agg,
	}
	agg.checkpoint = sketchValue{
		sketch: sdk.NewDDSketch(cfg),
		self:   agg,
	}
	return agg
}

// NewDefaultConfig returns a new, default DDSketch config.
//
// TODO: Should the Config constructor set minValue to -Inf to
// when the descriptor has absolute=false?  This requires providing
// values for alpha and maxNumBins, apparently.
func NewDefaultConfig() *Config {
	return sdk.NewDefaultConfig()
}

// Checkpoint saves the current state and resets the current state to
// the empty set, taking a lock to prevent concurrent Update() calls.
func (c *Aggregator) Checkpoint(_ *metric.Descriptor) {
	replace := sdk.NewDDSketch(c.cfg)

	c.lock.Lock()
	c.checkpoint.sketch = c.current.sketch
	c.current.sketch = replace
	c.lock.Unlock()
}

// Update adds the recorded measurement to the current data set.
// Update takes a lock to prevent concurrent Update() and Checkpoint()
// calls.
func (c *Aggregator) Update(_ context.Context, number metric.Number, desc *metric.Descriptor) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.current.sketch.Add(number.CoerceToFloat64(desc.NumberKind()))
	return nil
}

// Merge combines two sketches into one.
func (c *Aggregator) Merge(oa export.Aggregator, d *metric.Descriptor) error {
	o, _ := oa.(*Aggregator)
	if o == nil {
		return aggregator.NewInconsistentMergeError(c, oa)
	}

	c.current.sketch.Merge(o.checkpoint.sketch)
	return nil
}

func (c *Aggregator) Swap() {
	c.current, c.checkpoint = c.checkpoint, c.current
}

func (c *Aggregator) CheckpointedValue() aggregation.Aggregation {
	return &c.checkpoint
}

func (c *Aggregator) AccumulatedValue() aggregation.Aggregation {
	return &c.current
}

func (c *Aggregator) toNumber(f float64) metric.Number {
	if c.kind == metric.Float64NumberKind {
		return metric.NewFloat64Number(f)
	}
	return metric.NewInt64Number(int64(f))
}

// Kind returns aggregation.SketchKind
func (s *sketchValue) Kind() aggregation.Kind {
	return aggregation.SketchKind
}

// Sum returns the sum of values in the checkpoint.
func (s *sketchValue) Sum() (metric.Number, error) {
	return s.self.toNumber(s.sketch.Sum()), nil
}

// Count returns the number of values in the checkpoint.
func (s *sketchValue) Count() (int64, error) {
	return s.sketch.Count(), nil
}

// Max returns the maximum value in the checkpoint.
func (s *sketchValue) Max() (metric.Number, error) {
	return s.Quantile(1)
}

// Min returns the minimum value in the checkpoint.
func (s *sketchValue) Min() (metric.Number, error) {
	return s.Quantile(0)
}

// Quantile returns the estimated quantile of data in the checkpoint.
// It is an error if `q` is less than 0 or greated than 1.
func (s *sketchValue) Quantile(q float64) (metric.Number, error) {
	if s.sketch.Count() == 0 {
		return metric.Number(0), aggregation.ErrNoData
	}
	f := s.sketch.Quantile(q)
	if math.IsNaN(f) {
		return metric.Number(0), aggregation.ErrInvalidQuantile
	}
	return s.self.toNumber(f), nil
}
