package metric

import (
	"sync"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/internal/views"
	"go.opentelemetry.io/otel/sdk/resource"
)

type (
	Config struct {
		res   *resource.Resource
		views []views.View
	}

	Option func(cfg *Config)

	provider struct {
		cfg Config

		lock   sync.Mutex
		meters map[instrumentation.Library]*meter
	}

	meter struct {
	}
)

func WithResource(res *resource.Resource) Option {
	return func(cfg *Config) {
		cfg.res = res
	}
}

func WithView(view views.View) Option {
	return func(cfg *Config) {
		cfg.views = append(cfg.views, view)
	}
}

func New(opts ...Option) metric.MeterProvider {
	cfg := Config{
		res: resource.Default(),
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	return &provider{
		cfg:    cfg,
		meters: map[instrumentation.Library]*meter{},
	}
}

func (p *provider) Meter(name string, opts ...metric.MeterOption) metric.Meter {
	cfg := metric.NewMeterConfig(opts...)
	lib := instrumentation.Library{
		Name:      name,
		Version:   cfg.Version(),
		SchemaURL: cfg.SchemaURL(),
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	if m, ok := p.meters[lib]; ok {
		return m
	}

}
