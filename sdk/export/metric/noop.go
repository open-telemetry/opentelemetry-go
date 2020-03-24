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

package metric

// NoopLabelExporter does no encoding at all.
type NoopLabelExporter struct{}

// Encode is a part of an implementation of the LabelEncoder
// interface. It returns an empty string.
func (NoopLabelExporter) Encode(LabelIterator) string {
	return ""
}

// ID is a part of an implementation of the LabelEncoder interface.
func (NoopLabelExporter) ID() int64 {
	// special reserved number for no op label encoder, see
	// labelExporterIDCounter variable docs
	return 1
}
