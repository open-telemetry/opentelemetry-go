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

package view

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/internal/instrument"
)

var matchInstrument = instrument.Description{
	Library: instrumentation.Library{
		Name:      "bar",
		Version:   "v1.0.0",
		SchemaURL: "stuff.test/",
	},
	InstrumentName:        "foo",
	InstrumentDescription: "",
}
var noMatchInstrument = instrument.Description{
	Library: instrumentation.Library{
		Name:      "notfoo",
		Version:   "v0.x.0",
		SchemaURL: "notstuff.test/",
	},
	InstrumentName:        "notstuff",
	InstrumentDescription: "",
}

var emptyDescription = instrument.Description{}

func TestConfig_Transform(t *testing.T) {
	tests := []struct {
		name     string
		options  []Option
		match    instrument.Description
		notMatch instrument.Description
	}{
		{
			name: "instrument name",
			options: []Option{
				MatchInstrumentName("foo"),
			},
			match:    matchInstrument,
			notMatch: emptyDescription,
		},
		{
			name: "Library name",
			options: []Option{
				MatchInstrumentationLibrary(instrumentation.Library{
					Name: "bar",
				}),
			},
			match:    matchInstrument,
			notMatch: emptyDescription,
		},
		{
			name: "Library version",
			options: []Option{
				MatchInstrumentationLibrary(instrumentation.Library{
					Version: "v1.0.0",
				}),
			},

			match:    matchInstrument,
			notMatch: emptyDescription,
		},
		{
			name: "Library SchemaURL",
			options: []Option{
				MatchInstrumentationLibrary(instrumentation.Library{
					SchemaURL: "stuff.test/",
				}),
			},
			match:    matchInstrument,
			notMatch: emptyDescription,
		},
		{
			name: "Regex",
			options: []Option{
				MatchInstrumentNameRegexp(regexp.MustCompile("^f.*")),
			},
			match:    matchInstrument,
			notMatch: emptyDescription,
		},
		{
			name: "composite literal name",
			options: []Option{
				MatchInstrumentName("foo"),
				MatchInstrumentationLibrary(instrumentation.Library{
					Name:      "bar",
					Version:   "v1.0.0",
					SchemaURL: "stuff.test/",
				}),
			},
			match:    matchInstrument,
			notMatch: emptyDescription,
		},
		{
			name: "composite regex name",
			options: []Option{
				MatchInstrumentNameRegexp(regexp.MustCompile("^f.*")),
				MatchInstrumentationLibrary(instrumentation.Library{
					Name:      "bar",
					Version:   "v1.0.0",
					SchemaURL: "stuff.test/",
				}),
			},
			match:    matchInstrument,
			notMatch: emptyDescription,
		},

		{
			name: "rename",
			options: []Option{
				MatchInstrumentName("foo"),
				WithName("baz"),
			},
			match: instrument.Description{
				Library: instrumentation.Library{
					Name:      "bar",
					Version:   "v1.0.0",
					SchemaURL: "stuff.test/",
				},
				InstrumentName:        "baz",
				InstrumentDescription: "",
			},
			notMatch: emptyDescription,
		},
		{
			name: "rename",
			options: []Option{
				MatchInstrumentName("foo"),
				WithDescription("descriptive stuff"),
			},
			match: instrument.Description{
				Library: instrumentation.Library{
					Name:      "bar",
					Version:   "v1.0.0",
					SchemaURL: "stuff.test/",
				},
				InstrumentName:        "foo",
				InstrumentDescription: "descriptive stuff",
			},
			notMatch: emptyDescription,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := New(tt.options...)
			require.NoError(t, err)

			t.Run("match", func(t *testing.T) {
				got, match := cfg.TransformInstrument(matchInstrument)
				assert.Equal(t, tt.match, got)
				assert.True(t, match)
			})

			t.Run("does not match", func(t *testing.T) {
				got, match := cfg.TransformInstrument(noMatchInstrument)
				assert.Equal(t, tt.notMatch, got)
				assert.False(t, match)
			})
		})
	}
}

func TestConfig_TransformAttributes(t *testing.T) {
	inputSet := attribute.NewSet(
		attribute.String("foo", "bar"),
		attribute.Int("power-level", 9001),
		attribute.Float64("lifeUniverseEverything", 42.0),
	)

	tests := []struct {
		name   string
		filter []attribute.Key
		want   attribute.Set
	}{
		{
			name:   "empty should match all",
			filter: []attribute.Key{},
			want:   inputSet,
		},
		{
			name: "Match 1",
			filter: []attribute.Key{
				attribute.Key("power-level"),
			},
			want: attribute.NewSet(
				attribute.Int("power-level", 9001),
			),
		},
		{
			name: "Match 2",
			filter: []attribute.Key{
				attribute.Key("foo"),
				attribute.Key("lifeUniverseEverything"),
			},
			want: attribute.NewSet(
				attribute.Float64("lifeUniverseEverything", 42.0),
				attribute.String("foo", "bar"),
			),
		},
		{
			name: "Don't match",
			filter: []attribute.Key{
				attribute.Key("nothing"),
			},
			want: attribute.NewSet(),
		},
		{
			name: "Match some",
			filter: []attribute.Key{
				attribute.Key("power-level"),
				attribute.Key("nothing"),
			},
			want: attribute.NewSet(
				attribute.Int("power-level", 9001),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := New(WithAttributeFilter(tt.filter...))
			require.NoError(t, err)

			got := cfg.TransformAttributes(inputSet)
			assert.Equal(t, got.Equivalent(), tt.want.Equivalent())
		})
	}
}
