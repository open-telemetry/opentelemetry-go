// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package xlog

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	logapi "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
)

func TestLoggerEnabled(t *testing.T) {
	p0 := newFltrProcessor(true)
	p1 := newFltrProcessor(true)
	p2WithDisabled := newFltrProcessor(false)

	emptyResource := resource.Empty()
	res := resource.NewSchemaless(attribute.String("key", "value"))

	testCases := []struct {
		name             string
		logger           logapi.Logger
		ctx              context.Context
		expected         bool
		expectedP0Params []EnabledParameters
		expectedP1Params []EnabledParameters
		expectedP2Params []EnabledParameters
	}{
		{
			name:     "NoProcessors",
			logger:   log.NewLoggerProvider().Logger(t.Name()),
			ctx:      context.Background(),
			expected: false,
		},
		{
			name: "WithProcessors",
			logger: log.NewLoggerProvider(
				log.WithProcessor(p0),
				log.WithProcessor(p1),
				log.WithResource(res),
			).Logger(t.Name()),
			ctx:      context.Background(),
			expected: true,
			expectedP0Params: []EnabledParameters{{
				Resource:             *res,
				InstrumentationScope: instrumentation.Scope{Name: t.Name()},
			}},
			expectedP1Params: nil,
		},
		{
			name: "WithDisabledProcessors",
			logger: log.NewLoggerProvider(
				log.WithProcessor(p2WithDisabled),
				log.WithResource(emptyResource),
			).Logger(t.Name()),
			ctx:              context.Background(),
			expected:         false,
			expectedP2Params: []EnabledParameters{{InstrumentationScope: instrumentation.Scope{Name: t.Name()}}},
		},
		{
			name: "ContainsDisabledProcessor",
			logger: log.NewLoggerProvider(
				log.WithProcessor(p2WithDisabled),
				log.WithProcessor(p0),
				log.WithResource(emptyResource),
			).Logger(t.Name()),
			ctx:              context.Background(),
			expected:         true,
			expectedP2Params: []EnabledParameters{{InstrumentationScope: instrumentation.Scope{Name: t.Name()}}},
			expectedP0Params: []EnabledParameters{{InstrumentationScope: instrumentation.Scope{Name: t.Name()}}},
		},
		{
			name: "WithNilContext",
			logger: log.NewLoggerProvider(
				log.WithProcessor(p0),
				log.WithProcessor(p1),
				log.WithResource(emptyResource),
			).Logger(t.Name()),
			ctx:              nil,
			expected:         true,
			expectedP0Params: []EnabledParameters{{InstrumentationScope: instrumentation.Scope{Name: t.Name()}}},
			expectedP1Params: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clean up the records before the test.
			p0.params = nil
			p1.params = nil
			p2WithDisabled.params = nil

			assert.Equal(t, tc.expected, tc.logger.Enabled(tc.ctx, logapi.EnabledParameters{}))
			assert.Equal(t, tc.expectedP0Params, p0.params)
			assert.Equal(t, tc.expectedP1Params, p1.params)
			assert.Equal(t, tc.expectedP2Params, p2WithDisabled.params)
		})
	}
}

type fltrProcessor struct {
	enabled bool
	params  []EnabledParameters
}

var _ FilterProcessor = (*fltrProcessor)(nil)

func newFltrProcessor(enabled bool) *fltrProcessor {
	return &fltrProcessor{
		enabled: enabled,
	}
}

func (p *fltrProcessor) Enabled(_ context.Context, params EnabledParameters) bool {
	p.params = append(p.params, params)
	return p.enabled
}

func (p *fltrProcessor) OnEmit(ctx context.Context, r *log.Record) error {
	return nil
}

func (p *fltrProcessor) Shutdown(context.Context) error {
	return nil
}

func (p *fltrProcessor) ForceFlush(context.Context) error {
	return nil
}
