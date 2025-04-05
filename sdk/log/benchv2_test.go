// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/log"

	"github.com/stretchr/testify/assert"
)

func BenchmarkV2Processor(b *testing.B) {
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
				return []LoggerProviderOption{WithProcessor(NewBatchProcessorV2(noopExporter{}))}
			},
		},
		{
			name: "BatchSimulateExport",
			f: func() []LoggerProviderOption {
				return []LoggerProviderOption{WithProcessor(NewBatchProcessorV2(mockDelayExporter{}))}
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
					WithProcessor(NewBatchProcessorV2(noopExporter{})),
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
					WithProcessor(NewBatchProcessorV2(noopExporter{})),
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
					WithProcessor(NewBatchProcessorV2(noopExporter{})),
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
