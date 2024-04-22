// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package logtest is a testing helper package.
package logtest_test

import (
	"context"
	"fmt"
	"io"
	"os"

	logapi "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/log/logtest"
)

func ExampleRecordBuilder() {
	exp := exporter{os.Stdout}
	b := logtest.RecordBuilder{}.SetInstrumentationScope(instrumentation.Scope{Name: "myapp"})
	r1 := b.SetBody(logapi.StringValue("foo")).Record()
	r2 := b.SetBody(logapi.StringValue("bar")).Record()

	_ = exp.Export(context.Background(), []log.Record{r1, r2})

	// Output:
	// scope=myapp msg=foo
	// scope=myapp msg=bar
}

// Compile time check exporter implements log.Exporter.
var _ log.Exporter = exporter{}

type exporter struct{ io.Writer }

func (e exporter) Export(ctx context.Context, records []log.Record) error {
	for i, r := range records {
		if i != 0 {
			if _, err := e.Write([]byte("\n")); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(e, "scope=%s msg=%s", r.InstrumentationScope().Name, r.Body().String()); err != nil {
			return err
		}
	}
	return nil
}

func (e exporter) Shutdown(context.Context) error {
	return nil
}

// appropriate error should be returned in these situations.
func (e exporter) ForceFlush(context.Context) error {
	return nil
}
