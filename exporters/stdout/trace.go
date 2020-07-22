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
	"fmt"

	"go.opentelemetry.io/otel/sdk/export/trace"
)

// Exporter is an implementation of trace.SpanSyncer that writes spans to stdout.
type traceExporter struct {
	config Config
}

// ExportSpan writes a SpanData in json format to stdout.
func (e *traceExporter) ExportSpan(ctx context.Context, data *trace.SpanData) {
	if e.config.DisableTraceExport {
		return
	}
	e.ExportSpans(ctx, []*trace.SpanData{data})
}

// ExportSpans writes SpanData in json format to stdout.
func (e *traceExporter) ExportSpans(ctx context.Context, data []*trace.SpanData) {
	if e.config.DisableTraceExport || len(data) == 0 {
		return
	}
	out, err := e.marshal(data)
	if err != nil {
		fmt.Fprintf(e.config.Writer, "error converting spanData to json: %v", err)
		return

	}
	fmt.Fprintln(e.config.Writer, string(out))
}

// marshal v with approriate indentation.
func (e *traceExporter) marshal(v interface{}) ([]byte, error) {
	if e.config.PrettyPrint {
		return json.MarshalIndent(v, "", "\t")
	}
	return json.Marshal(v)
}
