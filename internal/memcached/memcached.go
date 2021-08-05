package memcached

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"go.uber.org/zap"

	"github.com/MarioCarrion/todo-api/internal"
)

// Task ...
type Task struct {
	client *memcache.Client
	orig   Datastore
	logger *zap.Logger
}

type Datastore interface {
	Delete(ctx context.Context, id string) error
	Index(ctx context.Context, task internal.Task) error
	Search(ctx context.Context, args internal.SearchArgs) (internal.SearchResults, error)
}

// NewTask instantiates the Task repository.
func NewTask(client *memcache.Client, orig Datastore, logger *zap.Logger) *Task {
	return &Task{
		client: client,
		orig:   orig,
		logger: logger,
	}
}

// Index ...
func (t *Task) Index(ctx context.Context, task internal.Task) error {
	return t.orig.Index(ctx, task)
}

// Delete ...
func (t *Task) Delete(ctx context.Context, id string) error {
	return t.orig.Delete(ctx, id)
}

// Search ...
func (t *Task) Search(ctx context.Context, args internal.SearchArgs) (internal.SearchResults, error) {
	key := newKey(args)

	item, err := t.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			t.logger.Info("values NOT found", zap.String("key", string(key)))

			res, err := t.orig.Search(ctx, args)
			if err != nil {
				return internal.SearchResults{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "orig.Search")
			}

			var b bytes.Buffer

			if err := gob.NewEncoder(&b).Encode(&res); err == nil {
				t.logger.Info("settin value")

				_ = t.client.Set(&memcache.Item{
					Key:        key,
					Value:      b.Bytes(),
					Expiration: int32(time.Now().Add(25 * time.Second).Unix()),
				})
			}

			return res, err
		}

		return internal.SearchResults{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "client.Get")
	}

	t.logger.Info("values found", zap.String("key", string(key)))

	var res internal.SearchResults

	if err := gob.NewDecoder(bytes.NewReader(item.Value)).Decode(&res); err != nil {
		return internal.SearchResults{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "gob.NewDecoder")
	}

	return res, nil
}

func newKey(args internal.SearchArgs) string {
	var (
		description string
		priority    int8
		isDone      bool
	)

	if args.Description != nil {
		description = *args.Description
	}

	if args.Priority != nil {
		priority = int8(*args.Priority)
	}

	if args.IsDone != nil {
		isDone = *args.IsDone
	}

	return fmt.Sprintf("%s_%d_%t_%d_%d", description, priority, isDone, args.From, args.Size)
}
