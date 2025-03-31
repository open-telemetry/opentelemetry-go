// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package internal provides internal functionality for the schema package.
package internal // import "go.opentelemetry.io/otel/schema/internal"

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
)

// CheckFileFormatField validates the file format field according to the rules here:
// https://github.com/open-telemetry/oteps/blob/main/text/0152-telemetry-schemas.md#schema-file-format-number
func CheckFileFormatField(fileFormat string, supportedFormatMajor, supportedFormatMinor int) error {
	// Verify that the version number in the file is a semver.
	fileFormatParsed, err := semver.StrictNewVersion(fileFormat)
	if err != nil {
		return fmt.Errorf(
			"invalid schema file format version number %q (expected semver): %w",
			fileFormat, err,
		)
	}

	if supportedFormatMajor < 0 {
		return errors.New("major version should be positive")
	}
	if supportedFormatMinor < 0 {
		return errors.New("major version should be positive")
	}

	// Check that the major version number in the file is the same as what we expect.
	if fileFormatParsed.Major() != uint64(
		supportedFormatMajor,
	) { // nolint:gosec // Version can't be negative (overflow checked).
		return fmt.Errorf(
			"this library cannot parse file formats with major version other than %v",
			supportedFormatMajor,
		)
	}

	// Check that the file minor version number is not greater than
	// what is requested supports.
	if fileFormatParsed.Minor() > uint64(
		supportedFormatMinor,
	) { // nolint:gosec // Version can't be negative (overflow checked).
		supportedFormatMajorMinor := strconv.Itoa(supportedFormatMajor) + "." +
			strconv.Itoa(supportedFormatMinor) // 1.0

		return fmt.Errorf(
			"unsupported schema file format minor version number, expected no newer than %v, got %v",
			supportedFormatMajorMinor+".x", fileFormat,
		)
	}

	// Patch, prerelease and metadata version number does not matter, so we don't check it.

	return nil
}

// CheckSchemaURL verifies that schemaURL is valid.
func CheckSchemaURL(schemaURL string) error {
	if strings.TrimSpace(schemaURL) == "" {
		return errors.New("schema_url field is missing")
	}

	if _, err := url.Parse(schemaURL); err != nil {
		return fmt.Errorf("invalid URL specified in schema_url field: %w", err)
	}
	return nil
}
