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

package otel

import (
	"context"
	"testing"
)

type testSpan struct {
	noopSpan

	ID int8
}

func TestContextSpan(t *testing.T) {
	ctx := context.Background()
	got, empty := SpanFromContext(ctx), noopSpan{}
	if got != empty {
		t.Errorf("SpanFromContext returned %v from an empty context, want %v", got, empty)
	}

	want := testSpan{ID: 0}
	ctx = ContextWithSpan(ctx, want)
	if got, ok := ctx.Value(currentSpanKey).(testSpan); !ok {
		t.Errorf("failed to set context with %#v", want)
	} else if got != want {
		t.Errorf("got %#v from context with current set, want %#v", got, want)
	}

	if got := SpanFromContext(ctx); got != want {
		t.Errorf("SpanFromContext returned %v from a set context, want %v", got, want)
	}

	want = testSpan{ID: 1}
	ctx = ContextWithSpan(ctx, want)
	if got, ok := ctx.Value(currentSpanKey).(testSpan); !ok {
		t.Errorf("failed to set context with %#v", want)
	} else if got != want {
		t.Errorf("got %#v from context with current overridden, want %#v", got, want)
	}

	if got := SpanFromContext(ctx); got != want {
		t.Errorf("SpanFromContext returned %v from a set context, want %v", got, want)
	}
}

func TestContextRemoteSpanReference(t *testing.T) {
	ctx := context.Background()
	got, empty := RemoteSpanReferenceFromContext(ctx), SpanReference{}
	if got != empty {
		t.Errorf("RemoteSpanReferenceFromContext returned %v from an empty context, want %v", got, empty)
	}

	want := SpanReference{TraceID: [16]byte{1}, SpanID: [8]byte{42}}
	ctx = ContextWithRemoteSpanReference(ctx, want)
	if got, ok := ctx.Value(remoteContextKey).(SpanReference); !ok {
		t.Errorf("failed to set SpanReference with %#v", want)
	} else if got != want {
		t.Errorf("got %#v from context with remote set, want %#v", got, want)
	}

	if got := RemoteSpanReferenceFromContext(ctx); got != want {
		t.Errorf("RemoteSpanReferenceFromContext returned %v from a set context, want %v", got, want)
	}

	want = SpanReference{TraceID: [16]byte{1}, SpanID: [8]byte{43}}
	ctx = ContextWithRemoteSpanReference(ctx, want)
	if got, ok := ctx.Value(remoteContextKey).(SpanReference); !ok {
		t.Errorf("failed to set SpanReference with %#v", want)
	} else if got != want {
		t.Errorf("got %#v from context with remote overridden, want %#v", got, want)
	}

	if got := RemoteSpanReferenceFromContext(ctx); got != want {
		t.Errorf("RemoteSpanReferenceFromContext returned %v from a set context, want %v", got, want)
	}
}

func TestIsValid(t *testing.T) {
	for _, testcase := range []struct {
		name string
		tid  TraceID
		sid  SpanID
		want bool
	}{
		{
			name: "SpanReference.IsValid() returns true if sc has both an Trace ID and Span ID",
			tid:  [16]byte{1},
			sid:  [8]byte{42},
			want: true,
		}, {
			name: "SpanReference.IsValid() returns false if sc has neither an Trace ID nor Span ID",
			tid:  TraceID([16]byte{}),
			sid:  [8]byte{},
			want: false,
		}, {
			name: "SpanReference.IsValid() returns false if sc has a Span ID but not a Trace ID",
			tid:  TraceID([16]byte{}),
			sid:  [8]byte{42},
			want: false,
		}, {
			name: "SpanReference.IsValid() returns false if sc has a Trace ID but not a Span ID",
			tid:  TraceID([16]byte{1}),
			sid:  [8]byte{},
			want: false,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			sc := SpanReference{
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
		tid   TraceID
		valid bool
	}{
		{
			name:  "Valid TraceID",
			tid:   TraceID([16]byte{128, 241, 152, 238, 86, 52, 59, 168, 100, 254, 139, 42, 87, 211, 239, 247}),
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
			tid, err := TraceIDFromHex(testcase.hex)

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
		tid  TraceID
		want bool
	}{
		{
			name: "SpanReference.HasTraceID() returns true if both Low and High are nonzero",
			tid:  TraceID([16]byte{1}),
			want: true,
		}, {
			name: "SpanReference.HasTraceID() returns false if neither Low nor High are nonzero",
			tid:  TraceID{},
			want: false,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (sc SpanReference) HasTraceID() bool{}
			sc := SpanReference{TraceID: testcase.tid}
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
		sc   SpanReference
		want bool
	}{
		{
			name: "SpanReference.HasSpanID() returns true if self.SpanID != 0",
			sc:   SpanReference{SpanID: [8]byte{42}},
			want: true,
		}, {
			name: "SpanReference.HasSpanID() returns false if self.SpanID == 0",
			sc:   SpanReference{},
			want: false,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (sc SpanReference) HasSpanID() bool {}
			have := testcase.sc.HasSpanID()
			if have != testcase.want {
				t.Errorf("Want: %v, but have: %v", testcase.want, have)
			}
		})
	}
}

func TestSpanReferenceIsSampled(t *testing.T) {
	for _, testcase := range []struct {
		name string
		sc   SpanReference
		want bool
	}{
		{
			name: "sampled",
			sc: SpanReference{
				TraceID:    TraceID([16]byte{1}),
				TraceFlags: FlagsSampled,
			},
			want: true,
		}, {
			name: "unused bits are ignored, still not sampled",
			sc: SpanReference{
				TraceID:    TraceID([16]byte{1}),
				TraceFlags: ^FlagsSampled,
			},
			want: false,
		}, {
			name: "unused bits are ignored, still sampled",
			sc: SpanReference{
				TraceID:    TraceID([16]byte{1}),
				TraceFlags: FlagsSampled | ^FlagsSampled,
			},
			want: true,
		}, {
			name: "not sampled/default",
			sc:   SpanReference{TraceID: TraceID{}},
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
		tid  TraceID
		want string
	}{
		{
			name: "TraceID.String returns string representation of self.TraceID values > 0",
			tid:  TraceID([16]byte{255}),
			want: "ff000000000000000000000000000000",
		},
		{
			name: "TraceID.String returns string representation of self.TraceID values == 0",
			tid:  TraceID([16]byte{}),
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
		sid  SpanID
		want string
	}{
		{
			name: "SpanID.String returns string representation of self.SpanID values > 0",
			sid:  SpanID([8]byte{255}),
			want: "ff00000000000000",
		},
		{
			name: "SpanID.String returns string representation of self.SpanID values == 0",
			sid:  SpanID([8]byte{}),
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

func TestValidateSpanKind(t *testing.T) {
	tests := []struct {
		in   SpanKind
		want SpanKind
	}{
		{
			SpanKindUnspecified,
			SpanKindInternal,
		},
		{

			SpanKindInternal,
			SpanKindInternal,
		},
		{

			SpanKindServer,
			SpanKindServer,
		},
		{

			SpanKindClient,
			SpanKindClient,
		},
		{
			SpanKindProducer,
			SpanKindProducer,
		},
		{
			SpanKindConsumer,
			SpanKindConsumer,
		},
	}
	for _, test := range tests {
		if got := ValidateSpanKind(test.in); got != test.want {
			t.Errorf("ValidateSpanKind(%#v) = %#v, want %#v", test.in, got, test.want)
		}
	}
}

func TestSpanKindString(t *testing.T) {
	tests := []struct {
		in   SpanKind
		want string
	}{
		{
			SpanKindUnspecified,
			"unspecified",
		},
		{

			SpanKindInternal,
			"internal",
		},
		{

			SpanKindServer,
			"server",
		},
		{

			SpanKindClient,
			"client",
		},
		{
			SpanKindProducer,
			"producer",
		},
		{
			SpanKindConsumer,
			"consumer",
		},
	}
	for _, test := range tests {
		if got := test.in.String(); got != test.want {
			t.Errorf("%#v.String() = %#v, want %#v", test.in, got, test.want)
		}
	}
}
