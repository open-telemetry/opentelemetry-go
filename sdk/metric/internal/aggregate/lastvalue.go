// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/internal/exemplar"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// datapoint is timestamped measurement data.
type datapoint[N int64 | float64] struct {
	timestamp time.Time
	value     N
	res       exemplar.Reservoir[N]
}

func newLastValue[N int64 | float64](limit int, r func() exemplar.Reservoir[N]) *lastValue[N] {
	return &lastValue[N]{
		newRes: r,
		limit:  newLimiter[datapoint[N]](limit),
		values: make(map[attribute.Set]datapoint[N]),
	}
}

// lastValue summarizes a set of measurements as the last one made.
type lastValue[N int64 | float64] struct {
	sync.Mutex

	newRes func() exemplar.Reservoir[N]
	limit  limiter[datapoint[N]]
	values map[attribute.Set]datapoint[N]
}

func (s *lastValue[N]) measure(ctx context.Context, value N, fltrAttr attribute.Set, droppedAttr []attribute.KeyValue) {
	t := now()

	s.Lock()
	defer s.Unlock()

	attr := s.limit.Attributes(fltrAttr, s.values)
	d, ok := s.values[attr]
	if !ok {
		d.res = s.newRes()
	}

	d.timestamp = t
	d.value = value
	d.res.Offer(ctx, t, value, droppedAttr)

	s.values[attr] = d
}

func (s *lastValue[N]) computeAggregation(dest *[]metricdata.DataPoint[N]) {
	s.Lock()
	defer s.Unlock()

	n := len(s.values)
	*dest = reset(*dest, n, n)

	var i int
	for a, v := range s.values {
		(*dest)[i].Attributes = a
		// The event time is the only meaningful timestamp, StartTime is
		// ignored.
		(*dest)[i].Time = v.timestamp
		(*dest)[i].Value = v.value
		v.res.Collect(&(*dest)[i].Exemplars)
		i++
	}
	// Do not report stale values.
	clear(s.values)
}
