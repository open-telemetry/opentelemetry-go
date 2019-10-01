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

package sdk

import (
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/api/metric"
	"go.opentelemetry.io/api/trace"
	"go.opentelemetry.io/experimental/streaming/exporter"
)

type observerData struct {
	observer metric.Observer
	callback metric.ObserverCallback
}

type observersMap map[metric.DescriptorID]observerData

type sdk struct {
	exporter  *exporter.Exporter
	resources exporter.EventID

	observersLock sync.Mutex
	observers     atomic.Value // observersMap
}

type SDK interface {
	trace.Tracer
	metric.Meter
}

var _ SDK = &sdk{}

func New(observers ...exporter.Observer) SDK {
	return &sdk{
		exporter: exporter.NewExporter(observers...),
	}
}
