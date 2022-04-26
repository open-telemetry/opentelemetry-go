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

// Matchers

func MatchInstrumentName(name string) ClauseOption {
	return func(clause *ClauseConfig) {
		clause.instrumentName = name
	}
}

func MatchInstrumentNameRegexp(re *regexp.Regexp) ClauseOption {
	return func(clause *ClauseConfig) {
		clause.instrumentNameRegexp = re
	}
}

func MatchInstrumentKind(k sdkinstrument.Kind) ClauseOption {
	return func(clause *ClauseConfig) {
		clause.instrumentKind = k
	}
}

func MatchNumberKind(k number.Kind) ClauseOption {
	return func(clause *ClauseConfig) {
		clause.numberKind = k
	}
}

func MatchInstrumentationLibrary(lib instrumentation.Library) ClauseOption {
	return func(clause *ClauseConfig) {
		clause.library = lib
	}
}

// Properties

// WithKeys overwrites; nil is distinct from empty non-nil.
func WithKeys(keys []attribute.Key) ClauseOption {
	return func(clause *ClauseConfig) {
		clause.keys = keys
	}
}

func WithName(name string) ClauseOption {
	return func(clause *ClauseConfig) {
		clause.name = name
	}
}

func WithDescription(desc string) ClauseOption {
	return func(clause *ClauseConfig) {
		clause.description = desc
	}
}

func WithAggregation(kind aggregation.Kind) ClauseOption {
	return func(clause *ClauseConfig) {
		clause.aggregation = kind
	}
}

func WithTemporality(tempo aggregation.Temporality) ClauseOption {
	return func(clause *ClauseConfig) {
		clause.temporality = tempo
	}
}

func WithAggregatorConfig(acfg aggregator.Config) ClauseOption {
	return func(clause *ClauseConfig) {
		clause.acfg = acfg
	}
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
