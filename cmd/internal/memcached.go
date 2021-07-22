package internal

import (
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"

	"github.com/MarioCarrion/todo-api/internal/envvar"
)

func NewMemcached(conf *envvar.Configuration) (*memcache.Client, error) {
	host, err := conf.Get("MEMCACHED_HOST")
	if err != nil {
		return nil, fmt.Errorf("conf.Get MEMCACHED_HOST %w", err)
	}

	// XXX Assuming environment variable contains only one server
	client := memcache.New(host)

	if err := client.Ping(); err != nil {
		return nil, err
	}

	client.Timeout = 100 * time.Millisecond
	client.MaxIdleConns = 100

	return client, nil
}
