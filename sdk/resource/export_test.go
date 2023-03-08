// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

var (
	MapRuntimeOSToSemconvOSType = mapRuntimeOSToSemconvOSType
)
