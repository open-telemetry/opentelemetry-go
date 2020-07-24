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
	commonpb "go.opentelemetry.io/otel/internal/opentelemetry-proto-gen/common/v1"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/sdk/resource"
)

// Attributes transforms a slice of KeyValues into a slice of OTLP attribute key-values.
func Attributes(attrs []kv.KeyValue) []*commonpb.KeyValue {
	if len(attrs) == 0 {
		return nil
	}

	out := make([]*commonpb.KeyValue, 0, len(attrs))
	for _, kv := range attrs {
		out = append(out, toAttribute(kv))
	}
	return out
}

// ResourceAttributes transforms a Resource into a slice of OTLP attribute key-values.
func ResourceAttributes(resource *resource.Resource) []*commonpb.KeyValue {
	if resource.Len() == 0 {
		return nil
	}

	out := make([]*commonpb.KeyValue, 0, resource.Len())
	for iter := resource.Iter(); iter.Next(); {
		out = append(out, toAttribute(iter.Attribute()))
	}

	return out
}

func toAttribute(v kv.KeyValue) *commonpb.KeyValue {
	result := &commonpb.KeyValue{
		Key:   string(v.Key),
		Value: new(commonpb.AnyValue),
	}
	switch v.Value.Type() {
	case kv.BOOL:
		result.Value.Value = &commonpb.AnyValue_BoolValue{
			BoolValue: v.Value.AsBool(),
		}
	case kv.INT64, kv.INT32, kv.UINT32, kv.UINT64:
		result.Value.Value = &commonpb.AnyValue_IntValue{
			IntValue: v.Value.AsInt64(),
		}
	case kv.FLOAT32:
		result.Value.Value = &commonpb.AnyValue_DoubleValue{
			DoubleValue: float64(v.Value.AsFloat32()),
		}
	case kv.FLOAT64:
		result.Value.Value = &commonpb.AnyValue_DoubleValue{
			DoubleValue: v.Value.AsFloat64(),
		}
	case kv.STRING:
		result.Value.Value = &commonpb.AnyValue_StringValue{
			StringValue: v.Value.AsString(),
		}
	default:
		result.Value.Value = &commonpb.AnyValue_StringValue{
			StringValue: "INVALID",
		}
	}
	return result
}
