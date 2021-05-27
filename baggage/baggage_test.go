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

package baggage

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyRegExp(t *testing.T) {
	// ASCII only
	invalidKeyRune := []rune{
		'\x00', '\x01', '\x02', '\x03', '\x04', '\x05', '\x06', '\x07',
		'\x08', '\x09', '\x0A', '\x0B', '\x0C', '\x0D', '\x0E', '\x0F',
		'\x10', '\x11', '\x12', '\x13', '\x14', '\x15', '\x16', '\x17',
		'\x18', '\x19', '\x1A', '\x1B', '\x1C', '\x1D', '\x1E', '\x1F', ' ',
		'(', ')', '<', '>', '@', ',', ';', ':', '\\', '"', '/', '[', ']', '?',
		'=', '{', '}', '\x7F',
	}

	for _, ch := range invalidKeyRune {
		assert.NotRegexp(t, keyDef, fmt.Sprintf("%c", ch))
	}
}

func TestValueRegExp(t *testing.T) {
	// ASCII only
	invalidValueRune := []rune{
		'\x00', '\x01', '\x02', '\x03', '\x04', '\x05', '\x06', '\x07',
		'\x08', '\x09', '\x0A', '\x0B', '\x0C', '\x0D', '\x0E', '\x0F',
		'\x10', '\x11', '\x12', '\x13', '\x14', '\x15', '\x16', '\x17',
		'\x18', '\x19', '\x1A', '\x1B', '\x1C', '\x1D', '\x1E', '\x1F', ' ',
		'"', ',', ';', '\\', '\x7F',
	}

	for _, ch := range invalidValueRune {
		assert.NotRegexp(t, `^`+valueDef+`$`, fmt.Sprintf("invalid-%c-value", ch))
	}
}

func TestParseProperty(t *testing.T) {
	p := Property{key: "key", value: "value", hasValue: true}

	testcases := []struct {
		in       string
		expected Property
	}{
		{
			in:       "",
			expected: Property{},
		},
		{
			in: "key",
			expected: Property{
				key: "key",
			},
		},
		{
			in: "key=",
			expected: Property{
				key:      "key",
				hasValue: true,
			},
		},
		{
			in:       "key=value",
			expected: p,
		},
		{
			in:       " key=value ",
			expected: p,
		},
		{
			in:       "key = value",
			expected: p,
		},
		{
			in:       " key = value ",
			expected: p,
		},
		{
			in:       "\tkey=value",
			expected: p,
		},
	}

	for _, tc := range testcases {
		actual, err := parseProperty(tc.in)

		if !assert.NoError(t, err) {
			continue
		}

		assert.Equal(t, tc.expected.Key(), actual.Key(), tc.in)

		actualV, actualOk := actual.Value()
		expectedV, expectedOk := tc.expected.Value()
		assert.Equal(t, expectedOk, actualOk, tc.in)
		assert.Equal(t, expectedV, actualV, tc.in)
	}
}

func TestParsePropertyError(t *testing.T) {
	_, err := parseProperty(",;,")
	assert.ErrorIs(t, err, errInvalidProperty)
}

func TestNewKeyProperty(t *testing.T) {
	p, err := NewKeyProperty(" ")
	assert.ErrorIs(t, err, errInvalidKey)
	assert.Equal(t, Property{}, p)

	p, err = NewKeyProperty("key")
	assert.NoError(t, err)
	assert.Equal(t, Property{key: "key"}, p)
}

func TestNewKeyValueProperty(t *testing.T) {
	p, err := NewKeyValueProperty(" ", "")
	assert.ErrorIs(t, err, errInvalidKey)
	assert.Equal(t, Property{}, p)

	p, err = NewKeyValueProperty("key", ";")
	assert.ErrorIs(t, err, errInvalidValue)
	assert.Equal(t, Property{}, p)

	p, err = NewKeyValueProperty("key", "value")
	assert.NoError(t, err)
	assert.Equal(t, Property{key: "key", value: "value", hasValue: true}, p)
}

func TestPropertyValidate(t *testing.T) {
	p := Property{}
	assert.ErrorIs(t, p.validate(), errInvalidKey)

	p.key = "k"
	assert.NoError(t, p.validate())

	p.value = ";"
	assert.EqualError(t, p.validate(), "invalid property: inconsistent value")

	p.hasValue = true
	assert.ErrorIs(t, p.validate(), errInvalidValue)

	p.value = "v"
	assert.NoError(t, p.validate())
}

func TestBaggageParse(t *testing.T) {
	b := make([]rune, maxBytesPerBaggageString+1)
	for i := range b {
		b[i] = 'a'
	}
	tooLarge := string(b)

	b = make([]rune, maxBytesPerMembers+1)
	for i := range b {
		b[i] = 'a'
	}
	tooLargeMember := string(b)

	m := make([]string, maxMembers+1)
	for i := range m {
		m[i] = "a="
	}
	tooManyMembers := strings.Join(m, listDelimiter)

	testcases := []struct {
		name    string
		in      string
		baggage map[string]value
		err     error
	}{
		{
			name:    "empty value",
			in:      "",
			baggage: map[string]value(nil),
		},
		{
			name: "single member empty value no properties",
			in:   "foo=",
			baggage: map[string]value{
				"foo": {v: ""},
			},
		},
		{
			name: "single member no properties",
			in:   "foo=1",
			baggage: map[string]value{
				"foo": {v: "1"},
			},
		},
		{
			name: "single member empty value with properties",
			in:   "foo=;state=on;red",
			baggage: map[string]value{
				"foo": {
					v: "",
					p: properties{
						{key: "state", value: "on", hasValue: true},
						{key: "red"},
					},
				},
			},
		},
		{
			name: "single member with properties",
			in:   "foo=1;state=on;red",
			baggage: map[string]value{
				"foo": {
					v: "1",
					p: properties{
						{key: "state", value: "on", hasValue: true},
						{key: "red"},
					},
				},
			},
		},
		{
			name: "two members with properties",
			in:   "foo=1;state=on;red,bar=2;yellow",
			baggage: map[string]value{
				"foo": {
					v: "1",
					p: properties{
						{key: "state", value: "on", hasValue: true},
						{key: "red"},
					},
				},
				"bar": {
					v: "2",
					p: properties{{key: "yellow"}},
				},
			},
		},
		{
			// According to the OTel spec, last value wins.
			name: "duplicate key",
			in:   "foo=1;state=on;red,foo=2",
			baggage: map[string]value{
				"foo": {v: "2"},
			},
		},
		{
			name: "invalid member: empty key",
			in:   "foo=,,bar=",
			err:  errInvalidMember,
		},
		{
			name: "invalid member: invalid value",
			in:   "foo=\\",
			err:  errInvalidMember,
		},
		{
			name: "invalid property: invalid key",
			in:   "foo=1;=v",
			err:  errInvalidProperty,
		},
		{
			name: "invalid property: invalid value",
			in:   "foo=1;key=\\",
			err:  errInvalidProperty,
		},
		{
			name: "invalid baggage string: too large",
			in:   tooLarge,
			err:  errBaggageBytes,
		},
		{
			name: "invalid baggage string: member too large",
			in:   tooLargeMember,
			err:  errMemberBytes,
		},
		{
			name: "invalid baggage string: too many members",
			in:   tooManyMembers,
			err:  errMemberNumber,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := Parse(tc.in)
			assert.ErrorIs(t, err, tc.err)
			assert.Equal(t, Baggage{list: tc.baggage}, actual)
		})
	}
}

func TestBaggageString(t *testing.T) {
	testcases := []struct {
		name    string
		out     string
		baggage map[string]value
	}{
		{
			name:    "empty value",
			out:     "",
			baggage: map[string]value(nil),
		},
		{
			name: "single member empty value no properties",
			out:  "foo=",
			baggage: map[string]value{
				"foo": {v: ""},
			},
		},
		{
			name: "single member no properties",
			out:  "foo=1",
			baggage: map[string]value{
				"foo": {v: "1"},
			},
		},
		{
			name: "single member empty value with properties",
			out:  "foo=;red;state=on",
			baggage: map[string]value{
				"foo": {
					v: "",
					p: properties{
						{key: "state", value: "on", hasValue: true},
						{key: "red"},
					},
				},
			},
		},
		{
			name: "single member with properties",
			out:  "foo=1;red;state=on",
			baggage: map[string]value{
				"foo": {
					v: "1",
					p: properties{
						{key: "state", value: "on", hasValue: true},
						{key: "red"},
					},
				},
			},
		},
		{
			name: "two members with properties",
			out:  "foo=1;red;state=on,bar=2;yellow",
			baggage: map[string]value{
				"foo": {
					v: "1",
					p: properties{
						{key: "state", value: "on", hasValue: true},
						{key: "red"},
					},
				},
				"bar": {
					v: "2",
					p: properties{{key: "yellow"}},
				},
			},
		},
	}

	orderer := func(s string) string {
		members := strings.Split(s, listDelimiter)
		for i, m := range members {
			parts := strings.Split(m, propertyDelimiter)
			if len(parts) > 1 {
				sort.Strings(parts[1:])
				members[i] = strings.Join(parts, propertyDelimiter)
			}
		}
		return strings.Join(members, listDelimiter)
	}

	for _, tc := range testcases {
		b := Baggage{tc.baggage}
		assert.Equal(t, tc.out, orderer(b.String()))
	}
}
