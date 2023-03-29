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
	"fmt"
	"sort"

	"github.com/Masterminds/semver/v3"
	"go.opentelemetry.io/otel/attribute"
	ast10 "go.opentelemetry.io/otel/schema/v1.0/ast"
	"go.opentelemetry.io/otel/schema/v1.1/ast"
	"go.opentelemetry.io/otel/schema/v1.1/types"
)

// Upgrade upgrades attrs in place with schema.
func Upgrade(schema *ast.Schema, attrs []attribute.KeyValue) error {
	vers, err := versions(schema, nil, false)
	if err != nil {
		return fmt.Errorf("upgrade error: %w", err)
	}

	a := newAttributes(attrs)
	for _, v := range vers {
		vDef, ok := schema.Versions[v]
		if !ok {
			return fmt.Errorf("upgrade error: version parsing: %s", v)
		}
		f := a.RenameFunc()
		// Transformations in section "all" always are applied first.
		for _, c := range vDef.All.Changes {
			forEach(c.RenameAttributes.AttributeMap, f)
		}
		// Only other applicable section is for resources.
		for _, c := range vDef.Resources.Changes {
			forEach(c.RenameAttributes.AttributeMap, f)
		}
	}

	return nil
}

// Downgrade downgrade attrs to the schema version of url in place with schema.
func Downgrade(schema *ast.Schema, url string, attrs []attribute.KeyValue) error {
	var min *semver.Version
	if url != "" {
		var err error
		min, err = version(url)
		if err != nil {
			return fmt.Errorf("downgrade error: %w", err)
		}
	}

	vers, err := versions(schema, min, true)
	if err != nil {
		return fmt.Errorf("downgrade error: %w", err)
	}

	a := newAttributes(attrs)
	for _, v := range vers {
		vDef, ok := schema.Versions[v]
		if !ok {
			return fmt.Errorf("downgrade error: version parsing: %s", v)
		}
		f := a.UnrenameFunc()
		changes := vDef.Resources.Changes
		for i := len(changes) - 1; i >= 0; i-- {
			c := changes[i]
			forEach(c.RenameAttributes.AttributeMap, f)
		}
		// Downgraing means all transformations in section "all" always are
		// applied after the resource section is resolved.
		changes = vDef.All.Changes
		for i := len(changes) - 1; i >= 0; i-- {
			c := changes[i]
			forEach(c.RenameAttributes.AttributeMap, f)
		}
	}

	return nil
}

// versions returns the sorted versions contained in schema.
func versions(schema *ast.Schema, min *semver.Version, reverse bool) ([]types.TelemetryVersion, error) {
	// The transformations specified in each version are applied one by one.
	// Order the versions to ensure correct application.
	versions := make([]*semver.Version, 0, len(schema.Versions))
	for telV := range schema.Versions {
		v, err := semver.NewVersion(string(telV))
		if err != nil {
			return nil, fmt.Errorf("telemetry version %q: %w", telV, err)
		}
		versions = append(versions, v)
	}

	var sIface sort.Interface = semver.Collection(versions)
	if reverse {
		sIface = sort.Reverse(sIface)
	}
	sort.Sort(sIface)

	out := make([]types.TelemetryVersion, len(versions))
	for i := range versions {
		if min != nil && min.GreaterThan(versions[i]) {
			continue
		}
		out[i] = types.TelemetryVersion(versions[i].String())
	}
	return out, nil
}

func forEach(m ast10.AttributeMap, f func(string, string)) {
	for key, val := range m {
		f(key, val)
	}
}

type attributes struct {
	underlying []attribute.KeyValue

	index map[string]int
}

func newAttributes(attr []attribute.KeyValue) *attributes {
	// Lazy built index.
	return &attributes{underlying: attr}
}

func (a *attributes) init() {
	a.index = make(map[string]int, len(a.underlying))
	for i := range a.underlying {
		a.index[string(a.underlying[i].Key)] = i
	}
}

func (a *attributes) ReplaceKey(orig, repl string) {
	if a.index == nil {
		a.init()
	}

	i, ok := a.index[orig]
	if !ok {
		return
	}
	delete(a.index, orig)
	a.underlying[i].Key = attribute.Key(repl)
	a.index[repl] = i
}

func (a *attributes) RenameFunc() func(string, string) {
	return func(orig, repl string) { a.ReplaceKey(orig, repl) }
}

func (a *attributes) UnrenameFunc() func(string, string) {
	return func(orig, repl string) { a.ReplaceKey(repl, orig) }
}
