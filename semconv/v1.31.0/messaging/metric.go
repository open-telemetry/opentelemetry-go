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
// The messagingOperationName is the the system-specific name of the messaging
// operation.
//
// The messagingSystem is the the messaging system as identified by the client
// instrumentation.
//
// All additional attrs passed are included in the recorded value.
func (m ClientConsumedMessages) Add(
    ctx context.Context,
    incr int64,
	operationName string,
	system SystemAttr,
	attrs ...ClientConsumedMessagesAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("messaging.operation.name", operationName),
				attribute.String("messaging.system", string(system)),
			)...,
		),
	)
}

func (m ClientConsumedMessages) conv(in []ClientConsumedMessagesAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.clientConsumedMessagesAttr()
	}
	return out
}

// ClientConsumedMessagesAttr is an optional attribute for the
// ClientConsumedMessages instrument.
type ClientConsumedMessagesAttr interface {
    clientConsumedMessagesAttr() attribute.KeyValue
}

type clientConsumedMessagesAttr struct {
	kv attribute.KeyValue
}

func (a clientConsumedMessagesAttr) clientConsumedMessagesAttr() attribute.KeyValue {
    return a.kv
}

// ErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ClientConsumedMessages) ErrorTypeAttr(val ErrorTypeAttr) ClientConsumedMessagesAttr {
	return clientConsumedMessagesAttr{kv: attribute.String("error.type", string(val))}
}

// ConsumerGroupName returns an optional attribute for the
// "messaging.consumer.group.name" semantic convention. It represents the name of
// the consumer group with which a consumer is associated.
func (ClientConsumedMessages) ConsumerGroupNameAttr(val string) ClientConsumedMessagesAttr {
	return clientConsumedMessagesAttr{kv: attribute.String("messaging.consumer.group.name", val)}
}

// DestinationName returns an optional attribute for the
// "messaging.destination.name" semantic convention. It represents the message
// destination name.
func (ClientConsumedMessages) DestinationNameAttr(val string) ClientConsumedMessagesAttr {
	return clientConsumedMessagesAttr{kv: attribute.String("messaging.destination.name", val)}
}

// DestinationSubscriptionName returns an optional attribute for the
// "messaging.destination.subscription.name" semantic convention. It represents
// the name of the destination subscription from which a message is consumed.
func (ClientConsumedMessages) DestinationSubscriptionNameAttr(val string) ClientConsumedMessagesAttr {
	return clientConsumedMessagesAttr{kv: attribute.String("messaging.destination.subscription.name", val)}
}

// DestinationTemplate returns an optional attribute for the
// "messaging.destination.template" semantic convention. It represents the low
// cardinality representation of the messaging destination name.
func (ClientConsumedMessages) DestinationTemplateAttr(val string) ClientConsumedMessagesAttr {
	return clientConsumedMessagesAttr{kv: attribute.String("messaging.destination.template", val)}
}

// ServerAddress returns an optional attribute for the "server.address" semantic
// convention. It represents the server domain name if available without reverse
// DNS lookup; otherwise, IP address or Unix domain socket name.
func (ClientConsumedMessages) ServerAddressAttr(val string) ClientConsumedMessagesAttr {
	return clientConsumedMessagesAttr{kv: attribute.String("server.address", val)}
}

// DestinationPartitionID returns an optional attribute for the
// "messaging.destination.partition.id" semantic convention. It represents the
// identifier of the partition messages are sent to or received from, unique
// within the `messaging.destination.name`.
func (ClientConsumedMessages) DestinationPartitionIDAttr(val string) ClientConsumedMessagesAttr {
	return clientConsumedMessagesAttr{kv: attribute.String("messaging.destination.partition.id", val)}
}

// ServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientConsumedMessages) ServerPortAttr(val int) ClientConsumedMessagesAttr {
	return clientConsumedMessagesAttr{kv: attribute.Int("server.port", val)}
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
// The messagingOperationName is the the system-specific name of the messaging
// operation.
//
// The messagingSystem is the the messaging system as identified by the client
// instrumentation.
//
// All additional attrs passed are included in the recorded value.
func (m ClientOperationDuration) Record(
    ctx context.Context,
    val float64,
	operationName string,
	system SystemAttr,
	attrs ...ClientOperationDurationAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("messaging.operation.name", operationName),
				attribute.String("messaging.system", string(system)),
			)...,
		),
	)
}

func (m ClientOperationDuration) conv(in []ClientOperationDurationAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.clientOperationDurationAttr()
	}
	return out
}

// ClientOperationDurationAttr is an optional attribute for the
// ClientOperationDuration instrument.
type ClientOperationDurationAttr interface {
    clientOperationDurationAttr() attribute.KeyValue
}

type clientOperationDurationAttr struct {
	kv attribute.KeyValue
}

func (a clientOperationDurationAttr) clientOperationDurationAttr() attribute.KeyValue {
    return a.kv
}

// ErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ClientOperationDuration) ErrorTypeAttr(val ErrorTypeAttr) ClientOperationDurationAttr {
	return clientOperationDurationAttr{kv: attribute.String("error.type", string(val))}
}

// ConsumerGroupName returns an optional attribute for the
// "messaging.consumer.group.name" semantic convention. It represents the name of
// the consumer group with which a consumer is associated.
func (ClientOperationDuration) ConsumerGroupNameAttr(val string) ClientOperationDurationAttr {
	return clientOperationDurationAttr{kv: attribute.String("messaging.consumer.group.name", val)}
}

// DestinationName returns an optional attribute for the
// "messaging.destination.name" semantic convention. It represents the message
// destination name.
func (ClientOperationDuration) DestinationNameAttr(val string) ClientOperationDurationAttr {
	return clientOperationDurationAttr{kv: attribute.String("messaging.destination.name", val)}
}

// DestinationSubscriptionName returns an optional attribute for the
// "messaging.destination.subscription.name" semantic convention. It represents
// the name of the destination subscription from which a message is consumed.
func (ClientOperationDuration) DestinationSubscriptionNameAttr(val string) ClientOperationDurationAttr {
	return clientOperationDurationAttr{kv: attribute.String("messaging.destination.subscription.name", val)}
}

// DestinationTemplate returns an optional attribute for the
// "messaging.destination.template" semantic convention. It represents the low
// cardinality representation of the messaging destination name.
func (ClientOperationDuration) DestinationTemplateAttr(val string) ClientOperationDurationAttr {
	return clientOperationDurationAttr{kv: attribute.String("messaging.destination.template", val)}
}

// OperationType returns an optional attribute for the "messaging.operation.type"
// semantic convention. It represents a string identifying the type of the
// messaging operation.
func (ClientOperationDuration) OperationTypeAttr(val OperationTypeAttr) ClientOperationDurationAttr {
	return clientOperationDurationAttr{kv: attribute.String("messaging.operation.type", string(val))}
}

// ServerAddress returns an optional attribute for the "server.address" semantic
// convention. It represents the server domain name if available without reverse
// DNS lookup; otherwise, IP address or Unix domain socket name.
func (ClientOperationDuration) ServerAddressAttr(val string) ClientOperationDurationAttr {
	return clientOperationDurationAttr{kv: attribute.String("server.address", val)}
}

// DestinationPartitionID returns an optional attribute for the
// "messaging.destination.partition.id" semantic convention. It represents the
// identifier of the partition messages are sent to or received from, unique
// within the `messaging.destination.name`.
func (ClientOperationDuration) DestinationPartitionIDAttr(val string) ClientOperationDurationAttr {
	return clientOperationDurationAttr{kv: attribute.String("messaging.destination.partition.id", val)}
}

// ServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientOperationDuration) ServerPortAttr(val int) ClientOperationDurationAttr {
	return clientOperationDurationAttr{kv: attribute.Int("server.port", val)}
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
// The messagingOperationName is the the system-specific name of the messaging
// operation.
//
// The messagingSystem is the the messaging system as identified by the client
// instrumentation.
//
// All additional attrs passed are included in the recorded value.
func (m ClientSentMessages) Add(
    ctx context.Context,
    incr int64,
	operationName string,
	system SystemAttr,
	attrs ...ClientSentMessagesAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("messaging.operation.name", operationName),
				attribute.String("messaging.system", string(system)),
			)...,
		),
	)
}

func (m ClientSentMessages) conv(in []ClientSentMessagesAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.clientSentMessagesAttr()
	}
	return out
}

// ClientSentMessagesAttr is an optional attribute for the ClientSentMessages
// instrument.
type ClientSentMessagesAttr interface {
    clientSentMessagesAttr() attribute.KeyValue
}

type clientSentMessagesAttr struct {
	kv attribute.KeyValue
}

func (a clientSentMessagesAttr) clientSentMessagesAttr() attribute.KeyValue {
    return a.kv
}

// ErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ClientSentMessages) ErrorTypeAttr(val ErrorTypeAttr) ClientSentMessagesAttr {
	return clientSentMessagesAttr{kv: attribute.String("error.type", string(val))}
}

// DestinationName returns an optional attribute for the
// "messaging.destination.name" semantic convention. It represents the message
// destination name.
func (ClientSentMessages) DestinationNameAttr(val string) ClientSentMessagesAttr {
	return clientSentMessagesAttr{kv: attribute.String("messaging.destination.name", val)}
}

// DestinationTemplate returns an optional attribute for the
// "messaging.destination.template" semantic convention. It represents the low
// cardinality representation of the messaging destination name.
func (ClientSentMessages) DestinationTemplateAttr(val string) ClientSentMessagesAttr {
	return clientSentMessagesAttr{kv: attribute.String("messaging.destination.template", val)}
}

// ServerAddress returns an optional attribute for the "server.address" semantic
// convention. It represents the server domain name if available without reverse
// DNS lookup; otherwise, IP address or Unix domain socket name.
func (ClientSentMessages) ServerAddressAttr(val string) ClientSentMessagesAttr {
	return clientSentMessagesAttr{kv: attribute.String("server.address", val)}
}

// DestinationPartitionID returns an optional attribute for the
// "messaging.destination.partition.id" semantic convention. It represents the
// identifier of the partition messages are sent to or received from, unique
// within the `messaging.destination.name`.
func (ClientSentMessages) DestinationPartitionIDAttr(val string) ClientSentMessagesAttr {
	return clientSentMessagesAttr{kv: attribute.String("messaging.destination.partition.id", val)}
}

// ServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientSentMessages) ServerPortAttr(val int) ClientSentMessagesAttr {
	return clientSentMessagesAttr{kv: attribute.Int("server.port", val)}
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
// The messagingOperationName is the the system-specific name of the messaging
// operation.
//
// The messagingSystem is the the messaging system as identified by the client
// instrumentation.
//
// All additional attrs passed are included in the recorded value.
func (m ProcessDuration) Record(
    ctx context.Context,
    val float64,
	operationName string,
	system SystemAttr,
	attrs ...ProcessDurationAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("messaging.operation.name", operationName),
				attribute.String("messaging.system", string(system)),
			)...,
		),
	)
}

func (m ProcessDuration) conv(in []ProcessDurationAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.processDurationAttr()
	}
	return out
}

// ProcessDurationAttr is an optional attribute for the ProcessDuration
// instrument.
type ProcessDurationAttr interface {
    processDurationAttr() attribute.KeyValue
}

type processDurationAttr struct {
	kv attribute.KeyValue
}

func (a processDurationAttr) processDurationAttr() attribute.KeyValue {
    return a.kv
}

// ErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ProcessDuration) ErrorTypeAttr(val ErrorTypeAttr) ProcessDurationAttr {
	return processDurationAttr{kv: attribute.String("error.type", string(val))}
}

// ConsumerGroupName returns an optional attribute for the
// "messaging.consumer.group.name" semantic convention. It represents the name of
// the consumer group with which a consumer is associated.
func (ProcessDuration) ConsumerGroupNameAttr(val string) ProcessDurationAttr {
	return processDurationAttr{kv: attribute.String("messaging.consumer.group.name", val)}
}

// DestinationName returns an optional attribute for the
// "messaging.destination.name" semantic convention. It represents the message
// destination name.
func (ProcessDuration) DestinationNameAttr(val string) ProcessDurationAttr {
	return processDurationAttr{kv: attribute.String("messaging.destination.name", val)}
}

// DestinationSubscriptionName returns an optional attribute for the
// "messaging.destination.subscription.name" semantic convention. It represents
// the name of the destination subscription from which a message is consumed.
func (ProcessDuration) DestinationSubscriptionNameAttr(val string) ProcessDurationAttr {
	return processDurationAttr{kv: attribute.String("messaging.destination.subscription.name", val)}
}

// DestinationTemplate returns an optional attribute for the
// "messaging.destination.template" semantic convention. It represents the low
// cardinality representation of the messaging destination name.
func (ProcessDuration) DestinationTemplateAttr(val string) ProcessDurationAttr {
	return processDurationAttr{kv: attribute.String("messaging.destination.template", val)}
}

// ServerAddress returns an optional attribute for the "server.address" semantic
// convention. It represents the server domain name if available without reverse
// DNS lookup; otherwise, IP address or Unix domain socket name.
func (ProcessDuration) ServerAddressAttr(val string) ProcessDurationAttr {
	return processDurationAttr{kv: attribute.String("server.address", val)}
}

// DestinationPartitionID returns an optional attribute for the
// "messaging.destination.partition.id" semantic convention. It represents the
// identifier of the partition messages are sent to or received from, unique
// within the `messaging.destination.name`.
func (ProcessDuration) DestinationPartitionIDAttr(val string) ProcessDurationAttr {
	return processDurationAttr{kv: attribute.String("messaging.destination.partition.id", val)}
}

// ServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ProcessDuration) ServerPortAttr(val int) ProcessDurationAttr {
	return processDurationAttr{kv: attribute.Int("server.port", val)}
}