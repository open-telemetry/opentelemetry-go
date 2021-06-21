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

package otlpmetric

import (
	"context"
	"errors"
	"sync"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/internal/metrictransform"
	"go.opentelemetry.io/otel/metric"
	metricsdk "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
)

var (
	errAlreadyStarted = errors.New("already started")
)

// Exporter exports metrics data in the OTLP wire format.
type Exporter struct {
	client             Client
	exportKindSelector metricsdk.ExportKindSelector

	mu      sync.RWMutex
	started bool

	startOnce sync.Once
	stopOnce  sync.Once
}

// Export exports a batch of metrics.
func (e *Exporter) Export(ctx context.Context, checkpointSet metricsdk.CheckpointSet) error {
	rms, err := metrictransform.CheckpointSet(ctx, e, checkpointSet, 1)
	if err != nil {
		return err
	}
	if len(rms) == 0 {
		return nil
	}

	return e.client.UploadMetrics(ctx, rms)
}

// Start establishes a connection to the receiving endpoint.
func (e *Exporter) Start(ctx context.Context) error {
	var err = errAlreadyStarted
	e.startOnce.Do(func() {
		e.mu.Lock()
		e.started = true
		e.mu.Unlock()
		err = e.client.Start(ctx)
	})

	return err
}

// Shutdown flushes all exports and closes all connections to the receiving endpoint.
func (e *Exporter) Shutdown(ctx context.Context) error {

	e.mu.RLock()
	started := e.started
	e.mu.RUnlock()

	if !started {
		return nil
	}

	var err error

	e.stopOnce.Do(func() {
		err = e.client.Stop(ctx)
		e.mu.Lock()
		e.started = false
		e.mu.Unlock()
	})

	return err
}

func (e *Exporter) ExportKindFor(descriptor *metric.Descriptor, aggregatorKind aggregation.Kind) metricsdk.ExportKind {
	return e.exportKindSelector.ExportKindFor(descriptor, aggregatorKind)
}

var _ metricsdk.Exporter = (*Exporter)(nil)

// New constructs a new Exporter and starts it.
func New(ctx context.Context, client Client, opts ...Option) (*Exporter, error) {
	exp := NewUnstarted(client, opts...)
	if err := exp.Start(ctx); err != nil {
		return nil, err
	}
	return exp, nil
}

// NewUnstarted constructs a new Exporter and does not start it.
func NewUnstarted(client Client, opts ...Option) *Exporter {
	cfg := config{
		// Note: the default ExportKindSelector is specified
		// as Cumulative:
		// https://github.com/open-telemetry/opentelemetry-specification/issues/731
		exportKindSelector: metricsdk.CumulativeExportKindSelector(),
	}

	for _, opt := range opts {
		opt.apply(&cfg)
	}

	e := &Exporter{
		client:             client,
		exportKindSelector: cfg.exportKindSelector,
	}

	return e
}
