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

package cmd // import "go.opentelemetry.io/otel/sdk/resource/internal/schema/generate/cmd"

import (
	"bytes"
	"embed"
	"fmt"
	"go/format"
	"io"
	"os"
	"sort"
	"text/template"

	"github.com/Masterminds/semver/v3"

	ast10 "go.opentelemetry.io/otel/schema/v1.0/ast"
	"go.opentelemetry.io/otel/schema/v1.1/ast"
	"go.opentelemetry.io/otel/schema/v1.1/types"
)

// Template source.
const src = "transforms.go.tmpl"

//go:embed transforms.go.tmpl
var tmpl embed.FS

// Render renders the passed Schema into the formatted code at dest.
func Render(dest string, s *ast.Schema) error {
	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to open destination %q: %w", dest, err)
	}
	defer f.Close()

	data, err := entries(s)
	if err != nil {
		return fmt.Errorf("failed to parse schema: %w", err)
	}

	err = render(f, data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}
	return nil
}

type entry struct {
	Version   *semver.Version
	All       ast10.Attributes
	Resources ast10.Attributes
}

func newEntry(ver types.TelemetryVersion, def ast.VersionDef) (entry, error) {
	v, err := semver.NewVersion(string(ver))
	return entry{
		Version:   v,
		All:       def.All,
		Resources: def.Resources,
	}, err
}

func entries(s *ast.Schema) ([]entry, error) {
	if s == nil {
		return nil, nil
	}

	es := make([]entry, 0, len(s.Versions))
	for v, def := range s.Versions {
		e, err := newEntry(v, def)
		if err != nil {
			return nil, err
		}

		es = append(es, e)
	}

	sort.Slice(es, func(i, j int) bool {
		return es[i].Version.LessThan(es[j].Version)
	})

	return es, nil
}

func render(dest io.Writer, data any) error {
	t, err := template.ParseFS(tmpl, src)
	if err != nil {
		return err
	}

	var out bytes.Buffer
	err = t.Execute(&out, data)
	if err != nil {
		return err
	}

	src, err := format.Source(out.Bytes())
	if err != nil {
		return err
	}

	_, err = dest.Write(src)
	return err
}
