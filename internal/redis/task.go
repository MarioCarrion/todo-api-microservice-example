package redis

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
)

// Task represents the repository used for publishing Task records.
type Task struct {
	client *redis.Client
}

// NewTask instantiates the Task repository.
func NewTask(client *redis.Client) *Task {
	return &Task{
		client: client,
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

func (t *Task) publish(ctx context.Context, channel string, event any) error {
	var b bytes.Buffer

	if err := json.NewEncoder(&b).Encode(event); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.Encode")
	}

	res := t.client.Publish(ctx, channel, b.Bytes())
	if err := res.Err(); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "client.Publish")
	}

	return nil
}
