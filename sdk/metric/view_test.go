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
	completeIP = InstrumentProperties{
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
		criteria   InstrumentProperties
		matches    []InstrumentProperties
		notMatches []InstrumentProperties
	}{
		{
			name:       "Empty",
			notMatches: []InstrumentProperties{{}, {Name: "foo"}, completeIP},
		},
		{
			name:     "Name/Exact",
			criteria: InstrumentProperties{Name: "foo"},
			matches:  []InstrumentProperties{{Name: "foo"}, completeIP},
			notMatches: []InstrumentProperties{
				{},
				{Name: "bar"},
				{Name: "foobar"},
				{Name: "barfoo"},
				{Name: "ffooo"},
			},
		},
		{
			name:     "Name/Wildcard/*",
			criteria: InstrumentProperties{Name: "*"},
			matches: []InstrumentProperties{
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
			criteria: InstrumentProperties{Name: "?oo"},
			matches: []InstrumentProperties{
				{Name: "foo"},
				{Name: "1oo"},
				completeIP,
			},
			notMatches: []InstrumentProperties{
				{},
				{Name: "bar"},
				{Name: "foobar"},
				{Name: "barfoo"},
				{Name: "barfoobaz"},
			},
		},
		{
			name:     "Name/Wildcard/Back?",
			criteria: InstrumentProperties{Name: "fo?"},
			matches: []InstrumentProperties{
				{Name: "foo"},
				{Name: "fo1"},
				completeIP,
			},
			notMatches: []InstrumentProperties{
				{},
				{Name: "bar"},
				{Name: "foobar"},
				{Name: "barfoo"},
				{Name: "barfoobaz"},
			},
		},
		{
			name:     "Name/Wildcard/Front*",
			criteria: InstrumentProperties{Name: "*foo"},
			matches: []InstrumentProperties{
				{Name: "foo"},
				{Name: "123foo"},
				{Name: "barfoo"},
				completeIP,
			},
			notMatches: []InstrumentProperties{
				{},
				{Name: "bar"},
				{Name: "foobar"},
				{Name: "barfoobaz"},
			},
		},
		{
			name:     "Name/Wildcard/Back*",
			criteria: InstrumentProperties{Name: "foo*"},
			matches: []InstrumentProperties{
				{Name: "foo"},
				{Name: "foo1"},
				{Name: "foobar"},
				completeIP,
			},
			notMatches: []InstrumentProperties{
				{},
				{Name: "bar"},
				{Name: "barfoo"},
				{Name: "barfoobaz"},
			},
		},
		{
			name:     "Name/Wildcard/FrontBack*",
			criteria: InstrumentProperties{Name: "*foo*"},
			matches: []InstrumentProperties{
				{Name: "foo"},
				{Name: "foo1"},
				{Name: "1foo"},
				{Name: "1foo1"},
				{Name: "foobar"},
				{Name: "barfoobaz"},
				completeIP,
			},
			notMatches: []InstrumentProperties{
				{},
				{Name: "bar"},
			},
		},
		{
			name:     "Name/Wildcard/Front**",
			criteria: InstrumentProperties{Name: "**foo"},
			matches: []InstrumentProperties{
				{Name: "foo"},
				{Name: "123foo"},
				{Name: "barfoo"},
				{Name: "afoo"},
				completeIP,
			},
			notMatches: []InstrumentProperties{
				{},
				{Name: "bar"},
				{Name: "foobar"},
				{Name: "barfoobaz"},
			},
		},
		{
			name:     "Name/Wildcard/Back**",
			criteria: InstrumentProperties{Name: "foo**"},
			matches: []InstrumentProperties{
				{Name: "foo"},
				{Name: "foo1"},
				{Name: "fooa"},
				{Name: "foobar"},
				completeIP,
			},
			notMatches: []InstrumentProperties{
				{},
				{Name: "bar"},
				{Name: "barfoo"},
				{Name: "barfoobaz"},
			},
		},
		{
			name:     "Name/Wildcard/Front*?",
			criteria: InstrumentProperties{Name: "*?oo"},
			matches: []InstrumentProperties{
				{Name: "foo"},
				{Name: "123foo"},
				{Name: "barfoo"},
				{Name: "afoo"},
				completeIP,
			},
			notMatches: []InstrumentProperties{
				{},
				{Name: "fo"},
				{Name: "bar"},
				{Name: "foobar"},
				{Name: "barfoobaz"},
			},
		},
		{
			name:     "Name/Wildcard/Back*?",
			criteria: InstrumentProperties{Name: "fo*?"},
			matches: []InstrumentProperties{
				{Name: "foo"},
				{Name: "foo1"},
				{Name: "fooa"},
				{Name: "foobar"},
				completeIP,
			},
			notMatches: []InstrumentProperties{
				{},
				{Name: "bar"},
				{Name: "barfoo"},
				{Name: "barfoobaz"},
			},
		},
		{
			name:     "Name/Wildcard/Front?*",
			criteria: InstrumentProperties{Name: "?*oo"},
			matches: []InstrumentProperties{
				{Name: "foo"},
				{Name: "123foo"},
				{Name: "barfoo"},
				{Name: "afoo"},
				completeIP,
			},
			notMatches: []InstrumentProperties{
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
			criteria: InstrumentProperties{Name: "fo?*"},
			matches: []InstrumentProperties{
				{Name: "foo"},
				{Name: "foo1"},
				{Name: "fooa"},
				{Name: "foobar"},
				completeIP,
			},
			notMatches: []InstrumentProperties{
				{},
				{Name: "fo"},
				{Name: "bar"},
				{Name: "barfoo"},
				{Name: "barfoobaz"},
			},
		},
		{
			name:     "Name/Wildcard/Middle*",
			criteria: InstrumentProperties{Name: "f*o"},
			matches: []InstrumentProperties{
				{Name: "fo"},
				{Name: "foo"},
				{Name: "fooo"},
				{Name: "fo12baro"},
				completeIP,
			},
			notMatches: []InstrumentProperties{
				{},
				{Name: "bar"},
				{Name: "barfoo"},
				{Name: "barfoobaz"},
			},
		},
		{
			name:     "Name/Wildcard/Middle?",
			criteria: InstrumentProperties{Name: "f?o"},
			matches: []InstrumentProperties{
				{Name: "foo"},
				{Name: "f1o"},
				completeIP,
			},
			notMatches: []InstrumentProperties{
				{},
				{Name: "fo"},
				{Name: "fooo"},
				{Name: "fo12baro"},
				{Name: "bar"},
			},
		},
		{
			name:     "Name/Wildcard/MetaCharacters",
			criteria: InstrumentProperties{Name: "*.+()|[]{}^$-_?"},
			matches: []InstrumentProperties{
				{Name: "aa.+()|[]{}^$-_b"},
				{Name: ".+()|[]{}^$-_b"},
			},
			notMatches: []InstrumentProperties{
				{},
				{Name: "foo"},
				{Name: ".+()|[]{}^$-_"},
			},
		},
		{
			name:     "Description",
			criteria: InstrumentProperties{Description: "foo desc"},
			matches:  []InstrumentProperties{{Description: "foo desc"}, completeIP},
			notMatches: []InstrumentProperties{
				{},
				{Description: "foo"},
				{Description: "desc"},
			},
		},
		{
			name:     "Kind",
			criteria: InstrumentProperties{Kind: InstrumentKindSyncCounter},
			matches:  []InstrumentProperties{{Kind: InstrumentKindSyncCounter}, completeIP},
			notMatches: []InstrumentProperties{
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
			criteria: InstrumentProperties{Unit: unit.Bytes},
			matches:  []InstrumentProperties{{Unit: unit.Bytes}, completeIP},
			notMatches: []InstrumentProperties{
				{},
				{Unit: unit.Dimensionless},
				{Unit: unit.Unit("K")},
			},
		},
		{
			name: "ScopeName",
			criteria: InstrumentProperties{
				Scope: instrumentation.Scope{Name: "TestNewViewMatch"},
			},
			matches: []InstrumentProperties{
				{Scope: instrumentation.Scope{Name: "TestNewViewMatch"}},
				completeIP,
			},
			notMatches: []InstrumentProperties{
				{},
				{Scope: instrumentation.Scope{Name: "PrefixTestNewViewMatch"}},
				{Scope: instrumentation.Scope{Name: "TestNewViewMatchSuffix"}},
				{Scope: instrumentation.Scope{Name: "alt"}},
			},
		},
		{
			name: "ScopeVersion",
			criteria: InstrumentProperties{
				Scope: instrumentation.Scope{Version: "v0.1.0"},
			},
			matches: []InstrumentProperties{
				{Scope: instrumentation.Scope{Version: "v0.1.0"}},
				completeIP,
			},
			notMatches: []InstrumentProperties{
				{},
				{Scope: instrumentation.Scope{Version: "v0.1.0-RC1"}},
				{Scope: instrumentation.Scope{Version: "v0.1.1"}},
			},
		},
		{
			name: "ScopeSchemaURL",
			criteria: InstrumentProperties{
				Scope: instrumentation.Scope{SchemaURL: schemaURL},
			},
			matches: []InstrumentProperties{
				{Scope: instrumentation.Scope{SchemaURL: schemaURL}},
				completeIP,
			},
			notMatches: []InstrumentProperties{
				{},
				{Scope: instrumentation.Scope{SchemaURL: schemaURL + "/path"}},
				{Scope: instrumentation.Scope{SchemaURL: "https://go.dev"}},
			},
		},
		{
			name: "Scope",
			criteria: InstrumentProperties{
				Scope: instrumentation.Scope{
					Name:      "TestNewViewMatch",
					Version:   "v0.1.0",
					SchemaURL: schemaURL,
				},
			},
			matches: []InstrumentProperties{
				{
					Scope: instrumentation.Scope{
						Name:      "TestNewViewMatch",
						Version:   "v0.1.0",
						SchemaURL: schemaURL,
					},
				},
				completeIP,
			},
			notMatches: []InstrumentProperties{
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
			matches:  []InstrumentProperties{completeIP},
			notMatches: []InstrumentProperties{
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
		v := NewView(test.criteria, DataStream{})
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
		mask DataStream
		want func(InstrumentProperties) DataStream
	}{
		{
			name: "Nothing",
			want: func(ip InstrumentProperties) DataStream {
				return DataStream{InstrumentProperties: ip}
			},
		},
		{
			name: "Name",
			mask: DataStream{
				InstrumentProperties: InstrumentProperties{Name: alt},
			},
			want: func(ip InstrumentProperties) DataStream {
				ip.Name = alt
				return DataStream{InstrumentProperties: ip}
			},
		},
		{
			name: "Description",
			mask: DataStream{
				InstrumentProperties: InstrumentProperties{
					Description: alt,
				},
			},
			want: func(ip InstrumentProperties) DataStream {
				ip.Description = alt
				return DataStream{InstrumentProperties: ip}
			},
		},
		{
			name: "Kind",
			mask: DataStream{
				InstrumentProperties: InstrumentProperties{
					Kind: InstrumentKindAsyncUpDownCounter,
				},
			},
			want: func(ip InstrumentProperties) DataStream {
				ip.Kind = InstrumentKindAsyncUpDownCounter
				return DataStream{InstrumentProperties: ip}
			},
		},
		{
			name: "Unit",
			mask: DataStream{
				InstrumentProperties: InstrumentProperties{
					Unit: unit.Dimensionless,
				},
			},
			want: func(ip InstrumentProperties) DataStream {
				ip.Unit = unit.Dimensionless
				return DataStream{InstrumentProperties: ip}
			},
		},
		{
			name: "ScopeName",
			mask: DataStream{
				InstrumentProperties: InstrumentProperties{
					Scope: instrumentation.Scope{Name: alt},
				},
			},
			want: func(ip InstrumentProperties) DataStream {
				ip.Scope.Name = alt
				return DataStream{InstrumentProperties: ip}
			},
		},
		{
			name: "ScopeVersion",
			mask: DataStream{
				InstrumentProperties: InstrumentProperties{
					Scope: instrumentation.Scope{Version: alt},
				},
			},
			want: func(ip InstrumentProperties) DataStream {
				ip.Scope.Version = alt
				return DataStream{InstrumentProperties: ip}
			},
		},
		{
			name: "ScopeSchemaURL",
			mask: DataStream{
				InstrumentProperties: InstrumentProperties{
					Scope: instrumentation.Scope{SchemaURL: alt},
				},
			},
			want: func(ip InstrumentProperties) DataStream {
				ip.Scope.SchemaURL = alt
				return DataStream{InstrumentProperties: ip}
			},
		},
		{
			name: "Scope",
			mask: DataStream{
				InstrumentProperties: InstrumentProperties{
					Scope: instrumentation.Scope{
						Name:      "Alt Scope Name",
						Version:   "1.1.1",
						SchemaURL: "https://go.dev",
					},
				},
			},
			want: func(ip InstrumentProperties) DataStream {
				ip.Scope.Name = "Alt Scope Name"
				ip.Scope.Version = "1.1.1"
				ip.Scope.SchemaURL = "https://go.dev"
				return DataStream{InstrumentProperties: ip}
			},
		},
		{
			name: "Aggregation",
			mask: DataStream{
				Aggregation: aggregation.LastValue{},
			},
			want: func(ip InstrumentProperties) DataStream {
				return DataStream{
					InstrumentProperties: ip,
					Aggregation:          aggregation.LastValue{},
				}
			},
		},
		{
			name: "Complete",
			mask: DataStream{
				InstrumentProperties: InstrumentProperties{
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
			want: func(ip InstrumentProperties) DataStream {
				ip.Name = alt
				ip.Description = alt
				ip.Kind = InstrumentKindAsyncUpDownCounter
				ip.Unit = unit.Dimensionless
				ip.Scope.Name = alt
				ip.Scope.Version = alt
				ip.Scope.SchemaURL = alt
				return DataStream{
					InstrumentProperties: ip,
					Aggregation:          aggregation.LastValue{},
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
		mask := DataStream{AttributeFilter: filter}
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
	mask := DataStream{Aggregation: agg}
	got, match := NewView(completeIP, mask)(completeIP)
	require.True(t, match, "view did not match exact criteria")
	assert.Nil(t, got.Aggregation, "erroring aggregation used")
	assert.Equal(t, 1, l.ErrorN())
}
