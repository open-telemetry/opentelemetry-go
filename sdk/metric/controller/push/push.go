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
	sdk      *sdk.SDK
	batcher  export.Batcher
	exporter export.Exporter
	ticker   *time.Ticker
	wg       sync.WaitGroup
	ch       chan struct{}
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
		sdk:      sdk.New(batcher, lencoder),
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
	c.wg.Add(1)
	go c.run()
}

func (c *Controller) Stop() {
	close(c.ch)
	c.wg.Wait()

	c.tick()
}

func (c *Controller) run() {
	for {
		select {
		case <-c.ch:
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
		// TODO: report this error
		_ = err
	}
}
