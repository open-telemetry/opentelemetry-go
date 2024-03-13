// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"

	ast10 "go.opentelemetry.io/otel/schema/v1.0/ast"
	types10 "go.opentelemetry.io/otel/schema/v1.0/types"
	ast11 "go.opentelemetry.io/otel/schema/v1.1/ast"
	types11 "go.opentelemetry.io/otel/schema/v1.1/types"
)

func TestParseSchemaFile(t *testing.T) {
	ts, err := ParseFile("testdata/valid-example.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, ts)
	assert.EqualValues(
		t, &ast11.Schema{
			FileFormat: "1.1.0",
			SchemaURL:  "https://opentelemetry.io/schemas/1.1.0",
			Versions: map[types11.TelemetryVersion]ast11.VersionDef{
				"1.0.0": {},

				"1.1.0": {
					All: ast10.Attributes{
						Changes: []ast10.AttributeChange{
							{
								RenameAttributes: &ast10.RenameAttributes{
									AttributeMap: ast10.AttributeMap{
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

					Resources: ast10.Attributes{
						Changes: []ast10.AttributeChange{
							{
								RenameAttributes: &ast10.RenameAttributes{
									AttributeMap: ast10.AttributeMap{
										"telemetry.auto.version": "telemetry.auto_instr.version",
									},
								},
							},
						},
					},

					Spans: ast10.Spans{
						Changes: []ast10.SpansChange{
							{
								RenameAttributes: &ast10.AttributeMapForSpans{
									AttributeMap: ast10.AttributeMap{
										"peer.service": "peer.service.name",
									},
									ApplyToSpans: []types10.SpanName{"HTTP GET"},
								},
							},
						},
					},

					SpanEvents: ast10.SpanEvents{
						Changes: []ast10.SpanEventsChange{
							{
								RenameEvents: &ast10.RenameSpanEvents{
									EventNameMap: map[string]string{
										"exception.stacktrace": "exception.stack_trace",
									},
								},
							},
							{
								RenameAttributes: &ast10.RenameSpanEventAttributes{
									ApplyToEvents: []types10.EventName{"exception.stack_trace"},
									AttributeMap: ast10.AttributeMap{
										"peer.service": "peer.service.name",
									},
								},
							},
						},
					},

					Logs: ast10.Logs{
						Changes: []ast10.LogsChange{
							{
								RenameAttributes: &ast10.RenameAttributes{
									AttributeMap: map[string]string{
										"process.executable_name": "process.executable.name",
									},
								},
							},
						},
					},

					Metrics: ast11.Metrics{
						Changes: []ast11.MetricsChange{
							{
								RenameAttributes: &ast10.AttributeMapForMetrics{
									AttributeMap: map[string]string{
										"http.status_code": "http.response_status_code",
									},
								},
							},
							{
								RenameMetrics: map[types10.MetricName]types10.MetricName{
									"container.cpu.usage.total":  "cpu.usage.total",
									"container.memory.usage.max": "memory.usage.max",
								},
							},
							{
								RenameAttributes: &ast10.AttributeMapForMetrics{
									ApplyToMetrics: []types10.MetricName{
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
								Split: &ast11.SplitMetric{
									ApplyToMetric: "system.paging.operations",
									ByAttribute:   "direction",
									MetricsFromAttributes: map[types10.MetricName]types11.AttributeValue{
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

func TestFailParseFileUnsupportedFileFormat(t *testing.T) {
	ts, err := ParseFile("testdata/unsupported-file-format.yaml")
	assert.ErrorContains(t, err, "unsupported schema file format minor version number")
	assert.Nil(t, ts)
}

func TestFailParseFileUnknownField(t *testing.T) {
	ts, err := ParseFile("testdata/unknown-field.yaml")
	assert.ErrorContains(t, err, "field Resources not found in type ast.VersionDef")
	assert.Nil(t, ts)
}
