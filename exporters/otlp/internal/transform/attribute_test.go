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

package transform

import (
	"testing"

	commonpb "github.com/open-telemetry/opentelemetry-proto/gen/go/common/v1"
	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"
)

func TestAttributes(t *testing.T) {
	for _, test := range []struct {
		attrs    []core.KeyValue
		expected []*commonpb.AttributeKeyValue
	}{
		{nil, nil},
		{
			[]core.KeyValue{
				key.Int("int to int", 123),
				key.Uint("uint to int", 1234),
				key.Int32("int32 to int", 12345),
				key.Uint32("uint32 to int", 123456),
				key.Int64("int64 to int64", 1234567),
				key.Uint64("uint64 to int64", 12345678),
				key.Float32("float32 to double", 3.14),
				key.Float32("float64 to double", 1.61),
				key.String("string to string", "string"),
				key.Bool("bool to bool", true),
			},
			[]*commonpb.AttributeKeyValue{
				{
					Key:      "int to int",
					Type:     commonpb.AttributeKeyValue_INT,
					IntValue: 123,
				},
				{
					Key:      "uint to int",
					Type:     commonpb.AttributeKeyValue_INT,
					IntValue: 1234,
				},
				{
					Key:      "int32 to int",
					Type:     commonpb.AttributeKeyValue_INT,
					IntValue: 12345,
				},
				{
					Key:      "uint32 to int",
					Type:     commonpb.AttributeKeyValue_INT,
					IntValue: 123456,
				},
				{
					Key:      "int64 to int64",
					Type:     commonpb.AttributeKeyValue_INT,
					IntValue: 1234567,
				},
				{
					Key:      "uint64 to int64",
					Type:     commonpb.AttributeKeyValue_INT,
					IntValue: 12345678,
				},
				{
					Key:         "float32 to double",
					Type:        commonpb.AttributeKeyValue_DOUBLE,
					DoubleValue: 3.14,
				},
				{
					Key:         "float64 to double",
					Type:        commonpb.AttributeKeyValue_DOUBLE,
					DoubleValue: 1.61,
				},
				{
					Key:         "string to string",
					Type:        commonpb.AttributeKeyValue_STRING,
					StringValue: "string",
				},
				{
					Key:       "bool to bool",
					Type:      commonpb.AttributeKeyValue_BOOL,
					BoolValue: true,
				},
			},
		},
	} {
		got := Attributes(test.attrs)
		if !assert.Len(t, got, len(test.expected)) {
			continue
		}
		for i, actual := range got {
			if actual.Type == commonpb.AttributeKeyValue_DOUBLE {
				if !assert.InDelta(t, test.expected[i].DoubleValue, actual.DoubleValue, 0.01) {
					continue
				}
				test.expected[i].DoubleValue = actual.DoubleValue
			}
			assert.Equal(t, test.expected[i], actual)
		}
	}
}
