// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"
	"slices"
	"sync/atomic"
)

type instruction struct {
	Record *[]Record
	Flush  chan [][]Record
}

type testExporter struct {
	// Err is the error returned by all methods of the testExporter.
	Err error

	// Counts of method calls.
	exportN, shutdownN, forceFlushN *int32

	input chan instruction
	done  chan struct{}
}

func newTestExporter(err error) *testExporter {
	e := &testExporter{
		Err:         err,
		exportN:     new(int32),
		shutdownN:   new(int32),
		forceFlushN: new(int32),
		input:       make(chan instruction),
	}
	e.done = run(e.input)

	return e
}

func run(input chan instruction) chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)

		var records [][]Record
		for in := range input {
			if in.Record != nil {
				records = append(records, *in.Record)
			}
			if in.Flush != nil {
				cp := slices.Clone(records)
				records = records[:0]
				in.Flush <- cp
			}
		}
	}()
	return done
}

func (e *testExporter) Records() [][]Record {
	out := make(chan [][]Record, 1)
	e.input <- instruction{Flush: out}
	return <-out
}

func (e *testExporter) Export(ctx context.Context, r []Record) error {
	atomic.AddInt32(e.exportN, 1)
	e.input <- instruction{Record: &r}
	return e.Err
}

func (e *testExporter) ExportN() int {
	return int(atomic.LoadInt32(e.exportN))
}

func (e *testExporter) Stop() {
	close(e.input)
	<-e.done
}

func (e *testExporter) Shutdown(ctx context.Context) error {
	atomic.AddInt32(e.shutdownN, 1)
	return e.Err
}

func (e *testExporter) ShutdownN() int {
	return int(atomic.LoadInt32(e.shutdownN))
}

func (e *testExporter) ForceFlush(ctx context.Context) error {
	atomic.AddInt32(e.forceFlushN, 1)
	return e.Err
}

func (e *testExporter) ForceFlushN() int {
	return int(atomic.LoadInt32(e.forceFlushN))
}
