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

package trace

import (
	"context"
	"fmt"
	"testing"

	"go.opentelemetry.io/otel/label"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testSpan struct {
	noopSpan

	ID byte
}

func (s testSpan) SpanContext() SpanContext { return SpanContext{SpanID: [8]byte{s.ID}} }

func TestContextSpan(t *testing.T) {
	testCases := []struct {
		name         string
		context      context.Context
		expectedSpan Span
	}{
		{
			name:         "empty context",
			context:      context.Background(),
			expectedSpan: noopSpan{},
		},
		{
			name:         "span 0",
			context:      ContextWithSpan(context.Background(), testSpan{ID: 0}),
			expectedSpan: testSpan{ID: 0},
		},
		{
			name:         "span 1",
			context:      ContextWithSpan(context.Background(), testSpan{ID: 1}),
			expectedSpan: testSpan{ID: 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			span := SpanFromContext(tc.context)
			assert.Equal(t, tc.expectedSpan, span)

			if _, ok := tc.expectedSpan.(noopSpan); !ok {
				span, ok := tc.context.Value(currentSpanKey).(testSpan)
				assert.True(t, ok)
				assert.Equal(t, tc.expectedSpan.(testSpan), span)
			}
		})
	}
}

func TestContextRemoteSpanContext(t *testing.T) {
	ctx := context.Background()
	got, empty := RemoteSpanContextFromContext(ctx), SpanContext{}
	if !assertSpanContextEqual(got, empty) {
		t.Errorf("RemoteSpanContextFromContext returned %v from an empty context, want %v", got, empty)
	}

	want := SpanContext{TraceID: [16]byte{1}, SpanID: [8]byte{42}}
	ctx = ContextWithRemoteSpanContext(ctx, want)
	if got, ok := ctx.Value(remoteContextKey).(SpanContext); !ok {
		t.Errorf("failed to set SpanContext with %#v", want)
	} else if !assertSpanContextEqual(got, want) {
		t.Errorf("got %#v from context with remote set, want %#v", got, want)
	}

	if got := RemoteSpanContextFromContext(ctx); !assertSpanContextEqual(got, want) {
		t.Errorf("RemoteSpanContextFromContext returned %v from a set context, want %v", got, want)
	}

	want = SpanContext{TraceID: [16]byte{1}, SpanID: [8]byte{43}}
	ctx = ContextWithRemoteSpanContext(ctx, want)
	if got, ok := ctx.Value(remoteContextKey).(SpanContext); !ok {
		t.Errorf("failed to set SpanContext with %#v", want)
	} else if !assertSpanContextEqual(got, want) {
		t.Errorf("got %#v from context with remote set, want %#v", got, want)
	}

	got = RemoteSpanContextFromContext(ctx)
	if !assertSpanContextEqual(got, want) {
		t.Errorf("RemoteSpanContextFromContext returned %v from a set context, want %v", got, want)
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
			name: "SpanContext.IsValid() returns true if sc has both an Trace ID and Span ID",
			tid:  [16]byte{1},
			sid:  [8]byte{42},
			want: true,
		}, {
			name: "SpanContext.IsValid() returns false if sc has neither an Trace ID nor Span ID",
			tid:  TraceID([16]byte{}),
			sid:  [8]byte{},
			want: false,
		}, {
			name: "SpanContext.IsValid() returns false if sc has a Span ID but not a Trace ID",
			tid:  TraceID([16]byte{}),
			sid:  [8]byte{42},
			want: false,
		}, {
			name: "SpanContext.IsValid() returns false if sc has a Trace ID but not a Span ID",
			tid:  TraceID([16]byte{1}),
			sid:  [8]byte{},
			want: false,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			sc := SpanContext{
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
			name: "SpanContext.HasTraceID() returns true if both Low and High are nonzero",
			tid:  TraceID([16]byte{1}),
			want: true,
		}, {
			name: "SpanContext.HasTraceID() returns false if neither Low nor High are nonzero",
			tid:  TraceID{},
			want: false,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (sc SpanContext) HasTraceID() bool{}
			sc := SpanContext{TraceID: testcase.tid}
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
		sc   SpanContext
		want bool
	}{
		{
			name: "SpanContext.HasSpanID() returns true if self.SpanID != 0",
			sc:   SpanContext{SpanID: [8]byte{42}},
			want: true,
		}, {
			name: "SpanContext.HasSpanID() returns false if self.SpanID == 0",
			sc:   SpanContext{},
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
		sc   SpanContext
		want bool
	}{
		{
			name: "sampled",
			sc: SpanContext{
				TraceID:    TraceID([16]byte{1}),
				TraceFlags: FlagsSampled,
			},
			want: true,
		}, {
			name: "unused bits are ignored, still not sampled",
			sc: SpanContext{
				TraceID:    TraceID([16]byte{1}),
				TraceFlags: ^FlagsSampled,
			},
			want: false,
		}, {
			name: "unused bits are ignored, still sampled",
			sc: SpanContext{
				TraceID:    TraceID([16]byte{1}),
				TraceFlags: FlagsSampled | ^FlagsSampled,
			},
			want: true,
		}, {
			name: "not sampled/default",
			sc:   SpanContext{TraceID: TraceID{}},
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

func TestSpanContextFromContext(t *testing.T) {
	testCases := []struct {
		name                string
		context             context.Context
		expectedSpanContext SpanContext
	}{
		{
			name:    "empty context",
			context: context.Background(),
		},
		{
			name:                "span 1",
			context:             ContextWithSpan(context.Background(), testSpan{ID: 1}),
			expectedSpanContext: SpanContext{SpanID: [8]byte{1}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			spanContext := SpanContextFromContext(tc.context)
			assert.Equal(t, tc.expectedSpanContext, spanContext)
		})
	}
}

func TestTraceStateString(t *testing.T) {
	testCases := []struct {
		name        string
		traceState  TraceState
		expectedStr string
	}{
		{
			name: "Non-empty trace state",
			traceState: TraceState{
				kvs: []label.KeyValue{
					label.String("key1", "val1"),
					label.String("key2", "val2"),
					label.String("key3@vendor", "val3"),
				},
			},
			expectedStr: "key1=val1,key2=val2,key3@vendor=val3",
		},
		{
			name:        "Empty trace state",
			traceState:  TraceState{},
			expectedStr: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedStr, tc.traceState.String())
		})
	}
}

func TestTraceStateGet(t *testing.T) {
	testCases := []struct {
		name          string
		traceState    TraceState
		key           label.Key
		expectedValue string
	}{
		{
			name:          "OK case",
			traceState:    TraceState{kvsWithMaxMembers},
			key:           "key16",
			expectedValue: "value16",
		},
		{
			name:          "Not found",
			traceState:    TraceState{kvsWithMaxMembers},
			key:           "keyxx",
			expectedValue: "",
		},
		{
			name:          "Invalid key",
			traceState:    TraceState{kvsWithMaxMembers},
			key:           "key!",
			expectedValue: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			kv := tc.traceState.Get(tc.key)
			assert.Equal(t, tc.expectedValue, kv.AsString())
		})
	}
}

func TestTraceStateDelete(t *testing.T) {
	testCases := []struct {
		name               string
		traceState         TraceState
		key                label.Key
		expectedTraceState TraceState
		expectedErr        error
	}{
		{
			name: "OK case",
			traceState: TraceState{
				kvs: []label.KeyValue{
					label.String("key1", "val1"),
					label.String("key2", "val2"),
					label.String("key3", "val3"),
				},
			},
			key: "key2",
			expectedTraceState: TraceState{
				kvs: []label.KeyValue{
					label.String("key1", "val1"),
					label.String("key3", "val3"),
				},
			},
		},
		{
			name: "Non-existing key",
			traceState: TraceState{
				kvs: []label.KeyValue{
					label.String("key1", "val1"),
					label.String("key2", "val2"),
					label.String("key3", "val3"),
				},
			},
			key: "keyx",
			expectedTraceState: TraceState{
				kvs: []label.KeyValue{
					label.String("key1", "val1"),
					label.String("key2", "val2"),
					label.String("key3", "val3"),
				},
			},
		},
		{
			name: "Invalid key",
			traceState: TraceState{
				kvs: []label.KeyValue{
					label.String("key1", "val1"),
					label.String("key2", "val2"),
					label.String("key3", "val3"),
				},
			},
			key: "in va lid",
			expectedTraceState: TraceState{
				kvs: []label.KeyValue{
					label.String("key1", "val1"),
					label.String("key2", "val2"),
					label.String("key3", "val3"),
				},
			},
			expectedErr: errInvalidTraceStateKeyValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.traceState.Delete(tc.key)
			if tc.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
				assert.Equal(t, tc.traceState, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedTraceState, result)
			}
		})
	}
}

func TestTraceStateInsert(t *testing.T) {
	testCases := []struct {
		name               string
		traceState         TraceState
		keyValue           label.KeyValue
		expectedTraceState TraceState
		expectedErr        error
	}{
		{
			name: "OK case - add new",
			traceState: TraceState{
				kvs: []label.KeyValue{
					label.String("key1", "val1"),
					label.String("key2", "val2"),
					label.String("key3", "val3"),
				},
			},
			keyValue: label.String("key4@vendor", "val4"),
			expectedTraceState: TraceState{
				kvs: []label.KeyValue{
					label.String("key4@vendor", "val4"),
					label.String("key1", "val1"),
					label.String("key2", "val2"),
					label.String("key3", "val3"),
				},
			},
		},
		{
			name: "OK case - replace",
			traceState: TraceState{
				kvs: []label.KeyValue{
					label.String("key1", "val1"),
					label.String("key2", "val2"),
					label.String("key3", "val3"),
				},
			},
			keyValue: label.String("key2", "valX"),
			expectedTraceState: TraceState{
				kvs: []label.KeyValue{
					label.String("key2", "valX"),
					label.String("key1", "val1"),
					label.String("key3", "val3"),
				},
			},
		},
		{
			name: "Invalid key/value",
			traceState: TraceState{
				kvs: []label.KeyValue{
					label.String("key1", "val1"),
				},
			},
			keyValue: label.String("key!", "val!"),
			expectedTraceState: TraceState{
				kvs: []label.KeyValue{
					label.String("key1", "val1"),
				},
			},
			expectedErr: errInvalidTraceStateKeyValue,
		},
		{
			name:               "Too many entries",
			traceState:         TraceState{kvsWithMaxMembers},
			keyValue:           label.String("keyx", "valx"),
			expectedTraceState: TraceState{kvsWithMaxMembers},
			expectedErr:        errInvalidTraceStateMembersNumber,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.traceState.Insert(tc.keyValue)
			if tc.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
				assert.Equal(t, tc.traceState, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedTraceState, result)
			}
		})
	}
}

func TestTraceStateFromKeyValues(t *testing.T) {
	testCases := []struct {
		name               string
		kvs                []label.KeyValue
		expectedTraceState TraceState
		expectedErr        error
	}{
		{
			name:               "OK case",
			kvs:                kvsWithMaxMembers,
			expectedTraceState: TraceState{kvsWithMaxMembers},
		},
		{
			name:               "OK case (empty)",
			expectedTraceState: TraceState{},
		},
		{
			name: "Too many entries",
			kvs: func() []label.KeyValue {
				kvs := kvsWithMaxMembers
				kvs = append(kvs, label.String("keyx", "valX"))
				return kvs
			}(),
			expectedTraceState: TraceState{},
			expectedErr:        errInvalidTraceStateMembersNumber,
		},
		{
			name: "Duplicate",
			kvs: []label.KeyValue{
				label.String("key1", "val1"),
				label.String("key1", "val2"),
			},
			expectedTraceState: TraceState{},
			expectedErr:        errInvalidTraceStateDuplicate,
		},
		{
			name: "Invalid key/value",
			kvs: []label.KeyValue{
				label.String("key!", "val!"),
			},
			expectedTraceState: TraceState{},
			expectedErr:        errInvalidTraceStateKeyValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := TraceStateFromKeyValues(tc.kvs...)
			if tc.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, TraceState{}, result)
				assert.Equal(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, tc.expectedTraceState)
				assert.Equal(t, tc.expectedTraceState, result)
			}
		})
	}

}

func assertSpanContextEqual(got SpanContext, want SpanContext) bool {
	return got.SpanID == want.SpanID &&
		got.TraceID == want.TraceID &&
		got.TraceFlags == want.TraceFlags &&
		assertTraceStateEqual(got.TraceState, want.TraceState)
}

func assertTraceStateEqual(got TraceState, want TraceState) bool {
	if len(got.kvs) != len(want.kvs) {
		return false
	}

	for i, kv := range got.kvs {
		if kv != want.kvs[i] {
			return false
		}
	}

	return true
}

var kvsWithMaxMembers = func() []label.KeyValue {
	kvs := make([]label.KeyValue, traceStateMaxListMembers)
	for i := 0; i < traceStateMaxListMembers; i++ {
		kvs[i] = label.String(fmt.Sprintf("key%d", i+1),
			fmt.Sprintf("value%d", i+1))
	}
	return kvs
}()
