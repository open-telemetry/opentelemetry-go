package reader

import (
	"unsafe"
	
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
	"go.opentelemetry.io/otel/sdk/metric/views"
)

type (
	Config struct {
		views          []views.View
		hasDefaultView bool
	}

	Option func(config *Config)

	Reader struct {
		config   Config
		storage  map[mapkey]unsafe.Pointer
		tempSort attribute.Sortable
	}

	mapkey struct {
		*sdkapi.Descriptor
		attribute.Set
	}
)

func WithView(view views.View) Option {
	return func(cfg *Config) {
		cfg.views = append(cfg.views, view)
	}
}

func WithDefaultView(hasDefaultView bool) Option {
	return func(cfg *Config) {
		cfg.hasDefaultView = hasDefaultView
	}
}

func NewConfig(opts ...Option) Config {
	cfg := Config{
		hasDefaultView: true,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

func New(config Config) *Reader {
	return &Reader{
		config: config,
		storage: map[mapkey]unsafe.Pointer{},
	}
}

func (r *Reader) Views() []views.View {
	return r.config.views
}

func (r *Reader) HasDefaultView() bool {
	return r.config.hasDefaultView
}

func FindOrCreate[N number.Any, Storage, Config any, Methods aggregator.Methods[N, Storage, Config]](
	reader *Reader,
	desc *sdkapi.Descriptor,
	attrs attribute.Set,
	aggConfig *Config,
) *Storage {
	var methods Methods
	mk := mapkey{
		Descriptor: desc,
		Set: attrs,
	}
	uptr, has := reader.storage[mk]

	if has {
		return (*Storage)(uptr)
	}
	ns := new(Storage)
	methods.Init(ns, *aggConfig)
	reader.storage[mk] = unsafe.Pointer(ns)
	return ns
}
