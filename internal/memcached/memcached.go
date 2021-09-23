package memcached

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/bradfitz/gomemcache/memcache"

	"github.com/MarioCarrion/todo-api/internal"
)

// XXX: "delete" and "set" intentionally ignore errors, a better approach
// would be to implement an unexported "client" type defining all the three
// methods defined in this file that includes a Circuit Breaker for logging
// errors and retry as needed.
// See https://youtu.be/UnL2iGcD7vE for more details about that pattern.

func deleteTask(client *memcache.Client, key string) {
	_ = client.Delete(key)
}

func getTask(client *memcache.Client, key string, target interface{}) error {
	item, err := client.Get(key)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "client.Get")
	}

	if err := gob.NewDecoder(bytes.NewReader(item.Value)).Decode(target); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "gob.NewDecoder")
	}

	return nil
}

func setTask(client *memcache.Client, key string, value interface{}, expiration time.Duration) {
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(value); err != nil {
		return
	}

	_ = client.Set(&memcache.Item{
		Key:        key,
		Value:      b.Bytes(),
		Expiration: int32(time.Now().Add(expiration).Unix()),
	})
}
