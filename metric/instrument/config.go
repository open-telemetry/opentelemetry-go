package instrument

import "go.opentelemetry.io/otel/metric/unit"

// Config contains options for metric instrument descriptors.
type Config struct {
	description string
	unit        unit.Unit
}

// Description describes the instrument in human-readable terms.
func (cfg Config) Description() string {
	return cfg.description
}

// Unit describes the measurement unit for a instrument.
func (cfg Config) Unit() unit.Unit {
	return cfg.unit
}

// Option is an interface for applying metric instrument options.
type Option interface {
	// ApplyMeter is used to set a Option value of a
	// Config.
	applyInstrument(*Config)
}

// NewConfig creates a new Config
// and applies all the given options.
func NewConfig(opts ...Option) Config {
	var config Config
	for _, o := range opts {
		o.applyInstrument(&config)
	}
	return config
}

type optionFunc func(*Config)

func (fn optionFunc) applyInstrument(cfg *Config) {
	fn(cfg)
}

// WithDescription applies provided description.
func WithDescription(desc string) Option {
	return optionFunc(func(cfg *Config) {
		cfg.description = desc
	})
}

// WithUnit applies provided unit.
func WithUnit(unit unit.Unit) Option {
	return optionFunc(func(cfg *Config) {
		cfg.unit = unit
	})
}
