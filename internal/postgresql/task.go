package postgresql

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/MarioCarrion/todo-api/internal"
	"github.com/MarioCarrion/todo-api/internal/postgresql/db"
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
func (t *Task) Create(ctx context.Context, description string, priority internal.Priority, dates internal.Dates) (internal.Task, error) {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Task.Create")
	span.SetAttributes(attribute.String("db.system", "postgresql"))

	defer span.End()

	// XXX: `ID` and `IsDone` make no sense when creating new records, that's why those are ignored.
	// XXX: We will revisit the number of received arguments in future episodes.
	// XXX: We are intentionally NOT SUPPORTING `SubTasks` and `Categories` JUST YET.

	id, err := t.q.InsertTask(ctx, db.InsertTaskParams{
		Description: description,
		Priority:    newPriority(priority),
		StartDate:   newNullTime(dates.Start),
		DueDate:     newNullTime(dates.Due),
	})
	if err != nil {
		return internal.Task{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "insert task")
	}

	return internal.Task{
		ID:          id.String(),
		Description: description,
		Priority:    priority,
		Dates:       dates,
	}, nil
}

// Delete deletes the existing record matching the id.
func (t *Task) Delete(ctx context.Context, id string) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Task.Delete")
	span.SetAttributes(attribute.String("db.system", "postgresql"))

	defer span.End()

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
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Task.Find")
	span.SetAttributes(attribute.String("db.system", "postgresql"))

	defer span.End()

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
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "Task.Update")
	span.SetAttributes(attribute.String("db.system", "postgresql"))

	defer span.End()

	// XXX: We will revisit the number of received arguments in future episodes.
	val, err := uuid.Parse(id)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "invalid uuid")
	}

	if _, err := t.q.UpdateTask(ctx, db.UpdateTaskParams{
		ID:          val,
		Description: description,
		Priority:    newPriority(priority),
		StartDate:   newNullTime(dates.Start),
		DueDate:     newNullTime(dates.Due),
		Done:        isDone,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return internal.WrapErrorf(err, internal.ErrorCodeNotFound, "task not found")
		}

		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "update task")
	}

	return nil
}
