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
	"sync"
	"sync/atomic"
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
	ExportSpan(context.Context, interface{})
}

// BatchExporter is a type for functions that receive sampled trace spans.
//
// The ExportSpans method is called asynchronously. However BatchExporter should
// not take forever to process the spans.
//
// The SpanData should not be modified.
type BatchExporter interface {
	ExportSpans(context.Context, []interface{})
}

type exportersMap map[Type][]interface{}

var (
	exporterMu sync.Mutex
	exporters  atomic.Value
)

// Register adds to the list of Exporters that will receive sampled trace
// spans.
func Register(t Type, e interface{}) {
	exporterMu.Lock()
	defer exporterMu.Unlock()

	nm := make(exportersMap)
	if old, ok := exporters.Load().(exportersMap); ok {
		for k, v := range old {
			nm[k] = v
		}
	}

	nm[t] = append(nm[t], e)
	exporters.Store(nm)
}

// Unregister removes from the list of Exporters the Exporter that was
// registered with the given name.
func Unregister(e interface{}) {
	exporterMu.Lock()
	defer exporterMu.Unlock()

	nm := make(exportersMap)
	if old, ok := exporters.Load().(exportersMap); ok {
		for k, tv := range old {
			for i, v := range tv {
				if e == v {
					tv[i] = tv[len(tv)-1]
					tv = tv[:len(tv)-1]
				}
			}

			nm[k] = tv
		}
	}

	exporters.Store(nm)
}

// Load loads all the registered exporters matching a specific type
func Load(t Type) []interface{} {
	exporterMu.Lock()
	defer exporterMu.Unlock()

	e, _ := exporters.Load().(exportersMap)
	return e[t]
}
