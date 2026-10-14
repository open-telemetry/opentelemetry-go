// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package finish

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Func marks a metric series as finished at t.
//
// attrs is the complete, unfiltered attribute collection supplied by the
// instrument caller. Implementations must not retain or modify attrs.
type Func func(attrs []attribute.KeyValue, t time.Time)

type int64Counter struct {
	metric.Int64Counter
	finishers []Func
}

// NewInt64Counter decorates counter with Finish support.
func NewInt64Counter(counter metric.Int64Counter, finishers []Func) metric.Int64Counter {
	return &int64Counter{Int64Counter: counter, finishers: finishers}
}

func (i *int64Counter) Finish(_ context.Context, attrs ...attribute.KeyValue) {
	t := time.Now()
	for _, f := range i.finishers {
		f(attrs, t)
	}
}
