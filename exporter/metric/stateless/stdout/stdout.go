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

package stdout

import (
	"context"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/sdk/export"
)

type Exporter struct {
}

var _ export.MetricExporter = &Exporter{}

func New() *Exporter {
	return &Exporter{}
}

func (*Exporter) Export(_ context.Context, producer export.MetricProducer) {
	producer.Foreach(func(agg export.MetricAggregator, desc *export.Descriptor, labels []core.KeyValue) {
		// fmt.Printf("%s %s\n",
	})
}
