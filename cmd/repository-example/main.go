package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/MarioCarrion/todo-api/internal"
	"github.com/MarioCarrion/todo-api/internal/envvar"
	"github.com/MarioCarrion/todo-api/internal/envvar/vault"
	"github.com/MarioCarrion/todo-api/internal/postgresql"
	"github.com/MarioCarrion/todo-api/internal/service"
)

func main() {
	var env string

	flag.StringVar(&env, "env", "", "Environment Variables filename")
	flag.Parse()

	if err := envvar.Load(env); err != nil {
		log.Fatalln("Couldn't load configuration", err)
	}

	conf := envvar.New(newVaultProvider())

	//-

	db := newDB(conf)
	defer db.Close()

	//-

	repo := postgresql.NewTask(db) // Task Repository
	svc := service.NewTask(repo)   // Task Application Service
	task, err := svc.Create(context.Background(), "new task", internal.PriorityLow, internal.Dates{})

	fmt.Printf("NEW task %#v, err %s\n", task, err)

	//-

	if err := svc.Update(context.Background(),
		task.ID,
		"changed task",
		internal.PriorityHigh,
		internal.Dates{
			Due: time.Now().Add(2 * time.Hour),
		},
		false); err != nil {
		log.Fatalln("couldn't update task", err)
	}

	updatedTask, err := svc.Task(context.Background(), task.ID)
	if err != nil {
		log.Fatalln("couldn't find task", err)
	}

	fmt.Printf("UPDATED task %#v, err %s\n", updatedTask, err)
}

func newDB(conf *envvar.Configuration) *sql.DB {
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
		log.Fatalln("Couldn't open DB", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalln("Couldn't ping DB", err)
	}

	return db
}

func newVaultProvider() *vault.Provider {
	// XXX: We will revisit this code in future episodes replacing it with another solution
	vaultPath := os.Getenv("VAULT_PATH")
	vaultToken := os.Getenv("VAULT_TOKEN")
	vaultAddress := os.Getenv("VAULT_ADDRESS")
	// XXX: -

	provider, err := vault.New(vaultToken, vaultAddress, vaultPath)
	if err != nil {
		log.Fatalln("Couldn't load provider", err)
	}

	return provider
}
