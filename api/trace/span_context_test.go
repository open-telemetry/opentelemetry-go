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

package trace_test

import (
	"testing"

	"go.opentelemetry.io/otel/api/trace"
)

func TestIsValid(t *testing.T) {
	for _, testcase := range []struct {
		name string
		tid  trace.ID
		sid  trace.SpanID
		want bool
	}{
		{
			name: "SpanContext.IsValid() returns true if sc has both an Trace ID and Span ID",
			tid:  [16]byte{1},
			sid:  [8]byte{42},
			want: true,
		}, {
			name: "SpanContext.IsValid() returns false if sc has neither an Trace ID nor Span ID",
			tid:  trace.ID([16]byte{}),
			sid:  [8]byte{},
			want: false,
		}, {
			name: "SpanContext.IsValid() returns false if sc has a Span ID but not a Trace ID",
			tid:  trace.ID([16]byte{}),
			sid:  [8]byte{42},
			want: false,
		}, {
			name: "SpanContext.IsValid() returns false if sc has a Trace ID but not a Span ID",
			tid:  trace.ID([16]byte{1}),
			sid:  [8]byte{},
			want: false,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			sc := trace.SpanContext{
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

func TestIsValidFromHex(t *testing.T) {
	for _, testcase := range []struct {
		name  string
		hex   string
		tid   trace.ID
		valid bool
	}{
		{
			name:  "Valid TraceID",
			tid:   trace.ID([16]byte{128, 241, 152, 238, 86, 52, 59, 168, 100, 254, 139, 42, 87, 211, 239, 247}),
			hex:   "80f198ee56343ba864fe8b2a57d3eff7",
			valid: true,
		}, {
			name:  "Invalid TraceID with invalid length",
			hex:   "80f198ee56343ba864fe8b2a57d3eff",
			valid: false,
		}, {
			name:  "Invalid TraceID with invalid char",
			hex:   "80f198ee56343ba864fe8b2a57d3efg7",
			valid: false,
		}, {
			name:  "Invalid TraceID with uppercase",
			hex:   "80f198ee56343ba864fe8b2a57d3efF7",
			valid: false,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			tid, err := trace.IDFromHex(testcase.hex)

			if testcase.valid && err != nil {
				t.Errorf("Expected TraceID %s to be valid but end with error %s", testcase.hex, err.Error())
			}

			if !testcase.valid && err == nil {
				t.Errorf("Expected TraceID %s to be invalid but end no error", testcase.hex)
			}

			if tid != testcase.tid {
				t.Errorf("Want: %v, but have: %v", testcase.tid, tid)
			}
		})
	}
}

func TestHasTraceID(t *testing.T) {
	for _, testcase := range []struct {
		name string
		tid  trace.ID
		want bool
	}{
		{
			name: "SpanContext.HasTraceID() returns true if both Low and High are nonzero",
			tid:  trace.ID([16]byte{1}),
			want: true,
		}, {
			name: "SpanContext.HasTraceID() returns false if neither Low nor High are nonzero",
			tid:  trace.ID{},
			want: false,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (sc SpanContext) HasTraceID() bool{}
			sc := trace.SpanContext{TraceID: testcase.tid}
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
		sc   trace.SpanContext
		want bool
	}{
		{
			name: "SpanContext.HasSpanID() returns true if self.SpanID != 0",
			sc:   trace.SpanContext{SpanID: [8]byte{42}},
			want: true,
		}, {
			name: "SpanContext.HasSpanID() returns false if self.SpanID == 0",
			sc:   trace.SpanContext{},
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

func TestSpanContextIsSampled(t *testing.T) {
	for _, testcase := range []struct {
		name string
		sc   trace.SpanContext
		want bool
	}{
		{
			name: "sampled",
			sc: trace.SpanContext{
				TraceID:    trace.ID([16]byte{1}),
				TraceFlags: trace.FlagsSampled,
			},
			want: true,
		}, {
			name: "sampled plus unused",
			sc: trace.SpanContext{
				TraceID:    trace.ID([16]byte{1}),
				TraceFlags: trace.FlagsSampled | trace.FlagsUnused,
			},
			want: true,
		}, {
			name: "not sampled/default",
			sc:   trace.SpanContext{TraceID: trace.ID{}},
			want: false,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			have := testcase.sc.IsSampled()
			if have != testcase.want {
				t.Errorf("Want: %v, but have: %v", testcase.want, have)
			}
		})
	}
}

func TestStringTraceID(t *testing.T) {
	for _, testcase := range []struct {
		name string
		tid  trace.ID
		want string
	}{
		{
			name: "TraceID.String returns string representation of self.TraceID values > 0",
			tid:  trace.ID([16]byte{255}),
			want: "ff000000000000000000000000000000",
		},
		{
			name: "TraceID.String returns string representation of self.TraceID values == 0",
			tid:  trace.ID([16]byte{}),
			want: "00000000000000000000000000000000",
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (t TraceID) String() string {}
			have := testcase.tid.String()
			if have != testcase.want {
				t.Errorf("Want: %s, but have: %s", testcase.want, have)
			}
		})
	}
}

func TestStringSpanID(t *testing.T) {
	for _, testcase := range []struct {
		name string
		sid  trace.SpanID
		want string
	}{
		{
			name: "SpanID.String returns string representation of self.SpanID values > 0",
			sid:  trace.SpanID([8]byte{255}),
			want: "ff00000000000000",
		},
		{
			name: "SpanID.String returns string representation of self.SpanID values == 0",
			sid:  trace.SpanID([8]byte{}),
			want: "0000000000000000",
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (t TraceID) String() string {}
			have := testcase.sid.String()
			if have != testcase.want {
				t.Errorf("Want: %s, but have: %s", testcase.want, have)
			}
		})
	}
}
