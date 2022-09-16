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

//go:build go1.18
// +build go1.18

package otest // import "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/internal/otest"

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	cpb "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	mpb "go.opentelemetry.io/proto/otlp/metrics/v1"
)

type client struct {
	storage *Storage
}

func (c *client) Collect() *Storage {
	return c.storage
}

func (c *client) UploadMetrics(ctx context.Context, rm *mpb.ResourceMetrics) error {
	c.storage.Add(&cpb.ExportMetricsServiceRequest{
		ResourceMetrics: []*mpb.ResourceMetrics{rm},
	})
	return ctx.Err()
}

func (c *client) ForceFlush(ctx context.Context) error { return ctx.Err() }
func (c *client) Shutdown(ctx context.Context) error   { return ctx.Err() }

func TestClientTests(t *testing.T) {
	factory := func() (otlpmetric.Client, Collector) {
		c := &client{storage: NewStorage()}
		return c, c
	}

	t.Run("Integration", RunClientTests(factory))
}
