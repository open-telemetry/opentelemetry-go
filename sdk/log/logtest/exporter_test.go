// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logtest

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

func TestInMemoryExporterImplementsExporter(t *testing.T) {
	assert.Implements(t, (*sdklog.Exporter)(nil), NewInMemoryExporter())
}

func TestNewInMemoryExporter(t *testing.T) {
	imsb := NewInMemoryExporter()

	require.NoError(t, imsb.Export(context.Background(), nil))
	assert.Len(t, imsb.GetRecords(), 0)

	input := make([]sdklog.Record, 10)
	for i := 0; i < 10; i++ {
		input[i] = sdklog.Record{}
		input[i].SetBody(log.StringValue(fmt.Sprintf("record %d", i)))
	}
	require.NoError(t, imsb.Export(context.Background(), input))
	sds := imsb.GetRecords()
	assert.Len(t, sds, 10)
	for i, sd := range sds {
		assert.Equal(t, input[i], sd)
	}
	imsb.Reset()
	// Ensure that operations on the internal storage does not change the previously returned value.
	assert.Len(t, sds, 10)
	assert.Len(t, imsb.GetRecords(), 0)

	require.NoError(t, imsb.Export(context.Background(), input[0:1]))
	sds = imsb.GetRecords()
	assert.Len(t, sds, 1)
	assert.Equal(t, input[0], sds[0])
}

func TestInMemoryExporterReset(t *testing.T) {
	imsb := NewInMemoryExporter()

	input := make([]sdklog.Record, 10)
	for i := 0; i < 10; i++ {
		input[i] = sdklog.Record{}
		input[i].SetBody(log.StringValue(fmt.Sprintf("record %d", i)))
	}
	require.NoError(t, imsb.Export(context.Background(), input))

	assert.Len(t, imsb.GetRecords(), 10)
	imsb.Reset()
	assert.Len(t, imsb.GetRecords(), 0)
}
