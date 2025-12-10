package redis_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/go-cmp/cmp"
	"github.com/testcontainers/testcontainers-go"
	redistest "github.com/testcontainers/testcontainers-go/modules/redis"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
	redistask "github.com/MarioCarrion/todo-api-microservice-example/internal/redis"
)

// dockerImage must match the docker image listed in `compose.redis.yml`.
const dockerImage = "redis:7.0.9-alpine3.17"

func TestMain(m *testing.M) {
	client := setupClient()
	if client.err != nil {
		panic(fmt.Sprintf("Failed to set up Redis client: %v", client.err))
	}

	code := m.Run()

	if err := client.Teardown(); err != nil {
		panic(fmt.Sprintf("Failed to close Redis client: %v", err))
	}

	os.Exit(code)
}

func TestTask_All(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		newPubsub func(t *testing.T, client *redis.Client) *redis.PubSub
		call      func(t *testing.T, client *redis.Client)
		verify    func(t *testing.T, msg *redis.Message)
	}{
		{
			name: "Created",
			newPubsub: func(t *testing.T, client *redis.Client) *redis.PubSub {
				t.Helper()

				pubsub := client.Subscribe(t.Context(), redistask.TaskCreatedChannel)
				t.Cleanup(func() {
					_ = pubsub.Close()
				})

				_, err := pubsub.Receive(t.Context())
				if err != nil {
					t.Fatalf("Failed to receive from pubsub: %v", err)
				}

				return pubsub
			},
			call: func(t *testing.T, client *redis.Client) {
				t.Helper()

				taskPub := redistask.NewTask(client)

				task := internal.Task{
					ID:          "test-123",
					IsDone:      false,
					Description: "Test task",
				}

				if err := taskPub.Created(t.Context(), task); err != nil {
					t.Fatalf("Failed to publish created event: %v", err)
				}
			},
			verify: func(t *testing.T, msg *redis.Message) {
				t.Helper()

				var got internal.Task

				if err := json.NewDecoder(strings.NewReader(msg.Payload)).Decode(&got); err != nil {
					t.Fatalf("Failed to decode message payload: %v", err)
				}

				task := internal.Task{
					ID:          "test-123",
					IsDone:      false,
					Description: "Test task",
				}

				if diff := cmp.Diff(got, task); diff != "" {
					t.Fatalf("Received task is not the same as the created one: %s", diff)
				}
			},
		},
		{
			name: "Updated",
			newPubsub: func(t *testing.T, client *redis.Client) *redis.PubSub {
				t.Helper()

				pubsub := client.Subscribe(t.Context(), redistask.TaskUpdatedChannel)
				t.Cleanup(func() {
					_ = pubsub.Close()
				})

				_, err := pubsub.Receive(t.Context())
				if err != nil {
					t.Fatalf("Failed to receive from pubsub: %v", err)
				}

				return pubsub
			},
			call: func(t *testing.T, client *redis.Client) {
				t.Helper()

				taskPub := redistask.NewTask(client)

				task := internal.Task{
					ID:          "test-123",
					IsDone:      true,
					Description: "Update description",
				}

				if err := taskPub.Updated(t.Context(), task); err != nil {
					t.Fatalf("Failed to publish created event: %v", err)
				}
			},
			verify: func(t *testing.T, msg *redis.Message) {
				t.Helper()

				var got internal.Task

				if err := json.NewDecoder(strings.NewReader(msg.Payload)).Decode(&got); err != nil {
					t.Fatalf("Failed to decode message payload: %v", err)
				}

				task := internal.Task{
					ID:          "test-123",
					IsDone:      true,
					Description: "Update description",
				}

				if diff := cmp.Diff(got, task); diff != "" {
					t.Fatalf("Received task is not the same as the created one: %s", diff)
				}
			},
		},
		{
			name: "Deleted",
			newPubsub: func(t *testing.T, client *redis.Client) *redis.PubSub {
				t.Helper()

				pubsub := client.Subscribe(t.Context(), redistask.TaskDeletedChannel)
				t.Cleanup(func() {
					_ = pubsub.Close()
				})

				_, err := pubsub.Receive(t.Context())
				if err != nil {
					t.Fatalf("Failed to receive from pubsub: %v", err)
				}

				return pubsub
			},
			call: func(t *testing.T, client *redis.Client) {
				t.Helper()

				taskPub := redistask.NewTask(client)

				if err := taskPub.Deleted(t.Context(), "test-delete"); err != nil {
					t.Fatalf("Failed to publish deleted event: %v", err)
				}
			},
			verify: func(t *testing.T, msg *redis.Message) {
				t.Helper()

				var got string

				if err := json.NewDecoder(strings.NewReader(msg.Payload)).Decode(&got); err != nil {
					t.Fatalf("Failed to decode message payload: %v", err)
				}

				id := "test-delete"

				if got != id {
					t.Fatalf("Received task ID is not the same as the deleted one: %s, actual: %s", id, got)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := setupClient()
			if client.err != nil {
				t.Fatalf("Failed to setupClient: %v", client.err)
			}

			pubsub := tt.newPubsub(t, client.redis)

			ch := pubsub.Channel()

			tt.call(t, client.redis)

			time.AfterFunc(time.Second, func() {
				_ = pubsub.Close()
			})

			var (
				count int
				msg   *redis.Message
			)

			for msg = range ch {
				count++
			}

			if count != 1 {
				t.Fatalf("Task event was not received after creation")
			}

			tt.verify(t, msg)
		})
	}
}

//---

var setupClient = sync.OnceValue(func() RedisClient { //nolint: gochecknoglobals
	var res RedisClient

	ctx := context.Background()

	container, err := redistest.Run(ctx, dockerImage)
	if err != nil {
		res.err = fmt.Errorf("failed to start redis container: %w", err)

		return res
	}
	res.container = container

	connStr, err := container.ConnectionString(context.Background())
	if err != nil {
		res.err = fmt.Errorf("failed to get connection string: %w", err)

		return res
	}

	opt, err := redis.ParseURL(connStr)
	if err != nil {
		res.err = fmt.Errorf("failed to parse connection string: %w", err)

		return res
	}

	client := redis.NewClient(opt)
	if status := client.Ping(ctx); status.Err() != nil {
		res.err = fmt.Errorf("failed to ping: %w", err)

		return res
	}
	res.redis = client

	return res
})

type RedisClient struct {
	redis     *redis.Client
	container *redistest.RedisContainer
	err       error
}

func (r *RedisClient) Teardown() error {
	if r.container != nil {
		if err := testcontainers.TerminateContainer(r.container); err != nil {
			return fmt.Errorf("failed to terminate container: %w", err)
		}
	}

	if r.err != nil {
		return r.err
	}

	return nil
}
