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

package core_test

import (
	"testing"

	"go.opentelemetry.io/api/core"
)

func TestIsValid(t *testing.T) {
	for _, testcase := range []struct {
		name string
		tid  core.TraceID
		sid  uint64
		want bool
	}{
		{
			name: "SpanContext.IsValid() returns true if sc has both an Trace ID and Span ID",
			tid:  core.TraceID{High: uint64(42)},
			sid:  uint64(42),
			want: true,
		}, {
			name: "SpanContext.IsValid() returns false if sc has neither an Trace ID nor Span ID",
			tid:  core.TraceID{High: uint64(0)},
			sid:  uint64(0),
			want: false,
		}, {
			name: "SpanContext.IsValid() returns false if sc has a Span ID but not a Trace ID",
			tid:  core.TraceID{High: uint64(0)},
			sid:  uint64(42),
			want: false,
		}, {
			name: "SpanContext.IsValid() returns false if sc has a Trace ID but not a Span ID",
			tid:  core.TraceID{High: uint64(42)},
			sid:  uint64(0),
			want: false,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			sc := core.SpanContext{
				TraceID: testcase.tid,
				SpanID:  testcase.sid,
			}
			have := sc.IsValid()
			if have != testcase.want {
				t.Errorf("Want: %v, but have: %v", testcase.want, have)
			}
		})
	}
}

func TestHasTraceID(t *testing.T) {
	for _, testcase := range []struct {
		name string
		tid  core.TraceID
		want bool
	}{
		{
			name: "SpanContext.HasTraceID() returns true if both Low and High are nonzero",
			tid:  core.TraceID{High: uint64(42), Low: uint64(42)},
			want: true,
		}, {
			name: "SpanContext.HasTraceID() returns false if neither Low nor High are nonzero",
			tid:  core.TraceID{},
			want: false,
		}, {
			name: "SpanContext.HasTraceID() returns true if High != 0",
			tid:  core.TraceID{High: uint64(42)},
			want: true,
		}, {
			name: "SpanContext.HasTraceID() returns true if Low != 0",
			tid:  core.TraceID{Low: uint64(42)},
			want: true,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (sc SpanContext) HasTraceID() bool{}
			sc := core.SpanContext{TraceID: testcase.tid}
			have := sc.HasTraceID()
			if have != testcase.want {
				t.Errorf("Want: %v, but have: %v", testcase.want, have)
			}
		})
	}
}

func TestHasSpanID(t *testing.T) {
	for _, testcase := range []struct {
		name string
		sc   core.SpanContext
		want bool
	}{
		{
			name: "SpanContext.HasSpanID() returns true if self.SpanID != 0",
			sc:   core.SpanContext{SpanID: uint64(42)},
			want: true,
		}, {
			name: "SpanContext.HasSpanID() returns false if self.SpanID == 0",
			sc:   core.SpanContext{},
			want: false,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (sc SpanContext) HasSpanID() bool {}
			have := testcase.sc.HasSpanID()
			if have != testcase.want {
				t.Errorf("Want: %v, but have: %v", testcase.want, have)
			}
		})
	}
}

func TestSpanIDString(t *testing.T) {
	for _, testcase := range []struct {
		name string
		sc   core.SpanContext
		want string
	}{
		{
			name: "SpanContext.SpanIDString returns string representation of self.TraceID values > 0",
			sc:   core.SpanContext{SpanID: uint64(42)},
			want: `000000000000002a`,
		}, {
			name: "SpanContext.SpanIDString returns string representation of self.TraceID values == 0",
			sc:   core.SpanContext{},
			want: `0000000000000000`,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (sc SpanContext) SpanIDString() string {}
			have := testcase.sc.SpanIDString()
			if have != testcase.want {
				t.Errorf("Want: %s, but have: %s", testcase.want, have)
			}
		})
	}
}

func TestTraceIDString(t *testing.T) {
	for _, testcase := range []struct {
		name string
		sc   core.SpanContext
		want string
	}{
		{
			name: "SpanContext.TraceIDString returns string representation of self.TraceID values > 0",
			sc: core.SpanContext{
				TraceID: core.TraceID{
					High: uint64(42),
					Low:  uint64(42),
				},
			},
			want: `000000000000002a000000000000002a`,
		}, {
			name: "SpanContext.TraceIDString returns string representation of self.TraceID values == 0",
			sc:   core.SpanContext{TraceID: core.TraceID{}},
			want: `00000000000000000000000000000000`,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (sc SpanContext) TraceIDString() string {}
			have := testcase.sc.TraceIDString()
			if have != testcase.want {
				t.Errorf("Want: %s, but have: %s", testcase.want, have)
			}
		})
	}
}

func TestSpanContextIsSampled(t *testing.T) {
	for _, testcase := range []struct {
		name string
		sc   core.SpanContext
		want bool
	}{
		{
			name: "sampled",
			sc: core.SpanContext{
				TraceID: core.TraceID{
					High: uint64(42),
					Low:  uint64(42),
				},
				TraceOptions: core.TraceOptionSampled,
			},
			want: true,
		}, {
			name: "sampled plus unused",
			sc: core.SpanContext{
				TraceID: core.TraceID{
					High: uint64(42),
					Low:  uint64(42),
				},
				TraceOptions: core.TraceOptionSampled | core.TraceOptionUnused,
			},
			want: true,
		}, {
			name: "not sampled/default",
			sc:   core.SpanContext{TraceID: core.TraceID{}},
			want: false,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (sc SpanContext) TraceIDString() string {}
			have := testcase.sc.IsSampled()
			if have != testcase.want {
				t.Errorf("Want: %v, but have: %v", testcase.want, have)
			}
		})
	}
}
