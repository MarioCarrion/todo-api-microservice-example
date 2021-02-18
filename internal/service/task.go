package service

import (
	"context"
	"fmt"

	"github.com/MarioCarrion/todo-api/internal"
)

// TaskRepository defines the datastore handling persisting Task records.
type TaskRepository interface {
	Create(ctx context.Context, description string, priority internal.Priority, dates internal.Dates) (internal.Task, error)
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
	// XXX: We will revisit the number of received arguments in future episodes.
	task, err := t.repo.Create(ctx, description, priority, dates)
	if err != nil {
		return internal.Task{}, fmt.Errorf("repo create: %w", err)
	}

	return task, nil
}

// Task gets an existing Task from the datastore.
func (t *Task) Task(ctx context.Context, id string) (internal.Task, error) {
	// XXX: We will revisit the number of received arguments in future episodes.
	task, err := t.repo.Find(ctx, id)
	if err != nil {
		return internal.Task{}, fmt.Errorf("repo find: %w", err)
	}

	return task, nil
}

// Update updates an existing Task in the datastore.
func (t *Task) Update(ctx context.Context, id string, description string, priority internal.Priority, dates internal.Dates, isDone bool) error {
	// XXX: We will revisit the number of received arguments in future episodes.
	if err := t.repo.Update(ctx, id, description, priority, dates, isDone); err != nil {
		return fmt.Errorf("repo update: %w", err)
	}

	return nil
}
