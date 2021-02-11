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

// The crosslink tool generates and maintains replace directives in all
// the go.mod files within this repository. Some directives are superfluous
// (e.g. because the replaced module doesn't occur in the dependency tree),
// but we generate them anyway for the sake of consistency (#1529 tracks
// pruning this to a mininal set).
//
// In particular, we generate a replace directive from each module to itself
// (i.e., the target path "./"). This is actually necessary in the presence of
// cyclic dependencies between modules.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/tabwriter"
)

type repo string

func findRepoRoot() (repo, error) {
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

		return repo(dir), nil
	}
}

type mod struct {
	filePath   string
	importPath string
}

func (r repo) findModules() (mods, error) {
	var results []mod
	err := filepath.Walk(string(r), func(path string, info os.FileInfo, err error) error {
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

type mods []mod

func (m mods) print(w io.Writer) error {
	tw := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
	if _, err := fmt.Fprintln(tw, "FILE PATH\tIMPORT PATH"); err != nil {
		return err
	}
	for _, m := range m {
		if _, err := fmt.Fprintf(tw, "%s\t%s\n", m.filePath, m.importPath); err != nil {
			return err
		}
	}
	return tw.Flush()
}

func (m mods) crossLink() error {
	for _, from := range m {
		args := []string{"mod", "edit"}

		for _, to := range m {
			localPath, err := filepath.Rel(from.filePath, to.filePath)
			if err != nil {
				return err
			}
			if localPath == "." || localPath == ".." {
				localPath += "/"
			} else if !strings.HasPrefix(localPath, "..") {
				localPath = "./" + localPath
			}
			args = append(args, "-replace", to.importPath+"="+localPath)
		}

		cmd := exec.Command("go", args...)
		cmd.Dir = from.filePath
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Println(string(out))
			return err
		}
	}
	return nil
}

func main() {
	repoRoot, err := findRepoRoot()
	if err != nil {
		log.Fatalf("unable to find repo root: %v", err)
	}

	mods, err := repoRoot.findModules()
	if err != nil {
		log.Fatalf("unable to list modules: %v", err)
	}

	if err := mods.print(os.Stdout); err != nil {
		log.Fatalf("unable to print modules: %v", err)
	}

	if err := mods.crossLink(); err != nil {
		log.Fatalf("unable to crosslink: %v", err)
	}
}
