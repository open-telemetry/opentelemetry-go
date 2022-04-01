package views

import (
	"regexp"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregation/histogram"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/sdkinstrument"
)

type (
	View struct {
		cfg Config
	}

	keyFilter map[attribute.Key]struct{}

	Config struct {
		// Matchers for the instrument
		instrumentName       string
		instrumentNameRegexp *regexp.Regexp
		instrumentKind       sdkinstrument.Kind
		numberKind           number.Kind
		library              instrumentation.Library

		// Properties of the view
		keys         keyFilter
		name         string
		description  string
		aggregation  aggregation.Kind
		temporality  aggregation.Temporality
		histoOptions []histogram.Option
	}

	Option func(cfg *Config)
)

const (
	unsetKind       = sdkinstrument.Kind(-1)
	unsetNumberKind = number.Kind(-1)
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

func MatchKind(k sdkinstrument.Kind) Option {
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

func WithKeys(keys []attribute.Key) Option {
	return func(cfg *Config) {
		if len(keys) == 0 {
			cfg.keys = nil
			return
		}
		if cfg.keys == nil {
			cfg.keys = keyFilter{}
		}
		for _, key := range keys {
			cfg.keys[key] = struct{}{}
		}
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

func WithTemporality(tempo aggregation.Temporality) Option {
	return func(cfg *Config) {
		cfg.temporality = tempo
	}
}

func New(opts ...Option) View {
	cfg := Config{
		instrumentKind: unsetKind,
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

func (ks keyFilter) filter(kv attribute.KeyValue) bool {
	_, has := ks[kv.Key]
	return has
}

func (v View) Keys() attribute.Filter {
	return v.cfg.keys.filter
}

func (v View) Description() string {
	return v.cfg.description
}

func (v View) Aggregation() aggregation.Kind {
	return v.cfg.aggregation
}

func (v View) Temporality() aggregation.Temporality {
	return v.cfg.temporality
}

func (v View) HistogramOptions() []histogram.Option {
	return v.cfg.histoOptions
}

func stringMismatch(test, value string) bool {
	return test != "" && test != value
}

func ikindMismatch(test, value sdkinstrument.Kind) bool {
	return test != unsetKind && test != value
}

func nkindMismatch(test, value number.Kind) bool {
	return test != unsetNumberKind && test != value
}

func regexpMismatch(test *regexp.Regexp, value string) bool {
	return test != nil && test.MatchString(value)
}

func (v View) Matches(lib instrumentation.Library, desc sdkinstrument.Descriptor) bool {
	return !stringMismatch(v.cfg.library.Name, lib.Name) &&
		!stringMismatch(v.cfg.library.Version, lib.Version) &&
		!stringMismatch(v.cfg.library.SchemaURL, lib.SchemaURL) &&
		!stringMismatch(v.cfg.instrumentName, desc.Name) &&
		!ikindMismatch(v.cfg.instrumentKind, desc.Kind) &&
		!nkindMismatch(v.cfg.numberKind, desc.NumberKind) &&
		!regexpMismatch(v.cfg.instrumentNameRegexp, desc.Name)
}
