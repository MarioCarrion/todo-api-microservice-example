package internal

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"

	// Initialize "pgx".
	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/MarioCarrion/todo-api/internal"
	"github.com/MarioCarrion/todo-api/internal/envvar"
)

// NewPostgreSQL instantiates the PostgreSQL database using configuration defined in environment variables.
func NewPostgreSQL(conf *envvar.Configuration) (*sql.DB, error) {
	get := func(v string) string {
		res, err := conf.Get(v)
		if err != nil {
			log.Fatalf("Couldn't get configuration value for %s: %s", v, err)
		}

		return res
	}

	// XXX: We will revisit this code in future episodes replacing it with another solution
	databaseHost := get("DATABASE_HOST")
	databasePort := get("DATABASE_PORT")
	databaseUsername := get("DATABASE_USERNAME")
	databasePassword := get("DATABASE_PASSWORD")
	databaseName := get("DATABASE_NAME")
	databaseSSLMode := get("DATABASE_SSLMODE")
	// XXX: -

	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(databaseUsername, databasePassword),
		Host:   fmt.Sprintf("%s:%s", databaseHost, databasePort),
		Path:   databaseName,
	}

	q := dsn.Query()
	q.Add("sslmode", databaseSSLMode)

	dsn.RawQuery = q.Encode()

	db, err := sql.Open("pgx", dsn.String())
	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "sql.Open")
	}

	if err := db.Ping(); err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "db.Ping")
	}

	return db, nil
}
