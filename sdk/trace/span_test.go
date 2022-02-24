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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
)

func TestTruncateAttr(t *testing.T) {
	const key = "key"

	tests := []struct {
		limit      int
		attr, want attribute.KeyValue
	}{
		{
			limit: -1,
			attr:  attribute.String(key, "value"),
			want:  attribute.String(key, "value"),
		},
		{
			limit: -1,
			attr:  attribute.StringSlice(key, []string{"value-0", "value-1"}),
			want:  attribute.StringSlice(key, []string{"value-0", "value-1"}),
		},
		{
			limit: 0,
			attr:  attribute.String(key, "value"),
			want:  attribute.String(key, ""),
		},
		{
			limit: 0,
			attr:  attribute.StringSlice(key, []string{"value-0", "value-1"}),
			want:  attribute.StringSlice(key, []string{"", ""}),
		},
		{
			limit: 1,
			attr:  attribute.String(key, "value"),
			want:  attribute.String(key, "v"),
		},
		{
			limit: 1,
			attr:  attribute.StringSlice(key, []string{"value-0", "value-1"}),
			want:  attribute.StringSlice(key, []string{"v", "v"}),
		},
		{
			limit: 5,
			attr:  attribute.String(key, "value"),
			want:  attribute.String(key, "value"),
		},
		{
			limit: 7,
			attr:  attribute.StringSlice(key, []string{"value-0", "value-1"}),
			want:  attribute.StringSlice(key, []string{"value-0", "value-1"}),
		},
		{
			limit: 6,
			attr:  attribute.StringSlice(key, []string{"value", "value-1"}),
			want:  attribute.StringSlice(key, []string{"value", "value-"}),
		},
		{
			limit: 128,
			attr:  attribute.String(key, "value"),
			want:  attribute.String(key, "value"),
		},
		{
			limit: 128,
			attr:  attribute.StringSlice(key, []string{"value-0", "value-1"}),
			want:  attribute.StringSlice(key, []string{"value-0", "value-1"}),
		},
	}

	for _, test := range tests {
		name := fmt.Sprintf("%s->%s(limit:%d)", test.attr.Key, test.attr.Value.Emit(), test.limit)
		t.Run(name, func(t *testing.T) {
			truncateAttr(test.limit, &test.attr)
			assert.Equal(t, test.want, test.attr)
		})
	}
}
