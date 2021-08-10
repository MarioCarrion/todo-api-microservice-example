module github.com/MarioCarrion/todo-api

go 1.16

require (
	github.com/bradfitz/gomemcache v0.0.0-20190913173617-a41fca850d0b
	github.com/confluentinc/confluent-kafka-go v1.7.0
	github.com/deepmap/oapi-codegen v1.8.2
	github.com/didip/tollbooth/v6 v6.1.1
	github.com/elastic/go-elasticsearch/v7 v7.14.0
	github.com/getkin/kin-openapi v0.69.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-redis/redis/v8 v8.11.2
	github.com/golang-migrate/migrate/v4 v4.14.1
	github.com/google/go-cmp v0.5.6
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/vault/api v1.1.1
	github.com/jackc/pgx/v4 v4.13.0
	github.com/joho/godotenv v1.3.0
	github.com/mercari/go-circuitbreaker v0.0.0-20201130021310-aff740600e91
	github.com/ory/dockertest/v3 v3.7.0
	github.com/streadway/amqp v1.0.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.20.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.20.0
	go.opentelemetry.io/contrib/instrumentation/runtime v0.20.0
	go.opentelemetry.io/otel v0.20.0
	go.opentelemetry.io/otel/exporters/metric/prometheus v0.20.0
	go.opentelemetry.io/otel/exporters/stdout v0.20.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.20.0
	go.opentelemetry.io/otel/metric v0.20.0
	go.opentelemetry.io/otel/sdk v0.20.0
	go.opentelemetry.io/otel/trace v0.20.0
	go.uber.org/zap v1.18.1
	goa.design/model v1.7.6
)
