// Package rabbitmq implements the RabbitMQ repository to publish events.
package rabbitmq

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
)

const (
	// TaskCreatedMessageType is the routing key used when a Task is created.
	TaskCreatedMessageType = "Task.Created"

	// TaskDeletedMessageType is the routing key used when a Task is deleted.
	TaskDeletedMessageType = "Task.Deleted"

	// TaskUpdatedMessageType is the routing key used when a Task is updated.
	TaskUpdatedMessageType = "Task.Updated"

	// ExchangeName is the name of the exchange used for Task messages.
	ExchangeName = "Tasks"

	// RoutingKeyWildcard is the wildcard used to subscribe to all Task events.
	RoutingKeyWildcard = "Task.*"
)

// Task represents the repository used for publishing Task records.
type Task struct {
	ch *amqp.Channel
}

// NewTask instantiates the Task repository.
func NewTask(channel *amqp.Channel) *Task {
	return &Task{
		ch: channel,
	}
}

// Created publishes a message indicating a task was created.
func (t *Task) Created(ctx context.Context, task internal.Task) error {
	return t.publish(ctx, TaskCreatedMessageType, task)
}

// Deleted publishes a message indicating a task was deleted.
func (t *Task) Deleted(ctx context.Context, id string) error {
	return t.publish(ctx, TaskDeletedMessageType, id)
}

// Updated publishes a message indicating a task was updated.
func (t *Task) Updated(ctx context.Context, task internal.Task) error {
	return t.publish(ctx, TaskUpdatedMessageType, task)
}

func (t *Task) publish(ctx context.Context, routingKey string, event any) error {
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(event); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "gob.Encode")
	}

	err := t.ch.PublishWithContext(ctx,
		ExchangeName, // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			AppId:       "tasks-rest-server",
			ContentType: "application/x-encoding-gob", // XXX: We will revisit this in future episodes
			Body:        b.Bytes(),
			Timestamp:   time.Now(),
		})
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "ch.Publish")
	}

	return nil
}
