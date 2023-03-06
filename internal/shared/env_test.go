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

package shared

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	dft := "default value"
	testCases := []struct {
		name string
		act  func() string
		envs map[string]string
		want string
	}{
		{
			name: "simple",
			act:  func() string { return envString(dft, "ENV_ONE") },
			envs: map[string]string{
				"ENV_ONE": "one",
			},
			want: "one",
		},
		{
			name: "first has precedence",
			act:  func() string { return envString(dft, "ENV_ONE", "ENV_TWO") },
			envs: map[string]string{
				"ENV_ONE": "one",
				"ENV_TWO": "two",
			},
			want: "one",
		},
		{
			name: "returns first not empty",
			act:  func() string { return envString(dft, "ENV_ONE", "ENV_TWO", "ENV_THREE") },
			envs: map[string]string{
				"ENV_TWO":   "",
				"ENV_THREE": "three",
			},
			want: "three",
		},
		{
			name: "returns default if all is empty",
			act:  func() string { return envString(dft, "ENV_ONE", "ENV_TWO", "ENV_THREE") },
			envs: map[string]string{
				"ENV_TWO": "",
			},
			want: dft,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.envs {
				t.Setenv(k, v)
			}

			got := tc.act()

			assert.Equal(t, tc.want, got)
		})
	}
}

func TestInt(t *testing.T) {
	dft := 999
	testCases := []struct {
		name string
		act  func() int
		envs map[string]string
		want int
	}{
		{
			name: "simple",
			act:  func() int { return envInt(dft, "ENV_ONE") },
			envs: map[string]string{
				"ENV_ONE": "1",
			},
			want: 1,
		},
		{
			name: "first has precedence",
			act:  func() int { return envInt(dft, "ENV_ONE", "ENV_TWO") },
			envs: map[string]string{
				"ENV_ONE": "1",
				"ENV_TWO": "2",
			},
			want: 1,
		},
		{
			name: "returns first not empty",
			act:  func() int { return envInt(dft, "ENV_ONE", "ENV_TWO", "ENV_THREE") },
			envs: map[string]string{
				"ENV_TWO":   "",
				"ENV_THREE": "3",
			},
			want: 3,
		},
		{
			name: "returns default if all is empty",
			act:  func() int { return envInt(dft, "ENV_ONE", "ENV_TWO", "ENV_THREE") },
			envs: map[string]string{
				"ENV_TWO": "",
			},
			want: dft,
		},
		{
			name: "returns the first valid value",
			act:  func() int { return envInt(dft, "ENV_ONE", "ENV_TWO", "ENV_THREE") },
			envs: map[string]string{
				"ENV_ONE": "bad",
				"ENV_TWO": "-2",
			},
			want: -2,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.envs {
				t.Setenv(k, v)
			}

			got := tc.act()

			assert.Equal(t, tc.want, got)
		})
	}
}

func TestDuration(t *testing.T) {
	dft := time.Minute
	testCases := []struct {
		name string
		act  func() time.Duration
		envs map[string]string
		want time.Duration
	}{
		{
			name: "simple",
			act:  func() time.Duration { return envDuration(dft, "ENV_ONE") },
			envs: map[string]string{
				"ENV_ONE": "1",
			},
			want: 1 * time.Millisecond,
		},
		{
			name: "first has precedence",
			act:  func() time.Duration { return envDuration(dft, "ENV_ONE", "ENV_TWO") },
			envs: map[string]string{
				"ENV_ONE": "1",
				"ENV_TWO": "2",
			},
			want: 1 * time.Millisecond,
		},
		{
			name: "returns first not empty",
			act:  func() time.Duration { return envDuration(dft, "ENV_ONE", "ENV_TWO", "ENV_THREE") },
			envs: map[string]string{
				"ENV_TWO":   "",
				"ENV_THREE": "3",
			},
			want: 3 * time.Millisecond,
		},
		{
			name: "returns default if all is empty",
			act:  func() time.Duration { return envDuration(dft, "ENV_ONE", "ENV_TWO", "ENV_THREE") },
			envs: map[string]string{
				"ENV_TWO": "",
			},
			want: dft,
		},
		{
			name: "returns the first valid value",
			act:  func() time.Duration { return envDuration(dft, "ENV_ONE", "ENV_TWO", "ENV_THREE") },
			envs: map[string]string{
				"ENV_ONE":   "bad",
				"ENV_TWO":   "-2",
				"ENV_THREE": "3",
			},
			want: 3 * time.Millisecond,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.envs {
				t.Setenv(k, v)
			}

			got := tc.act()

			assert.Equal(t, tc.want, got)
		})
	}
}
