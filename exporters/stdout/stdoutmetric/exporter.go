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

package stdoutmetric // import "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"

import (
	"context"
	"encoding/json"
	"sync"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// New creates an Exporter with the passed options.
func New(options ...Option) (*Exporter, error) {
	cfg, err := newConfig(options...)
	if err != nil {
		return nil, err
	}

	enc := json.NewEncoder(cfg.Writer)
	if cfg.PrettyPrint {
		enc.SetIndent("", "\t")
	}

	return &Exporter{
		encoder: enc,
	}, nil
}

// Exporter is an implementation of metric.Exporter that writes ResourceMetrics to stdout.
type Exporter struct {
	encoder   *json.Encoder
	encoderMu sync.Mutex

	stoppedMu sync.RWMutex
	stopped   bool
}

var _ metric.Exporter = &Exporter{}

// Export serializes and transmits metric data to a receiver.
//
// This is called synchronously, there is no concurrency safety
// requirement. Because of this, it is critical that all timeouts and
// cancellations of the passed context be honored.
//
// All retry logic must be contained in this function. The SDK does not
// implement any retry logic. All errors returned by this function are
// considered unrecoverable and will be reported to a configured error
// Handler.
func (e *Exporter) Export(_ context.Context, data metricdata.ResourceMetrics) error {
	e.stoppedMu.RLock()
	stopped := e.stopped
	e.stoppedMu.RUnlock()
	if stopped {
		return nil
	}

	e.encoderMu.Lock()
	defer e.encoderMu.Unlock()

	return e.encoder.Encode(data)
}

// ForceFlush flushes any metric data held by an exporter.
//
// The deadline or cancellation of the passed context must be honored. An
// appropriate error should be returned in these situations.
func (e *Exporter) ForceFlush(_ context.Context) error {
	// This exporter doesn't hold any data.
	return nil
}

// Shutdown flushes all metric data held by an exporter and releases any
// held computational resources.
//
// The deadline or cancellation of the passed context must be honored. An
// appropriate error should be returned in these situations.
//
// After Shutdown is called, calls to Export will perform no operation and
// instead will return an error indicating the shutdown state.
func (e *Exporter) Shutdown(ctx context.Context) error {
	e.stoppedMu.Lock()
	e.stopped = true
	e.stoppedMu.Unlock()

	return ctx.Err()
}
