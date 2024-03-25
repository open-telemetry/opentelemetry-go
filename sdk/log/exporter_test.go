// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log

import "context"

type testExporter struct {
	// Err is the error returned by all methods of the testExporter.
	Err error

	// Counts of method calls.
	ExportN, ShutdownN, ForceFlushN int
	// Records are the Records passed to export.
	Records [][]Record
}

func (e *testExporter) Export(ctx context.Context, r []Record) error {
	e.ExportN++
	e.Records = append(e.Records, r)
	return e.Err
}

func (e *testExporter) Shutdown(ctx context.Context) error {
	e.ShutdownN++
	return e.Err
}

func (e *testExporter) ForceFlush(ctx context.Context) error {
	e.ForceFlushN++
	return e.Err
}
