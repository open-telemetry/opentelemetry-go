// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"sync"
	"time"

	"go.opentelemetry.io/otel/sdk/metric/exemplar"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

var exemplarPool = sync.Pool{
	New: func() any { return new([]exemplar.Exemplar) },
}

func collectExemplars[N int64 | float64](out *[]metricdata.Exemplar[N], f func(*[]exemplar.Exemplar)) {
	collectExemplarsAfter[N](out, time.Time{}, f)
}

func collectExemplarsAfter[N int64 | float64](
	out *[]metricdata.Exemplar[N],
	startTime time.Time,
	f func(*[]exemplar.Exemplar),
) {
	dest := exemplarPool.Get().(*[]exemplar.Exemplar)
	defer func() {
		clear(*dest) // Erase elements to let GC collect objects.
		*dest = (*dest)[:0]
		exemplarPool.Put(dest)
	}()

	*dest = reset(*dest, len(*out), cap(*out))

	f(dest)

	*out = reset(*out, len(*dest), cap(*dest))
	i := 0
	for _, e := range *dest {
		if !startTime.IsZero() && e.Time.Before(startTime) {
			continue
		}
		(*out)[i].FilteredAttributes = e.FilteredAttributes
		(*out)[i].Time = e.Time
		(*out)[i].SpanID = e.SpanID
		(*out)[i].TraceID = e.TraceID

		switch e.Value.Type() {
		case exemplar.Int64ValueType:
			(*out)[i].Value = N(e.Value.Int64())
		case exemplar.Float64ValueType:
			(*out)[i].Value = N(e.Value.Float64())
		}
		i++
	}
	*out = (*out)[:i]
}
