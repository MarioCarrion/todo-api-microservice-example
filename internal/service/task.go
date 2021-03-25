package service

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"

	"github.com/MarioCarrion/todo-api/internal"
)

// TaskRepository defines the datastore handling persisting Task records.
type TaskRepository interface {
	Create(ctx context.Context, description string, priority internal.Priority, dates internal.Dates) (internal.Task, error)
	Delete(ctx context.Context, id string) error
	Find(ctx context.Context, id string) (internal.Task, error)
	Update(ctx context.Context, id string, description string, priority internal.Priority, dates internal.Dates, isDone bool) error
}

// Task defines the application service in charge of interacting with Tasks.
type Task struct {
	repo TaskRepository
}

// NewTask ...
func NewTask(repo TaskRepository) *Task {
	return &Task{
		repo: repo,
	}
}

// Create stores a new record.
func (t *Task) Create(ctx context.Context, description string, priority internal.Priority, dates internal.Dates) (internal.Task, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Task.Create")
	defer span.End()

	// XXX: We will revisit the number of received arguments in future episodes.
	task, err := t.repo.Create(ctx, description, priority, dates)
	if err != nil {
		return internal.Task{}, fmt.Errorf("repo create: %w", err)
	}

	return task, nil
}

// Delete removes an existing Task from the datastore.
func (t *Task) Delete(ctx context.Context, id string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Task.Delete")
	defer span.End()

	// XXX: We will revisit the number of received arguments in future episodes.
	if err := t.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("repo delete: %w", err)
	}

	return nil
}

// Task gets an existing Task from the datastore.
func (t *Task) Task(ctx context.Context, id string) (internal.Task, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Task.Task")
	defer span.End()

	// XXX: We will revisit the number of received arguments in future episodes.
	task, err := t.repo.Find(ctx, id)
	if err != nil {
		return internal.Task{}, fmt.Errorf("repo find: %w", err)
	}

	return task, nil
}

// Update updates an existing Task in the datastore.
func (t *Task) Update(ctx context.Context, id string, description string, priority internal.Priority, dates internal.Dates, isDone bool) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Task.Update")
	defer span.End()

	// XXX: We will revisit the number of received arguments in future episodes.
	if err := t.repo.Update(ctx, id, description, priority, dates, isDone); err != nil {
		return fmt.Errorf("repo update: %w", err)
	}

	return nil
}
