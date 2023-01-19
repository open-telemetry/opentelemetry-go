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

package stdoutlog // import "go.opentelemetry.io/otel/exporters/stdout/stdoutlog"

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/log/logtest"
)

var zeroTime time.Time

var _ log.LogRecordExporter = &Exporter{}

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
		encoder:    enc,
		timestamps: cfg.Timestamps,
	}, nil
}

// Exporter is an implementation of trace.SpanSyncer that writes spans to stdout.
type Exporter struct {
	encoder    *json.Encoder
	encoderMu  sync.Mutex
	timestamps bool

	stoppedMu sync.RWMutex
	stopped   bool
}

// ExportLogRecords writes spans in json format to stdout.
func (e *Exporter) ExportLogRecords(ctx context.Context, spans []log.ReadOnlyLogRecord) error {
	e.stoppedMu.RLock()
	stopped := e.stopped
	e.stoppedMu.RUnlock()
	if stopped {
		return nil
	}

	if len(spans) == 0 {
		return nil
	}

	stubs := logtest.LogRecordStubsFromReadOnlyLogRecords(spans)

	e.encoderMu.Lock()
	defer e.encoderMu.Unlock()
	for i := range stubs {
		stub := &stubs[i]
		// Remove timestamps
		if !e.timestamps {
			stub.Timestamp = zeroTime
		}

		// Encode span stubs, one by one
		if err := e.encoder.Encode(stub); err != nil {
			return err
		}
	}
	return nil
}

// Shutdown is called to stop the exporter, it preforms no action.
func (e *Exporter) Shutdown(ctx context.Context) error {
	e.stoppedMu.Lock()
	e.stopped = true
	e.stoppedMu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return nil
}

// MarshalLog is the marshaling function used by the logging system to represent this exporter.
func (e *Exporter) MarshalLog() interface{} {
	return struct {
		Type           string
		WithTimestamps bool
	}{
		Type:           "stdout",
		WithTimestamps: e.timestamps,
	}
}
