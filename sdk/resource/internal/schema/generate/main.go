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
	"embed"
	"flag"
	"log"

	sUtil "go.opentelemetry.io/otel/schema/v1.1"
	"go.opentelemetry.io/otel/schema/v1.1/ast"
	"go.opentelemetry.io/otel/sdk/resource/internal/schema/generate/cmd"
)

const schema = "schema/schema.yaml"

//go:embed schema/schema.yaml
var schemaFs embed.FS

func load() (*ast.Schema, error) {
	f, err := schemaFs.Open(schema)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return sUtil.Parse(f)
}

func main() {
	flag.Parse()
	dest := flag.Arg(0)
	if dest == "" {
		log.Fatalln("empty desination")
	}

	s, err := load()
	if err != nil {
		log.Fatalf("failed to load schema: %s", err)
	}

	if err := cmd.Run(dest, s); err != nil {
		log.Fatal(err)
	}
}
