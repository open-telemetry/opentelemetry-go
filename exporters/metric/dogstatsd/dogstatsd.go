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

package dogstatsd // import "go.opentelemetry.io/otel/exporters/metric/dogstatsd"

import (
	"bytes"
	"time"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/metric/internal/statsd"

	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/batcher/ungrouped"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

type (
	Config = statsd.Config

	// Exporter implements a dogstatsd-format statsd exporter,
	// which encodes label sets as independent fields in the
	// output.
	//
	// TODO: find a link for this syntax.  It's been copied out of
	// code, not a specification:
	//
	// https://github.com/stripe/veneur/blob/master/sinks/datadog/datadog.go
	Exporter struct {
		*statsd.Exporter

		labelEncoder *statsd.LabelEncoder
	}
)

var (
	_ export.Exporter = &Exporter{}
)

// NewRawExporter returns a new Dogstatsd-syntax exporter for use in a pipeline.
func NewRawExporter(config Config) (*Exporter, error) {
	exp := &Exporter{
		labelEncoder: statsd.NewLabelEncoder(),
	}

	var err error
	exp.Exporter, err = statsd.NewExporter(config, exp)
	return exp, err
}

// InstallNewPipeline instantiates a NewExportPipeline and registers it globally.
// Typically called as:
//
// 	pipeline, err := dogstatsd.InstallNewPipeline(dogstatsd.Config{...})
// 	if err != nil {
// 		...
// 	}
// 	defer pipeline.Stop()
// 	... Done
func InstallNewPipeline(config Config) (*push.Controller, error) {
	controller, err := NewExportPipeline(config, time.Minute)
	if err != nil {
		return controller, err
	}
	global.SetMeterProvider(controller)
	return controller, err
}

// NewExportPipeline sets up a complete export pipeline with the recommended setup,
// chaining a NewRawExporter into the recommended selectors and batchers.
func NewExportPipeline(config Config, period time.Duration) (*push.Controller, error) {
	selector := simple.NewWithExactMeasure()
	exporter, err := NewRawExporter(config)
	if err != nil {
		return nil, err
	}

	// The ungrouped batcher ensures that the export sees the full
	// set of labels as dogstatsd tags.
	batcher := ungrouped.New(selector, exporter.labelEncoder, false)

	pusher := push.New(batcher, exporter, period)
	pusher.Start()

	return pusher, nil
}

// AppendName is part of the stats-internal adapter interface.
func (*Exporter) AppendName(rec export.Record, buf *bytes.Buffer) {
	_, _ = buf.WriteString(rec.Descriptor().Name())
}

// AppendTags is part of the stats-internal adapter interface.
func (e *Exporter) AppendTags(rec export.Record, buf *bytes.Buffer) {
	encoded := rec.Labels().Encoded(e.labelEncoder)
	_, _ = buf.WriteString(encoded)
}
