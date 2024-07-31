// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"
	"testing"
	"time"

	"go.opentelemetry.io/otel/log"

	"github.com/stretchr/testify/assert"
)

func BenchmarkProcessor(b *testing.B) {
	for _, tc := range []struct {
		name string
		f    func() []LoggerProviderOption
	}{
		{
			name: "Simple",
			f: func() []LoggerProviderOption {
				return []LoggerProviderOption{WithProcessor(NewSimpleProcessor(noopExporter{}))}
			},
		},
		{
			name: "Batch",
			f: func() []LoggerProviderOption {
				return []LoggerProviderOption{WithProcessor(NewBatchProcessor(noopExporter{}))}
			},
		},
		{
			name: "SetTimestampSimple",
			f: func() []LoggerProviderOption {
				return []LoggerProviderOption{
					WithProcessor(timestampProcessor{}),
					WithProcessor(NewSimpleProcessor(noopExporter{})),
				}
			},
		},
		{
			name: "SetTimestampBatch",
			f: func() []LoggerProviderOption {
				return []LoggerProviderOption{
					WithProcessor(timestampProcessor{}),
					WithProcessor(NewBatchProcessor(noopExporter{})),
				}
			},
		},
		{
			name: "AddAttributesSimple",
			f: func() []LoggerProviderOption {
				return []LoggerProviderOption{
					WithProcessor(attrAddProcessor{}),
					WithProcessor(NewSimpleProcessor(noopExporter{})),
				}
			},
		},
		{
			name: "AddAttributesBatch",
			f: func() []LoggerProviderOption {
				return []LoggerProviderOption{
					WithProcessor(attrAddProcessor{}),
					WithProcessor(NewBatchProcessor(noopExporter{})),
				}
			},
		},
		{
			name: "SetAttributesSimple",
			f: func() []LoggerProviderOption {
				return []LoggerProviderOption{
					WithProcessor(attrSetDecorator{}),
					WithProcessor(NewSimpleProcessor(noopExporter{})),
				}
			},
		},
		{
			name: "SetAttributesBatch",
			f: func() []LoggerProviderOption {
				return []LoggerProviderOption{
					WithProcessor(attrSetDecorator{}),
					WithProcessor(NewBatchProcessor(noopExporter{})),
				}
			},
		},
	} {
		b.Run(tc.name, func(b *testing.B) {
			provider := NewLoggerProvider(tc.f()...)
			b.Cleanup(func() { assert.NoError(b, provider.Shutdown(context.Background())) })
			logger := provider.Logger(b.Name())

			b.ReportAllocs()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					r := log.Record{}
					r.SetBody(log.StringValue("message"))
					r.SetSeverity(log.SeverityInfo)
					r.AddAttributes(
						log.String("foo", "bar"),
						log.Float64("float", 3.14),
						log.Int("int", 123),
						log.Bool("bool", true),
					)
					logger.Emit(context.Background(), r)
				}
			})
		})
	}
}

type timestampProcessor struct{}

func (p timestampProcessor) OnEmit(ctx context.Context, r *Record) error {
	r.SetObservedTimestamp(time.Date(1988, time.November, 17, 0, 0, 0, 0, time.UTC))
	return nil
}

func (p timestampProcessor) Enabled(context.Context, Record) bool {
	return true
}

func (p timestampProcessor) Shutdown(ctx context.Context) error {
	return nil
}

func (p timestampProcessor) ForceFlush(ctx context.Context) error {
	return nil
}

type attrAddProcessor struct{}

func (p attrAddProcessor) OnEmit(ctx context.Context, r *Record) error {
	r.AddAttributes(log.String("add", "me"))
	return nil
}

func (p attrAddProcessor) Enabled(context.Context, Record) bool {
	return true
}

func (p attrAddProcessor) Shutdown(ctx context.Context) error {
	return nil
}

func (p attrAddProcessor) ForceFlush(ctx context.Context) error {
	return nil
}

type attrSetDecorator struct{}

func (p attrSetDecorator) OnEmit(ctx context.Context, r *Record) error {
	r.SetAttributes(log.String("replace", "me"))
	return nil
}

func (p attrSetDecorator) Enabled(context.Context, Record) bool {
	return true
}

func (p attrSetDecorator) Shutdown(ctx context.Context) error {
	return nil
}

func (p attrSetDecorator) ForceFlush(ctx context.Context) error {
	return nil
}
