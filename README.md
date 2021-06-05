# "ToDo API" Microservice Example

## Introduction

Welcome! 👋

This is an educational repository that includes a microservice written in Go. It is used as the principal example of my video series: [Building Microservices in Go](https://www.youtube.com/playlist?list=PL7yAAGMOat_Fn8sAXIk0WyBfK_sT1pohu).

This repository **is not** a template **nor** a framework, it's a **collection of patterns and guidelines** I've successfully used to deliver enterprise microservices when using Go, and just like with everything in Software Development some trade-offs were made.

My end goal with this project is to help you learn another way to structure your Go project with 3 final goals:

1. It is _enterprise_, meant to last for years,
1. It allows a team to collaborate efficiently with little friction, and
1. It is as idiomatic as possible.

Join the fun at [https://youtube.com/MarioCarrion](https://www.youtube.com/c/MarioCarrion).

## Domain Driven Design

This project uses a lot of the ideas introduced by Eric Evans in his book [Domain Driven Design](https://www.domainlanguage.com/), I do encourage reading that book but before I think reading [Domain-Driven Design Distilled](https://smile.amazon.com/Domain-Driven-Design-Distilled-Vaughn-Vernon/dp/0134434420/) makes more sense, also there's a free to download [DDD Reference](https://www.domainlanguage.com/ddd/reference/) available as well.

On YouTube I created [a playlist](https://www.youtube.com/playlist?list=PL7yAAGMOat_GJqfTdM9PBdTRSH7jXs6mI) that includes some of my favorite talks and webinars, feel free to explore that as well.

## Project Structure

Talking specifically about microservices **only**, the structure I like to recommend is the following, everything using `<` and `>` depends on the domain being implemented and the bounded context being defined.

- [ ] `build/`: defines the code used for creating infrastructure as well as docker containers.
  - [ ] `<cloud-providers>/`: define concrete cloud provider.
  - [ ] `<executableN>/`: contains a Dockerfile used for building the binary.
- [ ] `cmd/`
  - [ ] `<primary-server>/`: uses primary database.
  - [ ] `<replica-server>/`: uses readonly databases.
  - [ ] `<binaryN>/`
- [X] `db/`
  - [X] `migrations/`: contains database migrations.
  - [ ] `seeds/`: contains file meant to populate basic database values.
- [ ] `internal/`: defines the _core domain_.
  - [ ] `<datastoreN>/`: a concrete _repository_ used by the domain, for example `postgresql`
  - [ ] `http/`: defines HTTP Handlers.
  - [ ] `service/`: orchestrates use cases and manages transactions.
- [X] `pkg/` public API meant to be imported by other Go package.

There are cases where requiring a new bounded context is needed, in those cases the recommendation would be to
define a package like `internal/<bounded-context>` that then should follow the same structure, for example:

* `internal/<bounded-context>/`
  * `internal/<bounded-context>/<datastoreN>`
  * `internal/<bounded-context>/http`
  * `internal/<bounded-context>/service`

## Tools

```
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.14.1
go install github.com/kyleconroy/sqlc/cmd/sqlc@v1.6.0
go install github.com/maxbrunsfeld/counterfeiter/v6@v6.3.0
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.5.1
go install goa.design/model/cmd/mdl@v1.7.6
go install goa.design/model/cmd/stz@v1.7.6
```

## Features

Icons meaning:

* <img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video"> means a link to Youtube video.
* <img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/link.svg" width="20" height="20" alt="Blog post"> means a link to Blog post.

In no particular order:

- [X] Project Layout [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/LUvid5TJ81Y) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/link.svg" width="20" height="20" alt="Blog post">](https://mariocarrion.com/2021/03/21/golang-microservices-domain-driven-design-project-layout.html)
- [X] Dependency Injection [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/Z89UU4vSayY)
- [X] [Secure Configuration](docs/SECURE\_CONFIGURATION.md)
  - [X] Using [Hashicorp Vault](https://www.hashicorp.com/products/vault) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/7UmJR0dOkjM)
  - [X] Using [AWS SSM](https://aws.amazon.com/systems-manager/features/#Parameter_Store) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/link.svg" width="20" height="20" alt="Blog post">](https://mariocarrion.com/2019/11/20/golang-aws-ssm-env-vars-package.html)
- [ ] Infrastructure as code
- [X] [Metrics, Traces and Logging using OpenTelemetry](docs/METRICS\_TRACES\_LOGGING.md) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/bytCFQJ43DE)
- [ ] Caching
  - [X] Memcached [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/yKI-sy70PwA) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/link.svg" width="20" height="20" alt="Blog post">](https://mariocarrion.com/2021/01/30/tips-building-microservices-in-go-golang-caching-memcached.html)
  - [ ] Redis [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/wj6-w0DLKRw)
- [X] Persistent Storage
  - [X] Repository Pattern [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/Z89UU4vSayY) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/link.svg" width="20" height="20" alt="Blog post">](https://mariocarrion.com/2021/04/04/golang-microservices-repository-pattern.html)
  - [X] Database migrations [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/EavdaeUmn64) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/link.svg" width="20" height="20" alt="Blog post">](https://mariocarrion.com/2021/01/10/golang-tools-for-database-schema-migrations.html)
  - [ ] MySQL
  - [X] [PostgreSQL](docs/PERSISTENT\_STORAGE.md) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/link.svg" width="20" height="20" alt="Blog post">](https://mariocarrion.com/2021/03/02/tips-building-microservices-in-go-golang-databases-postgresql-conclusion.html)
      - [`jmoiron/sqlx`](https://github.com/jmoiron/sqlx), [`jackc/pgx`](https://github.com/jackc/pgx) and [`database/sql`](https://pkg.go.dev/database/sql) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/l8t6UKM1kro) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/link.svg" width="20" height="20" alt="Blog post">](https://mariocarrion.com/2021/02/05/tips-building-microservices-in-go-golang-databases-postgresql-part1.html)
      - [`go-gorm/gorm`](https://github.com/go-gorm/gorm) and [`volatiletech/sqlboiler`](https://github.com/volatiletech/sqlboiler) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/CT2v0Xas8Sc) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/link.svg" width="20" height="20" alt="Blog post">](https://mariocarrion.com/2021/02/10/tips-building-microservices-in-go-golang-databases-postgresql-gorm-orm.html)
      - [`Masterminds/squirrel`](https://github.com/Masterminds/squirrel) and [`kyleconroy/sqlc`](https://github.com/kyleconroy/sqlc) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/CT2v0Xas8Sc) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/link.svg" width="20" height="20" alt="Blog post">](https://mariocarrion.com/2021/02/23/tips-building-microservices-in-go-golang-databases-postgresql-sqlc-squirrel.html)
- [ ] REST APIs
  - [X] HTTP Handlers [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/CLdxwJCvTZE) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/link.svg" width="20" height="20" alt="Blog post">](https://mariocarrion.com/2021/04/18/golang-microservices-rest-api-http-handlers.html)
  - [X] Custom JSON Types [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/UmVYkEYm4hw) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/link.svg" width="20" height="20" alt="Blog post">](https://mariocarrion.com/2021/04/28/golang-microservices-rest-api-custom-json-type.html)
  - [ ] Versioning [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/4THy4iBQpFA)
  - [X] Error Handling [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/uQOfXL6IFmQ)
  - [X] [OpenAPI 3 and Swagger-UI](docs/OPENAPI3\_SWAGGER.md) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/HwtOAc0M08o) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/link.svg" width="20" height="20" alt="Blog post">](https://mariocarrion.com/2021/05/02/golang-microservices-rest-api-openapi3-swagger-ui.html)
  - [ ] Authorization
- [ ] Events and Messaging
  - [ ] [Apache Kafka](https://kafka.apache.org/) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/jr7OULxYm0A)
  - [ ] [RabbitMQ](https://www.rabbitmq.com/) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/L0yJxCKrkIY)
  - [ ] [Redis](https://redis.io/) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/uzuwZNiN4Y8)
- [ ] Testing
  - [X] Type-safe mocks with [`maxbrunsfeld/counterfeiter`](https://github.com/maxbrunsfeld/counterfeiter) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/ENqwq64TsDk) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/link.svg" width="20" height="20" alt="Blog post">](https://mariocarrion.com/2019/06/24/golang-tools-counterfeiter.html)
  - [X] Equality with [`google/go-cmp`](https://github.com/google/go-cmp) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/ae15DzSwNnU) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/link.svg" width="20" height="20" alt="Blog post">](https://mariocarrion.com/2021/01/22/go-package-equality-google-go-cmp.html)
  - [X] Integration tests for Datastores with [`ory/dockertest`](https://github.com/ory/dockertest) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/a-CCceqerhg) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/link.svg" width="20" height="20" alt="Blog post">](https://mariocarrion.com/2021/03/14/golang-package-testing-datastores-ory-dockertest.html)
  - [X] REST APIs [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/lMrWO7OUMdY) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/link.svg" width="20" height="20" alt="Blog post">](https://mariocarrion.com/2021/04/25/golang-microservices-rest-api-testing.html)
- [ ] Containerization using Docker [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/u_ayzie9pAQ)
- [ ] Graceful Shutdown [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/VXxe7-b5euo)
- [ ] Search Engine using [ElasticSearch](https://www.elastic.co/elasticsearch/) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/ZrdbQRYst5E)
- [ ] Whatever else I forgot to include

## More ideas

* [2016: Peter Bourgon's: Repository structure](https://peter.bourgon.org/go-best-practices-2016/#repository-structure)
* [2016: Ben Johnson's: Standard Package Layout](https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1)
* [2017: William Kennedy's: Design Philosophy On Packaging](https://www.ardanlabs.com/blog/2017/02/design-philosophy-on-packaging.html)
* [2017: Jaana Dogan's: Style guideline for Go packages](https://rakyll.org/style-packages/)
* [2018: Kat Zien - How Do You Structure Your Go Apps](https://www.youtube.com/watch?v=oL6JBUk6tj0)

## Docker Containers

Please notice in order to run this project locally you need to run a few programs in advance, if you use Docker please refer to the concrete instructions in [`docs/`](docs/) for more details.

There's also a [docker-compose.yml](docker-compose.yml), covered in [Building Microservices In Go: Containerization with Docker](https://youtu.be/u_ayzie9pAQ), however like I mentioned in the video you have to execute `docker-compose` in multiple steps.

Notice that because of the way RabbitMQ and Kafka are being used they are sort of competing with each other, so at the moment we either have to enable Kafka and disable RabbitMQ or the other way around in both the code and the `docker-compose.yml` file, in either case there are Dockerfiles and services defined that cover building and running them.

* Run `docker-compose up`, here both _rest-server_ and _elasticsearch-indexer_ services will fail because the `postgres`, `rabbitmq`, `elasticsearch` and `kafka` services take too long to start.
    * If you're planning to use RabbitMQ, run `docker-compose up rest-server elasticsearch-indexer-rabbitmq`.
    * If you're planning to use Kafka, run `docker-compose up rest-server elasticsearch-indexer-kafka`.
    * If you're planning to use Redis, run `docker-compose up rest-server elasticsearch-indexer-redis`.
* For building the service images you can use:
    * `rest-server` image: `docker-compose build rest-server`.
    * `elasticsearch-indexer-rabbitmq` image: `docker-compose build elasticsearch-indexer-rabbitmq`.
    * `elasticsearch-indexer-kafka` image: `docker-compose build elasticsearch-indexer-kafka`.
    * `elasticsearch-indexer-redis` image: `docker-compose build elasticsearch-indexer-redis`.
* Run `docker-compose run rest-server migrate -path /api/migrations/ -database postgres://user:password@postgres:5432/dbname?sslmode=disable up` to finally have everything working correctly.

## Diagrams

To start a local HTTP server that serves a graphical editor:

```
mdl serve github.com/MarioCarrion/todo-api/internal/doc -dir docs/diagrams/
```

To generate JSON artifact for uploading to [structurizr](https://structurizr.com/):

```
stz gen github.com/MarioCarrion/todo-api/internal/doc
```
