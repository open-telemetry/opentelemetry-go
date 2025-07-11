// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resource

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
)

func TestDetectOnePair(t *testing.T) {
	t.Setenv(resourceAttrKey, "key=value")

	detector := &fromEnv{}
	res, err := detector.Detect(context.Background())
	require.NoError(t, err)
	assert.Equal(t, NewSchemaless(attribute.String("key", "value")), res)
}

func TestDetectURIEncodingOnePair(t *testing.T) {
	t.Setenv(resourceAttrKey, "key=x+y+z?q=123")

	detector := &fromEnv{}
	res, err := detector.Detect(context.Background())
	require.NoError(t, err)
	assert.Equal(t, NewSchemaless(attribute.String("key", "x+y+z?q=123")), res)
}

func TestDetectMultiPairs(t *testing.T) {
	t.Setenv("x", "1")
	t.Setenv(resourceAttrKey, "key=value, k = v , a= x, a=z, b=c%2Fd")

	detector := &fromEnv{}
	res, err := detector.Detect(context.Background())
	require.NoError(t, err)
	assert.Equal(t, NewSchemaless(
		attribute.String("key", "value"),
		attribute.String("k", "v"),
		attribute.String("a", "x"),
		attribute.String("a", "z"),
		attribute.String("b", "c/d"),
	), res)
}

func TestDetectURIEncodingMultiPairs(t *testing.T) {
	t.Setenv("x", "1")
	t.Setenv(resourceAttrKey, "key=x+y+z,namespace=localhost/test&verify")
	detector := &fromEnv{}
	res, err := detector.Detect(context.Background())
	require.NoError(t, err)
	assert.Equal(t, NewSchemaless(
		attribute.String("key", "x+y+z"),
		attribute.String("namespace", "localhost/test&verify"),
	), res)
}

func TestEmpty(t *testing.T) {
	t.Setenv(resourceAttrKey, "   ")
	detector := &fromEnv{}
	res, err := detector.Detect(context.Background())
	require.NoError(t, err)
	assert.Equal(t, Empty(), res)
}

func TestNoResourceAttributesSet(t *testing.T) {
	t.Setenv(svcNameKey, "bar")
	detector := &fromEnv{}
	res, err := detector.Detect(context.Background())
	require.NoError(t, err)
	assert.Equal(t, res, NewSchemaless(
		semconv.ServiceName("bar"),
	))
}

func TestMissingKeyError(t *testing.T) {
	t.Setenv(resourceAttrKey, "key=value,key")
	detector := &fromEnv{}
	res, err := detector.Detect(context.Background())
	assert.Error(t, err)
	assert.Equal(t, err, fmt.Errorf("%w: %v", errMissingValue, "[key]"))
	assert.Equal(t, res, NewSchemaless(
		attribute.String("key", "value"),
	))
}

func TestInvalidPercentDecoding(t *testing.T) {
	t.Setenv(resourceAttrKey, "key=%invalid")
	detector := &fromEnv{}
	res, err := detector.Detect(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, NewSchemaless(
		attribute.String("key", "%invalid"),
	), res)
}

func TestDetectServiceNameFromEnv(t *testing.T) {
	t.Setenv(resourceAttrKey, "key=value,service.name=foo")
	t.Setenv(svcNameKey, "bar")
	detector := &fromEnv{}
	res, err := detector.Detect(context.Background())
	require.NoError(t, err)
	assert.Equal(t, res, NewSchemaless(
		attribute.String("key", "value"),
		semconv.ServiceName("bar"),
	))
}
