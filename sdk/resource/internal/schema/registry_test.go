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

package schema

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/schema/v1.1/ast"
)

var (
	schemaURL  = "https://opentelemetry.io/schemas/1.0.0"
	schemaYAML = `
file_format: 1.1.0
schema_url: ` + schemaURL + `
versions:
  1.0.0:
`
)

func TestRegistryGet(t *testing.T) {
	msg := new(string)
	*msg = schemaYAML
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, *msg)
	}))
	t.Cleanup(ts.Close)

	ctx := context.Background()
	reg := NewRegistry(ts.Client())

	assertGet := func(s *ast.Schema, err error) {
		t.Helper()
		require.NoError(t, err)
		require.NotNil(t, s)
		assert.Equal(t, schemaURL, s.SchemaURL)
	}

	assertGet(reg.Get(ctx, ts.URL))

	// Cache miss.
	assertGet(reg.Get(ctx, ts.URL+"/extra"))

	// Cache hit. This will fail to parse if the HTTP request is actually made.
	*msg = "first"
	assertGet(reg.Get(ctx, ts.URL))
}
