// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// These benchmarks are based on slog/internal/benchmarks.
//
// They test a complete log record, from the user's call to its return.

package internal

import (
	"context"
	"io"
	"testing"
	"time"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/noop"
	"go.opentelemetry.io/otel/trace"

	"github.com/go-logr/logr"
	"golang.org/x/exp/slog"
)

var (
	ctx            = trace.ContextWithSpanContext(context.Background(), trace.NewSpanContext(trace.SpanContextConfig{TraceID: [16]byte{1}, SpanID: [8]byte{42}}))
	testTimestamp  = time.Date(1988, time.November, 17, 0, 0, 0, 0, time.UTC)
	testBodyString = "log message"
	testBody       = log.StringValue(testBodyString)
	testSeverity   = log.SeverityInfo
	testFloat      = 1.2345
	testString     = "7e3b3b2aaeff56a7108fe11e154200dd/7819479873059528190"
	testInt        = 32768
	testBool       = true
)

// WriterLogger is an optimistic version of a real logger, doing real-world
// tasks as fast as possible . This gives us an upper bound
// on handler performance, so we can evaluate the (logger-independent) core
// activity of the package in an end-to-end context without concern that a
// slow logger implementation is skewing the results. The writerLogger
// allocates memory only when using strconv.
func BenchmarkEmit(b *testing.B) {
	for _, tc := range []struct {
		name   string
		logger log.Logger
	}{
		{"noop", noop.Logger{}},
		{"writer", &writerLogger{w: io.Discard}},
	} {
		b.Run(tc.name, func(b *testing.B) {
			for _, call := range []struct {
				name string
				f    func()
			}{
				{
					"no attrs",
					func() {
						r := log.Record{}
						r.SetTimestamp(testTimestamp)
						r.SetSeverity(testSeverity)
						r.SetBody(testBody)
						tc.logger.Emit(ctx, r)
					},
				},
				{
					"3 attrs",
					func() {
						r := log.Record{}
						r.SetTimestamp(testTimestamp)
						r.SetSeverity(testSeverity)
						r.SetBody(testBody)
						r.AddAttributes(
							log.String("string", testString),
							log.Float64("float", testFloat),
							log.Int("int", testInt),
						)
						tc.logger.Emit(ctx, r)
					},
				},
				{
					// The number should match nAttrsInline in record.go and in slog/record.go.
					// This should exercise the code path where no allocations
					// happen in Record or Attr. If there are allocations, they
					// should only be from strconv used in writerLogger.
					"5 attrs",
					func() {
						r := log.Record{}
						r.SetTimestamp(testTimestamp)
						r.SetSeverity(testSeverity)
						r.SetBody(testBody)
						r.AddAttributes(
							log.String("string", testString),
							log.Float64("float", testFloat),
							log.Int("int", testInt),
							log.Bool("bool", testBool),
							log.String("string", testString),
						)
						tc.logger.Emit(ctx, r)
					},
				},
				{
					"10 attrs",
					func() {
						r := log.Record{}
						r.SetTimestamp(testTimestamp)
						r.SetSeverity(testSeverity)
						r.SetBody(testBody)
						r.AddAttributes(
							log.String("string", testString),
							log.Float64("float", testFloat),
							log.Int("int", testInt),
							log.Bool("bool", testBool),
							log.String("string", testString),
							log.String("string", testString),
							log.Float64("float", testFloat),
							log.Int("int", testInt),
							log.Bool("bool", testBool),
							log.String("string", testString),
						)
						tc.logger.Emit(ctx, r)
					},
				},
				{
					"40 attrs",
					func() {
						r := log.Record{}
						r.SetTimestamp(testTimestamp)
						r.SetSeverity(testSeverity)
						r.SetBody(testBody)
						r.AddAttributes(
							log.String("string", testString),
							log.Float64("float", testFloat),
							log.Int("int", testInt),
							log.Bool("bool", testBool),
							log.String("string", testString),
							log.String("string", testString),
							log.Float64("float", testFloat),
							log.Int("int", testInt),
							log.Bool("bool", testBool),
							log.String("string", testString),
							log.String("string", testString),
							log.Float64("float", testFloat),
							log.Int("int", testInt),
							log.Bool("bool", testBool),
							log.String("string", testString),
							log.String("string", testString),
							log.Float64("float", testFloat),
							log.Int("int", testInt),
							log.Bool("bool", testBool),
							log.String("string", testString),
							log.String("string", testString),
							log.Float64("float", testFloat),
							log.Int("int", testInt),
							log.Bool("bool", testBool),
							log.String("string", testString),
							log.String("string", testString),
							log.Float64("float", testFloat),
							log.Int("int", testInt),
							log.Bool("bool", testBool),
							log.String("string", testString),
							log.String("string", testString),
							log.Float64("float", testFloat),
							log.Int("int", testInt),
							log.Bool("bool", testBool),
							log.String("string", testString),
							log.String("string", testString),
							log.Float64("float", testFloat),
							log.Int("int", testInt),
							log.Bool("bool", testBool),
							log.String("string", testString),
						)
						tc.logger.Emit(ctx, r)
					},
				},
			} {
				b.Run(call.name, func(b *testing.B) {
					b.ReportAllocs()
					b.RunParallel(func(pb *testing.PB) {
						for pb.Next() {
							call.f()
						}
					})
				})
			}
		})
	}
}

func BenchmarkSlog(b *testing.B) {
	logger := slog.New(&slogHandler{noop.Logger{}})
	for _, call := range []struct {
		name string
		f    func()
	}{
		{
			"no attrs",
			func() {
				logger.LogAttrs(ctx, slog.LevelInfo, testBodyString)
			},
		},
		{
			"3 attrs",
			func() {
				logger.LogAttrs(ctx, slog.LevelInfo, testBodyString,
					slog.String("string", testString),
					slog.Float64("float", testFloat),
					slog.Int("int", testInt),
				)
			},
		},
		{
			"5 attrs",
			func() {
				logger.LogAttrs(ctx, slog.LevelInfo, testBodyString,
					slog.String("string", testString),
					slog.Float64("float", testFloat),
					slog.Int("int", testInt),
					slog.Bool("bool", testBool),
					slog.String("string", testString),
				)
			},
		},
		{
			"10 attrs",
			func() {
				logger.LogAttrs(ctx, slog.LevelInfo, testBodyString,
					slog.String("string", testString),
					slog.Float64("float", testFloat),
					slog.Int("int", testInt),
					slog.Bool("bool", testBool),
					slog.String("string", testString),
					slog.String("string", testString),
					slog.Float64("float", testFloat),
					slog.Int("int", testInt),
					slog.Bool("bool", testBool),
					slog.String("string", testString),
				)
			},
		},
		{
			"40 attrs",
			func() {
				logger.LogAttrs(ctx, slog.LevelInfo, testBodyString,
					slog.String("string", testString),
					slog.Float64("float", testFloat),
					slog.Int("int", testInt),
					slog.Bool("bool", testBool),
					slog.String("string", testString),
					slog.String("string", testString),
					slog.Float64("float", testFloat),
					slog.Int("int", testInt),
					slog.Bool("bool", testBool),
					slog.String("string", testString),
					slog.String("string", testString),
					slog.Float64("float", testFloat),
					slog.Int("int", testInt),
					slog.Bool("bool", testBool),
					slog.String("string", testString),
					slog.String("string", testString),
					slog.Float64("float", testFloat),
					slog.Int("int", testInt),
					slog.Bool("bool", testBool),
					slog.String("string", testString),
					slog.String("string", testString),
					slog.Float64("float", testFloat),
					slog.Int("int", testInt),
					slog.Bool("bool", testBool),
					slog.String("string", testString),
					slog.String("string", testString),
					slog.Float64("float", testFloat),
					slog.Int("int", testInt),
					slog.Bool("bool", testBool),
					slog.String("string", testString),
					slog.String("string", testString),
					slog.Float64("float", testFloat),
					slog.Int("int", testInt),
					slog.Bool("bool", testBool),
					slog.String("string", testString),
					slog.String("string", testString),
					slog.Float64("float", testFloat),
					slog.Int("int", testInt),
					slog.Bool("bool", testBool),
					slog.String("string", testString),
				)
			},
		},
	} {
		b.Run(call.name, func(b *testing.B) {
			b.ReportAllocs()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					call.f()
				}
			})
		})
	}
}

func BenchmarkLogr(b *testing.B) {
	logger := logr.New(&logrSink{noop.Logger{}})
	for _, call := range []struct {
		name string
		f    func()
	}{
		{
			"no attrs",
			func() {
				logger.Info(testBodyString)
			},
		},
		{
			"3 attrs",
			func() {
				logger.Info(testBodyString,
					"string", testString,
					"float", testFloat,
					"int", testInt,
				)
			},
		},
		{
			// The number should match nAttrsInline in record.go.
			// This should exercise the code path where no allocations
			// happen in Record or Attr. If there are allocations, they
			// should only be from strconv used in writerLogger.
			"5 attrs",
			func() {
				logger.Info(testBodyString,
					"string", testString,
					"float", testFloat,
					"int", testInt,
					"bool", testBool,
					"string", testString,
				)
			},
		},
		{
			"10 attrs",
			func() {
				logger.Info(testBodyString,
					"string", testString,
					"float", testFloat,
					"int", testInt,
					"bool", testBool,
					"string", testString,
					"string", testString,
					"float", testFloat,
					"int", testInt,
					"bool", testBool,
					"string", testString,
				)
			},
		},
		{
			"40 attrs",
			func() {
				logger.Info(testBodyString,
					"string", testString,
					"float", testFloat,
					"int", testInt,
					"bool", testBool,
					"string", testString,
					"string", testString,
					"float", testFloat,
					"int", testInt,
					"bool", testBool,
					"string", testString,
					"string", testString,
					"float", testFloat,
					"int", testInt,
					"bool", testBool,
					"string", testString,
					"string", testString,
					"float", testFloat,
					"int", testInt,
					"bool", testBool,
					"string", testString,
					"string", testString,
					"float", testFloat,
					"int", testInt,
					"bool", testBool,
					"string", testString,
					"string", testString,
					"float", testFloat,
					"int", testInt,
					"bool", testBool,
					"string", testString,
					"string", testString,
					"float", testFloat,
					"int", testInt,
					"bool", testBool,
					"string", testString,
					"string", testString,
					"float", testFloat,
					"int", testInt,
					"bool", testBool,
					"string", testString,
				)
			},
		},
	} {
		b.Run(call.name, func(b *testing.B) {
			b.ReportAllocs()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					call.f()
				}
			})
		})
	}
}
