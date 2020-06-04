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

/*
* The oterror package provides unified error interactions in OpenTelemetry.
* This includes providing standardized errors common to OpenTelemetry (APIs,
* SDKs, and exporters). Additionally it provides an API for unified error
* handling in OpenTelemetry.
*
* The unified error handling interface is used for any error that
* OpenTelemetry component are not able to remediate on their own, instead
* handling them in a uniform and user-defined way.
 */

package oterror
