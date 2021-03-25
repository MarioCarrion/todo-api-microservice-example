module github.com/MarioCarrion/todo-api

go 1.16

require (
	github.com/deepmap/oapi-codegen v1.5.1
	github.com/getkin/kin-openapi v0.37.0
	github.com/ghodss/yaml v1.0.0
	github.com/golang-migrate/migrate/v4 v4.14.1
	github.com/google/go-cmp v0.5.5
	github.com/google/uuid v1.2.0
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/vault/api v1.0.4
	github.com/jackc/pgx/v4 v4.10.1
	github.com/joho/godotenv v1.3.0
	github.com/ory/dockertest/v3 v3.6.3
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.19.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.19.0
	go.opentelemetry.io/contrib/instrumentation/runtime v0.19.0
	go.opentelemetry.io/otel v0.19.0
	go.opentelemetry.io/otel/exporters/metric/prometheus v0.19.0
	go.opentelemetry.io/otel/exporters/stdout v0.19.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.19.0
	go.opentelemetry.io/otel/metric v0.19.0
	go.opentelemetry.io/otel/sdk v0.19.0
	go.opentelemetry.io/otel/sdk/metric v0.19.0
	go.opentelemetry.io/otel/trace v0.19.0
	go.uber.org/zap v1.16.0
)
