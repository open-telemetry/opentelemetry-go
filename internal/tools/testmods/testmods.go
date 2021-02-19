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

// The testmods tool runs go test for all Go modules within this repository.
// Any arguments passed will be forwarded to the go test command.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var excludeMods = map[string]bool{
	"go.opentelemetry.io/otel/internal/tools": true,
}

func gitRoot(start string) (string, error) {
	dir := start
	for {
		_, err := os.Stat(filepath.Join(dir, ".git"))
		if errors.Is(err, os.ErrNotExist) {
			dir = filepath.Dir(dir)
			// From https://golang.org/pkg/path/filepath/#Dir:
			// The returned path does not end in a separator unless it is the root directory.
			if strings.HasSuffix(dir, string(filepath.Separator)) {
				return "", fmt.Errorf("unable to find git repository enclosing working dir %s", start)
			}
			continue
		}

		if err != nil {
			return "", err
		}

		return dir, nil
	}
}

type mod struct {
	filePath   string
	importPath string
}

func listModules(root string) ([]mod, error) {
	var results []mod
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}

		_, err = os.Stat(filepath.Join(path, "go.mod"))
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		if err != nil {
			return err
		}

		cmd := exec.Command("go", "mod", "edit", "-json")
		cmd.Dir = path
		out, err := cmd.Output()
		if err != nil {
			return err
		}

		var result struct {
			Module struct {
				Path string
			}
		}
		err = json.Unmarshal(out, &result)
		if err != nil {
			return err
		}

		results = append(results, mod{
			filePath:   path,
			importPath: result.Module.Path,
		})
		return nil
	})

	return results, err
}

func testmod(m mod, arg ...string) error {
	cmd := exec.Command("go", append([]string{"test"}, arg...)...)
	cmd.Dir = m.filePath
	out, err := cmd.CombinedOutput()
	fmt.Print(string(out))
	return err
}

func main() {
	start, err := os.Getwd()
	if err != nil {
		log.Fatalf("unable to get current directory: %v", err)
	}

	root, err := gitRoot(start)
	if err != nil {
		log.Fatalf("unable to find git root: %v", err)
	}

	mods, err := listModules(root)
	if err != nil {
		log.Fatalf("unable to find all Go modules: %v", err)
	}

	// Zero-length slice with the same underlying data.
	filtered := mods[:0]
	for _, m := range mods {
		if _, exclude := excludeMods[m.importPath]; !exclude {
			filtered = append(filtered, m)
		}
	}

	for _, m := range filtered {
		if err := testmod(m, os.Args[1:]...); err != nil {
			log.Fatalf("unable to test %s: %v\n", m.importPath, err)
		}
	}
}
