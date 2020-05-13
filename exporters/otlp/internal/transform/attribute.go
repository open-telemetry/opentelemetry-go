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

	"go.opentelemetry.io/otel/api/kv/value"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/sdk/resource"
)

// Attributes transforms a slice of KeyValues into a slice of OTLP attribute key-values.
func Attributes(attrs []kv.KeyValue) []*commonpb.AttributeKeyValue {
	if len(attrs) == 0 {
		return nil
	}

	out := make([]*commonpb.AttributeKeyValue, 0, len(attrs))
	for _, kv := range attrs {
		out = append(out, toAttribute(kv))
	}
	return out
}

// ResourceAttributes transforms a Resource into a slice of OTLP attribute key-values.
func ResourceAttributes(resource *resource.Resource) []*commonpb.AttributeKeyValue {
	if resource.Len() == 0 {
		return nil
	}

	out := make([]*commonpb.AttributeKeyValue, 0, resource.Len())
	for iter := resource.Iter(); iter.Next(); {
		out = append(out, toAttribute(iter.Attribute()))
	}

	return out
}

func toAttribute(v kv.KeyValue) *commonpb.AttributeKeyValue {
	switch v.Value.Type() {
	case value.BOOL:
		return &commonpb.AttributeKeyValue{
			Key:       string(v.Key),
			Type:      commonpb.AttributeKeyValue_BOOL,
			BoolValue: v.Value.AsBool(),
		}
	case value.INT64, value.INT32, value.UINT32, value.UINT64:
		return &commonpb.AttributeKeyValue{
			Key:      string(v.Key),
			Type:     commonpb.AttributeKeyValue_INT,
			IntValue: v.Value.AsInt64(),
		}
	case value.FLOAT32:
		return &commonpb.AttributeKeyValue{
			Key:         string(v.Key),
			Type:        commonpb.AttributeKeyValue_DOUBLE,
			DoubleValue: float64(v.Value.AsFloat32()),
		}
	case value.FLOAT64:
		return &commonpb.AttributeKeyValue{
			Key:         string(v.Key),
			Type:        commonpb.AttributeKeyValue_DOUBLE,
			DoubleValue: v.Value.AsFloat64(),
		}
	case value.STRING:
		return &commonpb.AttributeKeyValue{
			Key:         string(v.Key),
			Type:        commonpb.AttributeKeyValue_STRING,
			StringValue: v.Value.AsString(),
		}
	default:
		return &commonpb.AttributeKeyValue{
			Key:         string(v.Key),
			Type:        commonpb.AttributeKeyValue_STRING,
			StringValue: "INVALID",
		}
	}
}
