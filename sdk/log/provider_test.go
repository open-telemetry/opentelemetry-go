// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
)

type processor struct {
	Name string
	Err  error

	shutdownCalls   int
	forceFlushCalls int

	records []Record
}

func newProcessor(name string) *processor {
	return &processor{Name: name}
}

func (p *processor) OnEmit(ctx context.Context, r Record) error {
	if p.Err != nil {
		return p.Err
	}

	p.records = append(p.records, r)
	return nil
}

func (*processor) Enabled(context.Context, Record) bool {
	return true
}

func (p *processor) Shutdown(context.Context) error {
	p.shutdownCalls++
	return p.Err
}

func (p *processor) ForceFlush(context.Context) error {
	p.forceFlushCalls++
	return p.Err
}

func TestNewLoggerProviderConfiguration(t *testing.T) {
	t.Cleanup(func(orig otel.ErrorHandler) func() {
		otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
			t.Log(err)
		}))
		return func() { otel.SetErrorHandler(orig) }
	}(otel.GetErrorHandler()))

	res := resource.NewSchemaless(attribute.String("key", "value"))
	p0, p1 := newProcessor("0"), newProcessor("1")
	attrCntLim := 12
	attrValLenLim := 21

	testcases := []struct {
		name    string
		envars  map[string]string
		options []LoggerProviderOption
		want    *LoggerProvider
	}{
		{
			name: "Defaults",
			want: &LoggerProvider{
				resource:                  resource.Default(),
				attributeCountLimit:       defaultAttrCntLim,
				attributeValueLengthLimit: defaultAttrValLenLim,
			},
		},
		{
			name: "Options",
			options: []LoggerProviderOption{
				WithResource(res),
				WithProcessor(p0),
				WithProcessor(p1),
				WithAttributeCountLimit(attrCntLim),
				WithAttributeValueLengthLimit(attrValLenLim),
			},
			want: &LoggerProvider{
				resource:                  res,
				processors:                []Processor{p0, p1},
				attributeCountLimit:       attrCntLim,
				attributeValueLengthLimit: attrValLenLim,
			},
		},
		{
			name: "Environment",
			envars: map[string]string{
				envarAttrCntLim:    strconv.Itoa(attrCntLim),
				envarAttrValLenLim: strconv.Itoa(attrValLenLim),
			},
			want: &LoggerProvider{
				resource:                  resource.Default(),
				attributeCountLimit:       attrCntLim,
				attributeValueLengthLimit: attrValLenLim,
			},
		},
		{
			name: "InvalidEnvironment",
			envars: map[string]string{
				envarAttrCntLim:    "invalid attributeCountLimit",
				envarAttrValLenLim: "invalid attributeValueLengthLimit",
			},
			want: &LoggerProvider{
				resource:                  resource.Default(),
				attributeCountLimit:       defaultAttrCntLim,
				attributeValueLengthLimit: defaultAttrValLenLim,
			},
		},
		{
			name: "Precedence",
			envars: map[string]string{
				envarAttrCntLim:    strconv.Itoa(100),
				envarAttrValLenLim: strconv.Itoa(101),
			},
			options: []LoggerProviderOption{
				// These override the environment variables.
				WithAttributeCountLimit(attrCntLim),
				WithAttributeValueLengthLimit(attrValLenLim),
			},
			want: &LoggerProvider{
				resource:                  resource.Default(),
				attributeCountLimit:       attrCntLim,
				attributeValueLengthLimit: attrValLenLim,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			for key, value := range tc.envars {
				t.Setenv(key, value)
			}
			assert.Equal(t, tc.want, NewLoggerProvider(tc.options...))
		})
	}
}

func TestLimitValueFailsOpen(t *testing.T) {
	var l limit
	assert.Equal(t, -1, l.Value(), "limit value should default to unlimited")
}
