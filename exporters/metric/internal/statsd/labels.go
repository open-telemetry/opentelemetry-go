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

	export "go.opentelemetry.io/otel/sdk/export/metric"
)

// LabelEncoder encodes metric labels in the dogstatsd syntax.
//
// TODO: find a link for this syntax.  It's been copied out of code,
// not a specification:
//
// https://github.com/stripe/veneur/blob/master/sinks/datadog/datadog.go
type LabelEncoder struct {
	pool sync.Pool
}

// sameCheck is used to test whether label encoders are the same.
type sameCheck interface {
	isStatsd()
}

var _ export.LabelEncoder = &LabelEncoder{}

// NewLabelEncoder returns a new encoder for dogstatsd-syntax metric
// labels.
func NewLabelEncoder() *LabelEncoder {
	return &LabelEncoder{
		pool: sync.Pool{
			New: func() interface{} {
				return &bytes.Buffer{}
			},
		},
	}
}

// Encode emits a string like "|#key1:value1,key2:value2".
func (e *LabelEncoder) Encode(iter export.LabelIterator) string {
	buf := e.pool.Get().(*bytes.Buffer)
	defer e.pool.Put(buf)
	buf.Reset()

	delimiter := "|#"

	for iter.Next() {
		kv := iter.Label()
		_, _ = buf.WriteString(delimiter)
		_, _ = buf.WriteString(string(kv.Key))
		_, _ = buf.WriteRune(':')
		_, _ = buf.WriteString(kv.Value.Emit())
		delimiter = ","
	}
	return buf.String()
}

func (e *LabelEncoder) isStatsd() {}

// ForceEncode returns a statsd label encoding, even if the exported
// labels were encoded by a different type of encoder.  Returns a
// boolean to indicate whether the labels were in fact re-encoded, to
// test for (and warn about) efficiency.
func (e *LabelEncoder) ForceEncode(labels export.Labels) (string, bool) {
	if _, ok := labels.Encoder().(sameCheck); ok {
		return labels.Encoded(), false
	}

	return e.Encode(labels.Iter()), true
}
