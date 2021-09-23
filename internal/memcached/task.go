package memcached

import (
	"context"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"go.uber.org/zap"

	"github.com/MarioCarrion/todo-api/internal"
)

type Task struct {
	client     *memcache.Client
	orig       TaskStore
	expiration time.Duration
	logger     *zap.Logger
}

type TaskStore interface {
	Create(ctx context.Context, description string, priority internal.Priority, dates internal.Dates) (internal.Task, error)
	Delete(ctx context.Context, id string) error
	Find(ctx context.Context, id string) (internal.Task, error)
	Update(ctx context.Context, id string, description string, priority internal.Priority, dates internal.Dates, isDone bool) error
}

func NewTask(client *memcache.Client, orig TaskStore, logger *zap.Logger) *Task {
	return &Task{
		client:     client,
		orig:       orig,
		expiration: 10 * time.Minute,
		logger:     logger,
	}
}

func (t *Task) Create(ctx context.Context, description string, priority internal.Priority, dates internal.Dates) (internal.Task, error) {
	task, err := t.orig.Create(ctx, description, priority, dates)
	if err != nil {
		return internal.Task{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.Create")
	}

	// Write-Through Caching

	t.logger.Info("Create: setting value")

	setTask(t.client, task.ID, &task, t.expiration)

	return task, nil
}

func (t *Task) Delete(ctx context.Context, id string) error {
	if err := t.orig.Delete(ctx, id); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.Delete")
	}

	deleteTask(t.client, id)

	return nil
}

func (t *Task) Find(ctx context.Context, id string) (internal.Task, error) {
	var res internal.Task

	t.logger.Info("Find: get value")

	if err := getTask(t.client, id, &res); err == nil {
		return res, nil
	}

	t.logger.Info("Find: not found, let's cache it")

	// Cache-Aside Caching

	res, err := t.orig.Find(ctx, id)
	if err != nil {
		return res, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.Find")
	}

	setTask(t.client, res.ID, &res, t.expiration)

	return res, nil
}

func (t *Task) Update(ctx context.Context, id string, description string, priority internal.Priority, dates internal.Dates, isDone bool) error {
	if err := t.orig.Update(ctx, id, description, priority, dates, isDone); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.Update")
	}

	// Write-Through Caching

	t.logger.Info("Update: setting value")

	// Update cache

	// XXX:
	// What if any of the following instructions fail? We may end up with stale
	// values

	deleteTask(t.client, id) // XXX

	task, err := t.orig.Find(ctx, id)
	if err != nil { // XXX
		return nil //nolint: nilerr
	}

	setTask(t.client, task.ID, &task, t.expiration) // XXX

	return nil
}
