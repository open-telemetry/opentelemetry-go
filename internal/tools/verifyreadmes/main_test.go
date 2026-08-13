// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

// dirEntry returns the [fs.DirEntry] for path the same way [filepath.WalkDir]
// builds the entry it passes to its callback.
func dirEntry(t *testing.T, path string) fs.DirEntry {
	t.Helper()

	info, err := os.Lstat(path)
	if err != nil {
		t.Fatalf("os.Lstat(%q): %v", path, err)
	}

	return fs.FileInfoToDirEntry(info)
}

// writeFile writes content to path, creating any missing parent directories.
func writeFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		t.Fatalf("os.MkdirAll(%q): %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("os.WriteFile(%q): %v", path, err)
	}
}

// TestVerifyReadme uses paths relative to a temporary working directory.
// verifyReadme matches excludedDirs against the whole path
func TestVerifyReadme(t *testing.T) {
	t.Run("propagates incoming walk error", func(t *testing.T) {
		wantErr := errors.New("walk error")

		got := verifyReadme("", nil, wantErr)
		if !errors.Is(got, wantErr) {
			t.Fatalf("verifyReadme() = %v; want %v", got, wantErr)
		}
	})

	t.Run("ignores non-go.mod files", func(t *testing.T) {
		t.Chdir(t.TempDir())

		path := "foo.txt"
		writeFile(t, path, "x")

		if err := verifyReadme(path, dirEntry(t, path), nil); err != nil {
			t.Fatalf("verifyReadme() = %v; want nil", err)
		}
	})

	t.Run("ignores directories named go.mod", func(t *testing.T) {
		t.Chdir(t.TempDir())

		path := "go.mod"
		if err := os.Mkdir(path, 0o750); err != nil {
			t.Fatalf("os.Mkdir(%q): %v", path, err)
		}

		if err := verifyReadme(path, dirEntry(t, path), nil); err != nil {
			t.Fatalf("verifyReadme() = %v; want nil", err)
		}
	})

	t.Run("ignores excluded directories", func(t *testing.T) {
		t.Chdir(t.TempDir())

		// "internal" is the only part of this path matching excludedDirs.
		modPath := filepath.Join("internal", "go.mod")
		writeFile(t, modPath, "module example.com/internal\n")

		if err := verifyReadme(modPath, dirEntry(t, modPath), nil); err != nil {
			t.Fatalf("verifyReadme() = %v; want nil", err)
		}
	})

	t.Run("returns nil when README.md exists next to go.mod", func(t *testing.T) {
		t.Chdir(t.TempDir())

		modDir := "mod"
		modPath := filepath.Join(modDir, "go.mod")
		writeFile(t, modPath, "module example.com/mod\n")
		writeFile(t, filepath.Join(modDir, readmeFilename), "# README\n")

		if err := verifyReadme(modPath, dirEntry(t, modPath), nil); err != nil {
			t.Fatalf("verifyReadme() = %v; want nil", err)
		}
	})

	t.Run("returns expected error when README.md is missing", func(t *testing.T) {
		t.Chdir(t.TempDir())

		modDir := "mod"
		modPath := filepath.Join(modDir, "go.mod")
		writeFile(t, modPath, "module example.com/mod\n")

		err := verifyReadme(modPath, dirEntry(t, modPath), nil)
		if err == nil {
			t.Fatalf("verifyReadme() = nil; want error")
		}

		want := fmt.Sprintf("couldn't find %s for %q", readmeFilename, modDir)
		if err.Error() != want {
			t.Fatalf("verifyReadme() error = %q; want %q", err.Error(), want)
		}
	})
}
