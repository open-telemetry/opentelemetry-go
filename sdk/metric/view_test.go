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

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"fmt"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
)

var (
	schemaURL  = "https://opentelemetry.io/schemas/1.0.0"
	completeIP = Instrument{
		Name:        "foo",
		Description: "foo desc",
		Kind:        InstrumentKindSyncCounter,
		Unit:        unit.Bytes,
		Scope: instrumentation.Scope{
			Name:      "TestNewViewMatch",
			Version:   "v0.1.0",
			SchemaURL: schemaURL,
		},
	}
)

func TestNewViewMatch(t *testing.T) {
	tests := []struct {
		name       string
		criteria   Instrument
		matches    []Instrument
		notMatches []Instrument
	}{
		{
			name:       "Empty",
			notMatches: []Instrument{{}, {Name: "foo"}, completeIP},
		},
		{
			name:     "Name/Exact",
			criteria: Instrument{Name: "foo"},
			matches:  []Instrument{{Name: "foo"}, completeIP},
			notMatches: []Instrument{
				{},
				{Name: "bar"},
				{Name: "foobar"},
				{Name: "barfoo"},
				{Name: "ffooo"},
			},
		},
		{
			name:     "Name/Wildcard/*",
			criteria: Instrument{Name: "*"},
			matches: []Instrument{
				{},
				{Name: "foo"},
				{Name: "foobar"},
				{Name: "barfoo"},
				{Name: "barfoobaz"},
				completeIP,
			},
		},
		{
			name:     "Name/Wildcard/Front?",
			criteria: Instrument{Name: "?oo"},
			matches: []Instrument{
				{Name: "foo"},
				{Name: "1oo"},
				completeIP,
			},
			notMatches: []Instrument{
				{},
				{Name: "bar"},
				{Name: "foobar"},
				{Name: "barfoo"},
				{Name: "barfoobaz"},
			},
		},
		{
			name:     "Name/Wildcard/Back?",
			criteria: Instrument{Name: "fo?"},
			matches: []Instrument{
				{Name: "foo"},
				{Name: "fo1"},
				completeIP,
			},
			notMatches: []Instrument{
				{},
				{Name: "bar"},
				{Name: "foobar"},
				{Name: "barfoo"},
				{Name: "barfoobaz"},
			},
		},
		{
			name:     "Name/Wildcard/Front*",
			criteria: Instrument{Name: "*foo"},
			matches: []Instrument{
				{Name: "foo"},
				{Name: "123foo"},
				{Name: "barfoo"},
				completeIP,
			},
			notMatches: []Instrument{
				{},
				{Name: "bar"},
				{Name: "foobar"},
				{Name: "barfoobaz"},
			},
		},
		{
			name:     "Name/Wildcard/Back*",
			criteria: Instrument{Name: "foo*"},
			matches: []Instrument{
				{Name: "foo"},
				{Name: "foo1"},
				{Name: "foobar"},
				completeIP,
			},
			notMatches: []Instrument{
				{},
				{Name: "bar"},
				{Name: "barfoo"},
				{Name: "barfoobaz"},
			},
		},
		{
			name:     "Name/Wildcard/FrontBack*",
			criteria: Instrument{Name: "*foo*"},
			matches: []Instrument{
				{Name: "foo"},
				{Name: "foo1"},
				{Name: "1foo"},
				{Name: "1foo1"},
				{Name: "foobar"},
				{Name: "barfoobaz"},
				completeIP,
			},
			notMatches: []Instrument{
				{},
				{Name: "bar"},
			},
		},
		{
			name:     "Name/Wildcard/Front**",
			criteria: Instrument{Name: "**foo"},
			matches: []Instrument{
				{Name: "foo"},
				{Name: "123foo"},
				{Name: "barfoo"},
				{Name: "afoo"},
				completeIP,
			},
			notMatches: []Instrument{
				{},
				{Name: "bar"},
				{Name: "foobar"},
				{Name: "barfoobaz"},
			},
		},
		{
			name:     "Name/Wildcard/Back**",
			criteria: Instrument{Name: "foo**"},
			matches: []Instrument{
				{Name: "foo"},
				{Name: "foo1"},
				{Name: "fooa"},
				{Name: "foobar"},
				completeIP,
			},
			notMatches: []Instrument{
				{},
				{Name: "bar"},
				{Name: "barfoo"},
				{Name: "barfoobaz"},
			},
		},
		{
			name:     "Name/Wildcard/Front*?",
			criteria: Instrument{Name: "*?oo"},
			matches: []Instrument{
				{Name: "foo"},
				{Name: "123foo"},
				{Name: "barfoo"},
				{Name: "afoo"},
				completeIP,
			},
			notMatches: []Instrument{
				{},
				{Name: "fo"},
				{Name: "bar"},
				{Name: "foobar"},
				{Name: "barfoobaz"},
			},
		},
		{
			name:     "Name/Wildcard/Back*?",
			criteria: Instrument{Name: "fo*?"},
			matches: []Instrument{
				{Name: "foo"},
				{Name: "foo1"},
				{Name: "fooa"},
				{Name: "foobar"},
				completeIP,
			},
			notMatches: []Instrument{
				{},
				{Name: "bar"},
				{Name: "barfoo"},
				{Name: "barfoobaz"},
			},
		},
		{
			name:     "Name/Wildcard/Front?*",
			criteria: Instrument{Name: "?*oo"},
			matches: []Instrument{
				{Name: "foo"},
				{Name: "123foo"},
				{Name: "barfoo"},
				{Name: "afoo"},
				completeIP,
			},
			notMatches: []Instrument{
				{},
				{Name: "oo"},
				{Name: "fo"},
				{Name: "bar"},
				{Name: "foobar"},
				{Name: "barfoobaz"},
			},
		},
		{
			name:     "Name/Wildcard/Back?*",
			criteria: Instrument{Name: "fo?*"},
			matches: []Instrument{
				{Name: "foo"},
				{Name: "foo1"},
				{Name: "fooa"},
				{Name: "foobar"},
				completeIP,
			},
			notMatches: []Instrument{
				{},
				{Name: "fo"},
				{Name: "bar"},
				{Name: "barfoo"},
				{Name: "barfoobaz"},
			},
		},
		{
			name:     "Name/Wildcard/Middle*",
			criteria: Instrument{Name: "f*o"},
			matches: []Instrument{
				{Name: "fo"},
				{Name: "foo"},
				{Name: "fooo"},
				{Name: "fo12baro"},
				completeIP,
			},
			notMatches: []Instrument{
				{},
				{Name: "bar"},
				{Name: "barfoo"},
				{Name: "barfoobaz"},
			},
		},
		{
			name:     "Name/Wildcard/Middle?",
			criteria: Instrument{Name: "f?o"},
			matches: []Instrument{
				{Name: "foo"},
				{Name: "f1o"},
				completeIP,
			},
			notMatches: []Instrument{
				{},
				{Name: "fo"},
				{Name: "fooo"},
				{Name: "fo12baro"},
				{Name: "bar"},
			},
		},
		{
			name:     "Name/Wildcard/MetaCharacters",
			criteria: Instrument{Name: "*.+()|[]{}^$-_?"},
			matches: []Instrument{
				{Name: "aa.+()|[]{}^$-_b"},
				{Name: ".+()|[]{}^$-_b"},
			},
			notMatches: []Instrument{
				{},
				{Name: "foo"},
				{Name: ".+()|[]{}^$-_"},
			},
		},
		{
			name:     "Description",
			criteria: Instrument{Description: "foo desc"},
			matches:  []Instrument{{Description: "foo desc"}, completeIP},
			notMatches: []Instrument{
				{},
				{Description: "foo"},
				{Description: "desc"},
			},
		},
		{
			name:     "Kind",
			criteria: Instrument{Kind: InstrumentKindSyncCounter},
			matches:  []Instrument{{Kind: InstrumentKindSyncCounter}, completeIP},
			notMatches: []Instrument{
				{},
				{Kind: InstrumentKindSyncUpDownCounter},
				{Kind: InstrumentKindSyncHistogram},
				{Kind: InstrumentKindAsyncCounter},
				{Kind: InstrumentKindAsyncUpDownCounter},
				{Kind: InstrumentKindAsyncGauge},
			},
		},
		{
			name:     "Unit",
			criteria: Instrument{Unit: unit.Bytes},
			matches:  []Instrument{{Unit: unit.Bytes}, completeIP},
			notMatches: []Instrument{
				{},
				{Unit: unit.Dimensionless},
				{Unit: unit.Unit("K")},
			},
		},
		{
			name: "ScopeName",
			criteria: Instrument{
				Scope: instrumentation.Scope{Name: "TestNewViewMatch"},
			},
			matches: []Instrument{
				{Scope: instrumentation.Scope{Name: "TestNewViewMatch"}},
				completeIP,
			},
			notMatches: []Instrument{
				{},
				{Scope: instrumentation.Scope{Name: "PrefixTestNewViewMatch"}},
				{Scope: instrumentation.Scope{Name: "TestNewViewMatchSuffix"}},
				{Scope: instrumentation.Scope{Name: "alt"}},
			},
		},
		{
			name: "ScopeVersion",
			criteria: Instrument{
				Scope: instrumentation.Scope{Version: "v0.1.0"},
			},
			matches: []Instrument{
				{Scope: instrumentation.Scope{Version: "v0.1.0"}},
				completeIP,
			},
			notMatches: []Instrument{
				{},
				{Scope: instrumentation.Scope{Version: "v0.1.0-RC1"}},
				{Scope: instrumentation.Scope{Version: "v0.1.1"}},
			},
		},
		{
			name: "ScopeSchemaURL",
			criteria: Instrument{
				Scope: instrumentation.Scope{SchemaURL: schemaURL},
			},
			matches: []Instrument{
				{Scope: instrumentation.Scope{SchemaURL: schemaURL}},
				completeIP,
			},
			notMatches: []Instrument{
				{},
				{Scope: instrumentation.Scope{SchemaURL: schemaURL + "/path"}},
				{Scope: instrumentation.Scope{SchemaURL: "https://go.dev"}},
			},
		},
		{
			name: "Scope",
			criteria: Instrument{
				Scope: instrumentation.Scope{
					Name:      "TestNewViewMatch",
					Version:   "v0.1.0",
					SchemaURL: schemaURL,
				},
			},
			matches: []Instrument{
				{
					Scope: instrumentation.Scope{
						Name:      "TestNewViewMatch",
						Version:   "v0.1.0",
						SchemaURL: schemaURL,
					},
				},
				completeIP,
			},
			notMatches: []Instrument{
				{},
				{
					Scope: instrumentation.Scope{
						Name:      "CompleteMisMatch",
						Version:   "v0.2.0",
						SchemaURL: "https://go.dev",
					},
				},
				{
					Scope: instrumentation.Scope{
						Name:      "NameMisMatch",
						Version:   "v0.1.0",
						SchemaURL: schemaURL,
					},
				},
			},
		},
		{
			name:     "Complete",
			criteria: completeIP,
			matches:  []Instrument{completeIP},
			notMatches: []Instrument{
				{},
				{Name: "foo"},
				{
					Name:        "Wrong Name",
					Description: "foo desc",
					Kind:        InstrumentKindSyncCounter,
					Unit:        unit.Bytes,
					Scope: instrumentation.Scope{
						Name:      "TestNewViewMatch",
						Version:   "v0.1.0",
						SchemaURL: schemaURL,
					},
				},
				{
					Name:        "foo",
					Description: "Wrong Description",
					Kind:        InstrumentKindSyncCounter,
					Unit:        unit.Bytes,
					Scope: instrumentation.Scope{
						Name:      "TestNewViewMatch",
						Version:   "v0.1.0",
						SchemaURL: schemaURL,
					},
				},
				{
					Name:        "foo",
					Description: "foo desc",
					Kind:        InstrumentKindAsyncUpDownCounter,
					Unit:        unit.Bytes,
					Scope: instrumentation.Scope{
						Name:      "TestNewViewMatch",
						Version:   "v0.1.0",
						SchemaURL: schemaURL,
					},
				},
				{
					Name:        "foo",
					Description: "foo desc",
					Kind:        InstrumentKindSyncCounter,
					Unit:        unit.Dimensionless,
					Scope: instrumentation.Scope{
						Name:      "TestNewViewMatch",
						Version:   "v0.1.0",
						SchemaURL: schemaURL,
					},
				},
				{
					Name:        "foo",
					Description: "foo desc",
					Kind:        InstrumentKindSyncCounter,
					Unit:        unit.Bytes,
					Scope: instrumentation.Scope{
						Name:      "Wrong Scope Name",
						Version:   "v0.1.0",
						SchemaURL: schemaURL,
					},
				},
				{
					Name:        "foo",
					Description: "foo desc",
					Kind:        InstrumentKindSyncCounter,
					Unit:        unit.Bytes,
					Scope: instrumentation.Scope{
						Name:      "TestNewViewMatch",
						Version:   "v1.4.3",
						SchemaURL: schemaURL,
					},
				},
				{
					Name:        "foo",
					Description: "foo desc",
					Kind:        InstrumentKindSyncCounter,
					Unit:        unit.Bytes,
					Scope: instrumentation.Scope{
						Name:      "TestNewViewMatch",
						Version:   "v0.1.0",
						SchemaURL: "https://go.dev",
					},
				},
			},
		},
	}

	for _, test := range tests {
		v := NewView(test.criteria, Stream{})
		t.Run(test.name, func(t *testing.T) {
			for _, instrument := range test.matches {
				_, matches := v(instrument)
				assert.Truef(t, matches, "view does not match %#v", instrument)
			}

			for _, instrument := range test.notMatches {
				_, matches := v(instrument)
				assert.Falsef(t, matches, "view matches %#v", instrument)
			}
		})
	}
}

func TestNewViewReplace(t *testing.T) {
	alt := "alternative value"
	tests := []struct {
		name string
		mask Stream
		want func(Instrument) Stream
	}{
		{
			name: "Nothing",
			want: func(ip Instrument) Stream {
				return Stream{Instrument: ip}
			},
		},
		{
			name: "Name",
			mask: Stream{
				Instrument: Instrument{Name: alt},
			},
			want: func(ip Instrument) Stream {
				ip.Name = alt
				return Stream{Instrument: ip}
			},
		},
		{
			name: "Description",
			mask: Stream{
				Instrument: Instrument{
					Description: alt,
				},
			},
			want: func(ip Instrument) Stream {
				ip.Description = alt
				return Stream{Instrument: ip}
			},
		},
		{
			name: "Kind",
			mask: Stream{
				Instrument: Instrument{
					Kind: InstrumentKindAsyncUpDownCounter,
				},
			},
			want: func(ip Instrument) Stream {
				ip.Kind = InstrumentKindAsyncUpDownCounter
				return Stream{Instrument: ip}
			},
		},
		{
			name: "Unit",
			mask: Stream{
				Instrument: Instrument{
					Unit: unit.Dimensionless,
				},
			},
			want: func(ip Instrument) Stream {
				ip.Unit = unit.Dimensionless
				return Stream{Instrument: ip}
			},
		},
		{
			name: "ScopeName",
			mask: Stream{
				Instrument: Instrument{
					Scope: instrumentation.Scope{Name: alt},
				},
			},
			want: func(ip Instrument) Stream {
				ip.Scope.Name = alt
				return Stream{Instrument: ip}
			},
		},
		{
			name: "ScopeVersion",
			mask: Stream{
				Instrument: Instrument{
					Scope: instrumentation.Scope{Version: alt},
				},
			},
			want: func(ip Instrument) Stream {
				ip.Scope.Version = alt
				return Stream{Instrument: ip}
			},
		},
		{
			name: "ScopeSchemaURL",
			mask: Stream{
				Instrument: Instrument{
					Scope: instrumentation.Scope{SchemaURL: alt},
				},
			},
			want: func(ip Instrument) Stream {
				ip.Scope.SchemaURL = alt
				return Stream{Instrument: ip}
			},
		},
		{
			name: "Scope",
			mask: Stream{
				Instrument: Instrument{
					Scope: instrumentation.Scope{
						Name:      "Alt Scope Name",
						Version:   "1.1.1",
						SchemaURL: "https://go.dev",
					},
				},
			},
			want: func(ip Instrument) Stream {
				ip.Scope.Name = "Alt Scope Name"
				ip.Scope.Version = "1.1.1"
				ip.Scope.SchemaURL = "https://go.dev"
				return Stream{Instrument: ip}
			},
		},
		{
			name: "Aggregation",
			mask: Stream{
				Aggregation: aggregation.LastValue{},
			},
			want: func(ip Instrument) Stream {
				return Stream{
					Instrument:  ip,
					Aggregation: aggregation.LastValue{},
				}
			},
		},
		{
			name: "Complete",
			mask: Stream{
				Instrument: Instrument{
					Name:        alt,
					Description: alt,
					Kind:        InstrumentKindAsyncUpDownCounter,
					Unit:        unit.Dimensionless,
					Scope: instrumentation.Scope{
						Name:      alt,
						Version:   alt,
						SchemaURL: alt,
					},
				},
				Aggregation: aggregation.LastValue{},
			},
			want: func(ip Instrument) Stream {
				ip.Name = alt
				ip.Description = alt
				ip.Kind = InstrumentKindAsyncUpDownCounter
				ip.Unit = unit.Dimensionless
				ip.Scope.Name = alt
				ip.Scope.Version = alt
				ip.Scope.SchemaURL = alt
				return Stream{
					Instrument:  ip,
					Aggregation: aggregation.LastValue{},
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, match := NewView(completeIP, test.mask)(completeIP)
			require.True(t, match, "view did not match exact criteria")
			assert.Equal(t, test.want(completeIP), got)
		})
	}

	// Go does not allow for the comparison of function values, even their
	// addresses. Therefore, the AttributeFilter field needs an alternative
	// testing strategy.
	t.Run("AttributeFilter", func(t *testing.T) {
		allowed := attribute.String("key", "val")
		filter := func(kv attribute.KeyValue) bool {
			return kv == allowed
		}
		mask := Stream{AttributeFilter: filter}
		got, match := NewView(completeIP, mask)(completeIP)
		require.True(t, match, "view did not match exact criteria")
		require.NotNil(t, got.AttributeFilter, "AttributeFilter not set")
		assert.True(t, got.AttributeFilter(allowed), "wrong AttributeFilter")
		other := attribute.String("key", "other val")
		assert.False(t, got.AttributeFilter(other), "wrong AttributeFilter")
	})
}

type badAgg struct {
	aggregation.Aggregation

	err error
}

func (a badAgg) Copy() aggregation.Aggregation {
	return a
}

func (a badAgg) Err() error {
	return a.err
}

func TestNewViewAggregationErrorLogged(t *testing.T) {
	tLog := testr.NewWithOptions(t, testr.Options{Verbosity: 6})
	l := &logCounter{LogSink: tLog.GetSink()}
	otel.SetLogger(logr.New(l))

	agg := badAgg{err: assert.AnError}
	mask := Stream{Aggregation: agg}
	got, match := NewView(completeIP, mask)(completeIP)
	require.True(t, match, "view did not match exact criteria")
	assert.Nil(t, got.Aggregation, "erroring aggregation used")
	assert.Equal(t, 1, l.ErrorN())
}

func ExampleNewView() {
	// Rename the "latency" instrument from the v0.34.0 version of the "http"
	// instrumentation library as "request.latency".
	v := NewView(Instrument{
		Name: "latency",
		Scope: instrumentation.Scope{
			Name:    "http",
			Version: "v0.34.0",
		},
	}, Stream{
		Instrument: Instrument{Name: "request.latency"},
	})

	stream, _ := v(Instrument{
		Name:        "latency",
		Description: "request latency",
		Unit:        unit.Milliseconds,
		Kind:        InstrumentKindSyncCounter,
		Scope: instrumentation.Scope{
			Name:      "http",
			Version:   "v0.34.0",
			SchemaURL: "https://opentelemetry.io/schemas/1.0.0",
		},
	})
	fmt.Println("name:", stream.Name)
	fmt.Println("description:", stream.Description)
	fmt.Println("unit:", stream.Unit)
	fmt.Println("kind:", stream.Kind)
	fmt.Println("scope:", stream.Scope)
	// Output:
	// name: request.latency
	// description: request latency
	// unit: ms
	// kind: 1
	// scope: {http v0.34.0 https://opentelemetry.io/schemas/1.0.0}
}

func ExampleNewView_drop() {
	// Set the drop aggregator for all instrumentation from the "db" library.
	v := NewView(
		Instrument{Scope: instrumentation.Scope{Name: "db"}},
		Stream{Aggregation: aggregation.Drop{}},
	)

	stream, _ := v(Instrument{
		Name:  "queries",
		Kind:  InstrumentKindSyncCounter,
		Scope: instrumentation.Scope{Name: "db", Version: "v0.4.0"},
	})
	fmt.Println("name:", stream.Name)
	fmt.Println("scope:", stream.Scope)
	fmt.Printf("aggregation: %#v", stream.Aggregation)
	// Output:
	// name: queries
	// scope: {db v0.4.0 }
	// aggregation: aggregation.Drop{}
}

func ExampleNewView_wildcard() {
	// Set unit to milliseconds for any instrument with a name suffix of ".ms".
	v := NewView(
		Instrument{Name: "*.ms"},
		Stream{
			Instrument: Instrument{
				Unit: unit.Milliseconds,
			},
		},
	)

	stream, _ := v(Instrument{
		Name: "computation.time.ms",
		Unit: unit.Dimensionless,
	})
	fmt.Println("name:", stream.Name)
	fmt.Println("unit:", stream.Unit)
	// Output:
	// name: computation.time.ms
	// unit: ms
}
