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

package asyncint64 // import "go.opentelemetry.io/otel/metric/instrument/asyncint64"

import "go.opentelemetry.io/otel/metric/instrument"

// Counter is an instrument used to asynchronously record increasing int64
// measurements once per a measurement collection cycle. The Observe method is
// used to record the measured state of the instrument when it is called.
// Implementations will assume the observed value to be the cumulative sum of
// the count.
//
// Warning: methods may be added to this interface in minor releases.
type Counter interface{ instrument.Int64Observer }

// UpDownCounter is an instrument used to asynchronously record int64
// measurements once per a measurement collection cycle. The Observe method is
// used to record the measured state of the instrument when it is called.
// Implementations will assume the observed value to be the cumulative sum of
// the count.
//
// Warning: methods may be added to this interface in minor releases.
type UpDownCounter interface{ instrument.Int64Observer }

// Gauge is an instrument used to asynchronously record instantaneous int64
// measurements once per a measurement collection cycle.
//
// Warning: methods may be added to this interface in minor releases.
type Gauge interface{ instrument.Int64Observer }
