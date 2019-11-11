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

	Exporter struct {
		*statsd.Exporter
		*statsd.LabelEncoder
	}
)

var (
	_ export.Exporter     = &Exporter{}
	_ export.LabelEncoder = &Exporter{}
)

func New(config Config) (*Exporter, error) {
	exp := &Exporter{
		LabelEncoder: statsd.NewLabelEncoder(),
	}

	var err error
	exp.Exporter, err = statsd.NewExporter(config, exp)
	return exp, err
}

func (*Exporter) AppendName(rec export.Record, buf *bytes.Buffer) {
	_, _ = buf.WriteString(rec.Descriptor().Name())
}

func (e *Exporter) AppendTags(rec export.Record, buf *bytes.Buffer) {
	labels := rec.Labels()

	if labels.Encoder() != e {
		// TODO: This case could be handled by directly
		// encoding the labels at this point, but presently it
		// should not occur.
		panic("Should have self-encoded labels")
	}

	_, _ = buf.WriteString(labels.Encoded())
}
