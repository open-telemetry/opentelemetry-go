package periodic

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric/reader"
)

type exporter struct {
	ticker  *time.Ticker
	timeout time.Duration

	reader reader.ManualReader

	done chan struct{}
}

type config struct {
	timeout time.Duration
}

func newConfig(opts ...Option) config {
	cfg := config{
		timeout: time.Minute,
	}

	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}

	return cfg
}

type Option interface {
	apply(config) config
}

type optionFunc func(config) config

func (o optionFunc) apply(cfg config) config {
	return o(cfg)
}

func WithTimeout(to time.Duration) Option {
	return optionFunc(func(cfg config) config {
		cfg.timeout = to
		return cfg
	})
}

var _ reader.Reader = &exporter{}

func New(period time.Duration, exp reader.Exporter, opts ...Option) reader.Reader {
	cfg := newConfig(opts...)

	e := &exporter{

		ticker:  time.NewTicker(period),
		timeout: cfg.timeout,

		reader: *reader.NewManualReader(exp),

		done: make(chan struct{}),
	}

	go func() {
		for {
			select {
			case <-e.ticker.C:
				e.collect(context.Background())
			case <-e.done:
				e.ticker.Stop()
				return
			}
		}
	}()

	return e
}

func (e *exporter) Register(prod reader.Producer) {
	e.reader.Register(prod)
}

func (e *exporter) collect(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	err := e.reader.Collect(ctx, nil)
	if err != nil {
		otel.Handle(err)
	}
}

func (e *exporter) Flush(ctx context.Context) error {
	return e.reader.Flush(ctx)
}

func (e *exporter) Shutdown(ctx context.Context) error {
	close(e.done)
	e.reader.Shutdown(ctx)

	return nil
}
