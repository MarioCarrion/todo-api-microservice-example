package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/streadway/amqp"
	"go.uber.org/zap"

	// "go.opentelemetry.io/otel/api/correlation"

	"github.com/MarioCarrion/todo-api/cmd/internal"
	internaldomain "github.com/MarioCarrion/todo-api/internal"
	"github.com/MarioCarrion/todo-api/internal/elasticsearch"
	"github.com/MarioCarrion/todo-api/internal/envvar"
)

const rabbitMQConsumerName = "elasticsearch-indexer"

func main() {
	var env string

	flag.StringVar(&env, "env", "", "Environment Variables filename")
	flag.Parse()

	errC, err := run(env)
	if err != nil {
		log.Fatalf("Couldn't run: %s", err)
	}

	if err := <-errC; err != nil {
		log.Fatalf("Error while running: %s", err)
	}
}

func run(env string) (<-chan error, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("zap.NewProduction %w", err)
	}

	if err := envvar.Load(env); err != nil {
		return nil, fmt.Errorf("envvar.Load %w", err)
	}

	vault, err := internal.NewVaultProvider()
	if err != nil {
		return nil, fmt.Errorf("internal.NewVaultProvider %w", err)
	}

	conf := envvar.New(vault)

	//-

	es, err := internal.NewElasticSearch(conf)
	if err != nil {
		return nil, fmt.Errorf("internal.NewElasticSearch %w", err)
	}

	rmq, err := internal.NewRabbitMQ(conf)
	if err != nil {
		return nil, fmt.Errorf("newRabbitMQ %w", err)
	}

	//-

	_, err = internal.NewOTExporter(conf)
	if err != nil {
		return nil, fmt.Errorf("newOTExporter %w", err)
	}

	//-

	srv := &Server{
		logger: logger,
		rmq:    rmq,
		task:   elasticsearch.NewTask(es),
		done:   make(chan struct{}),
	}

	errC := make(chan error, 1)

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-ctx.Done()

		logger.Info("Shutdown signal received")

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		defer func() {
			logger.Sync()
			rmq.Close()
			stop()
			cancel()
			close(errC)
		}()

		if err := srv.Shutdown(ctxTimeout); err != nil {
			errC <- err
		}

		logger.Info("Shutdown completed")
	}()

	go func() {
		logger.Info("Listening and serving")

		if err := srv.ListenAndServe(); err != nil {
			errC <- err
		}
	}()

	return errC, nil
}

type Server struct {
	logger *zap.Logger
	rmq    *internal.RabbitMQ
	task   *elasticsearch.Task
	queue  amqp.Queue
	done   chan struct{}
}

// ListenAndServe ...
func (s *Server) ListenAndServe() error {
	// XXX: Dead Letter Exchange will be implemented in future episodes
	q, err := s.rmq.Channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("channel.QueueDeclare %w", err)
	}

	err = s.rmq.Channel.QueueBind(
		q.Name,          // queue name
		"tasks.event.*", // routing key
		"tasks",         // exchange
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("channel.QueueBind %w", err)
	}

	msgs, err := s.rmq.Channel.Consume(
		q.Name,               // queue
		rabbitMQConsumerName, // consumer
		false,                // auto-ack
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
	)
	if err != nil {
		return fmt.Errorf("channel.Consume %w", err)
	}

	go func() {
		for msg := range msgs {
			s.logger.Info(fmt.Sprintf("Received message: %s", msg.RoutingKey))

			var nack bool

			// XXX: Instrumentation to be added in a future episode

			// XXX: We will revisit defining these topics in a better way in future episodes
			switch msg.RoutingKey {
			case "tasks.event.updated", "tasks.event.created":
				task, err := decodeTask(msg.Body)
				if err != nil {
					nack = true
					return
				}

				if err := s.task.Index(context.Background(), task); err != nil {
					nack = true
				}
			case "tasks.event.deleted":
				id, err := decodeID(msg.Body)

				if err != nil {
					nack = true
					return
				}

				if err := s.task.Delete(context.Background(), id); err != nil {
					nack = true
				}
			default:
				nack = true
			}

			if nack {
				s.logger.Info("NAcking :(")
				msg.Nack(false, nack)
			} else {
				s.logger.Info("Acking :)")
				msg.Ack(false)
			}
		}

		s.logger.Info("No more messages to consume. Exiting.")

		s.done <- struct{}{}
	}()

	return nil
}

// Shutdown ...
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server")

	s.rmq.Channel.Cancel(rabbitMQConsumerName, false)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context.Done: %w", ctx.Err())

		case <-s.done:
			return nil
		}
	}
}

func decodeTask(b []byte) (internaldomain.Task, error) {
	var res internaldomain.Task

	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&res); err != nil {
		return internaldomain.Task{}, internaldomain.WrapErrorf(err, internaldomain.ErrorCodeUnknown, "gob.Decode")
	}

	return res, nil
}

func decodeID(b []byte) (string, error) {
	var res string

	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&res); err != nil {
		return "", internaldomain.WrapErrorf(err, internaldomain.ErrorCodeUnknown, "gob.Decode")
	}

	return res, nil
}
