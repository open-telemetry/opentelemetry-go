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

package metric

import (
	"bytes"
	"sync"

	"go.opentelemetry.io/otel/api/core"
)

// escapeChar is used to ensure uniqueness of the label encoding where
// keys or values contain either '=' or ','.  Since there is no parser
// needed for this encoding and its only requirement is to be unique,
// this choice is arbitrary.  Users will see these in some exporters
// (e.g., stdout), so the backslash ('\') is used a conventional choice.
const escapeChar = '\\'

type defaultLabelEncoder struct {
	// pool is a pool of labelset builders.  The buffers in this
	// pool grow to a size that most label encodings will not
	// allocate new memory.
	pool sync.Pool // *bytes.Buffer
}

var _ LabelEncoder = &defaultLabelEncoder{}

// NewDefaultLabelEncoder returns a label encoder that encodes labels
// in such a way that each escaped label's key is followed by an equal
// sign and then by an escaped label's value. All key-value pairs are
// separated by a comma.
//
// Escaping is done by prepending a backslash before either a
// backslash, equal sign or a comma.
func NewDefaultLabelEncoder() LabelEncoder {
	return &defaultLabelEncoder{
		pool: sync.Pool{
			New: func() interface{} {
				return &bytes.Buffer{}
			},
		},
	}
}

// Encode is a part of an implementation of the LabelEncoder
// interface.
func (d *defaultLabelEncoder) Encode(iter LabelIterator) string {
	buf := d.pool.Get().(*bytes.Buffer)
	defer d.pool.Put(buf)
	buf.Reset()

	for iter.Next() {
		i, kv := iter.IndexedLabel()
		if i > 0 {
			_, _ = buf.WriteRune(',')
		}
		copyAndEscape(buf, string(kv.Key))

		_, _ = buf.WriteRune('=')

		if kv.Value.Type() == core.STRING {
			copyAndEscape(buf, kv.Value.AsString())
		} else {
			_, _ = buf.WriteString(kv.Value.Emit())
		}
	}
	return buf.String()
}

// ID is a part of an implementation of the LabelEncoder interface.
func (*defaultLabelEncoder) ID() int64 {
	return defaultLabelEncoderID
}

func copyAndEscape(buf *bytes.Buffer, val string) {
	for _, ch := range val {
		switch ch {
		case '=', ',', escapeChar:
			buf.WriteRune(escapeChar)
		}
		buf.WriteRune(ch)
	}
}
