// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package schema

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/schema/v1.0/ast"
	"go.opentelemetry.io/otel/schema/v1.0/types"
)

func TestParseSchemaFile(t *testing.T) {
	ts, err := ParseFile("testdata/valid-example.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, ts)
	assert.EqualValues(
		t, &ast.Schema{
			FileFormat: "1.0.0",
			SchemaURL:  "https://opentelemetry.io/schemas/1.1.0",
			Versions: map[types.TelemetryVersion]ast.VersionDef{
				"1.0.0": {},

				"1.1.0": {
					All: ast.Attributes{
						Changes: []ast.AttributeChange{
							{
								RenameAttributes: &ast.RenameAttributes{
									AttributeMap: ast.AttributeMap{
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

					Resources: ast.Attributes{
						Changes: []ast.AttributeChange{
							{
								RenameAttributes: &ast.RenameAttributes{
									AttributeMap: ast.AttributeMap{
										"telemetry.auto.version": "telemetry.auto_instr.version",
									},
								},
							},
						},
					},

					Spans: ast.Spans{
						Changes: []ast.SpansChange{
							{
								RenameAttributes: &ast.AttributeMapForSpans{
									AttributeMap: ast.AttributeMap{
										"peer.service": "peer.service.name",
									},
									ApplyToSpans: []types.SpanName{"HTTP GET"},
								},
							},
						},
					},

					SpanEvents: ast.SpanEvents{
						Changes: []ast.SpanEventsChange{
							{
								RenameEvents: &ast.RenameSpanEvents{
									EventNameMap: map[string]string{
										"exception.stacktrace": "exception.stack_trace",
									},
								},
							},
							{
								RenameAttributes: &ast.RenameSpanEventAttributes{
									ApplyToEvents: []types.EventName{"exception.stack_trace"},
									AttributeMap: ast.AttributeMap{
										"peer.service": "peer.service.name",
									},
								},
							},
						},
					},

					Logs: ast.Logs{
						Changes: []ast.LogsChange{
							{
								RenameAttributes: &ast.RenameAttributes{
									AttributeMap: map[string]string{
										"process.executable_name": "process.executable.name",
									},
								},
							},
						},
					},

					Metrics: ast.Metrics{
						Changes: []ast.MetricsChange{
							{
								RenameAttributes: &ast.AttributeMapForMetrics{
									AttributeMap: map[string]string{
										"http.status_code": "http.response_status_code",
									},
								},
							},
							{
								RenameMetrics: map[types.MetricName]types.MetricName{
									"container.cpu.usage.total":  "cpu.usage.total",
									"container.memory.usage.max": "memory.usage.max",
								},
							},
							{
								RenameAttributes: &ast.AttributeMapForMetrics{
									ApplyToMetrics: []types.MetricName{
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
						},
					},
				},
			},
		}, ts,
	)
}

func TestFailParseSchemaFile(t *testing.T) {
	ts, err := ParseFile("testdata/unsupported-file-format.yaml")
	assert.Error(t, err)
	assert.Nil(t, ts)

	ts, err = ParseFile("testdata/invalid-schema-url.yaml")
	assert.Error(t, err)
	assert.Nil(t, ts)

	ts, err = ParseFile("testdata/unknown-field.yaml")
	assert.ErrorContains(t, err, "field Resources not found in type ast.VersionDef")
	assert.Nil(t, ts)
}

func TestFailParseSchema(t *testing.T) {
	_, err := Parse(bytes.NewReader([]byte("")))
	assert.Error(t, err)

	_, err = Parse(bytes.NewReader([]byte("invalid yaml")))
	assert.Error(t, err)

	_, err = Parse(bytes.NewReader([]byte("file_format: 1.0.0")))
	assert.Error(t, err)
}
