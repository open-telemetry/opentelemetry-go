// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

/*
Package profile provides an experimental OpenTelemetry Profiles SDK for Go.

The Profiles signal is currently Alpha in the OpenTelemetry specification
(https://opentelemetry.io/docs/specs/otel/profiles/). This package may change
in backwards-incompatible ways between releases.

The entry point is [NewProfilerProvider], which coordinates [Producer]
implementations (such as [RuntimeProducer] for Go runtime heap, allocation,
goroutine, mutex, block, and CPU profiles) and [Reader] implementations
([PeriodicReader] and [ManualReader]) that collect and export
[go.opentelemetry.io/otel/sdk/profile/profiledata.ResourceProfiles] through an
[Exporter].
*/
package profile
