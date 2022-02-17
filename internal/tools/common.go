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

// Package tools provides helper functions used in scripts within the
// internal/tools module, as well as imports needed for a build with the
// "tools" build tag.
package tools // import "go.opentelemetry.io/otel/internal/tools"

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"

	"golang.org/x/mod/modfile"
)

// Repo represents a git repository.
type Repo string

// FindRepoRoot retrieves the root of the repository containing the current
// working directory. Beginning at the current working directory (dir), the
// algorithm checks if joining the ".git" suffix, such as "dir.get", is a
// valid file. Otherwise, it will continue checking the dir's parent directory
// until it reaches the repo root or returns an error if it cannot be found.
func FindRepoRoot() (Repo, error) {
	start, err := os.Getwd()
	if err != nil {
		return "", err
	}

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

		return Repo(dir), nil
	}
}

// FindModules returns all Go modules contained in Repo r.
func (r Repo) FindModules() ([]*modfile.File, error) {
	var results []*modfile.File
	err := filepath.Walk(string(r), func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			// Walk failed to walk into this directory. Stop walking and
			// signal this error.
			return walkErr
		}

		if !info.IsDir() {
			return nil
		}

		goMod := filepath.Join(path, "go.mod")
		f, err := os.Open(goMod)
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		if err != nil {
			return err
		}

		var b bytes.Buffer
		io.Copy(&b, f)
		if err = f.Close(); err != nil {
			return err
		}

		mFile, err := modfile.Parse(goMod, b.Bytes(), nil)
		if err != nil {
			return err
		}
		results = append(results, mFile)
		return nil
	})

	sort.SliceStable(results, func(i, j int) bool {
		return results[i].Syntax.Name < results[j].Syntax.Name
	})

	return results, err
}

func PrintModFiles(w io.Writer, mFiles []*modfile.File) error {
	tw := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
	if _, err := fmt.Fprintln(tw, "FILE PATH\tIMPORT PATH"); err != nil {
		return err
	}
	for _, m := range mFiles {
		if _, err := fmt.Fprintf(tw, "%s\t%s\n", m.Syntax.Name, m.Module.Mod.Path); err != nil {
			return err
		}
	}
	return tw.Flush()
}
