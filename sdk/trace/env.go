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

package trace // import "go.opentelemetry.io/otel/sdk/trace"

import (
	"os"
	"strconv"

	"go.opentelemetry.io/otel/internal/global"
)

// Environment variable names
const (
	// EnvBatchSpanProcessorScheduleDelay
	// Delay interval between two consecutive exports.
	// i.e. 5000
	EnvBatchSpanProcessorScheduleDelay = "OTEL_BSP_SCHEDULE_DELAY"
	// EnvBatchSpanProcessorExportTimeout
	// Maximum allowed time to export data.
	// i.e. 3000
	EnvBatchSpanProcessorExportTimeout = "OTEL_BSP_EXPORT_TIMEOUT"
	// EnvBatchSpanProcessorMaxQueueSize
	// Maximum queue size
	// i.e. 2048
	EnvBatchSpanProcessorMaxQueueSize = "OTEL_BSP_MAX_QUEUE_SIZE"
	// EnvBatchSpanProcessorMaxExportBatchSize
	// Maximum batch size
	// Note: Must be less than or equal to EnvBatchSpanProcessorMaxQueueSize
	// i.e. 512
	EnvBatchSpanProcessorMaxExportBatchSize = "OTEL_BSP_MAX_EXPORT_BATCH_SIZE"
)

// intEnvOr returns an env variable's numeric value if it is exists (and valid) or the default if not
func intEnvOr(key string, defaultValue int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		global.Info("Got invalid value, number value expected.", key, value)
		return defaultValue
	}

	return intValue
}
