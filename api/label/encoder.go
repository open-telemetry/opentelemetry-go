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

package label

import (
	"bytes"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel/api/core"
)

type (
	// Encoder is a mechanism.
	Encoder interface {
		// Encode is called (concurrently) in instrumentation context.
		//
		// The expectation is that when setting up an export pipeline
		// both the batcher and the exporter will use the same label
		// encoder to avoid the duplicate computation of the encoded
		// labels in the export path.
		Encode(Iterator) string

		// ID should return a unique positive number associated with
		// the label encoder. Stateless label encoders could return
		// the same number regardless of an instance, stateful label
		// encoders should return a number depending on their state.
		ID() EncoderID
	}

	EncoderID struct {
		value int64
	}

	defaultLabelEncoder struct {
		// pool is a pool of labelset builders.  The buffers in this
		// pool grow to a size that most label encodings will not
		// allocate new memory.
		pool sync.Pool // *bytes.Buffer
	}
)

// escapeChar is used to ensure uniqueness of the label encoding where
// keys or values contain either '=' or ','.  Since there is no parser
// needed for this encoding and its only requirement is to be unique,
// this choice is arbitrary.  Users will see these in some exporters
// (e.g., stdout), so the backslash ('\') is used a conventional choice.
const escapeChar = '\\'

var (
	_ Encoder = &defaultLabelEncoder{}

	// labelEncoderIDCounter is for generating IDs for other label
	// encoders.
	encoderIDCounter int64 = 1

	defaultEncoderID = NewEncoderID()
)

// NewEncoderID returns a unique label encoder ID. It should be
// called once per each type of label encoder. Preferably in init() or
// in var definition.
func NewEncoderID() EncoderID {
	return EncoderID{value: atomic.AddInt64(&encoderIDCounter, 1)}
}

// NewDefaultEncoder returns a label encoder that encodes labels
// in such a way that each escaped label's key is followed by an equal
// sign and then by an escaped label's value. All key-value pairs are
// separated by a comma.
//
// Escaping is done by prepending a backslash before either a
// backslash, equal sign or a comma.
func NewDefaultEncoder() Encoder {
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
func (d *defaultLabelEncoder) Encode(iter Iterator) string {
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
func (*defaultLabelEncoder) ID() EncoderID {
	return defaultEncoderID
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

func (id EncoderID) Valid() bool {
	return id.value > 0
}
