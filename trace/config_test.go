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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
)

func TestNewSpanConfig(t *testing.T) {
	k1v1 := attribute.String("key1", "value1")
	k1v2 := attribute.String("key1", "value2")
	k2v2 := attribute.String("key2", "value2")

	timestamp0 := time.Unix(0, 0)
	timestamp1 := time.Unix(0, 0)

	link1 := Link{
		SpanContext: SpanContext{traceID: TraceID([16]byte{1, 1}), spanID: SpanID{3}},
		Attributes:  []attribute.KeyValue{k1v1},
	}
	link2 := Link{
		SpanContext: SpanContext{traceID: TraceID([16]byte{1, 1}), spanID: SpanID{3}},
		Attributes:  []attribute.KeyValue{k1v2, k2v2},
	}

	tests := []struct {
		options  []SpanOption
		expected *SpanConfig
	}{
		{
			// No non-zero-values should be set.
			[]SpanOption{},
			new(SpanConfig),
		},
		{
			[]SpanOption{
				WithAttributes(k1v1),
			},
			&SpanConfig{
				Attributes: []attribute.KeyValue{k1v1},
			},
		},
		{
			// Multiple calls should append not overwrite.
			[]SpanOption{
				WithAttributes(k1v1),
				WithAttributes(k1v2),
				WithAttributes(k2v2),
			},
			&SpanConfig{
				// No uniqueness is guaranteed by the API.
				Attributes: []attribute.KeyValue{k1v1, k1v2, k2v2},
			},
		},
		{
			[]SpanOption{
				WithAttributes(k1v1, k1v2, k2v2),
			},
			&SpanConfig{
				// No uniqueness is guaranteed by the API.
				Attributes: []attribute.KeyValue{k1v1, k1v2, k2v2},
			},
		},
		{
			[]SpanOption{
				WithTimestamp(timestamp0),
			},
			&SpanConfig{
				Timestamp: timestamp0,
			},
		},
		{
			[]SpanOption{
				// Multiple calls overwrites with last-one-wins.
				WithTimestamp(timestamp0),
				WithTimestamp(timestamp1),
			},
			&SpanConfig{
				Timestamp: timestamp1,
			},
		},
		{
			[]SpanOption{
				WithLinks(link1),
			},
			&SpanConfig{
				Links: []Link{link1},
			},
		},
		{
			[]SpanOption{
				// Multiple calls should append not overwrite.
				WithLinks(link1),
				WithLinks(link1, link2),
			},
			&SpanConfig{
				// No uniqueness is guaranteed by the API.
				Links: []Link{link1, link1, link2},
			},
		},
		{
			[]SpanOption{
				WithNewRoot(),
			},
			&SpanConfig{
				NewRoot: true,
			},
		},
		{
			[]SpanOption{
				// Multiple calls should not change NewRoot state.
				WithNewRoot(),
				WithNewRoot(),
			},
			&SpanConfig{
				NewRoot: true,
			},
		},
		{
			[]SpanOption{
				WithSpanKind(SpanKindConsumer),
			},
			&SpanConfig{
				SpanKind: SpanKindConsumer,
			},
		},
		{
			[]SpanOption{
				// Multiple calls overwrites with last-one-wins.
				WithSpanKind(SpanKindClient),
				WithSpanKind(SpanKindConsumer),
			},
			&SpanConfig{
				SpanKind: SpanKindConsumer,
			},
		},
		{
			// Everything should work together.
			[]SpanOption{
				WithAttributes(k1v1),
				WithTimestamp(timestamp0),
				WithLinks(link1, link2),
				WithNewRoot(),
				WithSpanKind(SpanKindConsumer),
			},
			&SpanConfig{
				Attributes: []attribute.KeyValue{k1v1},
				Timestamp:  timestamp0,
				Links:      []Link{link1, link2},
				NewRoot:    true,
				SpanKind:   SpanKindConsumer,
			},
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, NewSpanConfig(test.options...))
	}
}

func TestTracerConfig(t *testing.T) {
	v1 := "semver:0.0.1"
	v2 := "semver:1.0.0"
	tests := []struct {
		options  []TracerOption
		expected *TracerConfig
	}{
		{
			// No non-zero-values should be set.
			[]TracerOption{},
			new(TracerConfig),
		},
		{
			[]TracerOption{
				WithInstrumentationVersion(v1),
			},
			&TracerConfig{
				InstrumentationVersion: v1,
			},
		},
		{
			[]TracerOption{
				// Multiple calls should overwrite.
				WithInstrumentationVersion(v1),
				WithInstrumentationVersion(v2),
			},
			&TracerConfig{
				InstrumentationVersion: v2,
			},
		},
	}
	for _, test := range tests {
		config := NewTracerConfig(test.options...)
		assert.Equal(t, test.expected, config)
	}
}
