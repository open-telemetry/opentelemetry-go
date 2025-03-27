// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package schema provides functionality and types for OpenTelemetry schemas.
package schema // import "go.opentelemetry.io/otel/schema/v1.1"

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"

	"go.opentelemetry.io/otel/schema/internal"
	"go.opentelemetry.io/otel/schema/v1.1/ast"
)

// Major file version number that this library supports.
const supportedFormatMajor = 1

// Maximum minor version number that this library supports.
const supportedFormatMinor = 1

// ParseFile a schema file. schemaFilePath is the file path.
func ParseFile(schemaFilePath string) (*ast.Schema, error) {
	file, err := os.Open(schemaFilePath)
	if err != nil {
		return nil, err
	}
	return Parse(file)
}

// Parse a schema file. schemaFileContent is the readable content of the schema file.
func Parse(schemaFileContent io.Reader) (*ast.Schema, error) {
	var ts ast.Schema
	d := yaml.NewDecoder(schemaFileContent)
	d.KnownFields(true)
	err := d.Decode(&ts)
	if err != nil {
		return nil, err
	}

	err = internal.CheckFileFormatField(ts.FileFormat, supportedFormatMajor, supportedFormatMinor)
	if err != nil {
		return nil, err
	}

	err = internal.CheckSchemaURL(ts.SchemaURL)
	if err != nil {
		return nil, err
	}

	return &ts, nil
}
