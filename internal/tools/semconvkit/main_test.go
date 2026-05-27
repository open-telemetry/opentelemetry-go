// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewMigrationIgnoresUnexportedDeclarations(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	prev := filepath.Join(root, "v1.40.0")
	cur := filepath.Join(root, "v1.41.0")

	writeGoFile(t, prev, "semconv", `package semconv

func ErrorType(err error) {}
func DeploymentEnvironmentName(val string) {}
`)
	writeGoFile(t, cur, "semconv", `package semconv

func ErrorType(err error) {}
func errorType(err error) {}
`)

	m, err := newMigration(cur, prev)
	if err != nil {
		t.Fatalf("newMigration() error = %v", err)
	}

	if len(m.Renames) != 0 {
		t.Fatalf("newMigration() renames = %#v, want none", m.Renames)
	}

	wantRemoval := []string{"DeploymentEnvironmentName"}
	if len(m.Removals) != len(wantRemoval) || m.Removals[0] != wantRemoval[0] {
		t.Fatalf("newMigration() removals = %#v, want %#v", m.Removals, wantRemoval)
	}
}

// TestNoInvalidObservableHistogramTypes asserts that the semconv code
// generator never emits types containing "HistogramObservable" or
// "ObservableObservable", which would indicate the histogram-exclusion
// rule in metric.go.j2 was accidentally removed or bypassed.
func TestNoInvalidObservableHistogramTypes(t *testing.T) {
	t.Parallel()

	semconvDir := filepath.Join("..", "..", "..", "semconv")

	forbidden := [][]byte{
		[]byte("HistogramObservable"),
		[]byte("ObservableObservable"),
	}

	err := filepath.WalkDir(semconvDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Base(path) != "metric.go" || !strings.HasSuffix(filepath.Base(filepath.Dir(path)), "conv") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		for _, term := range forbidden {
			if bytes.Contains(data, term) {
				t.Errorf("%s: contains forbidden substring %q", path, term)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("WalkDir(%q) error = %v", semconvDir, err)
	}
}

func writeGoFile(t *testing.T, dir, pkg, src string) {
	t.Helper()

	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll(%q) error = %v", dir, err)
	}

	path := filepath.Join(dir, "decls.go")
	if err := os.WriteFile(path, []byte(src), 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", path, err)
	}
}
