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

package format

import (
	"strings"

	"go.opentelemetry.io/experimental/streaming/exporter/reader/format"
	"go.opentelemetry.io/experimental/streaming/exporter/spandata"
)

func AppendSpan(buf *strings.Builder, data *spandata.Span) {
	for _, event := range data.Events {
		format.AppendEvent(buf, event)
	}
}

func SpanToString(data *spandata.Span) string {
	var buf strings.Builder
	AppendSpan(&buf, data)
	return buf.String()
}
