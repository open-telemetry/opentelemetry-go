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
		f    func() []Processor
	}{
		{
			name: "Simple",
			f: func() []Processor {
				return []Processor{
					NewSimpleProcessor(noopExporter{}),
				}
			},
		},
		{
			name: "Batch",
			f: func() []Processor {
				return []Processor{
					NewBatchProcessor(noopExporter{}),
				}
			},
		},
		{
			name: "ModifyTimestampSimple",
			f: func() []Processor {
				return []Processor{
					timestampSetter{},
					NewSimpleProcessor(noopExporter{}),
				}
			},
		},
		{
			name: "ModifyTimestampBatch",
			f: func() []Processor {
				return []Processor{
					timestampSetter{},
					NewBatchProcessor(noopExporter{}),
				}
			},
		},
		{
			name: "ModifyAttributesSimple",
			f: func() []Processor {
				return []Processor{
					attrSetter{},
					NewSimpleProcessor(noopExporter{}),
				}
			},
		},
		{
			name: "ModifyAttributesBatch",
			f: func() []Processor {
				return []Processor{
					attrSetter{},
					NewBatchProcessor(noopExporter{}),
				}
			},
		},
	} {
		b.Run(tc.name, func(b *testing.B) {
			var opts []LoggerProviderOption
			for _, p := range tc.f() {
				opts = append(opts, WithProcessor(p))
			}
			provider := NewLoggerProvider(opts...)
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

type timestampSetter struct{}

func (e timestampSetter) OnEmit(ctx context.Context, r Record) (Record, error) {
	r = r.Clone()
	r.SetObservedTimestamp(time.Date(1988, time.November, 17, 0, 0, 0, 0, time.UTC))
	return r, nil
}

func (e timestampSetter) Enabled(context.Context, Record) bool {
	return true
}

func (e timestampSetter) Shutdown(ctx context.Context) error {
	return nil
}

func (e timestampSetter) ForceFlush(ctx context.Context) error {
	return nil
}

type attrSetter struct{}

func (e attrSetter) OnEmit(ctx context.Context, r Record) (Record, error) {
	r = r.Clone()
	r.SetAttributes(log.String("replace", "me"))
	return r, nil
}

func (e attrSetter) Enabled(context.Context, Record) bool {
	return true
}

func (e attrSetter) Shutdown(ctx context.Context) error {
	return nil
}

func (e attrSetter) ForceFlush(ctx context.Context) error {
	return nil
}
