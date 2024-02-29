// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resource // import "go.opentelemetry.io/otel/sdk/resource"

var (
	SetDefaultOSProviders           = setDefaultOSProviders
	SetOSProviders                  = setOSProviders
	SetDefaultRuntimeProviders      = setDefaultRuntimeProviders
	SetRuntimeProviders             = setRuntimeProviders
	SetDefaultUserProviders         = setDefaultUserProviders
	SetUserProviders                = setUserProviders
	SetDefaultOSDescriptionProvider = setDefaultOSDescriptionProvider
	SetOSDescriptionProvider        = setOSDescriptionProvider
	SetDefaultContainerProviders    = setDefaultContainerProviders
	SetContainerProviders           = setContainerProviders
)

var (
	CommandArgs = commandArgs
	RuntimeName = runtimeName
	RuntimeOS   = runtimeOS
	RuntimeArch = runtimeArch
)

var MapRuntimeOSToSemconvOSType = mapRuntimeOSToSemconvOSType
