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

package trace

// SimpleSpanProcessor implements SpanProcessor interfaces. It is used by
// exporters to receive SpanData synchronously when span is finished.
type SimpleSpanProcessor struct {
	exporter Exporter
}

var _ SpanProcessor = (*SimpleSpanProcessor)(nil)

// NewSimpleSpanProcessor creates a new instance of SimpleSpanProcessor
// for a given exporter.
func NewSimpleSpanProcessor(exporter Exporter) *SimpleSpanProcessor {
	ssp := &SimpleSpanProcessor{
		exporter: exporter,
	}
	return ssp
}

// OnStart method does nothing.
func (ssp *SimpleSpanProcessor) OnStart(sd *SpanData) {
}

// OnEnd method exports SpanData using associated exporter.
func (ssp *SimpleSpanProcessor) OnEnd(sd *SpanData) {
	if ssp.exporter != nil {
		ssp.exporter.ExportSpan(sd)
	}
}

// Shutdown method does nothing. There is no data to cleanup.
func (ssp *SimpleSpanProcessor) Shutdown() {
}
