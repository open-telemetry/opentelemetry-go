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

	"go.opentelemetry.io/otel/label"
	export "go.opentelemetry.io/otel/sdk/export/trace"
)

func TestAttributesMap(t *testing.T) {
	attrMap := newAttributesMap(128)

	for i := 0; i < 256; i++ {
		attrMap.add(label.Int(fmt.Sprintf("test-key-%d", i), i))
	}
	if attrMap.capacity != 128 {
		t.Fatalf("bad capacity: %v", attrMap.capacity)
	}

	if attrMap.droppedCount != 128 {
		t.Fatalf("bad dropped count: %v", attrMap.droppedCount)
	}

	for i := 0; i < 128; i++ {
		_, ok := attrMap.attributes[label.Key(fmt.Sprintf("test-key-%d", i))]
		if ok {
			t.Fatal("should be dropped")
		}
	}
	for i := 128; i < 256; i++ {
		_, ok := attrMap.attributes[label.Key(fmt.Sprintf("test-key-%d", i))]
		if !ok {
			t.Fatal("should not be dropped")
		}
	}
}

func TestAttributesMapGetOldestRemoveOldest(t *testing.T) {
	attrMap := newAttributesMap(128)

	for i := 0; i < 128; i++ {
		attrMap.add(label.Int(fmt.Sprintf("test-key-%d", i), i))
	}

	attrMap.removeOldest()
	attrMap.removeOldest()
	attrMap.removeOldest()

	for i := 0; i < 3; i++ {
		_, ok := attrMap.attributes[label.Key(fmt.Sprintf("test-key-%d", i))]
		if ok {
			t.Fatal("should be removed")
		}
	}
}

func TestAttributesMapToSpanData(t *testing.T) {
	attrMap := newAttributesMap(128)

	for i := 0; i < 128; i++ {
		attrMap.add(label.Int(fmt.Sprintf("test-key-%d", i), i))
	}

	sd := &export.SpanData{}

	attrMap.toSpanData(sd)

	if attrMap.droppedCount != sd.DroppedAttributeCount {
		t.Fatalf("droppedCount in map and span are not equal: %d != %d",
			attrMap.droppedCount,
			sd.DroppedAttributeCount)
	}

	if len(attrMap.attributes) != len(sd.Attributes) {
		t.Fatalf("amount of elements in map and span is not equal: %d != %d",
			len(attrMap.attributes),
			len(sd.Attributes))
	}
}
