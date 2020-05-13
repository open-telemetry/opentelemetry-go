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
	"testing"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/label"

	"github.com/stretchr/testify/require"
)

type testCase struct {
	kvs      []kv.KeyValue
	encoding string
}

func expect(enc string, kvs ...kv.KeyValue) testCase {
	return testCase{
		kvs:      kvs,
		encoding: enc,
	}
}

func TestSetDedup(t *testing.T) {
	cases := []testCase{
		expect("A=B", kv.String("A", "2"), kv.String("A", "B")),
		expect("A=B", kv.String("A", "2"), kv.Int("A", 1), kv.String("A", "B")),
		expect("A=B", kv.String("A", "B"), kv.String("A", "C"), kv.String("A", "D"), kv.String("A", "B")),

		expect("A=B,C=D", kv.String("A", "1"), kv.String("C", "D"), kv.String("A", "B")),
		expect("A=B,C=D", kv.String("A", "2"), kv.String("A", "B"), kv.String("C", "D")),
		expect("A=B,C=D", kv.Float64("C", 1.2), kv.String("A", "2"), kv.String("A", "B"), kv.String("C", "D")),
		expect("A=B,C=D", kv.String("C", "D"), kv.String("A", "B"), kv.String("A", "C"), kv.String("A", "D"), kv.String("A", "B")),
		expect("A=B,C=D", kv.String("A", "B"), kv.String("C", "D"), kv.String("A", "C"), kv.String("A", "D"), kv.String("A", "B")),
		expect("A=B,C=D", kv.String("A", "B"), kv.String("A", "C"), kv.String("A", "D"), kv.String("A", "B"), kv.String("C", "D")),
	}
	enc := label.DefaultEncoder()

	s2d := map[string][]label.Distinct{}
	d2s := map[label.Distinct][]string{}

	for _, tc := range cases {
		cpy := make([]kv.KeyValue, len(tc.kvs))
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
