// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/log"
)

type mockDelayExporter struct{}

func (mockDelayExporter) Export(context.Context, []Record) error {
	time.Sleep(time.Millisecond * 5)
	return nil
}

func (mockDelayExporter) Shutdown(context.Context) error { return nil }

func (mockDelayExporter) ForceFlush(context.Context) error { return nil }

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
			name: "BatchSimulateExport",
			f: func() []LoggerProviderOption {
				return []LoggerProviderOption{WithProcessor(NewBatchProcessor(mockDelayExporter{}))}
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

func (timestampProcessor) OnEmit(_ context.Context, r *Record) error {
	r.SetObservedTimestamp(time.Date(1988, time.November, 17, 0, 0, 0, 0, time.UTC))
	return nil
}

func (timestampProcessor) Enabled(context.Context, Record) bool {
	return true
}

func (timestampProcessor) Shutdown(context.Context) error {
	return nil
}

func (timestampProcessor) ForceFlush(context.Context) error {
	return nil
}

type attrAddProcessor struct{}

func (attrAddProcessor) OnEmit(_ context.Context, r *Record) error {
	r.AddAttributes(log.String("add", "me"))
	return nil
}

func (attrAddProcessor) Enabled(context.Context, Record) bool {
	return true
}

func (attrAddProcessor) Shutdown(context.Context) error {
	return nil
}

func (attrAddProcessor) ForceFlush(context.Context) error {
	return nil
}

type attrSetDecorator struct{}

func (attrSetDecorator) OnEmit(_ context.Context, r *Record) error {
	r.SetAttributes(log.String("replace", "me"))
	return nil
}

func (attrSetDecorator) Enabled(context.Context, Record) bool {
	return true
}

func (attrSetDecorator) Shutdown(context.Context) error {
	return nil
}

func (attrSetDecorator) ForceFlush(context.Context) error {
	return nil
}
