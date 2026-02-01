// Package kafka implements the Kafka repository to publish events.
package kafka

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
)

const (
	// TaskCreatedMessageType is the channel used when a Task is created.
	TaskCreatedMessageType = "Task.Created"

	// TaskDeletedMessageType is the channel used when a Task is deleted.
	TaskDeletedMessageType = "Task.Deleted"

	// TaskUpdatedMessageType is the channel used when a Task is updated.
	TaskUpdatedMessageType = "Task.Updated"
)

// Task represents the Message Broker publisher used to publish Task records.
type Task struct {
	producer  *kafka.Producer
	topicName string
}

type event struct {
	Type  string        `json:"type"`
	Value internal.Task `json:"value"`
}

// NewTask instantiates the Task message broker publisher.
func NewTask(producer *kafka.Producer, topicName string) *Task {
	return &Task{
		topicName: topicName,
		producer:  producer,
	}
}

// Created publishes a message indicating a task was created.
func (t *Task) Created(ctx context.Context, task internal.Task) error {
	return t.publish(ctx, TaskCreatedMessageType, task)
}

// Deleted publishes a message indicating a task was deleted.
func (t *Task) Deleted(ctx context.Context, id string) error {
	return t.publish(ctx, TaskDeletedMessageType, internal.Task{ID: id})
}

// Updated publishes a message indicating a task was updated.
func (t *Task) Updated(ctx context.Context, task internal.Task) error {
	return t.publish(ctx, TaskUpdatedMessageType, task)
}

func (t *Task) publish(_ context.Context, msgType string, task internal.Task) error {
	var b bytes.Buffer

	evt := event{
		Type:  msgType,
		Value: task,
	}

	if err := json.NewEncoder(&b).Encode(evt); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.Encode")
	}

	if err := t.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &t.topicName,
			Partition: kafka.PartitionAny,
		},
		Value: b.Bytes(),
	}, nil); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "product.Producer")
	}

	return nil
}
