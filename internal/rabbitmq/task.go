package rabbitmq

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	"github.com/streadway/amqp"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
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
	return t.publish(ctx, "Task.Created", task)
}

// Deleted publishes a message indicating a task was deleted.
func (t *Task) Deleted(ctx context.Context, id string) error {
	return t.publish(ctx, "Task.Deleted", id)
}

// Updated publishes a message indicating a task was updated.
func (t *Task) Updated(ctx context.Context, task internal.Task) error {
	return t.publish(ctx, "Task.Updated", task)
}

func (t *Task) publish(_ context.Context, routingKey string, event interface{}) error {
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(event); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "gob.Encode")
	}

	err := t.ch.Publish(
		"tasks",    // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
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
