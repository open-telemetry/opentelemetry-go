// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVerifyReadmeIgnoresHiddenDotDirectories(t *testing.T) {
	t.Parallel()

	modDir := filepath.Join(t.TempDir(), ".hidden", "mod")
	if err := os.MkdirAll(modDir, 0o755); err != nil {
		t.Fatalf("MkdirAll(%q) error = %v", modDir, err)
	}

	goMod := filepath.Join(modDir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module example.com/hidden\n\ngo 1.25.0\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", goMod, err)
	}

	// Use a backslash-separated walk path so this exercises normalization on
	// every test runner, not only on Windows.
	walkPath := strings.ReplaceAll(filepath.ToSlash(goMod), "/", "\\")

	if err := verifyReadme(walkPath, dirEntry(t, goMod), nil); err != nil {
		t.Fatalf("verifyReadme(%q) error = %v, want nil", walkPath, err)
	}
}

func dirEntry(t *testing.T, path string) os.DirEntry {
	t.Helper()

	dir := filepath.Dir(path)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir(%q) error = %v", dir, err)
	}

	base := filepath.Base(path)
	for _, entry := range entries {
		if entry.Name() == base {
			return entry
		}
	}

	t.Fatalf("ReadDir(%q) missing %q", dir, base)
	return nil
}
