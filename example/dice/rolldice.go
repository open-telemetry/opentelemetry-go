// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	// TODO: https://github.com/open-telemetry/opentelemetry-go/issues/5801
	// "go.opentelemetry.io/contrib/bridges/otelslog".
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const name = "go.opentelemetry.io/otel/example/dice"

var (
	tracer  = otel.Tracer(name)
	meter   = otel.Meter(name)
	logger  = slog.New(slog.NewJSONHandler(os.Stdout, nil)) // TODO: logger  = otelslog.NewLogger(name).
	rollCnt metric.Int64Counter
)

func init() {
	var err error
	rollCnt, err = meter.Int64Counter("dice.rolls",
		metric.WithDescription("The number of rolls by roll value"),
		metric.WithUnit("{roll}"))
	if err != nil {
		panic(err)
	}
}
