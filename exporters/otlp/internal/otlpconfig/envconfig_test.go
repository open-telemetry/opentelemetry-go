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

package otlpconfig

import (
	"reflect"
	"testing"
)

func TestStringToHeader(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  map[string]string
	}{
		{
			name:  "simple test",
			value: "userId=alice",
			want:  map[string]string{"userId": "alice"},
		},
		{
			name:  "simple test with spaces",
			value: " userId = alice  ",
			want:  map[string]string{"userId": "alice"},
		},
		{
			name:  "multiples headers encoded",
			value: "userId=alice,serverNode=DF%3A28,isProduction=false",
			want: map[string]string{
				"userId":       "alice",
				"serverNode":   "DF:28",
				"isProduction": "false",
			},
		},
		{
			name:  "invalid headers format",
			value: "userId:alice",
			want:  map[string]string{},
		},
		{
			name:  "invalid key",
			value: "%XX=missing,userId=alice",
			want: map[string]string{
				"userId": "alice",
			},
		},
		{
			name:  "invalid value",
			value: "missing=%XX,userId=alice",
			want: map[string]string{
				"userId": "alice",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringToHeader(tt.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stringToHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
