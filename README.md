# "ToDo API" Microservice Example

## Introduction

Welcome! 👋

This is an educational repository that includes a microservice written in Go. It is used as the principal example of my video series: [Building Microservices in Go/Golang](https://www.youtube.com/playlist?list=PL7yAAGMOat_Fn8sAXIk0WyBfK_sT1pohu).

It's a collection of patterns and guidelines I've successfully used to deliver enterprise microservices when using Go.

The whole purpose of this project is to give you an idea about structuring your Go project with 3 principal goals:

1. It is _enterprise_, meant to last for years,
1. It allows a team to collaborate efficiently with little friction, and
1. It is as idiomatic as possible.

Join the fun at https://youtube.com/MarioCarrion

## Domain Driven Design

This project uses a lot of the ideas introduced by Eric Evans in his book [Domain Driven Design](https://www.domainlanguage.com/), I do encourage reading that book but before I think reading [Domain-Driven Design Distilled](https://smile.amazon.com/Domain-Driven-Design-Distilled-Vaughn-Vernon/dp/0134434420/) makes more sense, also there's a free to download [DDD Reference](https://www.domainlanguage.com/ddd/reference/) available as well.

On YouTube I created [a playlist](https://www.youtube.com/playlist?list=PL7yAAGMOat_GJqfTdM9PBdTRSH7jXs6mI) that includes some of my favorite talks and webinars, feel free to explore that as well.

## Project Structure

Talking specifically about microservices **only**, the structure I like to recommended is the following, everything using `<` and `>` depends on the domain being implemented and the bounded context being defined.

- [ ] `build/`: defines the code used for creating infrastructure
  - [ ] `<cloud-providers>/`: define concrete cloud provider 
  - [ ] `<executableN>/`: contains a Dockerfile used for building the binary
- [ ] `cmd/`
  - [ ] `<primary-server>/`: uses primary database.
  - [ ] `<replica-server>/`: uses readonly databases.
  - [ ] `<binaryN>/`
- [X] `db/`
  - [X] `migrations/`: contains database migrations
  - [ ] `seeds/`: contains file meant to populate basic database values
- [ ] `internal/`: defines the _core domain_
  - [ ] `<datastoreN>/`: a concrete _repository_ used by the domain, for example `postgresql`
  - [ ] `http/`: defines HTTP Handlers
  - [ ] `service/`: orchestrates use cases and manages transactions.
- [X] `pkg/` public API meant to imported by other Go package

## Tools

```
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.14.1
go install github.com/kyleconroy/sqlc/cmd/sqlc@v1.6.0
go install github.com/maxbrunsfeld/counterfeiter/v6@v6.3.0
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.5.1
```

## Features

Icons meaning:

* <img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video"> means a link to Youtube video.
* <img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/link.svg" width="20" height="20" alt="Blog post"> means a link to Blog post.

In no particular order:

- [X] Database migrations [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/EavdaeUmn64) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/link.svg" width="20" height="20" alt="Blog post">](https://mariocarrion.com/2021/01/10/golang-tools-for-database-schema-migrations.html)
- [X] Repositories [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/Z89UU4vSayY)
- [X] Project Layout [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/LUvid5TJ81Y)
- [X] Dependency Injection [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/Z89UU4vSayY)
- [X] Secure Configuration [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/7UmJR0dOkjM)
  - [X] Using AWS SSM [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/link.svg" width="20" height="20" alt="Blog post">](https://mariocarrion.com/2019/11/20/golang-aws-ssm-env-vars-package.html)
- [X] [OpenAPI/Swagger](docs/SWAGGER.md) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/HwtOAc0M08o)
- [ ] Infrastructure as code
- [ ] [Metrics, Traces and Logging](docs/METRICS\_TRACES\_LOGGING.md)
- [X] Error Handling [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/uQOfXL6IFmQ)
- [ ] Caching
  - [X] Memcached [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/yKI-sy70PwA) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/link.svg" width="20" height="20" alt="Blog post">](https://mariocarrion.com/2021/01/30/tips-building-microservices-in-go-golang-caching-memcached.html)
  - [ ] Redis [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/wj6-w0DLKRw)
- [X] Persistent Storage (using PostgreSQL) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/link.svg" width="20" height="20" alt="Blog post">](https://mariocarrion.com/2021/03/02/tips-building-microservices-in-go-golang-databases-postgresql-conclusion.html)
  - [`jmoiron/sqlx`](https://github.com/jmoiron/sqlx), [`jackc/pgx`](https://github.com/jackc/pgx) and [`database/sql`](https://pkg.go.dev/database/sql) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/l8t6UKM1kro) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/link.svg" width="20" height="20" alt="Blog post">](https://mariocarrion.com/2021/02/05/tips-building-microservices-in-go-golang-databases-postgresql-part1.html)
  - [`go-gorm/gorm`](https://github.com/go-gorm/gorm) and [`volatiletech/sqlboiler`](https://github.com/volatiletech/sqlboiler) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/CT2v0Xas8Sc) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/link.svg" width="20" height="20" alt="Blog post">](https://mariocarrion.com/2021/02/10/tips-building-microservices-in-go-golang-databases-postgresql-gorm-orm.html)
  - [`Masterminds/squirrel`](https://github.com/Masterminds/squirrel) and [`kyleconroy/sqlc`](https://github.com/kyleconroy/sqlc) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/CT2v0Xas8Sc) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/link.svg" width="20" height="20" alt="Blog post">](https://mariocarrion.com/2021/02/23/tips-building-microservices-in-go-golang-databases-postgresql-sqlc-squirrel.html)
- [ ] Authorization
- [ ] REST APIs
  - [ ] Custom JSON Types [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/UmVYkEYm4hw)
  - [ ] Versioning [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/4THy4iBQpFA)
- [ ] Events
- [ ] Testing
  - [X] Type-safe mocks with [`maxbrunsfeld/counterfeiter`](https://github.com/maxbrunsfeld/counterfeiter) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/ENqwq64TsDk) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/link.svg" width="20" height="20" alt="Blog post">](https://mariocarrion.com/2019/06/24/golang-tools-counterfeiter.html)
  - [X] Equality with [`google/go-cmp`](https://github.com/google/go-cmp) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/ae15DzSwNnU) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/link.svg" width="20" height="20" alt="Blog post">](https://mariocarrion.com/2021/01/22/go-package-equality-google-go-cmp.html)
  - [ ] REST APIs [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/lMrWO7OUMdY)
  - [ ] Integration tests for Datastores with [`ory/dockertest`](https://github.com/ory/dockertest) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/youtube.svg" width="20" height="20" alt="YouTube video">](https://youtu.be/a-CCceqerhg) [<img src="https://github.com/MarioCarrion/MarioCarrion/blob/main/link.svg" width="20" height="20" alt="Blog post">](https://mariocarrion.com/2021/03/14/golang-package-testing-datastores-ory-dockertest.html)
- [ ] Whatever else I forgot to include

## More ideas

* [2016: Peter Bourgon's: Repository structure](https://peter.bourgon.org/go-best-practices-2016/#repository-structure)
* [2016: Ben Johnson's: Standard Package Layout](https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1)
* [2017: William Kennedy's: Design Philosophy On Packaging](https://www.ardanlabs.com/blog/2017/02/design-philosophy-on-packaging.html)
* [2017: Jaana Dogan's: Style guideline for Go packages](https://rakyll.org/style-packages/)
* [2018: Kat Zien - How Do You Structure Your Go Apps](https://www.youtube.com/watch?v=oL6JBUk6tj0)

## Docker Containers

### PostgreSQL

Used as repository for persisting data.

```
docker run \
  -d \
  -e POSTGRES_HOST_AUTH_METHOD=trust \
  -e POSTGRES_USER=user \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=dbname \
  -p 5432:5432 \
  postgres:12.5-alpine
```

### Vault

Used as repository for retrieving secrets for configuration values.

```
docker run \
  -d \
  --cap-add=IPC_LOCK \
  -e 'VAULT_DEV_ROOT_TOKEN_ID=myroot' \
  -e 'VAULT_DEV_LISTEN_ADDRESS=0.0.0.0:8300' \
  -p 8300:8300 \
  vault:1.6.2
```
