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
	"time"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"

	"github.com/MarioCarrion/todo-api/internal/envvar"
	"github.com/MarioCarrion/todo-api/internal/envvar/vault"
	"github.com/MarioCarrion/todo-api/internal/postgresql"
	"github.com/MarioCarrion/todo-api/internal/rest"
	"github.com/MarioCarrion/todo-api/internal/service"
)

//go:embed static
var content embed.FS

func main() {
	var env string

	flag.StringVar(&env, "env", "", "Environment Variables filename")
	flag.Parse()

	if err := envvar.Load(env); err != nil {
		log.Fatalln("Couldn't load configuration", err)
	}

	conf := envvar.New(newVaultProvider())

	//-

	promExporter := initTracer(conf)

	defer func() {
		_ = initMeter().Stop(context.Background())
	}()

	if err := runtime.Start(
		runtime.WithMinimumReadMemStatsInterval(time.Second),
	); err != nil {
		panic(err)
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	middleware := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info(r.Method,
				zap.Time("time", time.Now()),
				zap.String("url", r.URL.String()),
			)

			h.ServeHTTP(w, r)
		})
	}

	//-

	db := newDB(conf)
	defer db.Close()

	//-

	repo := postgresql.NewTask(db) // Task Repository
	svc := service.NewTask(repo)   // Task Application Service

	//-

	r := mux.NewRouter()
	r.Handle("/metrics", promExporter)

	r.Use(otelmux.Middleware("todo-api-server"))

	//-

	rest.RegisterOpenAPI(r)
	rest.NewTaskHandler(svc).Register(r)

	fsys, _ := fs.Sub(content, "static")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(fsys))))

	//-

	address := "0.0.0.0:9234"

	srv := &http.Server{
		Handler:           middleware(r),
		Addr:              address,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       1 * time.Second,
	}

	log.Println("Starting server", address)

	log.Fatal(srv.ListenAndServe())
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

//-

func initMeter() *controller.Controller {
	pusher, err := stdout.InstallNewPipeline([]stdout.Option{
		stdout.WithPrettyPrint(),
	}, nil)
	if err != nil {
		log.Panicf("Couldn't initialize metric stdout exporter %v", err)
	}

	return pusher
}

func initTracer(conf *envvar.Configuration) *prometheus.Exporter {
	promExporter, err := prometheus.NewExportPipeline(prometheus.Config{})
	if err != nil {
		log.Fatalln("Couldn't initialize Prometheus exporter", err)
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
		log.Fatalln("Couldn't initialize jaeger", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSyncer(jaegerExporter),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return promExporter
}
