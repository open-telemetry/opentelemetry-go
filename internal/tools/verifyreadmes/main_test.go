// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestVerifyReadme(t *testing.T) {
	t.Run("propagates incoming walk error", func(t *testing.T) {
		wantErr := errors.New("walk error")

		got := verifyReadme("", nil, wantErr)
		if !errors.Is(got, wantErr) {
			t.Fatalf("errors.Is(got, wantErr) = false; got=%v want=%v", got, wantErr)
		}
	})

	t.Run("ignores non-go.mod files", func(t *testing.T) {
		dir := t.TempDir()
		p := filepath.Join(dir, "foo.txt")
		if err := os.WriteFile(p, []byte("x"), 0o600); err != nil {
			t.Fatalf("os.WriteFile(%q): %v", p, err)
		}
		info, err := os.Stat(p)
		if err != nil {
			t.Fatalf("os.Stat(%q): %v", p, err)
		}

		if err := verifyReadme(p, info, nil); err != nil {
			t.Fatalf("verifyReadme() = %v; want nil", err)
		}
	})

	t.Run("returns nil when README.md exists next to go.mod", func(t *testing.T) {
		dir := t.TempDir()
		modPath := filepath.Join(dir, "go.mod")
		readmePath := filepath.Join(dir, readmeFilename)

		if err := os.WriteFile(modPath, []byte("module example.com/test\n"), 0o600); err != nil {
			t.Fatalf("os.WriteFile(%q): %v", modPath, err)
		}
		if err := os.WriteFile(readmePath, []byte("# README\n"), 0o600); err != nil {
			t.Fatalf("os.WriteFile(%q): %v", readmePath, err)
		}
		info, err := os.Stat(modPath)
		if err != nil {
			t.Fatalf("os.Stat(%q): %v", modPath, err)
		}

		if err := verifyReadme(modPath, info, nil); err != nil {
			t.Fatalf("verifyReadme() = %v; want nil", err)
		}
	})

	t.Run("ignores excluded directories", func(t *testing.T) {
		dir := t.TempDir()
		modDir := filepath.Join(dir, "internal", "pkg")
		if err := os.MkdirAll(modDir, 0o755); err != nil {
			t.Fatalf("os.MkdirAll(%q): %v", modDir, err)
		}
		modPath := filepath.Join(modDir, "go.mod")
		if err := os.WriteFile(modPath, []byte("module example.com/internal/pkg\n"), 0o600); err != nil {
			t.Fatalf("os.WriteFile(%q): %v", modPath, err)
		}
		info, err := os.Stat(modPath)
		if err != nil {
			t.Fatalf("os.Stat(%q): %v", modPath, err)
		}

		if err := verifyReadme(modPath, info, nil); err != nil {
			t.Fatalf("verifyReadme() = %v; want nil", err)
		}
	})

	t.Run("returns expected error when README.md is missing", func(t *testing.T) {
		dir := t.TempDir()
		modPath := filepath.Join(dir, "go.mod")

		if err := os.WriteFile(modPath, []byte("module example.com/missingreadme\n"), 0o600); err != nil {
			t.Fatalf("os.WriteFile(%q): %v", modPath, err)
		}
		info, err := os.Stat(modPath)
		if err != nil {
			t.Fatalf("os.Stat(%q): %v", modPath, err)
		}

		err = verifyReadme(modPath, info, nil)
		if err == nil {
			t.Fatalf("verifyReadme() = nil; want error")
		}

		want := fmt.Sprintf("couldn't find %s for %q", readmeFilename, dir)
		if err.Error() != want {
			t.Fatalf("verifyReadme() error = %q; want %q", err.Error(), want)
		}
	})
}
