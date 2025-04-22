// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated from semantic convention specification. DO NOT EDIT.

package semconv // import "go.opentelemetry.io/otel/semconv/v1.31.0/messaging"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// ErrorTypeAttr is an attribute conforming to the error.type semantic
// conventions. It represents the describes a class of error the operation ended
// with.
type ErrorTypeAttr string

var (
	// ErrorTypeOther is a fallback error value to be used when the instrumentation
	// doesn't define a custom value.
	ErrorTypeOther ErrorTypeAttr = "_OTHER"
)

// OperationTypeAttr is an attribute conforming to the messaging.operation.type
// semantic conventions. It represents a string identifying the type of the
// messaging operation.
type OperationTypeAttr string

var (
	// OperationTypeCreate is a message is created. "Create" spans always refer to a
	// single message and are used to provide a unique creation context for messages
	// in batch sending scenarios.
	OperationTypeCreate OperationTypeAttr = "create"
	// OperationTypeSend is the one or more messages are provided for sending to an
	// intermediary. If a single message is sent, the context of the "Send" span can
	// be used as the creation context and no "Create" span needs to be created.
	OperationTypeSend OperationTypeAttr = "send"
	// OperationTypeReceive is the one or more messages are requested by a consumer.
	// This operation refers to pull-based scenarios, where consumers explicitly
	// call methods of messaging SDKs to receive messages.
	OperationTypeReceive OperationTypeAttr = "receive"
	// OperationTypeProcess is the one or more messages are processed by a consumer.
	OperationTypeProcess OperationTypeAttr = "process"
	// OperationTypeSettle is the one or more messages are settled.
	OperationTypeSettle OperationTypeAttr = "settle"
)

// SystemAttr is an attribute conforming to the messaging.system semantic
// conventions. It represents the messaging system as identified by the client
// instrumentation.
type SystemAttr string

var (
	// SystemActivemq is the apache ActiveMQ.
	SystemActivemq SystemAttr = "activemq"
	// SystemAWSSqs is the amazon Simple Queue Service (SQS).
	SystemAWSSqs SystemAttr = "aws_sqs"
	// SystemEventgrid is the azure Event Grid.
	SystemEventgrid SystemAttr = "eventgrid"
	// SystemEventhubs is the azure Event Hubs.
	SystemEventhubs SystemAttr = "eventhubs"
	// SystemServicebus is the azure Service Bus.
	SystemServicebus SystemAttr = "servicebus"
	// SystemGCPPubsub is the google Cloud Pub/Sub.
	SystemGCPPubsub SystemAttr = "gcp_pubsub"
	// SystemJms is the java Message Service.
	SystemJms SystemAttr = "jms"
	// SystemKafka is the apache Kafka.
	SystemKafka SystemAttr = "kafka"
	// SystemRabbitmq is the rabbitMQ.
	SystemRabbitmq SystemAttr = "rabbitmq"
	// SystemRocketmq is the apache RocketMQ.
	SystemRocketmq SystemAttr = "rocketmq"
	// SystemPulsar is the apache Pulsar.
	SystemPulsar SystemAttr = "pulsar"
)

// ClientConsumedMessages is an instrument used to record metric values
// conforming to the "messaging.client.consumed.messages" semantic conventions.
// It represents the number of messages that were delivered to the application.
type ClientConsumedMessages struct {
	inst metric.Int64Counter
}

// NewClientConsumedMessages returns a new ClientConsumedMessages instrument.
func NewClientConsumedMessages(m metric.Meter) (ClientConsumedMessages, error) {
	i, err := m.Int64Counter(
	    "messaging.client.consumed.messages",
	    metric.WithDescription("Number of messages that were delivered to the application."),
	    metric.WithUnit("{message}"),
	)
	if err != nil {
	    return ClientConsumedMessages{}, err
	}
	return ClientConsumedMessages{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ClientConsumedMessages) Name() string {
	return "messaging.client.consumed.messages"
}

// Unit returns the semantic convention unit of the instrument
func (ClientConsumedMessages) Unit() string {
	return "{message}"
}

// Description returns the semantic convention description of the instrument
func (ClientConsumedMessages) Description() string {
	return "Number of messages that were delivered to the application."
}

// Add adds incr to the existing count.
//
// The operationName is the the system-specific name of the messaging operation.
//
// The system is the the messaging system as identified by the client
// instrumentation.
//
// All additional attrs passed are included in the recorded value.
func (m ClientConsumedMessages) Add(
	ctx context.Context,
	incr int64,
	operationName string,
	system SystemAttr,
	attrs ...attribute.KeyValue,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("messaging.operation.name", operationName),
				attribute.String("messaging.system", string(system)),
			)...,
		),
	)
}

// AttrErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ClientConsumedMessages) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// AttrConsumerGroupName returns an optional attribute for the
// "messaging.consumer.group.name" semantic convention. It represents the name of
// the consumer group with which a consumer is associated.
func (ClientConsumedMessages) AttrConsumerGroupName(val string) attribute.KeyValue {
	return attribute.String("messaging.consumer.group.name", val)
}

// AttrDestinationName returns an optional attribute for the
// "messaging.destination.name" semantic convention. It represents the message
// destination name.
func (ClientConsumedMessages) AttrDestinationName(val string) attribute.KeyValue {
	return attribute.String("messaging.destination.name", val)
}

// AttrDestinationSubscriptionName returns an optional attribute for the
// "messaging.destination.subscription.name" semantic convention. It represents
// the name of the destination subscription from which a message is consumed.
func (ClientConsumedMessages) AttrDestinationSubscriptionName(val string) attribute.KeyValue {
	return attribute.String("messaging.destination.subscription.name", val)
}

// AttrDestinationTemplate returns an optional attribute for the
// "messaging.destination.template" semantic convention. It represents the low
// cardinality representation of the messaging destination name.
func (ClientConsumedMessages) AttrDestinationTemplate(val string) attribute.KeyValue {
	return attribute.String("messaging.destination.template", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the server domain name if available without
// reverse DNS lookup; otherwise, IP address or Unix domain socket name.
func (ClientConsumedMessages) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// AttrDestinationPartitionID returns an optional attribute for the
// "messaging.destination.partition.id" semantic convention. It represents the
// identifier of the partition messages are sent to or received from, unique
// within the `messaging.destination.name`.
func (ClientConsumedMessages) AttrDestinationPartitionID(val string) attribute.KeyValue {
	return attribute.String("messaging.destination.partition.id", val)
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientConsumedMessages) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// ClientOperationDuration is an instrument used to record metric values
// conforming to the "messaging.client.operation.duration" semantic conventions.
// It represents the duration of messaging operation initiated by a producer or
// consumer client.
type ClientOperationDuration struct {
	inst metric.Float64Histogram
}

// NewClientOperationDuration returns a new ClientOperationDuration instrument.
func NewClientOperationDuration(m metric.Meter) (ClientOperationDuration, error) {
	i, err := m.Float64Histogram(
	    "messaging.client.operation.duration",
	    metric.WithDescription("Duration of messaging operation initiated by a producer or consumer client."),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return ClientOperationDuration{}, err
	}
	return ClientOperationDuration{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ClientOperationDuration) Name() string {
	return "messaging.client.operation.duration"
}

// Unit returns the semantic convention unit of the instrument
func (ClientOperationDuration) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (ClientOperationDuration) Description() string {
	return "Duration of messaging operation initiated by a producer or consumer client."
}

// Record records val to the current distribution.
//
// The operationName is the the system-specific name of the messaging operation.
//
// The system is the the messaging system as identified by the client
// instrumentation.
//
// All additional attrs passed are included in the recorded value.
func (m ClientOperationDuration) Record(
	ctx context.Context,
	val float64,
	operationName string,
	system SystemAttr,
	attrs ...attribute.KeyValue,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("messaging.operation.name", operationName),
				attribute.String("messaging.system", string(system)),
			)...,
		),
	)
}

// AttrErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ClientOperationDuration) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// AttrConsumerGroupName returns an optional attribute for the
// "messaging.consumer.group.name" semantic convention. It represents the name of
// the consumer group with which a consumer is associated.
func (ClientOperationDuration) AttrConsumerGroupName(val string) attribute.KeyValue {
	return attribute.String("messaging.consumer.group.name", val)
}

// AttrDestinationName returns an optional attribute for the
// "messaging.destination.name" semantic convention. It represents the message
// destination name.
func (ClientOperationDuration) AttrDestinationName(val string) attribute.KeyValue {
	return attribute.String("messaging.destination.name", val)
}

// AttrDestinationSubscriptionName returns an optional attribute for the
// "messaging.destination.subscription.name" semantic convention. It represents
// the name of the destination subscription from which a message is consumed.
func (ClientOperationDuration) AttrDestinationSubscriptionName(val string) attribute.KeyValue {
	return attribute.String("messaging.destination.subscription.name", val)
}

// AttrDestinationTemplate returns an optional attribute for the
// "messaging.destination.template" semantic convention. It represents the low
// cardinality representation of the messaging destination name.
func (ClientOperationDuration) AttrDestinationTemplate(val string) attribute.KeyValue {
	return attribute.String("messaging.destination.template", val)
}

// AttrOperationType returns an optional attribute for the
// "messaging.operation.type" semantic convention. It represents a string
// identifying the type of the messaging operation.
func (ClientOperationDuration) AttrOperationType(val OperationTypeAttr) attribute.KeyValue {
	return attribute.String("messaging.operation.type", string(val))
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the server domain name if available without
// reverse DNS lookup; otherwise, IP address or Unix domain socket name.
func (ClientOperationDuration) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// AttrDestinationPartitionID returns an optional attribute for the
// "messaging.destination.partition.id" semantic convention. It represents the
// identifier of the partition messages are sent to or received from, unique
// within the `messaging.destination.name`.
func (ClientOperationDuration) AttrDestinationPartitionID(val string) attribute.KeyValue {
	return attribute.String("messaging.destination.partition.id", val)
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientOperationDuration) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// ClientSentMessages is an instrument used to record metric values conforming to
// the "messaging.client.sent.messages" semantic conventions. It represents the
// number of messages producer attempted to send to the broker.
type ClientSentMessages struct {
	inst metric.Int64Counter
}

// NewClientSentMessages returns a new ClientSentMessages instrument.
func NewClientSentMessages(m metric.Meter) (ClientSentMessages, error) {
	i, err := m.Int64Counter(
	    "messaging.client.sent.messages",
	    metric.WithDescription("Number of messages producer attempted to send to the broker."),
	    metric.WithUnit("{message}"),
	)
	if err != nil {
	    return ClientSentMessages{}, err
	}
	return ClientSentMessages{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ClientSentMessages) Name() string {
	return "messaging.client.sent.messages"
}

// Unit returns the semantic convention unit of the instrument
func (ClientSentMessages) Unit() string {
	return "{message}"
}

// Description returns the semantic convention description of the instrument
func (ClientSentMessages) Description() string {
	return "Number of messages producer attempted to send to the broker."
}

// Add adds incr to the existing count.
//
// The operationName is the the system-specific name of the messaging operation.
//
// The system is the the messaging system as identified by the client
// instrumentation.
//
// All additional attrs passed are included in the recorded value.
func (m ClientSentMessages) Add(
	ctx context.Context,
	incr int64,
	operationName string,
	system SystemAttr,
	attrs ...attribute.KeyValue,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("messaging.operation.name", operationName),
				attribute.String("messaging.system", string(system)),
			)...,
		),
	)
}

// AttrErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ClientSentMessages) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// AttrDestinationName returns an optional attribute for the
// "messaging.destination.name" semantic convention. It represents the message
// destination name.
func (ClientSentMessages) AttrDestinationName(val string) attribute.KeyValue {
	return attribute.String("messaging.destination.name", val)
}

// AttrDestinationTemplate returns an optional attribute for the
// "messaging.destination.template" semantic convention. It represents the low
// cardinality representation of the messaging destination name.
func (ClientSentMessages) AttrDestinationTemplate(val string) attribute.KeyValue {
	return attribute.String("messaging.destination.template", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the server domain name if available without
// reverse DNS lookup; otherwise, IP address or Unix domain socket name.
func (ClientSentMessages) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// AttrDestinationPartitionID returns an optional attribute for the
// "messaging.destination.partition.id" semantic convention. It represents the
// identifier of the partition messages are sent to or received from, unique
// within the `messaging.destination.name`.
func (ClientSentMessages) AttrDestinationPartitionID(val string) attribute.KeyValue {
	return attribute.String("messaging.destination.partition.id", val)
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientSentMessages) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// ProcessDuration is an instrument used to record metric values conforming to
// the "messaging.process.duration" semantic conventions. It represents the
// duration of processing operation.
type ProcessDuration struct {
	inst metric.Float64Histogram
}

// NewProcessDuration returns a new ProcessDuration instrument.
func NewProcessDuration(m metric.Meter) (ProcessDuration, error) {
	i, err := m.Float64Histogram(
	    "messaging.process.duration",
	    metric.WithDescription("Duration of processing operation."),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return ProcessDuration{}, err
	}
	return ProcessDuration{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ProcessDuration) Name() string {
	return "messaging.process.duration"
}

// Unit returns the semantic convention unit of the instrument
func (ProcessDuration) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (ProcessDuration) Description() string {
	return "Duration of processing operation."
}

// Record records val to the current distribution.
//
// The operationName is the the system-specific name of the messaging operation.
//
// The system is the the messaging system as identified by the client
// instrumentation.
//
// All additional attrs passed are included in the recorded value.
func (m ProcessDuration) Record(
	ctx context.Context,
	val float64,
	operationName string,
	system SystemAttr,
	attrs ...attribute.KeyValue,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("messaging.operation.name", operationName),
				attribute.String("messaging.system", string(system)),
			)...,
		),
	)
}

// AttrErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ProcessDuration) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// AttrConsumerGroupName returns an optional attribute for the
// "messaging.consumer.group.name" semantic convention. It represents the name of
// the consumer group with which a consumer is associated.
func (ProcessDuration) AttrConsumerGroupName(val string) attribute.KeyValue {
	return attribute.String("messaging.consumer.group.name", val)
}

// AttrDestinationName returns an optional attribute for the
// "messaging.destination.name" semantic convention. It represents the message
// destination name.
func (ProcessDuration) AttrDestinationName(val string) attribute.KeyValue {
	return attribute.String("messaging.destination.name", val)
}

// AttrDestinationSubscriptionName returns an optional attribute for the
// "messaging.destination.subscription.name" semantic convention. It represents
// the name of the destination subscription from which a message is consumed.
func (ProcessDuration) AttrDestinationSubscriptionName(val string) attribute.KeyValue {
	return attribute.String("messaging.destination.subscription.name", val)
}

// AttrDestinationTemplate returns an optional attribute for the
// "messaging.destination.template" semantic convention. It represents the low
// cardinality representation of the messaging destination name.
func (ProcessDuration) AttrDestinationTemplate(val string) attribute.KeyValue {
	return attribute.String("messaging.destination.template", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the server domain name if available without
// reverse DNS lookup; otherwise, IP address or Unix domain socket name.
func (ProcessDuration) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// AttrDestinationPartitionID returns an optional attribute for the
// "messaging.destination.partition.id" semantic convention. It represents the
// identifier of the partition messages are sent to or received from, unique
// within the `messaging.destination.name`.
func (ProcessDuration) AttrDestinationPartitionID(val string) attribute.KeyValue {
	return attribute.String("messaging.destination.partition.id", val)
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ProcessDuration) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}