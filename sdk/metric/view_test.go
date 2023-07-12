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
	"regexp"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/funcr"
	"github.com/go-logr/logr/testr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
)

var (
	schemaURL  = "https://opentelemetry.io/schemas/1.0.0"
	completeIP = Instrument{
		Name:        "foo",
		Description: "foo desc",
		Kind:        InstrumentKindCounter,
		Unit:        "By",
		Scope: instrumentation.Scope{
			Name:      "TestNewViewMatch",
			Version:   "v0.1.0",
			SchemaURL: schemaURL,
		},
	}
)

func scope(name, ver, url string) instrumentation.Scope {
	return instrumentation.Scope{Name: name, Version: ver, SchemaURL: url}
}

func testNewViewMatchName() func(t *testing.T) {
	tests := []struct {
		name     string
		criteria string
		match    []string
		notMatch []string
	}{
		{
			name:     "Exact",
			criteria: "foo",
			match:    []string{"foo"},
			notMatch: []string{"", "bar", "foobar", "barfoo", "ffooo"},
		},
		{
			name:     "Wildcard/*",
			criteria: "*",
			match:    []string{"", "foo", "foobar", "barfoo", "barfoobaz"},
		},
		{
			name:     "Wildcard/Front?",
			criteria: "?oo",
			match:    []string{"foo", "1oo"},
			notMatch: []string{"", "bar", "foobar", "barfoo", "barfoobaz"},
		},
		{
			name:     "Wildcard/Back?",
			criteria: "fo?",
			match:    []string{"foo", "fo1"},
			notMatch: []string{"", "bar", "foobar", "barfoo", "barfoobaz"},
		},
		{
			name:     "Wildcard/Front*",
			criteria: "*foo",
			match:    []string{"foo", "123foo", "barfoo"},
			notMatch: []string{"", "bar", "foobar", "barfoobaz"},
		},
		{
			name:     "Wildcard/Back*",
			criteria: "foo*",
			match:    []string{"foo", "foo1", "foobar"},
			notMatch: []string{"", "bar", "barfoo", "barfoobaz"},
		},
		{
			name:     "Wildcard/FrontBack*",
			criteria: "*foo*",
			match:    []string{"foo", "foo1", "1foo", "1foo1", "foobar", "barfoobaz"},
			notMatch: []string{"", "bar"},
		},
		{
			name:     "Wildcard/Front**",
			criteria: "**foo",
			match:    []string{"foo", "123foo", "barfoo", "afoo"},
			notMatch: []string{"", "bar", "foobar", "barfoobaz"},
		},
		{
			name:     "Wildcard/Back**",
			criteria: "foo**",
			match:    []string{"foo", "foo1", "fooa", "foobar"},
			notMatch: []string{"", "bar", "barfoo", "barfoobaz"},
		},
		{
			name:     "Wildcard/Front*?",
			criteria: "*?oo",
			match:    []string{"foo", "123foo", "barfoo", "afoo"},
			notMatch: []string{"", "fo", "bar", "foobar", "barfoobaz"},
		},
		{
			name:     "Wildcard/Back*?",
			criteria: "fo*?",
			match:    []string{"foo", "foo1", "fooa", "foobar"},
			notMatch: []string{"", "bar", "barfoo", "barfoobaz"},
		},
		{
			name:     "Wildcard/Front?*",
			criteria: "?*oo",
			match:    []string{"foo", "123foo", "barfoo", "afoo"},
			notMatch: []string{"", "oo", "fo", "bar", "foobar", "barfoobaz"},
		},
		{
			name:     "Wildcard/Back?*",
			criteria: "fo?*",
			match:    []string{"foo", "foo1", "fooa", "foobar"},
			notMatch: []string{"", "fo", "bar", "barfoo", "barfoobaz"},
		},
		{
			name:     "Wildcard/Middle*",
			criteria: "f*o",
			match:    []string{"fo", "foo", "fooo", "fo12baro"},
			notMatch: []string{"", "bar", "barfoo", "barfoobaz"},
		},
		{
			name:     "Wildcard/Middle?",
			criteria: "f?o",
			match:    []string{"foo", "f1o"},
			notMatch: []string{"", "fo", "fooo", "fo12baro", "bar"},
		},
		{
			name:     "Wildcard/MetaCharacters",
			criteria: "*.+()|[]{}^$-_?",
			match:    []string{"aa.+()|[]{}^$-_b", ".+()|[]{}^$-_b"},
			notMatch: []string{"", "foo", ".+()|[]{}^$-_"},
		},
	}

	return func(t *testing.T) {
		for _, test := range tests {
			v := NewView(Instrument{Name: test.criteria}, Stream{})
			t.Run(test.name, func(t *testing.T) {
				for _, n := range test.match {
					_, matches := v(Instrument{Name: n})
					assert.Truef(t, matches, "%s does not match %s", test.criteria, n)
				}
				for _, n := range test.notMatch {
					_, matches := v(Instrument{Name: n})
					assert.Falsef(t, matches, "%s matches %s", test.criteria, n)
				}
			})
		}
	}
}

func TestNewViewMatch(t *testing.T) {
	// Avoid boilerplate for name match testing.
	t.Run("Name", testNewViewMatchName())

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
			name:       "Description",
			criteria:   Instrument{Description: "foo desc"},
			matches:    []Instrument{{Description: "foo desc"}, completeIP},
			notMatches: []Instrument{{}, {Description: "foo"}, {Description: "desc"}},
		},
		{
			name:     "Kind",
			criteria: Instrument{Kind: InstrumentKindCounter},
			matches:  []Instrument{{Kind: InstrumentKindCounter}, completeIP},
			notMatches: []Instrument{
				{},
				{Kind: InstrumentKindUpDownCounter},
				{Kind: InstrumentKindHistogram},
				{Kind: InstrumentKindObservableCounter},
				{Kind: InstrumentKindObservableUpDownCounter},
				{Kind: InstrumentKindObservableGauge},
			},
		},
		{
			name:     "Unit",
			criteria: Instrument{Unit: "By"},
			matches:  []Instrument{{Unit: "By"}, completeIP},
			notMatches: []Instrument{
				{},
				{Unit: "1"},
				{Unit: "K"},
			},
		},
		{
			name:     "ScopeName",
			criteria: Instrument{Scope: scope("TestNewViewMatch", "", "")},
			matches: []Instrument{
				{Scope: scope("TestNewViewMatch", "", "")},
				completeIP,
			},
			notMatches: []Instrument{
				{},
				{Scope: scope("PrefixTestNewViewMatch", "", "")},
				{Scope: scope("TestNewViewMatchSuffix", "", "")},
				{Scope: scope("alt", "", "")},
			},
		},
		{
			name:     "ScopeVersion",
			criteria: Instrument{Scope: scope("", "v0.1.0", "")},
			matches: []Instrument{
				{Scope: scope("", "v0.1.0", "")},
				completeIP,
			},
			notMatches: []Instrument{
				{},
				{Scope: scope("", "v0.1.0-RC1", "")},
				{Scope: scope("", "v0.1.1", "")},
			},
		},
		{
			name:     "ScopeSchemaURL",
			criteria: Instrument{Scope: scope("", "", schemaURL)},
			matches: []Instrument{
				{Scope: scope("", "", schemaURL)},
				completeIP,
			},
			notMatches: []Instrument{
				{},
				{Scope: scope("", "", schemaURL+"/path")},
				{Scope: scope("", "", "https://go.dev")},
			},
		},
		{
			name:     "Scope",
			criteria: Instrument{Scope: scope("TestNewViewMatch", "v0.1.0", schemaURL)},
			matches: []Instrument{
				{Scope: scope("TestNewViewMatch", "v0.1.0", schemaURL)},
				completeIP,
			},
			notMatches: []Instrument{
				{},
				{Scope: scope("CompleteMisMatch", "v0.2.0", "https://go.dev")},
				{Scope: scope("NameMisMatch", "v0.1.0", schemaURL)},
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
					Kind:        InstrumentKindCounter,
					Unit:        "By",
					Scope:       scope("TestNewViewMatch", "v0.1.0", schemaURL),
				},
				{
					Name:        "foo",
					Description: "Wrong Description",
					Kind:        InstrumentKindCounter,
					Unit:        "By",
					Scope:       scope("TestNewViewMatch", "v0.1.0", schemaURL),
				},
				{
					Name:        "foo",
					Description: "foo desc",
					Kind:        InstrumentKindObservableUpDownCounter,
					Unit:        "By",
					Scope:       scope("TestNewViewMatch", "v0.1.0", schemaURL),
				},
				{
					Name:        "foo",
					Description: "foo desc",
					Kind:        InstrumentKindCounter,
					Unit:        "1",
					Scope:       scope("TestNewViewMatch", "v0.1.0", schemaURL),
				},
				{
					Name:        "foo",
					Description: "foo desc",
					Kind:        InstrumentKindCounter,
					Unit:        "By",
					Scope:       scope("Wrong Scope Name", "v0.1.0", schemaURL),
				},
				{
					Name:        "foo",
					Description: "foo desc",
					Kind:        InstrumentKindCounter,
					Unit:        "By",
					Scope:       scope("TestNewViewMatch", "v1.4.3", schemaURL),
				},
				{
					Name:        "foo",
					Description: "foo desc",
					Kind:        InstrumentKindCounter,
					Unit:        "By",
					Scope:       scope("TestNewViewMatch", "v0.1.0", "https://go.dev"),
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
			want: func(i Instrument) Stream {
				return Stream{
					Name:        i.Name,
					Description: i.Description,
					Unit:        i.Unit,
				}
			},
		},
		{
			name: "Name",
			mask: Stream{Name: alt},
			want: func(i Instrument) Stream {
				return Stream{
					Name:        alt,
					Description: i.Description,
					Unit:        i.Unit,
				}
			},
		},
		{
			name: "Description",
			mask: Stream{Description: alt},
			want: func(i Instrument) Stream {
				return Stream{
					Name:        i.Name,
					Description: alt,
					Unit:        i.Unit,
				}
			},
		},
		{
			name: "Unit",
			mask: Stream{Unit: "1"},
			want: func(i Instrument) Stream {
				return Stream{
					Name:        i.Name,
					Description: i.Description,
					Unit:        "1",
				}
			},
		},
		{
			name: "Aggregation",
			mask: Stream{Aggregation: aggregation.LastValue{}},
			want: func(i Instrument) Stream {
				return Stream{
					Name:        i.Name,
					Description: i.Description,
					Unit:        i.Unit,
					Aggregation: aggregation.LastValue{},
				}
			},
		},
		{
			name: "AttributeKeys",
			mask: Stream{AllowAttributeKeys: []attribute.Key{"test"}},
			want: func(i Instrument) Stream {
				return Stream{
					Name:               i.Name,
					Description:        i.Description,
					Unit:               i.Unit,
					AllowAttributeKeys: []attribute.Key{"test"},
				}
			},
		},
		{
			name: "Complete",
			mask: Stream{
				Name:        alt,
				Description: alt,
				Unit:        "1",
				Aggregation: aggregation.LastValue{},
			},
			want: func(i Instrument) Stream {
				return Stream{
					Name:        alt,
					Description: alt,
					Unit:        "1",
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
}

type badAgg struct {
	aggregation.Aggregation
	err error
}

func (a badAgg) Copy() aggregation.Aggregation { return a }

func (a badAgg) Err() error { return a.err }

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

func TestNewViewEmptyViewErrorLogged(t *testing.T) {
	var got string
	otel.SetLogger(funcr.New(func(_, args string) {
		got = args
	}, funcr.Options{Verbosity: 6}))

	_ = NewView(Instrument{}, Stream{})
	assert.Contains(t, got, errEmptyView.Error())
}

func TestNewViewMultiInstMatchErrorLogged(t *testing.T) {
	var got string
	otel.SetLogger(funcr.New(func(_, args string) {
		got = args
	}, funcr.Options{Verbosity: 6}))

	_ = NewView(Instrument{
		Name: "*", // Wildcard match name (multiple instruments).
	}, Stream{
		Name: "non-empty",
	})
	assert.Contains(t, got, errMultiInst.Error())
}

func ExampleNewView() {
	// Create a view that renames the "latency" instrument from the v0.34.0
	// version of the "http" instrumentation library as "request.latency".
	view := NewView(Instrument{
		Name: "latency",
		Scope: instrumentation.Scope{
			Name:    "http",
			Version: "v0.34.0",
		},
	}, Stream{Name: "request.latency"})

	// The created view can then be registered with the OpenTelemetry metric
	// SDK using the WithView option. Below is an example of how the view will
	// function in the SDK for certain instruments.

	stream, _ := view(Instrument{
		Name:        "latency",
		Description: "request latency",
		Unit:        "ms",
		Kind:        InstrumentKindCounter,
		Scope: instrumentation.Scope{
			Name:      "http",
			Version:   "v0.34.0",
			SchemaURL: "https://opentelemetry.io/schemas/1.0.0",
		},
	})
	fmt.Println("name:", stream.Name)
	fmt.Println("description:", stream.Description)
	fmt.Println("unit:", stream.Unit)
	// Output:
	// name: request.latency
	// description: request latency
	// unit: ms
}

func ExampleNewView_drop() {
	// Create a view that sets the drop aggregator for all instrumentation from
	// the "db" library, effectively turning-off all instrumentation from that
	// library.
	view := NewView(
		Instrument{Scope: instrumentation.Scope{Name: "db"}},
		Stream{Aggregation: aggregation.Drop{}},
	)

	// The created view can then be registered with the OpenTelemetry metric
	// SDK using the WithView option. Below is an example of how the view will
	// function in the SDK for certain instruments.

	stream, _ := view(Instrument{
		Name:  "queries",
		Kind:  InstrumentKindCounter,
		Scope: instrumentation.Scope{Name: "db", Version: "v0.4.0"},
	})
	fmt.Println("name:", stream.Name)
	fmt.Printf("aggregation: %#v", stream.Aggregation)
	// Output:
	// name: queries
	// aggregation: aggregation.Drop{}
}

func ExampleNewView_wildcard() {
	// Create a view that sets unit to milliseconds for any instrument with a
	// name suffix of ".ms".
	view := NewView(
		Instrument{Name: "*.ms"},
		Stream{Unit: "ms"},
	)

	// The created view can then be registered with the OpenTelemetry metric
	// SDK using the WithView option. Below is an example of how the view
	// function in the SDK for certain instruments.

	stream, _ := view(Instrument{
		Name: "computation.time.ms",
		Unit: "1",
	})
	fmt.Println("name:", stream.Name)
	fmt.Println("unit:", stream.Unit)
	// Output:
	// name: computation.time.ms
	// unit: ms
}

func ExampleView() {
	// The NewView function provides convenient creation of common Views
	// construction. However, it is limited in what it can create.
	//
	// When NewView is not able to provide the functionally needed, a custom
	// View can be constructed directly. Here a custom View is constructed that
	// uses Go's regular expression matching to ensure all data stream names
	// have a suffix of the units it uses.

	re := regexp.MustCompile(`[._](ms|byte)$`)
	var view View = func(i Instrument) (Stream, bool) {
		s := Stream{Name: i.Name, Description: i.Description, Unit: i.Unit}
		// Any instrument that does not have a unit suffix defined, but has a
		// dimensional unit defined, update the name with a unit suffix.
		if re.MatchString(i.Name) {
			return s, false
		}
		switch i.Unit {
		case "ms":
			s.Name += ".ms"
		case "By":
			s.Name += ".byte"
		default:
			return s, false
		}
		return s, true
	}

	// The created view can then be registered with the OpenTelemetry metric
	// SDK using the WithView option. Below is an example of how the view will
	// function in the SDK for certain instruments.

	stream, _ := view(Instrument{
		Name: "computation.time.ms",
		Unit: "ms",
	})
	fmt.Println("name:", stream.Name)

	stream, _ = view(Instrument{
		Name: "heap.size",
		Unit: "By",
	})
	fmt.Println("name:", stream.Name)
	// Output:
	// name: computation.time.ms
	// name: heap.size.byte
}
