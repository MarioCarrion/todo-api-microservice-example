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
		end   pgtype.Timestamp
	)

	var dates *internal.Dates

	if params.Dates != nil {
		start = newTimestamp(params.Dates.Start)
		end = newTimestamp(params.Dates.Due)

		dates = &internal.Dates{
			Start: params.Dates.Start,
			Due:   params.Dates.Due,
		}
	}

	// XXX: `ID` and `IsDone` make no sense when creating new records, that's why those are ignored.
	// XXX: We are intentionally NOT SUPPORTING `SubTasks` and `Categories` JUST YET.
	newID, err := t.q.InsertTask(ctx, db.InsertTaskParams{
		Description: params.Description,
		Priority:    newPriority(params.Priority),
		StartDate:   start,
		DueDate:     end,
	})
	if err != nil {
		return internal.Task{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "insert task")
	}

	return internal.Task{
		ID:          newID,
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
	if res.StartDate.Valid {
		dates = &internal.Dates{
			Start: &res.StartDate.Time,
		}
	}

	if res.DueDate.Valid {
		dates = newDates(dates)
		dates.Due = &res.DueDate.Time
	}

	return internal.Task{
		ID:          res.ID,
		Description: res.Description,
		Priority:    &priority,
		Dates:       dates,
		IsDone:      res.Done,
	}, nil
}

// Update updates the existing record with new values.
func (t *Task) Update(ctx context.Context, id string, description string, priority internal.Priority, dates internal.Dates, isDone bool) error {
	// XXX: We will revisit the number of received arguments in future episodes.
	val, err := uuid.Parse(id)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "invalid uuid")
	}

	if _, err := t.q.UpdateTask(ctx, db.UpdateTaskParams{
		ID:          val,
		Description: description,
		Priority:    newPriority(&priority),
		StartDate:   newTimestamp(dates.Start),
		DueDate:     newTimestamp(dates.Due),
		Done:        isDone,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return internal.WrapErrorf(err, internal.ErrorCodeNotFound, "task not found")
		}

		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "update task")
	}

	return nil
}

func newDates(dates *internal.Dates) *internal.Dates {
	if dates == nil {
		return &internal.Dates{}
	}
	return dates
}
