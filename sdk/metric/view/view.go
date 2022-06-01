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

//go:build go1.17
// +build go1.17

package view // import "go.opentelemetry.io/otel/sdk/metric/view"

import (
	"fmt"
	"regexp"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/internal/instrument"
)

// View provides users with the flexibility to customize the metrics that are
// output by the SDK. A View can be used to:
//
// * Ignore Instruments.
// * Change the name of an Instrument.
// * Change the aggregation of an Instrument.
// * Customize which attribute(s) are to be reported by the Instrument.
//
// An empty config will match all instruments, and do no modifications.
type View struct {
	instrumentName       string
	instrumentNameRegexp *regexp.Regexp
	scope                instrumentation.Library

	filterKeys  map[attribute.Key]struct{}
	viewName    string
	description string
	// TODO: Aggregation selection
}

// New returns a new configured view Config.
func New(opts ...Option) (View, error) {
	v := View{}

	for _, opt := range opts {
		v = opt.apply(v)
	}

	if v.instrumentName == "" && v.viewName != "" {
		return View{}, fmt.Errorf("must use an instrument name with a view name")
	}

	return v, nil
}

// TransformInstrument will check if an instrument description matches this view
// and will convert it if it does.
func (v View) TransformInstrument(desc instrument.Description) (outputDesc instrument.Description, match bool) {
	if !v.matchName(desc.InstrumentName) ||
		!v.matchLibraryName(desc.Library.Name) ||
		!v.matchLibrarySchemaURL(desc.Library.SchemaURL) ||
		!v.matchLibraryVersion(desc.Library.Version) {
		return instrument.Description{}, false
	}
	if v.viewName != "" {
		desc.InstrumentName = v.viewName
	}
	if v.description != "" {
		desc.InstrumentDescription = v.description
	}
	return desc, true
}

// TransformAttributes filters an attribute set to the keys in the config. If no
// filter was provided the original set is returned.
func (v View) TransformAttributes(input attribute.Set) attribute.Set {
	if len(v.filterKeys) == 0 {
		return input
	}
	out, _ := input.Filter(attribute.Filter(func(kv attribute.KeyValue) bool {
		_, ok := v.filterKeys[kv.Key]
		return ok
	}))
	return out
}

// TODO: Provide Transform* for AggregationKind (#2816)

func (v View) matchName(name string) bool {
	if v.instrumentNameRegexp != nil {
		return v.instrumentNameRegexp.MatchString(name)
	}
	return v.instrumentName == "" || name == v.instrumentName
}

func (v View) matchLibraryName(name string) bool {
	return v.scope.Name == "" || name == v.scope.Name
}

func (v View) matchLibraryVersion(version string) bool {
	return v.scope.Version == "" || version == v.scope.Version
}

func (v View) matchLibrarySchemaURL(schemaURL string) bool {
	return v.scope.SchemaURL == "" || schemaURL == v.scope.SchemaURL
}

// Option applies a Configuration option value to a view Config. All options
// will be used together to determine match and transforms.
type Option interface {
	apply(View) View
}

type optionFunc func(View) View

func (f optionFunc) apply(v View) View {
	return f(v)
}

// MatchInstrumentName will do an exact match of the name of the instrument.
// Not compatible with MatchInstrumentNameRegexp.
func MatchInstrumentName(name string) Option {
	return optionFunc(func(v View) View {
		v.instrumentName = name
		v.instrumentNameRegexp = nil
		return v
	})
}

// MatchInstrumentNameRegexp will match any instrument with the provided
// regexp.
// Not compatible with MatchInstrumentName or WithName.
func MatchInstrumentNameRegexp(re *regexp.Regexp) Option {
	return optionFunc(func(v View) View {
		v.instrumentNameRegexp = re
		v.instrumentName = ""
		return v
	})
}

// TODO: Implement when InstrumentKind and NumberKind are defined
// func MatchInstrumentKind(k sdkinstrument.Kind) Option {
// 	return optionFunc(func(v Config) Config {
// 		v.instrumentKind = k
// 		return v
// 	})
// }

// func MatchNumberKind(k number.Kind) Option {
// 	return optionFunc(func(v Config) Config {
// 		v.numberKind = k
// 		return v
// 	})
// }

// MatchInstrumentationLibrary will do an exact match on any
// instrumentation.Library field that is not blank ("").
func MatchInstrumentationLibrary(lib instrumentation.Library) Option {
	return optionFunc(func(v View) View {
		v.scope = lib
		return v
	})
}

// WithName will rename the instrument the view matches. If not used or empty the
// instrument name will not be changed. Not compatible with
// MatchInstrumentNameRegexp.
func WithName(name string) Option {
	return optionFunc(func(v View) View {
		v.viewName = name
		return v
	})
}

// WithDescription will change the description of the instruments the view
// matches. If not used or empty the description will not be changed.
func WithDescription(desc string) Option {
	return optionFunc(func(v View) View {
		v.description = desc
		return v
	})
}

// WithAttributeFilter will select attributes that have a matching key.  If not used
// or empty no filter will be applied.
func WithAttributeFilter(keys ...attribute.Key) Option {
	return optionFunc(func(v View) View {
		if v.filterKeys == nil {
			v.filterKeys = map[attribute.Key]struct{}{}
		}
		for _, key := range keys {
			v.filterKeys[key] = struct{}{}
		}
		return v
	})
}

// TODO: Implement when Aggregations and Temporalities are defined
// func WithAggregation(kind aggregation.Kind) Option {
// 	return optionFunc(func(v Config) Config {
// 		v.aggregation = kind
// 		return v
// 	})
// }
