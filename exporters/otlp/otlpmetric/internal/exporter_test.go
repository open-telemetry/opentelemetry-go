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

package internal // import "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/internal"

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	mpb "go.opentelemetry.io/proto/otlp/metrics/v1"
)

type client struct {
	// n is incremented by all Client methods. If these methods are called
	// concurrently this should fail tests run with the race detector.
	n int
}

func (c *client) Temporality(k metric.InstrumentKind) metricdata.Temporality {
	return metric.DefaultTemporalitySelector(k)
}

func (c *client) Aggregation(k metric.InstrumentKind) aggregation.Aggregation {
	return metric.DefaultAggregationSelector(k)
}

func (c *client) UploadMetrics(context.Context, *mpb.ResourceMetrics) error {
	c.n++
	return nil
}

func (c *client) ForceFlush(context.Context) error {
	c.n++
	return nil
}

func (c *client) Shutdown(context.Context) error {
	c.n++
	return nil
}

func TestExporterClientConcurrentSafe(t *testing.T) {
	const goroutines = 5

	exp := New(&client{})
	rm := new(metricdata.ResourceMetrics)
	ctx := context.Background()

	done := make(chan struct{})
	first := make(chan struct{}, goroutines)
	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			assert.NoError(t, exp.Export(ctx, rm))
			assert.NoError(t, exp.ForceFlush(ctx))
			// Ensure some work is done before shutting down.
			first <- struct{}{}

			for {
				_ = exp.Export(ctx, rm)
				_ = exp.ForceFlush(ctx)

				select {
				case <-done:
					return
				default:
				}
			}
		}()
	}

	for i := 0; i < goroutines; i++ {
		<-first
	}
	close(first)
	assert.NoError(t, exp.Shutdown(ctx))
	assert.ErrorIs(t, exp.Shutdown(ctx), errShutdown)

	close(done)
	wg.Wait()
}
