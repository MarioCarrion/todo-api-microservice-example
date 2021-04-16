package main

import (
	"context"
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	esv7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"

	"github.com/MarioCarrion/todo-api/internal/elasticsearch"
	"github.com/MarioCarrion/todo-api/internal/envvar"
	"github.com/MarioCarrion/todo-api/internal/envvar/vault"
	"github.com/MarioCarrion/todo-api/internal/postgresql"
	"github.com/MarioCarrion/todo-api/internal/rest"
	"github.com/MarioCarrion/todo-api/internal/service"
)

//go:embed static
var content embed.FS

func main() {
	var env, address string

	flag.StringVar(&env, "env", "", "Environment Variables filename")
	flag.StringVar(&address, "address", ":9234", "HTTP Server Address")
	flag.Parse()

	errC, err := run(env, address)
	if err != nil {
		log.Fatalf("Couldn't run: %s", err)
	}

	if err := <-errC; err != nil {
		log.Fatalf("Error while running: %s", err)
	}
}

func run(env, address string) (<-chan error, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("zap.NewProduction %w", err)
	}

	if err := envvar.Load(env); err != nil {
		return nil, fmt.Errorf("envvar.Load %w", err)
	}

	vault, err := newVaultProvider()
	if err != nil {
		return nil, fmt.Errorf("newVaultProvider %w", err)
	}

	conf := envvar.New(vault)

	//-

	db, err := newDB(conf)
	if err != nil {
		return nil, fmt.Errorf("newDB %w", err)
	}

	es, err := newElasticSearch(conf)
	if err != nil {
		return nil, fmt.Errorf("newElasticSearch %w", err)
	}

	//-

	promExporter, err := newOTExporter(conf)
	if err != nil {
		return nil, fmt.Errorf("newOTExporter %w", err)
	}

	logging := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info(r.Method,
				zap.Time("time", time.Now()),
				zap.String("url", r.URL.String()),
			)

			h.ServeHTTP(w, r)
		})
	}

	//-

	errC := make(chan error, 1)

	srv := newServer(address, db, es, promExporter, otelmux.Middleware("todo-api-server"), logging)

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	// XXX: When using Go 1.15 or older
	// sc := make(chan os.Signal, 1)
	// signal.Notify(sc,
	// 	os.Interrupt,
	// 	syscall.SIGTERM,
	// 	syscall.SIGQUIT)

	go func() {
		// <-sc // XXX: When using Go 1.15 or older
		<-ctx.Done()

		logger.Info("Shutdown signal received")

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer func() {
			logger.Sync()
			db.Close()
			stop()
			cancel()
			close(errC)
		}()

		srv.SetKeepAlivesEnabled(false)

		if err := srv.Shutdown(ctxTimeout); err != nil {
			errC <- err
		}

		logger.Info("Shutdown completed")
	}()

	go func() {
		logger.Info("Listening and serving", zap.String("address", address))

		// "ListenAndServe always returns a non-nil error. After Shutdown or Close, the returned error is
		// ErrServerClosed."
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errC <- err
		}
	}()

	return errC, nil
}

func newServer(address string, db *sql.DB, es *esv7.Client, metrics http.Handler, mws ...mux.MiddlewareFunc) *http.Server {
	r := mux.NewRouter()

	for _, mw := range mws {
		r.Use(mw)
	}

	//-

	repo := postgresql.NewTask(db)
	search := elasticsearch.NewTask(es)
	svc := service.NewTask(repo, search)

	rest.RegisterOpenAPI(r)
	rest.NewTaskHandler(svc).Register(r)

	//-

	fsys, _ := fs.Sub(content, "static")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(fsys))))

	r.Handle("/metrics", metrics)

	//-

	return &http.Server{
		Handler:           r,
		Addr:              address,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       1 * time.Second,
	}
}

func newDB(conf *envvar.Configuration) (*sql.DB, error) {
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
		return nil, fmt.Errorf("sql.Open %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping %w", err)
	}

	return db, nil
}

func newVaultProvider() (*vault.Provider, error) {
	// XXX: We will revisit this code in future episodes replacing it with another solution
	vaultPath := os.Getenv("VAULT_PATH")
	vaultToken := os.Getenv("VAULT_TOKEN")
	vaultAddress := os.Getenv("VAULT_ADDRESS")
	// XXX: -

	provider, err := vault.New(vaultToken, vaultAddress, vaultPath)
	if err != nil {
		return nil, fmt.Errorf("vault.New %w", err)
	}

	return provider, nil
}

//-

func newOTExporter(conf *envvar.Configuration) (*prometheus.Exporter, error) {
	if err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second)); err != nil {
		return nil, fmt.Errorf("runtime.Start %w", err)
	}

	promExporter, err := prometheus.NewExportPipeline(prometheus.Config{})
	if err != nil {
		return nil, fmt.Errorf("prometheus.NewExportPipeline %w", err)
	}

	global.SetMeterProvider(promExporter.MeterProvider())

	//-

	jaegerEndpoint, _ := conf.Get("JAEGER_ENDPOINT")

	jaegerExporter, err := jaeger.NewRawExporter(
		jaeger.WithCollectorEndpoint(jaegerEndpoint),
		jaeger.WithSDKOptions(sdktrace.WithSampler(sdktrace.AlwaysSample())),
		jaeger.WithProcessFromEnv(),
	)
	if err != nil {
		return nil, fmt.Errorf("jaeger.NewRawExporter %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSyncer(jaegerExporter),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return promExporter, nil
}

func newElasticSearch(conf *envvar.Configuration) (*esv7.Client, error) {
	es, err := esv7.NewDefaultClient()
	if err != nil {
		return nil, fmt.Errorf("elasticsearch.Open %w", err)
	}

	res, err := es.Info()
	if err != nil {
		return nil, fmt.Errorf("es.Info %w", err)
	}
	defer res.Body.Close()

	return es, nil
}
