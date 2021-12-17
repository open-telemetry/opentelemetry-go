package views

import (
	"regexp"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/sdkapi"
	"go.opentelemetry.io/otel/metric/sdkapi/number"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
)

type (
	View struct {
		cfg Config
	}

	Config struct {
		// Matchers for the instrument
		instrumentName       string
		instrumentNameRegexp *regexp.Regexp
		instrumentKind       sdkapi.InstrumentKind
		numberKind           number.Kind
		library              instrumentation.Library

		// Properties of the view
		attributes  []attribute.Key
		name        string
		description string
		aggregation aggregation.Kind
	}

	Option func(cfg *Config)
)

const (
	unsetInstrumentKind = sdkapi.InstrumentKind(-1)
	unsetNumberKind     = number.Kind(-1)
)

// Matchers

func MatchInstrumentName(name string) Option {
	return func(cfg *Config) {
		cfg.instrumentName = name
	}
}

func MatchInstrumentNameRegexp(re *regexp.Regexp) Option {
	return func(cfg *Config) {
		cfg.instrumentNameRegexp = re
	}
}

func MatchInstrumentKind(k sdkapi.InstrumentKind) Option {
	return func(cfg *Config) {
		cfg.instrumentKind = k
	}
}

func MatchNumberKind(k number.Kind) Option {
	return func(cfg *Config) {
		cfg.numberKind = k
	}
}

func MatchInstrumentationLibrary(lib instrumentation.Library) Option {
	return func(cfg *Config) {
		cfg.library = lib
	}
}

// Properties

func WithAttributes(keys ...attribute.Key) Option {
	return func(cfg *Config) {
		cfg.attributes = keys
	}
}

func WithName(name string) Option {
	return func(cfg *Config) {
		cfg.name = name
	}
}

func WithDescription(desc string) Option {
	return func(cfg *Config) {
		cfg.description = desc
	}
}

func WithAggregation(kind aggregation.Kind) Option {
	return func(cfg *Config) {
		cfg.aggregation = kind
	}
}

func New(opts ...Option) View {
	cfg := Config{
		instrumentKind: unsetInstrumentKind,
		numberKind:     unsetNumberKind,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return View{
		cfg: cfg,
	}
}

// IsSingleInstrument is a requirement when HasName().
func (v View) IsSingleInstrument() bool {
	return v.cfg.instrumentName != ""
}

// HasName implies IsSingleInstrument SHOULD be required.
func (v View) HasName() bool {
	return v.cfg.name != ""
}

func (v View) Name() string {
	return v.cfg.name
}

func (v View) Attributes() []attribute.Key {
	return v.cfg.attributes
}

func (v View) Description() string {
	return v.cfg.description
}

func (v View) Aggregation() aggregation.Kind {
	return v.cfg.aggregation
}

func stringMismatch(test, value string) bool {
	return test != "" && test != value
}

func ikindMismatch(test, value sdkapi.InstrumentKind) bool {
	return test != unsetInstrumentKind && test != value
}

func nkindMismatch(test, value number.Kind) bool {
	return test != unsetNumberKind && test != value
}

func regexpMismatch(test *regexp.Regexp, value string) bool {
	return test != nil && test.MatchString(value)
}

func (v View) Matches(lib instrumentation.Library, desc sdkapi.Descriptor) bool {
	return !stringMismatch(v.cfg.library.Name, lib.Name) &&
		!stringMismatch(v.cfg.library.Version, lib.Version) &&
		!stringMismatch(v.cfg.library.SchemaURL, lib.SchemaURL) &&
		!stringMismatch(v.cfg.instrumentName, desc.Name()) &&
		!ikindMismatch(v.cfg.instrumentKind, desc.InstrumentKind()) &&
		!nkindMismatch(v.cfg.numberKind, desc.NumberKind()) &&
		!regexpMismatch(v.cfg.instrumentNameRegexp, desc.Name())
}
