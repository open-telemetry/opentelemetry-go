// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resource

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	ottest "go.opentelemetry.io/otel/internal/internaltest"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
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
		semconv.ServiceNameKey.String("bar"),
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
		semconv.ServiceNameKey.String("bar"),
	))
}
