// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package propagation // import "go.opentelemetry.io/otel/propagation"

import (
	"os"
	"strings"
)

// EnvCarrier is a TextMapCarrier that uses the environment variables as a
// storage medium for propagated key-value pairs. The keys are uppercased
// before being used to access the environment variables.
// This is useful for propagating values that are set in the environment
// and need to be accessed by different processes or services.
// The keys are uppercased to avoid case sensitivity issues across different
// operating systems and environments.
type EnvCarrier struct {
	// SetEnvFunc is a function that sets the environment variable.
	// Usually, you want to set the environment variable for processes
	// that are spawned by the current process.
	// By default implementation, it does nothing.
	SetEnvFunc func(key, value string) error
}

var _ TextMapCarrier = EnvCarrier{}

// Get returns the value associated with the passed key.
// The key is uppercased before being used to access the environment variable.
func (EnvCarrier) Get(key string) string {
	k := strings.ToUpper(key)
	return os.Getenv(k)
}

// Set stores the key-value pair in the environment variable.
// The key is uppercased before being used to set the environment variable.
// If SetEnvFunc is not set, this method does nothing.
func (e EnvCarrier) Set(key, value string) {
	if e.SetEnvFunc == nil {
		return
	}
	k := strings.ToUpper(key)
	_ = e.SetEnvFunc(k, value)
}

// Keys lists the keys stored in this carrier.
// This returns all the keys in the environment variables.
func (EnvCarrier) Keys() []string {
	keys := make([]string, 0, len(os.Environ()))
	for _, kv := range os.Environ() {
		kvPair := strings.SplitN(kv, "=", 2)
		if len(kvPair) < 1 {
			continue
		}
		keys = append(keys, kvPair[0])
	}
	return keys
}
