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
	"log"
)

const path = "./testdata/valid-example.yaml"

func ExampleParseFile() {
	s, err := ParseFile(path)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("file format:", s.FileFormat)
	fmt.Println("schema URL:", s.SchemaURL)
	// Output:
	// file format: 1.1.0
	// schema URL: https://opentelemetry.io/schemas/1.1.0
}
