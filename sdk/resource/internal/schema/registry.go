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
	"net/http"

	sUtil "go.opentelemetry.io/otel/schema/v1.1"
	"go.opentelemetry.io/otel/schema/v1.1/ast"
)

type cache struct {
	data map[string]*ast.Schema
}

func (c *cache) lookup(key string, f func() (*ast.Schema, error)) (*ast.Schema, error) {
	if c.data == nil {
		s, err := f()
		if err != nil {
			return nil, err
		}
		c.data = map[string]*ast.Schema{key: s}
		return s, nil
	}

	if s, ok := c.data[key]; ok {
		return s, nil
	}

	s, err := f()
	if err != nil {
		return nil, err
	}
	c.data = map[string]*ast.Schema{key: s}
	return s, nil
}

// Registry hold a registration of schema files. It will cache any schema files
// it gets from external URLs.
type Registry struct {
	client *http.Client

	cache cache
}

// NewRegistry returns a Registry that uses the HTTP client. If client is nil
// it will use the default client from "net/http".
func NewRegistry(client *http.Client) *Registry {
	if client == nil {
		client = http.DefaultClient
	}
	return &Registry{client: client}
}

// Get returns the Schema at the target schemaURL using the registry client.
func (r *Registry) Get(ctx context.Context, schemaURL string) (*ast.Schema, error) {
	return r.cache.lookup(schemaURL, func() (*ast.Schema, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, schemaURL, http.NoBody)
		if err != nil {
			return nil, err
		}
		resp, err := r.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		return sUtil.Parse(resp.Body)
	})
}
