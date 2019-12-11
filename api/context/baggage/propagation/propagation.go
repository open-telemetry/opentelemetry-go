package propagation

import (
	"context"
	"net/url"
	"strings"

	"go.opentelemetry.io/otel/api/context/baggage"
	"go.opentelemetry.io/otel/api/context/propagation"
	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"
)

type ctxEntriesType struct{}

var (
	// CorrelationContextHeader is specified by W3C.
	CorrelationContextHeader = "Correlation-Context"

	ctxEntriesKey = &ctxEntriesType{}
)

// CorrelationContext propagates Key:Values in W3C TraceContext format.
type CorrelationContext struct{}

var _ propagation.HTTPPropagator = CorrelationContext{}

// WithMap enters a baggage.Map into a new Context.
func WithMap(ctx context.Context, m baggage.Map) context.Context {
	return context.WithValue(ctx, ctxEntriesKey, m)
}

// WithMap enters a key:value set into a new Context.
func NewContext(ctx context.Context, keyvalues ...core.KeyValue) context.Context {
	return WithMap(ctx, FromContext(ctx).Apply(baggage.MapUpdate{
		MultiKV: keyvalues,
	}))
}

// FromContext gets the current baggage.Map from a Context.
func FromContext(ctx context.Context) baggage.Map {
	if m, ok := ctx.Value(ctxEntriesKey).(baggage.Map); ok {
		return m
	}
	return baggage.NewEmptyMap()
}

// Inject implements HTTPInjector.
func (CorrelationContext) Inject(ctx context.Context, supplier propagation.HTTPSupplier) {
	correlationCtx := FromContext(ctx)
	firstIter := true
	var headerValueBuilder strings.Builder
	correlationCtx.Foreach(func(kv core.KeyValue) bool {
		if !firstIter {
			headerValueBuilder.WriteRune(',')
		}
		firstIter = false
		headerValueBuilder.WriteString(url.QueryEscape(strings.TrimSpace((string)(kv.Key))))
		headerValueBuilder.WriteRune('=')
		headerValueBuilder.WriteString(url.QueryEscape(strings.TrimSpace(kv.Value.Emit())))
		return true
	})
	if headerValueBuilder.Len() > 0 {
		headerString := headerValueBuilder.String()
		supplier.Set(CorrelationContextHeader, headerString)
	}
}

// Inject implements HTTPExtractor.
func (CorrelationContext) Extract(ctx context.Context, supplier propagation.HTTPSupplier) context.Context {
	correlationContext := supplier.Get(CorrelationContextHeader)
	if correlationContext == "" {
		return WithMap(ctx, baggage.NewEmptyMap())
	}

	contextValues := strings.Split(correlationContext, ",")
	keyValues := make([]core.KeyValue, 0, len(contextValues))
	for _, contextValue := range contextValues {
		valueAndProps := strings.Split(contextValue, ";")
		if len(valueAndProps) < 1 {
			continue
		}
		nameValue := strings.Split(valueAndProps[0], "=")
		if len(nameValue) < 2 {
			continue
		}
		name, err := url.QueryUnescape(nameValue[0])
		if err != nil {
			continue
		}
		trimmedName := strings.TrimSpace(name)
		value, err := url.QueryUnescape(nameValue[1])
		if err != nil {
			continue
		}
		trimmedValue := strings.TrimSpace(value)

		// TODO (skaris): properties defiend https://w3c.github.io/correlation-context/, are currently
		// just put as part of the value.
		var trimmedValueWithProps strings.Builder
		trimmedValueWithProps.WriteString(trimmedValue)
		for _, prop := range valueAndProps[1:] {
			trimmedValueWithProps.WriteRune(';')
			trimmedValueWithProps.WriteString(prop)
		}

		keyValues = append(keyValues, key.New(trimmedName).String(trimmedValueWithProps.String()))
	}
	return WithMap(ctx, baggage.NewMap(baggage.MapUpdate{
		MultiKV: keyValues,
	}))
}

// GetAllKeys implements HTTPPropagator.
func (CorrelationContext) GetAllKeys() []string {
	return []string{CorrelationContextHeader}
}
