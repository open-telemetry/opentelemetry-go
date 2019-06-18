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

	"github.com/open-telemetry/opentelemetry-go/api/unit"
	"strconv"
)

type (
	ScopeID struct {
		EventID
		SpanContext
	}

	SpanContext struct {
		TraceIDHigh uint64
		TraceIDLow  uint64
		SpanID      uint64
	}

	EventID uint64

	BaseMeasure interface {
		Name() string
		Description() string
		Unit() unit.Unit

		DefinitionID() EventID
	}

	Measure interface {
		BaseMeasure

		M(float64) Measurement
		V(float64) KeyValue
	}

	Measurement struct {
		// NOTE: If we add a ScopeID field this can carry
		// pre-aggregated measures via the stats.Record API.
		Measure Measure
		Value   float64
		ScopeID ScopeID
	}

	Key interface {
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

	KeyValue struct {
		Key   Key
		Value Value
	}

	ValueType int

	Value struct {
		Type    ValueType
		Bool    bool
		Int64   int64
		Uint64  uint64
		Float64 float64
		String  string
		Bytes   []byte

		// TODO Lazy value type?
	}

	MutatorOp int

	Mutator struct {
		MutatorOp
		KeyValue
		MeasureMetadata
	}

	MeasureMetadata struct {
		MaxHops int // -1 == infinite, 0 == do not propagate

		// TODO time to live?
	}

	Status int32
)

// OpenCensus uses google.golang.org/grpc/ Following is copied from there.
const (
	// OK is returned on success.
	OK Status = 0

	// Canceled indicates the operation was canceled (typically by the caller).
	Canceled Status = 1

	// Unknown error. An example of where this error may be returned is
	// if a Status value received from another address space belongs to
	// an error-space that is not known in this address space. Also
	// errors raised by APIs that do not return enough error information
	// may be converted to this error.
	Unknown Status = 2

	// InvalidArgument indicates client specified an invalid argument.
	// Note that this differs from FailedPrecondition. It indicates arguments
	// that are problematic regardless of the state of the system
	// (e.g., a malformed file name).
	InvalidArgument Status = 3

	// DeadlineExceeded means operation expired before completion.
	// For operations that change the state of the system, this error may be
	// returned even if the operation has completed successfully. For
	// example, a successful response from a server could have been delayed
	// long enough for the deadline to expire.
	DeadlineExceeded Status = 4

	// NotFound means some requested entity (e.g., file or directory) was
	// not found.
	NotFound Status = 5

	// AlreadyExists means an attempt to create an entity failed because one
	// already exists.
	AlreadyExists Status = 6

	// PermissionDenied indicates the caller does not have permission to
	// execute the specified operation. It must not be used for rejections
	// caused by exhausting some resource (use ResourceExhausted
	// instead for those errors). It must not be
	// used if the caller cannot be identified (use Unauthenticated
	// instead for those errors).
	PermissionDenied Status = 7

	// ResourceExhausted indicates some resource has been exhausted, perhaps
	// a per-user quota, or perhaps the entire file system is out of space.
	ResourceExhausted Status = 8

	// FailedPrecondition indicates operation was rejected because the
	// system is not in a state required for the operation's execution.
	// For example, directory to be deleted may be non-empty, an rmdir
	// operation is applied to a non-directory, etc.
	//
	// A litmus test that may help a service implementor in deciding
	// between FailedPrecondition, Aborted, and Unavailable:
	//  (a) Use Unavailable if the client can retry just the failing call.
	//  (b) Use Aborted if the client should retry at a higher-level
	//      (e.g., restarting a read-modify-write sequence).
	//  (c) Use FailedPrecondition if the client should not retry until
	//      the system state has been explicitly fixed. E.g., if an "rmdir"
	//      fails because the directory is non-empty, FailedPrecondition
	//      should be returned since the client should not retry unless
	//      they have first fixed up the directory by deleting files from it.
	//  (d) Use FailedPrecondition if the client performs conditional
	//      REST Get/Update/Delete on a resource and the resource on the
	//      server does not match the condition. E.g., conflicting
	//      read-modify-write on the same resource.
	FailedPrecondition Status = 9

	// Aborted indicates the operation was aborted, typically due to a
	// concurrency issue like sequencer check failures, transaction aborts,
	// etc.
	//
	// See litmus test above for deciding between FailedPrecondition,
	// Aborted, and Unavailable.
	Aborted Status = 10

	// OutOfRange means operation was attempted past the valid range.
	// E.g., seeking or reading past end of file.
	//
	// Unlike InvalidArgument, this error indicates a problem that may
	// be fixed if the system state changes. For example, a 32-bit file
	// system will generate InvalidArgument if asked to read at an
	// offset that is not in the range [0,2^32-1], but it will generate
	// OutOfRange if asked to read from an offset past the current
	// file size.
	//
	// There is a fair bit of overlap between FailedPrecondition and
	// OutOfRange. We recommend using OutOfRange (the more specific
	// error) when it applies so that callers who are iterating through
	// a space can easily look for an OutOfRange error to detect when
	// they are done.
	OutOfRange Status = 11

	// Unimplemented indicates operation is not implemented or not
	// supported/enabled in this service.
	Unimplemented Status = 12

	// Internal errors. Means some invariants expected by underlying
	// system has been broken. If you see one of these errors,
	// something is very broken.
	Internal Status = 13

	// Unavailable indicates the service is currently unavailable.
	// This is a most likely a transient condition and may be corrected
	// by retrying with a backoff.
	//
	// See litmus test above for deciding between FailedPrecondition,
	// Aborted, and Unavailable.
	Unavailable Status = 14

	// DataLoss indicates unrecoverable data loss or corruption.
	DataLoss Status = 15

	// Unauthenticated indicates the request does not have valid
	// authentication credentials for the operation.
	Unauthenticated Status = 16

	_maxStatus = 17
)

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

func StatusToString(s Status) string {
	switch c := s; c {
	case OK:
		return "OK"
	case Canceled:
		return "CANCELLED"
	case Unknown:
		return "UNKNOWN"
	case InvalidArgument:
		return "INVALID_ARGUMENT"
	case DeadlineExceeded:
		return "DEADLINE_EXCEEDED"
	case NotFound:
		return "NOT_FOUND"
	case AlreadyExists:
		return "ALREADY_EXISTS"
	case PermissionDenied:
		return "PERMISSION_DENIED"
	case ResourceExhausted:
		return "RESOURCE_EXHAUSTED"
	case FailedPrecondition:
		return "FAILED_PRECONDITION"
	case Aborted:
		return "ABORTED"
	case OutOfRange:
		return "OUT_OF_RANGE"
	case Unimplemented:
		return "UNIMPLEMENTED"
	case Internal:
		return "INTERNAL"
	case Unavailable:
		return "UNAVAILABLE"
	case DataLoss:
		return "DATA_LOSS"
	case Unauthenticated:
		return "UNAUTHENTICATED"
	default:
		return "STATUS_" + strconv.FormatInt(int64(c), 10)
	}
}
