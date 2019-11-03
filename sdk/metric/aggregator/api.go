// Copyright 2019, OpenTelemetry Authors
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

package aggregator // import "go.opentelemetry.io/otel/sdk/metric/aggregator"

import (
	"time"

	"go.opentelemetry.io/otel/api/core"
)

// TODO: Add Min() support to maxsumcount?  It's the same as
// Quantile(0) but cheap to compute like Max().

type (
	Sum interface {
		Sum() core.Number
	}

	Count interface {
		Count() core.Number
	}

	Max interface {
		Max() core.Number
	}

	Quantile interface {
		Quantile() core.Number
	}

	LastValue interface {
		LastValue() core.Number
		Timestamp() time.Time
	}

	MaxSumCount interface {
		Sum
		Count
		Max
	}

	Distribution interface {
		MaxSumCount
		Quantile
	}
)
