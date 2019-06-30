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

package core

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"

	"github.com/open-telemetry/opentelemetry-go/api/unit"
)

type ScopeID struct {
	EventID
	SpanContext
}

type SpanContext struct {
	TraceIDHigh uint64
	TraceIDLow  uint64
	SpanID      uint64
}

type EventID uint64

type BaseMeasure interface {
	Name() string
	Description() string
	Unit() unit.Unit

	DefinitionID() EventID
}

type Measure interface {
	BaseMeasure

	M(float64) Measurement
	V(float64) KeyValue
}

type Measurement struct {
	// NOTE: If we add a ScopeID field this can carry
	// pre-aggregated measures via the stats.Record API.
	Measure Measure
	Value   float64
	ScopeID ScopeID
}

type Key interface {
	BaseMeasure

	Value(ctx context.Context) KeyValue

	Bool(v bool) KeyValue

	Int(v int) KeyValue
	Int32(v int32) KeyValue
	Int64(v int64) KeyValue

	Uint(v uint) KeyValue
	Uint32(v uint32) KeyValue
	Uint64(v uint64) KeyValue

	Float32(v float32) KeyValue
	Float64(v float64) KeyValue

	String(v string) KeyValue
	Bytes(v []byte) KeyValue
}

type KeyValue struct {
	Key   Key
	Value Value
}

type ValueType int

type Value struct {
	Type    ValueType
	Bool    bool
	Int64   int64
	Uint64  uint64
	Float64 float64
	String  string
	Bytes   []byte

	// TODO Lazy value type?
}

type MutatorOp int

type Mutator struct {
	MutatorOp
	KeyValue
	MeasureMetadata
}

type MeasureMetadata struct {
	MaxHops int // -1 == infinite, 0 == do not propagate

	// TODO time to live?
}

const (
	INVALID ValueType = iota
	BOOL
	INT32
	INT64
	UINT32
	UINT64
	FLOAT32
	FLOAT64
	STRING
	BYTES

	INSERT MutatorOp = iota
	UPDATE
	UPSERT
	DELETE
)

var (
	// INVALID_SPAN_CONTEXT is meant for internal use to return invalid span context during error
	// conditions.
	INVALID_SPAN_CONTEXT = SpanContext{}
)

func (sc SpanContext) HasTraceID() bool {
	return sc.TraceIDHigh != 0 || sc.TraceIDLow != 0
}

func (sc SpanContext) HasSpanID() bool {
	return sc.SpanID != 0
}

func (sc SpanContext) SpanIDString() string {
	p := fmt.Sprintf("%.16x", sc.SpanID)
	return p[0:3] + ".." + p[13:16]
}

func (sc SpanContext) TraceIDString() string {
	p1 := fmt.Sprintf("%.16x", sc.TraceIDHigh)
	p2 := fmt.Sprintf("%.16x", sc.TraceIDLow)
	return p1[0:3] + ".." + p2[13:16]
}

// TODO make this a lazy one-time conversion.
func (v Value) Emit() string {
	switch v.Type {
	case BOOL:
		return fmt.Sprint(v.Bool)
	case INT32, INT64:
		return fmt.Sprint(v.Int64)
	case UINT32, UINT64:
		return fmt.Sprint(v.Uint64)
	case FLOAT32, FLOAT64:
		return fmt.Sprint(v.Float64)
	case STRING:
		return v.String
	case BYTES:
		return string(v.Bytes)
	}
	return "unknown"
}

func (m Mutator) WithMaxHops(hops int) Mutator {
	m.MaxHops = hops
	return m
}

func (e EventID) Scope() ScopeID {
	return ScopeID{
		EventID: e,
	}
}

func (s SpanContext) Scope() ScopeID {
	return ScopeID{
		SpanContext: s,
	}
}

func (m Measurement) With(id ScopeID) Measurement {
	m.ScopeID = id
	return m
}

func GrpcCodeToString(c codes.Code) string {
	return c.String()
}
