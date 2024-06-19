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
	ottest "go.opentelemetry.io/otel/sdk/internal/internaltest"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func TestDetectOnePair(t *testing.T) {
	store, err := ottest.SetEnvVariables(map[string]string{
		resourceAttrKey: "key=value",
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, store.Restore()) }()

	detector := &fromEnv{}
	res, err := detector.Detect(context.Background())
	require.NoError(t, err)
	assert.Equal(t, NewSchemaless(attribute.String("key", "value")), res)
}

func TestDetectURIEncodingOnePair(t *testing.T) {
	store, err := ottest.SetEnvVariables(map[string]string{
		resourceAttrKey: "key=x+y+z?q=123",
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, store.Restore()) }()

	detector := &fromEnv{}
	res, err := detector.Detect(context.Background())
	require.NoError(t, err)
	assert.Equal(t, NewSchemaless(attribute.String("key", "x+y+z?q=123")), res)
}

func TestDetectMultiPairs(t *testing.T) {
	store, err := ottest.SetEnvVariables(map[string]string{
		"x":             "1",
		resourceAttrKey: "key=value, k = v , a= x, a=z, b=c%2Fd",
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, store.Restore()) }()

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
	store, err := ottest.SetEnvVariables(map[string]string{
		"x":             "1",
		resourceAttrKey: "key=x+y+z,namespace=localhost/test&verify",
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, store.Restore()) }()

	detector := &fromEnv{}
	res, err := detector.Detect(context.Background())
	require.NoError(t, err)
	assert.Equal(t, NewSchemaless(
		attribute.String("key", "x+y+z"),
		attribute.String("namespace", "localhost/test&verify"),
	), res)
}

func TestEmpty(t *testing.T) {
	store, err := ottest.SetEnvVariables(map[string]string{
		resourceAttrKey: "   ",
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, store.Restore()) }()

	detector := &fromEnv{}
	res, err := detector.Detect(context.Background())
	require.NoError(t, err)
	assert.Equal(t, Empty(), res)
}

func TestNoResourceAttributesSet(t *testing.T) {
	store, err := ottest.SetEnvVariables(map[string]string{
		svcNameKey: "bar",
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, store.Restore()) }()
	detector := &fromEnv{}
	res, err := detector.Detect(context.Background())
	require.NoError(t, err)
	assert.Equal(t, res, NewSchemaless(
		semconv.ServiceName("bar"),
	))
}

func TestMissingKeyError(t *testing.T) {
	store, err := ottest.SetEnvVariables(map[string]string{
		resourceAttrKey: "key=value,key",
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, store.Restore()) }()

	detector := &fromEnv{}
	res, err := detector.Detect(context.Background())
	assert.Error(t, err)
	assert.Equal(t, err, fmt.Errorf("%w: %v", errMissingValue, "[key]"))
	assert.Equal(t, res, NewSchemaless(
		attribute.String("key", "value"),
	))
}

func TestInvalidPercentDecoding(t *testing.T) {
	store, err := ottest.SetEnvVariables(map[string]string{
		resourceAttrKey: "key=%invalid",
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, store.Restore()) }()

	detector := &fromEnv{}
	res, err := detector.Detect(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, NewSchemaless(
		attribute.String("key", "%invalid"),
	), res)
}

func TestDetectServiceNameFromEnv(t *testing.T) {
	store, err := ottest.SetEnvVariables(map[string]string{
		resourceAttrKey: "key=value,service.name=foo",
		svcNameKey:      "bar",
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, store.Restore()) }()

	detector := &fromEnv{}
	res, err := detector.Detect(context.Background())
	require.NoError(t, err)
	assert.Equal(t, res, NewSchemaless(
		attribute.String("key", "value"),
		semconv.ServiceName("bar"),
	))
}
