package view

import (
	"regexp"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/sdkinstrument"
)

const (
	unsetInstrumentKind = sdkinstrument.Kind(-1)
	unsetNumberKind     = number.Kind(-1)
)

// ClauseOption applies a configuration option value to a view Config.
type ClauseOption interface {
	apply(ClauseConfig) ClauseConfig
}

// clauseOptionFunction makes a functional ClauseOption out of a function object.
type clauseOptionFunction func(cfg ClauseConfig) ClauseConfig

// apply implements ClauseOption.
func (of clauseOptionFunction) apply(in ClauseConfig) ClauseConfig {
	return of(in)
}

// Matchers

func MatchInstrumentName(name string) ClauseOption {
	return clauseOptionFunction(func(clause ClauseConfig) ClauseConfig {
		clause.instrumentName = name
		return clause
	})
}

func MatchInstrumentNameRegexp(re *regexp.Regexp) ClauseOption {
	return clauseOptionFunction(func(clause ClauseConfig) ClauseConfig {
		clause.instrumentNameRegexp = re
		return clause
	})
}

func MatchInstrumentKind(k sdkinstrument.Kind) ClauseOption {
	return clauseOptionFunction(func(clause ClauseConfig) ClauseConfig {
		clause.instrumentKind = k
		return clause
	})
}

func MatchNumberKind(k number.Kind) ClauseOption {
	return clauseOptionFunction(func(clause ClauseConfig) ClauseConfig {
		clause.numberKind = k
		return clause
	})
}

func MatchInstrumentationLibrary(lib instrumentation.Library) ClauseOption {
	return clauseOptionFunction(func(clause ClauseConfig) ClauseConfig {
		clause.library = lib
		return clause
	})
}

// Properties

// WithKeys overwrites; nil is distinct from empty non-nil.
func WithKeys(keys []attribute.Key) ClauseOption {
	return clauseOptionFunction(func(clause ClauseConfig) ClauseConfig {
		clause.keys = keys
		return clause
	})
}

func WithName(name string) ClauseOption {
	return clauseOptionFunction(func(clause ClauseConfig) ClauseConfig {
		clause.name = name
		return clause
	})
}

func WithDescription(desc string) ClauseOption {
	return clauseOptionFunction(func(clause ClauseConfig) ClauseConfig {
		clause.description = desc
		return clause
	})
}

func WithAggregation(kind aggregation.Kind) ClauseOption {
	return clauseOptionFunction(func(clause ClauseConfig) ClauseConfig {
		clause.aggregation = kind
		return clause
	})
}

func WithTemporality(tempo aggregation.Temporality) ClauseOption {
	return clauseOptionFunction(func(clause ClauseConfig) ClauseConfig {
		clause.temporality = tempo
		return clause
	})
}

func WithAggregatorConfig(acfg aggregator.Config) ClauseOption {
	return clauseOptionFunction(func(clause ClauseConfig) ClauseConfig {
		clause.acfg = acfg
		return clause
	})
}

// IsSingleInstrument is a requirement when HasName().
func (c *ClauseConfig) IsSingleInstrument() bool {
	return c.instrumentName != ""
}

// HasName implies IsSingleInstrument SHOULD be required.
func (c *ClauseConfig) HasName() bool {
	return c.name != ""
}

func (c *ClauseConfig) Name() string {
	return c.name
}

func (c *ClauseConfig) Keys() []attribute.Key {
	return c.keys
}

func (c *ClauseConfig) Description() string {
	return c.description
}

func (c *ClauseConfig) Aggregation() aggregation.Kind {
	return c.aggregation
}

func (c *ClauseConfig) Temporality() aggregation.Temporality {
	return c.temporality
}

func (c *ClauseConfig) AggregatorConfig() aggregator.Config {
	return c.acfg
}

func stringMismatch(test, value string) bool {
	return test != "" && test != value
}

func ikindMismatch(test, value sdkinstrument.Kind) bool {
	return test != unsetInstrumentKind && test != value
}

func nkindMismatch(test, value number.Kind) bool {
	return test != unsetNumberKind && test != value
}

func regexpMismatch(test *regexp.Regexp, value string) bool {
	return test != nil && test.MatchString(value)
}

func (c *ClauseConfig) Matches(lib instrumentation.Library, desc sdkinstrument.Descriptor) bool {
	return !stringMismatch(c.library.Name, lib.Name) &&
		!stringMismatch(c.library.Version, lib.Version) &&
		!stringMismatch(c.library.SchemaURL, lib.SchemaURL) &&
		!stringMismatch(c.instrumentName, desc.Name) &&
		!ikindMismatch(c.instrumentKind, desc.Kind) &&
		!nkindMismatch(c.numberKind, desc.NumberKind) &&
		!regexpMismatch(c.instrumentNameRegexp, desc.Name)
}
