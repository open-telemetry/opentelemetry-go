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

package label_test

import (
	"regexp"
	"testing"

	"go.opentelemetry.io/otel/label"

	"github.com/stretchr/testify/require"
)

type testCase struct {
	kvs []label.KeyValue

	keyRe *regexp.Regexp

	encoding string
	fullEnc  string
}

func expect(enc string, kvs ...label.KeyValue) testCase {
	return testCase{
		kvs:      kvs,
		encoding: enc,
	}
}

func expectFiltered(enc, filter, fullEnc string, kvs ...label.KeyValue) testCase {
	return testCase{
		kvs:      kvs,
		keyRe:    regexp.MustCompile(filter),
		encoding: enc,
		fullEnc:  fullEnc,
	}
}

func TestSetDedup(t *testing.T) {
	cases := []testCase{
		expect("A=B", label.String("A", "2"), label.String("A", "B")),
		expect("A=B", label.String("A", "2"), label.Int("A", 1), label.String("A", "B")),
		expect("A=B", label.String("A", "B"), label.String("A", "C"), label.String("A", "D"), label.String("A", "B")),

		expect("A=B,C=D", label.String("A", "1"), label.String("C", "D"), label.String("A", "B")),
		expect("A=B,C=D", label.String("A", "2"), label.String("A", "B"), label.String("C", "D")),
		expect("A=B,C=D", label.Float64("C", 1.2), label.String("A", "2"), label.String("A", "B"), label.String("C", "D")),
		expect("A=B,C=D", label.String("C", "D"), label.String("A", "B"), label.String("A", "C"), label.String("A", "D"), label.String("A", "B")),
		expect("A=B,C=D", label.String("A", "B"), label.String("C", "D"), label.String("A", "C"), label.String("A", "D"), label.String("A", "B")),
		expect("A=B,C=D", label.String("A", "B"), label.String("A", "C"), label.String("A", "D"), label.String("A", "B"), label.String("C", "D")),
	}
	enc := label.DefaultEncoder()

	s2d := map[string][]label.Distinct{}
	d2s := map[label.Distinct][]string{}

	for _, tc := range cases {
		cpy := make([]label.KeyValue, len(tc.kvs))
		copy(cpy, tc.kvs)
		sl := label.NewSet(cpy...)

		// Ensure that the input was reordered but no elements went missing.
		require.ElementsMatch(t, tc.kvs, cpy)

		str := sl.Encoded(enc)
		equ := sl.Equivalent()

		s2d[str] = append(s2d[str], equ)
		d2s[equ] = append(d2s[equ], str)

		require.Equal(t, tc.encoding, str)
	}

	for s, d := range s2d {
		// No other Distinct values are equal to this.
		for s2, d2 := range s2d {
			if s2 == s {
				continue
			}
			for _, elt := range d {
				for _, otherDistinct := range d2 {
					require.NotEqual(t, otherDistinct, elt)
				}
			}
		}
		for _, strings := range d2s {
			if strings[0] == s {
				continue
			}
			for _, otherString := range strings {
				require.NotEqual(t, otherString, s)
			}
		}
	}

	for d, s := range d2s {
		// No other Distinct values are equal to this.
		for d2, s2 := range d2s {
			if d2 == d {
				continue
			}
			for _, elt := range s {
				for _, otherDistinct := range s2 {
					require.NotEqual(t, otherDistinct, elt)
				}
			}
		}
		for _, distincts := range s2d {
			if distincts[0] == d {
				continue
			}
			for _, otherDistinct := range distincts {
				require.NotEqual(t, otherDistinct, d)
			}
		}
	}
}

func TestUniqueness(t *testing.T) {
	short := []label.KeyValue{
		label.String("A", "0"),
		label.String("B", "2"),
		label.String("A", "1"),
	}
	long := []label.KeyValue{
		label.String("B", "2"),
		label.String("C", "5"),
		label.String("B", "2"),
		label.String("C", "1"),
		label.String("A", "4"),
		label.String("C", "3"),
		label.String("A", "1"),
	}
	cases := []testCase{
		expectFiltered("A=1", "^A$", "B=2", short...),
		expectFiltered("B=2", "^B$", "A=1", short...),
		expectFiltered("A=1,B=2", "^A|B$", "", short...),
		expectFiltered("", "^C", "A=1,B=2", short...),

		expectFiltered("A=1,C=3", "A|C", "B=2", long...),
		expectFiltered("B=2,C=3", "C|B", "A=1", long...),
		expectFiltered("C=3", "C", "A=1,B=2", long...),
		expectFiltered("", "D", "A=1,B=2,C=3", long...),
	}
	enc := label.DefaultEncoder()

	for _, tc := range cases {
		cpy := make([]label.KeyValue, len(tc.kvs))
		copy(cpy, tc.kvs)
		distinct, uniq := label.NewSetWithFiltered(cpy, func(label label.KeyValue) bool {
			return tc.keyRe.MatchString(string(label.Key))
		})

		full := label.NewSet(uniq...)

		require.Equal(t, tc.encoding, distinct.Encoded(enc))
		require.Equal(t, tc.fullEnc, full.Encoded(enc))
	}
}

func TestLookup(t *testing.T) {
	set := label.NewSet(label.Int("C", 3), label.Int("A", 1), label.Int("B", 2))

	value, has := set.Value("C")
	require.True(t, has)
	require.Equal(t, int64(3), value.AsInt64())

	value, has = set.Value("B")
	require.True(t, has)
	require.Equal(t, int64(2), value.AsInt64())

	value, has = set.Value("A")
	require.True(t, has)
	require.Equal(t, int64(1), value.AsInt64())

	value, has = set.Value("D")
	require.False(t, has)
}
