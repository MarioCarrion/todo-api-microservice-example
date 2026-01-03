// Package postgresql implements the PostgreSQL repository for tasks.
package postgresql

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
	"github.com/MarioCarrion/todo-api-microservice-example/internal/postgresql/db"
)

// Task represents the repository used for interacting with Task records.
type Task struct {
	q *db.Queries
}

// NewTask instantiates the Task repository.
func NewTask(d db.DBTX) *Task {
	return &Task{
		q: db.New(d),
	}
}

// Create inserts a new task record.
func (t *Task) Create(ctx context.Context, params internal.CreateParams) (internal.Task, error) {
	var (
		start pgtype.Timestamp
		due   pgtype.Timestamp
		dates *internal.Dates
	)

	if params.Dates != nil {
		start = newTimestamp(params.Dates.Start)
		due = newTimestamp(params.Dates.Due)

		dates = &internal.Dates{
			Start: params.Dates.Start,
			Due:   params.Dates.Due,
		}
	}

	// TODO: We are intentionally NOT SUPPORTING `SubTasks` and `Categories` JUST YET.
	newID, err := t.q.InsertTask(ctx, db.InsertTaskParams{
		Description: params.Description,
		Priority:    newPriority(params.Priority),
		StartDate:   start,
		DueDate:     due,
	})
	if err != nil {
		return internal.Task{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "insert task")
	}

	return internal.Task{
		ID:          newID.String(),
		Description: params.Description,
		Priority:    params.Priority,
		Dates:       dates,
	}, nil
}

// Delete deletes the existing record matching the id.
func (t *Task) Delete(ctx context.Context, id string) error {
	val, err := uuid.Parse(id)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "invalid uuid")
	}

	_, err = t.q.DeleteTask(ctx, val)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return internal.WrapErrorf(err, internal.ErrorCodeNotFound, "task not found")
		}

		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "delete task")
	}

	return nil
}

// Find returns the requested task by searching its id.
func (t *Task) Find(ctx context.Context, id string) (internal.Task, error) {
	val, err := uuid.Parse(id)
	if err != nil {
		return internal.Task{}, internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "invalid uuid")
	}

	res, err := t.q.SelectTask(ctx, val)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return internal.Task{}, internal.WrapErrorf(err, internal.ErrorCodeNotFound, "task not found")
		}

		return internal.Task{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "select task")
	}

	priority, err := convertPriority(res.Priority)
	if err != nil {
		return internal.Task{}, internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "convert priority")
	}

	var dates *internal.Dates
	if res.StartDate.Valid || res.DueDate.Valid {
		dates = &internal.Dates{}

		if res.StartDate.Valid {
			dates.Start = &res.StartDate.Time
		}

		if res.DueDate.Valid {
			dates.Due = &res.DueDate.Time
		}
	}

	return internal.Task{
		ID:          res.ID.String(),
		Description: res.Description,
		Priority:    &priority,
		Dates:       dates,
		IsDone:      res.Done,
	}, nil
}

// Update updates the existing record with new values.
func (t *Task) Update(ctx context.Context, id string, params internal.UpdateParams) error {
	// XXX: We will revisit the number of received arguments in future episodes.
	val, err := uuid.Parse(id)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "invalid uuid")
	}

	var (
		start pgtype.Timestamp
		due   pgtype.Timestamp
	)

	if params.Dates != nil {
		if params.Dates.Start != nil {
			start = newTimestamp(params.Dates.Start)
		}

		if params.Dates.Due != nil {
			due = newTimestamp(params.Dates.Due)
		}
	}

	if _, err := t.q.UpdateTask(ctx, db.UpdateTaskParams{
		ID:          val,
		Description: internal.PointerToValue(params.Description),
		Priority:    newPriority(params.Priority),
		StartDate:   start,
		DueDate:     due,
		Done:        internal.PointerToValue(params.IsDone),
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return internal.WrapErrorf(err, internal.ErrorCodeNotFound, "task not found")
		}

		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "update task")
	}

	return nil
}
