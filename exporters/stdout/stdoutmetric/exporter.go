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
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// exporter is an OpenTelemetry metric exporter.
type exporter struct {
	encVal atomic.Value // encoderHolder

	shutdownOnce sync.Once
}

// New returns a configured metric Exporter.
//
// If no options are passed, the default exporter returned will use a JSON
// encoder with tab indentations.
func New(options ...Option) (metric.Exporter, error) {
	cfg := newConfig(options...)
	exp := &exporter{}
	exp.encVal.Store(*cfg.encoder)
	return exp, nil
}

func (e *exporter) Export(ctx context.Context, data metricdata.ResourceMetrics) error {
	select {
	case <-ctx.Done():
		// Don't do anything if the context has already timed out.
		return ctx.Err()
	default:
		// Context is still valid, continue.
	}

	return e.encVal.Load().(encoderHolder).Encode(data)
}

func (e *exporter) ForceFlush(context.Context) error {
	// exporter holds no state, nothing to flush.
	return nil
}

func (e *exporter) Shutdown(context.Context) error {
	e.shutdownOnce.Do(func() {
		e.encVal.Store(encoderHolder{
			encoder: shutdownEncoder{},
		})
	})
	return nil
}
