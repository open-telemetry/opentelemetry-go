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
	"regexp"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/internal/instrument"
)

// Config contains Configuration options for a view.
type Config struct {
	// TODO (#2837): implement.
	instrumentName       string
	instrumentNameRegexp *regexp.Regexp
	library              instrumentation.Library

	filterKeys []attribute.Key
	// additionalAttributes []attribute.KeyValue
	viewName    string
	description string
}

func (cfg Config) Transform(desc instrument.Description) instrument.Description {
	if !cfg.matchName(desc.InstrumentName) ||
		!cfg.matchLibraryName(desc.Library.Name) ||
		!cfg.matchLibrarySchemaURL(desc.Library.SchemaURL) ||
		!cfg.matchLibraryVersion(desc.Library.Version) {
		return instrument.EmptyDescription
	}
	if cfg.viewName != "" {
		desc.InstrumentName = cfg.viewName
	}
	if cfg.description != "" {
		desc.InstrumentDescription = cfg.description
	}
	return desc
}

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

// TODO (#2837): add getter functions for all the internal fields of a Config.

// Option applies a Configuration option value to a view Config.
type Option interface {
	apply(Config) Config
}

type optionFunc func(Config) Config

func (f optionFunc) apply(cfg Config) Config {
	return f(cfg)
}

// TODO (#2837): implement view match options.
// TODO (#2837): implement view annotation options.

// New returns a new and Configured view Config.
func New(opts ...Option) (Config, error) {
	cfg := Config{}

	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}
	// TODO invariant checking eg. can't rename with regex

	return cfg, nil
}

func MatchInstrumentName(name string) Option {
	return optionFunc(func(cfg Config) Config {
		cfg.instrumentName = name
		return cfg
	})
}

func MatchInstrumentNameRegexp(re *regexp.Regexp) Option {
	return optionFunc(func(cfg Config) Config {
		cfg.instrumentNameRegexp = re
		return cfg
	})
}

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

func MatchInstrumentationLibrary(lib instrumentation.Library) Option {
	return optionFunc(func(cfg Config) Config {
		cfg.library = lib
		return cfg
	})
}

// Properties

func WithKeyFilter(keys ...attribute.Key) Option {
	return optionFunc(func(cfg Config) Config {
		cfg.filterKeys = append(cfg.filterKeys, keys...)
		return cfg
	})
}

// func WithAdditionalAttributes(attrs ...attribute.KeyValue) Option {
// 	return optionFunc(func(cfg Config) Config {
// 		cfg.additionalAttributes = append(cfg.additionalAttributes, attrs...)
// 		return cfg
// 	})
// }

func WithName(name string) Option {
	return optionFunc(func(cfg Config) Config {
		cfg.viewName = name
		return cfg
	})
}

func WithDescription(desc string) Option {
	return optionFunc(func(cfg Config) Config {
		cfg.description = desc
		return cfg
	})
}

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

// func WithAggregatorConfig(acfg aggregator.Config) Option {
// 	return optionFunc(func(cfg Config) Config {
// 		cfg.acfg = acfg
// 		return cfg
// 	})
// }
