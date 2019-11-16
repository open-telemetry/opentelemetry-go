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

package statsd

import (
	"bytes"
	"sync"

	"go.opentelemetry.io/otel/api/core"
	export "go.opentelemetry.io/otel/sdk/export/metric"
)

type LabelEncoder struct {
	pool sync.Pool
}

var _ export.LabelEncoder = &LabelEncoder{}

func NewLabelEncoder() *LabelEncoder {
	return &LabelEncoder{
		pool: sync.Pool{
			New: func() interface{} {
				return &bytes.Buffer{}
			},
		},
	}
}

func (e *LabelEncoder) Encode(labels []core.KeyValue) string {
	buf := e.pool.Get().(*bytes.Buffer)
	defer e.pool.Put(buf)
	buf.Reset()

	delimiter := "|#"

	for _, kv := range labels {
		_, _ = buf.WriteString(delimiter)
		_, _ = buf.WriteString(string(kv.Key))
		_, _ = buf.WriteRune(':')
		_, _ = buf.WriteString(kv.Value.Emit())
		delimiter = ","
	}
	return buf.String()
}
