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
	"go/format"
	"io"
	"text/template"
)

// Template source.
const src = "transforms.go.tmpl"

//go:embed transforms.go.tmpl
var tmpl embed.FS

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
