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

package metric // import "go.opentelemetry.io/otel/metric"

import (
	"go.opentelemetry.io/otel/label"
)

//go:generate stringer -type=ViewLabelConfig

// ViewLabelConfig is the option to control how to track the set of labels
// for a given view.
// See available values for more details.
type ViewLabelConfig uint8

const (
	// Drop all label keys, and aggregate all measurements from the metric
	// instrument together regardless of their label sets.
	//
	// Note: the DropAll is the same as LabelKeys with empty keys.
	DropAll ViewLabelConfig = iota

	// Record all label keys. This is the default option without explicitly
	// registering a view.
	Ungroup

	// Specify a set of label keys to track at view creation time, and drop
	// other keys from the label sets of recorded measurements before
	// aggregating.
	//
	// Note: some label values may be undefined if the measurements don't
	// provide all label values.
	LabelKeys
)

// View supports configuring non-default aggregation behaviors on the level
// of an individual synchronous instrument.
type View struct {
	inst               SyncImpl
	labelConfig        ViewLabelConfig
	labelKeys          []label.Key
	aggregatorSelector AggregatorSelector
}

// SyncImpl returns the associated synchronous instrument interface for a given view.
func (v View) SyncImpl() SyncImpl {
	return v.inst
}

// LabelConfig returns the label tracking option for a given view.
func (v View) LabelConfig() ViewLabelConfig {
	return v.labelConfig
}

// LabelKeys returns the pre-configured label keys for a given view.
func (v View) LabelKeys() []label.Key {
	return v.labelKeys
}

// AggregatorFactory returns the AggregatorFactory instance for a given view.
func (v View) AggregatorFactory() AggregatorSelector {
	return v.aggregatorSelector
}

// NewView returns a view object for the given parameters.
func NewView(inst SyncImpl, labelConfig ViewLabelConfig, labelKeys []label.Key, selector AggregatorSelector) View {
	return View{
		inst:               inst,
		labelConfig:        labelConfig,
		labelKeys:          labelKeys,
		aggregatorSelector: selector,
	}
}
