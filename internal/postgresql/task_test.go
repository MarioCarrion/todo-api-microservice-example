package postgresql_test

import (
	"context"
	"errors"
	"os"
	"path"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/tern/v2/migrate"
	"github.com/testcontainers/testcontainers-go"
	tpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
	"github.com/MarioCarrion/todo-api-microservice-example/internal/postgresql"
)

func TestTask_Create(t *testing.T) {
	t.Parallel()

	t.Run("Create: OK", func(t *testing.T) {
		t.Parallel()

		task, err := postgresql.NewTask(newDB(t)).Create(t.Context(),
			internal.CreateParams{
				Description: "test",
				Priority:    internal.ValueToPointer(internal.PriorityNone),
			})
		if err != nil {
			t.Fatalf("expected no error, got %s", err)
		}

		if task.ID == "" {
			t.Fatalf("expected valid record, got empty value")
		}
	})

	t.Run("Create: ERR", func(t *testing.T) {
		t.Parallel()

		_, err := postgresql.NewTask(newDB(t)).Create(t.Context(),
			internal.CreateParams{
				Description: "",
				Priority:    internal.ValueToPointer(internal.Priority(-1)),
			})
		if err == nil { // because of invalid priority
			t.Fatalf("expected error, got no value")
		}

		var ierr *internal.Error
		if !errors.As(err, &ierr) || ierr.Code() != internal.ErrorCodeUnknown {
			t.Fatalf("expected %T error, got %T : %v", ierr, err, err)
		}
	})
}

func TestTask_Delete(t *testing.T) {
	t.Parallel()

	t.Run("Delete: OK", func(t *testing.T) {
		t.Parallel()

		store := postgresql.NewTask(newDB(t))

		createdTask, err := store.Create(t.Context(), internal.CreateParams{
			Description: "test",
			Priority:    internal.ValueToPointer(internal.PriorityNone),
		})
		if err != nil {
			t.Fatalf("expected no error, got %s", err)
		}

		if err := store.Delete(t.Context(), createdTask.ID); err != nil {
			t.Fatalf("expected no error, got %s", err)
		}

		if _, err = store.Find(t.Context(), createdTask.ID); !errors.Is(err, pgx.ErrNoRows) {
			t.Fatalf("expected no error, got %s", err)
		}
	})

	t.Run("Update: ERR uuid", func(t *testing.T) {
		t.Parallel()

		err := postgresql.NewTask(newDB(t)).Delete(t.Context(), "x")

		if err == nil {
			t.Fatalf("expected error, got not value")
		}

		var ierr *internal.Error
		if !errors.As(err, &ierr) || ierr.Code() != internal.ErrorCodeInvalidArgument {
			t.Fatalf("expected %T error, got %T : %v", ierr, err, err)
		}
	})

	t.Run("Delete: ERR not found", func(t *testing.T) {
		t.Parallel()

		err := postgresql.NewTask(newDB(t)).Delete(t.Context(), "44633fe3-b039-4fb3-a35f-a57fe3c906c7")

		var ierr *internal.Error
		if !errors.As(err, &ierr) || ierr.Code() != internal.ErrorCodeNotFound {
			t.Fatalf("expected %T error, got %T : %v", ierr, err, err)
		}
	})
}

func TestTask_Find(t *testing.T) {
	t.Parallel()

	t.Run("Find: OK", func(t *testing.T) {
		t.Parallel()

		now := time.Now().UTC().Truncate(time.Minute)

		store := postgresql.NewTask(newDB(t))

		originalTask, err := store.Create(t.Context(), internal.CreateParams{
			Description: "test",
			Priority:    internal.ValueToPointer(internal.PriorityNone),
			Dates: &internal.Dates{
				Start: internal.ValueToPointer(now),
				Due:   internal.ValueToPointer(now),
			},
		})
		if err != nil {
			t.Fatalf("expected no error, got %s", err)
		}

		actualTask, err := store.Find(t.Context(), originalTask.ID)
		if err != nil {
			t.Fatalf("expected no error, got %s", err)
		}

		if !cmp.Equal(originalTask, actualTask) {
			t.Fatalf("expected result does not match: %s", cmp.Diff(originalTask, actualTask))
		}
	})

	t.Run("Find: ERR uuid", func(t *testing.T) {
		t.Parallel()

		_, err := postgresql.NewTask(newDB(t)).Find(t.Context(), "x")
		if err == nil {
			t.Fatalf("expected error, got not value")
		}

		var ierr *internal.Error
		if !errors.As(err, &ierr) || ierr.Code() != internal.ErrorCodeInvalidArgument {
			t.Fatalf("expected %T error, got %T : %v", ierr, err, err)
		}
	})

	t.Run("Find: ERR not found", func(t *testing.T) {
		t.Parallel()

		_, err := postgresql.NewTask(newDB(t)).Find(t.Context(), "44633fe3-b039-4fb3-a35f-a57fe3c906c7")
		if err == nil {
			t.Fatalf("expected error, got not value")
		}

		var ierr *internal.Error
		if !errors.As(err, &ierr) || ierr.Code() != internal.ErrorCodeNotFound {
			t.Fatalf("expected %T error, got %T : %v", ierr, err, err)
		}
	})
}

func TestTask_Update(t *testing.T) {
	t.Parallel()

	t.Run("Update: OK", func(t *testing.T) {
		t.Parallel()

		store := postgresql.NewTask(newDB(t))

		originalTask, err := store.Create(t.Context(), internal.CreateParams{
			Description: "test",
			Priority:    internal.ValueToPointer(internal.PriorityNone),
		})
		if err != nil {
			t.Fatalf("expected no error, got %s", err)
		}

		now := time.Now().UTC().Truncate(time.Minute)

		originalTask.Description = "changed"
		originalTask.Dates = &internal.Dates{
			Start: internal.ValueToPointer(now),
			Due:   internal.ValueToPointer(now),
		}
		originalTask.Priority = internal.ValueToPointer(internal.PriorityHigh)

		params := internal.UpdateParams{
			Description: &originalTask.Description,
			Priority:    originalTask.Priority,
			Dates:       originalTask.Dates,
			IsDone:      &originalTask.IsDone,
		}

		if err := store.Update(t.Context(), originalTask.ID, params); err != nil {
			t.Fatalf("expected no error, got %s", err)
		}

		actualTask, err := store.Find(t.Context(), originalTask.ID)
		if err != nil {
			t.Fatalf("expected no error, got %s", err)
		}

		opts := cmp.Comparer(func(x, y time.Time) bool {
			return x.Unix() == y.Unix()
		})

		if !cmp.Equal(originalTask, actualTask, opts) {
			t.Fatalf("expected result does not match: %s", cmp.Diff(originalTask, actualTask))
		}
	})

	t.Run("Update: ERR uuid", func(t *testing.T) {
		t.Parallel()

		params := internal.UpdateParams{
			Description: internal.ValueToPointer("x"),
			Priority:    internal.ValueToPointer(internal.PriorityNone),
			Dates:       &internal.Dates{},
			IsDone:      new(bool),
		}

		err := postgresql.NewTask(newDB(t)).Update(t.Context(), "x", params)
		if err == nil {
			t.Fatalf("expected error, got not value")
		}

		var ierr *internal.Error
		if !errors.As(err, &ierr) || ierr.Code() != internal.ErrorCodeInvalidArgument {
			t.Fatalf("expected %T error, got %T : %v", ierr, err, err)
		}
	})

	t.Run("Update: ERR invalid priority", func(t *testing.T) {
		t.Parallel()

		store := postgresql.NewTask(newDB(t))
		dates := internal.Dates{}

		task, err := store.Create(t.Context(), internal.CreateParams{
			Description: "test",
			Priority:    internal.ValueToPointer(internal.PriorityNone),
			Dates:       &dates,
		})
		if err != nil {
			t.Fatalf("expected no error, got %s", err)
		}

		params := internal.UpdateParams{
			Priority: internal.ValueToPointer(internal.Priority(-1)),
			Dates:    &internal.Dates{},
			IsDone:   new(bool),
		}

		err = postgresql.NewTask(newDB(t)).Update(t.Context(),
			task.ID,
			params,
		)
		if err == nil {
			t.Fatalf("expected error, got not value")
		}

		var ierr *internal.Error
		if !errors.As(err, &ierr) || ierr.Code() != internal.ErrorCodeUnknown {
			t.Fatalf("expected %T error, got %T : %v", ierr, err, err)
		}
	})

	t.Run("Update: ERR not found", func(t *testing.T) {
		t.Parallel()

		params := internal.UpdateParams{
			Priority: internal.ValueToPointer(internal.PriorityNone),
			Dates:    &internal.Dates{},
			IsDone:   new(bool),
		}

		err := postgresql.NewTask(newDB(t)).Update(t.Context(),
			"44633fe3-b039-4fb3-a35f-a57fe3c906c7",
			params)
		if err == nil {
			t.Fatalf("expected error, got not value")
		}

		var ierr *internal.Error
		if !errors.As(err, &ierr) || ierr.Code() != internal.ErrorCodeNotFound {
			t.Fatalf("expected %T error, got %T : %v", ierr, err, err)
		}
	})
}

func newDB(tb testing.TB) *pgxpool.Pool {
	tb.Helper()

	const (
		username = "username"
		password = "password"
		dbName   = "todo"
	)

	//- Run container and verify it works as expected

	container, err := tpostgres.Run(tb.Context(),
		"postgres:17.4-bookworm",
		tpostgres.WithDatabase(dbName),
		tpostgres.WithUsername(username),
		tpostgres.WithPassword(password),
		tpostgres.BasicWaitStrategies(),
	)
	if err != nil {
		tb.Fatalf("Failed to run container: %s", err)
	}

	tb.Cleanup(func() {
		if err := testcontainers.TerminateContainer(container); err != nil {
			tb.Logf("Failed to terminate container: %s", err)
		}
	})

	host, err := container.ConnectionString(tb.Context())
	if err != nil {
		tb.Fatalf("Failed to get host address: %s", err)
	}

	ctx, cFunc := context.WithDeadline(tb.Context(), time.Now().Add(500*time.Millisecond))
	tb.Cleanup(cFunc)

	tb.Logf("Waiting for database to be ready at: %s", host)

	db, err := pgx.Connect(ctx, host)
	if err != nil {
		tb.Fatalf("Failed to connect to database: %s", err)
	}

	tb.Cleanup(func() {
		_ = db.Close(ctx)
	})

	//- DB Migrations

	migrator, err := migrate.NewMigrator(ctx, db, "public.schema_version")
	if err != nil {
		tb.Fatalf("Failed to migrate (1): %s", err)
	}

	err = migrator.LoadMigrations(os.DirFS(path.Join("..", "..", "db", "migrations")))
	if err != nil {
		tb.Fatalf("Failed to migrate (2): %s", err)
	}

	if err = migrator.Migrate(tb.Context()); err != nil {
		tb.Fatalf("Failed to migrate (3): %s", err)
	}

	//- Initialize DB Pool

	dbpool, err := pgxpool.New(tb.Context(), host)
	if err != nil {
		tb.Fatalf("Failed to open DB Pool: %s", err)
	}

	tb.Cleanup(dbpool.Close)

	return dbpool
}
