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

package push

import (
	"context"
	"time"

	"go.opentelemetry.io/api/metric"
	"go.opentelemetry.io/sdk/export"
	sdk "go.opentelemetry.io/sdk/metric"
)

// Controller organizes a periodic push of metric data.
type Controller struct {
	sdk      *sdk.SDK
	batcher  export.MetricBatcher
	exporter export.MetricExporter
	ticker   *time.Ticker
	ch       chan struct{}
}

var _ metric.Provider = &Controller{}

// New constructs a Controller, an implementation of metric.Provider,
// using the provider batcher, exporter, period.  The batcher itself
// is configured with aggregation policy selection.
func New(batcher export.MetricBatcher, exporter export.MetricExporter, period time.Duration) *Controller {
	return &Controller{
		sdk:      sdk.New(batcher),
		batcher:  batcher,
		exporter: exporter,
		ticker:   time.NewTicker(period),
		ch:       make(chan struct{}),
	}
}

func (c *Controller) GetMeter(name string) metric.Meter {
	return c.sdk
}

func (c *Controller) Start() {
	go c.run()
}

func (c *Controller) Stop() {
	close(c.ch)

	// TODO wait for the last run, flush, etc.
}

func (c *Controller) run() {
	for {
		select {
		case <-c.ch:
			return
		case <-c.ticker.C:
			c.tick()
		}
	}
}

func (c *Controller) tick() {
	ctx := context.Background()
	c.sdk.Collect(ctx)
	c.exporter.Export(ctx, c.batcher.ReadCheckpoint())
}
