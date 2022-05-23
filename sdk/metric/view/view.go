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

// Config contains Configuration options for a view. An empty config will match
// all instruments, and do no modifications.
type Config struct {
	instrumentName       string
	instrumentNameRegexp *regexp.Regexp
	library              instrumentation.Library

	filterKeys  map[attribute.Key]struct{}
	viewName    string
	description string
	// TODO: Aggregation selection
}

// New returns a new configured view Config.
func New(opts ...Option) (Config, error) {
	cfg := Config{}

	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}

	if cfg.instrumentName != "" && cfg.instrumentNameRegexp != nil {
		return Config{}, fmt.Errorf("can not use name selector and name regex in the same view")
	}
	if cfg.instrumentNameRegexp != nil && cfg.viewName != "" {
		return Config{}, fmt.Errorf("can not use name regex with a view name")
	}

	return cfg, nil
}

// TransformInstrument will check if an instrument description matches this view
// and will convert it if it does.
func (cfg Config) TransformInstrument(desc instrument.Description) (outputDesc instrument.Description, match bool) {
	if !cfg.matchName(desc.InstrumentName) ||
		!cfg.matchLibraryName(desc.Library.Name) ||
		!cfg.matchLibrarySchemaURL(desc.Library.SchemaURL) ||
		!cfg.matchLibraryVersion(desc.Library.Version) {
		return instrument.EmptyDescription, false
	}
	if cfg.viewName != "" {
		desc.InstrumentName = cfg.viewName
	}
	if cfg.description != "" {
		desc.InstrumentDescription = cfg.description
	}
	return desc, true
}

// TransformAttributes filters an attribute set to the keys in the config. If no
// filter was provided the original set is returned.
func (cfg Config) TransformAttributes(input attribute.Set) attribute.Set {
	if len(cfg.filterKeys) == 0 {
		return input
	}
	out, _ := input.Filter(attribute.Filter(func(kv attribute.KeyValue) bool {
		_, ok := cfg.filterKeys[kv.Key]
		return ok
	}))
	return out
}

// TODO: Provide Transfrom* for AggregationKind (#2816)

func (cfg Config) matchName(name string) bool {
	if cfg.instrumentNameRegexp != nil {
		return cfg.instrumentNameRegexp.MatchString(name)
	}
	return cfg.instrumentName == "" || name == cfg.instrumentName
}
func (cfg Config) matchLibraryName(name string) bool {
	return cfg.library.Name == "" || name == cfg.library.Name
}
func (cfg Config) matchLibraryVersion(version string) bool {
	return cfg.library.Version == "" || version == cfg.library.Version
}
func (cfg Config) matchLibrarySchemaURL(schemaURL string) bool {
	return cfg.library.SchemaURL == "" || schemaURL == cfg.library.SchemaURL
}

// Option applies a Configuration option value to a view Config. All options
// will be used together to determine match and transforms.
type Option interface {
	apply(Config) Config
}

type optionFunc func(Config) Config

func (f optionFunc) apply(cfg Config) Config {
	return f(cfg)
}

// MatchInstrumentName will do an exact match of the name of the instrument.
// Not compatible with MatchInstrumentNameRegexp.
func MatchInstrumentName(name string) Option {
	return optionFunc(func(cfg Config) Config {
		cfg.instrumentName = name
		return cfg
	})
}

// Match MatchInstrumentNameRegexp will match any instrument with the provided
// regexp.
// Not compatible with MatchInstrumentName or WithName
func MatchInstrumentNameRegexp(re *regexp.Regexp) Option {
	return optionFunc(func(cfg Config) Config {
		cfg.instrumentNameRegexp = re
		return cfg
	})
}

// TODO: Implement when InstrumentKind and NumberKind are defined
// func MatchInstrumentKind(k sdkinstrument.Kind) Option {
// 	return optionFunc(func(cfg Config) Config {
// 		cfg.instrumentKind = k
// 		return cfg
// 	})
// }

// func MatchNumberKind(k number.Kind) Option {
// 	return optionFunc(func(cfg Config) Config {
// 		cfg.numberKind = k
// 		return cfg
// 	})
// }

// MatchInstrumentationLibrary will do an exact match on any
// instrumentation.Library field that is not blank ("").
func MatchInstrumentationLibrary(lib instrumentation.Library) Option {
	return optionFunc(func(cfg Config) Config {
		cfg.library = lib
		return cfg
	})
}

// WithName will rename the instrument the view matches. If not used or empty the
// instrument name will not be changed. Not compatible with
// MatchInstrumentNameRegexp.
func WithName(name string) Option {
	return optionFunc(func(cfg Config) Config {
		cfg.viewName = name
		return cfg
	})
}

// WithDescription will change the description of the instruments the view
// matches. If not used or empty the description will not be changed.
func WithDescription(desc string) Option {
	return optionFunc(func(cfg Config) Config {
		cfg.description = desc
		return cfg
	})
}

// WithKeyFilter will select attributes that have a matching key.  If not used
// or empty no filter will be applied.
func WithKeyFilter(keys ...attribute.Key) Option {
	return optionFunc(func(cfg Config) Config {
		if cfg.filterKeys == nil {
			cfg.filterKeys = map[attribute.Key]struct{}{}
		}
		for _, key := range keys {
			cfg.filterKeys[key] = struct{}{}
		}
		return cfg
	})
}

// TODO: Implement when Aggregations and Temporalities are defined
// func WithAggregation(kind aggregation.Kind) Option {
// 	return optionFunc(func(cfg Config) Config {
// 		cfg.aggregation = kind
// 		return cfg
// 	})
// }

// func WithTemporality(tempo aggregation.Temporality) Option {
// 	return optionFunc(func(cfg Config) Config {
// 		cfg.temporality = tempo
// 		return cfg
// 	})
// }
