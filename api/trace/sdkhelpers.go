package trace

import (
	"context"

	"go.opentelemetry.io/otel/api/core"
)

type TracerWithNamespace interface {
	// Start a span.
	Start(ctx context.Context, name core.Name, opts ...StartOption) (context.Context, Span)

	// WithSpan wraps the execution of the fn function with a span.
	// It starts a new span, sets it as an active span in the context,
	// executes the fn function and closes the span before returning the result of fn.
	WithSpan(
		ctx context.Context,
		name core.Name,
		fn func(ctx context.Context) error,
	) error
}
