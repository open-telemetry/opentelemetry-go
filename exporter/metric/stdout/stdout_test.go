package stdout_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/exporter/metric/stdout"
	"go.opentelemetry.io/otel/exporter/metric/test"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/counter"
	aggtest "go.opentelemetry.io/otel/sdk/metric/aggregator/test"
)

type testFixture struct {
	t        *testing.T
	ctx      context.Context
	exporter *stdout.Exporter
	output   *bytes.Buffer
}

func newFixture(t *testing.T, options stdout.Options) testFixture {
	buf := &bytes.Buffer{}
	options.File = buf
	options.DoNotPrintTime = true
	exp, err := stdout.New(options)
	if err != nil {
		t.Fatal("Error building fixture: ", err)
	}
	return testFixture{
		t:        t,
		ctx:      context.Background(),
		exporter: exp,
		output:   buf,
	}
}

func (fix testFixture) Output() string {
	return strings.TrimSpace(fix.output.String())
}

func (fix testFixture) Export(producer export.Producer) {
	err := fix.exporter.Export(fix.ctx, producer)
	if err != nil {
		fix.t.Error("export failed: ", err)
	}
}

func TestStdoutInvalidQuantile(t *testing.T) {
	_, err := stdout.New(stdout.Options{
		Quantiles: []float64{1.1, 0.9},
	})
	require.Error(t, err, "Invalid quantile error expected")
	require.Equal(t, aggregator.ErrInvalidQuantile, err)
}

func TestStdoutCounterFormat(t *testing.T) {
	fix := newFixture(t, stdout.Options{})

	producer := test.NewProducer(sdk.DefaultLabelEncoder())

	desc := export.NewDescriptor("test.name", export.CounterKind, nil, "", "", core.Int64NumberKind, false)
	cagg := counter.New()
	aggtest.CheckedUpdate(fix.t, cagg, core.NewInt64Number(123), desc)
	cagg.Checkpoint(fix.ctx, desc)

	producer.Add(desc, cagg, key.String("A", "B"), key.String("C", "D"))

	fix.Export(producer)

	require.Equal(t, `{"updates":[{"name":"test.name{A=B,C=D}","sum":"123"}]}`, fix.Output())
}
