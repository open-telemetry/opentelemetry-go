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

package otlp // import "go.opentelemetry.io/otel/exporters/otlp"

import (
	"context"
	"sync"

	metricsdk "go.opentelemetry.io/otel/sdk/export/metric"
)

// ProtocolDriver is an interface used by OTLP exporter. It's
// responsible for connecting to and disconnecting from the collector,
// and for transforming traces and metrics into wire format and
// transmitting them to the collector.
type ProtocolDriver interface {
	// Start should establish connection(s) to endpoint(s). It is
	// called just once by the exporter, so the implementation
	// does not need to worry about idempotence and locking.
	Start(ctx context.Context) error
	// Stop should close the connections. The function is called
	// only once by the exporter, so the implementation does not
	// need to worry about idempotence, but it may be called
	// concurrently with ExportMetrics or ExportTraces, so proper
	// locking is required. The function serves as a
	// synchronization point - after the function returns, the
	// process of closing connections is assumed to be finished.
	Stop(ctx context.Context) error
	// ExportMetrics should transform the passed metrics to the
	// wire format and send it to the collector. May be called
	// concurrently with ExportTraces, so the manager needs to
	// take this into account by doing proper locking.
	ExportMetrics(ctx context.Context, cps metricsdk.CheckpointSet, selector metricsdk.ExportKindSelector) error
}

// SplitConfig is used to configure a split driver.
type SplitConfig struct {
	// ForMetrics driver will be used for sending metrics to the
	// collector.
	ForMetrics ProtocolDriver
}

type splitDriver struct {
	metric ProtocolDriver
}

// noopDriver implements the ProtocolDriver interface and
// is used internally to implement split drivers that do not have
// all drivers configured.
type noopDriver struct{}

var _ ProtocolDriver = (*noopDriver)(nil)

var _ ProtocolDriver = (*splitDriver)(nil)

// NewSplitDriver creates a protocol driver which may contain multiple
// protocol drivers and will forward signals to the appropriate driver.
func NewSplitDriver(opts ...SplitDriverOption) ProtocolDriver {
	driver := splitDriver{
		metric: &noopDriver{},
	}
	for _, opt := range opts {
		opt.apply(&driver)
	}
	return &driver
}

// Start implements ProtocolDriver. It starts all drivers at the same
// time.
func (d *splitDriver) Start(ctx context.Context) error {
	wg := sync.WaitGroup{}
	wg.Add(1)
	var (
		metricErr error
	)
	go func() {
		defer wg.Done()
		metricErr = d.metric.Start(ctx)
	}()
	wg.Wait()
	if metricErr != nil {
		return metricErr
	}
	return nil
}

// Stop implements ProtocolDriver. It stops all drivers at the same
// time.
func (d *splitDriver) Stop(ctx context.Context) error {
	wg := sync.WaitGroup{}
	wg.Add(1)
	var (
		metricErr error
	)
	go func() {
		defer wg.Done()
		metricErr = d.metric.Stop(ctx)
	}()
	wg.Wait()
	if metricErr != nil {
		return metricErr
	}
	return nil
}

// ExportMetrics implements ProtocolDriver. It forwards the call to
// the driver used for sending metrics.
func (d *splitDriver) ExportMetrics(ctx context.Context, cps metricsdk.CheckpointSet, selector metricsdk.ExportKindSelector) error {
	return d.metric.ExportMetrics(ctx, cps, selector)
}

// Start does nothing.
func (d *noopDriver) Start(ctx context.Context) error {
	return nil
}

// Stop does nothing.
func (d *noopDriver) Stop(ctx context.Context) error {
	return nil
}

// ExportMetrics does nothing.
func (d *noopDriver) ExportMetrics(ctx context.Context, cps metricsdk.CheckpointSet, selector metricsdk.ExportKindSelector) error {
	return nil
}
