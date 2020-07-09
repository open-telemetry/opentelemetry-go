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

package detect

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/sdk/resource"
)

func TestDetectOnePairTrue(t *testing.T) {
	os.Setenv(envVar, "key=value")

	detector := &FromEnv{}
	res, err := detector.Detect(context.Background())
	require.NoError(t, err)
	assert.Equal(t, resource.New(kv.String("key", "value")), res)
}

func TestDetectMultiplePairsTrue(t *testing.T) {
	os.Setenv(envVar, "key= value, k = v ")

	detector := &FromEnv{}
	res, err := detector.Detect(context.Background())
	require.NoError(t, err)
	assert.Equal(t, resource.New(kv.String("key", "value"), kv.String("k", "v")), res)
}

func TestDetectFalse(t *testing.T) {
	os.Setenv(envVar, "")

	detector := &FromEnv{}
	res, err := detector.Detect(context.Background())
	require.NoError(t, err)
	assert.Equal(t, resource.Empty(), res)
}

func TestDetectError(t *testing.T) {
	os.Setenv(envVar, "key=value,key")

	detector := &FromEnv{}
	res, err := detector.Detect(context.Background())
	assert.Error(t, err)
	assert.Equal(t, resource.Empty(), res)
}
