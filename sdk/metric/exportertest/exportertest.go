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

package exportertest

import (
	"go.opentelemetry.io/otel/api/label"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	processortest "go.opentelemetry.io/otel/sdk/metric/processor/test"
)

type (
	// Exporter is a testing implementation of export.Exporter that
	// assembles its results as a map[string]float64.
	Exporter struct {
		export.ExportKindSelector
	}
)

// NewExporter returns a new testing Exporter implementation.
// Verify exporter outputs using Values(), e.g.,:
//
//     require.EqualValues(t, map[string]float64{
//         "counter.sum/A=1,B=2/R=V": 100,
//     }, exporter.Values())
//
// Where in the example A=1,B=2 is the encoded labels and R=V is the
// encoded resource value.
func NewExporter(selector export.ExportKindSelector) *Exporter {
	return &Exporter{
		ExportKindSelector: selector,
	}
}

// Values returns the mapping from label set to point values for the
// accumulations that were processed.  Point values are chosen as
// either the Sum or the LastValue, whichever is implemented.  (All
// the built-in Aggregators implement one of these interfaces.)
func (e *Exporter) Values(ckpt export.CheckpointSet, enc label.Encoder) map[string]float64 {
	output := processortest.NewOutput(enc)
	err := ckpt.ForEach(e.ExportKindSelector, func(r export.Record) error {
		return output.AddRecord(r)
	})
	if err != nil {
		panic("Error in Values()")
	}
	return output.Map()
}
