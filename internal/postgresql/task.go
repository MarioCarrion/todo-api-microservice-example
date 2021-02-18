package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/MarioCarrion/todo-api/internal"
)

// Task represents the repository used for interacting with Task records.
type Task struct {
	q *Queries
}

// NewTask instantiates the Task repository.
func NewTask(db *sql.DB) *Task {
	return &Task{
		q: New(db),
	}
}

// Create inserts a new task record.
func (t *Task) Create(ctx context.Context, description string, priority internal.Priority, dates internal.Dates) (internal.Task, error) {
	// XXX: `ID` and `IsDone` make no sense when creating new records, that's why those are ignored.
	// XXX: We will revisit the number of received arguments in future episodes.
	// XXX: We are intentionally NOT SUPPORTING `SubTasks` and `Categories` JUST YET.

	id, err := t.q.InsertTask(ctx, InsertTaskParams{
		Description: description,
		Priority:    newPriority(priority),
		StartDate:   newNullTime(dates.Start),
		DueDate:     newNullTime(dates.Due),
	})
	if err != nil {
		return internal.Task{}, fmt.Errorf("insert: %w", err)
	}

	return internal.Task{
		ID:          id.String(),
		Description: description,
		Priority:    priority,
		Dates:       dates,
	}, nil
}

// Find returns the requested task by searching its id.
func (t *Task) Find(ctx context.Context, id string) (internal.Task, error) {
	val, err := uuid.Parse(id)
	if err != nil {
		return internal.Task{}, fmt.Errorf("uuid parse: %w", err)
	}

	res, err := t.q.SelectTask(ctx, val)
	if err != nil {
		return internal.Task{}, fmt.Errorf("select: %w", err)
	}

	priority, err := convertPriority(res.Priority)
	if err != nil {
		return internal.Task{}, fmt.Errorf("convert priority: %w", err)
	}

	return internal.Task{
		ID:          res.ID.String(),
		Description: res.Description,
		Priority:    priority,
		Dates: internal.Dates{
			Start: res.StartDate.Time,
			Due:   res.DueDate.Time,
		},
		IsDone: res.Done,
	}, nil
}

// Update updates the existing record with new values.
func (t *Task) Update(ctx context.Context, id string, description string, priority internal.Priority, dates internal.Dates, isDone bool) error {
	// XXX: We will revisit the number of received arguments in future episodes.
	val, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("uuid parse: %w", err)
	}

	if err := t.q.UpdateTask(ctx, UpdateTaskParams{
		ID:          val,
		Description: description,
		Priority:    newPriority(priority),
		StartDate:   newNullTime(dates.Start),
		DueDate:     newNullTime(dates.Due),
		Done:        isDone,
	}); err != nil {
		return fmt.Errorf("update: %w", err)
	}

	return nil
}
