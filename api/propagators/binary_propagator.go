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

package propagators

import (
	"go.opentelemetry.io/otel/api/core"
)

// BinaryFormat is an interface that specifies methods to convert SpanContext
// to/from byte array.
type BinaryFormat interface {
	// ToBytes serializes span context into a byte array and returns the array.
	ToBytes(sc core.SpanContext) []byte

	// FromBytes de-serializes byte array into span context and returns the span context.
	FromBytes([]byte) core.SpanContext
}

var _ BinaryFormat = binary{}

type binary struct{}

// Binary creates a new propagator. The propagator implements
// ToBytes and FromBytes method to transform SpanContext to/from byte array.
func Binary() BinaryFormat {
	return binary{}
}

// ToBytes implements ToBytes method of propagators.BinaryFormat.
// It serializes core.SpanContext into a byte array.
func (bp binary) ToBytes(sc core.SpanContext) []byte {
	if sc == core.EmptySpanContext() {
		return nil
	}
	var b [29]byte
	copy(b[2:18], sc.TraceID[:])
	b[18] = 1
	copy(b[19:27], sc.SpanID[:])
	b[27] = 2
	b[28] = sc.TraceFlags
	return b[:]
}

// FromBytes implements FromBytes method of propagators.BinaryFormat.
// It de-serializes bytes into core.SpanContext.
func (bp binary) FromBytes(b []byte) (sc core.SpanContext) {
	if len(b) == 0 {
		return core.EmptySpanContext()
	}
	b = b[1:]
	if len(b) >= 17 && b[0] == 0 {
		copy(sc.TraceID[:], b[1:17])
		b = b[17:]
	} else {
		return core.EmptySpanContext()
	}
	if len(b) >= 9 && b[0] == 1 {
		copy(sc.SpanID[:], b[1:9])
		b = b[9:]
	}
	if len(b) >= 2 && b[0] == 2 {
		sc.TraceFlags = b[1]
	}
	if sc.IsValid() {
		return sc
	}
	return core.EmptySpanContext()
}
