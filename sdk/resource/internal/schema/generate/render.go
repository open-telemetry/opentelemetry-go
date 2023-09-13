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

package main

import (
	"bytes"
	"embed"
	"go/format"
	"os"
	"text/template"
)

// Template source.
const src = "templates/schema.go.tmpl"

//go:embed templates/schema.go.tmpl
var tmpl embed.FS

// render renders src from f in dest using data.
//
// All src will be rendered in the same file-tree with the same names (except
// for any ".tmpl" suffixes) as found in the Renderer's fs.FS.
func render(dest string, data any) error {
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

	f, err := os.Create(dest)
	if err != nil {
		return err
	}

	_, err = f.Write(src)
	if err != nil {
		return err
	}

	return f.Close()
}
