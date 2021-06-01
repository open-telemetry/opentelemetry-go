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
	"math/rand"
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

func TestNewEmptyBaggage(t *testing.T) {
	b, err := New()
	assert.NoError(t, err)
	assert.Equal(t, Baggage{}, b)
}

func TestNewBaggage(t *testing.T) {
	b, err := New(Member{key: "k"})
	assert.NoError(t, err)
	assert.Equal(t, Baggage{list: map[string]value{"k": {}}}, b)
}

func TestNewBaggageWithDuplicates(t *testing.T) {
	m := make([]Member, maxMembers+1)
	for i := range m {
		// Duplicates are collapsed.
		m[i] = Member{key: "a"}
	}
	b, err := New(m...)
	assert.NoError(t, err)
	assert.Equal(t, Baggage{list: map[string]value{"a": {}}}, b)
}

func TestNewBaggageErrorInvalidMember(t *testing.T) {
	_, err := New(Member{key: ""})
	assert.ErrorIs(t, err, errInvalidKey)
}

func key(n int) string {
	r := []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = r[rand.Intn(len(r))]
	}
	return string(b)
}

func TestNewBaggageErrorTooManyBytes(t *testing.T) {
	m := make([]Member, (maxBytesPerBaggageString/maxBytesPerMembers)+1)
	for i := range m {
		m[i] = Member{key: key(maxBytesPerMembers)}
	}
	_, err := New(m...)
	assert.ErrorIs(t, err, errBaggageBytes)
}

func TestNewBaggageErrorTooManyMembers(t *testing.T) {
	m := make([]Member, maxMembers+1)
	for i := range m {
		m[i] = Member{key: fmt.Sprintf("%d", i)}
	}
	_, err := New(m...)
	assert.ErrorIs(t, err, errMemberNumber)
}

func TestBaggageParse(t *testing.T) {
	tooLarge := key(maxBytesPerBaggageString + 1)

	tooLargeMember := key(maxBytesPerMembers + 1)

	m := make([]string, maxMembers+1)
	for i := range m {
		m[i] = fmt.Sprintf("a%d=", i)
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
			name: "URL encoded value",
			out:  "foo=1%3D1",
			baggage: map[string]value{
				"foo": {v: "1=1"},
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
			// Properties are "opaque values" meaning they are sent as they
			// are set and no encoding is performed.
			out: "foo=1;red;state=on;z=z=z",
			baggage: map[string]value{
				"foo": {
					v: "1",
					p: properties{
						{key: "state", value: "on", hasValue: true},
						{key: "red"},
						{key: "z", value: "z=z", hasValue: true},
					},
				},
			},
		},
		{
			name: "two members with properties",
			out:  "bar=2;yellow,foo=1;red;state=on",
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
		sort.Strings(members)
		return strings.Join(members, listDelimiter)
	}

	for _, tc := range testcases {
		b := Baggage{tc.baggage}
		assert.Equal(t, tc.out, orderer(b.String()))
	}
}

func TestBaggageLen(t *testing.T) {
	b := Baggage{}
	assert.Equal(t, 0, b.Len())

	b.list = make(map[string]value, 1)
	assert.Equal(t, 0, b.Len())

	b.list["k"] = value{}
	assert.Equal(t, 1, b.Len())
}

func TestBaggageDeleteMember(t *testing.T) {
	key := "k"

	b0 := Baggage{}
	b1 := b0.DeleteMember(key)
	assert.NotContains(t, b1.list, key)

	b0 = Baggage{list: map[string]value{
		key:     {},
		"other": {},
	}}
	b1 = b0.DeleteMember(key)
	assert.Contains(t, b0.list, key)
	assert.NotContains(t, b1.list, key)
}

func TestBaggageSetMemberError(t *testing.T) {
	_, err := Baggage{}.SetMember(Member{})
	assert.ErrorIs(t, err, errInvalidMember)
}

func TestBaggageSetMember(t *testing.T) {
	b0 := Baggage{}

	key := "k"
	m := Member{key: key}
	b1, err := b0.SetMember(m)
	assert.NoError(t, err)
	assert.NotContains(t, b0.list, key)
	assert.Equal(t, value{}, b1.list[key])
	assert.Equal(t, 0, len(b0.list))
	assert.Equal(t, 1, len(b1.list))

	m.value = "v"
	b2, err := b1.SetMember(m)
	assert.NoError(t, err)
	assert.Equal(t, value{}, b1.list[key])
	assert.Equal(t, value{v: "v"}, b2.list[key])
	assert.Equal(t, 1, len(b1.list))
	assert.Equal(t, 1, len(b2.list))

	p := properties{{key: "p"}}
	m.properties = p
	b3, err := b2.SetMember(m)
	assert.NoError(t, err)
	assert.Equal(t, value{v: "v"}, b2.list[key])
	assert.Equal(t, value{v: "v", p: p}, b3.list[key])
	assert.Equal(t, 1, len(b2.list))
	assert.Equal(t, 1, len(b3.list))

	// The returned baggage needs to be immutable and should use a copy of the
	// properties slice.
	p[0] = Property{key: "different"}
	assert.Equal(t, value{v: "v", p: properties{{key: "p"}}}, b3.list[key])
	// Reset for below.
	p[0] = Property{key: "p"}

	m = Member{key: "another"}
	b4, err := b3.SetMember(m)
	assert.NoError(t, err)
	assert.Equal(t, value{v: "v", p: p}, b3.list[key])
	assert.NotContains(t, b3.list, m.key)
	assert.Equal(t, value{v: "v", p: p}, b4.list[key])
	assert.Equal(t, value{}, b4.list[m.key])
	assert.Equal(t, 1, len(b3.list))
	assert.Equal(t, 2, len(b4.list))
}

func TestNilBaggageMembers(t *testing.T) {
	assert.Nil(t, Baggage{}.Members())
}

func TestBaggageMembers(t *testing.T) {
	members := []Member{
		{
			key:   "foo",
			value: "1",
			properties: properties{
				{key: "state", value: "on", hasValue: true},
				{key: "red"},
			},
		},
		{
			key:   "bar",
			value: "2",
			properties: properties{
				{key: "yellow"},
			},
		},
	}

	baggage := Baggage{list: map[string]value{
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
	}}

	assert.ElementsMatch(t, members, baggage.Members())
}

func TestBaggageMember(t *testing.T) {
	baggage := Baggage{list: map[string]value{"foo": {v: "1"}}}
	assert.Equal(t, Member{key: "foo", value: "1"}, baggage.Member("foo"))
	assert.Equal(t, Member{}, baggage.Member("bar"))
}

func TestMemberKey(t *testing.T) {
	m := Member{}
	assert.Equal(t, "", m.Key(), "even invalid values should be returned")

	key := "k"
	m.key = key
	assert.Equal(t, key, m.Key())
}

func TestMemberValue(t *testing.T) {
	m := Member{key: "k", value: "\\"}
	assert.Equal(t, "\\", m.Value(), "even invalid values should be returned")

	value := "v"
	m.value = value
	assert.Equal(t, value, m.Value())
}

func TestMemberProperties(t *testing.T) {
	m := Member{key: "k", value: "v"}
	assert.Nil(t, m.Properties())

	p := []Property{{key: "foo"}}
	m.properties = properties(p)
	got := m.Properties()
	assert.Equal(t, p, got)

	// Returned slice needs to be a copy so the original is immutable.
	got[0] = Property{key: "bar"}
	assert.NotEqual(t, m.properties, got)
}

func TestMemberValidation(t *testing.T) {
	m := Member{}
	assert.ErrorIs(t, m.validate(), errInvalidKey)

	m.key, m.value = "k", "\\"
	assert.ErrorIs(t, m.validate(), errInvalidValue)

	m.value = "v"
	assert.NoError(t, m.validate())
}

func TestNewMember(t *testing.T) {
	m, err := NewMember("", "")
	assert.ErrorIs(t, err, errInvalidKey)
	assert.Equal(t, Member{}, m)

	key, val := "k", "v"
	p := Property{key: "foo"}
	m, err = NewMember(key, val, p)
	assert.NoError(t, err)
	expected := Member{key: key, value: val, properties: properties{{key: "foo"}}}
	assert.Equal(t, expected, m)

	// Ensure new member is immutable.
	p.key = "bar"
	assert.Equal(t, expected, m)
}

func TestPropertiesValidate(t *testing.T) {
	p := properties{{}}
	assert.ErrorIs(t, p.validate(), errInvalidKey)

	p[0].key = "foo"
	assert.NoError(t, p.validate())

	p = append(p, Property{key: "bar"})
	assert.NoError(t, p.validate())
}
