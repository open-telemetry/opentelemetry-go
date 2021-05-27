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
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Taken from the W3C tests:
// https://github.com/w3c/trace-context/blob/dcd3ad9b7d6ac36f70ff3739874b73c11b0302a1/test/test_data.json
var testcases = []struct {
	in         string
	tracestate TraceState
	out        string
	err        error
}{
	{
		in:  "foo=1,foo=1",
		err: errDuplicate,
	},
	{
		in:  "foo=1,foo=2",
		err: errDuplicate,
	},
	{
		in:  "foo =1",
		err: errInvalidMember,
	},
	{
		in:  "FOO=1",
		err: errInvalidMember,
	},
	{
		in:  "foo.bar=1",
		err: errInvalidMember,
	},
	{
		in:  "foo@=1,bar=2",
		err: errInvalidMember,
	},
	{
		in:  "@foo=1,bar=2",
		err: errInvalidMember,
	},
	{
		in:  "foo@@bar=1,bar=2",
		err: errInvalidMember,
	},
	{
		in:  "foo@bar@baz=1,bar=2",
		err: errInvalidMember,
	},
	{
		in:  "foo=1,zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz=1",
		err: errInvalidMember,
	},
	{
		in:  "foo=1,tttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttt@v=1",
		err: errInvalidMember,
	},
	{
		in:  "foo=1,t@vvvvvvvvvvvvvvv=1",
		err: errInvalidMember,
	},
	{
		in:  "foo=bar=baz",
		err: errInvalidMember,
	},
	{
		in:  "foo=,bar=3",
		err: errInvalidMember,
	},
	{
		in:  "bar01=01,bar02=02,bar03=03,bar04=04,bar05=05,bar06=06,bar07=07,bar08=08,bar09=09,bar10=10,bar11=11,bar12=12,bar13=13,bar14=14,bar15=15,bar16=16,bar17=17,bar18=18,bar19=19,bar20=20,bar21=21,bar22=22,bar23=23,bar24=24,bar25=25,bar26=26,bar27=27,bar28=28,bar29=29,bar30=30,bar31=31,bar32=32,bar33=33",
		err: errMemberNumber,
	},
	{
		in:  "abcdefghijklmnopqrstuvwxyz0123456789_-*/= !\"#$%&'()*+-./0123456789:;<>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~",
		out: "abcdefghijklmnopqrstuvwxyz0123456789_-*/= !\"#$%&'()*+-./0123456789:;<>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~",
		tracestate: TraceState{list: []member{
			{
				Key:   "abcdefghijklmnopqrstuvwxyz0123456789_-*/",
				Value: " !\"#$%&'()*+-./0123456789:;<>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~",
			},
		}},
	},
	{
		in:  "abcdefghijklmnopqrstuvwxyz0123456789_-*/@a-z0-9_-*/= !\"#$%&'()*+-./0123456789:;<>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~",
		out: "abcdefghijklmnopqrstuvwxyz0123456789_-*/@a-z0-9_-*/= !\"#$%&'()*+-./0123456789:;<>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~",
		tracestate: TraceState{list: []member{
			{
				Key:   "abcdefghijklmnopqrstuvwxyz0123456789_-*/@a-z0-9_-*/",
				Value: " !\"#$%&'()*+-./0123456789:;<>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~",
			},
		}},
	},
	{
		// Empty input should result in no error and a zero value
		// TraceState being returned, that TraceState should be encoded as an
		// empty string.
	},
	{
		in:  "foo=1",
		out: "foo=1",
		tracestate: TraceState{list: []member{
			{Key: "foo", Value: "1"},
		}},
	},
	{
		in:  "foo=1,",
		out: "foo=1",
		tracestate: TraceState{list: []member{
			{Key: "foo", Value: "1"},
		}},
	},
	{
		in:  "foo=1,bar=2",
		out: "foo=1,bar=2",
		tracestate: TraceState{list: []member{
			{Key: "foo", Value: "1"},
			{Key: "bar", Value: "2"},
		}},
	},
	{
		in:  "foo=1,zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz=1",
		out: "foo=1,zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz=1",
		tracestate: TraceState{list: []member{
			{
				Key:   "foo",
				Value: "1",
			},
			{
				Key:   "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz",
				Value: "1",
			},
		}},
	},
	{
		in:  "foo=1,ttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttt@vvvvvvvvvvvvvv=1",
		out: "foo=1,ttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttt@vvvvvvvvvvvvvv=1",
		tracestate: TraceState{list: []member{
			{
				Key:   "foo",
				Value: "1",
			},
			{
				Key:   "ttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttt@vvvvvvvvvvvvvv",
				Value: "1",
			},
		}},
	},
	{
		in:  "bar01=01,bar02=02,bar03=03,bar04=04,bar05=05,bar06=06,bar07=07,bar08=08,bar09=09,bar10=10,bar11=11,bar12=12,bar13=13,bar14=14,bar15=15,bar16=16,bar17=17,bar18=18,bar19=19,bar20=20,bar21=21,bar22=22,bar23=23,bar24=24,bar25=25,bar26=26,bar27=27,bar28=28,bar29=29,bar30=30,bar31=31,bar32=32",
		out: "bar01=01,bar02=02,bar03=03,bar04=04,bar05=05,bar06=06,bar07=07,bar08=08,bar09=09,bar10=10,bar11=11,bar12=12,bar13=13,bar14=14,bar15=15,bar16=16,bar17=17,bar18=18,bar19=19,bar20=20,bar21=21,bar22=22,bar23=23,bar24=24,bar25=25,bar26=26,bar27=27,bar28=28,bar29=29,bar30=30,bar31=31,bar32=32",
		tracestate: TraceState{list: []member{
			{Key: "bar01", Value: "01"},
			{Key: "bar02", Value: "02"},
			{Key: "bar03", Value: "03"},
			{Key: "bar04", Value: "04"},
			{Key: "bar05", Value: "05"},
			{Key: "bar06", Value: "06"},
			{Key: "bar07", Value: "07"},
			{Key: "bar08", Value: "08"},
			{Key: "bar09", Value: "09"},
			{Key: "bar10", Value: "10"},
			{Key: "bar11", Value: "11"},
			{Key: "bar12", Value: "12"},
			{Key: "bar13", Value: "13"},
			{Key: "bar14", Value: "14"},
			{Key: "bar15", Value: "15"},
			{Key: "bar16", Value: "16"},
			{Key: "bar17", Value: "17"},
			{Key: "bar18", Value: "18"},
			{Key: "bar19", Value: "19"},
			{Key: "bar20", Value: "20"},
			{Key: "bar21", Value: "21"},
			{Key: "bar22", Value: "22"},
			{Key: "bar23", Value: "23"},
			{Key: "bar24", Value: "24"},
			{Key: "bar25", Value: "25"},
			{Key: "bar26", Value: "26"},
			{Key: "bar27", Value: "27"},
			{Key: "bar28", Value: "28"},
			{Key: "bar29", Value: "29"},
			{Key: "bar30", Value: "30"},
			{Key: "bar31", Value: "31"},
			{Key: "bar32", Value: "32"},
		}},
	},
	{
		in:  "foo=1,bar=2,rojo=1,congo=2,baz=3",
		out: "foo=1,bar=2,rojo=1,congo=2,baz=3",
		tracestate: TraceState{list: []member{
			{Key: "foo", Value: "1"},
			{Key: "bar", Value: "2"},
			{Key: "rojo", Value: "1"},
			{Key: "congo", Value: "2"},
			{Key: "baz", Value: "3"},
		}},
	},
	{
		in:  "foo=1 \t , \t bar=2, \t baz=3",
		out: "foo=1,bar=2,baz=3",
		tracestate: TraceState{list: []member{
			{Key: "foo", Value: "1"},
			{Key: "bar", Value: "2"},
			{Key: "baz", Value: "3"},
		}},
	},
	{
		in:  "foo=1\t \t,\t \tbar=2,\t \tbaz=3",
		out: "foo=1,bar=2,baz=3",
		tracestate: TraceState{list: []member{
			{Key: "foo", Value: "1"},
			{Key: "bar", Value: "2"},
			{Key: "baz", Value: "3"},
		}},
	},
	{
		in:  "foo=1 ",
		out: "foo=1",
		tracestate: TraceState{list: []member{
			{Key: "foo", Value: "1"},
		}},
	},
	{
		in:  "foo=1\t",
		out: "foo=1",
		tracestate: TraceState{list: []member{
			{Key: "foo", Value: "1"},
		}},
	},
	{
		in:  "foo=1 \t",
		out: "foo=1",
		tracestate: TraceState{list: []member{
			{Key: "foo", Value: "1"},
		}},
	},
}

var maxMembers = func() TraceState {
	members := make([]member, maxListMembers)
	for i := 0; i < maxListMembers; i++ {
		members[i] = member{
			Key:   fmt.Sprintf("key%d", i+1),
			Value: fmt.Sprintf("value%d", i+1),
		}
	}
	return TraceState{list: members}
}()

func TestParseTraceState(t *testing.T) {
	for _, tc := range testcases {
		got, err := ParseTraceState(tc.in)
		assert.Equal(t, tc.tracestate, got)
		if tc.err != nil {
			assert.ErrorIs(t, err, tc.err, tc.in)
		} else {
			assert.NoError(t, err, tc.in)
		}
	}
}

func TestTraceStateString(t *testing.T) {
	for _, tc := range testcases {
		if tc.err != nil {
			// Only test non-zero value TraceState.
			continue
		}

		assert.Equal(t, tc.out, tc.tracestate.String())
	}
}

func TestTraceStateMarshalJSON(t *testing.T) {
	for _, tc := range testcases {
		if tc.err != nil {
			// Only test non-zero value TraceState.
			continue
		}

		// Encode UTF-8.
		expected, err := json.Marshal(tc.out)
		require.NoError(t, err)

		actual, err := json.Marshal(tc.tracestate)
		require.NoError(t, err)

		assert.Equal(t, expected, actual)
	}
}

func TestTraceStateGet(t *testing.T) {
	testCases := []struct {
		name     string
		key      string
		expected string
	}{
		{
			name:     "OK case",
			key:      "key16",
			expected: "value16",
		},
		{
			name:     "not found",
			key:      "keyxx",
			expected: "",
		},
		{
			name:     "invalid W3C key",
			key:      "key!",
			expected: "",
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expected, maxMembers.Get(tc.key), tc.name)
	}
}

func TestTraceStateDelete(t *testing.T) {
	ts := TraceState{list: []member{
		{Key: "key1", Value: "val1"},
		{Key: "key2", Value: "val2"},
		{Key: "key3", Value: "val3"},
	}}

	testCases := []struct {
		name     string
		key      string
		expected TraceState
	}{
		{
			name: "OK case",
			key:  "key2",
			expected: TraceState{list: []member{
				{Key: "key1", Value: "val1"},
				{Key: "key3", Value: "val3"},
			}},
		},
		{
			name: "Non-existing key",
			key:  "keyx",
			expected: TraceState{list: []member{
				{Key: "key1", Value: "val1"},
				{Key: "key2", Value: "val2"},
				{Key: "key3", Value: "val3"},
			}},
		},
		{
			name: "Invalid key",
			key:  "in va lid",
			expected: TraceState{list: []member{
				{Key: "key1", Value: "val1"},
				{Key: "key2", Value: "val2"},
				{Key: "key3", Value: "val3"},
			}},
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expected, ts.Delete(tc.key), tc.name)
	}
}

func TestTraceStateInsert(t *testing.T) {
	ts := TraceState{list: []member{
		{Key: "key1", Value: "val1"},
		{Key: "key2", Value: "val2"},
		{Key: "key3", Value: "val3"},
	}}

	testCases := []struct {
		name       string
		tracestate TraceState
		key, value string
		expected   TraceState
		err        error
	}{
		{
			name:       "add new",
			tracestate: ts,
			key:        "key4@vendor",
			value:      "val4",
			expected: TraceState{list: []member{
				{Key: "key4@vendor", Value: "val4"},
				{Key: "key1", Value: "val1"},
				{Key: "key2", Value: "val2"},
				{Key: "key3", Value: "val3"},
			}},
		},
		{
			name:       "replace",
			tracestate: ts,
			key:        "key2",
			value:      "valX",
			expected: TraceState{list: []member{
				{Key: "key2", Value: "valX"},
				{Key: "key1", Value: "val1"},
				{Key: "key3", Value: "val3"},
			}},
		},
		{
			name:       "invalid key",
			tracestate: ts,
			key:        "key!",
			value:      "val",
			expected:   ts,
			err:        errInvalidKey,
		},
		{
			name:       "invalid value",
			tracestate: ts,
			key:        "key",
			value:      "v=l",
			expected:   ts,
			err:        errInvalidValue,
		},
		{
			name:       "invalid key/value",
			tracestate: ts,
			key:        "key!",
			value:      "v=l",
			expected:   ts,
			err:        errInvalidKey,
		},
		{
			name:       "too many entries",
			tracestate: maxMembers,
			key:        "keyx",
			value:      "valx",
			expected:   maxMembers,
			err:        errMemberNumber,
		},
	}

	for _, tc := range testCases {
		actual, err := tc.tracestate.Insert(tc.key, tc.value)
		assert.ErrorIs(t, err, tc.err, tc.name)
		if tc.err != nil {
			assert.Equal(t, tc.tracestate, actual, tc.name)
		} else {
			assert.Equal(t, tc.expected, actual, tc.name)
		}
	}
}

func TestTraceStateLen(t *testing.T) {
	ts := TraceState{}
	assert.Equal(t, 0, ts.Len(), "zero value TraceState is empty")

	key := "key"
	ts = TraceState{list: []member{{key, "value"}}}
	assert.Equal(t, 1, ts.Len(), "TraceState with one value")
}

func TestTraceStateImmutable(t *testing.T) {
	k0, v0 := "k0", "v0"
	ts0 := TraceState{list: []member{{k0, v0}}}
	assert.Equal(t, v0, ts0.Get(k0))

	// Insert should not modify the original.
	k1, v1 := "k1", "v1"
	ts1, err := ts0.Insert(k1, v1)
	require.NoError(t, err)
	assert.Equal(t, v0, ts0.Get(k0))
	assert.Equal(t, "", ts0.Get(k1))
	assert.Equal(t, v0, ts1.Get(k0))
	assert.Equal(t, v1, ts1.Get(k1))

	// Update should not modify the original.
	v2 := "v2"
	ts2, err := ts1.Insert(k1, v2)
	require.NoError(t, err)
	assert.Equal(t, v0, ts0.Get(k0))
	assert.Equal(t, "", ts0.Get(k1))
	assert.Equal(t, v0, ts1.Get(k0))
	assert.Equal(t, v1, ts1.Get(k1))
	assert.Equal(t, v0, ts2.Get(k0))
	assert.Equal(t, v2, ts2.Get(k1))

	// Delete should not modify the original.
	ts3 := ts2.Delete(k0)
	assert.Equal(t, v0, ts0.Get(k0))
	assert.Equal(t, v0, ts1.Get(k0))
	assert.Equal(t, v0, ts2.Get(k0))
	assert.Equal(t, "", ts3.Get(k0))
}
