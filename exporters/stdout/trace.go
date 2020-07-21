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

package stdout

import (
	"context"
	"encoding/json"

	"go.opentelemetry.io/otel/sdk/export/trace"
)

// Exporter is an implementation of trace.SpanSyncer that writes spans to stdout.
type traceExporter struct {
	config Config
}

// ExportSpan writes a SpanData in json format to stdout.
func (e *traceExporter) ExportSpan(ctx context.Context, data *trace.SpanData) {
	var jsonSpan []byte
	var err error
	if e.config.PrettyPrint {
		jsonSpan, err = json.MarshalIndent(data, "", "\t")
	} else {
		jsonSpan, err = json.Marshal(data)
	}
	if err != nil {
		// ignore writer failures for now
		_, _ = e.config.Writer.Write([]byte("Error converting spanData to json: " + err.Error()))
		return
	}
	// ignore writer failures for now
	_, _ = e.config.Writer.Write(append(jsonSpan, byte('\n')))
}

// ExportSpans writes SpanData in json format to stdout.
func (e *traceExporter) ExportSpans(ctx context.Context, data []*trace.SpanData) {
	for _, sd := range data {
		e.ExportSpan(ctx, sd)
	}
}
