// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import "strings"

// TemporalityPreference defines the user's desired temporality for metrics instruments.
type TemporalityPreference string

const (
	// TemporalityPreferenceDefault indicates the SDK's default temporality should be used.
	TemporalityPreferenceDefault TemporalityPreference = ""

	// TemporalityPreferenceCumulative indicates a cumulative temporality should be used.
	TemporalityPreferenceCumulative TemporalityPreference = "cumulative"

	// TemporalityPreferenceDelta indicates a delta temporality should be used.
	TemporalityPreferenceDelta TemporalityPreference = "delta"

	// TemporalityPreferenceLowMemory indicates temporality preference that optimizes memory use.
	TemporalityPreferenceLowMemory TemporalityPreference = "lowmemory"
)

// IsValid returns true whether the preference is a valid string constant.
func (t TemporalityPreference) IsValid() bool {
	switch t.lowercase() {
	case TemporalityPreferenceCumulative,
		TemporalityPreferenceDelta,
		TemporalityPreferenceLowMemory,
		TemporalityPreferenceDefault:
		return true
	default:
		return false
	}
}

func (t TemporalityPreference) lowercase() TemporalityPreference {
	return TemporalityPreference(strings.ToLower(string(t)))
}

// TemporalitySelectorForPreference will return a TemporalitySelector for the given preference.
func TemporalitySelectorForPreference(t TemporalityPreference) TemporalitySelector {
	switch t.lowercase() {
	case TemporalityPreferenceCumulative:
		return CumulativeTemporalitySelector
	case TemporalityPreferenceDelta:
		return DeltaTemporalitySelector
	case TemporalityPreferenceLowMemory:
		return LowMemoryTemporalitySelector
	default:
		return DefaultTemporalitySelector
	}
}
