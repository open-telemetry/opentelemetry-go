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

package dogstatsd // import "go.opentelemetry.io/otel/exporter/metric/dogstatsd"

import (
	"bytes"

	"go.opentelemetry.io/otel/exporter/metric/internal/statsd"
	export "go.opentelemetry.io/otel/sdk/export/metric"
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
		*statsd.LabelEncoder

		ReencodedLabelsCount int
	}
)

var (
	_ export.Exporter     = &Exporter{}
	_ export.LabelEncoder = &Exporter{}
)

// New returns a new Dogstatsd-syntax exporter.  This type implements
// the metric.LabelEncoder interface, allowing the SDK's unique label
// encoding to be pre-computed for the exporter and stored in the
// LabelSet.
func New(config Config) (*Exporter, error) {
	exp := &Exporter{
		LabelEncoder: statsd.NewLabelEncoder(),
	}

	var err error
	exp.Exporter, err = statsd.NewExporter(config, exp)
	return exp, err
}

// AppendName is part of the stats-internal adapter interface.
func (*Exporter) AppendName(rec export.Record, buf *bytes.Buffer) {
	_, _ = buf.WriteString(rec.Descriptor().Name())
}

// AppendTags is part of the stats-internal adapter interface.
func (e *Exporter) AppendTags(rec export.Record, buf *bytes.Buffer) {
	encoded, inefficient := e.LabelEncoder.ForceEncode(rec.Labels())
	_, _ = buf.WriteString(encoded)

	if inefficient {
		e.ReencodedLabelsCount++
	}
}
