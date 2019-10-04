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

package exporter

import (
	"context"
)

// Type defines all the types of exporters we have available
type Type int32

// All the exporter types
const (
	ExporterTypeSync = iota
	ExporterTypeBatch
)

// SyncExporter is a synchronous exporter
type SyncExporter interface {
	ExportSpan(context.Context, *SpanData)
}

// BatchExporter is a type for functions that receive sampled trace spans.
//
// The ExportSpans method is called asynchronously. However BatchExporter should
// not take forever to process the spans.
//
// The SpanData should not be modified.
type BatchExporter interface {
	ExportSpans(context.Context, []*SpanData)
}
