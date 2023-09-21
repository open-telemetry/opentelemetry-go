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

package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSchemaFile(t *testing.T) {
	ts, err := ParseFile("testdata/valid-example.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, ts)
	assert.EqualValues(
		t, &Schema{
			FileFormat: "1.1.0",
			SchemaURL:  "https://opentelemetry.io/schemas/1.1.0",
			Versions: map[string]Changeset{
				"1.0.0": {},

				"1.1.0": {
					All: All{
						Changes: []AllChange{
							{
								RenameAttributes: &RenameAttributes{
									AttributeMap: map[string]string{
										"k8s.cluster.name":     "kubernetes.cluster.name",
										"k8s.namespace.name":   "kubernetes.namespace.name",
										"k8s.node.name":        "kubernetes.node.name",
										"k8s.node.uid":         "kubernetes.node.uid",
										"k8s.pod.name":         "kubernetes.pod.name",
										"k8s.pod.uid":          "kubernetes.pod.uid",
										"k8s.container.name":   "kubernetes.container.name",
										"k8s.replicaset.name":  "kubernetes.replicaset.name",
										"k8s.replicaset.uid":   "kubernetes.replicaset.uid",
										"k8s.cronjob.name":     "kubernetes.cronjob.name",
										"k8s.cronjob.uid":      "kubernetes.cronjob.uid",
										"k8s.job.name":         "kubernetes.job.name",
										"k8s.job.uid":          "kubernetes.job.uid",
										"k8s.statefulset.name": "kubernetes.statefulset.name",
										"k8s.statefulset.uid":  "kubernetes.statefulset.uid",
										"k8s.daemonset.name":   "kubernetes.daemonset.name",
										"k8s.daemonset.uid":    "kubernetes.daemonset.uid",
										"k8s.deployment.name":  "kubernetes.deployment.name",
										"k8s.deployment.uid":   "kubernetes.deployment.uid",
										"service.namespace":    "service.namespace.name",
									},
								},
							},
						},
					},

					Resources: Resources{
						Changes: []ResourcesChange{
							{
								RenameAttributes: &RenameAttributes{
									AttributeMap: map[string]string{
										"telemetry.auto.version": "telemetry.auto_instr.version",
									},
								},
							},
						},
					},

					Spans: Spans{
						Changes: []SpansChange{
							{
								RenameAttributes: &RenameSpansAttributes{
									AttributeMap: map[string]string{
										"peer.service": "peer.service.name",
									},
									ApplyToSpans: []string{"HTTP GET"},
								},
							},
						},
					},

					SpanEvents: SpanEvents{
						Changes: []SpanEventsChange{
							{
								RenameEvents: &RenameSpanEvents{
									EventNameMap: map[string]string{
										"exception.stacktrace": "exception.stack_trace",
									},
								},
							},
							{
								RenameAttributes: &RenameSpanEventsAttributes{
									ApplyToEvents: []string{"exception.stack_trace"},
									AttributeMap: map[string]string{
										"peer.service": "peer.service.name",
									},
								},
							},
						},
					},

					Logs: Logs{
						Changes: []LogsChange{
							{
								RenameAttributes: &RenameAttributes{
									AttributeMap: map[string]string{
										"process.executable_name": "process.executable.name",
									},
								},
							},
						},
					},

					Metrics: Metrics{
						Changes: []MetricsChange{
							{
								RenameAttributes: &RenameMetricsAttributes{
									AttributeMap: map[string]string{
										"http.status_code": "http.response_status_code",
									},
								},
							},
							{
								RenameMetrics: map[string]string{
									"container.cpu.usage.total":  "cpu.usage.total",
									"container.memory.usage.max": "memory.usage.max",
								},
							},
							{
								RenameAttributes: &RenameMetricsAttributes{
									ApplyToMetrics: []string{
										"system.cpu.utilization",
										"system.memory.usage",
										"system.memory.utilization",
										"system.paging.usage",
									},
									AttributeMap: map[string]string{
										"status": "state",
									},
								},
							},
							{
								Split: &SplitMetric{
									ApplyToMetric: "system.paging.operations",
									ByAttribute:   "direction",
									MetricsFromAttributes: map[string]any{
										"system.paging.operations.in":  "in",
										"system.paging.operations.out": "out",
									},
								},
							},
						},
					},
				},
			},
		}, ts,
	)
}

func TestFailParseFileUnknownField(t *testing.T) {
	ts, err := ParseFile("testdata/unknown-field.yaml")
	assert.ErrorContains(t, err, "field Resources not found in type schema.Changeset")
	assert.Nil(t, ts)
}
