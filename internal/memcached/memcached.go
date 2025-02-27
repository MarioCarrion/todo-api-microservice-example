package memcached

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	"github.com/bradfitz/gomemcache/memcache"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
)

// XXX: "delete" and "set" intentionally ignore errors, a better approach
// would be to implement an unexported "client" type defining all the three
// methods defined in this file that includes a Circuit Breaker for logging
// errors and retry as needed.
// See https://youtu.be/UnL2iGcD7vE for more details about that pattern.

func deleteTask(_ context.Context, client *memcache.Client, key string) {
	_ = client.Delete(key)
}

func getTask(_ context.Context, client *memcache.Client, key string, target any) error {
	item, err := client.Get(key)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "client.Get")
	}

	if err := gob.NewDecoder(bytes.NewReader(item.Value)).Decode(target); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "gob.NewDecoder")
	}

	return nil
}

func setTask(_ context.Context, client *memcache.Client, key string, value any, expiration time.Duration) {
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(value); err != nil {
		return
	}

	_ = client.Set(&memcache.Item{
		Key:        key,
		Value:      b.Bytes(),
		Expiration: int32(time.Now().Add(expiration).Unix()), //nolint:gosec
	})
}
