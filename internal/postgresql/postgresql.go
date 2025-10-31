package postgresql

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
	"github.com/MarioCarrion/todo-api-microservice-example/internal/postgresql/db"
)

//go:generate sqlc generate

func convertPriority(priority db.Priority) (internal.Priority, error) {
	switch priority {
	case db.PriorityNone:
		return internal.PriorityNone, nil
	case db.PriorityLow:
		return internal.PriorityLow, nil
	case db.PriorityMedium:
		return internal.PriorityMedium, nil
	case db.PriorityHigh:
		return internal.PriorityHigh, nil
	}

	return internal.Priority(-1), fmt.Errorf("unknown value: %s", priority)
}

func newTimestamp(t *time.Time) pgtype.Timestamp {
	if t == nil {
		return pgtype.Timestamp{
			Valid: false,
		}
	}

	return pgtype.Timestamp{
		Time:  *t,
		Valid: !t.IsZero(),
	}
}

func newPriority(p *internal.Priority) db.Priority {
	if p == nil {
		return db.PriorityNone
	}

	switch *p {
	case internal.PriorityNone:
		return db.PriorityNone
	case internal.PriorityLow:
		return db.PriorityLow
	case internal.PriorityMedium:
		return db.PriorityMedium
	case internal.PriorityHigh:
		return db.PriorityHigh
	}

	// XXX: because we are using an enum type, postgres will fail with the following value.

	return "invalid"
}
