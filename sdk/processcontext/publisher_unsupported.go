// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build !linux

package processcontext

import (
	"errors"

	"go.opentelemetry.io/otel/sdk/resource"
)

var errUnsupported = errors.New("processcontext: not supported on this platform")

type publisher struct{}

func newPublisher(_ *resource.Resource, _ config) (*publisher, error) {
	return nil, errUnsupported
}

func (p *publisher) update(_ *resource.Resource) error { return errUnsupported }

func (p *publisher) shutdown() {}
