// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package verifyreadmes is used to verify that all go modules in the repository
// have a README.md file.
package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestVerifyReadme(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "test-readme")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testCases := []struct {
		name  string
		setup func(t *testing.T, dir string) (string, os.FileInfo)
		want  bool
	}{
		{
			name: "README exists",
			setup: func(t *testing.T, dir string) (string, os.FileInfo) {
				path := filepath.Join(dir, "with-readme")
				if err := os.Mkdir(path, 0755); err != nil {
					t.Fatal(err)
				}
				goModPath := filepath.Join(path, "go.mod")
				if _, err := os.Create(goModPath); err != nil {
					t.Fatal(err)
				}
				if _, err := os.Create(filepath.Join(path, "README.md")); err != nil {
					t.Fatal(err)
				}
				info, err := os.Stat(goModPath)
				if err != nil {
					t.Fatal(err)
				}
				return goModPath, info
			},
			want: false,
		},
		{
			name: "Excluded directory",
			setup: func(t *testing.T, dir string) (string, os.FileInfo) {
				path := filepath.Join(dir, "internal", "excluded")
				if err := os.MkdirAll(path, 0755); err != nil {
					t.Fatal(err)
				}
				goModPath := filepath.Join(path, "go.mod")
				if _, err := os.Create(goModPath); err != nil {
					t.Fatal(err)
				}
				info, err := os.Stat(goModPath)
				if err != nil {
					t.Fatal(err)
				}
				return goModPath, info
			},
			want: false,
		},
		{
			name: "No go.mod",
			setup: func(t *testing.T, dir string) (string, os.FileInfo) {
				path := filepath.Join(dir, "no-go-mod")
				if err := os.Mkdir(path, 0755); err != nil {
					t.Fatal(err)
				}
				readmePath := filepath.Join(path, "README.md")
				if _, err := os.Create(readmePath); err != nil {
					t.Fatal(err)
				}
				info, err := os.Stat(readmePath)
				if err != nil {
					t.Fatal(err)
				}
				return readmePath, info
			},
			want: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			path, info := tc.setup(t, tmpDir)
			got := verifyReadme(path, info, nil)
			if (err != nil) != tc.want {
				t.Errorf("Expected error: %v, got: %v", tc.want, got)
			}
		})
	}
}
