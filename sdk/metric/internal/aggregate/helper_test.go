// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
)

func measureDeltaHistogram[N int64 | float64](h *deltaHistogram[N], ctx context.Context, value N, attr attribute.Set) {
	h.measure(ctx, value, newLazyFilteredSet(attr, nil))
}

func measureCumulativeHistogram[N int64 | float64](
	h *cumulativeHistogram[N],
	ctx context.Context,
	value N,
	attr attribute.Set,
) {
	h.measure(ctx, value, newLazyFilteredSet(attr, nil))
}

func measureExpoHistogram[N int64 | float64](h *expoHistogram[N], ctx context.Context, value N, attr attribute.Set) {
	h.measure(ctx, value, newLazyFilteredSet(attr, nil))
}
