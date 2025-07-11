// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package baggage

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/internal/baggage"
)

// Seed with a static value to ensure deterministic results.
var rng = rand.New(rand.NewChaCha8([32]byte{}))

func TestValidateKeyChar(t *testing.T) {
	// ASCII only
	invalidKeyRune := []rune{
		'\x00', '\x01', '\x02', '\x03', '\x04', '\x05', '\x06', '\x07',
		'\x08', '\x09', '\x0A', '\x0B', '\x0C', '\x0D', '\x0E', '\x0F',
		'\x10', '\x11', '\x12', '\x13', '\x14', '\x15', '\x16', '\x17',
		'\x18', '\x19', '\x1A', '\x1B', '\x1C', '\x1D', '\x1E', '\x1F', ' ',
		'(', ')', '<', '>', '@', ',', ';', ':', '\\', '"', '/', '[', ']', '?',
		'=', '{', '}', '\x7F', 2 >> 20, '\x80',
	}

	for _, ch := range invalidKeyRune {
		assert.False(t, validateKeyChar(ch))
	}
}

func TestValidateValueChar(t *testing.T) {
	// ASCII only
	invalidValueRune := []rune{
		'\x00', '\x01', '\x02', '\x03', '\x04', '\x05', '\x06', '\x07',
		'\x08', '\x09', '\x0A', '\x0B', '\x0C', '\x0D', '\x0E', '\x0F',
		'\x10', '\x11', '\x12', '\x13', '\x14', '\x15', '\x16', '\x17',
		'\x18', '\x19', '\x1A', '\x1B', '\x1C', '\x1D', '\x1E', '\x1F', ' ',
		'"', ',', ';', '\\', '\x7F', '\x80',
	}

	for _, ch := range invalidValueRune {
		assert.False(t, validateValueChar(ch))
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
	assert.NoError(t, err)
	assert.Equal(t, Property{key: " "}, p)

	p, err = NewKeyProperty("key")
	assert.NoError(t, err)
	assert.Equal(t, Property{key: "key"}, p)

	// UTF-8 key
	p, err = NewKeyProperty("B% 💼")
	assert.NoError(t, err)
	assert.Equal(t, Property{key: "B% 💼"}, p)

	// Invalid UTF-8 key
	p, err = NewKeyProperty(string([]byte{255}))
	assert.ErrorIs(t, err, errInvalidKey)
	assert.Equal(t, Property{}, p)
}

func TestNewKeyValueProperty(t *testing.T) {
	p, err := NewKeyValueProperty(" ", "value")
	assert.ErrorIs(t, err, errInvalidKey)
	assert.Equal(t, Property{}, p)

	p, err = NewKeyValueProperty("key", ";")
	assert.ErrorIs(t, err, errInvalidValue)
	assert.Equal(t, Property{}, p)

	// it won't use percent decoding for key
	p, err = NewKeyValueProperty("%zzzzz", "value")
	assert.NoError(t, err)
	assert.Equal(t, Property{key: "%zzzzz", value: "value", hasValue: true}, p)

	// wrong value with wrong decoding
	p, err = NewKeyValueProperty("key", "%zzzzz")
	assert.ErrorIs(t, err, errInvalidValue)
	assert.Equal(t, Property{}, p)

	p, err = NewKeyValueProperty("key", "value")
	assert.NoError(t, err)
	assert.Equal(t, Property{key: "key", value: "value", hasValue: true}, p)

	// Percent-encoded value
	p, err = NewKeyValueProperty("key", "%C4%85%C5%9B%C4%87")
	assert.NoError(t, err)
	assert.Equal(t, Property{key: "key", value: "ąść", hasValue: true}, p)
}

func TestNewKeyValuePropertyRaw(t *testing.T) {
	// Empty key
	p, err := NewKeyValuePropertyRaw("", " ")
	assert.ErrorIs(t, err, errInvalidKey)
	assert.Equal(t, Property{}, p)

	// Empty value
	// Empty string is also a valid UTF-8 string.
	p, err = NewKeyValuePropertyRaw(" ", "")
	assert.NoError(t, err)
	assert.Equal(t, Property{key: " ", hasValue: true}, p)

	// Space value
	p, err = NewKeyValuePropertyRaw(" ", " ")
	assert.NoError(t, err)
	assert.Equal(t, Property{key: " ", value: " ", hasValue: true}, p)

	p, err = NewKeyValuePropertyRaw("B% 💼", "Witaj Świecie!")
	assert.NoError(t, err)
	assert.Equal(t, Property{key: "B% 💼", value: "Witaj Świecie!", hasValue: true}, p)
}

func TestPropertyValidate(t *testing.T) {
	p := Property{}
	assert.ErrorIs(t, p.validate(), errInvalidKey)

	p.key = "k"
	assert.NoError(t, p.validate())

	p.value = "v"
	assert.EqualError(t, p.validate(), "invalid property: inconsistent value")

	p.hasValue = true
	assert.NoError(t, p.validate())

	// Invalid value
	p.value = string([]byte{255})
	assert.ErrorIs(t, p.validate(), errInvalidValue)
}

func TestNewEmptyBaggage(t *testing.T) {
	b, err := New()
	assert.NoError(t, err)
	assert.Equal(t, Baggage{}, b)
}

func TestNewBaggage(t *testing.T) {
	b, err := New(Member{key: "k", hasData: true})
	assert.NoError(t, err)
	assert.Equal(t, Baggage{list: baggage.List{"k": {}}}, b)
}

func TestNewBaggageWithDuplicates(t *testing.T) {
	// Having this many members would normally cause this to error, but since
	// these are duplicates of the same key they will be collapsed into a
	// single entry.
	m := make([]Member, maxMembers+1)
	for i := range m {
		// Duplicates are collapsed.
		m[i] = Member{
			key:     "a",
			value:   fmt.Sprintf("%d", i),
			hasData: true,
		}
	}
	b, err := New(m...)
	assert.NoError(t, err)

	// Ensure that the last-one-wins by verifying the value.
	v := fmt.Sprintf("%d", maxMembers)
	want := Baggage{list: baggage.List{"a": {Value: v}}}
	assert.Equal(t, want, b)
}

func TestNewBaggageErrorEmptyMember(t *testing.T) {
	_, err := New(Member{})
	assert.ErrorIs(t, err, errInvalidMember)
}

func key(n int) string {
	r := []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = r[rng.IntN(len(r))]
	}
	return string(b)
}

func TestNewBaggageErrorTooManyBytes(t *testing.T) {
	m := make([]Member, (maxBytesPerBaggageString/maxBytesPerMembers)+1)
	for i := range m {
		m[i] = Member{key: key(maxBytesPerMembers), hasData: true}
	}
	_, err := New(m...)
	assert.ErrorIs(t, err, errBaggageBytes)
}

func TestNewBaggageErrorTooManyMembers(t *testing.T) {
	m := make([]Member, maxMembers+1)
	for i := range m {
		m[i] = Member{key: fmt.Sprintf("%d", i), hasData: true}
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
		name string
		in   string
		want baggage.List
		err  error
	}{
		{
			name: "empty value",
			in:   "",
			want: baggage.List(nil),
		},
		{
			name: "single member empty value no properties",
			in:   "foo=",
			want: baggage.List{
				"foo": {Value: ""},
			},
		},
		{
			name: "single member no properties",
			in:   "foo=1",
			want: baggage.List{
				"foo": {Value: "1"},
			},
		},
		{
			name: "single member no properties plus",
			in:   "foo=1+1",
			want: baggage.List{
				"foo": {Value: "1+1"},
			},
		},
		{
			name: "single member no properties plus encoded",
			in:   "foo=1%2B1",
			want: baggage.List{
				"foo": {Value: "1+1"},
			},
		},
		{
			name: "single member no properties slash",
			in:   "foo=1/1",
			want: baggage.List{
				"foo": {Value: "1/1"},
			},
		},
		{
			name: "single member no properties slash encoded",
			in:   "foo=1%2F1",
			want: baggage.List{
				"foo": {Value: "1/1"},
			},
		},
		{
			name: "single member no properties equals",
			in:   "foo=1=1",
			want: baggage.List{
				"foo": {Value: "1=1"},
			},
		},
		{
			name: "single member no properties equals encoded",
			in:   "foo=1%3D1",
			want: baggage.List{
				"foo": {Value: "1=1"},
			},
		},
		{
			name: "single member with spaces",
			in:   " foo \t= 1\t\t ",
			want: baggage.List{
				"foo": {Value: "1"},
			},
		},
		{
			name: "single member empty value with properties",
			in:   "foo=;state=on;red",
			want: baggage.List{
				"foo": {
					Value: "",
					Properties: []baggage.Property{
						{Key: "state", Value: "on", HasValue: true},
						{Key: "red"},
					},
				},
			},
		},
		{
			name: "single member with properties",
			in:   "foo=1;state=on;red",
			want: baggage.List{
				"foo": {
					Value: "1",
					Properties: []baggage.Property{
						{Key: "state", Value: "on", HasValue: true},
						{Key: "red"},
					},
				},
			},
		},
		{
			name: "single member with value containing equal signs",
			in:   "foo=0=0=0",
			want: baggage.List{
				"foo": {Value: "0=0=0"},
			},
		},
		{
			name: "two members with properties",
			in:   "foo=1;state=on;red,bar=2;yellow",
			want: baggage.List{
				"foo": {
					Value: "1",
					Properties: []baggage.Property{
						{Key: "state", Value: "on", HasValue: true},
						{Key: "red"},
					},
				},
				"bar": {
					Value:      "2",
					Properties: []baggage.Property{{Key: "yellow"}},
				},
			},
		},
		{
			// According to the OTel spec, last value wins.
			name: "duplicate key",
			in:   "foo=1;state=on;red,foo=2",
			want: baggage.List{
				"foo": {Value: "2"},
			},
		},
		{
			name: "= value",
			in:   "key==",
			want: baggage.List{
				"key": {Value: "="},
			},
		},
		{
			name: "encoded ASCII string",
			in:   "key1=val%252%2C",
			want: baggage.List{
				"key1": {Value: "val%2,"},
			},
		},
		{
			name: "encoded property",
			in:   "key1=;bar=val%252%2C",
			want: baggage.List{
				"key1": {
					Properties: []baggage.Property{{Key: "bar", HasValue: true, Value: "val%2,"}},
				},
			},
		},
		{
			name: "encoded UTF-8 string",
			in:   "foo=%C4%85%C5%9B%C4%87",
			want: baggage.List{
				"foo": {Value: "ąść"},
			},
		},
		{
			name: "encoded UTF-8 string in key",
			in:   "a=b,%C4%85%C5%9B%C4%87=%C4%85%C5%9B%C4%87",
			want: baggage.List{
				"a":                  {Value: "b"},
				"%C4%85%C5%9B%C4%87": {Value: "ąść"},
			},
		},
		{
			name: "encoded UTF-8 string in property",
			in:   "a=b,%C4%85%C5%9B%C4%87=%C4%85%C5%9B%C4%87;%C4%85%C5%9B%C4%87=%C4%85%C5%9B%C4%87",
			want: baggage.List{
				"a": {Value: "b"},
				"%C4%85%C5%9B%C4%87": {Value: "ąść", Properties: []baggage.Property{
					{Key: "%C4%85%C5%9B%C4%87", HasValue: true, Value: "ąść"},
				}},
			},
		},
		{
			name: "invalid member: empty",
			in:   "foo=,,bar=",
			err:  errInvalidMember,
		},
		{
			name: "invalid member: no key",
			in:   "=foo",
			err:  errInvalidKey,
		},
		{
			name: "invalid member: no value",
			in:   "foo",
			err:  errInvalidMember,
		},
		{
			name: "invalid member: invalid key",
			in:   "\\=value",
			err:  errInvalidKey,
		},
		{
			name: "invalid member: invalid value",
			in:   "foo=\\",
			err:  errInvalidValue,
		},
		{
			name: "invalid member: improper url encoded value",
			in:   "key1=val%",
			err:  errInvalidValue,
		},
		{
			name: "invalid property: no key",
			in:   "foo=1;=v",
			err:  errInvalidProperty,
		},
		{
			name: "invalid property: invalid key",
			in:   "foo=1;key\\=v",
			err:  errInvalidProperty,
		},
		{
			name: "invalid property: invalid value",
			in:   "foo=1;key=\\",
			err:  errInvalidProperty,
		},
		{
			name: "invalid property: improper url encoded value",
			in:   "foo=1;key=val%",
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
		{
			name: "percent-encoded octet sequences do not match the UTF-8 encoding scheme",
			in:   "k=aa%ffcc;p=d%fff",
			want: baggage.List{
				"k": {
					Value: "aa�cc",
					Properties: []baggage.Property{
						{Key: "p", Value: "d�f", HasValue: true},
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := Parse(tc.in)
			assert.ErrorIs(t, err, tc.err)
			assert.Equal(t, Baggage{list: tc.want}, actual)
		})
	}
}

func TestBaggageParseValue(t *testing.T) {
	testcases := []struct {
		name          string
		in            string
		valueWant     string
		valueWantSize int
	}{
		{
			name:          "percent encoded octet sequence matches UTF-8 encoding scheme",
			in:            "k=aa%26cc",
			valueWant:     "aa&cc",
			valueWantSize: 5,
		},
		{
			name:          "percent encoded octet sequence doesn't match UTF-8 encoding scheme",
			in:            "k=aa%ffcc",
			valueWant:     "aa�cc",
			valueWantSize: 7,
		},
		{
			name:          "multiple percent encoded octet sequences don't match UTF-8 encoding scheme",
			in:            "k=aa%ffcc%fedd%fa",
			valueWant:     "aa�cc�dd�",
			valueWantSize: 15,
		},
		{
			name:          "raw value",
			in:            "k=aacc",
			valueWant:     "aacc",
			valueWantSize: 4,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := Parse(tc.in)
			assert.NoError(t, err)

			val := b.Members()[0].Value()

			assert.Equal(t, tc.valueWant, val)
			assert.Equal(t, len(val), tc.valueWantSize)
			assert.True(t, utf8.ValidString(val))
		})
	}
}

func TestBaggageString(t *testing.T) {
	testcases := []struct {
		name    string
		out     string
		baggage baggage.List
	}{
		{
			name:    "empty value",
			out:     "",
			baggage: baggage.List(nil),
		},
		{
			name: "single member empty value no properties",
			out:  "foo=",
			baggage: baggage.List{
				"foo": {Value: ""},
			},
		},
		{
			name: "single member no properties",
			out:  "foo=1",
			baggage: baggage.List{
				"foo": {Value: "1"},
			},
		},
		{
			name: "Encoded value",
			// Allowed value characters are:
			//
			//   %x21 / %x23-2B / %x2D-3A / %x3C-5B / %x5D-7E
			//
			// Meaning, US-ASCII characters excluding CTLs, whitespace,
			// DQUOTE, comma, semicolon, and backslash. All excluded
			// characters need to be percent encoded.
			out: "foo=%00%01%02%03%04%05%06%07%08%09%0A%0B%0C%0D%0E%0F%10%11%12%13%14%15%16%17%18%19%1A%1B%1C%1D%1E%1F%20!%22#$%25&'()*+%2C-./0123456789:%3B<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[%5C]^_`abcdefghijklmnopqrstuvwxyz{|}~%7F",
			baggage: baggage.List{
				"foo": {Value: func() string {
					// All US-ASCII characters.
					b := [128]byte{}
					for i := range b {
						b[i] = byte(i)
					}
					return string(b[:])
				}()},
			},
		},
		{
			name: "non-ASCII UTF-8 string",
			out:  "foo=%C4%85%C5%9B%C4%87",
			baggage: baggage.List{
				"foo": {Value: "ąść"},
			},
		},
		{
			name: "Encoded property value",
			out:  "foo=;bar=%20",
			baggage: baggage.List{
				"foo": {
					Properties: []baggage.Property{
						{Key: "bar", Value: " ", HasValue: true},
					},
				},
			},
		},
		{
			name: "plus",
			out:  "foo=1+1",
			baggage: baggage.List{
				"foo": {Value: "1+1"},
			},
		},
		{
			name: "equal",
			out:  "foo=1=1",
			baggage: baggage.List{
				"foo": {Value: "1=1"},
			},
		},
		{
			name: "single member empty value with properties",
			out:  "foo=;red;state=on",
			baggage: baggage.List{
				"foo": {
					Value: "",
					Properties: []baggage.Property{
						{Key: "state", Value: "on", HasValue: true},
						{Key: "red"},
					},
				},
			},
		},
		{
			name: "single member with properties",
			// Properties are "opaque values" meaning they are sent as they
			// are set and no encoding is performed.
			out: "foo=1;red;state=on;z=z=z",
			baggage: baggage.List{
				"foo": {
					Value: "1",
					Properties: []baggage.Property{
						{Key: "state", Value: "on", HasValue: true},
						{Key: "red"},
						{Key: "z", Value: "z=z", HasValue: true},
					},
				},
			},
		},
		{
			name: "two members with properties",
			out:  "bar=2;yellow,foo=1;red;state=on",
			baggage: baggage.List{
				"foo": {
					Value: "1",
					Properties: []baggage.Property{
						{Key: "state", Value: "on", HasValue: true},
						{Key: "red"},
					},
				},
				"bar": {
					Value:      "2",
					Properties: []baggage.Property{{Key: "yellow"}},
				},
			},
		},
		{
			// W3C does not allow percent-encoded keys.
			// The baggage that has percent-encoded keys should be ignored.
			name: "utf-8 key and value",
			out:  "foo=B%25%20%F0%9F%92%BC-2;foo-1=B%25%20%F0%9F%92%BC-4;foo-2",
			baggage: baggage.List{
				"ąść": {
					Value: "B% 💼",
					Properties: []baggage.Property{
						{Key: "ąść-1", Value: "B% 💼-1", HasValue: true},
						{Key: "ąść-2"},
					},
				},
				"foo": {
					Value: "B% 💼-2",
					Properties: []baggage.Property{
						{Key: "ąść", Value: "B% 💼-3", HasValue: true},
						{Key: "foo-1", Value: "B% 💼-4", HasValue: true},
						{Key: "foo-2"},
					},
				},
			},
		},
	}

	orderer := func(s string) string {
		members := strings.Split(s, listDelimiter)
		for i, m := range members {
			parts := strings.Split(m, propertyDelimiter)
			if len(parts) > 1 {
				slices.Sort(parts[1:])
				members[i] = strings.Join(parts, propertyDelimiter)
			}
		}
		slices.Sort(members)
		return strings.Join(members, listDelimiter)
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			b := Baggage{tc.baggage}
			assert.Equal(t, tc.out, orderer(b.String()))
		})
	}
}

func TestBaggageLen(t *testing.T) {
	b := Baggage{}
	assert.Equal(t, 0, b.Len())

	b.list = make(baggage.List, 1)
	assert.Equal(t, 0, b.Len())

	b.list["k"] = baggage.Item{}
	assert.Equal(t, 1, b.Len())
}

func TestBaggageDeleteMember(t *testing.T) {
	key := "k"

	b0 := Baggage{}
	b1 := b0.DeleteMember(key)
	assert.NotContains(t, b1.list, key)

	b0 = Baggage{list: baggage.List{
		key:     {},
		"other": {},
	}}
	b1 = b0.DeleteMember(key)
	assert.Contains(t, b0.list, key)
	assert.NotContains(t, b1.list, key)
}

func TestBaggageSetMemberEmpty(t *testing.T) {
	_, err := Baggage{}.SetMember(Member{})
	assert.ErrorIs(t, err, errInvalidMember)
}

func TestBaggageSetMember(t *testing.T) {
	b0 := Baggage{}

	key := "k"
	m := Member{key: key, hasData: true}
	b1, err := b0.SetMember(m)
	assert.NoError(t, err)
	assert.NotContains(t, b0.list, key)
	assert.Equal(t, baggage.Item{}, b1.list[key])
	assert.Empty(t, b0.list)
	assert.Len(t, b1.list, 1)

	m.value = "v"
	b2, err := b1.SetMember(m)
	assert.NoError(t, err)
	assert.Equal(t, baggage.Item{}, b1.list[key])
	assert.Equal(t, baggage.Item{Value: "v"}, b2.list[key])
	assert.Len(t, b1.list, 1)
	assert.Len(t, b2.list, 1)

	p := properties{{key: "p"}}
	m.properties = p
	b3, err := b2.SetMember(m)
	assert.NoError(t, err)
	assert.Equal(t, baggage.Item{Value: "v"}, b2.list[key])
	assert.Equal(t, baggage.Item{Value: "v", Properties: []baggage.Property{{Key: "p"}}}, b3.list[key])
	assert.Len(t, b2.list, 1)
	assert.Len(t, b3.list, 1)

	// The returned baggage needs to be immutable and should use a copy of the
	// properties slice.
	p[0] = Property{key: "different"}
	assert.Equal(t, baggage.Item{Value: "v", Properties: []baggage.Property{{Key: "p"}}}, b3.list[key])
	// Reset for below.
	p[0] = Property{key: "p"}

	m = Member{key: "another", hasData: true}
	b4, err := b3.SetMember(m)
	assert.NoError(t, err)
	assert.Equal(t, baggage.Item{Value: "v", Properties: []baggage.Property{{Key: "p"}}}, b3.list[key])
	assert.NotContains(t, b3.list, m.key)
	assert.Equal(t, baggage.Item{Value: "v", Properties: []baggage.Property{{Key: "p"}}}, b4.list[key])
	assert.Equal(t, baggage.Item{}, b4.list[m.key])
	assert.Len(t, b3.list, 1)
	assert.Len(t, b4.list, 2)
}

func TestBaggageSetFalseMember(t *testing.T) {
	b0 := Baggage{}

	key := "k"
	m := Member{key: key, hasData: false}
	b1, err := b0.SetMember(m)
	assert.Error(t, err)
	assert.NotContains(t, b0.list, key)
	assert.Equal(t, baggage.Item{}, b1.list[key])
	assert.Empty(t, b0.list)
	assert.Empty(t, b1.list)

	m.value = "v"
	b2, err := b1.SetMember(m)
	assert.Error(t, err)
	assert.Equal(t, baggage.Item{}, b1.list[key])
	assert.Equal(t, baggage.Item{Value: ""}, b2.list[key])
	assert.Empty(t, b1.list)
	assert.Empty(t, b2.list)
}

func TestBaggageSetFalseMembers(t *testing.T) {
	b0 := Baggage{}

	key := "k"
	m := Member{key: key, hasData: true}
	b1, err := b0.SetMember(m)
	assert.NoError(t, err)
	assert.NotContains(t, b0.list, key)
	assert.Equal(t, baggage.Item{}, b1.list[key])
	assert.Empty(t, b0.list)
	assert.Len(t, b1.list, 1)

	m.value = "v"
	b2, err := b1.SetMember(m)
	assert.NoError(t, err)
	assert.Equal(t, baggage.Item{}, b1.list[key])
	assert.Equal(t, baggage.Item{Value: "v"}, b2.list[key])
	assert.Len(t, b1.list, 1)
	assert.Len(t, b2.list, 1)

	p := properties{{key: "p"}}
	m.properties = p
	b3, err := b2.SetMember(m)
	assert.NoError(t, err)
	assert.Equal(t, baggage.Item{Value: "v"}, b2.list[key])
	assert.Equal(t, baggage.Item{Value: "v", Properties: []baggage.Property{{Key: "p"}}}, b3.list[key])
	assert.Len(t, b2.list, 1)
	assert.Len(t, b3.list, 1)

	// The returned baggage needs to be immutable and should use a copy of the
	// properties slice.
	p[0] = Property{key: "different"}
	assert.Equal(t, baggage.Item{Value: "v", Properties: []baggage.Property{{Key: "p"}}}, b3.list[key])
	// Reset for below.
	p[0] = Property{key: "p"}

	m = Member{key: "another"}
	b4, err := b3.SetMember(m)
	assert.Error(t, err)
	assert.Equal(t, baggage.Item{Value: "v", Properties: []baggage.Property{{Key: "p"}}}, b3.list[key])
	assert.NotContains(t, b3.list, m.key)
	assert.Equal(t, baggage.Item{Value: "v", Properties: []baggage.Property{{Key: "p"}}}, b4.list[key])
	assert.Equal(t, baggage.Item{}, b4.list[m.key])
	assert.Len(t, b3.list, 1)
	assert.Len(t, b4.list, 1)
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
			hasData: true,
		},
		{
			key:   "bar",
			value: "2",
			properties: properties{
				{key: "yellow"},
			},
			hasData: true,
		},
	}

	bag := Baggage{list: baggage.List{
		"foo": {
			Value: "1",
			Properties: []baggage.Property{
				{Key: "state", Value: "on", HasValue: true},
				{Key: "red"},
			},
		},
		"bar": {
			Value:      "2",
			Properties: []baggage.Property{{Key: "yellow"}},
		},
	}}

	assert.ElementsMatch(t, members, bag.Members())
}

func TestBaggageMember(t *testing.T) {
	bag := Baggage{list: baggage.List{"foo": {Value: "1"}}}
	assert.Equal(t, Member{key: "foo", value: "1", hasData: true}, bag.Member("foo"))
	assert.Equal(t, Member{}, bag.Member("bar"))
}

func TestMemberKey(t *testing.T) {
	m := Member{}
	assert.Empty(t, m.Key(), "even invalid values should be returned")

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
	m := Member{hasData: false}
	assert.ErrorIs(t, m.validate(), errInvalidMember)

	m.hasData = true
	assert.ErrorIs(t, m.validate(), errInvalidKey)

	// Invalid UTF-8 in value
	m.key, m.value = "k", string([]byte{255})
	assert.ErrorIs(t, m.validate(), errInvalidValue)

	m.key, m.value = "k", "\\"
	assert.NoError(t, m.validate())
}

func TestNewMember(t *testing.T) {
	m, err := NewMember("", "")
	assert.ErrorIs(t, err, errInvalidKey)
	assert.Equal(t, Member{hasData: false}, m)

	key, val := "k", "v"
	p := Property{key: "foo"}
	m, err = NewMember(key, val, p)
	assert.NoError(t, err)
	expected := Member{
		key:        key,
		value:      val,
		properties: properties{{key: "foo"}},
		hasData:    true,
	}
	assert.Equal(t, expected, m)

	// it won't use percent decoding for key
	key = "%3B"
	m, err = NewMember(key, val, p)
	assert.NoError(t, err)
	expected = Member{
		key:        key,
		value:      val,
		properties: properties{{key: "foo"}},
		hasData:    true,
	}
	assert.Equal(t, expected, m)

	// wrong value with invalid token
	key = "k"
	val = ";"
	_, err = NewMember(key, val, p)
	assert.ErrorIs(t, err, errInvalidValue)

	// wrong value with wrong decoding
	val = "%zzzzz"
	_, err = NewMember(key, val, p)
	assert.ErrorIs(t, err, errInvalidValue)

	// value should be decoded
	val = "%3B"
	m, err = NewMember(key, val, p)
	expected = Member{
		key:        key,
		value:      ";",
		properties: properties{{key: "foo"}},
		hasData:    true,
	}
	assert.NoError(t, err)
	assert.Equal(t, expected, m)

	// Ensure new member is immutable.
	p.key = "bar"
	assert.Equal(t, expected, m)
}

func TestNewMemberRaw(t *testing.T) {
	m, err := NewMemberRaw("", "")
	assert.ErrorIs(t, err, errInvalidKey)
	assert.Equal(t, Member{hasData: false}, m)

	key, val := "k", "v"
	p := Property{key: "foo"}
	m, err = NewMemberRaw(key, val, p)
	assert.NoError(t, err)
	expected := Member{
		key:        key,
		value:      val,
		properties: properties{{key: "foo"}},
		hasData:    true,
	}
	assert.Equal(t, expected, m)

	// Ensure new member is immutable.
	p.key = "bar"
	assert.Equal(t, expected, m)
}

func TestBaggageUTF8(t *testing.T) {
	testCases := map[string]string{
		"ąść": "B% 💼",

		// Case sensitive
		"a": "a",
		"A": "A",
	}

	var members []Member
	for k, v := range testCases {
		m, err := NewMemberRaw(k, v)
		require.NoError(t, err)

		members = append(members, m)
	}

	b, err := New(members...)
	require.NoError(t, err)

	for k, v := range testCases {
		assert.Equal(t, v, b.Member(k).Value())
	}
}

func TestPropertiesValidate(t *testing.T) {
	p := properties{{}}
	assert.ErrorIs(t, p.validate(), errInvalidKey)

	p[0].key = "foo"
	assert.NoError(t, p.validate())

	p = append(p, Property{key: "bar"})
	assert.NoError(t, p.validate())
}

func TestMemberString(t *testing.T) {
	// normal key value pair
	member, _ := NewMemberRaw("key", "value")
	memberStr := member.String()
	assert.Equal(t, "key=value", memberStr)

	// encoded value
	member, _ = NewMemberRaw("key", "; ")
	memberStr = member.String()
	assert.Equal(t, "key=%3B%20", memberStr)
}

var benchBaggage Baggage

func BenchmarkNew(b *testing.B) {
	mem1, _ := NewMemberRaw("key1", "val1")
	mem2, _ := NewMemberRaw("key2", "val2")
	mem3, _ := NewMemberRaw("key3", "val3")
	mem4, _ := NewMemberRaw("key4", "val4")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchBaggage, _ = New(mem1, mem2, mem3, mem4)
	}
}

var benchMember Member

func BenchmarkNewMemberRaw(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		benchMember, _ = NewMemberRaw("key", "value")
	}
}

func BenchmarkParse(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		benchBaggage, _ = Parse(
			"userId=alice,serverNode = DF28 , isProduction = false,hasProp=stuff;propKey;propWValue=value, invalidUtf8=pr%ffo%ffp%fcValue",
		)
	}
}

func BenchmarkString(b *testing.B) {
	var members []Member
	addMember := func(k, v string) {
		m, err := NewMemberRaw(k, valueEscape(v))
		require.NoError(b, err)
		members = append(members, m)
	}

	addMember("key1", "val1")
	addMember("key2", " ;,%")
	addMember("B% 💼", "Witaj świecie!")
	addMember("key4", strings.Repeat("Hello world!", 10))

	bg, err := New(members...)
	require.NoError(b, err)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = bg.String()
	}
}

func BenchmarkValueEscape(b *testing.B) {
	testCases := []struct {
		name string
		in   string
	}{
		{name: "nothing to escape", in: "value"},
		{name: "requires escaping", in: " ;,%"},
		{name: "long value", in: strings.Repeat("Hello world!", 20)},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_ = valueEscape(tc.in)
			}
		})
	}
}

func BenchmarkMemberString(b *testing.B) {
	alphabet := "abcdefghijklmnopqrstuvwxyz"
	props := make([]Property, len(alphabet))
	for i, r := range alphabet {
		props[i] = Property{key: string(r)}
	}
	member, err := NewMember(alphabet, alphabet, props...)
	require.NoError(b, err)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = member.String()
	}
}
