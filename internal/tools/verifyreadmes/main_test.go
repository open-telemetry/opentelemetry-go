// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestVerifyReadme(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	withReadme := filepath.Join(root, "with-readme")
	missingReadme := filepath.Join(root, "missing-readme")
	readmeOnly := filepath.Join(root, "readme-only")

	writeFile(t, filepath.Join(withReadme, "go.mod"))
	writeFile(t, filepath.Join(withReadme, readmeFilename))
	writeFile(t, filepath.Join(missingReadme, "go.mod"))
	writeFile(t, filepath.Join(readmeOnly, readmeFilename))

	tests := []struct {
		name    string
		path    string
		walkErr error
		wantErr string
	}{
		{
			name: "with readme",
			path: filepath.Join(withReadme, "go.mod"),
		},
		{
			name:    "missing readme",
			path:    filepath.Join(missingReadme, "go.mod"),
			wantErr: `couldn't find README.md for "` + missingReadme + `"`,
		},
		{
			name: "non go mod file",
			path: filepath.Join(readmeOnly, readmeFilename),
		},
		{
			name:    "walk error",
			path:    filepath.Join(withReadme, "go.mod"),
			walkErr: errors.New("walk failed"),
			wantErr: "walk failed",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			entry := dirEntry(t, test.path)

			err := verifyReadme(test.path, entry, test.walkErr)
			if test.wantErr == "" {
				if err != nil {
					t.Fatalf("verifyReadme(%q) error = %v, want nil", test.path, err)
				}
				return
			}

			if err == nil {
				t.Fatalf("verifyReadme(%q) error = nil, want %q", test.path, test.wantErr)
			}
			if err.Error() != test.wantErr {
				t.Fatalf("verifyReadme(%q) error = %q, want %q", test.path, err.Error(), test.wantErr)
			}
		})
	}
}

func TestVerifyReadmeIgnoresExcludedDirectories(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	paths := []string{
		filepath.Join(root, "internal", "pkg", "go.mod"),
		filepath.Join(root, "test", "pkg", "go.mod"),
		filepath.Join(root, "example", "pkg", "go.mod"),
		filepath.Join(root, ".hidden", "pkg", "go.mod"),
	}

	for _, path := range paths {
		writeFile(t, path)

		err := verifyReadme(path, dirEntry(t, path), nil)
		if err != nil {
			t.Fatalf("verifyReadme(%q) error = %v, want nil", path, err)
		}
	}
}

func dirEntry(t *testing.T, path string) fs.DirEntry {
	t.Helper()

	info, err := os.Lstat(path)
	if err != nil {
		t.Fatalf("Lstat(%q) error = %v", path, err)
	}
	return fs.FileInfoToDirEntry(info)
}

func writeFile(t *testing.T, path string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll(%q) error = %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, nil, 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", path, err)
	}
}
