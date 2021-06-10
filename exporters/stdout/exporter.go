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

package stdout // import "go.opentelemetry.io/otel/exporters/stdout"

import (
	"go.opentelemetry.io/otel/sdk/export/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Exporter struct {
	traceExporter
	metricExporter
}

var (
	_ metric.Exporter       = &Exporter{}
	_ sdktrace.SpanExporter = &Exporter{}
)

// New creates an Exporter with the passed options.
func New(options ...Option) (*Exporter, error) {
	cfg, err := newConfig(options...)
	if err != nil {
		return nil, err
	}
	return &Exporter{
		traceExporter:  traceExporter{config: cfg},
		metricExporter: metricExporter{cfg},
	}, nil
}
