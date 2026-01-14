package main

import (
	"context"
	"embed"
	"errors"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/didip/tollbooth/v6"
	"github.com/didip/tollbooth/v6/limiter"
	esv7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/MarioCarrion/todo-api-microservice-example/cmd/internal"
	internaldomain "github.com/MarioCarrion/todo-api-microservice-example/internal"
	"github.com/MarioCarrion/todo-api-microservice-example/internal/elasticsearch"
	"github.com/MarioCarrion/todo-api-microservice-example/internal/envvar"
	"github.com/MarioCarrion/todo-api-microservice-example/internal/memcached"
	"github.com/MarioCarrion/todo-api-microservice-example/internal/postgresql"
	"github.com/MarioCarrion/todo-api-microservice-example/internal/rest"
	"github.com/MarioCarrion/todo-api-microservice-example/internal/service"
)

//go:embed static
var content embed.FS

// MessageBrokerPublisher represents the type that indicates the different Message Brokers supported.
type MessageBrokerPublisher interface {
	Publisher() service.TaskMessageBrokerPublisher
	Close() error
}

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
		return nil, internaldomain.WrapErrorf(err, internaldomain.ErrorCodeUnknown, "zap.NewProduction")
	}

	if err := envvar.Load(env); err != nil {
		return nil, internaldomain.WrapErrorf(err, internaldomain.ErrorCodeUnknown, "envvar.Load")
	}

	vault, err := internal.NewVaultProvider()
	if err != nil {
		return nil, internaldomain.WrapErrorf(err, internaldomain.ErrorCodeUnknown, "internal.NewVaultProvider")
	}

	conf := envvar.New(vault)

	//-

	pool, err := internal.NewPostgreSQL(conf)
	if err != nil {
		return nil, internaldomain.WrapErrorf(err, internaldomain.ErrorCodeUnknown, "internal.NewPostgreSQL")
	}

	esClient, err := internal.NewElasticSearch(conf)
	if err != nil {
		return nil, internaldomain.WrapErrorf(err, internaldomain.ErrorCodeUnknown, "internal.NewElasticSearch")
	}

	memcached, err := internal.NewMemcached(conf)
	if err != nil {
		return nil, internaldomain.WrapErrorf(err, internaldomain.ErrorCodeUnknown, "internal.NewMemcached")
	}

	msgBroker, err := NewMessageBrokerPublisher(conf)
	if err != nil {
		return nil, internaldomain.WrapErrorf(err, internaldomain.ErrorCodeUnknown, "NewMessageBroker")
	}

	//-

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

	srv := newServer(serverConfig{
		Address:       address,
		DB:            pool,
		ElasticSearch: esClient,
		Middlewares:   []rest.MiddlewareFunc{logging},
		Logger:        logger,
		Memcached:     memcached,
		MessageBroker: msgBroker,
	})

	errC := make(chan error, 1)

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-ctx.Done()

		logger.Info("Shutdown signal received")

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer func() {
			_ = logger.Sync()

			pool.Close()
			srv.Close()
			_ = msgBroker.Close()

			stop()
			cancel()
			close(errC)
		}()

		srv.SetKeepAlivesEnabled(false)

		if err := srv.Shutdown(ctxTimeout); err != nil { //nolint: contextcheck
			errC <- err
		}

		logger.Info("Shutdown completed")
	}()

	go func() {
		logger.Info("Listening and serving", zap.String("address", address))

		// "ListenAndServe always returns a non-nil error. After Shutdown or Close, the returned error is
		// ErrServerClosed."
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errC <- err
		}
	}()

	return errC, nil
}

type serverConfig struct {
	Address       string
	DB            *pgxpool.Pool
	ElasticSearch *esv7.Client
	Memcached     *memcache.Client
	Middlewares   []rest.MiddlewareFunc
	Logger        *zap.Logger
	MessageBroker MessageBrokerPublisher
}

func newServer(conf serverConfig) *http.Server {
	repo := postgresql.NewTask(conf.DB)
	mrepo := memcached.NewTask(conf.Memcached, repo, conf.Logger)

	search := elasticsearch.NewTask(conf.ElasticSearch)
	msearch := memcached.NewSearchableTask(conf.Memcached, search)

	svc := service.NewTask(conf.Logger, mrepo, msearch, conf.MessageBroker.Publisher())

	taskHandler := rest.NewTaskHandler(svc)

	router := http.NewServeMux()

	fsys, _ := fs.Sub(content, "static")
	router.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(fsys))))

	strictHandler := rest.NewStrictHandler(taskHandler, nil)

	options := rest.StdHTTPServerOptions{
		BaseRouter:  router,
		Middlewares: conf.Middlewares,
		ErrorHandlerFunc: func(w http.ResponseWriter, _ *http.Request, err error) {
			switch {
			case errors.Is(err, context.Canceled):
				// Client canceled the request; treat as a bad request from the client's perspective.
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			case errors.Is(err, context.DeadlineExceeded):
				// Request timed out; indicate a gateway timeout.
				http.Error(w, http.StatusText(http.StatusGatewayTimeout), http.StatusGatewayTimeout)
			default:
				// Log internal error details but do not expose them to the client.
				if conf.Logger != nil {
					conf.Logger.Error("request failed", zap.Error(err))
				}
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		},
	}

	handler := rest.HandlerWithOptions(strictHandler, options)

	//-

	lmt := tollbooth.NewLimiter(3, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Second})

	lmtmw := tollbooth.LimitHandler(lmt, handler)

	//-

	return &http.Server{
		Handler:           lmtmw,
		Addr:              conf.Address,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       1 * time.Second,
	}
}
