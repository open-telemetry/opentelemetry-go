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
	commonpb "github.com/open-telemetry/opentelemetry-proto/gen/go/common/v1"

	"go.opentelemetry.io/otel/api/core"
)

// Attributes transforms a slice of KeyValues into a slice of OTLP attribute key-values.
func Attributes(attrs []core.KeyValue) []*commonpb.AttributeKeyValue {
	if len(attrs) == 0 {
		return nil
	}

	out := make([]*commonpb.AttributeKeyValue, 0, len(attrs))
	for _, v := range attrs {
		switch v.Value.Type() {
		case core.BOOL:
			out = append(out, &commonpb.AttributeKeyValue{
				Key:       string(v.Key),
				Type:      commonpb.AttributeKeyValue_BOOL,
				BoolValue: v.Value.AsBool(),
			})
		case core.INT64, core.INT32, core.UINT32, core.UINT64:
			out = append(out, &commonpb.AttributeKeyValue{
				Key:      string(v.Key),
				Type:     commonpb.AttributeKeyValue_INT,
				IntValue: v.Value.AsInt64(),
			})
		case core.FLOAT32:
			f32 := v.Value.AsFloat32()
			out = append(out, &commonpb.AttributeKeyValue{
				Key:         string(v.Key),
				Type:        commonpb.AttributeKeyValue_DOUBLE,
				DoubleValue: float64(f32),
			})
		case core.FLOAT64:
			out = append(out, &commonpb.AttributeKeyValue{
				Key:         string(v.Key),
				Type:        commonpb.AttributeKeyValue_DOUBLE,
				DoubleValue: v.Value.AsFloat64(),
			})
		case core.STRING:
			out = append(out, &commonpb.AttributeKeyValue{
				Key:         string(v.Key),
				Type:        commonpb.AttributeKeyValue_STRING,
				StringValue: v.Value.AsString(),
			})
		}
	}
	return out
}
